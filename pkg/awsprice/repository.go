package awsprice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Repository struct {
	Region   []string   `json:"region"`
	Internal RecordList `json:"internal"`
}

func NewRepository(path string) (*Repository, error) {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	var repo Repository
	if err := json.Unmarshal(read, &repo); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	return &repo, nil
}

func (r *Repository) SelectAll() RecordList {
	return r.Internal
}

func (r *Repository) FindMinimumInstanceType(record *Record) (*Record, error) {
	instanceType := record.InstanceType
	familiy := instanceType[:strings.LastIndex(instanceType, ".")]

	if strings.Contains(instanceType, "cache") {
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
