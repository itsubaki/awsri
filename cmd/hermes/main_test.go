package main

import (
	"fmt"
	"testing"
)

func TestGenerateInput(t *testing.T) {
	input := &Input{
		Forecast: []*Forecast{
			{
				AccountID: "012345678901",
				Alias:     "example",
				Region:    "ap-northeast-1",
				UsageType: "APN1-BoxUsage:c4.2xlarge",
				Platform:  "Linux/Unix",
				InstanceNum: []InstanceNum{
					{
						Date:        "2019-01",
						InstanceNum: 100,
					},
					{
						Date:        "2018-02",
						InstanceNum: 100,
					},
				},
			},
		},
	}

	fmt.Println(input.JSON())
}

func TestCache(t *testing.T) {
	input := &Input{
		Forecast: []*Forecast{
			{
				AccountID:   "012345678901",
				Alias:       "example",
				Region:      "ap-northeast-1",
				UsageType:   "APN1-NodeUsage:cache.r3.4xlarge",
				CacheEngine: "Redis",
				InstanceNum: []InstanceNum{
					{
						Date:        "2019-01",
						InstanceNum: 100,
					},
					{
						Date:        "2019-02",
						InstanceNum: 100,
					},
					{
						Date:        "2019-03",
						InstanceNum: 100,
					},
					{
						Date:        "2019-04",
						InstanceNum: 100,
					},
					{
						Date:        "2019-05",
						InstanceNum: 100,
					},
					{
						Date:        "2019-06",
						InstanceNum: 100,
					},
					{
						Date:        "2019-07",
						InstanceNum: 100,
					},
					{
						Date:        "2019-08",
						InstanceNum: 100,
					},
					{
						Date:        "2019-09",
						InstanceNum: 100,
					},
					{
						Date:        "2019-10",
						InstanceNum: 100,
					},
					{
						Date:        "2019-11",
						InstanceNum: 100,
					},
					{
						Date:        "2019-12",
						InstanceNum: 100,
					},
				},
			},
		},
	}

	merged := Merge(input.Forecast)
	fmt.Println(merged[0])

	fmt.Println(Recommended(merged))

}
