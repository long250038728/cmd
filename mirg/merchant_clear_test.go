package mirg

import (
	"fmt"
	"github.com/long250038728/web/tool/sliceconv"
	"testing"
)

func TestOrderCheckOrderContentClear(t *testing.T) {
	db, readDb := NewDb()
	var ids []int32
	if err := readDb.Table("zby_order_check").Where("id <= 1547044 and JSON_LENGTH(order_content) > 0").Select("id").Find(&ids).Error; err != nil {
		fmt.Println("zby_order_check", err)
	}

	t.Log(len(ids))

	for _, list := range sliceconv.Chunk(ids, 1000) {
		res := db.Table("zby_order_check").Where("id in (?)", list).Update("order_content", "{}")
		fmt.Println(res.RowsAffected)
	}
}

func TestMerchantClear(t *testing.T) {
	merchantId := 0
	MerchantShopId := []int32{1}
	newMerchantId := -0

	var ids []int32
	db, readDb := NewDb()

	t.Run("zby_sale_order", func(t *testing.T) {
		if err := readDb.Table("zby_sale_order_goods").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_order_goods").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_sale_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_sale_order_old_exchange", func(t *testing.T) {
		if err := readDb.Table("zby_sale_order_old_exchange").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_order_old_exchange").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_sale_order_new_exchange", func(t *testing.T) {
		if err := readDb.Table("zby_sale_order_new_exchange").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_order_new_exchange").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_recycle_order", func(t *testing.T) {
		if err := readDb.Table("zby_recycle_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_recycle_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_recycle_order_goods").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_recycle_order_goods").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_sale_pay_log", func(t *testing.T) {
		if err := readDb.Table("zby_sale_pay_log").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_pay_log").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_refund_order", func(t *testing.T) {
		if err := readDb.Table("zby_refund_order_goods").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_refund_order_goods").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_refund_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_refund_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_sale_order_history", func(t *testing.T) {
		if err := readDb.Table("zby_sale_order_history").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_order_history").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_goods_stock", func(t *testing.T) {
		if err := readDb.Table("zby_goods_stock").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_goods_stock").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_allocation", func(t *testing.T) {
		if err := readDb.Table("zby_stock_allocation_record").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_allocation_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_stock_allocation_order").Where("merchant_id = ? AND (out_merchant_shop_id in (?) or (in_merchant_shop_id in (?)))", merchantId, MerchantShopId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_allocation_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_conversion_record", func(t *testing.T) {
		if err := readDb.Table("zby_stock_conversion_record").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_conversion_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_stock_conversion_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_conversion_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_modify_record", func(t *testing.T) {
		if err := readDb.Table("zby_stock_modify_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_modify_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_stock_modify_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_modify_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

	})

	t.Run("zby_stock_modify_type_record", func(t *testing.T) {
		if err := readDb.Table("zby_stock_modify_type_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_modify_type_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_stock_modify_type_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_modify_type_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

	})

	t.Run("zby_stock_price_order", func(t *testing.T) {
		if err := readDb.Table("zby_stock_price_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_price_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_stock_price_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_price_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_check_order", func(t *testing.T) {
		if err := readDb.Table("zby_stock_check_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_check_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
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
			res := db.Table("zby_stock_check_record_part_1").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_stock_check_record_part_2").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&recordIds).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(recordIds, 1000) {
			res := db.Table("zby_stock_check_record_part_2").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_stock_check_record_part_3").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&recordIds).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(recordIds, 1000) {
			res := db.Table("zby_stock_check_record_part_3").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_modify_order_v2", func(t *testing.T) {
		strs := make([]string, 0, len(MerchantShopId))
		for _, MerchantShopIdInt := range MerchantShopId {
			strs = append(strs, fmt.Sprintf("%d", MerchantShopIdInt))
		}
		if err := readDb.Table("zby_stock_modify_order_v2").Where("merchant_id = ? AND (JSON_CONTAINS(JSON_ARRAY(?),JSON_KEYS(merchant_shop_agg)))", merchantId, strs).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_modify_order_v2").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
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
			res := db.Table("zby_stock_modify_order_v2_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

	})

	t.Run("zby_goods_import_order", func(t *testing.T) {
		// or merchant_shop_id = 0
		if err := readDb.Table("zby_goods_import_order").Where("merchant_id = ? AND (merchant_shop_id in (?))", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_goods_import_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
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
			res := db.Table("zby_goods_import_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

	})

	t.Run("zby_stock_export_order", func(t *testing.T) {
		if err := readDb.Table("zby_stock_export_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_export_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_stock_export_record").Where("merchant_id = ? AND export_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_export_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_add_order", func(t *testing.T) {
		if err := readDb.Table("zby_stock_add_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_add_order").Where("id in (?)", list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_stock_add_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_add_record").Where("id in (?)", list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_sale_performance", func(t *testing.T) {
		if err := readDb.Table("zby_sale_performance").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_sale_performance").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_stock_lossover_order", func(t *testing.T) {
		if err := readDb.Table("zby_stock_lossover_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_stock_lossover_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	//==============================

	t.Run("zby_servicing_order", func(t *testing.T) {
		if err := readDb.Table("zby_servicing_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_servicing_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	//t.Run("zby_goods_stock_extend", func(t *testing.T) {
	//	if err := readDb.Table("zby_goods_stock_extend").Where("merchant_id = ?", merchantId).Select("id").Find(&ids).Error; err != nil {
	//		fmt.Println("order", err)
	//	}
	//	for _, list := range sliceconv.Chunk(ids, 1000) {
	//		res := db.Table("zby_goods_stock_extend").Where("id in (?)", list).Update("merchant_id", newMerchantId)
	//		fmt.Println(res.RowsAffected)
	//	}
	//})

	//==============================

	t.Run("zby_old_stock", func(t *testing.T) {
		if err := readDb.Table("zby_old_stock").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_old_stock_export_order", func(t *testing.T) {
		if err := readDb.Table("zby_old_stock_export_order").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock_export_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_old_stock_export_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("goods", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock_export_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_financial_profit_statement", func(t *testing.T) {
		if err := readDb.Table("zby_financial_profit_statement").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_financial_profit_statement").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_income_expenditure_record", func(t *testing.T) {
		if err := readDb.Table("zby_income_expenditure_record").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_income_expenditure_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_old_stock_period", func(t *testing.T) {
		if err := readDb.Table("zby_old_stock_period").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock_period").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_old_stock_auditor", func(t *testing.T) {
		if err := readDb.Table("zby_old_stock_auditor").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock_auditor").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_old_stock_delivery_order", func(t *testing.T) {
		for _, shopId := range MerchantShopId {
			if err := readDb.Table("zby_old_stock_delivery_order").Where("merchant_id = ? AND find_in_set(?,merchant_shop_ids)", merchantId, shopId).Select("id").Find(&ids).Error; err != nil {
				fmt.Println("order", err)
			}
			for _, list := range sliceconv.Chunk(ids, 1000) {
				res := db.Table("zby_old_stock_delivery_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
				fmt.Println(res.RowsAffected)
			}

			if len(ids) == 0 {
				continue
			}

			if err := readDb.Table("zby_old_stock_delivery_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
				fmt.Println("goods", err)
			}
			for _, list := range sliceconv.Chunk(ids, 1000) {
				res := db.Table("zby_old_stock_delivery_record").Where("id in (?)", list).Update("merchant_id", newMerchantId)
				fmt.Println(res.RowsAffected)
			}
		}
	})

	t.Run("zby_old_stock_export_type", func(t *testing.T) {
		//if err := readDb.Table("zby_old_stock_export_type").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
		//	fmt.Println("order", err)
		//}
		//for _, list := range sliceconv.Chunk(ids, 1000) {
		//	res := db.Table("zby_old_stock_export_type").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
		//	fmt.Println(res.RowsAffected)
		//}

		if err := readDb.Table("zby_old_stock_export_order").Where("merchant_id = ? AND merchant_shop_ids in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock_export_order").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if len(ids) == 0 {
			return
		}

		if err := readDb.Table("zby_old_stock_export_record").Where("merchant_id = ? AND order_id in (?)", merchantId, ids).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_old_stock_export_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	//==============================

	t.Run("zby_policy", func(t *testing.T) {
		if err := readDb.Table("zby_policy").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_policy").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_coupon_customer_record", func(t *testing.T) {
		if err := readDb.Table("zby_coupon_customer_record").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			// merchant_id 定义为整形(更新为0)
			res := db.Table("zby_coupon_customer_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", 0)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_bonus_shop_record", func(t *testing.T) {
		if err := readDb.Table("zby_bonus_shop_record").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_bonus_shop_record").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}
	})

	t.Run("zby_customer", func(t *testing.T) {
		MerchantShopId = append(MerchantShopId, 0)

		if err := readDb.Table("zby_customer").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_customer").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_customer_bp_log").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_customer_bp_log").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
			fmt.Println(res.RowsAffected)
		}

		if err := readDb.Table("zby_customer_bp_detail").Where("merchant_id = ? AND merchant_shop_id in (?)", merchantId, MerchantShopId).Select("id").Find(&ids).Error; err != nil {
			fmt.Println("order", err)
		}
		for _, list := range sliceconv.Chunk(ids, 1000) {
			res := db.Table("zby_customer_bp_detail").Where("merchant_id = ? AND id in (?)", merchantId, list).Update("merchant_id", newMerchantId)
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
//goods_stock
//```
//POST /sale_performance/_update_by_query?wait_for_completion=true
//{
//     "script":{
//        "lang":"painless",
//        "source":"ctx._source.merchant_id = 0"
//      },
//        "query": {
//            "bool": {
//                    "must":[
//                            {
//                             "term":{
//                                    "merchant_id":361
//                                }
//
//                            },
//                            {
//                                "terms":{
//                                    "merchant_shop_id":[1472]
//                                }
//
//                            }
//
//                    ]
//
//                }
//            },
//    "track_total_hits":true
//}
//```

//sale_performance  有个查询不加商户id

//```
//POST /sale_performance/_update_by_query?wait_for_completion=true
//{
//     "script":{
//        "lang":"painless",
//        "source":"ctx._source.merchant_shop_id = 0"
//      },
//        "query": {
//            "bool": {
//                    "must":[
//
//                            {
//                                "terms":{
//                                    "merchant_shop_id":[1267, 1268, 1269, 1270, 1272]
//                                }
//
//                            }
//
//                    ]
//
//                }
//            },
//    "track_total_hits":true
//}
//```

//es
//sale_order_report   单据
//sale_order_record_report 单据明细
//import_sale_report.2023  进销存
//import_sale_report.2024  进销存
//import_sale_report.2025 进销存
//old_stock_import_sale_report 旧料入库跟销售
//zby_stock_allocation_record  调拨
//zby_customer 会员
//sale_performance  业绩
//sale_performance_rel  业绩关系
//stock_price_record 金价
//zby_goods_import_record 入库
//zby_stock_export_record 出库
