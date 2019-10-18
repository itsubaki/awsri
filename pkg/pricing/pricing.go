package pricing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var URL = []string{
	Compute,
	Database,
	Cache,
	Redshift,
}

var BaseURL = "https://pricing.us-east-1.amazonaws.com"
var Compute = fmt.Sprintf("%s%s", BaseURL, "/offers/v1.0/aws/AmazonEC2/current/region_index.json")
var Database = fmt.Sprintf("%s%s", BaseURL, "/offers/v1.0/aws/AmazonRDS/current/region_index.json")
var Cache = fmt.Sprintf("%s%s", BaseURL, "/offers/v1.0/aws/AmazonElastiCache/current/region_index.json")
var Redshift = fmt.Sprintf("%s%s", BaseURL, "/offers/v1.0/aws/AmazonRedshift/current/region_index.json")

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
	Terms           map[string]map[string]map[string]Term `json:"terms"`
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

type Price struct {
	Version                 string  `json:"version,omitempty"`                   // common
	SKU                     string  `json:"sku,omitempty"`                       // common
	OfferTermCode           string  `json:"offer_term_code,omitempty"`           // common
	Region                  string  `json:"region,omitempty"`                    // common
	InstanceType            string  `json:"instance_type,omitempty"`             // common
	UsageType               string  `json:"usage_type,omitempty"`                // common
	LeaseContractLength     string  `json:"lease_contract_length,omitempty"`     // common
	PurchaseOption          string  `json:"purchase_option,omitempty"`           // common
	OnDemand                float64 `json:"ondemand,omitempty"`                  // common
	ReservedQuantity        float64 `json:"reserved_quantity,omitempty"`         // common
	ReservedHrs             float64 `json:"reserved_hours,omitempty"`            // common
	Tenancy                 string  `json:"tenancy,omitempty"`                   // compute: Shared, Host, Dedicated
	PreInstalled            string  `json:"pre_installed,omitempty"`             // compute: SQL Web, SQL Ent, SQL Std, NA
	Operation               string  `json:"operation,omitempty"`                 // compute
	OperatingSystem         string  `json:"operating_system,omitempty"`          // compute: Windows, Linux, SUSE, RHEL
	CacheEngine             string  `json:"cache_engine,omitempty"`              // cache
	DatabaseEngine          string  `json:"database_engine,omitempty"`           // database
	OfferingClass           string  `json:"offering_class,omitempty"`            // compute, database
	NormalizationSizeFactor string  `json:"normalization_size_factor,omitempty"` // compute, database
}

func (p Price) ID() string {
	return fmt.Sprintf(
		"%s_%s_%s_%s%s%s_%s_%s_%s",
		p.UsageType,
		p.LeaseContractLength,
		p.PurchaseOption,
		p.OperatingSystem,
		p.CacheEngine,
		p.DatabaseEngine,
		p.Tenancy,
		p.PreInstalled,
		p.OfferingClass,
	)
}

func (p Price) String() string {
	return p.JSON()
}

func (p Price) JSON() string {
	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func (p Price) DiscountRate() float64 {
	month := 12
	if p.LeaseContractLength == "3yr" {
		month = 12 * 3
	}

	ond, res := 0.0, p.ReservedQuantity
	for i := 1; i < month+1; i++ {
		ond, res = ond+p.OnDemand*24*float64(Days[i%12]), res+p.ReservedHrs*24*float64(Days[i%12])
	}

	if ond == 0.0 {
		return 0.0
	}

	return (ond - res) / ond
}

func (p Price) BreakEvenPoint() int {
	month := 12
	if p.LeaseContractLength == "3yr" {
		month = 12 * 3
	}

	res := p.ReservedQuantity
	for i := 1; i < month+1; i++ {
		res = res + p.ReservedHrs*24*float64(Days[i%12])
	}

	out, ond := 0, 0.0
	for i := 1; i < month+1; i++ {
		ond = ond + p.OnDemand*24*float64(Days[i%12])
		if ond > res {
			out = i
			break
		}
	}

	return out
}

func Fetch(url, region string) (map[string]Price, error) {
	return FetchWithClient(url, region, http.DefaultClient)
}

func FetchWithClient(url, region string, client *http.Client) (map[string]Price, error) {
	var input InputPrice
	{
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("get %s: %v", url, err)
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
			return nil, fmt.Errorf("get %s: %v", url, err)
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read body: %v", err)
		}

		if err := json.Unmarshal(buf, &list); err != nil {
			return nil, fmt.Errorf("unmarshal: %v", err)
		}
	}

	return fetch(region, list)
}

func fetch(region string, list PriceList) (map[string]Price, error) {
	p := make(map[string]Price)
	{
		for _, t := range list.Terms["Reserved"] {
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

				// k is SKU.OfferingTermCode. it is unique.
				p[k] = Price{
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

		for _, t := range list.Terms["OnDemand"] {
			for k, v := range t { // 1
				for kk, pp := range p {
					if !strings.HasPrefix(k, pp.SKU) {
						continue
					}

					for _, vv := range v.PriceDimensions { // 1
						hrs, _ := strconv.ParseFloat(vv.PricePerUnit.USD, 64)
						p[kk] = Price{
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

	out := make(map[string]Price)
	for _, pp := range list.Products {
		for k, v := range p {
			if !strings.HasPrefix(k, pp.SKU) {
				continue
			}

			out[k] = Price{
				Version:                 v.Version,
				SKU:                     v.SKU,
				OfferTermCode:           v.OfferTermCode,
				Region:                  region,
				InstanceType:            pp.Attributes["instanceType"],
				UsageType:               pp.Attributes["usagetype"],
				Tenancy:                 pp.Attributes["tenancy"],
				PreInstalled:            pp.Attributes["preInstalledSw"],
				OperatingSystem:         pp.Attributes["operatingSystem"],
				Operation:               pp.Attributes["operation"],
				CacheEngine:             pp.Attributes["cacheEngine"],
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
