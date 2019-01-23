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
		Start: aws.String("2018-10-01"),
		End:   aws.String("2018-11-01"),
	}

	c := costexplorer.New(session.Must(session.NewSession()))

	{
		input := costexplorer.GetReservationCoverageInput{
			Granularity: aws.String("MONTHLY"),
			TimePeriod:  period,
		}

		out, err := c.GetReservationCoverage(&input)
		if err != nil {
			t.Errorf("get reservation coverage: %v", err)
		}

		fmt.Println(out)
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

func TestLinkedAccount(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	period := &costexplorer.DateInterval{
		Start: aws.String("2018-11-01"),
		End:   aws.String("2018-12-01"),
	}

	list, err := New().GetLinkedAccount(period)
	if err != nil {
		t.Errorf("get linked account: %v", err)
	}

	if len(list) < 1 {
		t.Errorf("linked account is empty")
	}
}

func TestUsageType(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	period := &costexplorer.DateInterval{
		Start: aws.String("2018-11-01"),
		End:   aws.String("2018-12-01"),
	}

	list, err := New().GetUsageType(period)
	if err != nil {
		t.Errorf("get usage type: %v", err)
	}

	if len(list) < 1 {
		t.Errorf("usage type is empty")
	}
}

func TestGetUsageQuantity(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	period := &costexplorer.DateInterval{
		Start: aws.String("2018-11-01"),
		End:   aws.String("2018-11-02"),
	}

	list, err := New().GetUsageQuantity(period)
	if err != nil {
		t.Errorf("get usage quantity: %v", err)
	}

	if len(list) < 1 {
		t.Errorf("usage quantity is empty")
	}
}
