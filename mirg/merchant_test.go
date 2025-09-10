package mirg

import (
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/excel"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/sliceconv"
	"testing"
	"time"
)

func TestBB(t *testing.T) {
	MerchantAction(1)
}

type EModel struct {
	Telephone    string `json:"telephone"`
	MerchantShop string `json:"merchant_shop"`
}

type CustomerLog struct {
	MerchantId     int32  `json:"merchant_id"`
	MerchantShopId int32  `json:"merchant_shop_id"`
	CustomerId     int32  `json:"customer_id"`
	CustomerName   string `json:"customer_name"`
	Type           int32  `json:"type"`
	Comment        string `json:"comment"`
	CreateTime     string `json:"create_time"`
}

var EHeader = []excel.Header{
	{Key: "telephone", Name: "手机号", Type: "string"},
	{Key: "merchant_shop", Name: "转移门店", Type: "string"},
}

func EExcel(path, sheet string) ([]*EModel, error) {
	var data []*EModel
	r := excel.NewRead(path)
	defer r.Close()
	err := r.Read(sheet, EHeader, &data)

	if err != nil {
		return nil, err
	}
	return data, nil
}

type MerchantShop struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

//func roundToNDecimal(f float64, n int) float64 {
//	pow := math.Pow(10, float64(n))
//	return math.Round(f*pow) / pow
//}

// ====================================================================

//type Customer struct {
//	Id             int32  `json:"id"`
//	Name           string `json:"name"`
//	MerchantId     int32  `json:"merchant_id"`
//	BrandId        int32  `json:"brand_id"`
//	MerchantShopId int32  `json:"merchant_shop_id"`
//	Telephone      string `json:"telephone"`
//}
//
//type CustomerBpLog struct {
//	MerchantId     int32 `json:"merchant_id"`
//	BrandId        int32 `json:"brand_id"`
//	MerchantShopId int32 `json:"merchant_shop_id"`
//
//	CustomerId   int32  `json:"customer_id"`
//	CustomerName string `json:"customer_name"`
//
//	PointTotal   float64 `json:"point_total"`
//	Point        float64 `json:"point"`
//	Type         int32   `json:"type"`
//	Comment      string  `json:"comment"`
//	CreateTime   int32   `json:"create_time"`
//	ActivityName string  `json:"activity_name"`
//	ActivityId   int32   `json:"activity_id"`
//	AdminUserId  int32   `json:"admin_user_id"`
//}

//====================================================================

func MerchantAction(accessToken int) {
	var merchantId int32 = 3
	BrandId := 1
	Path := "/Users/linlong/Desktop/a.xlsx"
	sheet := "Sheet2"

	// 获取表格信息
	data, err := EExcel(Path, sheet)
	if err != nil {
		panic(err)
	}

	telHash := make(map[string][]string)
	for _, c := range data {
		telHash[c.MerchantShop] = append(telHash[c.MerchantShop], c.Telephone)
	}

	// 获取会员信息
	var ormConfig orm.Config
	configurator.NewYaml().MustLoad("./config/emperor/db.yaml", &ormConfig)
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

	for shop, val := range telHash {
		merchantShopId, ok := shopHash[shop]
		if !ok {
			panic("cccc")
		}

		for _, chuck := range sliceconv.Chunk(val, 1000) {
			chuckCustomers := make([]*Customer, 0, 1000)
			chuckCustomerLogs := make([]*CustomerLog, 0, 100)
			if err := db.Where("merchant_id = ?", merchantId).
				Where("brand_id = ?", BrandId).
				Where("status = ?", 1).
				Where("telephone in (?)", chuck).
				Find(&chuckCustomers).Error; err != nil {
				panic(err)
			}

			ids := sliceconv.Extract(datchuckCustomersa, func(d *Customer) int32 {
				return d.Id
			})

			row := db.Table("zby_customer").Where("id in (?)", ids).Updates(map[string]any{
				"merchant_shop_id": id,
				"update_time":      time.Now().Unix(),
			}).RowsAffected
			fmt.Println(row)

			for _, c := range chuckCustomers {
				chuckCustomerLogs = append(chuckCustomerLogs, &CustomerLog{
					MerchantId:     merchantId,
					MerchantShopId: merchantShopId,
					CustomerId:     c.Id,
					CustomerName:   c.Name,
					Type:           3,
					Comment:        "从" + shopNameHash[c.MerchantShopId] + "迁移到" + shop,
					CreateTime:     "2025-09-15 12:00:00",
				})
			}

			err = db.Create(chuckCustomerLogs).Error
			if err != nil {
				continue
			}
		}
	}
}
