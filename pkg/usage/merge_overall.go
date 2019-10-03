package usage

import "fmt"

func MergeOverall(q []Quantity) []Quantity {
	merged := make(map[string]Quantity)
	for i := range q {
		hash := fmt.Sprintf(
			"%s%s%s%s%s",
			q[i].UsageType,
			q[i].Platform,
			q[i].CacheEngine,
			q[i].DatabaseEngine,
			q[i].Date,
		)

		v, ok := merged[hash]
		if !ok {
			merged[hash] = Quantity{
				Region:         q[i].Region,
				UsageType:      q[i].UsageType,
				Platform:       q[i].Platform,
				CacheEngine:    q[i].CacheEngine,
				DatabaseEngine: q[i].DatabaseEngine,
				Date:           q[i].Date,
				InstanceHour:   q[i].InstanceHour,
				InstanceNum:    q[i].InstanceNum,
			}
			continue
		}

		merged[hash] = Quantity{
			Region:         q[i].Region,
			UsageType:      q[i].UsageType,
			Platform:       q[i].Platform,
			CacheEngine:    q[i].CacheEngine,
			DatabaseEngine: q[i].DatabaseEngine,
			Date:           q[i].Date,
			InstanceHour:   q[i].InstanceHour + v.InstanceHour,
			InstanceNum:    q[i].InstanceNum + v.InstanceNum,
		}
	}

	out := make([]Quantity, 0)
	for k := range merged {
		out = append(out, merged[k])
	}

	return out
}
