package mirg

import (
	"testing"
)

func TestCustomerMerchantShopMove(t *testing.T) {
	MerchantAction(false, false)
}

func TestMerchantStaffAction(t *testing.T) {
	MerchantStaffAction(false, false)
}

func TestTestJoinHttp(t *testing.T) {
	TestJoinHttp()
}
