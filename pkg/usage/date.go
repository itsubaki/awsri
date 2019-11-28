package usage

import (
	"fmt"
	"sort"
	"time"
)

type Date struct {
	Start string
	End   string
}

func (d Date) YYYYMM() string {
	return d.Start[:7]
}

func Last12Months() []Date {
	return LastNMonths(12)
}

func LastNMonths(n int) []Date {
	return LastNMonthsWith(time.Now(), n)
}

func LastNMonthsWith(now time.Time, n int) []Date {
	if n < 1 || n > 12 {
		panic(fmt.Sprintf("parameter=%v is not in 0 < n < 13", n))
	}

	month := make([]time.Time, 0)
	for i := 1; i < n+1; i++ {
		month = append(month, now.AddDate(0, -i, -now.Day()+1))
	}

	tmp := make([]Date, 0)
	for _, m := range month {
		tmp = append(tmp, Date{
			Start: m.Format("2006-01") + "-01",
			End:   m.AddDate(0, 1, 0).Format("2006-01") + "-01",
		})
	}

	out := make([]Date, 0)
	for i := len(tmp) - 1; i > -1; i-- {
		out = append(out, tmp[i])
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Start > out[j].Start })

	return out
}
