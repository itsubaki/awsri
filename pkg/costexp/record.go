package costexp

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"sort"
)

type RecordList []*Record

func (list RecordList) JSON() string {
	bytea, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

func (list RecordList) Pretty() string {
	bytea, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	if err := json.Indent(&out, bytea, "", " "); err != nil {
		panic(err)
	}

	return string(out.Bytes())
}

type Record struct {
	AccountID      string  `json:"account_id"`
	Description    string  `json:"description"`
	Region         string  `json:"region"`
	UsageType      string  `json:"usage_type"`
	Platform       string  `json:"platform,omitempty"`        // compute
	DatabaseEngine string  `json:"database_engine,omitempty"` // database
	CacheEngine    string  `json:"cache_engine,omitempty"`    // cache
	Date           string  `json:"date"`
	InstanceHour   float64 `json:"instance_hour"`
	InstanceNum    float64 `json:"instance_num"`
}

func (u *Record) Hash() string {
	bytea, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	sha := sha256.Sum256(bytea)
	hash := hex.EncodeToString(sha[:])
	return hash
}

func (u *Record) JSON() string {
	bytea, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}

func (list RecordList) Unique(fieldname string) []string {
	uniq := make(map[string]bool)
	for i := range list {
		ref := reflect.ValueOf(*list[i]).FieldByName(fieldname)
		val := ref.Interface().(string)
		if len(val) > 0 {
			uniq[val] = true
		}
	}

	out := []string{}
	for k := range uniq {
		out = append(out, k)
	}

	return out
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

func (list RecordList) CacheEngine(engine string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].CacheEngine != engine {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) DatabaseEngine(engine string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].DatabaseEngine != engine {
			continue
		}
		ret = append(ret, list[i])
	}

	return ret
}

func (list RecordList) Sort() RecordList {
	ret := append(RecordList{}, list...)

	sort.SliceStable(ret, func(i, j int) bool { return ret[i].UsageType < ret[j].UsageType })
	sort.SliceStable(ret, func(i, j int) bool { return ret[i].Platform < ret[j].Platform })
	sort.SliceStable(ret, func(i, j int) bool { return ret[i].CacheEngine < ret[j].CacheEngine })
	sort.SliceStable(ret, func(i, j int) bool { return ret[i].DatabaseEngine < ret[j].DatabaseEngine })
	sort.SliceStable(ret, func(i, j int) bool { return ret[i].AccountID < ret[j].AccountID })

	return ret
}
