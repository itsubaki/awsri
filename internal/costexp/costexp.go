package costexp

type OutputCostExp struct {
}

type UsageQuantityList []*UsageQuantity

type UsageQuantity struct {
	AccountID       string  `json:"account_id"`
	Date            string  `json:"date"`
	UsageType       string  `json:"usage_type"`
	OperatingSystem string  `json:"operating_system,omitempty"`
	Engine          string  `json:"engine,omitempty"`
	InstanceHour    float64 `json:"instance_hour"`
	InstanceNum     float64 `json:"instance_num"`
}
