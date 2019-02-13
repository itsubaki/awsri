package cmd

import (
	"fmt"
	"testing"
)

func TestGenerateInput(t *testing.T) {
	input := &Input{
		Forecast: ForecastList{
			{
				AccountID: "012345678901",
				Alias:     "projectA",
				Region:    "ap-northeast-1",
				UsageType: "APN1-BoxUsage:c4.2xlarge",
				Platform:  "Linux/Unix",
				InstanceNum: InstanceNumList{
					{Date: "2019-01", InstanceNum: 200},
					{Date: "2019-02", InstanceNum: 150},
					{Date: "2019-03", InstanceNum: 80},
					{Date: "2019-04", InstanceNum: 80},
					{Date: "2019-05", InstanceNum: 150},
					{Date: "2019-06", InstanceNum: 80},
					{Date: "2019-07", InstanceNum: 80},
					{Date: "2019-08", InstanceNum: 150},
					{Date: "2019-09", InstanceNum: 80},
					{Date: "2019-10", InstanceNum: 80},
					{Date: "2019-11", InstanceNum: 80},
					{Date: "2019-12", InstanceNum: 150},
				},
			},
			{
				AccountID:      "012345678901",
				Alias:          "projectA",
				Region:         "ap-northeast-1",
				UsageType:      "APN1-InstanceUsage:db.r3.xlarge",
				DatabaseEngine: "Aurora MySQL",
				InstanceNum: InstanceNumList{
					{Date: "2019-01", InstanceNum: 100},
					{Date: "2019-02", InstanceNum: 100},
					{Date: "2019-03", InstanceNum: 100},
					{Date: "2019-04", InstanceNum: 100},
					{Date: "2019-05", InstanceNum: 100},
					{Date: "2019-06", InstanceNum: 100},
					{Date: "2019-07", InstanceNum: 100},
					{Date: "2019-08", InstanceNum: 100},
					{Date: "2019-09", InstanceNum: 100},
					{Date: "2019-10", InstanceNum: 100},
					{Date: "2019-11", InstanceNum: 100},
					{Date: "2019-12", InstanceNum: 100},
				},
			},
			{
				AccountID:   "012345678901",
				Alias:       "projectA",
				Region:      "ap-northeast-1",
				UsageType:   "APN1-NodeUsage:cache.r3.4xlarge",
				CacheEngine: "Redis",
				InstanceNum: InstanceNumList{
					{Date: "2019-01", InstanceNum: 100},
					{Date: "2019-02", InstanceNum: 100},
					{Date: "2019-03", InstanceNum: 100},
					{Date: "2019-04", InstanceNum: 100},
					{Date: "2019-05", InstanceNum: 100},
					{Date: "2019-06", InstanceNum: 100},
					{Date: "2019-07", InstanceNum: 100},
					{Date: "2019-08", InstanceNum: 100},
					{Date: "2019-09", InstanceNum: 100},
					{Date: "2019-10", InstanceNum: 100},
					{Date: "2019-11", InstanceNum: 100},
					{Date: "2019-12", InstanceNum: 100},
				},
			},
			{
				AccountID: "987654321098",
				Alias:     "projectB",
				Region:    "ap-northeast-1",
				UsageType: "APN1-BoxUsage:c4.2xlarge",
				Platform:  "Linux/Unix",
				InstanceNum: InstanceNumList{
					{Date: "2019-01", InstanceNum: 100},
					{Date: "2019-02", InstanceNum: 50},
					{Date: "2019-03", InstanceNum: 50},
					{Date: "2019-04", InstanceNum: 50},
					{Date: "2019-05", InstanceNum: 50},
					{Date: "2019-06", InstanceNum: 50},
					{Date: "2019-07", InstanceNum: 50},
					{Date: "2019-08", InstanceNum: 100},
					{Date: "2019-09", InstanceNum: 50},
					{Date: "2019-10", InstanceNum: 50},
					{Date: "2019-11", InstanceNum: 50},
					{Date: "2019-12", InstanceNum: 80},
				},
			},
		},
	}

	fmt.Println(input.JSON())
}
