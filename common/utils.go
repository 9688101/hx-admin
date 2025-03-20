package common

import (
	"fmt"

	"github.com/9688101/hx-admin/global"
)

func LogQuota(quota int64) string {
	if global.DisplayInCurrencyEnabled {
		return fmt.Sprintf("＄%.6f 额度", float64(quota)/global.QuotaPerUnit)
	} else {
		return fmt.Sprintf("%d 点额度", quota)
	}
}
