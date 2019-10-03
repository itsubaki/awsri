package usage

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/itsubaki/hermes/pkg/hermes"
	"github.com/itsubaki/hermes/pkg/pricing"

	"github.com/itsubaki/hermes/pkg/usage"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
	region := c.StringSlice("region")
	dir := c.GlobalString("dir")
	format := c.String("format")
	normalize := c.Bool("normalize")
	merge := c.Bool("merge")
	overall := c.Bool("overall")
	monthly := c.Bool("monthly")

	date := usage.Last12Months()
	quantity, err := usage.Deserialize(dir, date)
	if err != nil {
		fmt.Errorf("deserialize usage: %v", err)
		os.Exit(1)
	}

	if normalize {
		plist, err := pricing.Deserialize("/var/tmp/hermes", region)
		if err != nil {
			fmt.Errorf("desirialize pricing: %v", err)
		}

		family := pricing.Family(plist)
		mini := pricing.Minimum(family, plist)

		quantity = hermes.Normalize(quantity, mini)
	}

	if merge && overall {
		quantity = usage.MergeOverall(quantity)
	}

	if merge && !overall {
		quantity = usage.Merge(quantity)
	}

	if format == "json" && !monthly {
		sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].Date < quantity[j].Date })
		sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].DatabaseEngine < quantity[j].DatabaseEngine })
		sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].CacheEngine < quantity[j].CacheEngine })
		sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].Platform < quantity[j].Platform })
		sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].UsageType < quantity[j].UsageType })
		sort.SliceStable(quantity, func(i, j int) bool { return quantity[i].AccountID < quantity[j].AccountID })

		for _, q := range quantity {
			bytes, err := json.Marshal(q)
			if err != nil {
				fmt.Printf("marshal: %v", err)
				os.Exit(1)
			}

			fmt.Println(string(bytes))
		}

		return
	}

	if format == "json" && monthly {
		mq := usage.Monthly(quantity)
		for _, q := range mq {
			bytes, err := json.Marshal(q)
			if err != nil {
				fmt.Printf("marshal: %v", err)
				os.Exit(1)
			}

			fmt.Println(string(bytes))
		}
	}

	//if format == "csv" {
	//	fmt.Printf("accountID, description, region, usage_type, os/engine, ")
	//	for i := range date {
	//		fmt.Printf("%s, ", date[i].YYYYMM())
	//	}
	//	fmt.Println()
	//
	//	quantity = usage.MergeGroupBy(quantity)
	//	month := usage.Monthly(quantity)
	//	for _, v := range month {
	//		fmt.Printf("%s, %s, ", v[0].AccountID, v[0].Description)
	//		fmt.Printf("%s, %s, ", v[0].Region, v[0].UsageType)
	//		fmt.Printf("%s, ", fmt.Sprintf("%s%s%s", v[0].Platform, v[0].CacheEngine, v[0].DatabaseEngine))
	//
	//		for _, d := range date {
	//			found := false
	//			for _, q := range v {
	//				if d.YYYYMM() == q.Date {
	//					fmt.Printf("%.3f, ", q.InstanceNum)
	//					found = true
	//					break
	//				}
	//			}
	//
	//			if !found {
	//				fmt.Printf("0.0, ")
	//			}
	//		}
	//		fmt.Println()
	//	}
	//}
}
