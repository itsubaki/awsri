package costexp

import (
	"encoding/json"
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
	AccountID    string  `json:"account_id"`
	Date         string  `json:"date"`
	UsageType    string  `json:"usage_type"`
	Platform     string  `json:"platform,omitempty"`
	Engine       string  `json:"engine,omitempty"`
	InstanceHour float64 `json:"instance_hour"`
	InstanceNum  float64 `json:"instance_num"`
}

func (u *UsageQuantity) String() string {
	bytea, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	return string(bytea)
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
		// ec2
		{
			ec2UsageType := []string{}
			for i := range usageType {
				if !strings.Contains(usageType[i], "BoxUsage") {
					continue
				}
				ec2UsageType = append(ec2UsageType, usageType[i])
			}

			ec2, err := c.getUsageQuantity(&getUsageQuantityInput{
				AccountID: linkedAccount[i].AccountID,
				Dimension: "PLATFORM",
				UsageType: ec2UsageType,
				Period:    period,
			})

			if err != nil {
				return out, fmt.Errorf("get usage quantity: %v", err)
			}
			out = append(out, ec2...)
		}

		// cache
		{
			cacheUsageType := []string{}
			for i := range usageType {
				if !strings.Contains(usageType[i], "NodeUsage") {
					continue
				}
				cacheUsageType = append(cacheUsageType, usageType[i])
			}

			cache, err := c.getUsageQuantity(&getUsageQuantityInput{
				AccountID: linkedAccount[i].AccountID,
				Dimension: "CACHE_ENGINE",
				UsageType: cacheUsageType,
				Period:    period,
			})

			if err != nil {
				return out, fmt.Errorf("get usage quantity: %v", err)
			}
			out = append(out, cache...)
		}

		// db
		{
			dbUsageType := []string{}
			for i := range usageType {
				if !strings.Contains(usageType[i], "InstanceUsage") && !strings.Contains(usageType[i], "Multi-AZUsage") {

					continue
				}
				dbUsageType = append(dbUsageType, usageType[i])
			}

			db, err := c.getUsageQuantity(&getUsageQuantityInput{
				AccountID: linkedAccount[i].AccountID,
				Dimension: "DATABASE_ENGINE",
				UsageType: dbUsageType,
				Period:    period,
			})

			if err != nil {
				return out, fmt.Errorf("get usage quantity: %v", err)
			}
			out = append(out, db...)
		}
	}

	return out, nil
}

type getUsageQuantityInput struct {
	AccountID string
	Dimension string
	UsageType []string
	Period    *costexplorer.DateInterval
}

func (c *CostExp) getUsageQuantity(in *getUsageQuantityInput) (UsageQuantityList, error) {
	out := UsageQuantityList{}

	and := []*costexplorer.Expression{}
	and = append(and, &costexplorer.Expression{
		Dimensions: &costexplorer.DimensionValues{
			Key:    aws.String("LINKED_ACCOUNT"),
			Values: []*string{aws.String(in.AccountID)},
		},
	})

	or := []*costexplorer.Expression{}
	for i := range in.UsageType {
		or = append(or, &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String("USAGE_TYPE"),
				Values: []*string{aws.String(in.UsageType[i])},
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
				Key:  aws.String(in.Dimension),
				Type: aws.String("DIMENSION"),
			},
		},
		TimePeriod: in.Period,
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

			hrs, num := GetInstanceHourAndNum(amount, *in.Period.Start)
			index := strings.LastIndex(*in.Period.Start, "-")
			date := string(*in.Period.Start)[:index]
			q := &UsageQuantity{
				AccountID:    in.AccountID,
				Date:         date,
				UsageType:    *g.Keys[0],
				InstanceHour: hrs,
				InstanceNum:  num,
			}

			if in.Dimension == "PLATFORM" {
				q.Platform = *g.Keys[1]
			} else {
				q.Engine = *g.Keys[1]
			}

			out = append(out, q)
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
		"-BoxUsage:",
		"-NodeUsage:",
		"-InstanceUsage:",
		"-Multi-AZUsage:",
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
