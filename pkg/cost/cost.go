package cost

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/itsubaki/hermes/pkg/usage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type AccountCost struct {
	AccountID        string `json:"account_id"`
	Description      string `json:"description"`
	Date             string `json:"date,omitempty"`
	AmortizedCost    Cost   `json:"amortized_cost"`
	NetAmortizedCost Cost   `json:"net_amortized_cost"`
	UnblendedCost    Cost   `json:"unblended_cost"`
	NetUnblendedCost Cost   `json:"net_unblended_cost"`
	BlendedCost      Cost   `json:"blended_cost"`
}

type Cost struct {
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

func (a AccountCost) String() string {
	return a.JSON()
}

func (a AccountCost) JSON() string {
	bytes, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func Fetch(start, end string) ([]AccountCost, error) {
	input := costexplorer.GetCostAndUsageInput{
		Metrics: []*string{
			aws.String("NetAmortizedCost"),
			aws.String("NetUnblendedCost"),
			aws.String("UnblendedCost"),
			aws.String("AmortizedCost"),
			aws.String("BlendedCost"),
		},
		Granularity: aws.String("MONTHLY"),
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Key:  aws.String("LINKED_ACCOUNT"),
				Type: aws.String("DIMENSION"),
			},
		},
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(start),
			End:   aws.String(end),
		},
	}

	c := costexplorer.New(session.Must(session.NewSession()))
	cost, err := c.GetCostAndUsage(&input)
	if err != nil {
		return []AccountCost{}, fmt.Errorf("get cost and usage: %v", err)
	}

	index := strings.LastIndex(start, "-")
	date := string(start)[:index]

	out := make([]AccountCost, 0)
	for _, r := range cost.ResultsByTime {
		for _, g := range r.Groups {
			out = append(out, AccountCost{
				AccountID: *g.Keys[0],
				Date:      date,
				AmortizedCost: Cost{
					Amount: *g.Metrics["AmortizedCost"].Amount,
					Unit:   *g.Metrics["AmortizedCost"].Unit,
				},
				NetAmortizedCost: Cost{
					Amount: *g.Metrics["NetAmortizedCost"].Amount,
					Unit:   *g.Metrics["NetAmortizedCost"].Unit,
				},
				UnblendedCost: Cost{
					Amount: *g.Metrics["UnblendedCost"].Amount,
					Unit:   *g.Metrics["UnblendedCost"].Unit,
				},
				NetUnblendedCost: Cost{
					Amount: *g.Metrics["NetUnblendedCost"].Amount,
					Unit:   *g.Metrics["NetUnblendedCost"].Unit,
				},
				BlendedCost: Cost{
					Amount: *g.Metrics["BlendedCost"].Amount,
					Unit:   *g.Metrics["BlendedCost"].Unit,
				},
			})

		}
	}

	a, err := usage.FetchLinkedAccount(start, end)
	if err != nil {
		return []AccountCost{}, fmt.Errorf("get linked account: %v", err)
	}

	for i := range out {
		for _, aa := range a {
			if out[i].AccountID != aa.ID {
				continue
			}
			out[i].Description = aa.Description
		}
	}

	return out, nil
}
