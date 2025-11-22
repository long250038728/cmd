package mirg

import (
	"context"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/server/http"
	"time"
)

var baseUrl string = "https://applet.zhubaoe.cn"
var configPath = "./config/online/db.yaml"

//var baseUrl string = "https://mini.zhubaoe.cn"
//var configPath = "./config/test/db.yaml"

var idWhere = "<="

// CustomerSync 会员同步
func CustomerSync(merchantId, brandId, minId, maxId, status int32) {
	var ormConfig orm.Config
	configurator.NewYaml().MustLoad(configPath, &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}
	var ids []int32
	err = db.Table("zby_customer").Select("id").Where("merchant_id = ?  and brand_id = ? and status = 1 and id > ? and id < ?", merchantId, brandId, minId, maxId).Find(&ids).Error
	if err != nil {
		panic(err)
	}
	if len(ids) == 0 {
		panic(errors.New("会员没有数据"))
	}

	fmt.Println(fmt.Sprintf("=======count: %d=======", len(ids)))

	if status == 2 {
		httpClient := http.NewClient(http.SetTimeout(time.Second * 5))
		adder := fmt.Sprintf("%s/lmcrm/lfx_customer/push", baseUrl)
		for _, id := range ids {
			ctx := context.Background()
			data := map[string]any{
				"merchant_id": merchantId,
				"customer_id": id,
			}
			b, _, err := httpClient.Post(ctx, adder, data)
			if err != nil {
				fmt.Println(string(b))
				continue
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// OrderSaleSync 订单同步
func OrderSaleSync(merchantId, miniId, status int32) {
	var ormConfig orm.Config
	configurator.NewYaml().MustLoad(configPath, &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}

	var ids []int32
	err = db.Table("zby_sale_order").Select("id").Where(fmt.Sprintf("merchant_id = ? and status = 1 and id %s ? and customer_id > 0", idWhere), merchantId, miniId).Find(&ids).Error
	if err != nil {
		panic(err)
	}
	if len(ids) == 0 {
		panic(errors.New("销售单据没有数据"))
	}

	fmt.Println(fmt.Sprintf("=======count: %d=======", len(ids)))

	if status == 2 {
		httpClient := http.NewClient(http.SetTimeout(time.Second * 5))
		adder := fmt.Sprintf("%s/lmcrm/lfx_sale/push", baseUrl)
		for _, id := range ids {
			ctx := context.Background()
			data := map[string]any{
				"merchant_id": merchantId,
				"order_id":    id,
				"order_type":  8,
				"status":      1,
			}
			b, _, err := httpClient.Post(ctx, adder, data)
			if err != nil {
				fmt.Println(string(b))
				continue
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// OrderRefundSync 订单同步
func OrderRefundSync(merchantId, miniId, status int32) {
	var ormConfig orm.Config
	configurator.NewYaml().MustLoad(configPath, &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}

	var ids []int32
	err = db.Table("zby_refund_order").Select("id").Where(fmt.Sprintf("merchant_id = ? and status = 1 and id %s  ? and customer_id > 0", idWhere), merchantId, miniId).Find(&ids).Error
	if err != nil {
		panic(err)
	}
	if len(ids) == 0 {
		panic(errors.New("退货单据没有数据"))
	}

	fmt.Println(fmt.Sprintf("=======count: %d=======", len(ids)))

	if status == 2 {
		httpClient := http.NewClient(http.SetTimeout(time.Second * 5))
		adder := fmt.Sprintf("%s/lmcrm/lfx_sale/push", baseUrl)
		for _, id := range ids {
			ctx := context.Background()
			data := map[string]any{
				"merchant_id": merchantId,
				"order_id":    id,
				"order_type":  9,
				"status":      1,
			}
			b, _, err := httpClient.Post(ctx, adder, data)
			if err != nil {
				fmt.Println(string(b))
				continue
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// OrderRecycleSync 订单同步
func OrderRecycleSync(merchantId, miniId, status int32) {
	var ormConfig orm.Config
	configurator.NewYaml().MustLoad(configPath, &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}

	var ids []int32
	err = db.Table("zby_reshape_order").Select("id").Where(fmt.Sprintf("merchant_id = ? and status = 1 and id %s  ? and customer_id > 0", idWhere), merchantId, miniId).Find(&ids).Error
	if err != nil {
		panic(err)
	}
	if len(ids) == 0 {
		panic(errors.New("回收单据没有数据"))
	}

	if status == 1 {
		fmt.Println(ids)
	}

	if status == 2 {
		httpClient := http.NewClient(http.SetTimeout(time.Second * 5))
		adder := fmt.Sprintf("%s/lmcrm/lfx_sale/push", baseUrl)
		for _, id := range ids {
			ctx := context.Background()
			data := map[string]any{
				"merchant_id": merchantId,
				"order_id":    id,
				"order_type":  17,
				"status":      1,
			}
			b, _, err := httpClient.Post(ctx, adder, data)
			if err != nil {
				fmt.Println(string(b))
				continue
			}
			time.Sleep(time.Second)
		}
	}
}

// CustomerBpSync 积分同步
func CustomerBpSync(merchantId, miniId, status int32) {
	var ormConfig orm.Config
	configurator.NewYaml().MustLoad(configPath, &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}

	var bps []*CustomerBpLog
	err = db.Table("zby_customer_bp_log").Where(fmt.Sprintf("merchant_id = ? AND id %s  ?", idWhere), merchantId, miniId).Find(&bps).Error
	if err != nil {
		panic(err)
	}
	if len(bps) == 0 {
		panic(errors.New("会员积分没有数据"))
	}

	fmt.Println(fmt.Sprintf("=======count: %d=======", len(bps)))

	if status == 2 {
		httpClient := http.NewClient(http.SetTimeout(time.Second * 5))
		adder := fmt.Sprintf("%s/lmcrm/lfx_customer_bp/push", baseUrl)

		for index, bp := range bps {
			fmt.Println(fmt.Sprintf("=======%d======= :%d", index, bp.Id))

			ctx := context.Background()
			data := map[string]any{
				"merchant_id":      bp.MerchantId,
				"merchant_shop_id": bp.MerchantShopId,
				"customer_id":      bp.CustomerId,
				"pay_bonus":        bp.Point,
				"comment":          bp.Comment,
			}
			b, _, err := httpClient.Post(ctx, adder, data)
			if err != nil {
				fmt.Println(string(b))
				continue
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}
