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

func (list RecordList) ProductDescription(desc string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].ProductDescription != desc {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) ContainsProductDescription(desc string) RecordList {
	ret := RecordList{}

	for i := range list {
		if !strings.Contains(list[i].ProductDescription, desc) {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}
