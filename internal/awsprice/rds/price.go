package rds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var BaseURL = "https://pricing.us-east-1.amazonaws.com"

type InputPrice struct {
	FormatVersion   string               `json:"formatVersion"`
	Disclaimer      string               `json:"disclaimer"`
	PublicationDate string               `json:"publicationDate"`
	Regions         map[string]RegionUrl `json:"regions"`
}

type RegionUrl struct {
	RegionCode        string `json:"regionCode"`
	CurrentVersionUrl string `json:"currentVersionUrl"`
}

type PriceList struct {
	FormatVersion   string                                `json:"formatVersion"`
	Disclaimer      string                                `json:"disclaimer"`
	OfferCode       string                                `json:"offerCode"`
	Version         string                                `json:"version"`
	PublicationDate string                                `json:"publicationDate"`
	Products        map[string]Product                    `json:"products"`
	Term            map[string]map[string]map[string]Term `json:"terms"`
}

type Term struct {
	SKU             string                     `json:"sku"`
	OfferTermCode   string                     `json:"offerTermCode"`
	EffectiveDate   string                     `json:"effectiveDate"`
	PriceDimensions map[string]PriceDimensions `json:"priceDimensions"`
	TermAttributes  TermAttributes             `json:"termAttributes"`
}

type TermAttributes struct {
	LeaseContractLength string `json:"LeaseContractLength"`
	OfferingClass       string `json:"OfferingClass"`
	PurchaseOption      string `json:"PurchaseOption"`
}

type PriceDimensions struct {
	RateCode     string       `json:"rateCode"`
	Description  string       `json:"description"`
	BeginRange   string       `json:"beginRange"`
	EndRange     string       `json:"endRange"`
	Unit         string       `json:"unit"`
	PricePerUnit PricePerUnit `json:"pricePerUnit"`
	AppliesTo    []string     `json:"appliesTo"`
}

type PricePerUnit struct {
	USD string `json:"USD"`
}

type Product struct {
	SKU           string            `json:"sku"`
	ProductFamily string            `json:"productFamily"`
	Attributes    map[string]string `json:"attributes"`
}

type OutputPrice struct {
	Version                 string  // common
	SKU                     string  // common
	OfferTermCode           string  // common
	Region                  string  // common
	InstanceType            string  // common
	UsageType               string  // common
	LeaseContractLength     string  // common
	PurchaseOption          string  // common
	OnDemand                float64 // common
	ReservedQuantity        float64 // common
	ReservedHrs             float64 // common
	OfferingClass           string  // ec2, rds
	NormalizationSizeFactor string  // ec2, rds
	DatabaseEngine          string  // rds
}

func GetPrice(region string) (map[string]OutputPrice, error) {
	return GetPriceWithClient(region, http.DefaultClient)
}

func GetPriceWithClient(region string, client *http.Client) (map[string]OutputPrice, error) {
	var input InputPrice
	{
		url := fmt.Sprintf("%s/offers/v1.0/aws/AmazonRDS/current/region_index.json", BaseURL)
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("get %s: %v", BaseURL, err)
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read body: %v", err)
		}

		if err := json.Unmarshal(buf, &input); err != nil {
			return nil, fmt.Errorf("unmarshal: %v", err)
		}
	}

	var list PriceList
	{
		url := fmt.Sprintf("%s%s", BaseURL, input.Regions[region].CurrentVersionUrl)
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("get %s: %v", BaseURL, err)
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read body: %v", err)
		}

		if err := json.Unmarshal(buf, &list); err != nil {
			return nil, fmt.Errorf("unmarshal: %v", err)
		}
	}

	return usage(region, list)
}

func usage(region string, list PriceList) (map[string]OutputPrice, error) {
	p := make(map[string]OutputPrice)
	{
		for _, t := range list.Term["Reserved"] {
			for k, v := range t {
				var q, h float64
				for _, vv := range v.PriceDimensions {
					if vv.Unit == "Quantity" {
						q, _ = strconv.ParseFloat(vv.PricePerUnit.USD, 64)
					}
					if vv.Unit == "Hrs" {
						h, _ = strconv.ParseFloat(vv.PricePerUnit.USD, 64)
					}
				}

				p[k] = OutputPrice{
					Version:             list.Version,
					SKU:                 v.SKU,
					OfferTermCode:       v.OfferTermCode,
					LeaseContractLength: v.TermAttributes.LeaseContractLength,
					PurchaseOption:      v.TermAttributes.PurchaseOption,
					OfferingClass:       v.TermAttributes.OfferingClass,
					ReservedHrs:         h,
					ReservedQuantity:    q,
				}
			}
		}

		for _, t := range list.Term["OnDemand"] {
			for k, v := range t { // 1
				for kk, pp := range p {
					if !strings.HasPrefix(k, pp.SKU) {
						continue
					}

					for _, vv := range v.PriceDimensions { // 1
						hrs, _ := strconv.ParseFloat(vv.PricePerUnit.USD, 64)
						p[kk] = OutputPrice{
							Version:             p[kk].Version,
							SKU:                 p[kk].SKU,
							OfferTermCode:       p[kk].OfferTermCode,
							LeaseContractLength: p[kk].LeaseContractLength,
							PurchaseOption:      p[kk].PurchaseOption,
							OfferingClass:       p[kk].OfferingClass,
							ReservedHrs:         p[kk].ReservedHrs,
							ReservedQuantity:    p[kk].ReservedQuantity,
							OnDemand:            hrs,
						}
					}
				}
			}
		}
	}

	out := make(map[string]OutputPrice)
	for _, pp := range list.Products {
		for k, v := range p {
			if !strings.HasPrefix(k, pp.SKU) {
				continue
			}

			out[k] = OutputPrice{
				Version:                 v.Version,
				SKU:                     v.SKU,
				OfferTermCode:           v.OfferTermCode,
				Region:                  region,
				InstanceType:            pp.Attributes["instanceType"],
				UsageType:               pp.Attributes["usagetype"],
				DatabaseEngine:          pp.Attributes["databaseEngine"],
				LeaseContractLength:     v.LeaseContractLength,
				PurchaseOption:          v.PurchaseOption,
				OfferingClass:           v.OfferingClass,
				OnDemand:                v.OnDemand,
				ReservedQuantity:        v.ReservedQuantity,
				ReservedHrs:             v.ReservedHrs,
				NormalizationSizeFactor: pp.Attributes["normalizationSizeFactor"],
			}
		}
	}

	return out, nil
}
