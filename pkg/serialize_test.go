package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/itsubaki/awsri/internal/awsprice/cache"
	"github.com/itsubaki/awsri/internal/awsprice/ec2"
	"github.com/itsubaki/awsri/internal/awsprice/rds"
	"github.com/itsubaki/awsri/internal/costviz"
	"github.com/itsubaki/awsri/pkg/awsprice"
)

func TestAWSPriceSerialize(t *testing.T) {
	region := []string{
		"ap-northeast-1",
		"eu-central-1",
		"us-west-1",
		"us-west-2",
	}

	for _, r := range region {
		path := fmt.Sprintf("%s/%s/%s.out", os.Getenv("GOPATH"), "src/github.com/itsubaki/awsri/internal/_serialized", r)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}

		repo := &awsprice.Repository{}
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

		read, err := ioutil.ReadFile(path)
		if err != nil {
			t.Errorf("read file: %v", err)
		}

		var repo2 awsprice.Repository
		if err := json.Unmarshal(read, &repo2); err != nil {
			t.Errorf("unmarshal: %v", err)
		}

		for _, v := range repo2.SelectAll() {
			fmt.Println(v)
		}
	}
}

func TestCostVizSerialize(t *testing.T) {
	if len(os.Getenv("COSTVIZ_BASEURL")) < 1 {
		return
	}

	for _, id := range strings.Split(os.Getenv("COSTVIZ_ACCOUNTID"), ",") {
		path := fmt.Sprintf("%s/%s/%s.out", os.Getenv("GOPATH"), "src/github.com/itsubaki/awsri/internal/_serialized", id)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}

		c := &costviz.CostViz{
			BaseURL:   os.Getenv("COSTVIZ_BASEURL"),
			XApiKey:   os.Getenv("COSTVIZ_XAPIKEY"),
			AccountID: id,
			TableName: []string{
				"awsbilling_201806",
				"awsbilling_201807",
				"awsbilling_201808",
				"awsbilling_201809",
				"awsbilling_201810",
				"awsbilling_201811",
			},
		}

		u, err := c.GetUtilization()
		if err != nil {
			t.Error(err)
		}

		for _, uu := range u {
			fmt.Println(uu)
		}

		if err := Serialize(u, path); err != nil {
			t.Error(err)
		}
	}
}

func Serialize(r costviz.RecordList, path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}

	bytes, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}
