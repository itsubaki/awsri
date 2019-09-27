# hermes

[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/hermes?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/hermes)

 - AWS Cost Optimization Library

## Motivation

 In order to reduce AWS cost,
 It is necessary to effectively buy Reserved Instances.
 But AWS pricing is complicated and difficult.
 This library shows the RI that you should buy now,
 based on the future instance usage and the current RI purchase.

## Required

```
# set aws credential "example" with iam policy "hermes"

$ cat ~/.aws/credentials
[example]
aws_access_key_id = ********************
aws_secret_access_key = ****************************************
```

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "hermes",
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeReserved*",
        "rds:DescribeReserved*",
        "elasticache:DescribeReserved*",
        "organizations:List*",
        "organizations:Describe*",
        "ce:Get*"
      ],
      "Resource": "*"
    }
  ]
}
```

## Install

```
$ go get github.com/itsubaki/hermes
```

## API Example

```go
p, err := pricing.Fetch(pricing.Redshift, "ap-northeast-1")
if err != nil {
	fmt.Printf("fetch pricing: %v", err)
}

for _, v := range p {
	fmt.Printf("%#v\n", v)
}

{
  "Version": "20190730012138", 
  "SKU": "PDMPNVN5SPA5HWHH",
  "OfferTermCode": "6QCMYABX3D",
  "Region": "ap-northeast-1",
  "InstanceType": "ds1.8xlarge",
  "UsageType": "APN1-Node:dw.hs1.8xlarge",
  "LeaseContractLength": "1yr",
  "PurchaseOption": "All Upfront",
  "OnDemand": 9.52,
  "ReservedQuantity": 49020,
  "ReservedHrs": 0,
  "Tenancy": "",
  "PreInstalled": "",
  "OperatingSystem": "",
  "Operation": "RunComputeNode:0001",
  "CacheEngine": "", 
  "DatabaseEngine": "",
  "OfferingClass": "standard",
  "NormalizationSizeFactor": ""
}
...
```

```go
u, err := usage.Fetch("2019-08-01", "2019-09-01")
if err != nil {
    fmt.Printf("fetch usage: %v", err)
    os.Exit(1)
}

for i := range u {
	fmt.Printf("%#v\n", u[i])
}

{
  "account_id": "123456789012", 
  "description": "example",
  "region": "us-west-2",
  "usage_type": "USW2-NodeUsage:cache.t2.small",
  "cache_engine": "Redis",
  "date": "2019-08",
  "instance_hour": 101,
  "instance_num": 0.135752688172043
}
...
```

```go
price := []pricing.Price{
  pricing.Price{
  	Region:                  "ap-northeast-1",
  	UsageType:               "APN1-BoxUsage:c4.large",
  	Tenancy:                 "Shared",
  	PreInstalled:            "NA",
  	OperatingSystem:         "Linux",
  	OfferingClass:           "standard",
  	LeaseContractLength:     "1yr",
  	PurchaseOption:          "All Upfront",
  	OnDemand:                0.126,
  	ReservedQuantity:        738,
  	ReservedHrs:             0,
  	NormalizationSizeFactor: "4",
  },
  pricing.Price{
  	Region:                  "ap-northeast-1",
  	UsageType:               "APN1-BoxUsage:c4.xlarge",
  	Tenancy:                 "Shared",
  	PreInstalled:            "NA",
  	OperatingSystem:         "Linux",
  	OfferingClass:           "standard",
  	LeaseContractLength:     "1yr",
  	PurchaseOption:          "All Upfront",
  	OnDemand:                0.126 * 2,
  	ReservedQuantity:        738 * 2,
  	ReservedHrs:             0,
  	NormalizationSizeFactor: "8",
  },
  pricing.Price{
  	Region:                  "ap-northeast-1",
  	UsageType:               "APN1-BoxUsage:c4.2xlarge",
  	Tenancy:                 "Shared",
  	PreInstalled:            "NA",
  	OperatingSystem:         "Linux",
  	OfferingClass:           "standard",
  	LeaseContractLength:     "1yr",
  	PurchaseOption:          "All Upfront",
  	OnDemand:                0.126 * 4,
  	ReservedQuantity:        738 * 4,
  	ReservedHrs:             0,
  	NormalizationSizeFactor: "16",
  },
 }

