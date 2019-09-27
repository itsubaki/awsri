package usage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/itsubaki/hermes/pkg/usage"
)

func Hash(q usage.Quantity) string {
	tmp := usage.Quantity{
		AccountID:      q.AccountID,
		Description:    q.Description,
		Region:         q.Region,
		UsageType:      q.UsageType,
		Platform:       q.Platform,
		DatabaseEngine: q.DatabaseEngine,
		CacheEngine:    q.CacheEngine,
	}

	val, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}

	sha := sha256.Sum256(val)
	hash := hex.EncodeToString(sha[:])
	return hash
}
