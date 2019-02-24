package reserved

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/rds"
)

type FetchReservedRecordList func(ses *session.Session, region string) (RecordList, error)

func NewFetchReservedRecordList() []FetchReservedRecordList {
	return []FetchReservedRecordList{
		FetchReservedComputeRecordList,
		FetchReservedCacheRecordList,
		FetchReservedDatabaseRecordList}
}

func FetchReservedComputeRecordList(ses *session.Session, region string) (RecordList, error) {
	res, err := ec2.New(ses).DescribeReservedInstances(&ec2.DescribeReservedInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("describe reserved instances (region=%s): %v", region, err)
	}

	out := RecordList{}
	for _, i := range res.ReservedInstances {
		out = append(out, &Record{
			Region:             region,
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

	return out, nil
}

func FetchReservedCacheRecordList(ses *session.Session, region string) (RecordList, error) {
	out := RecordList{}

	client := elasticache.New(ses)
	var maker *string
	for {
		input := &elasticache.DescribeReservedCacheNodesInput{}
		if maker != nil {
			input.Marker = maker
		}

		res, err := client.DescribeReservedCacheNodes(input)
		if err != nil {
			return nil, fmt.Errorf("describe reserved cache node (region=%s): %v", region, err)
		}

		for _, i := range res.ReservedCacheNodes {
			out = append(out, &Record{
				Region:             region,
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

		if res.Marker == nil {
			break
		}
		maker = res.Marker
	}

	return out, nil
}

func FetchReservedDatabaseRecordList(ses *session.Session, region string) (RecordList, error) {
	out := RecordList{}

	client := rds.New(ses)
	var maker *string
	for {
		input := &rds.DescribeReservedDBInstancesInput{}
		if maker != nil {
			input.Marker = maker
		}

		res, err := client.DescribeReservedDBInstances(input)
		if err != nil {
			return nil, fmt.Errorf("describe reserved db instance (region=%s): %v", region, err)
		}

		for _, i := range res.ReservedDBInstances {
			out = append(out, &Record{
				Region:             region,
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

		if res.Marker == nil {
			break
		}
		maker = res.Marker
	}

	return out, nil
}
