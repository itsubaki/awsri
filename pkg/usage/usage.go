package usage

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type Account struct {
	ID          string
	Description string
}

type Quantity struct {
	AccountID      string  `json:"account_id,omitempty"`
	Description    string  `json:"description,omitempty"`
	Region         string  `json:"region,omitempty"`
	UsageType      string  `json:"usage_type,omitempty"`
	Platform       string  `json:"platform,omitempty"`
	CacheEngine    string  `json:"cache_engine,omitempty"`
	DatabaseEngine string  `json:"database_engine,omitempty"`
	Date           string  `json:"date,omitempty"`
	InstanceHour   float64 `json:"instance_hour,omitempty"`
	InstanceNum    float64 `json:"instance_num"`
}

type GetQuantityInput struct {
	AccountID   string
	Description string
	Dimension   string
	UsageType   []string
	Start       string
	End         string
}

func (q Quantity) OSEngine() string {
	return fmt.Sprintf("%s%s%s", OperatingSystem[q.Platform], q.CacheEngine, q.DatabaseEngine)
}

func (q Quantity) PFEngine() string {
	return fmt.Sprintf("%s%s%s", q.Platform, q.CacheEngine, q.DatabaseEngine)
}

func (q Quantity) String() string {
	return q.JSON()
}

func (q Quantity) JSON() string {
	bytes, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func Sort(quantity []Quantity) {
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].Date < quantity[j].Date })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].DatabaseEngine < quantity[j].DatabaseEngine })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].CacheEngine < quantity[j].CacheEngine })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].Platform < quantity[j].Platform })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].UsageType < quantity[j].UsageType })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].AccountID < quantity[j].AccountID })
}

type FetchFunc func(start, end string, account Account, usageType []string) ([]Quantity, error)

var FetchFuncList = []FetchFunc{
	fetchBoxUsage,
	//fetchSpotUsage,
	fetchNodeUsage,
	fetchInstanceUsage,
	fetchMultiAZUsage,
	fetchNode,
}

func Fetch(start, end string) ([]Quantity, error) {
	linkedAccount, err := FetchLinkedAccount(start, end)
	if err != nil {
		return nil, fmt.Errorf("get linked account: %v", err)
	}

	usageType, err := fetchUsageType(start, end)
	if err != nil {
		return nil, fmt.Errorf("get usage type: %v", err)
	}

	out := make([]Quantity, 0)
	for _, a := range linkedAccount {
		for _, f := range FetchFuncList {
			quantity, err := f(start, end, a, usageType)
			if err != nil {
				return nil, fmt.Errorf("get usage quantity: %v", err)
			}

			out = append(out, quantity...)
		}
	}

	return out, nil
}

func fetchBoxUsage(start, end string, account Account, usageType []string) ([]Quantity, error) {
	ut := make([]string, 0)
	for i := range usageType {
		if !strings.Contains(usageType[i], "BoxUsage:") {
			continue
		}
		ut = append(ut, usageType[i])
	}

	return fetchQuantity(&GetQuantityInput{
		AccountID:   account.ID,
		Description: account.Description,
		Dimension:   "PLATFORM",
		UsageType:   ut,
		Start:       start,
		End:         end,
	})
}

func fetchSpotUsage(start, end string, account Account, usageType []string) ([]Quantity, error) {
	ut := make([]string, 0)
	for i := range usageType {
		if !strings.Contains(usageType[i], "SpotUsage:") {
			continue
		}
		ut = append(ut, usageType[i])
	}

	return fetchQuantity(&GetQuantityInput{
		AccountID:   account.ID,
		Description: account.Description,
		Dimension:   "PLATFORM",
		UsageType:   ut,
		Start:       start,
		End:         end,
	})
}

func fetchNodeUsage(start, end string, account Account, usageType []string) ([]Quantity, error) {
	ut := make([]string, 0)
	for i := range usageType {
		if !strings.Contains(usageType[i], "NodeUsage:") {
			continue
		}
		ut = append(ut, usageType[i])
	}

	return fetchQuantity(&GetQuantityInput{
		AccountID:   account.ID,
		Description: account.Description,
		Dimension:   "CACHE_ENGINE",
		UsageType:   ut,
		Start:       start,
		End:         end,
	})
}

func fetchInstanceUsage(start, end string, account Account, usageType []string) ([]Quantity, error) {
	ut := make([]string, 0)
	for i := range usageType {
		if !strings.Contains(usageType[i], "InstanceUsage:") {
			continue
		}

		ut = append(ut, usageType[i])
	}

	return fetchQuantity(&GetQuantityInput{
		AccountID:   account.ID,
		Description: account.Description,
		Dimension:   "DATABASE_ENGINE",
		UsageType:   ut,
		Start:       start,
		End:         end,
	})
}

