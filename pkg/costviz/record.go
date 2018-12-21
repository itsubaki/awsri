package costviz

import "encoding/json"

type RecordList []*Record

type Record struct {
	AccountID       string  `json:"account_id"`
	Date            string  `json:"date"`
	ID              string  `json:"id"`
	UsageType       string  `json:"usage_type"`
	OperatingSystem string  `json:"operating_system,omitempty"`
	Engine          string  `json:"engine,omitempty"`
	InstanceHour    float64 `json:"instance_hour"`
	InstanceNum     float64 `json:"instance_num"`
}

func (u *Record) String() string {
	bytea, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}
