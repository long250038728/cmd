package model

type CustomerLog struct {
	MerchantId     int32  `json:"merchant_id"`
	MerchantShopId int32  `json:"merchant_shop_id"`
	CustomerId     int32  `json:"customer_id"`
	CustomerName   string `json:"customer_name"`
	Type           int32  `json:"type"`
	Comment        string `json:"comment"`
	CreateTime     string `json:"create_time"`
}

//=====================================================================

type CustomerBpLog struct {
	/*  */
	Id int32 `gorm:"primary_key;column:id;type:int(11);" json:"id"`
	/*  */
	MerchantId int32 `gorm:"column:merchant_id;type:int(11);" json:"merchant_id"`
	/* 品牌id */
	BrandId int32 `gorm:"column:brand_id;type:int(11);" json:"brand_id"`
	/*  */
	MerchantShopId int32 `gorm:"column:merchant_shop_id;type:int(11);" json:"merchant_shop_id"`
	/*  */
	CustomerId int32 `gorm:"column:customer_id;type:int(11);" json:"customer_id"`
	/*  */
	CustomerName string `gorm:"column:customer_name;type:varchar(128);" json:"customer_name"`
	/* 本地获得积分 */
	Point float64 `gorm:"column:point;type:decimal(12,2);" json:"point"`
	/* 获得总积分 */
	PointTotal float64 `gorm:"column:point_total;type:decimal(12,2);" json:"point_total"`
	/* 获得类型或者来源
	1-销售单 （收银端消费）+
	2-手动录入（手动调整）+
	3-连客积分 （转介绍顾客消费）+
	4-积分兑换 -
	5-激活微信会员卡+
	6-会员导入+
	7-积分抵现-
	8-取消兑换+
	9-退货-赠送积分收回-
	10-删单-赠送积分回收-
	11-新品换购-赠送积分收回
	12-积分清零
	13-退货-抵现积分退回
	14-删单-抵现积分退回
	15-新品换购-抵现积分退回
	16-盲盒抽奖 */
	Type int32 `gorm:"column:type;type:int(11);" json:"type"`
	/* 表示 该记录中的会员是否变更 （变更需要作废  更改订单信息 会导致出现多条数据导致消息多发）0 正常 1 已变更 */
	IsDel int32 `gorm:"column:is_del;type:tinyint(1);" json:"is_del"`
	/* 备注 */
	Comment string `gorm:"column:comment;type:varchar(256);" json:"comment"`
	/* 订单类型 */
	OrderType int32 `gorm:"column:order_type;type:tinyint(4);" json:"order_type"`
	/* 订单ID */
	OrderId int32 `gorm:"column:order_id;type:int(11);" json:"order_id"`
	/* 订单SN */
	OrderSn string `gorm:"column:order_sn;type:varchar(255);" json:"order_sn"`
	/* 订单商品ID */
	OrderGoodsId int32 `gorm:"column:order_goods_id;type:int(11);" json:"order_goods_id"`
	/* 条码 */
	StockCode string `gorm:"column:stock_code;type:varchar(128);" json:"stock_code"`
	/* 商品名称 */
	OrderGoodsName string `gorm:"column:order_goods_name;type:varchar(255);" json:"order_goods_name"`
	/* 活动类型 1-自定义 2-生日 3-纪念日 4会员日 */
	ActivityType int32 `gorm:"column:activity_type;type:tinyint(4);" json:"activity_type"`
	/* 活动id */
	ActivityId int32 `gorm:"column:activity_id;type:int(11);" json:"activity_id"`
	/* 活动名称 */
	ActivityName string `gorm:"column:activity_name;type:varchar(255);" json:"activity_name"`
	/* 过期时间 */
	ExpireTime int32 `gorm:"column:expire_time;type:int(11);" json:"expire_time"`
	/* 操作人id */
	AdminUserId int32 `gorm:"column:admin_user_id;type:int(11);" json:"admin_user_id"`
	/* 操作人name */
	AdminUserName string `gorm:"column:admin_user_name;type:varchar(255);" json:"admin_user_name"`
	/*  */
	CreateTime int32 `gorm:"column:create_time;type:int(11);" json:"create_time"`
	/*  */
	UpdateTime int32 `gorm:"column:update_time;type:int(11);" json:"update_time"`
	/*  */
	DetailId int32 `gorm:"column:detail_id;type:int(11);" json:"detail_id"`
}