func fetchMultiAZUsage(start, end string, account Account, usageType []string) ([]Quantity, error) {
	ut := make([]string, 0)
	for i := range usageType {
		if !strings.Contains(usageType[i], "Multi-AZUsage:") {
			continue
		}

		ut = append(ut, usageType[i])
	}

	return fetchQuantity(&GetQuantityInput{
		AccountID:   account.ID,
		Description: account.Description,
		Dimension:   "DATABASE_ENGINE",
		UsageType:   ut,
		Start:       start,
		End:         end,
	})
}

func fetchNode(start, end string, account Account, usageType []string) ([]Quantity, error) {
	ut := make([]string, 0)
	for i := range usageType {
		if !strings.Contains(usageType[i], "Node:") {
			continue
		}
		ut = append(ut, usageType[i])
	}

	q, err := fetchQuantity(&GetQuantityInput{
		AccountID:   account.ID,
		Description: account.Description,
		Dimension:   "DATABASE_ENGINE",
		UsageType:   ut,
		Start:       start,
		End:         end,
	})
	if err != nil {
		return make([]Quantity, 0), err
	}

	out := make([]Quantity, 0)
	for i := range q {
		if q[i].DatabaseEngine != "NoDatabaseEngine" {
			continue
		}
		out = append(out, q[i])
	}

	return out, nil
}

func fetchQuantity(in *GetQuantityInput) ([]Quantity, error) {
	and := make([]*costexplorer.Expression, 0)
	and = append(and, &costexplorer.Expression{
		Dimensions: &costexplorer.DimensionValues{
			Key:    aws.String("LINKED_ACCOUNT"),
			Values: []*string{aws.String(in.AccountID)},
		},
	})

	or := make([]*costexplorer.Expression, 0)
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
		TimePeriod: &costexplorer.DateInterval{
			Start: &in.Start,
			End:   &in.End,
		},
	}

	if len(or) > 1 {
		input.Filter = &costexplorer.Expression{
			And: append(and, &costexplorer.Expression{Or: or}),
		}
	}

	c := costexplorer.New(session.Must(session.NewSession()))
	usage, err := c.GetCostAndUsage(&input)
	if err != nil {
		return []Quantity{}, fmt.Errorf("get cost and usage. or=%v: %v", or, err)
	}

	out := make([]Quantity, 0)
	for _, r := range usage.ResultsByTime {
		for _, g := range r.Groups {
			amount := *g.Metrics["UsageQuantity"].Amount
			if amount == "0" {
				continue
			}

			hrs, _ := strconv.ParseFloat(amount, 64)
			month := strings.Split(in.Start, "-")[1]
			num := hrs / float64(24*Days[month])

			index := strings.LastIndex(in.Start, "-")
			date := string(in.Start)[:index]
			q := Quantity{
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

			region, ok := region[strings.Split(q.UsageType, "-")[0]]
			if !ok {
				continue
			}
			q.Region = region

			out = append(out, q)
		}
	}

	return out, nil
}

func fetchUsageType(start, end string) ([]string, error) {
	input := costexplorer.GetDimensionValuesInput{
		Dimension: aws.String("USAGE_TYPE"),
		TimePeriod: &costexplorer.DateInterval{
			Start: &start,
			End:   &end,
		},
	}

	c := costexplorer.New(session.Must(session.NewSession()))
	val, err := c.GetDimensionValues(&input)
	if err != nil {
		return []string{}, fmt.Errorf("get dimenstion value: %v", err)
	}

	out := make([]string, 0)
	for _, d := range val.DimensionValues {
		out = append(out, *d.Value)
	}

	return out, nil
}

func FetchLinkedAccount(start, end string) ([]Account, error) {
	input := costexplorer.GetDimensionValuesInput{
		Dimension: aws.String("LINKED_ACCOUNT"),
		TimePeriod: &costexplorer.DateInterval{
			Start: &start,
			End:   &end,
		},
	}

	c := costexplorer.New(session.Must(session.NewSession()))
	val, err := c.GetDimensionValues(&input)
	if err != nil {
		return []Account{}, fmt.Errorf("get dimension values: %v", err)
	}

	out := make([]Account, 0)
	for _, v := range val.DimensionValues {
		out = append(out, Account{
			ID:          *v.Value,
			Description: *v.Attributes["description"],
		})
	}

	return out, nil
}
