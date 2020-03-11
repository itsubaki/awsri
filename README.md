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
        "redshift:DescribeReserved*",
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
$ hermes pricing | jq .
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
$ hermes usage | jq .
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

```
$ hermes usage --format csv | column -t -s, | less -S
```


```
$ cat purchase.json | hermes recommend | jq .
{
  "region": "ap-northeast-1",
  "usage_type": "APN1-BoxUsage:c4.large",
  "platform": "Linux/UNIX",
  "instance_num": 1648
}
```

## API Example

```go
price := pricing.Price{
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
}

forecast := []usage.Quantity{
  {InstanceNum: 120},
  {InstanceNum: 110},
  {InstanceNum: 100},
  {InstanceNum: 90},
  {InstanceNum: 80},
  {InstanceNum: 70},
  {InstanceNum: 60},
  {InstanceNum: 50},
  {InstanceNum: 40},
  {InstanceNum: 30},
  {InstanceNum: 20},
  {InstanceNum: 10},
}

q, _ := BreakEvenPoint(forecast, price)
if q.InstanceNum != 40 {
  t.Errorf("%v", q.InstanceNum)
}
```
