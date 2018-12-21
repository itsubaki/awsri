package serialize

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/itsubaki/awsri/internal/awsprice/cache"
	"github.com/itsubaki/awsri/internal/awsprice/ec2"
	"github.com/itsubaki/awsri/internal/awsprice/rds"
	internal "github.com/itsubaki/awsri/internal/costexp"
	"github.com/itsubaki/awsri/pkg/awsprice"
	"github.com/itsubaki/awsri/pkg/costexp"
)

func TestSerializeCostExpOverall(t *testing.T) {
	os.Setenv("AWS_PROFILE", "aws")

	plist := []*costexplorer.DateInterval{
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

	for i := range plist {
		start := *plist[i].Start

		path := fmt.Sprintf(
			"%s/%s/%s_%s.out",
			os.Getenv("GOPATH"),
			"src/github.com/itsubaki/awsri/internal/_serialized/costexp",
			os.Getenv("AWS_PROFILE"),
			start[:7], // 2018-12-01 -> 2018-12
		)

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return
		}

		repo := &costexp.Repository{
			Profile: os.Getenv("AWS_PROFILE"),
		}

		q, err := internal.New().GetUsageQuantity(plist[i])
		if err != nil {
			t.Errorf("get usage quantity: %v", err)
		}

		for _, qq := range q {
			repo.Internal = append(repo.Internal, &costexp.Record{
				AccountID:    qq.AccountID,
				Date:         qq.Date,
				UsageType:    qq.UsageType,
				Platform:     qq.Platform,
				Engine:       qq.Engine,
				InstanceHour: qq.InstanceHour,
				InstanceNum:  qq.InstanceNum,
			})
		}

		bytes, err := json.Marshal(repo)
		if err != nil {
			t.Errorf("marshal: %v", err)
		}

		if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
			t.Errorf("write file: %v", err)
		}
	}
}

func TestSerializeCostExp(t *testing.T) {
	os.Setenv("AWS_PROFILE", "example")

	period := &costexplorer.DateInterval{
		Start: aws.String("2018-11-01"),
		End:   aws.String("2018-12-01"),
	}

	path := fmt.Sprintf(
		"%s/%s/%s_%s.out",
		os.Getenv("GOPATH"),
		"src/github.com/itsubaki/awsri/internal/_serialized/costexp",
		os.Getenv("AWS_PROFILE"),
		"2018-11",
	)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return
	}

	repo := &costexp.Repository{
		Profile: os.Getenv("AWS_PROFILE"),
	}

	q, err := internal.New().GetUsageQuantity(period)
	if err != nil {
		t.Errorf("get usage quantity: %v", err)
	}

	for _, qq := range q {
		repo.Internal = append(repo.Internal, &costexp.Record{
			AccountID:    qq.AccountID,
			Date:         qq.Date,
			UsageType:    qq.UsageType,
			Platform:     qq.Platform,
			Engine:       qq.Engine,
			InstanceHour: qq.InstanceHour,
			InstanceNum:  qq.InstanceNum,
		})
	}

	bytes, err := json.Marshal(repo)
	if err != nil {
		t.Errorf("marshal: %v", err)
	}

	if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
		t.Errorf("write file: %v", err)
	}
}

func TestSerializeAWSPrice(t *testing.T) {
	region := []string{
		"ap-northeast-1",
		"eu-central-1",
		"us-west-1",
		"us-west-2",
	}

	for _, r := range region {
		path := fmt.Sprintf("%s/%s/%s.out", os.Getenv("GOPATH"), "src/github.com/itsubaki/awsri/internal/_serialized/awsprice", r)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}

		repo := &awsprice.Repository{
			Region: r,
		}

		{
			price, err := ec2.ReadPrice(r)
			if err != nil {
				t.Errorf("read ec2 price file: %v", err)
			}

			for k := range price {
				v := price[k]
				repo.Internal = append(repo.Internal, &awsprice.Record{
					InstanceType:            v.InstanceType,
					LeaseContractLength:     v.LeaseContractLength,
					NormalizationSizeFactor: v.NormalizationSizeFactor,
					OfferTermCode:           v.OfferTermCode,
					OfferingClass:           v.OfferingClass,
					OnDemand:                v.OnDemand,
					OperatingSystem:         v.OperatingSystem,
					Operation:               v.Operation,
					PreInstalled:            v.PreInstalled,
					PurchaseOption:          v.PurchaseOption,
					Region:                  v.Region,
					ReservedHrs:             v.ReservedHrs,
					ReservedQuantity:        v.ReservedQuantity,
					SKU:                     v.SKU,
					Tenancy:                 v.Tenancy,
					UsageType:               v.UsageType,
				})
			}
		}

		{
			price, err := cache.ReadPrice(r)
			if err != nil {
				t.Errorf("read cache price file: %v", err)
			}
			for k := range price {
				v := price[k]
				repo.Internal = append(repo.Internal, &awsprice.Record{
					Engine:              v.Engine,
					InstanceType:        v.InstanceType,
					LeaseContractLength: v.LeaseContractLength,
					OfferTermCode:       v.OfferTermCode,
					OnDemand:            v.OnDemand,
					PurchaseOption:      v.PurchaseOption,
					Region:              v.Region,
					ReservedHrs:         v.ReservedHrs,
					ReservedQuantity:    v.ReservedQuantity,
					SKU:                 v.SKU,
					UsageType:           v.UsageType,
				})
			}
		}

		{
			price, err := rds.ReadPrice(r)
			if err != nil {
				t.Errorf("read cache price file: %v", err)
			}
			for k := range price {
				v := price[k]
				repo.Internal = append(repo.Internal, &awsprice.Record{
					Engine:              v.Engine,
					InstanceType:        v.InstanceType,
					LeaseContractLength: v.LeaseContractLength,
					OfferTermCode:       v.OfferTermCode,
					OnDemand:            v.OnDemand,
					PurchaseOption:      v.PurchaseOption,
					Region:              v.Region,
					ReservedHrs:         v.ReservedHrs,
					ReservedQuantity:    v.ReservedQuantity,
					SKU:                 v.SKU,
					UsageType:           v.UsageType,
				})
			}
		}

		bytes, err := json.Marshal(repo)
		if err != nil {
			t.Errorf("marshal: %v", err)
		}

		if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
			t.Errorf("write file: %v", err)
		}
	}
}
