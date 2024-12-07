package enums

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOperationType(t *testing.T) {
	tcs := []struct {
		name          string
		opertaionType int
		expectedEnum  OperationType
		expectedErr   error
	}{
		{
			name:          "Test ParseOperationType_Success",
			opertaionType: 1,
			expectedEnum:  NormalPurchase,
		},
		{
			name:          "Test ParseOperationType_Success",
			opertaionType: 3,
			expectedEnum:  Withdrawal,
		},
		{
			name:          "Test ParseOperationType_Failure",
			opertaionType: 5,
			expectedErr:   errors.New("sdd"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ot, err := ParseOperationType(tc.opertaionType)
			if tc.expectedErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.expectedEnum, ot)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAllowNegative(t *testing.T) {
	tcs := []struct {
		name          string
		opType        OperationType
		allowNegative bool
	}{
		{
			name:          "Test - Allow Negative True",
			allowNegative: true,
			opType:        NormalPurchase,
		},
		{
			name:          "Test Allow Negative True",
			allowNegative: true,
			opType:        Withdrawal,
		},
		{
			name:          "Test Allow Negative False",
			allowNegative: false,
			opType:        CreditVoucher,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			an := AllowNegative(tc.opType)

			require.Equal(t, tc.allowNegative, an)
		})
	}
}
