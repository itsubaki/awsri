package main

import (
	"fmt"
	"testing"
)

func TestGenerateInput(t *testing.T) {
	input := &ForecastList{
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
