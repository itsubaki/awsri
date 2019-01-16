package serialize

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func TestSerializeCostExp(t *testing.T) {
	date := []*costexplorer.DateInterval{
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
		{
			Start: aws.String("2018-12-01"),
			End:   aws.String("2019-01-01"),
		},
	}

	input := SerializeInput{
		Profile:   "example",
		Date:      date,
		OutputDir: "/var/tmp/hermes/costexp",
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
		Region:    region,
		OutputDir: "/var/tmp/hermes/awsprice",
	}

	if err := SerializeAWSPirice(&input); err != nil {
		t.Errorf("serialize aws price: %v", err)
	}
}

func TestSerializeReserved(t *testing.T) {
	input := SerializeReservedInput{
		Profile:   "example",
		OutputDir: "/var/tmp/hermes/reserved",
	}

	if err := SerializeReserved(&input); err != nil {
		t.Errorf("serialize reserved instance: %v", err)
	}
}
