package costexp

import "encoding/json"

type RecordList []*Record

type Record struct {
	AccountID    string  `json:"account_id"`
	Date         string  `json:"date"`
	UsageType    string  `json:"usage_type"`
	Platform     string  `json:"platform,omitempty"`
	Engine       string  `json:"engine,omitempty"`
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
