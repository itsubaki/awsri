package usage

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
)

func Normalize(q []Quantity, mini map[string]pricing.Tuple) []Quantity {
	out := make([]Quantity, 0)
	for i := range q {
		hash := fmt.Sprintf(
			"%s%s%s%s",
			q[i].UsageType,
			OperatingSystem[q[i].Platform],
			q[i].CacheEngine,
			q[i].DatabaseEngine,
		)

		v, ok := mini[hash]
		if !ok {
			out = append(out, q[i])
			continue
		}

		if len(v.Minimum.NormalizationSizeFactor) < 1 {
			out = append(out, q[i])
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

		out = append(out, Quantity{
			AccountID:      q[i].AccountID,
			Description:    q[i].Description,
			Region:         q[i].Region,
			UsageType:      v.Minimum.UsageType,
			Platform:       q[i].Platform,
			CacheEngine:    q[i].CacheEngine,
			DatabaseEngine: q[i].DatabaseEngine,
			Date:           q[i].Date,
			InstanceHour:   q[i].InstanceHour * scale,
			InstanceNum:    q[i].InstanceNum * scale,
		})
	}

	return out
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
