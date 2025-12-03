package mirg

import (
	"context"
	"fmt"
	"github.com/long250038728/cmd/mirg/model"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/excel"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/sliceconv"
	"math/rand"
	"time"
)

type excelModel struct {
	Name         string `json:"name"`
	Telephone    string `json:"telephone"`
	MerchantShop string `json:"merchant_shop"`
	StaffTel     string `json:"staff_tel"`
}

var excelHeader = []excel.Header{
	{Key: "name", Name: "客户姓名", Type: "string"},
	{Key: "telephone", Name: "客户手机号", Type: "string"},
	{Key: "merchant_shop", Name: "所属门店", Type: "string"},
	{Key: "staff_tel", Name: "归属员工手机号", Type: "string"},
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
	Id     int32  `json:"id"`
	Name   string `json:"name"`
	ShopSn string `json:"shop_sn"`
}

type AdminUser struct {
	Id     int32  `json:"id"`
	Mobile string `json:"mobile"`
}

var merchantConfigPath = "./config/online/db.yaml"
var merchantStaffUrl = "https://moss.zhubaoe.cn/api.php?s=/customer/updateJoinInfoList"

//====================================================================

func MerchantAction(isAddLog bool, isUpdate bool) {
	var merchantId int32 = 0
	BrandId := 0
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
		if "老凤祥衡水信发店" == shopName {
			continue
		}

		if !ok {
			fmt.Println(fmt.Errorf("%s", shopName))
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
		fmt.Println(len(noDbTelephone))
	}
}

func MerchantStaffAction(isUpdate bool) {
	var merchantId int32 = 0
	BrandId := 0
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
		adminHash[val.Mobile] = val.Id
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
		MerchantStaffToCustomer[c.StaffTel] = append(MerchantStaffToCustomer[c.StaffTel], c.Telephone)
	}

	//会员手机不在数据库中
	noDbTelephone := make([]string, 0, len(data))

	//====================================================================================

	for staffTel, val := range MerchantStaffToCustomer {
		adminId, ok := adminHash[staffTel]
		if !ok {
			fmt.Println(fmt.Errorf("%s", staffTel))
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
					fmt.Println("================", err.Error(), adminId, staffTel, index)
				}
			}
		}
	}

	if len(noDbTelephone) > 0 {
		//b, err := json.Marshal(noDbTelephone)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//fmt.Println(string(b))
		fmt.Println(len(noDbTelephone))
	}
}

func MerchantCustomerAddAction(isUpdate bool) {
	var merchantId int32 = 0
	var BrandId int32 = 1008
	Path := "/Users/linlong/Desktop/b.xlsx"
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
		shopHash[val.ShopSn] = val.Id
		shopNameHash[val.Id] = val.Name
	}

	//====================================================================================

	// 获取表格信息
	data, err := excelExcel(Path, sheet)
	if err != nil {
		panic(err)
	}

	moreTel := make(map[string]string, len(data))
	for _, c := range data {
		if _, ok := moreTel[c.Telephone]; ok {
			panic("会员重复")
		} else {
			moreTel[c.Telephone] = c.Telephone
		}
	}

	//生成数据 key为门店名称  value 为手机号
	MerchantShopToCustomer := make(map[string][]*excelModel)
	for _, c := range data {
		MerchantShopToCustomer[c.MerchantShop] = append(MerchantShopToCustomer[c.MerchantShop], c)
	}

	//====================================================================================

	for shopName, val := range MerchantShopToCustomer {
		shopName = merchantSellerIdChange(shopName)

		merchantShopId, ok := shopHash[shopName]
		if "老凤祥衡水信发店" == shopName {
			continue
		}

		if !ok {
			fmt.Println(fmt.Errorf("%s", shopName))
		}

		// 把数据切成1000为一次处理
		for index, chuck := range sliceconv.Chunk(val, 1000) {

			tels := sliceconv.Extract(chuck, func(d *excelModel) string {
				return d.Telephone
			})

			//查询数据库
			chuckCustomers := make([]*Customer, 0, 1000)
			if err := db.Where("merchant_id = ? AND brand_id = ? AND status = 1  and telephone in (?)", merchantId, BrandId, tels).Find(&chuckCustomers).Error; err != nil {
				panic(err)
			}

			if isUpdate {
				// 获取所有会员id && 更新会员
				// 数据库存在的手机号
				hasTel := sliceconv.Map(chuckCustomers, func(item *Customer) (key string, value interface{}) {
					return item.Telephone, 1
				})

				createList := make([]*Customer, 0, 1000)
				for _, c := range chuck {
					if _, ok := hasTel[c.Telephone]; !ok {
						createList = append(createList, &Customer{
							MerchantId:      merchantId,
							MerchantShopId:  merchantShopId,
							BrandId:         BrandId,
							Telephone:       c.Telephone,
							SuffixTelephone: c.Telephone[len(c.Telephone)-4:],
							Name:            c.Name,
							Status:          1,
							CreateTime:      int32(time.Now().Unix()),
							UpdateTime:      int32(time.Now().Unix()),
							BirthdayDt:      "0000-00-00",
							MarryDate:       "0000-00-00",
							LastBuyTime:     "0000-00-00 00:00:00",
							AddDatetime:     "0000-00-00 00:00:00",
							OriginPlatform:  4,
							OriginChannel:   4,
							Level:           1,
						})
					}
				}

				res := db.Table("zby_customer").CreateInBatches(createList, 500)
				if res.Error != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(shopName, index, len(chuckCustomers), len(chuck), res.RowsAffected)
			}

		}
	}
}

