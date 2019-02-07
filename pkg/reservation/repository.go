package reservation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Repository struct {
	Region   []string   `json:"region"`
	Internal RecordList `json:"internal"`
}

func New(region []string) (*Repository, error) {
	repo := NewRepository(region)
	return repo, repo.Fetch()
}

func NewRepository(region []string) *Repository {
	return &Repository{
		Region: region,
	}
}

func (repo *Repository) Fetch() error {
	return repo.FetchWithClient(http.DefaultClient)
}

func (repo *Repository) fetchEC2WithClient(client *http.Client) error {
	for _, r := range repo.Region {
		ses, err := session.NewSession(
			&aws.Config{
				Region:     aws.String(r),
				HTTPClient: client,
			},
		)
		if err != nil {
			return fmt.Errorf("new session (region=%s): %v", r, err)
		}

		output, err := ec2.New(ses).DescribeReservedInstances(&ec2.DescribeReservedInstancesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("state"), Values: []*string{aws.String("active")}},
			},
		})
		if err != nil {
			return fmt.Errorf("describe reserved instances (region=%s): %v", r, err)
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
				State:              *i.State,
			})
		}
	}

	return nil
}

func (repo *Repository) fetchCacheWithClient(client *http.Client) error {
	for _, r := range repo.Region {
		ses, err := session.NewSession(
			&aws.Config{
				Region:     aws.String(r),
				HTTPClient: client,
			},
		)
		if err != nil {
			return fmt.Errorf("new session (region=%s): %v", r, err)
		}

		client := elasticache.New(ses)
		var maker *string
		for {
			input := &elasticache.DescribeReservedCacheNodesInput{}
			if maker != nil {
				input.Marker = maker
			}

			output, err := client.DescribeReservedCacheNodes(input)
			if err != nil {
				return fmt.Errorf("describe reserved cache node (region=%s): %v", r, err)
			}

			for _, i := range output.ReservedCacheNodes {
				repo.Internal = append(repo.Internal, &Record{
					Region:             r,
					Duration:           *i.Duration,
					OfferingType:       *i.OfferingType,
					ProductDescription: *i.ProductDescription,
					CacheNodeType:      *i.CacheNodeType,
					CacheNodeCount:     *i.CacheNodeCount,
					Start:              *i.StartTime,
					State:              *i.State,
				})
			}

			if output.Marker == nil {
				break
			}
			maker = output.Marker
		}
	}

	return nil
}

func (repo *Repository) fetchRDSWithClient(client *http.Client) error {
	for _, r := range repo.Region {
		ses, err := session.NewSession(
			&aws.Config{
				Region:     aws.String(r),
				HTTPClient: client,
			},
		)
		if err != nil {
			return fmt.Errorf("new session (region=%s): %v", r, err)
		}

		client := rds.New(ses)
		var maker *string
		for {
			input := &rds.DescribeReservedDBInstancesInput{}
			if maker != nil {
				input.Marker = maker
			}

			output, err := client.DescribeReservedDBInstances(input)
			if err != nil {
				return fmt.Errorf("describe reserved db instance (region=%s): %v", r, err)
			}

			for _, i := range output.ReservedDBInstances {
				repo.Internal = append(repo.Internal, &Record{
					Region:             r,
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

			if output.Marker == nil {
				break
			}
			maker = output.Marker
		}
	}

	return nil
}

func (repo *Repository) FetchWithClient(client *http.Client) error {
	if err := repo.fetchEC2WithClient(client); err != nil {
		return err
	}

	if err := repo.fetchCacheWithClient(client); err != nil {
		return err
	}

	if err := repo.fetchRDSWithClient(client); err != nil {
		return err
	}

	return nil
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

func (repo *Repository) Write(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}

	bytes, err := repo.Serialize()
	if err != nil {
		return fmt.Errorf("serialize: %v", err)
	}

	if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}

func (repo *Repository) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(repo)
	if err != nil {
		return []byte{}, fmt.Errorf("marshal: %v", err)
	}

	return bytes, nil
}

func (repo *Repository) Deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, repo); err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	return nil
}

func (repo *Repository) SelectAll() RecordList {
	return repo.Internal
}

func (repo *Repository) FindByInstanceType(tipe string) RecordList {
	out := RecordList{}
	for i := range repo.Internal {
		if repo.Internal[i].InstanceType == tipe {
			out = append(out, repo.Internal[i])
		}
	}

	return out
}
