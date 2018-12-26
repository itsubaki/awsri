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

func (list RecordList) InstanceNumAvg() float64 {
	if len(list) == 0 {
		return 0
	}

	sum := 0.0
	for i := range list {
		sum = sum + list[i].InstanceNum
	}

	return sum / float64(len(list))
}

func (list RecordList) InstanceHourAvg() float64 {
	if len(list) == 0 {
		return 0
	}

	sum := 0.0
	for i := range list {
		sum = sum + list[i].InstanceHour
	}

	return sum / float64(len(list))
}

func (list RecordList) AccountID(accountID string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].AccountID != accountID {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) UsageType(usageType string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].UsageType != usageType {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) Platform(platform string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].Platform != platform {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) Date(date string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].Date != date {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) Engine(engine string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].Engine != engine {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (r RecordList) Sort() RecordList {
	list := append(RecordList{}, r...)

	sort.SliceStable(list, func(i, j int) bool { return list[i].UsageType < list[j].UsageType })
	sort.SliceStable(list, func(i, j int) bool { return list[i].Platform < list[j].Platform })
	sort.SliceStable(list, func(i, j int) bool { return list[i].Engine < list[j].Engine })
	sort.SliceStable(list, func(i, j int) bool { return list[i].AccountID < list[j].AccountID })

	return list
}
