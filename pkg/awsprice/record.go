package awsprice

import (
	"encoding/json"
	"math"
	"sort"
)

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
	OnDemand                float64 `json:"ondemand"`                            // common
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

func (r *Record) String() string {
	bytea, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

func (r *Record) BreakevenPointInMonth() int {
	month := 12
	if r.LeaseContractLength == "3yr" {
		month = 12 * 3
	}

	breakEvenPoint := 0
	res := r.ReservedQuantity
	ond := 0.0
	for i := 1; i < month+1; i++ {
		ond = ond + r.OnDemand*24*float64(GetDays(i))
		res = res + r.ReservedHrs*24*float64(GetDays(i))
		if ond > res {
			breakEvenPoint = i
			break
		}
	}

	return breakEvenPoint
}

type Recommended struct {
	Record                 *Record `json:"record"`
	BreakevenPointInMonth  int     `json:"breakevenpoint_in_month"`
	OnDemandInstanceNumAvg float64 `json:"ondemand_instance_num_avg"`
	ReservedInstanceNum    int64   `json:"reserved_instance_num"`
	FullOnDemandCost       float64 `json:"full_ondemand_cost"`
	ReservedAppliedCost    Cost    `json:"reserved_applied_cost"`
	ReservedQuantity       float64 `json:"reserved_quantity"`
	Subtraction            float64 `json:"subtraction"`
	DiscountRate           float64 `json:"discount_rate"`
}

func (r *Recommended) String() string {
	bytea, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

type Forecast struct {
	Month       string  `json:"month"` // 2018-11
	InstanceNum float64 `json:"instance_num"`
}

type ReservedAppliedCost struct {
	LeaseContractLength string  `json:"lease_contract_length"`
	PurchaseOption      string  `json:"purchase_option"`
	FullOnDemand        float64 `json:"full_ondemand"`
	ReservedApplied     Cost    `json:"reserved_applied"`
	ReservedQuantity    float64 `json:"reserved_quantity"`
	Subtraction         float64 `json:"subtraction"`
	DiscountRate        float64 `json:"discount_rate"`
}

type Cost struct {
	OnDemand float64 `json:"ondemand"`
	Reserved float64 `json:"reserved"`
	Total    float64 `json:"total"`
}

func (r *Record) Recommend(forecast []Forecast) *Recommended {
	ondemand, reserved := r.recommendedInstanceNum(forecast)
	cost := r.GetCost(ondemand, reserved)

	return &Recommended{
		Record:                 r,
		BreakevenPointInMonth:  r.BreakevenPointInMonth(),
		OnDemandInstanceNumAvg: ondemand,
		ReservedInstanceNum:    reserved,
		FullOnDemandCost:       cost.FullOnDemand,
		ReservedAppliedCost:    cost.ReservedApplied,
		ReservedQuantity:       cost.ReservedQuantity,
		Subtraction:            cost.Subtraction,
		DiscountRate:           cost.DiscountRate,
	}
}

func (r *Record) recommendedInstanceNum(forecast []Forecast) (float64, int64) {
	p := r.BreakevenPointInMonth()
	if len(forecast) < p {
		sum := 0.0
		for i := range forecast {
			sum = sum + forecast[i].InstanceNum
		}

		return sum / float64(len(forecast)), 0
	}

	tmp := append([]Forecast{}, forecast...)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i].InstanceNum > tmp[j].InstanceNum })
	rnum := int64(math.Floor(tmp[p-1].InstanceNum))

	sum := 0.0
	for i := range tmp {
		ond := tmp[i].InstanceNum - float64(rnum)
		if ond > 0 {
			sum = sum + ond
		}
	}

	return sum / float64(len(forecast)), rnum
}

// ondemandNum, reservedNum is Per Year  (LeaseContractLength=1yr)
// ondemandNum, reservedNum is Per 3Year (LeaseContractLength=3yr)
func (r *Record) GetCost(ondemandNum float64, reservedNum int64) *ReservedAppliedCost {
	full := r.GetAnnualCost().OnDemand * (ondemandNum + float64(reservedNum))
	ond := r.GetAnnualCost().OnDemand * ondemandNum
	res := r.GetAnnualCost().Reserved * float64(reservedNum)

	out := &ReservedAppliedCost{
		LeaseContractLength: r.LeaseContractLength,
		PurchaseOption:      r.PurchaseOption,
		FullOnDemand:        full,
		ReservedApplied: Cost{
			OnDemand: ond,
			Reserved: res,
			Total:    ond + res,
		},
		ReservedQuantity: r.ReservedQuantity * float64(reservedNum),
	}

	out.Subtraction = full - out.ReservedApplied.Total
	out.DiscountRate = 1.0 - (out.ReservedApplied.Total / full)

	return out
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

type AnnualCost struct {
	LeaseContractLength string  `json:"lease_contract_length"`
	PurchaseOption      string  `json:"purchase_option"`
	OnDemand            float64 `json:"ondemand"`
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
