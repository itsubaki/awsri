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

func TestCostExplorer(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")
	c := costexplorer.New(session.Must(session.NewSession()))

	period := &costexplorer.DateInterval{
		Start: aws.String("2018-01-01"),
		End:   aws.String("2018-11-30"),
	}

	input := costexplorer.GetDimensionValuesInput{
		Dimension:  aws.String("USAGE_TYPE"),
		TimePeriod: period,
	}

	usageType := []string{}
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
					usageType = append(usageType, *d.Value)
				}
			}
		}

		fmt.Println(usageType)
	}

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
			Filter:      &costexplorer.Expression{Or: or},
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
			Filter:      &costexplorer.Expression{Or: or},
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
			Filter:      &costexplorer.Expression{Or: or},
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
