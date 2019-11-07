package reservation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/itsubaki/hermes/pkg/usage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type Utilization struct {
	AccountID        string  `json:"account_id"`
	Description      string  `json:"description"`
	Region           string  `json:"region"`
	InstanceType     string  `json:"instance_type"`
	Platform         string  `json:"platform,omitempty"`
	CacheEngine      string  `json:"cache_engine,omitempty"`
	DatabaseEngine   string  `json:"database_engine,omitempty"`
	DeploymentOption string  `json:"deployment_option,omitempty"`
	Date             string  `json:"date"`
	Hours            float64 `json:"hours"`
	Num              float64 `json:"num"`
	Percentage       float64 `json:"percentage"`
}

func (u Utilization) UsageType() string {
	return fmt.Sprintf("%s-%s:%s", region[u.Region], u.Usage(), u.InstanceType)
}

func (u Utilization) PFEngine() string {
	return fmt.Sprintf("%s%s%s", u.Platform, u.CacheEngine, u.DatabaseEngine)
}

func (u Utilization) OSEngine() string {
	return fmt.Sprintf("%s%s%s", usage.OperatingSystem[u.Platform], u.CacheEngine, u.DatabaseEngine)
}

func (u Utilization) Usage() string {
	if len(u.Platform) > 0 {
		return "BoxUsage"
	}

	if len(u.CacheEngine) > 0 {
		return "NodeUsage"
	}

	if len(u.DatabaseEngine) > 0 && u.DeploymentOption == "Single-AZ" {
		return "InstanceUsage"
	}

	if len(u.DatabaseEngine) > 0 && u.DeploymentOption == "Multi-AZ" {
		return "Multi-AZUsage"
	}

	if strings.Contains(u.InstanceType, "ds1") || strings.Contains(u.InstanceType, "ds2") {
		return "Node"
	}

	if strings.Contains(u.InstanceType, "dc1") || strings.Contains(u.InstanceType, "dc2") {
		return "Node"
	}

	panic(fmt.Sprintf("invalid usage=%v", u))
}

func (u Utilization) String() string {
	return u.JSON()
}

func (u Utilization) JSON() string {
	bytes, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func (u Utilization) Pretty() string {
	b, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, b, "", " "); err != nil {
		panic(err)
	}

	return pretty.String()
}

// Service Filter are
// Amazon Elastic Compute Cloud - Compute
// Amazon Relational Database Service
// Amazon ElastiCache
// Amazon Redshift
// Amazon Elasticsearch Service
type fetchInputFunc func() (*costexplorer.Expression, []*costexplorer.GroupDefinition)

func fetchComputeInput() (*costexplorer.Expression, []*costexplorer.GroupDefinition) {
	return &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String("SERVICE"),
				Values: []*string{aws.String("Amazon Elastic Compute Cloud - Compute")},
			},
		}, []*costexplorer.GroupDefinition{
			{
				Key:  aws.String("INSTANCE_TYPE"),
				Type: aws.String("DIMENSION"),
			},
			{
				Key:  aws.String("REGION"),
				Type: aws.String("DIMENSION"),
			},
			{
				Key:  aws.String("PLATFORM"),
				Type: aws.String("DIMENSION"),
			},
		}
}

func fetchCacheInput() (*costexplorer.Expression, []*costexplorer.GroupDefinition) {
	return &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String("SERVICE"),
				Values: []*string{aws.String("Amazon ElastiCache")},
			},
		}, []*costexplorer.GroupDefinition{
			{
				Key:  aws.String("INSTANCE_TYPE"),
				Type: aws.String("DIMENSION"),
			},
			{
				Key:  aws.String("REGION"),
				Type: aws.String("DIMENSION"),
			},
			{
				Key:  aws.String("CACHE_ENGINE"),
				Type: aws.String("DIMENSION"),
			},
		}
}

func fetchDatabaseInput() (*costexplorer.Expression, []*costexplorer.GroupDefinition) {
	return &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String("SERVICE"),
				Values: []*string{aws.String("Amazon Relational Database Service")},
			},
		}, []*costexplorer.GroupDefinition{
			{
				Key:  aws.String("INSTANCE_TYPE"),
				Type: aws.String("DIMENSION"),
			},
			{
				Key:  aws.String("REGION"),
				Type: aws.String("DIMENSION"),
			},
			{
				Key:  aws.String("DATABASE_ENGINE"),
				Type: aws.String("DIMENSION"),
			},
			{
				Key:  aws.String("DEPLOYMENT_OPTION"),
				Type: aws.String("DIMENSION"),
			},
		}
}

