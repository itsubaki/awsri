package costexp

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type LinkedAccount struct {
	AccountID   string
	Description string
}

type UsageQuantityList []*UsageQuantity

type UsageQuantity struct {
	AccountID       string  `json:"account_id"`
	Date            string  `json:"date"`
	UsageType       string  `json:"usage_type"`
	OperatingSystem string  `json:"operating_system,omitempty"`
	Engine          string  `json:"engine,omitempty"`
	InstanceHour    float64 `json:"instance_hour"`
	InstanceNum     float64 `json:"instance_num"`
}

type CostExp struct {
	Client *costexplorer.CostExplorer
}

func New() *CostExp {
	return &CostExp{
		Client: costexplorer.New(session.Must(session.NewSession())),
	}
}

func (c *CostExp) GetUsageQuantity(period *costexplorer.DateInterval) (UsageQuantityList, error) {
	out := UsageQuantityList{}

	linkedAccount, err := c.GetLinkedAccount(period)
	if err != nil {
		return out, fmt.Errorf("get linked account: %v", err)
	}

	usageType, err := c.GetUsageType(period)
	if err != nil {
		return out, fmt.Errorf("get usage type: %v", err)
	}

	for i := range linkedAccount {
		and := []*costexplorer.Expression{}
		and = append(and, &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String("LINKED_ACCOUNT"),
				Values: []*string{aws.String(linkedAccount[i].AccountID)},
			},
		})

		or := []*costexplorer.Expression{}
		for i := range usageType {
			if !strings.Contains(usageType[i], "BoxUsage") {
				continue
			}

			or = append(or, &costexplorer.Expression{
				Dimensions: &costexplorer.DimensionValues{
					Key:    aws.String("USAGE_TYPE"),
					Values: []*string{aws.String(usageType[i])},
				},
			})
		}

		input := costexplorer.GetCostAndUsageInput{
			Filter: &costexplorer.Expression{
				And: append(and, &costexplorer.Expression{Or: or}),
			},
			Metrics:     []*string{aws.String("UsageQuantity")},
			Granularity: aws.String("MONTHLY"),
			GroupBy: []*costexplorer.GroupDefinition{
				{
					Key:  aws.String("USAGE_TYPE"),
					Type: aws.String("DIMENSION"),
				},
				{
					Key:  aws.String("PLATFORM"),
					Type: aws.String("DIMENSION"),
				},
			},
			TimePeriod: period,
		}

		usage, err := c.Client.GetCostAndUsage(&input)
		if err != nil {
			return out, fmt.Errorf("get cost and usage: %v", err)
		}

		for _, r := range usage.ResultsByTime {
			for _, g := range r.Groups {
				amount := *g.Metrics["UsageQuantity"].Amount
				if amount == "0" {
					continue
				}

				fmt.Printf("date=%v~%v, usage_type=%v, platform=%v, instance_hrs=%v\n", *period.Start, *period.End, *g.Keys[0], *g.Keys[1], amount)
			}
		}
	}

	return out, nil
}

func (c *CostExp) GetUsageType(period *costexplorer.DateInterval) ([]string, error) {
	usageTypeUnique := make(map[string]bool)
	input := costexplorer.GetDimensionValuesInput{
		Dimension:  aws.String("USAGE_TYPE"),
		TimePeriod: period,
	}

	val, err := c.Client.GetDimensionValues(&input)
	if err != nil {
		return []string{}, fmt.Errorf("get dimenstion value: %v", err)
	}

	filter := []string{
		"BoxUsage",
		"NodeUsage",
		"InstanceUsage",
		"Multi-AZUsage",
	}

	for _, d := range val.DimensionValues {
		for _, f := range filter {
			if strings.Contains(*d.Value, f) {
				usageTypeUnique[*d.Value] = true
			}
		}
	}

	out := []string{}
	for u := range usageTypeUnique {
		// remove BoxUsage:m4.large, APN1-BoxUsage
		if !strings.Contains(u, "-") {
			continue
		}
		if !strings.Contains(u, ":") {
			continue
		}

		out = append(out, u)
	}

	return out, nil
}

func (c *CostExp) GetLinkedAccount(period *costexplorer.DateInterval) ([]LinkedAccount, error) {
	out := []LinkedAccount{}

	input := costexplorer.GetDimensionValuesInput{
		Dimension:  aws.String("LINKED_ACCOUNT"),
		TimePeriod: period,
	}

	val, err := c.Client.GetDimensionValues(&input)
	if err != nil {
		return out, fmt.Errorf("get dimension value: %v", err)
	}

	for _, v := range val.DimensionValues {
		out = append(out, LinkedAccount{
			AccountID:   *v.Value,
			Description: *v.Attributes["description"],
		})
	}

	return out, nil
}
