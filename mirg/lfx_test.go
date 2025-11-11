package mirg

import "testing"

var merchantId int32 = 258
var status int32 = 2

/**
SELECT id FROM zby_customer WHERE merchant_id = 1028 order by id desc LIMIT 1;
SELECT id FROM zby_sale_order WHERE merchant_id = 1028 order by id desc LIMIT 1;
SELECT id FROM zby_refund_order WHERE merchant_id = 1028 order by id desc LIMIT 1;
SELECT id FROM zby_recycle_order  WHERE merchant_id = 1028 order by id desc LIMIT 1;
SELECT id FROM zby_customer_bp_log  WHERE merchant_id = 1028 order by id desc LIMIT 1;
*/

/**
-- 门店id
-- 969

-- 会员id
-- 239586

-- 订单id
-- 2489532

*/

func TestCustomerSync(t *testing.T) {
	CustomerSync(merchantId, 239595, status)
}

func TestOrderSaleSync(t *testing.T) {
	OrderSaleSync(merchantId, 2489532, status)
}

func TestOrderRefundSync(t *testing.T) {
	OrderRefundSync(merchantId, 1, status)
}

func TestOrderRecycleSync(t *testing.T) {
	OrderRecycleSync(merchantId, 1, status)
}

func TestCustomerBpSync(t *testing.T) {
	CustomerBpSync(merchantId, 91457, status)
}
