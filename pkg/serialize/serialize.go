package serialize

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/itsubaki/awsri/internal/awsprice/cache"
	"github.com/itsubaki/awsri/internal/awsprice/ec2"
	"github.com/itsubaki/awsri/internal/awsprice/rds"
	internal "github.com/itsubaki/awsri/internal/costexp"
	"github.com/itsubaki/awsri/pkg/awsprice"
	"github.com/itsubaki/awsri/pkg/costexp"
)

type SerializeInput struct {
	Profile   string
	Date      []*costexplorer.DateInterval
	OutputDir string
}

func Serialize(input *SerializeInput) error {
	os.Setenv("AWS_PROFILE", input.Profile)

	for i := range input.Date {
		// start[:7] => 2018-12-01 -> 2018-12
		start := *input.Date[i].Start
		path := fmt.Sprintf("%s/%s_%s.out", input.OutputDir, input.Profile, start[:7])
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}

		repo := &costexp.Repository{
			Profile: input.Profile,
		}

		q, err := internal.New().GetUsageQuantity(input.Date[i])
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

type SerializeAWSPriceInput struct {
	Region    []string
	OutputDir string
}

func SerializeAWSPirice(input *SerializeAWSPriceInput) error {
	for _, r := range input.Region {

		path := fmt.Sprintf("%s/%s.out", input.OutputDir, r)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}

		repo := &awsprice.Repository{
			Region: []string{r},
		}

		{
			price, err := ec2.ReadPrice(r)
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
			price, err := cache.ReadPrice(r)
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
			price, err := rds.ReadPrice(r)
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
