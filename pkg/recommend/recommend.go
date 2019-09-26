package recommend

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/itsubaki/hermes/pkg/usage"
)

func Recommend(quantity []usage.Quantity) ([]usage.Quantity, error) {
	tmp := make(map[string]map[string]map[string]float64)
	for _, q := range quantity {
		if _, ok := tmp[q.UsageType]; !ok {
			tmp[q.UsageType] = make(map[string]map[string]float64)
		}

		engine := fmt.Sprintf("%s%s%s", q.Platform, q.CacheEngine, q.DatabaseEngine)
		if _, ok := tmp[q.UsageType][engine]; !ok {
			tmp[q.UsageType][engine] = make(map[string]float64)
		}

		tmp[q.UsageType][engine][q.Date] = tmp[q.UsageType][engine][q.Date] + q.InstanceNum
	}

	for k, v := range tmp {
		fmt.Printf("%v: %v\n", k, v)
	}

	return []usage.Quantity{}, nil
}

func Hash(str string) string {
	val, err := json.Marshal(str)
	if err != nil {
		panic(err)
	}

	sha := sha256.Sum256(val)
	hash := hex.EncodeToString(sha[:])
	return hash
}
