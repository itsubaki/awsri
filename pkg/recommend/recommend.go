package recommend

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/itsubaki/hermes/pkg/pricing"
	"github.com/itsubaki/hermes/pkg/usage"
)

type Recommended struct {
	Price pricing.Price `json:"price"`
	Cost  Cost          `json:"cost"`
	Usage Usage         `json:"usage"`
}

func (r Recommended) String() string {
	return r.JSON()
}

func (r Recommended) JSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (r Recommended) Pretty() string {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, b, "", " "); err != nil {
		panic(err)
	}

	return pretty.String()
}

type Usage struct {
	TotalHours          float64 `json:"total_hours"`
	OnDemandHours       float64 `json:"ondemand_hours"`
	ReservedHours       float64 `json:"reserved_hours"`
	ReservedInstanceNum int     `json:"reserved_instance_num"`
}

type Cost struct {
	FullOnDemand    float64             `json:"full_ondemand"`
	Saving          float64             `json:"saving"`
	ReservedApplied ReservedAppliedCost `json:"reserved_applied"`
}

type ReservedAppliedCost struct {
	OnDemand float64      `json:"ondemand"`
	Reserved ReservedCost `json:"reserved"`
	Total    float64      `json:"total"`
}

type ReservedCost struct {
	Quantity float64 `json:"quantity"`
	Hours    float64 `json:"hours"`
}

func Do(monthly []usage.Quantity, price pricing.Price) Recommended {
	reserved, _ := BreakEvenPoint(monthly, price)

	reservedHours := 0.0
	for _, m := range monthly {
		reservedHours = reservedHours + reserved.InstanceNum*float64(24*usage.Days[strings.Split(m.Date, "-")[1]])
	}

	ondemandHours := 0.0
	for _, m := range monthly {
		onDemandNum := m.InstanceNum - reserved.InstanceNum
		if onDemandNum < 0 {
			continue
		}

		ondemandHours = ondemandHours + onDemandNum*float64(24*usage.Days[strings.Split(m.Date, "-")[1]])
	}

	totalHours := 0.0
	for _, m := range monthly {
		totalHours = totalHours + m.InstanceNum*float64(24*usage.Days[strings.Split(m.Date, "-")[1]])
	}

	fullOndemandCost := totalHours * price.OnDemand
	reservedOndemandCost := ondemandHours * price.OnDemand
	reservedQuantityCost := reserved.InstanceNum * price.ReservedQuantity
	reservedHoursCost := reservedHours * price.ReservedHrs
	reservedTotalCost := reservedOndemandCost + reservedQuantityCost + reservedHoursCost

	return Recommended{
		Price: price,
		Usage: Usage{
			TotalHours:          totalHours,
			OnDemandHours:       ondemandHours,
			ReservedHours:       reservedHours,
			ReservedInstanceNum: int(reserved.InstanceNum),
		},
		Cost: Cost{
			FullOnDemand: fullOndemandCost,
			Saving:       fullOndemandCost - reservedTotalCost,
			ReservedApplied: ReservedAppliedCost{
				OnDemand: reservedOndemandCost,
				Reserved: ReservedCost{
					Quantity: reservedQuantityCost,
					Hours:    reservedHoursCost,
				},
				Total: reservedTotalCost,
			},
		},
	}
}

func Do2(monthly []usage.Quantity, price pricing.Price) Recommended {
	totalHours := 0.0
	for _, m := range monthly {
		totalHours = totalHours + m.InstanceNum*float64(24*usage.Days[strings.Split(m.Date, "-")[1]])
	}

	reserved, _ := BreakEvenPoint(monthly, price)
	reservedHours := 0.0
	for _, m := range monthly {
		reservedHours = reservedHours + reserved.InstanceNum*float64(24*usage.Days[strings.Split(m.Date, "-")[1]])
	}
	ondemandHours := totalHours - reservedHours

	fullOndemandCost := totalHours * price.OnDemand
	reservedOndemandCost := ondemandHours * price.OnDemand
	reservedQuantityCost := reserved.InstanceNum * price.ReservedQuantity
	reservedHoursCost := reservedHours * price.ReservedHrs
	reservedTotalCost := reservedOndemandCost + reservedQuantityCost + reservedHoursCost

	return Recommended{
		Price: price,
		Usage: Usage{
			TotalHours:          totalHours,
			OnDemandHours:       ondemandHours,
			ReservedHours:       reservedHours,
			ReservedInstanceNum: int(reserved.InstanceNum),
		},
		Cost: Cost{
			FullOnDemand: fullOndemandCost,
			Saving:       fullOndemandCost - reservedTotalCost,
			ReservedApplied: ReservedAppliedCost{
				OnDemand: reservedOndemandCost,
				Reserved: ReservedCost{
					Quantity: reservedQuantityCost,
					Hours:    reservedHoursCost,
				},
				Total: reservedTotalCost,
			},
		},
	}
}
