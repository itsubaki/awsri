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

func (c *CostExp) GetUsageQuantity() (UsageQuantityList, error) {
	out := UsageQuantityList{}

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
