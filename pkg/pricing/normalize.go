package pricing

import (
	"fmt"
	"strings"
)

type NormalizeFunc func(repo *Repository, record *Record) (*Normalized, error)

func NewNormalizeFunc() []NormalizeFunc {
	return []NormalizeFunc{
		NormalizeCompute,
		NormalizeDatabase,
	}
}

func NormalizeCompute(repo *Repository, record *Record) (*Normalized, error) {
	if !record.Compute() {
		return nil, nil
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

	tmp := RecordList{}
	for i := range defined {
		suspect := fmt.Sprintf("%s.%s", familiy, defined[i])
		for j := range repo.Internal {
			if repo.Internal[j].InstanceType == suspect &&
				strings.LastIndex(repo.Internal[j].UsageType, ".") > 0 {
				tmp = append(tmp, repo.Internal[j])
			}
		}
		if len(tmp) > 0 {
			break
		}
	}

	if len(tmp) < 1 {
		return nil, fmt.Errorf("undefined instance type=%s family=%s. defined=%v", instanceType, familiy, defined)
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
		return nil, fmt.Errorf("invalid compute result set=%v", rs)
	}

	return &Normalized{Record: rs[0]}, nil
}

func NormalizeDatabase(repo *Repository, record *Record) (*Normalized, error) {
	if !record.Database() {
		return nil, nil
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

	tmp := RecordList{}
	for i := range defined {
		suspect := fmt.Sprintf("%s.%s", familiy, defined[i])
		for j := range repo.Internal {
			if repo.Internal[j].InstanceType == suspect &&
				repo.Internal[j].DatabaseEngine == record.DatabaseEngine &&
				strings.LastIndex(repo.Internal[j].UsageType, ".") > 0 {
				tmp = append(tmp, repo.Internal[j])
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

	return &Normalized{Record: rs[0]}, nil
}
