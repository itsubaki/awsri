package usage

func Merge(n []Quantity) []Quantity {
	merged := make(map[string]Quantity)
	for i := range n {
		v, ok := merged[n[i].Hash()]
		if !ok {
			merged[n[i].Hash()] = Quantity{
				Region:         n[i].Region,
				UsageType:      n[i].UsageType,
				Platform:       n[i].Platform,
				CacheEngine:    n[i].CacheEngine,
				DatabaseEngine: n[i].DatabaseEngine,
				Date:           n[i].Date,
				InstanceHour:   n[i].InstanceHour,
				InstanceNum:    n[i].InstanceNum,
			}
			continue
		}

		merged[n[i].Hash()] = Quantity{
			Region:         n[i].Region,
			UsageType:      n[i].UsageType,
			Platform:       n[i].Platform,
			CacheEngine:    n[i].CacheEngine,
			DatabaseEngine: n[i].DatabaseEngine,
			Date:           n[i].Date,
			InstanceHour:   n[i].InstanceHour + v.InstanceHour,
			InstanceNum:    n[i].InstanceNum + v.InstanceNum,
		}
	}

	out := make([]Quantity, 0)
	for k := range merged {
		out = append(out, merged[k])
	}

	return out
}
