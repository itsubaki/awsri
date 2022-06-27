package flag

import "strings"

func Split(value []string) []string {
	ret := make([]string, 0)
	for _, v := range value {
		s := strings.Split(v, ",")
		for _, ss := range s {
			ret = append(ret, strings.TrimSpace(ss))
		}
	}

	return ret
}
