package costexp

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func TestCostExp(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")
	c := costexplorer.New(session.Must(session.NewSession()))

	input := costexplorer.GetCostAndUsageInput{
		Metrics:     []*string{aws.String("UsageQuantity")},
		Granularity: aws.String("MONTHLY"),
		GroupBy: []*costexplorer.GroupDefinition{
			{Key: aws.String("SERVICE"), Type: aws.String("DIMENSION")},
		},
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String("2018-11-01"),
			End:   aws.String("2018-11-30"),
		},
	}

	out, err := c.GetCostAndUsage(&input)
	if err != nil {
		t.Errorf("%v", err)
	}

	fmt.Println(out)

}
