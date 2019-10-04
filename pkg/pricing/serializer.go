package pricing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func Serialize(dir, region string, price []Price) error {
	path := fmt.Sprintf("%s/pricing", dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	file := fmt.Sprintf("%s/%s.out", path, region)
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		return nil
	}

	bytes, err := json.Marshal(price)
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	if err := ioutil.WriteFile(file, bytes, os.ModePerm); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}

func Deserialize(dir string, region []string) ([]Price, error) {
	price := make([]Price, 0)
	for _, r := range region {
		file := fmt.Sprintf("%s/pricing/%s.out", dir, r)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return []Price{}, fmt.Errorf("file not found: %v", file)
		}

		read, err := ioutil.ReadFile(file)
		if err != nil {
			return []Price{}, fmt.Errorf("read %s: %v", file, err)
		}

		var p []Price
		if err := json.Unmarshal(read, &p); err != nil {
			return []Price{}, fmt.Errorf("unmarshal: %v", err)
		}

		price = append(price, p...)
	}

	sort.SliceStable(price, func(i, j int) bool { return price[i].Version < price[j].Version })
	sort.SliceStable(price, func(i, j int) bool { return price[i].Region < price[j].Region })
	sort.SliceStable(price, func(i, j int) bool { return price[i].InstanceType < price[j].InstanceType })
	sort.SliceStable(price, func(i, j int) bool { return price[i].LeaseContractLength < price[j].LeaseContractLength })
	sort.SliceStable(price, func(i, j int) bool { return price[i].PurchaseOption < price[j].PurchaseOption })

	return price, nil
}
