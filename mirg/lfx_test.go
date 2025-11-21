package mirg

import "testing"

var merchantId int32 = 1843
var status int32 = 1

/**
CREATE TABLE zby_customer_bp_detail_lfx
LIKE zby_customer_bp_detail;

-- 2. 将将要被删除的数据插入备份表
INSERT INTO zby_customer_bp_detail_lfx
SELECT *
FROM zby_customer_bp_detail
WHERE merchant_id = 1843;

-- 3. 删除原表中的这些数据
DELETE FROM zby_customer_bp_detail
WHERE merchant_id = 1843;
*/

/**
SELECT id FROM zby_customer WHERE merchant_id = 1843 order by id desc LIMIT 1;
SELECT id FROM zby_sale_order WHERE merchant_id = 1843 order by id desc LIMIT 1;
SELECT id FROM zby_refund_order WHERE merchant_id = 1843 order by id desc LIMIT 1;
SELECT id FROM zby_customer_bp_log  WHERE merchant_id = 1843 order by id desc LIMIT 1;
*/

func TestCustomerSync(t *testing.T) {
	CustomerSync(merchantId, 23573908, status)
}

func TestOrderSaleSync(t *testing.T) {
	OrderSaleSync(merchantId, 20137966, status)
}

func TestOrderRefundSync(t *testing.T) {
	OrderRefundSync(merchantId, 1843, status)
}

//func TestOrderRecycleSync(t *testing.T) {
//	OrderRecycleSync(merchantId, 1, status)
//}

func TestCustomerBpSync(t *testing.T) {
	CustomerBpSync(merchantId, 44267185, status)
}
