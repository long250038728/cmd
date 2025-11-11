package mirg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/long250038728/cmd/mirg/model"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/excel"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/sliceconv"
	"io/ioutil"
	http2 "net/http"
	"time"
)

type excelModel struct {
	Telephone    string `json:"telephone"`
	MerchantShop string `json:"merchant_shop"`
	StaffName    string `json:"staff_name"`
}

var excelHeader = []excel.Header{
	{Key: "telephone", Name: "客户手机号", Type: "string"},
	{Key: "merchant_shop", Name: "所属门店", Type: "string"},
	{Key: "staff_name", Name: "归属员工姓名", Type: "string"},
}

func excelExcel(path, sheet string) ([]*excelModel, error) {
	var data []*excelModel
	r := excel.NewRead(path)
	defer r.Close()
	err := r.Read(sheet, excelHeader, &data)

	if err != nil {
		return nil, err
	}
	return data, nil
}

type MerchantShop struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type AdminUser struct {
	Id       int32  `json:"id"`
	Username string `json:"username"`
}

var merchantConfigPath = "./config/test/db.yaml"
var merchantStaffUrl = "http://192.168.0.55/api.php?s=/customer/updateJoinInfoList"

//====================================================================

func MerchantAction(isAddLog bool, isUpdate bool) {
	var merchantId int32 = 3
	BrandId := 1
	Path := "/Users/linlong/Desktop/a.xlsx"
	sheet := "Sheet1"

	// 获取会员信息
	var ormConfig orm.Config
	configurator.NewYaml().MustLoad(merchantConfigPath, &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}

	//门店
	shops := make([]*MerchantShop, 0, 100)
	if err := db.Where("merchant_id = ?", merchantId).
		Where("brand_id = ?", BrandId).
		Find(&shops).Error; err != nil {
		panic(err)
	}
	shopHash := make(map[string]int32)
	shopNameHash := make(map[int32]string)
	for _, val := range shops {
		shopHash[val.Name] = val.Id
		shopNameHash[val.Id] = val.Name
	}

	//====================================================================================

	// 获取表格信息
	data, err := excelExcel(Path, sheet)
	if err != nil {
		panic(err)
	}

	//生成数据 key为门店名称  value 为手机号
	MerchantShopToCustomer := make(map[string][]string)
	for _, c := range data {
		MerchantShopToCustomer[c.MerchantShop] = append(MerchantShopToCustomer[c.MerchantShop], c.Telephone)
	}

	//会员手机不在数据库中
	noDbTelephone := make([]string, 0, len(data))

	//====================================================================================

	for shopName, val := range MerchantShopToCustomer {
		merchantShopId, ok := shopHash[shopName]
		if !ok {
			panic("cccc")
		}

		// 把数据切成1000为一次处理
		for index, chuck := range sliceconv.Chunk(val, 1000) {

			//查询数据库
			chuckCustomers := make([]*Customer, 0, 1000)
			if err := db.Where("merchant_id = ? AND brand_id = ? AND status = 1  and telephone in (?)", merchantId, BrandId, chuck).Find(&chuckCustomers).Error; err != nil {
				panic(err)
			}

			//计算没有存在的会员
			{
				//数据库存在的手机号
				hasTel := sliceconv.Map(chuckCustomers, func(item *Customer) (key string, value interface{}) {
					return item.Telephone, 1
				})

				for _, tel := range chuck {
					if _, ok := hasTel[tel]; !ok {
						noDbTelephone = append(noDbTelephone, tel)
					}
				}
			}

			if isUpdate {
				// 获取所有会员id && 更新会员
				ids := sliceconv.Extract(chuckCustomers, func(d *Customer) int32 {
					return d.Id
				})
				res := db.Table("zby_customer").Where("id in (?)", ids).Updates(map[string]any{
					"merchant_shop_id": merchantShopId,
					"update_time":      time.Now().Unix(),
				})
				if res.Error != nil {
					fmt.Println(err.Error())
				}
				//fmt.Println(merchantShopId, shopName, index, len(chuckCustomers), len(chuck))
				fmt.Println(shopName, index, len(chuckCustomers), len(chuck), res.RowsAffected)
			}

			// 是否添加log
			if isAddLog {
				chuckCustomerLogs := make([]*model.CustomerLog, 0, 100)
				for _, c := range chuckCustomers {
					oldShopName, ok := shopNameHash[c.MerchantShopId]
					if !ok {
						oldShopName = "-"
					}
					chuckCustomerLogs = append(chuckCustomerLogs, &model.CustomerLog{
						MerchantId:     merchantId,
						MerchantShopId: merchantShopId,
						CustomerId:     c.Id,
						CustomerName:   c.Name,
						Type:           3,
						Comment:        "从" + oldShopName + "迁移到" + shopName,
						CreateTime:     "2025-10-15 10:00:00",
					})
				}
				db.Create(chuckCustomerLogs)
			}
		}
	}

	if len(noDbTelephone) > 0 {
		b, err := json.Marshal(noDbTelephone)
		if err != nil {
			fmt.Println(string(b))
		}
	}
}

