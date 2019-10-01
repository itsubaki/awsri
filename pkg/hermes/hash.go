package hermes

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func Hash(str string) string {
	val, err := json.Marshal(str)
	if err != nil {
		panic(err)
	}

	sha := sha256.Sum256(val)
	hash := hex.EncodeToString(sha[:])

	return hash
}
