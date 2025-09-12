package mirg

import (
	"fmt"
	"github.com/long250038728/cmd/mirg/model"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/excel"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/sliceconv"
	"time"
)

type excelModel struct {
	Telephone    string `json:"telephone"`
	MerchantShop string `json:"merchant_shop"`
}

var excelHeader = []excel.Header{
	{Key: "telephone", Name: "手机号", Type: "string"},
	{Key: "merchant_shop", Name: "转移门店", Type: "string"},
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

//====================================================================

func MerchantAction() {
	var merchantId int32 = 3
	BrandId := 1
	Path := "/Users/linlong/Desktop/a.xlsx"
	sheet := "Sheet2"

	// 获取表格信息
	data, err := excelExcel(Path, sheet)
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

		for index, chuck := range sliceconv.Chunk(val, 1000) {
			chuckCustomers := make([]*Customer, 0, 1000)
			chuckCustomerLogs := make([]*model.CustomerLog, 0, 100)

			query := orm.NewBoolQuery().Must(
				orm.Eq("merchant_id", merchantId),
				orm.Eq("brand_id", BrandId),
				orm.Eq("status", 1),
				orm.In("telephone", chuck),
			)

			sql, args := query.Do()
			if err := db.Where(sql, args...).Find(&chuckCustomers).Error; err != nil {
				panic(err)
			}

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

			fmt.Println(shop, index, len(chuckCustomers), len(chuck), res.RowsAffected)

			for _, c := range chuckCustomers {
				chuckCustomerLogs = append(chuckCustomerLogs, &model.CustomerLog{
					MerchantId:     merchantId,
					MerchantShopId: merchantShopId,
					CustomerId:     c.Id,
					CustomerName:   c.Name,
					Type:           3,
					Comment:        "从" + shopNameHash[c.MerchantShopId] + "迁移到" + shop,
					CreateTime:     "2025-09-15 10:00:00",
				})
			}

			err = db.Create(chuckCustomerLogs).Error
			if err != nil {
				continue
			}
		}
	}
}
