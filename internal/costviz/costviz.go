package costviz

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strings"
)

type CostViz struct {
	BaseURL   string
	XApiKey   string
	AccountID string
	TableName []string
}

type UtilList []*Utilization

type Utilization struct {
	AccountID       string  `json:"account_id"`
	Date            string  `json:"date"`
	ID              string  `json:"id"`
	UsageType       string  `json:"usage_type"`
	OperatingSystem string  `json:"operating_system,omitempty"`
	Engine          string  `json:"engine,omitempty"`
	InstanceHour    float64 `json:"instance_hour"`
	InstanceNum     float64 `json:"instance_num"`
}

func (u *Utilization) String() string {
	bytea, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	return string(bytea)
}
func (v *CostViz) GetOSEngine() ([]string, error) {
	url := fmt.Sprintf("https://%s/linechart/?target=os_engine&summary=usageamount&stack=on&account=%s", v.BaseURL, v.AccountID)

	osengine := make(map[string]interface{})
	for i := range v.TableName {
		url = fmt.Sprintf("%s&tablename=%s", url, v.TableName[i])
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("x-api-key", v.XApiKey)

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return []string{}, fmt.Errorf("get %s: %v", url, err)
		}
		defer resp.Body.Close()

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []string{}, fmt.Errorf("read body: %v", err)
		}

		var tmp map[string]interface{}
		if err := json.Unmarshal(buf, &tmp); err != nil {
			return []string{}, fmt.Errorf("unmarshal: %v", err)
		}

		for k, v := range tmp {
			osengine[k] = v
		}
	}

	out := []string{}
	for k, _ := range osengine {
		out = append(out, k)
	}

	return out, nil
}
func (v *CostViz) GetUtilization() (UtilList, error) {
	osengine, err := v.GetOSEngine()
	if err != nil {
		return nil, fmt.Errorf("get os engine: %v", err)
	}

	out := UtilList{}
	for _, tableName := range v.TableName {
		for _, os := range osengine {
			u, err := v.getUtilization(tableName, os)
			if err != nil {
				return nil, fmt.Errorf("get utilization: %v", err)
			}
			out = append(out, u...)
		}
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (v *CostViz) getUtilization(tableName, os string) (UtilList, error) {
	hash := md5.Sum([]byte(os))
	hstr := hex.EncodeToString(hash[:])
	url := fmt.Sprintf(
		"https://%s/linechart/?target=usagetype&summary=usageamount&stack=on&account=%s&os_engine=%s&tablename=%s",
		v.BaseURL,
		v.AccountID,
		hstr,
		tableName,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-api-key", v.XApiKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get %s: %v", url, err)
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %v", err)
	}

	var usage map[string]map[string]float64
	if err := json.Unmarshal(buf, &usage); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	keys := []string{}
	for k := range usage {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := UtilList{}
	for _, k := range keys {
		max := 0.0
		for _, w := range usage[k] {
			max = math.Max(max, w)
		}

		r := &Utilization{
			ID:           fmt.Sprintf("%s:%s", k, strings.Replace(os, " ", "", -1)),
			AccountID:    v.AccountID,
			Date:         strings.Split(tableName, "_")[1], // awsbilling_201806 -> 201806
			UsageType:    k,                                // APN1-NodeUsage:cache.t2.micro
			InstanceHour: max,
		}

		date := r.Date
		days := Month[date[len(date)-2:]]
		r.InstanceNum = r.InstanceHour / float64(24*days)

		if strings.Contains(k, "cache.") || strings.Contains(k, "db.") {
			r.Engine = os
			out = append(out, r)
			continue
		}

		// ec2
		r.OperatingSystem = os
		out = append(out, r)
	}

	return out, nil
}
