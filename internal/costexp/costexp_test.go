package costexp

import (
	"fmt"
	"os"
	"testing"
)

//
// func TestReservationCoverage(t *testing.T) {
// 	os.Setenv("AWS_PROFILE", "example")
//
// 	period := &costexplorer.DateInterval{
// 		Start: aws.String("2019-01-01"),
// 		End:   aws.String("2019-02-01"),
// 	}
//
// 	c := costexplorer.New(session.Must(session.NewSession()))
//
// 	{
// 		input := costexplorer.GetReservationCoverageInput{
// 			Granularity: aws.String("MONTHLY"),
// 			TimePeriod:  period,
// 		}
//
// 		out, err := c.GetReservationCoverage(&input)
// 		if err != nil {
// 			t.Errorf("get reservation coverage: %v", err)
// 		}
//
// 		fmt.Println(out)
// 	}
//
// 	{
// 		input := costexplorer.GetReservationUtilizationInput{
// 			Granularity: aws.String("MONTHLY"),
// 			TimePeriod:  period,
// 		}
//
// 		out, err := c.GetReservationUtilization(&input)
// 		if err != nil {
// 			t.Errorf("get reservation coverage: %v", err)
// 		}
//
// 		fmt.Println(out)
// 	}
// }

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
