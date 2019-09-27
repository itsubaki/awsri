package recommend

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/itsubaki/hermes/pkg/usage"
)

func Recommend(quantity []usage.Quantity) ([]usage.Quantity, error) {
	merged := make(map[string]usage.Quantity)
	for _, q := range quantity {
		hash := Hash(fmt.Sprintf("%s%s%s%s%s", q.UsageType, q.Platform, q.CacheEngine, q.DatabaseEngine, q.Date))
		merged[hash] = usage.Quantity{
			Region:         q.Region,
			UsageType:      q.UsageType,
			Platform:       q.Platform,
			DatabaseEngine: q.DatabaseEngine,
			CacheEngine:    q.CacheEngine,
			Date:           q.Date,
			InstanceHour:   merged[hash].InstanceHour + q.InstanceHour,
			InstanceNum:    merged[hash].InstanceNum + q.InstanceNum,
		}
	}

	for _, v := range merged {
		fmt.Printf("%#v\n", v)
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
