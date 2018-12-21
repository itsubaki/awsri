package costexp

import (
	"encoding/json"
	"sort"
)

type RecordList []*Record

type Record struct {
	AccountID    string  `json:"account_id"`
	UsageType    string  `json:"usage_type"`
	Platform     string  `json:"platform,omitempty"`
	Engine       string  `json:"engine,omitempty"`
	Date         string  `json:"date"`
	InstanceHour float64 `json:"instance_hour"`
	InstanceNum  float64 `json:"instance_num"`
}

func (u *Record) String() string {
	bytea, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

func (r RecordList) Sort() RecordList {
	list := append(RecordList{}, r...)

	sort.SliceStable(list, func(i, j int) bool { return list[i].UsageType < list[j].UsageType })
	sort.SliceStable(list, func(i, j int) bool { return list[i].Platform < list[j].Platform })
	sort.SliceStable(list, func(i, j int) bool { return list[i].Engine < list[j].Engine })
	sort.SliceStable(list, func(i, j int) bool { return list[i].AccountID < list[j].AccountID })

	return list
}
