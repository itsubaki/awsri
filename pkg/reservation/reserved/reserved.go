package reserved

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/redshift"
)

type Reserved struct {
	ReservedID         string    `json:"reserved_id"`
	Region             string    `json:"region"`
	Duration           int64     `json:"duration"`
	OfferingType       string    `json:"offering_type"`
	OfferingClass      string    `json:"offering_class,omitempty"`
	ProductDescription string    `json:"product_description,omitempty"`
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

func (r Reserved) String() string {
	return r.JSON()
}

func (r Reserved) JSON() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

type fetchFunc func(region []string) ([]Reserved, error)

var fetchFuncList = []fetchFunc{
	fetchCompute,
	fetchCache,
	fetchDatabase,
	fetchRedshift,
}

func Fetch(region []string) ([]Reserved, error) {
	return FetchWith(region, []fetchFunc{
		fetchCompute,
		fetchCache,
		fetchDatabase,
		fetchRedshift,
	})
}

func FetchWith(region []string, fn []fetchFunc) ([]Reserved, error) {
	out := make([]Reserved, 0)
	for _, f := range fn {
		list, err := f(region)
		if err != nil {
			return out, fmt.Errorf("fetch reserved description: %v", err)
		}
		out = append(out, list...)
	}

	return out, nil
}

func fetchCompute(region []string) ([]Reserved, error) {
	out := make([]Reserved, 0)
	for _, r := range region {
		ses, err := session.NewSession(
			&aws.Config{
				Region: aws.String(r),
			},
		)
		if err != nil {
			return out, fmt.Errorf("new session (region=%s): %v", r, err)
		}

		desc, err := ec2.New(ses).DescribeReservedInstances(&ec2.DescribeReservedInstancesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("state"), Values: []*string{aws.String("active")}},
			},
		})
		if err != nil {
			return out, fmt.Errorf("describe reserved instances (region=%s): %v", r, err)
		}

		for _, i := range desc.ReservedInstances {
			out = append(out, Reserved{
				Region:             r,
				ReservedID:         *i.ReservedInstancesId,
				Duration:           *i.Duration,
				OfferingType:       *i.OfferingType,
				OfferingClass:      *i.OfferingClass,
				ProductDescription: *i.ProductDescription,
				InstanceType:       *i.InstanceType,
				InstanceCount:      *i.InstanceCount,
				Start:              *i.Start,
				State:              *i.State,
			})
		}
	}

	return out, nil
}

func fetchCache(region []string) ([]Reserved, error) {
	out := make([]Reserved, 0)
	for _, r := range region {
		ses, err := session.NewSession(
			&aws.Config{
				Region: aws.String(r),
			},
		)
		if err != nil {
			return out, fmt.Errorf("new session (region=%s): %v", r, err)
		}
		client := elasticache.New(ses)
		var maker *string
		for {
			input := elasticache.DescribeReservedCacheNodesInput{}
			if maker != nil {
				input.Marker = maker
			}

			desc, err := client.DescribeReservedCacheNodes(&input)
			if err != nil {
				return out, fmt.Errorf("describe reserved cache node (region=%s): %v", r, err)
			}

			for _, i := range desc.ReservedCacheNodes {
				out = append(out, Reserved{
					Region:             r,
					ReservedID:         *i.ReservedCacheNodeId,
					Duration:           *i.Duration,
					OfferingType:       *i.OfferingType,
					ProductDescription: *i.ProductDescription,
					CacheNodeType:      *i.CacheNodeType,
					CacheNodeCount:     *i.CacheNodeCount,
					Start:              *i.StartTime,
					State:              *i.State,
				})
			}

			if desc.Marker == nil {
				break
			}
			maker = desc.Marker
		}
	}

	return out, nil
}

func fetchDatabase(region []string) ([]Reserved, error) {
	out := make([]Reserved, 0)
	for _, r := range region {
		ses, err := session.NewSession(
			&aws.Config{
				Region: aws.String(r),
			},
		)
		if err != nil {
			return out, fmt.Errorf("new session (region=%s): %v", r, err)
		}

		client := rds.New(ses)
		var maker *string
		for {
			input := rds.DescribeReservedDBInstancesInput{}
			if maker != nil {
				input.Marker = maker
			}

			desc, err := client.DescribeReservedDBInstances(&input)
			if err != nil {
				return out, fmt.Errorf("describe reserved db instance (region=%s): %v", r, err)
			}

			for _, i := range desc.ReservedDBInstances {
				out = append(out, Reserved{
					Region:             r,
					ReservedID:         *i.ReservedDBInstanceId,
					Duration:           *i.Duration,
					OfferingType:       *i.OfferingType,
					ProductDescription: *i.ProductDescription,
					DBInstanceClass:    *i.DBInstanceClass,
					DBInstanceCount:    *i.DBInstanceCount,
					Start:              *i.StartTime,
					MultiAZ:            *i.MultiAZ,
					State:              *i.State,
				})
			}

			if desc.Marker == nil {
				break
			}
			maker = desc.Marker
		}
	}

	return out, nil
}

func fetchRedshift(region []string) ([]Reserved, error) {
	out := make([]Reserved, 0)
	for _, r := range region {
		ses, err := session.NewSession(
			&aws.Config{
				Region: aws.String(r),
			},
		)
		if err != nil {
			return out, fmt.Errorf("new session (region=%s): %v", r, err)
		}

		client := redshift.New(ses)
		var maker *string
		for {
			input := redshift.DescribeReservedNodesInput{}
			if maker != nil {
				input.Marker = maker
			}

			desc, err := client.DescribeReservedNodes(&input)
			if err != nil {
				return out, fmt.Errorf("describe reserved db instance (region=%s): %v", r, err)
			}

			for _, i := range desc.ReservedNodes {
				out = append(out, Reserved{
					Region:         r,
					ReservedID:     *i.ReservedNodeId,
					Duration:       *i.Duration,
					OfferingType:   *i.OfferingType,
					CacheNodeType:  *i.NodeType,
					CacheNodeCount: *i.NodeCount,
					Start:          *i.StartTime,
					State:          *i.State,
				})
			}

			if desc.Marker == nil {
				break
			}
			maker = desc.Marker
		}

	}

	return out, nil
}