func MerchantStaffAction(isAddLog bool, isUpdate bool) {
	var merchantId int32 = 3
	BrandId := 1
	Path := "/Users/linlong/Desktop/a.xlsx"
	sheet := "Sheet1"

	// 获取会员信息
	var ormConfig orm.Config
	configurator.NewYaml().MustLoad(merchantConfigPath, &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}

	//门店
	adminUsers := make([]*AdminUser, 0, 100)
	if err := db.Where("merchant_id = ?", merchantId).
		Find(&adminUsers).Error; err != nil {
		panic(err)
	}
	adminHash := make(map[string]int32)
	for _, val := range adminUsers {
		adminHash[val.Username] = val.Id
	}

	//====================================================================================

	// 获取表格信息
	data, err := excelExcel(Path, sheet)
	if err != nil {
		panic(err)
	}

	//生成数据 key为员工归属  value 为手机号
	MerchantStaffToCustomer := make(map[string][]string)
	for _, c := range data {
		MerchantStaffToCustomer[c.StaffName] = append(MerchantStaffToCustomer[c.StaffName], c.Telephone)
	}

	//会员手机不在数据库中
	noDbTelephone := make([]string, 0, len(data))

	//====================================================================================

	for staffName, val := range MerchantStaffToCustomer {
		adminId, ok := adminHash[staffName]
		if !ok {
			panic("cccc")
		}

		// 把数据切成1000为一次处理
		for index, chuck := range sliceconv.Chunk(val, 200) {

			//查询数据库
			chuckCustomers := make([]*Customer, 0, 200)
			if err := db.Where("merchant_id = ? AND brand_id = ? AND status = 1  and telephone in (?)", merchantId, BrandId, chuck).Find(&chuckCustomers).Error; err != nil {
				panic(err)
			}

			//计算没有存在的会员
			{
				//数据库存在的手机号
				hasTel := sliceconv.Map(chuckCustomers, func(item *Customer) (key string, value interface{}) {
					return item.Telephone, 1
				})

				for _, tel := range chuck {
					if _, ok := hasTel[tel]; !ok {
						noDbTelephone = append(noDbTelephone, tel)
					}
				}
			}

			if isUpdate {
				// 获取所有会员id && 更新会员
				list := sliceconv.Change(chuckCustomers, func(d *Customer) map[string]any {
					return map[string]any{
						"merchant_id":   d.MerchantId,
						"brand_id":      d.BrandId,
						"customer_id":   d.Id,
						"admin_user_id": adminId,
						"source":        3,
					}
				})

				httpClient := http.NewClient(http.SetTimeout(time.Second * 5))
				b, _, err := httpClient.Post(context.Background(), merchantStaffUrl, map[string]any{
					"type": 1,
					"list": list,
				})

				fmt.Println(string(b))
				if err != nil {
					fmt.Println("================", err.Error(), adminId, staffName, index)
				}
			}
		}
	}

	if len(noDbTelephone) > 0 {
		b, err := json.Marshal(noDbTelephone)
		if err != nil {
			fmt.Println(string(b))
		}
	}
}

func TestJoinHttp() {
	adminId := 1699
	chuckCustomers := []*Customer{
		{Id: 11276507, MerchantId: 168, BrandId: 543},
	}

	// 获取所有会员id && 更新会员
	list := sliceconv.Change(chuckCustomers, func(d *Customer) map[string]any {
		return map[string]any{
			"merchant_id":   d.MerchantId,
			"brand_id":      d.BrandId,
			"customer_id":   d.Id,
			"admin_user_id": adminId,
			"source":        3,
			"user_id":       0,
		}
	})

	data := map[string]any{
		"type": 1,
		"list": list,
		"s":    "123456",
	}
	b, _ := json.Marshal(data)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10000)
	defer cancel()

	request, err := http2.NewRequestWithContext(ctx, "POST", merchantStaffUrl, bytes.NewReader(b))
	request.Header.Set("Content-Type", "application/json")
	res, err := http2.DefaultClient.Do(request)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	if err != nil {
		fmt.Println("================", err.Error(), adminId)
	}
}
