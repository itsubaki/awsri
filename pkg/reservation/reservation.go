package reservation

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type Utilization struct {
	AccountID      string  `json:"account_id"`
	Region         string  `json:"region"`
	InstanceType   string  `json:"instance_type"`
	Platform       string  `json:"platform,omitempty"`
	CacheEngine    string  `json:"cache_engine,omitempty"`
	DatabaseEngine string  `json:"database_engine,omitempty"`
	Date           string  `json:"date"`
	Hours          float64 `json:"hours"`
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

// Service Filter are
// Amazon Elastic Compute Cloud - Compute
// Amazon Relational Database Service
// Amazon ElastiCache, Amazon Redshift
// Amazon Elasticsearch Service
type getInputFunc func() (string, string)

func getComputeInput() (string, string) {
	return "Amazon Elastic Compute Cloud - Compute", "PLATFORM"
}

func getCacheInput() (string, string) {
	return "Amazon ElastiCache", "CACHE_ENGINE"
}

func getDatabaseInput() (string, string) {
	return "Amazon Relational Database Service", "DATABASE_ENGINE"
}

var getInputFuncList = []getInputFunc{
	getComputeInput,
	getCacheInput,
	getDatabaseInput,
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

				u := Utilization{
					AccountID:    *g.Attributes["linkedAccount"],
					Region:       *g.Attributes["region"],
					InstanceType: *g.Attributes["instanceType"],
					Date:         date,
					Hours:        hours,
				}

				if g.Attributes["platform"] != nil {
					u.Platform = *g.Attributes["platform"]
				}

				if g.Attributes["cacheEngine"] != nil {
					u.CacheEngine = *g.Attributes["cacheEngine"]
				}

				if g.Attributes["databaseEngine"] != nil {
					u.DatabaseEngine = *g.Attributes["databaseEngine"]
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

func Fetch(start, end string) ([]Utilization, error) {
	out := make([]Utilization, 0)
	for _, f := range getInputFuncList {
		service, dimension := f()
		input := costexplorer.GetReservationCoverageInput{
			Filter: &costexplorer.Expression{
				Dimensions: &costexplorer.DimensionValues{
					Key:    aws.String("SERVICE"),
					Values: []*string{&service},
				},
			},
			Metrics: []*string{aws.String("Hour")},
			GroupBy: []*costexplorer.GroupDefinition{
				{
					Key:  aws.String("LINKED_ACCOUNT"),
					Type: aws.String("DIMENSION"),
				},
				{
					Key:  aws.String("INSTANCE_TYPE"),
					Type: aws.String("DIMENSION"),
				},
				{
					Key:  aws.String("REGION"),
					Type: aws.String("DIMENSION"),
				},
				{
					Key:  aws.String(dimension),
					Type: aws.String("DIMENSION"),
				},
			},
			TimePeriod: &costexplorer.DateInterval{
				Start: &start,
				End:   &end,
			},
		}

		u, err := fetch(input)
		if err != nil {
			return out, fmt.Errorf("fetch: %v", err)
		}

		out = append(out, u...)
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].Platform < out[j].Platform })
	sort.SliceStable(out, func(i, j int) bool { return out[i].CacheEngine < out[j].CacheEngine })
	sort.SliceStable(out, func(i, j int) bool { return out[i].DatabaseEngine < out[j].DatabaseEngine })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Region < out[j].Region })
	sort.SliceStable(out, func(i, j int) bool { return out[i].AccountID < out[j].AccountID })

	return out, nil
}
