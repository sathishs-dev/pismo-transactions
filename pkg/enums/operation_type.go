package enums

import "fmt"

type OperationType int

const (
	NormalPurchase OperationType = iota + 1
	PurchaseWithInstallments
	Withdrawal
	CreditVoucher
)

func AllowNegative(i OperationType) bool {
	switch i {
	case NormalPurchase, PurchaseWithInstallments, Withdrawal:
		return true
	case CreditVoucher:
		return false
	}

	return false
}

func ParseOperationType(i int) (OperationType, error) {
	switch i {
	case int(NormalPurchase):
		return NormalPurchase, nil
	case int(PurchaseWithInstallments):
		return PurchaseWithInstallments, nil
	case int(Withdrawal):
		return Withdrawal, nil
	case int(CreditVoucher):
		return CreditVoucher, nil
	}

	return NormalPurchase, fmt.Errorf("%d is not a valid operation type", i)
}
