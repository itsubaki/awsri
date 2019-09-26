package usage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	dir := c.GlobalString("dir")
	format := c.String("format")

	quantity := make([]usage.Quantity, 0)
	date := usage.Last12Months()
	for i := range date {
		file := fmt.Sprintf("%s/usage/%s.out", dir, date[i].YYYYMM())
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("file not found: %v", file)
			os.Exit(1)
		}

		read, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("read %s: %v", file, err)
			os.Exit(1)
		}

		var q []usage.Quantity
		if err := json.Unmarshal(read, &q); err != nil {
			fmt.Printf("unmarshal: %v", err)
			os.Exit(1)
		}

		quantity = append(quantity, q...)
	}

	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].AccountID < quantity[j].AccountID })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].Region < quantity[j].Region })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].UsageType < quantity[j].UsageType })
	sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].Date < quantity[j].Date })

	if format == "json" {
		bytes, err := json.Marshal(quantity)
		if err != nil {
			fmt.Printf("marshal: %v", err)
			os.Exit(1)
		}

		fmt.Println(string(bytes))
		return
	}

	if format == "csv" {
		tmp := make(map[string][]usage.Quantity)
		for _, q := range quantity {
			tmp[Hash(q)] = append(tmp[Hash(q)], q)
		}

		fmt.Printf("accountID, description, region, type, os/engine, ")
		for i := range date {
			fmt.Printf("%s, ", date[i].YYYYMM())
		}
		fmt.Println()

		for _, v := range tmp {
			fmt.Printf("%s, %s, ", v[0].AccountID, v[0].Description)
			fmt.Printf("%s, %s, ", v[0].Region, v[0].UsageType)
			fmt.Printf("%s, ", fmt.Sprintf("%s%s%s", v[0].Platform, v[0].CacheEngine, v[0].DatabaseEngine))

			for i := range date {
				found := false
				for _, q := range v {
					if date[i].YYYYMM() == q.Date {
						fmt.Printf("%.3f, ", q.InstanceNum)
						found = true
						break
					}
				}

				if !found {
					fmt.Printf("0.0, ")
				}
			}
			fmt.Println()
		}

		return
	}
}

type Quantity struct {
	Hash        string
	Quantity    usage.Quantity
	InstanceNum []float64
}

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
