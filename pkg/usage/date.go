package usage

import "time"

type Date struct {
	Start string
	End   string
}

func (d Date) YYYYMM() string {
	return d.Start[:7]
}

func Last12Months() []Date {
	month := make([]time.Time, 0)
	for i := 1; i < 13; i++ {
		month = append(month, time.Now().AddDate(0, -i, 0))
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

	return out
}