func fetchRedshiftInput() (*costexplorer.Expression, []*costexplorer.GroupDefinition) {
	return &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String("SERVICE"),
				Values: []*string{aws.String("Amazon Redshift")},
			},
		}, []*costexplorer.GroupDefinition{
			{
				Key:  aws.String("INSTANCE_TYPE"),
				Type: aws.String("DIMENSION"),
			},
			{
				Key:  aws.String("REGION"),
				Type: aws.String("DIMENSION"),
			},
		}
}

func Fetch(start, end string) ([]Utilization, error) {
	return FetchWith(start, end, []fetchInputFunc{
		fetchComputeInput,
		fetchCacheInput,
		fetchDatabaseInput,
		fetchRedshiftInput,
	})
}

func FetchWith(start, end string, fn []fetchInputFunc) ([]Utilization, error) {
	linked, err := usage.FetchLinkedAccount(start, end)
	if err != nil {
		return nil, fmt.Errorf("get linked account: %v", err)
	}

	out := make([]Utilization, 0)
	for _, f := range fn {
		for _, a := range linked {
			exp, groupby := f()

			and := make([]*costexplorer.Expression, 0)
			and = append(and, &costexplorer.Expression{
				Dimensions: &costexplorer.DimensionValues{
					Key:    aws.String("LINKED_ACCOUNT"),
					Values: []*string{aws.String(a.ID)},
				},
			})
			and = append(and, exp)

			input := costexplorer.GetReservationCoverageInput{
				Metrics: []*string{aws.String("Hour")},
				Filter: &costexplorer.Expression{
					And: and,
				},
				GroupBy: groupby,
				TimePeriod: &costexplorer.DateInterval{
					Start: &start,
					End:   &end,
				},
			}

			u, err := fetch(input)
			if err != nil {
				return out, fmt.Errorf("fetch: %v", err)
			}

			for i := range u {
				u[i].AccountID = a.ID
				u[i].Description = a.Description
			}

			out = append(out, u...)
		}
	}

	return out, nil
}

func fetch(input costexplorer.GetReservationCoverageInput) ([]Utilization, error) {
	out := make([]Utilization, 0)
	c := costexplorer.New(session.Must(session.NewSession()))

	var token *string
	for {
		input.NextPageToken = token

		rc, err := c.GetReservationCoverage(&input)
		if err != nil {
			return out, fmt.Errorf("get reservation coverage: %v", err)
		}

		for _, t := range rc.CoveragesByTime {
			for _, g := range t.Groups {
				if *g.Coverage.CoverageHours.ReservedHours == "0" {
					continue
				}

				index := strings.LastIndex(*input.TimePeriod.Start, "-")
				date := (*input.TimePeriod.Start)[:index]

				hours, err := strconv.ParseFloat(*g.Coverage.CoverageHours.ReservedHours, 64)
				if err != nil {
					return out, fmt.Errorf("parse float reserved hours: %v", err)
				}

				month := strings.Split(date, "-")[1]
				num := hours / float64(24*usage.Days[month])

				per, err := strconv.ParseFloat(*g.Coverage.CoverageHours.CoverageHoursPercentage, 64)
				if err != nil {
					return out, fmt.Errorf("parse float reserved hours percentage: %v", err)
				}

				u := Utilization{
					Region:       *g.Attributes["region"],
					InstanceType: *g.Attributes["instanceType"],
					Date:         date,
					Hours:        hours,
					Num:          num,
					Percentage:   per,
				}

				if g.Attributes["platform"] != nil {
					u.Platform = *g.Attributes["platform"]
				}

				if g.Attributes["cacheEngine"] != nil {
					u.CacheEngine = *g.Attributes["cacheEngine"]
				}

				if g.Attributes["databaseEngine"] != nil {
					u.DatabaseEngine = *g.Attributes["databaseEngine"]
					u.DeploymentOption = *g.Attributes["deploymentOption"]
				}

				out = append(out, u)
			}
		}

		if rc.NextPageToken == nil {
			break
		}
		token = rc.NextPageToken
	}

	return out, nil
}

func Sort(u []Utilization) {
	sort.SliceStable(u, func(i, j int) bool { return u[i].DeploymentOption < u[j].DeploymentOption })
	sort.SliceStable(u, func(i, j int) bool { return u[i].DatabaseEngine < u[j].DatabaseEngine })
	sort.SliceStable(u, func(i, j int) bool { return u[i].CacheEngine < u[j].CacheEngine })
	sort.SliceStable(u, func(i, j int) bool { return u[i].Platform < u[j].Platform })
	sort.SliceStable(u, func(i, j int) bool { return u[i].Region < u[j].Region })
	sort.SliceStable(u, func(i, j int) bool { return u[i].AccountID < u[j].AccountID })
}
