package awsprice

import "encoding/json"

type RecordList []*Record

func (list RecordList) Region(region string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].Region != region {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) InstanceType(tipe string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].InstanceType != tipe {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) Engine(engine string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].Engine != engine {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) LeaseContractLength(length string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].LeaseContractLength != length {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) Tenancy(tenancy string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].Tenancy != tenancy {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) PurchaseOption(purchase string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].PurchaseOption != purchase {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) PreInstalled(preinstalled string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].PreInstalled != preinstalled {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) OperatingSystem(os string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].OperatingSystem != os {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) OfferingClass(class string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].OfferingClass != class {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

type Record struct {
	SKU                     string  `json:"sku"`                                 // common
	OfferTermCode           string  `json:"offer_term_code"`                     // common
	Region                  string  `json:"region"`                              // common
	InstanceType            string  `json:"instance_type"`                       // common
	UsageType               string  `json:"usage_type"`                          // common
	LeaseContractLength     string  `json:"lease_contract_length"`               // common
	PurchaseOption          string  `json:"purchase_option"`                     // common
	OnDemand                float64 `json:"on_demand"`                           // common
	ReservedQuantity        float64 `json:"reserved_quantity"`                   // common
	ReservedHrs             float64 `json:"reserved_hrs"`                        // common
	Tenancy                 string  `json:"tenancy,omitempty"`                   // ec2: Shared, Host, Dedicated
	PreInstalled            string  `json:"pre_installed,omitempty"`             // ec2:  SQL Web, SQL Ent, SQL Std, NA
	OperatingSystem         string  `json:"operating_system,omitempty"`          // ec2:  Windows, Linux, SUSE, RHEL
	Operation               string  `json:"operation,omitempty"`                 // ec2
	OfferingClass           string  `json:"offering_class,omitempty"`            // ec2, rds
	NormalizationSizeFactor string  `json:"normalization_size_factor,omitempty"` // ec2, rds
	Engine                  string  `json:"engine,omitempty"`                    // rds, cache
}

// ondemandNum, reservedNum is Per Year  (LeaseContractLength=1yr)
// ondemandNum, reservedNum is Per 3Year (LeaseContractLength=3yr)
func (r *Record) ExpectedCost(ondemandNum, reservedNum int) *ExpectedCost {
	full := r.GetAnnualCost().OnDemand * float64(ondemandNum+reservedNum)
	ond := r.GetAnnualCost().OnDemand * float64(ondemandNum)
	res := r.GetAnnualCost().Reserved * float64(reservedNum)

	out := &ExpectedCost{
		LeaseContractLength: r.LeaseContractLength,
		PurchaseOption:      r.PurchaseOption,
		FullOnDemand: Cost{
			OnDemand: full,
			Reserved: 0.0,
			Total:    full,
		},
		ReservedApplied: Cost{
			OnDemand: ond,
			Reserved: res,
			Total:    ond + res,
		},
		ReservedQuantity: r.ReservedQuantity * float64(reservedNum),
	}

	out.Subtraction = out.FullOnDemand.Total - out.ReservedApplied.Total
	out.DiscountRate = 1.0 - (out.ReservedApplied.Total / out.FullOnDemand.Total)

	return out
}

type ExpectedCost struct {
	LeaseContractLength string  `json:"lease_contract_length"`
	PurchaseOption      string  `json:"purchase_option"`
	FullOnDemand        Cost    `json:"full_on_demand"`
	ReservedApplied     Cost    `json:"reserved_applied"`
	ReservedQuantity    float64 `json:"reserved_quantity"`
	Subtraction         float64 `json:"subtraction"`
	DiscountRate        float64 `json:"discount_rate"`
}

type Cost struct {
	OnDemand float64 `json:"on_demand"`
	Reserved float64 `json:"reserved"`
	Total    float64 `json:"total"`
}

func (r *ExpectedCost) String() string {
	bytea, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

func (r *Record) GetAnnualCost() *AnnualCost {
	ret := &AnnualCost{
		LeaseContractLength: r.LeaseContractLength,
		PurchaseOption:      r.PurchaseOption,
	}

	hrs := 365 * 24
	if r.LeaseContractLength == "3yr" {
		hrs = hrs * 3
	}

	ret.OnDemand = r.OnDemand * float64(hrs)
	ret.Reserved = r.ReservedQuantity + r.ReservedHrs*float64(hrs)
	ret.Subtraction = ret.OnDemand - ret.Reserved
	ret.DiscountRate = 1.0 - ret.Reserved/ret.OnDemand
	ret.ReservedQuantity = r.ReservedQuantity

	return ret
}

func (r *Record) String() string {
	bytea, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

type AnnualCost struct {
	LeaseContractLength string  `json:"lease_contract_length"`
	PurchaseOption      string  `json:"purchase_option"`
	OnDemand            float64 `json:"on_demand"`
	Reserved            float64 `json:"reserved"`
	ReservedQuantity    float64 `json:"reserved_quantity"`
	Subtraction         float64 `json:"subtraction"`
	DiscountRate        float64 `json:"discount_rate"`
}

func (r *AnnualCost) String() string {
	bytea, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}
