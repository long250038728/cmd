package mirg

import (
	"testing"
)

// TestMerchantCustomerAddAction 会员新增
func TestMerchantCustomerAddAction(t *testing.T) {
	MerchantCustomerAddAction(false)
}

// TestCustomerMerchantShopMove 门店修改
func TestCustomerMerchantShopMove(t *testing.T) {
	MerchantAction(false, false)
}

// TestMerchantStaffAction 归属修改
func TestMerchantStaffAction(t *testing.T) {
	MerchantStaffAction(false)
}

//func TestTestJoinHttp(t *testing.T) {
//	TestJoinHttp()
//}
