package reservation

import "fmt"

func Merge(u []Utilization) []Utilization {
	merged := make(map[string]Utilization)
	for i := range u {
		hash := fmt.Sprintf(
			"%s%s%s%s%s%s%s%s",
			u[i].AccountID,
			u[i].Region,
			u[i].InstanceType,
			u[i].Platform,
			u[i].CacheEngine,
			u[i].DatabaseEngine,
			u[i].DeploymentOption,
			u[i].Date,
		)

		v, ok := merged[hash]
		if !ok {
			merged[hash] = u[i]
			continue
		}

		merged[hash] = Utilization{
			AccountID:        u[i].AccountID,
			Description:      u[i].Description,
			Region:           u[i].Region,
			InstanceType:     u[i].InstanceType,
			Platform:         u[i].Platform,
			CacheEngine:      u[i].CacheEngine,
			DatabaseEngine:   u[i].DatabaseEngine,
			DeploymentOption: u[i].DeploymentOption,
			Date:             u[i].Date,
			Hours:            u[i].Hours + v.Hours,
			Percentage:       u[i].Percentage + v.Percentage,
		}
	}

	out := make([]Utilization, 0)
	for k := range merged {
		out = append(out, merged[k])
	}

	return out
}
