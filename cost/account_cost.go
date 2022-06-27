package cost

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/itsubaki/hermes/usage"
)

type AccountCost struct {
	AccountID        string `json:"account_id"`
	Description      string `json:"description"`
	Date             string `json:"date,omitempty"`
	Service          string `json:"service,omitempty"`
	RecordType       string `json:"record_type,omitempty"`
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

func Fetch(start, end string, metrics []string) ([]AccountCost, error) {
	return FetchWith(start, end, []string{}, metrics)
}

func FetchWith(start, end string, dim, met []string) ([]AccountCost, error) {
	metrics := make([]*string, 0)
	for _, m := range met {
		metrics = append(metrics, aws.String(m))
	}

	if len(metrics) < 1 {
		metrics = []*string{
			aws.String("NetAmortizedCost"),
			aws.String("NetUnblendedCost"),
			aws.String("UnblendedCost"),
			aws.String("AmortizedCost"),
			aws.String("BlendedCost"),
		}
	}

	input := costexplorer.GetCostAndUsageInput{
		Metrics:     metrics,
		Granularity: aws.String("MONTHLY"),
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Key:  aws.String("SERVICE"),
				Type: aws.String("DIMENSION"),
			},
			{
				Key:  aws.String("RECORD_TYPE"),
				Type: aws.String("DIMENSION"),
			},
		},
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(start),
			End:   aws.String(end),
		},
	}

	la, err := usage.FetchLinkedAccount(start, end)
	if err != nil {
		return []AccountCost{}, fmt.Errorf("get linked account: %v", err)
	}

	out := make([]AccountCost, 0)
	for _, a := range la {
		if len(dim) == 0 {
			input.Filter = &costexplorer.Expression{
				Dimensions: &costexplorer.DimensionValues{
					Key:    aws.String("LINKED_ACCOUNT"),
					Values: []*string{aws.String(a.ID)},
				},
			}
		}

		if len(dim) > 0 {
			or := make([]*costexplorer.Expression, 0)
			for _, w := range dim {
				or = append(or, &costexplorer.Expression{
					Dimensions: &costexplorer.DimensionValues{
						Key:    aws.String("SERVICE"),
						Values: []*string{aws.String(w)},
					},
				})
			}

			and := make([]*costexplorer.Expression, 0)
			and = append(and, &costexplorer.Expression{
				Dimensions: &costexplorer.DimensionValues{
					Key:    aws.String("LINKED_ACCOUNT"),
					Values: []*string{aws.String(a.ID)},
				},
			})
			and = append(and, &costexplorer.Expression{Or: or})

			input.Filter = &costexplorer.Expression{
				And: and,
			}
		}

		o, err := fetch(start, end, &input)
		if err != nil {
			return o, fmt.Errorf("fetch: %v", err)
		}

		for i := range o {
			o[i].AccountID = a.ID
			o[i].Description = a.Description
		}

		out = append(out, o...)
	}

	return out, nil
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

		for _, r := range cost.ResultsByTime {
			for _, g := range r.Groups {
				o := AccountCost{
					Service:    *g.Keys[0],
					RecordType: *g.Keys[1],
					Date:       start,
				}

				if v, ok := g.Metrics["NetAmortizedCost"]; ok {
					o.NetAmortizedCost = Cost{
						Amount: *v.Amount,
						Unit:   *v.Unit,
					}
				}

				if v, ok := g.Metrics["NetUnblendedCost"]; ok {
					o.NetUnblendedCost = Cost{
						Amount: *v.Amount,
						Unit:   *v.Unit,
					}
				}

				if v, ok := g.Metrics["AmortizedCost"]; ok {
					o.AmortizedCost = Cost{
						Amount: *v.Amount,
						Unit:   *v.Unit,
					}
				}

				if v, ok := g.Metrics["UnblendedCost"]; ok {
					o.UnblendedCost = Cost{
						Amount: *v.Amount,
						Unit:   *v.Unit,
					}
				}

				if v, ok := g.Metrics["BlendedCost"]; ok {
					o.BlendedCost = Cost{
						Amount: *v.Amount,
						Unit:   *v.Unit,
					}
				}

				out = append(out, o)
			}
		}

		if cost.NextPageToken == nil {
			break
		}
		token = cost.NextPageToken
	}

	return out, nil
}
