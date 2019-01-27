package awsprice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/itsubaki/hermes/internal/awsprice/cache"
	"github.com/itsubaki/hermes/internal/awsprice/ec2"
	"github.com/itsubaki/hermes/internal/awsprice/rds"
)

type Repository struct {
	Region   []string   `json:"region"`
	Internal RecordList `json:"internal"`
}

func NewRepository(region []string) (*Repository, error) {
	repo := &Repository{
		Region: region,
	}

	for _, r := range region {
		{
			price, err := ec2.ReadPrice(r)
			if err != nil {
				return nil, fmt.Errorf("read ec2 price file: %v", err)
			}

			for k := range price {
				v := price[k]
				repo.Internal = append(repo.Internal, &Record{
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
				return nil, fmt.Errorf("read cache price file: %v", err)
			}
			for k := range price {
				v := price[k]
				repo.Internal = append(repo.Internal, &Record{
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
				return nil, fmt.Errorf("read cache price file: %v", err)
			}
			for k := range price {
				v := price[k]
				repo.Internal = append(repo.Internal, &Record{
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
	}

	return repo, nil
}

func Read(path string) (*Repository, error) {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	repo := &Repository{}
	if err := repo.Deserialize(read); err != nil {
		return nil, fmt.Errorf("new repository: %v", err)
	}

	return repo, nil
}

func (r *Repository) Write(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}

	bytes, err := r.Serialize()
	if err != nil {
		return fmt.Errorf("serialize: %v", err)
	}

	if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}

func (r *Repository) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		return []byte{}, fmt.Errorf("marshal: %v", err)
	}

	return bytes, nil
}

func (r *Repository) Deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, r); err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	return nil
}

func (r *Repository) SelectAll() RecordList {
	return r.Internal
}

func (r *Repository) FindMinimumInstanceType(record *Record) (*Record, error) {
	if strings.Contains(record.InstanceType, "cache") {
		return nil, fmt.Errorf("invalid input. cache hasn't normalization size factor")
	}

	defined := []string{
		"nano",
		"micro",
		"small",
		"medium",
		"large",
		"xlarge",
	}

	instanceType := record.InstanceType
	familiy := instanceType[:strings.LastIndex(instanceType, ".")]

	if len(record.OperatingSystem) > 0 {
		tmp := RecordList{}
		for i := range defined {
			suspect := fmt.Sprintf("%s.%s", familiy, defined[i])
			for j := range r.Internal {
				if r.Internal[j].InstanceType == suspect {
					tmp = append(tmp, r.Internal[j])
				}
			}
			if len(tmp) > 0 {
				break
			}
		}

		if len(tmp) < 1 {
			return nil, fmt.Errorf("undefined instance type. defined=%v", defined)
		}

		usageType := fmt.Sprintf("%s%s",
			record.UsageType[:strings.LastIndex(record.UsageType, ".")],
			tmp[0].UsageType[strings.LastIndex(tmp[0].UsageType, "."):],
		)

		rs := tmp.UsageType(usageType).
			OperatingSystem(record.OperatingSystem).
			LeaseContractLength(record.LeaseContractLength).
			PurchaseOption(record.PurchaseOption).
			PreInstalled(record.PreInstalled).
			OfferingClass(record.OfferingClass).
			Region(record.Region)

		if len(rs) != 1 {
			return nil, fmt.Errorf("invalid ec2 result set=%v", rs)
		}

		return rs[0], nil
	}

	if len(record.DatabaseEngine) > 0 {
		tmp := RecordList{}
		for i := range defined {
			suspect := fmt.Sprintf("%s.%s", familiy, defined[i])
			for j := range r.Internal {
				if r.Internal[j].InstanceType == suspect &&
					r.Internal[j].DatabaseEngine == record.DatabaseEngine {
					tmp = append(tmp, r.Internal[j])
				}
			}
			if len(tmp) > 0 {
				break
			}
		}

		if len(tmp) < 1 {
			return nil, fmt.Errorf("undefined instance type. defined=%v", defined)
		}

		usageType := fmt.Sprintf("%s%s",
			record.UsageType[:strings.LastIndex(record.UsageType, ".")],
			tmp[0].UsageType[strings.LastIndex(tmp[0].UsageType, "."):],
		)

		rs := tmp.UsageType(usageType).
			DatabaseEngine(record.DatabaseEngine).
			LeaseContractLength(record.LeaseContractLength).
			PurchaseOption(record.PurchaseOption).
			Region(record.Region)

		if len(rs) != 1 {
			return nil, fmt.Errorf("invalid database result set=%v", rs)
		}

		return rs[0], nil
	}

	return nil, fmt.Errorf("invalid record=%v", record)
}

func (r *Repository) FindByInstanceType(tipe string) RecordList {
	out := RecordList{}
	for i := range r.Internal {
		if r.Internal[i].InstanceType == tipe {
			out = append(out, r.Internal[i])
		}
	}

	return out
}

func (r *Repository) FindByUsageType(tipe string) RecordList {
	out := RecordList{}
	for i := range r.Internal {
		if r.Internal[i].UsageType == tipe {
			out = append(out, r.Internal[i])
		}
	}

	return out
}

func (r *Repository) Recommend(record *Record, forecast []Forecast, strategy ...string) (*Recommended, error) {
	min, err := r.FindMinimumInstanceType(record)
	if err != nil {
		return nil, fmt.Errorf("find minimum instance type: %v", err)
	}

	rf64, err := strconv.ParseFloat(record.NormalizationSizeFactor, 64)
	if err != nil {
		return nil, fmt.Errorf("parse float normalization size factor in record: %v", err)
	}

	mf64, err := strconv.ParseFloat(min.NormalizationSizeFactor, 64)
	if err != nil {
		return nil, fmt.Errorf("parse float normalization size factor in minimum: %v", err)
	}

	scale := rf64 / mf64

	out := record.Recommend(forecast, strategy...)
	out.MinimumRecord = min
	out.MinimumReservedInstanceNum = float64(out.ReservedInstanceNum) * scale

	return out, nil
}
