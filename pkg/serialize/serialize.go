package serialize

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	awscache "github.com/aws/aws-sdk-go/service/elasticache"
	awsrds "github.com/aws/aws-sdk-go/service/rds"
	"github.com/itsubaki/hermes/internal/awsprice/cache"
	"github.com/itsubaki/hermes/internal/awsprice/ec2"
	"github.com/itsubaki/hermes/internal/awsprice/rds"
	internal "github.com/itsubaki/hermes/internal/costexp"
	"github.com/itsubaki/hermes/pkg/awsprice"
	"github.com/itsubaki/hermes/pkg/costexp"
	"github.com/itsubaki/hermes/pkg/reserved"
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
				AccountID:      qq.AccountID,
				Description:    qq.Description,
				Date:           qq.Date,
				UsageType:      qq.UsageType,
				Platform:       qq.Platform,
				CacheEngine:    qq.CacheEngine,
				DatabaseEngine: qq.DatabaseEngine,
				InstanceHour:   qq.InstanceHour,
				InstanceNum:    qq.InstanceNum,
			})
		}

		bytes, err := json.Marshal(repo)
		if err != nil {
			return fmt.Errorf("marshal: %v", err)
		}

		if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
			return fmt.Errorf("write file: %v", err)
		}

		fmt.Printf("write file: %s\n", path)
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
					CacheEngine:         v.CacheEngine,
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
					DatabaseEngine:          v.DatabaseEngine,
					InstanceType:            v.InstanceType,
					LeaseContractLength:     v.LeaseContractLength,
					NormalizationSizeFactor: v.NormalizationSizeFactor,
					OfferTermCode:           v.OfferTermCode,
					OnDemand:                v.OnDemand,
					PurchaseOption:          v.PurchaseOption,
					Region:                  v.Region,
					ReservedHrs:             v.ReservedHrs,
					ReservedQuantity:        v.ReservedQuantity,
					SKU:                     v.SKU,
					UsageType:               v.UsageType,
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

		fmt.Printf("write file: %s\n", path)
	}

	return nil
}

type SerializeReservedInput struct {
	Profile   string
	Region    []string
	OutputDir string
}

func SerializeReserved(input *SerializeReservedInput) error {
	path := fmt.Sprintf("%s/%s.out", input.OutputDir, input.Profile)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return err
	}

	repo := &reserved.Repository{
		Profile: input.Profile,
		Region:  input.Region,
	}

	for _, region := range input.Region {
		os.Setenv("AWS_PROFILE", input.Profile)
		os.Setenv("AWS_REGION", region)

		{
			client := awsec2.New(session.Must(session.NewSession()))
			output, err := client.DescribeReservedInstances(&awsec2.DescribeReservedInstancesInput{
				Filters: []*awsec2.Filter{
					{Name: aws.String("state"), Values: []*string{aws.String("active")}},
				},
			})
			if err != nil {
				return fmt.Errorf("describe reserved instances: %v", err)
			}

			for _, r := range output.ReservedInstances {
				repo.Internal = append(repo.Internal, &reserved.Record{
					Region:             region,
					Duration:           *r.Duration,
					OfferingType:       *r.OfferingType,
					OfferingClass:      *r.OfferingClass,
					ProductDescription: *r.ProductDescription,
					InstanceType:       *r.InstanceType,
					InstanceCount:      *r.InstanceCount,
					Start:              *r.Start,
				})
			}
		}

		{
			client := awscache.New(session.Must(session.NewSession()))
			var maker *string
			for {
				input := &awscache.DescribeReservedCacheNodesInput{}
				if maker != nil {
					input.Marker = maker
				}

				output, err := client.DescribeReservedCacheNodes(input)
				if err != nil {
					return fmt.Errorf("describe reserved cachenode: %v", err)
				}

				for _, r := range output.ReservedCacheNodes {
					if *r.State != "active" {
						continue
					}
					repo.Internal = append(repo.Internal, &reserved.Record{
						Region:             region,
						Duration:           *r.Duration,
						OfferingType:       *r.OfferingType,
						ProductDescription: *r.ProductDescription,
						CacheNodeType:      *r.CacheNodeType,
						CacheNodeCount:     *r.CacheNodeCount,
						Start:              *r.StartTime,
					})
				}

				if maker == nil {
					break
				}
			}
		}

		{
			client := awsrds.New(session.Must(session.NewSession()))
			var maker *string
			for {
				input := &awsrds.DescribeReservedDBInstancesInput{}
				if maker != nil {
					input.Marker = maker
				}

				output, err := client.DescribeReservedDBInstances(input)
				if err != nil {
					return fmt.Errorf("describe reserved db instance: %v", err)
				}

				for _, r := range output.ReservedDBInstances {
					if *r.State != "active" {
						continue
					}
					repo.Internal = append(repo.Internal, &reserved.Record{
						Region:             region,
						Duration:           *r.Duration,
						OfferingType:       *r.OfferingType,
						ProductDescription: *r.ProductDescription,
						DBInstanceClass:    *r.DBInstanceClass,
						DBInstanceCount:    *r.DBInstanceCount,
						Start:              *r.StartTime,
						MultiAZ:            *r.MultiAZ,
					})
				}

				if maker == nil {
					break
				}
			}
		}
	}

	bytes, err := json.Marshal(repo)
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	fmt.Printf("write file: %s\n", path)

	return nil
}
