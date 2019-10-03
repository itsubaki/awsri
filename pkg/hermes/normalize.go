package hermes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

func Normalize(quantity []usage.Quantity, mini map[string]pricing.Tuple) []usage.Quantity {
	n := make([]usage.Quantity, 0)
	for i := range quantity {
		hash := fmt.Sprintf(
			"%s%s%s%s",
			quantity[i].UsageType,
			OperatingSystem[quantity[i].Platform],
			quantity[i].CacheEngine,
			quantity[i].DatabaseEngine,
		)

		v, ok := mini[hash]
		if !ok {
			n = append(n, quantity[i])
			continue
		}

		if len(v.Minimum.NormalizationSizeFactor) < 1 {
			n = append(n, quantity[i])
			continue
		}

		s0, err := strconv.ParseFloat(v.Minimum.NormalizationSizeFactor, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid normalization size factor: %v", err))
		}

		s1, err := strconv.ParseFloat(v.Price.NormalizationSizeFactor, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid normalization size factor: %v", err))
		}

		scale := s1 / s0

		n = append(n, usage.Quantity{
			AccountID:      quantity[i].AccountID,
			Description:    quantity[i].Description,
			Region:         quantity[i].Region,
			UsageType:      v.Minimum.UsageType,
			Platform:       quantity[i].Platform,
			CacheEngine:    quantity[i].CacheEngine,
			DatabaseEngine: quantity[i].DatabaseEngine,
			Date:           quantity[i].Date,
			InstanceHour:   quantity[i].InstanceHour * scale,
			InstanceNum:    quantity[i].InstanceNum * scale,
		})
	}

	return n
}

// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/apply_ri.html
// Instance size flexibility does not apply to Reserved Instances
// that are purchased for a specific Availability Zone,
// bare metal instances,
// Reserved Instances with dedicated tenancy,
// and Reserved Instances for Windows,
// Windows with SQL Standard,
// Windows with SQL Server Enterprise,
// Windows with SQL Server Web,
// RHEL, and SLES.
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/apply_ri.html
// Instance size flexibility does not apply to Reserved Instances
// that are purchased for a specific Availability Zone,
// bare metal instances,
// Reserved Instances with dedicated tenancy,
// and Reserved Instances for Windows,
// Windows with SQL Standard,
// Windows with SQL Server Enterprise,
// Windows with SQL Server Web,
// RHEL, and SLES.
func HasFlexibility(p pricing.Price) bool {
	if strings.Contains(p.OperatingSystem, "Windows") {
		return false
	}

	if strings.Contains(p.OperatingSystem, "Red Hat Enterprise Linux") {
		return false
	}

	if strings.Contains(p.OperatingSystem, "SUSE Linux") {
		return false
	}

	if strings.Contains(p.Tenancy, "dedicated") {
		return false
	}

	if strings.Contains(p.InstanceType, "cache") {
		return false
	}

	return true
}
