package mirg

import (
	"fmt"
	"github.com/long250038728/web/tool/sliceconv"
	"testing"
)

func TestMerchantClear(t *testing.T) {
	merchantId := 432
	MerchantShopId := []int32{1633, 1634, 1635, 1636, 2151, 2775, 3576, 3866, 4563}

	var ids []int32

	db, readDb := NewDb()

	t.Run("zby_sale_order", func(t *testing.T) {
		if err := readDb.Table("zby_sale_order_goods").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_order_goods").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_sale_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_recycle_order", func(t *testing.T) {
		if err := readDb.Table("zby_recycle_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_recycle_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_recycle_order_goods").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_recycle_order_goods").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_sale_pay_log", func(t *testing.T) {
		if err := readDb.Table("zby_sale_pay_log").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_pay_log").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_refund_order", func(t *testing.T) {
		if err := readDb.Table("zby_refund_order_goods").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_refund_order_goods").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_refund_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_refund_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_sale_order_history", func(t *testing.T) {
		if err := readDb.Table("zby_sale_order_history").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_order_history").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_goods_stock", func(t *testing.T) {
		if err := readDb.Table("zby_goods_stock").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_goods_stock").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_allocation", func(t *testing.T) {
		if err := readDb.Table("zby_stock_allocation_record").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_allocation_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_stock_allocation_order").Where("merchant_id = ? AND (out_merchant_shop_id in (?) or (in_merchant_shop_id in (?)))", merchantId, MerchantShopId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_allocation_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_conversion_record", func(t *testing.T) {
		if err := readDb.Table("zby_stock_conversion_record").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_conversion_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_stock_conversion_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_conversion_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_modify_record", func(t *testing.T) {
		if err := readDb.Table("zby_stock_modify_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_modify_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_stock_modify_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_modify_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

	})

	t.Run("zby_stock_modify_type_record", func(t *testing.T) {
		if err := readDb.Table("zby_stock_modify_type_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_modify_type_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_stock_modify_type_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_modify_type_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

	})

	t.Run("zby_stock_price_order", func(t *testing.T) {
		if err := readDb.Table("zby_stock_price_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_price_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_stock_price_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_price_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_check_order", func(t *testing.T) {
		if err := readDb.Table("zby_stock_check_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_check_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		var recordIds []int32
		if err := readDb.Table("zby_stock_check_record_part_1").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&recordIds).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(recordIds, 1000) {
			res := db.Table("zby_stock_check_record_part_1").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_stock_check_record_part_2").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&recordIds).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(recordIds, 1000) {
			res := db.Table("zby_stock_check_record_part_2").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_stock_check_record_part_3").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&recordIds).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(recordIds, 1000) {
			res := db.Table("zby_stock_check_record_part_3").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_modify_order_v2", func(t *testing.T) {
		if err := readDb.Table("zby_stock_modify_order_v2").Where("merchant_id = ? AND (JSON_CONTAINS(JSON_ARRAY(?),JSON_KEYS(merchant_shop_agg)))", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_modify_order_v2").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		var recordIds []int32
		if err := readDb.Table("zby_stock_modify_order_v2_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&recordIds).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(recordIds, 1000) {
			res := db.Table("zby_stock_modify_order_v2_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

	})

	t.Run("zby_goods_import_order", func(t *testing.T) {
		if err := readDb.Table("zby_goods_import_order").Where("merchant_id = ? AND (merchant_shop_id in (?)  or merchant_shop_id = 0)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_goods_import_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		var recordIds []int32
		if err := readDb.Table("zby_goods_import_record").Where("merchant_id = ? AND import_id in (?)", merchantId, ids).Select("id").Find(&recordIds).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(recordIds, 1000) {
			res := db.Table("zby_goods_import_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

	})

	t.Run("zby_stock_export_order", func(t *testing.T) {
		if err := readDb.Table("zby_stock_export_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_export_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_stock_export_record").Where("merchant_id = ? AND export_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_export_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_add_order", func(t *testing.T) {
		if err := readDb.Table("zby_stock_add_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_add_order").Where("id in (?)", list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_stock_add_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_add_record").Where("id in (?)", list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_sale_performance", func(t *testing.T) {
		if err := readDb.Table("zby_sale_performance").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_performance").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	//==============================

	t.Run("zby_customer", func(t *testing.T) {
		if err := readDb.Table("zby_customer").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_customer").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)

			res = db.Table("zby_customer_bp_log").Where("merchant_id = ? AND customer_id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)

			res = db.Table("zby_customer_bp_detail").Where("merchant_id = ? AND customer_id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)

			res = db.Table("zby_customer_user_join").Where("merchant_id = ? AND customer_id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	//t.Run("zby_goods_stock_extend", func(t *testing.T) {
	//	if err := readDb.Table("zby_goods_stock_extend").Where("merchant_id = ?", merchantId).Select("id").Find(&ids).Error; err != nil {
	//		fmt.Println("order", err)
	//	}
	//	for _, list := range sliceconv.Chunk(ids, 1000) {
	//		res := db.Table("zby_goods_stock_extend").Where("id in (?)", list).Delete(nil)
	//		fmt.Println(res.RowsAffected)
	//	}
	//})

	//==============================

	t.Run("zby_old_stock", func(t *testing.T) {
		if err := readDb.Table("zby_old_stock").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_old_stock_export_order", func(t *testing.T) {
		if err := readDb.Table("zby_old_stock_export_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock_export_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_old_stock_export_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock_export_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_financial_profit_statement", func(t *testing.T) {
		if err := readDb.Table("zby_financial_profit_statement").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_financial_profit_statement").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_income_expenditure_record", func(t *testing.T) {
		if err := readDb.Table("zby_income_expenditure_record").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_income_expenditure_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_old_stock_period", func(t *testing.T) {
		if err := readDb.Table("zby_old_stock_period").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock_period").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_old_stock_auditor", func(t *testing.T) {
		if err := readDb.Table("zby_old_stock_auditor").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock_auditor").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	//==============================

	t.Run("zby_policy", func(t *testing.T) {
		if err := readDb.Table("zby_policy").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_policy").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_coupon_customer_record", func(t *testing.T) {
		if err := readDb.Table("zby_coupon_customer_record").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_coupon_customer_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_bonus_shop_record", func(t *testing.T) {
		if err := readDb.Table("zby_coupon_customer_record").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_bonus_shop_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Delete(nil)
			fmt.Println(res.RowsAffected)
		}
	})
}

//es
//sale_order_record_report
//import_sale_report.2023
//import_sale_report.2024
//import_sale_report.2025
//old_stock_import_sale_report
//sale_order_report
//zby_stock_allocation_record
//zby_customer
//sale_performance
//sale_performance_rel
//stock_price_record
//zby_goods_import_record
//zby_stock_export_record
