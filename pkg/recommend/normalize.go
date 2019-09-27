package recommend

import (
	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

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
func Normalize(quantity usage.Quantity, price []pricing.Price) []usage.Quantity {
	// usage.Quantity -> pricing.Price
	// flexibility ok?
	// find base size
	// calc power of
	// instance num x power

	return []usage.Quantity{}
}