monthly := MonthlyUsage(quantity)
recommended := make([]usage.Quantity, 0)
for _, p := range price {
    res, err := Recommend(monthly, p)
    if err != nil {
        t.Errorf("recommend: %v", err)
    }
    
    recommended = append(recommended, res)
}

for _, r := range recommended {
    fmt.Printf("%#v\n", r)
}

normalized := make([]usage.Quantity, 0)
for _, r := range recommended {
    n, err := Normalize(r, price)
    if err != nil {
        t.Errorf("recommend: %v", err)
    }

    normalized = append(normalized, n)
}

for _, r := range normalized {
    fmt.Printf("%#v\n", r)
}

usage.Quantity{UsageType:"APN1-BoxUsage:c4.large",   Platform:"Linux/UNIX", InstanceHour:72914.707223,       InstanceNum:98.0036387405914}
usage.Quantity{UsageType:"APN1-BoxUsage:c4.xlarge",  Platform:"Linux/UNIX", InstanceHour:39836.842499,       InstanceNum:55.32894791527778}
usage.Quantity{UsageType:"APN1-BoxUsage:c4.2xlarge", Platform:"Linux/UNIX", InstanceHour:480369.89635399997, InstanceNum:656.8305510524193}

usage.Quantity{UsageType:"APN1-BoxUsage:c4.large", Platform:"Linux/UNIX", InstanceHour:72914.707223,           InstanceNum:98.0036387405914}
usage.Quantity{UsageType:"APN1-BoxUsage:c4.large", Platform:"Linux/UNIX", InstanceHour:79673.684998,           InstanceNum:110.65789583055556}
usage.Quantity{UsageType:"APN1-BoxUsage:c4.large", Platform:"Linux/UNIX", InstanceHour:1.9214795854159999e+06, InstanceNum:2627.3222042096772}
```

## CommandLine Example

```
$ AWS_PROFILE=example hermes fetch
write: /var/tmp/hermes/pricing/ap-northeast-1.out
write: /var/tmp/hermes/pricing/us-west-2.out
write: /var/tmp/hermes/usage/2019-08.out
write: /var/tmp/hermes/usage/2019-07.out
write: /var/tmp/hermes/usage/2019-06.out
write: /var/tmp/hermes/usage/2019-04.out
write: /var/tmp/hermes/usage/2019-03.out
write: /var/tmp/hermes/usage/2019-02.out
write: /var/tmp/hermes/usage/2019-01.out
write: /var/tmp/hermes/usage/2018-12.out
write: /var/tmp/hermes/usage/2018-11.out
write: /var/tmp/hermes/usage/2018-10.out
write: /var/tmp/hermes/usage/2018-09.out
```

```
$ AWS_PROFILE=example hermes pricing | jq .
[ 
  {
    "Version": "20190730012138",
    "SKU": "PDMPNVN5SPA5HWHH",
    "OfferTermCode": "6QCMYABX3D",
    "Region": "ap-northeast-1",
    "InstanceType": "ds1.8xlarge",
    "UsageType": "APN1-Node:dw.hs1.8xlarge",
    "LeaseContractLength": "1yr",
    "PurchaseOption": "All Upfront",
    "OnDemand": 9.52,
    "ReservedQuantity": 49020,
    "ReservedHrs": 0,
    "Tenancy": "",
    "PreInstalled": "",
    "OperatingSystem": "",
    "Operation": "RunComputeNode:0001",
    "CacheEngine": "",
    "DatabaseEngine": "",
    "OfferingClass": "standard",
    "NormalizationSizeFactor": ""
  }
  ...
]
```

```
$ AWS_PROFILE=example hermes usage | jq .
[
  {
    "account_id": "123456789012",
    "description": "example",
    "region": "us-west-2",
    "usage_type": "USW2-NodeUsage:cache.t2.small",
    "cache_engine": "Redis",
    "date": "2019-08",
    "instance_hour": 101,
    "instance_num": 0.135752688172043
  }
  ...
]
```