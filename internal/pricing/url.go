package pricing

import "fmt"

var URL = []string{
	ComputeURL,
	DatabaseURL,
	CacheURL,
	RedshiftURL,
}

var BaseURL = "https://pricing.us-east-1.amazonaws.com"
var ComputeURL = fmt.Sprintf("%s%s", BaseURL, "/offers/v1.0/aws/AmazonEC2/current/region_index.json")
var DatabaseURL = fmt.Sprintf("%s%s", BaseURL, "/offers/v1.0/aws/AmazonRDS/current/region_index.json")
var CacheURL = fmt.Sprintf("%s%s", BaseURL, "/offers/v1.0/aws/AmazonElastiCache/current/region_index.json")
var RedshiftURL = fmt.Sprintf("%s%s", BaseURL, "/offers/v1.0/aws/AmazonRedshift/current/region_index.json")
