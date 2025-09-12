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
