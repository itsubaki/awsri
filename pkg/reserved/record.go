package reserved

import (
	"encoding/json"
	"strings"
	"time"
)

type RecordList []*Record

func (list RecordList) String() string {
	bytea, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

type Record struct {
	ReservedID         string    `json:"reserved_id"`
	Region             string    `json:"region"`
	Duration           int64     `json:"duration"`
	OfferingType       string    `json:"offering_type"`
	OfferingClass      string    `json:"offering_class,omitempty"`
	ProductDescription string    `json:"product_description"`
	InstanceType       string    `json:"instance_type,omitempty"`
	InstanceCount      int64     `json:"instance_count,omitempty"`
	CacheNodeType      string    `json:"cache_node_type,omitempty"`
	CacheNodeCount     int64     `json:"cache_node_count,omitempty"`
	DBInstanceClass    string    `json:"db_instance_class,omitempty"`
	DBInstanceCount    int64     `json:"db_instance_count,omitempty"`
	MultiAZ            bool      `json:"multi_az,omitempty"`
	Start              time.Time `json:"start"`
	State              string    `json:"state"`
}

func (r *Record) String() string {
	bytea, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

func (r *Record) Count() int64 {
	if r.InstanceCount > 0 {
		return r.InstanceCount
	}

	if r.DBInstanceCount > 0 {
		return r.DBInstanceCount
	}

	if r.CacheNodeCount > 0 {
		return r.CacheNodeCount
	}

	return 0
}

func (list RecordList) Region(region string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].Region != region {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) DBInstanceClass(class string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].DBInstanceClass != class {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) CacheNodeType(tipe string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].CacheNodeType != tipe {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) InstanceType(tipe string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].InstanceType != tipe {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) Duration(duration int64) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].Duration != duration {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) LeaseContractLength(length string) RecordList {
	ret := RecordList{}

	duration := 31536000
	if length == "3yr" {
		duration = 94608000
	}

	for i := range list {
		if list[i].Duration != int64(duration) {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) OfferingType(tipe string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].OfferingType != tipe {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) OfferingClass(class string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].OfferingClass != class {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) MultiAZ(multiaz bool) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].MultiAZ != multiaz {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) Active() RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].State != "active" {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

/*
ProductDescription is reserved instance product platform description.

https://docs.aws.amazon.com/cli/latest/reference/ec2/describe-reserved-instances-offerings.html
https://docs.aws.amazon.com/cli/latest/reference/rds/describe-reserved-db-instances.html
https://docs.aws.amazon.com/cli/latest/reference/elasticache/describe-reserved-cache-nodes.html

Instances that include (Amazon VPC) in the product platform description will only be displayed to EC2-Classic account holders and are for use with Amazon VPC.
Linux/UNIX
Linux/UNIX (Amazon VPC)
SUSE Linux
SUSE Linux (Amazon VPC)
Red Hat Enterprise Linux
Red Hat Enterprise Linux (Amazon VPC)
Windows
Windows (Amazon VPC)
Windows with SQL Server Standard
Windows with SQL Server Standard (Amazon VPC)
Windows with SQL Server Web
Windows with SQL Server Web (Amazon VPC)
Windows with SQL Server Enterprise
Windows with SQL Server Enterprise (Amazon VPC)
*/
func (list RecordList) ProductDescription(desc string) RecordList {
	ret := RecordList{}

	for i := range list {
		// desc == Linux
		if strings.Contains(list[i].ProductDescription, desc) {
			ret = append(ret, list[i])
			continue
		}

		// desc == RHEL
		if strings.Contains(list[i].ProductDescription, "Red Hat Enterprise Linux") && desc == "RHEL" {
			ret = append(ret, list[i])
			continue
		}

		// redis, aurora-mysql
		if list[i].ProductDescription == strings.ToLower(strings.Replace(desc, " ", "-", -1)) {
			ret = append(ret, list[i])
		}
	}
	return ret
}
