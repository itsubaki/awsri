package serialize

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/itsubaki/awsri/internal/awsprice"
	"github.com/itsubaki/awsri/internal/awsprice/cache"
	"github.com/itsubaki/awsri/internal/awsprice/ec2"
	"github.com/itsubaki/awsri/internal/awsprice/rds"
	"github.com/itsubaki/awsri/internal/costexp"
)

func Serialize(profile string, date []*costexplorer.DateInterval) error {
	os.Setenv("AWS_PROFILE", profile)

	for i := range date {
		start := *date[i].Start

		path := fmt.Sprintf(
			"%s/%s/%s_%s.out",
			os.Getenv("GOPATH"),
			"src/github.com/itsubaki/awsri/internal/_serialized/costexp",
			profile,
			start[:7], // 2018-12-01 -> 2018-12
		)

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}

		repo := &costexp.Repository{
			Profile: profile,
		}

		q, err := costexp.New().GetUsageQuantity(date[i])
		if err != nil {
			return fmt.Errorf("get usage quantity: %v", err)
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
			return fmt.Errorf("marshal: %v", err)
		}

		if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
			return fmt.Errorf("write file: %v", err)
		}
	}

	return nil
}

func SerializeAWSPirice(region []string) error {
	for _, r := range region {
		path := fmt.Sprintf(
			"%s/%s/%s.out",
			os.Getenv("GOPATH"),
			"src/github.com/itsubaki/awsri/internal/_serialized/awsprice",
			r,
		)

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}

		repo := &awsprice.Repository{
			Region: r,
		}

		{
			price, err := ec2.GetPrice(r)
			if err != nil {
				return fmt.Errorf("read ec2 price file: %v", err)
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
			price, err := cache.GetPrice(r)
			if err != nil {
				return fmt.Errorf("read cache price file: %v", err)
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
			price, err := rds.GetPrice(r)
			if err != nil {
				return fmt.Errorf("read cache price file: %v", err)
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
			return fmt.Errorf("marshal: %v", err)
		}

		if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
			return fmt.Errorf("write file: %v", err)
		}
	}

	return nil
}
