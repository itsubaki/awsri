package recommend

import (
	"github.com/itsubaki/hermes/pkg/pricing"
)

func BreakEvenPoint(p pricing.Price) int {
	month := 12
	if p.LeaseContractLength == "3yr" {
		month = 12 * 3
	}

	out, ond, res := 0, 0.0, p.ReservedQuantity
	for i := 1; i < month+1; i++ {
		ond, res = ond+p.OnDemand*24*float64(Days[i]), res+p.ReservedHrs*24*float64(Days[i])
		if ond > res {
			out = i
			break
		}
	}

	return out
}
