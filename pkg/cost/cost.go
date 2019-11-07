package cost

import (
	"bytes"
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
	Service          string `json:"service,omitempty"`
	UnblendedCost    Cost   `json:"unblended_cost"`     // volume discount for a single account
	BlendedCost      Cost   `json:"blended_cost"`       // volume discount across linked account
	AmortizedCost    Cost   `json:"amortized_cost"`     // unblended + ReservedInstances/12
	NetAmortizedCost Cost   `json:"net_amortized_cost"` // before discount
	NetUnblendedCost Cost   `json:"net_unblended_cost"` // before discount
}

type Cost struct {
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

func (a AccountCost) String() string {
	return a.JSON()
}

func (a AccountCost) JSON() string {
	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (a AccountCost) Pretty() string {
	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, b, "", " "); err != nil {
		panic(err)
	}

	return pretty.String()
}

func Fetch(start, end string) ([]AccountCost, error) {
	return FetchWith(start, end, []string{})
}

func FetchWith(start, end string, with []string) ([]AccountCost, error) {
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
			{
				Key:  aws.String("SERVICE"),
				Type: aws.String("DIMENSION"),
			},
		},
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(start),
			End:   aws.String(end),
		},
	}

	if len(with) == 1 {
		input.Filter = &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String("SERVICE"),
				Values: []*string{aws.String(with[0])},
			},
		}
	}

	if len(with) > 1 {
		or := make([]*costexplorer.Expression, 0)
		for _, w := range with {
			or = append(or, &costexplorer.Expression{
				Dimensions: &costexplorer.DimensionValues{
					Key:    aws.String("SERVICE"),
					Values: []*string{aws.String(w)},
				},
			})
		}

		input.Filter = &costexplorer.Expression{
			Or: or,
		}
	}

	return fetch(start, end, &input)
}

func fetch(start, end string, input *costexplorer.GetCostAndUsageInput) ([]AccountCost, error) {
	out := make([]AccountCost, 0)
	c := costexplorer.New(session.Must(session.NewSession()))

	var token *string
	for {
		input.NextPageToken = token

		cost, err := c.GetCostAndUsage(input)
		if err != nil {
			return []AccountCost{}, fmt.Errorf("get cost and usage: %v", err)
		}

		index := strings.LastIndex(start, "-")
		date := start[:index]

		for _, r := range cost.ResultsByTime {
			for _, g := range r.Groups {
				o := AccountCost{
					AccountID: *g.Keys[0],
					Service:   *g.Keys[1],
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
				}

				out = append(out, o)
			}
		}

		if cost.NextPageToken == nil {
			break
		}
		token = cost.NextPageToken
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
