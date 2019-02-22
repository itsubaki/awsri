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
	AccountID      string  `json:"account_id"`
	Description    string  `json:"description"`
	Region         string  `json:"region"`
	UsageType      string  `json:"usage_type"`
	Platform       string  `json:"platform,omitempty"`
	DatabaseEngine string  `json:"database_engine,omitempty"`
	CacheEngine    string  `json:"cache_engine,omitempty"`
	Date           string  `json:"date"`
	InstanceHour   float64 `json:"instance_hour"`
	InstanceNum    float64 `json:"instance_num"`
}

func (u *UsageQuantity) JSON() string {
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

func (c *CostExp) GetUsageQuantity(date *Date) (UsageQuantityList, error) {
	out := UsageQuantityList{}

	linkedAccount, err := c.GetLinkedAccount(date)
	if err != nil {
		return out, fmt.Errorf("get linked account: %v", err)
	}

	usageType, err := c.GetUsageType(date)
	if err != nil {
		return out, fmt.Errorf("get usage type: %v", err)
	}

	// compute
	for i := range linkedAccount {
		computeType := []string{}
		for i := range usageType {
			if !strings.Contains(usageType[i], "BoxUsage") {
				continue
			}
			computeType = append(computeType, usageType[i])
		}

		compute, err := c.getUsageQuantity(&getUsageQuantityInput{
			AccountID:   linkedAccount[i].AccountID,
			Description: linkedAccount[i].Description,
			Dimension:   "PLATFORM",
			UsageType:   computeType,
			Period: &costexplorer.DateInterval{
				Start: &date.Start,
				End:   &date.End,
			},
		})

		if err != nil {
			return out, fmt.Errorf("get compute usage quantity: %v", err)
		}

		out = append(out, compute...)
	}

	// cache
	for i := range linkedAccount {
		cacheUsageType := []string{}
		for i := range usageType {
			if !strings.Contains(usageType[i], "NodeUsage") {
				continue
			}
			cacheUsageType = append(cacheUsageType, usageType[i])
		}

		cache, err := c.getUsageQuantity(&getUsageQuantityInput{
			AccountID:   linkedAccount[i].AccountID,
			Description: linkedAccount[i].Description,
			Dimension:   "CACHE_ENGINE",
			UsageType:   cacheUsageType,
			Period: &costexplorer.DateInterval{
				Start: &date.Start,
				End:   &date.End,
			},
		})

		if err != nil {
			return out, fmt.Errorf("get cache usage quantity: %v", err)
		}

		out = append(out, cache...)
	}

	// database
	for i := range linkedAccount {
		databaseType := []string{}
		for i := range usageType {
			if !strings.Contains(usageType[i], "InstanceUsage") && !strings.Contains(usageType[i], "Multi-AZUsage") {
				continue
			}
			databaseType = append(databaseType, usageType[i])
		}

		db, err := c.getUsageQuantity(&getUsageQuantityInput{
			AccountID:   linkedAccount[i].AccountID,
			Description: linkedAccount[i].Description,
			Dimension:   "DATABASE_ENGINE",
			UsageType:   databaseType,
			Period: &costexplorer.DateInterval{
				Start: &date.Start,
				End:   &date.End,
			},
		})

		if err != nil {
			return out, fmt.Errorf("get database usage quantity: %v", err)
		}

		out = append(out, db...)
	}

	return out, nil
}

type getUsageQuantityInput struct {
	AccountID   string
	Description string
	Dimension   string
	UsageType   []string
	Period      *costexplorer.DateInterval
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

	if len(or) > 1 {
		input.Filter = &costexplorer.Expression{
			And: append(and, &costexplorer.Expression{Or: or}),
		}
	}

	usage, err := c.Client.GetCostAndUsage(&input)
	if err != nil {
		return out, fmt.Errorf("get cost and usage. or=%v: %v", or, err)
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
				Description:  in.Description,
				Date:         date,
				UsageType:    *g.Keys[0],
				InstanceHour: hrs,
				InstanceNum:  num,
			}

			if in.Dimension == "PLATFORM" {
				q.Platform = *g.Keys[1]
			}
			if in.Dimension == "CACHE_ENGINE" {
				q.CacheEngine = *g.Keys[1]
			}
			if in.Dimension == "DATABASE_ENGINE" {
				q.DatabaseEngine = *g.Keys[1]
			}

			region, ok := Region[strings.Split(q.UsageType, "-")[0]]
			if !ok {
				return nil, fmt.Errorf("region not found (usagetype=%s)", q.UsageType)
			}
			q.Region = region

			out = append(out, q)
		}
	}

	return out, nil
}

func (c *CostExp) GetUsageType(date *Date) ([]string, error) {
	usageTypeUnique := make(map[string]bool)
	input := costexplorer.GetDimensionValuesInput{
		Dimension: aws.String("USAGE_TYPE"),
		TimePeriod: &costexplorer.DateInterval{
			Start: &date.Start,
			End:   &date.End,
		},
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

func (c *CostExp) GetLinkedAccount(date *Date) ([]LinkedAccount, error) {
	out := []LinkedAccount{}

	input := costexplorer.GetDimensionValuesInput{
		Dimension: aws.String("LINKED_ACCOUNT"),
		TimePeriod: &costexplorer.DateInterval{
			Start: &date.Start,
			End:   &date.End,
		},
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
