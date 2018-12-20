package costexp

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func TestLinkedAccount(t *testing.T) {
	os.Setenv("AWS_PROFILE", "aws")
	period := &costexplorer.DateInterval{
		Start: aws.String("2018-01-01"),
		End:   aws.String("2018-11-30"),
	}

	c := New()
	list, err := c.GetLinkedAccount(period)
	if err != nil {
		t.Errorf("get linked account: %v", err)
	}

	if len(list) < 1 {
		t.Errorf("linked account is empty")
	}
}

func TestUsageType(t *testing.T) {
	os.Setenv("AWS_PROFILE", "aws")
	period := &costexplorer.DateInterval{
		Start: aws.String("2018-01-01"),
		End:   aws.String("2018-11-30"),
	}

	c := New()
	list, err := c.GetUsageType(period)
	if err != nil {
		t.Errorf("get usage type: %v", err)
	}

	if len(list) < 1 {
		t.Errorf("usage type is empty")
	}
}

func TestCostExplorer(t *testing.T) {
	os.Setenv("AWS_PROFILE", "aws")
	c := costexplorer.New(session.Must(session.NewSession()))

	period := &costexplorer.DateInterval{
		Start: aws.String("2018-11-01"),
		End:   aws.String("2018-11-30"),
	}

	linked := []string{}
	{
		input := costexplorer.GetDimensionValuesInput{
			Dimension:  aws.String("LINKED_ACCOUNT"),
			TimePeriod: period,
		}

		out, err := c.GetDimensionValues(&input)
		if err != nil {
			t.Error(err)
		}

		for _, v := range out.DimensionValues {
			linked = append(linked, *v.Value)
		}

		for _, v := range out.DimensionValues {
			fmt.Printf("%v %v\n", *v.Value, *v.Attributes["description"])
		}
	}

	usageTypeUnique := make(map[string]bool)
	for i := 0; i < 3; i++ {
		and := []*costexplorer.Expression{}
		and = append(and, &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String("LINKED_ACCOUNT"),
				Values: []*string{aws.String(linked[i])},
			},
		})

		input := costexplorer.GetDimensionValuesInput{
			Dimension:  aws.String("USAGE_TYPE"),
			TimePeriod: period,
		}

		{
			out, err := c.GetDimensionValues(&input)
			if err != nil {
				t.Errorf("%v", err)
			}

			filter := []string{
				"BoxUsage",
				"NodeUsage",
				"InstanceUsage",
				"Multi-AZUsage",
			}

			for _, d := range out.DimensionValues {
				for _, f := range filter {
					if strings.Contains(*d.Value, f) {
						usageTypeUnique[*d.Value] = true
					}
				}
			}
		}
	}

	usageType := []string{}
	for k := range usageTypeUnique {
		usageType = append(usageType, k)
	}

	for i := 0; i < 3; i++ {
		fmt.Println(linked[i])

		and := []*costexplorer.Expression{}
		and = append(and, &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String("LINKED_ACCOUNT"),
				Values: []*string{aws.String(linked[i])},
			},
		})

		// ec2
		{
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

			out, err := c.GetCostAndUsage(&input)
			if err != nil {
				t.Error(err)
			}

			for _, r := range out.ResultsByTime {
				start := *r.TimePeriod.Start
				end := *r.TimePeriod.End
				for _, g := range r.Groups {
					amount := *g.Metrics["UsageQuantity"].Amount
					if amount == "0" {
						continue
					}

					fmt.Printf("date=%v~%v, usage_type=%v, platform=%v, instance_hrs=%v\n", start, end, *g.Keys[0], *g.Keys[1], amount)
				}
			}
		}

		// elasticache
		{
			or := []*costexplorer.Expression{}
			for i := range usageType {
				if !strings.Contains(usageType[i], "NodeUsage") {
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
						Key:  aws.String("CACHE_ENGINE"),
						Type: aws.String("DIMENSION"),
					},
				},
				TimePeriod: period,
			}

			out, err := c.GetCostAndUsage(&input)
			if err != nil {
				t.Error(err)
			}

			for _, r := range out.ResultsByTime {
				start := *r.TimePeriod.Start
				end := *r.TimePeriod.End
				for _, g := range r.Groups {
					amount := *g.Metrics["UsageQuantity"].Amount
					if amount == "0" {
						continue
					}

					fmt.Printf("date=%v~%v, usage_type=%v, engine=%v, instance_hrs=%v\n", start, end, *g.Keys[0], *g.Keys[1], amount)
				}
			}
		}

		// rds
		{
			or := []*costexplorer.Expression{}
			for i := range usageType {
				if !strings.Contains(usageType[i], "InstanceUsage") && !strings.Contains(usageType[i], "Multi-AZUsage") {
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
						Key:  aws.String("DATABASE_ENGINE"),
						Type: aws.String("DIMENSION"),
					},
				},
				TimePeriod: period,
			}

			out, err := c.GetCostAndUsage(&input)
			if err != nil {
				t.Error(err)
			}

			for _, r := range out.ResultsByTime {
				start := *r.TimePeriod.Start
				end := *r.TimePeriod.End
				for _, g := range r.Groups {
					amount := *g.Metrics["UsageQuantity"].Amount
					if amount == "0" {
						continue
					}

					fmt.Printf("date=%v~%v, usage_type=%v, engine=%v, instance_hrs=%v\n", start, end, *g.Keys[0], *g.Keys[1], amount)
				}
			}
		}
	}
}
