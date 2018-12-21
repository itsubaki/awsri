package costexp

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func TestLinkedAccount(t *testing.T) {
	os.Setenv("AWS_PROFILE", "aws")
	period := &costexplorer.DateInterval{
		Start: aws.String("2018-11-01"),
		End:   aws.String("2018-12-01"),
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
		Start: aws.String("2018-11-01"),
		End:   aws.String("2018-12-01"),
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

func TestGetUsageQuantity(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")
	period := &costexplorer.DateInterval{
		Start: aws.String("2018-11-01"),
		End:   aws.String("2018-12-01"),
	}

	c := New()
	q, err := c.GetUsageQuantity(period)
	if err != nil {
		t.Errorf("get usage quantity: %v", err)
	}

	for _, qq := range q {
		fmt.Println(qq)
	}
}
