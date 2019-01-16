package reserved

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/itsubaki/hermes/pkg/awsprice"
)

func TestGetReserved(t *testing.T) {
	path := "/var/tmp/hermes/awsprice/ap-northeast-1.out"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("file not found: %v", path)
	}

	os.Setenv("AWS_PROFILE", "example")
	os.Setenv("AWS_REGION", "ap-northeast-1")

	client := ec2.New(session.Must(session.NewSession()))
	input := &ec2.DescribeReservedInstancesInput{
		Filters: []*ec2.Filter{
			{Name: aws.String("state"), Values: []*string{aws.String("active")}},
		},
	}
	output, err := client.DescribeReservedInstances(input)
	if err != nil {
		t.Errorf("%v", err)
	}

	if len(output.ReservedInstances) < 1 {
		return
	}

	r := output.ReservedInstances[0]

	yr := "1yr"
	if *r.Duration == 94608000 {
		yr = "3yr"
	}

	os := "Linux"
	if strings.Contains(*r.ProductDescription, "Windows") {
		os = "Windows"
	}

	repo, err := awsprice.NewRepository(path)
	if err != nil {
		t.Errorf("%v", err)
	}
	rs := repo.FindByInstanceType(*r.InstanceType).
		OfferingClass(*r.OfferingClass).
		PurchaseOption(*r.OfferingType).
		OperatingSystem(os).
		LeaseContractLength(yr)

	if len(rs) != 1 {
		t.Errorf("invalid resultset length")
	}
}
