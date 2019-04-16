package billing

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
)

type Record struct {
	AccountID        string `json:"account_id"`
	Description      string `json:"description"`
	Date             string `json:"date"`
	AmortizedCost    string `json:"amortized_cost"`
	BlendedCost      string `json:"blended_cost"`
	UnblendedCost    string `json:"unblended_cost"`
	NetAmortizedCost string `json:"net_amortized_cost"`
	NetUnblendedCost string `json:"net_unblended_cost"`
}

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

func (list RecordList) Description(description string) RecordList {
	ret := RecordList{}

	for i := range list {
		if list[i].Description != description {
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

func (list RecordList) Sort() {
	sort.SliceStable(list, func(i, j int) bool { return list[i].AccountID < list[j].AccountID })
	sort.SliceStable(list, func(i, j int) bool { return list[i].Date < list[j].Date })
}
