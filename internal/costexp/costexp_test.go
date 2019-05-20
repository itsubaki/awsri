package costexp

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func TestReservationCoverage(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	period := &costexplorer.DateInterval{
		Start: aws.String("2019-01-01"),
		End:   aws.String("2019-02-01"),
	}

	c := costexplorer.New(session.Must(session.NewSession()))

	{

		input := costexplorer.GetReservationCoverageInput{
			//Granularity: aws.String("MONTHLY"),
			GroupBy: []*costexplorer.GroupDefinition{
				{
					Key:  aws.String("REGION"),
					Type: aws.String("DIMENSION"),
				},
				{
					Key:  aws.String("INSTANCE_TYPE"),
					Type: aws.String("DIMENSION"),
				},
				{
					Key:  aws.String("LINKED_ACCOUNT"),
					Type: aws.String("DIMENSION"),
				},
				{
					Key:  aws.String("PLATFORM"),
					Type: aws.String("DIMENSION"),
				},
			},
			Filter: &costexplorer.Expression{
				Dimensions: &costexplorer.DimensionValues{
					Key: aws.String("SERVICE"),
					Values: []*string{
						aws.String("Amazon Elastic Compute Cloud - Compute"),
						//						aws.String("Amazon Relational Database Service"),
						//						aws.String("Amazon ElastiCache"),
						//						aws.String("Amazon Redshift"),
					},
				},
			},
			TimePeriod: period,
		}

		out, err := c.GetReservationCoverage(&input)
		if err != nil {
			t.Errorf("get reservation coverage: %v", err)
		}

		for _, c := range out.CoveragesByTime {
			for _, g := range c.Groups {
				if *g.Coverage.CoverageHours.ReservedHours == "0" {
					continue
				}

				fmt.Printf("%v %v %v %v %v\n",
					*g.Attributes["linkedAccount"],
					*g.Attributes["region"],
					*g.Attributes["instanceType"],
					*g.Attributes["platform"],
					*g.Coverage.CoverageHours.ReservedHours,
				)
			}
		}
	}

	{
		input := costexplorer.GetReservationUtilizationInput{
			Granularity: aws.String("MONTHLY"),
			TimePeriod:  period,
		}

		out, err := c.GetReservationUtilization(&input)
		if err != nil {
			t.Errorf("get reservation coverage: %v", err)
		}

		fmt.Println(out)
	}
}

func TestBilling(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	list, err := New().GetCost(&Date{
		Start: "2018-05-01",
		End:   "2018-06-01",
	})

	if err != nil {
		t.Errorf("get cost: %v", err)
	}

	for _, r := range list {
		fmt.Println(r)
	}
}

func TestLinkedAccount(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	list, err := New().GetLinkedAccount(&Date{
		Start: "2018-11-01",
		End:   "2018-12-01",
	})
	if err != nil {
		t.Errorf("get linked account: %v", err)
	}

	if len(list) < 1 {
		t.Errorf("linked account is empty")
	}
}

func TestUsageType(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	list, err := New().GetUsageType(&Date{
		Start: "2018-11-01",
		End:   "2018-12-01",
	})
	if err != nil {
		t.Errorf("get usage type: %v", err)
	}

	if len(list) < 1 {
		t.Errorf("usage type is empty")
	}
}

func TestGetUsageQuantity(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	list, err := New().GetUsageQuantity(&Date{
		Start: "2018-11-01",
		End:   "2018-11-02",
	})
	if err != nil {
		t.Errorf("get usage quantity: %v", err)
	}

	if len(list) < 1 {
		t.Errorf("usage quantity is empty")
	}
}
