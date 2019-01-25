package reserved

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/itsubaki/hermes/pkg/awsprice"
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

func (r *Record) Price(repo *awsprice.Repository) (*awsprice.Record, error) {
	yr := "1yr"
	if r.Duration == 94608000 {
		yr = "3yr"
	}

	os := "Linux"
	if strings.Contains(r.ProductDescription, "Windows") {
		os = "Windows"
	}

	rs := repo.FindByInstanceType(r.InstanceType).
		OfferingClass(r.OfferingClass).
		PurchaseOption(r.OfferingType).
		OperatingSystem(os).
		LeaseContractLength(yr).
		Region(r.Region)

	if len(rs) != 1 {
		return nil, fmt.Errorf("invalid query to awsprice repository")
	}

	return rs[0], nil
}
