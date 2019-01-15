package serialize

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func TestSerializeCostExp(t *testing.T) {
	date := []*costexplorer.DateInterval{
		{
			Start: aws.String("2017-12-01"),
			End:   aws.String("2018-01-01"),
		},
		{
			Start: aws.String("2018-01-01"),
			End:   aws.String("2018-02-01"),
		},
		{
			Start: aws.String("2018-02-01"),
			End:   aws.String("2018-03-01"),
		},
		{
			Start: aws.String("2018-03-01"),
			End:   aws.String("2018-04-01"),
		},
		{
			Start: aws.String("2018-04-01"),
			End:   aws.String("2018-05-01"),
		},
		{
			Start: aws.String("2018-05-01"),
			End:   aws.String("2018-06-01"),
		},
		{
			Start: aws.String("2018-06-01"),
			End:   aws.String("2018-07-01"),
		},
		{
			Start: aws.String("2018-07-01"),
			End:   aws.String("2018-08-01"),
		},
		{
			Start: aws.String("2018-08-01"),
			End:   aws.String("2018-09-01"),
		},
		{
			Start: aws.String("2018-09-01"),
			End:   aws.String("2018-10-01"),
		},
		{
			Start: aws.String("2018-10-01"),
			End:   aws.String("2018-11-01"),
		},
		{
			Start: aws.String("2018-11-01"),
			End:   aws.String("2018-12-01"),
		},
	}

	input := SerializeInput{
		Profile: "example",
		Date:    date,
		OutputDir: fmt.Sprintf(
			"%s/%s",
			os.Getenv("GOPATH"),
			"src/github.com/itsubaki/hermes/internal/_serialized/costexp",
		),
	}

	if err := Serialize(&input); err != nil {
		t.Errorf("serialize costexp: %v", err)
	}
}

func TestSerializeAWSPrice(t *testing.T) {
	region := []string{
		"ap-northeast-1",
		"eu-central-1",
		"us-west-1",
		"us-west-2",
	}

	input := SerializeAWSPriceInput{
		Region: region,
		OutputDir: fmt.Sprintf(
			"%s/%s",
			os.Getenv("GOPATH"),
			"src/github.com/itsubaki/hermes/internal/_serialized/awsprice",
		),
	}

	if err := SerializeAWSPirice(&input); err != nil {
		t.Errorf("serialize aws price: %v", err)
	}
}
