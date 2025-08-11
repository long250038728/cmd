package mirg

import (
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/persistence/orm"
)

type GoldSalePriceConfig struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`

	MerchantId     int32 `json:"merchant_id"`
	MerchantShopId int32 `json:"merchant_shop_id"`
	BrandId        int32 `json:"brand_id"`
}

type GoodsStock struct {
	Id int32 `json:"id"`
}

func GetGoldSalePriceConfig() {
	deleteType := 1

	var ormConfig orm.Config
	configurator.NewYaml().MustLoad("./config/online/db.yaml", &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}

	var list []*GoldSalePriceConfig
	if err := db.Raw("SELECT zby_gold_sale_price_config.id,zby_gold_sale_price_config.name,zby_merchant_shop.id as merchant_shop_id,zby_merchant_shop.merchant_id,zby_merchant_shop.brand_id FROM zby_gold_sale_price_config LEFT JOIN zby_merchant_shop ON zby_gold_sale_price_config.merchant_shop_id=zby_merchant_shop.id WHERE zby_gold_sale_price_config.`status`=1 AND zby_merchant_shop.brand_id !=zby_gold_sale_price_config.brand_id").Scan(&list).Error; err != nil {
		panic(err)
	}

	for _, item := range list {
		var realItem *GoldSalePriceConfig

		if err := db.Table("zby_gold_sale_price_config").Select("id,name").Where("merchant_id = ? AND merchant_shop_id = ? AND brand_id = ? AND name = ? and status = 1 ", item.MerchantId, item.MerchantShopId, item.BrandId, item.Name).Find(&realItem).Error; err != nil {
			fmt.Printf("%v+", err)
			continue
		}

		var stock []*GoodsStock
		if err := db.Table("zby_goods_stock").Where("merchant_id = ? and gold_sale_price_type = ?", item.MerchantId, item.Id).Find(&stock).Error; err != nil {
			fmt.Printf("%v+", err)
			continue
		}

		if deleteType == 1 {
			if len(stock) == 0 {
				db.Table("zby_gold_sale_price_config").Where("id = ?", item.Id).Update("status", 2)
			}
		}

		if deleteType == 2 {
			if realItem.Id == 0 {
				continue
			}
			db.Table("zby_goods_stock").Where("merchant_id = ? and gold_sale_price_type = ?", item.MerchantId, item.Id).Update("gold_sale_price_type", realItem.Id)
		}
	}
}