func merchantSellerIdChange(SellerID string) string {
	hash := map[string]string{
		"ZY002": "ZY067",
		"ZY003": "ZY085",
		"ZY004": "ZY057",
		"ZY006": "ZY022",
		"ZY011": "ZY005",
		"ZY013": "ZY010",
		"ZY014": "ZY020",
		"ZY016": "ZY099",
		"ZY017": "ZY073",
		"ZY019": "ZY015",
		"ZY024": "ZY092",
		"ZY040": "ZY083",
		"ZY041": "ZY022",
		"ZY045": "ZY022",
		"ZY046": "ZY022",
		"ZY050": "ZY091",
		"ZY053": "ZY036",
		"ZY054": "ZY048",
		"ZY055": "ZY083",
		"ZY058": "ZY051",
		"ZY060": "ZY034",
		"ZY072": "ZY039",
		"ZY077": "ZY102",
		"ZY079": "ZY001",
		"ZY081": "ZY083",
		"ZY082": "ZY067",
		"ZY086": "ZY110",
		"ZY089": "ZY057",
		"ZY090": "ZY073",
		"ZY093": "ZY109",
		"ZY098": "ZY049",
		"ZY106": "ZY048",
	}

	if newSellerID, ok := hash[SellerID]; ok {
		return newSellerID
	}
	return SellerID
}

//func TestJoinHttp() {
//	adminId := 1699
//	chuckCustomers := []*Customer{
//		{Id: 11276507, MerchantId: 168, BrandId: 543},
//	}
//
//	// 获取所有会员id && 更新会员
//	list := sliceconv.Change(chuckCustomers, func(d *Customer) map[string]any {
//		return map[string]any{
//			"merchant_id":   d.MerchantId,
//			"brand_id":      d.BrandId,
//			"customer_id":   d.Id,
//			"admin_user_id": adminId,
//			"source":        3,
//			"user_id":       0,
//		}
//	})
//
//	data := map[string]any{
//		"type": 1,
//		"list": list,
//		"s":    "123456",
//	}
//	httpClient := http.NewClient(http.SetTimeout(time.Second * 5))
//	b, _, err := httpClient.Post(context.Background(), merchantStaffUrl, data)
//
//	fmt.Println(string(b))
//	if err != nil {
//		fmt.Println("================", err.Error(), adminId)
//	}
//}

type CustomerE struct {
	Id       int32  `json:"id"`
	ShopName string `json:"shop_name"`
}

type AdminE struct {
	Tel      string `json:"tel"`
	ShopName string `json:"shop_name"`
}

var customerHeader = []excel.Header{
	{Key: "id", Name: "会员", Type: "int"},
	{Key: "shop_name", Name: "门店", Type: "string"},
}
var adminHeader = []excel.Header{
	{Key: "tel", Name: "手机号", Type: "string"},
	{Key: "shop_name", Name: "所属门店", Type: "string"},
}

func TestAAAAAA() {
	merchantId := 1843
	//BrandId := 0

	// 获取会员信息
	var ormConfig orm.Config
	configurator.NewYaml().MustLoad(merchantConfigPath, &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}

	ShopNameToIds := make(map[string][]int32)
	{
		TelToId := make(map[string]int32)
		{
			adminUsers := make([]*AdminUser, 0, 100)
			if err := db.Where("merchant_id = ?", merchantId).Find(&adminUsers).Error; err != nil {
				panic(err)
			}
			for _, a := range adminUsers {
				TelToId[a.Mobile] = a.Id
			}
		}

		adminList, err := aExcel("/Users/linlong/Desktop/admin.xlsx", "Sheet1")
		if err != nil {
			return
		}
		for _, a := range adminList {
			if _, ok := ShopNameToIds[a.ShopName]; !ok {
				ShopNameToIds[a.ShopName] = make([]int32, 0)
			}
			ShopNameToIds[a.ShopName] = append(ShopNameToIds[a.ShopName], TelToId[a.Tel])
		}
	}

	//=====

	customerList, err := cExcel("/Users/linlong/Desktop/customer.xlsx", "Sheet1")
	if err != nil {
		return
	}

	hash := make(map[string][]*CustomerE)
	for _, c := range customerList {
		if _, ok := hash[c.ShopName]; !ok {
			hash[c.ShopName] = make([]*CustomerE, 0)
		}
		hash[c.ShopName] = append(hash[c.ShopName], c)
	}

	for shopName, data := range hash {
		ids := ShopNameToIds[shopName]
		list := sliceconv.Change(data, func(d *CustomerE) map[string]any {
			i := rand.Int() % len(ids)
			return map[string]any{
				"merchant_id":   1843,
				"brand_id":      1008,
				"customer_id":   d.Id,
				"admin_user_id": ids[i],
				"source":        3,
			}
		})

		for _, chuck := range sliceconv.Chunk(list, 200) {
			httpClient := http.NewClient(http.SetTimeout(time.Second * 5))
			b, _, err := httpClient.Post(context.Background(), merchantStaffUrl, map[string]any{
				"type": 1,
				"list": chuck,
			})
			fmt.Println(string(b))
			if err != nil {
				fmt.Println("================", err.Error())
			}
		}
	}
}

func cExcel(path, sheet string) ([]*CustomerE, error) {
	var data []*CustomerE
	r := excel.NewRead(path)
	defer r.Close()
	err := r.Read(sheet, customerHeader, &data)

	if err != nil {
		return nil, err
	}
	return data, nil
}

func aExcel(path, sheet string) ([]*AdminE, error) {
	var data []*AdminE
	r := excel.NewRead(path)
	defer r.Close()
	err := r.Read(sheet, adminHeader, &data)

	if err != nil {
		return nil, err
	}
	return data, nil
}
