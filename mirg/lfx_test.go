package mirg

import "testing"

var merchantId int32 = 1028
var status int32 = 1

/**
SELECT id FROM zby_customer WHERE merchant_id = 1028 order by id desc LIMIT 1;
SELECT id FROM zby_sale_order WHERE merchant_id = 1028 order by id desc LIMIT 1;
SELECT id FROM zby_refund_order WHERE merchant_id = 1028 order by id desc LIMIT 1;
SELECT id FROM zby_recycle_order  WHERE merchant_id = 1028 order by id desc LIMIT 1;
SELECT id FROM zby_customer_bp_log  WHERE merchant_id = 1028 order by id desc LIMIT 1;
*/

func TestCustomerSync(t *testing.T) {
	CustomerSync(merchantId, 1, status)
}

func TestOrderSaleSync(t *testing.T) {
	OrderSaleSync(merchantId, 1, status)
}

func TestOrderRefundSync(t *testing.T) {
	OrderRefundSync(merchantId, 1, status)
}

func TestOrderRecycleSync(t *testing.T) {
	OrderRecycleSync(merchantId, 1, status)
}

func TestCustomerBpSync(t *testing.T) {
	CustomerBpSync(merchantId, 1, status)
}
