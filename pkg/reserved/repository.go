package reserved

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Repository struct {
	Profile  string     `json:"profile"`
	Region   []string   `json:"region"`
	Internal RecordList `json:"internal"`
}

func NewRepository(profile string, region []string) (*Repository, error) {
	repo := &Repository{
		Profile: profile,
		Region:  region,
	}

	os.Setenv("AWS_PROFILE", profile)
	for _, r := range region {
		os.Setenv("AWS_REGION", r)

		{
			client := ec2.New(session.Must(session.NewSession()))
			output, err := client.DescribeReservedInstances(&ec2.DescribeReservedInstancesInput{
				Filters: []*ec2.Filter{
					{Name: aws.String("state"), Values: []*string{aws.String("active")}},
				},
			})
			if err != nil {
				return nil, fmt.Errorf("describe reserved instances: %v", err)
			}

			for _, i := range output.ReservedInstances {
				repo.Internal = append(repo.Internal, &Record{
					Region:             r,
					Duration:           *i.Duration,
					OfferingType:       *i.OfferingType,
					OfferingClass:      *i.OfferingClass,
					ProductDescription: *i.ProductDescription,
					InstanceType:       *i.InstanceType,
					InstanceCount:      *i.InstanceCount,
					Start:              *i.Start,
				})
			}
		}

		{
			client := elasticache.New(session.Must(session.NewSession()))
			var maker *string
			for {
				input := &elasticache.DescribeReservedCacheNodesInput{}
				if maker != nil {
					input.Marker = maker
				}

				output, err := client.DescribeReservedCacheNodes(input)
				if err != nil {
					return nil, fmt.Errorf("describe reserved cachenode: %v", err)
				}

				for _, i := range output.ReservedCacheNodes {
					if *i.State != "active" {
						continue
					}

					repo.Internal = append(repo.Internal, &Record{
						Region:             r,
						Duration:           *i.Duration,
						OfferingType:       *i.OfferingType,
						ProductDescription: *i.ProductDescription,
						CacheNodeType:      *i.CacheNodeType,
						CacheNodeCount:     *i.CacheNodeCount,
						Start:              *i.StartTime,
					})
				}

				if maker == nil {
					break
				}
			}
		}

		{
			client := rds.New(session.Must(session.NewSession()))
			var maker *string
			for {
				input := &rds.DescribeReservedDBInstancesInput{}
				if maker != nil {
					input.Marker = maker
				}

				output, err := client.DescribeReservedDBInstances(input)
				if err != nil {
					return nil, fmt.Errorf("describe reserved db instance: %v", err)
				}

				for _, i := range output.ReservedDBInstances {
					if *i.State != "active" {
						continue
					}

					repo.Internal = append(repo.Internal, &Record{
						Region:             r,
						Duration:           *i.Duration,
						OfferingType:       *i.OfferingType,
						ProductDescription: *i.ProductDescription,
						DBInstanceClass:    *i.DBInstanceClass,
						DBInstanceCount:    *i.DBInstanceCount,
						Start:              *i.StartTime,
						MultiAZ:            *i.MultiAZ,
					})
				}

				if maker == nil {
					break
				}
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
