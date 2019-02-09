# hermes

[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/hermes?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/hermes)

 - aws reserved instance purchase recommendation library

## Motivation

 In order to reduce AWS cost,
 It is necessary to effectively buy Reserved Instance.
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
$ cd ${GOPATH}/src/github.com/itsubaki/hermes
$ make
```

## Example

```
$ cat test/forecast.json | hermes > test.out
$ column -s, -t < test.out | less -#2 -N -S

account_id,   alies,   usage_type,                      platform/engine, 2019-01, 2018-02, 2019-03, 2019-04, 2019-05, 2019-06, 2019-07, 2019-08, 2019-09, 2019-10, 2019-11, 2019-12,
123456789012, example, APN1-BoxUsage:c4.2xlarge,        Linux/UNIX,      100,     100,     100,     100,     100,     100,     100,     100, 100,     100,     100,     100,
123456789012, example, APN1-InstanceUsage:db.r3.xlarge, Aurora MySQL,    100,     100,     100,     100,     100,     100,     100,     100, 100,     100,     100,     100,
123456789012, example, APN1-NodeUsage:cache.r3.4xlarge, Redis,           100,     100,     100,     100,     100,     100,     100,     100, 100,     100,     100,     100,

                      usage_type,                      platform/engine, 2019-01, 2018-02, 2019-03, 2019-04, 2019-05, 2019-06, 2019-07, 2019-08, 2019-09, 2019-10, 2019-11, 2019-12,
                      APN1-BoxUsage:c4.2xlarge,        Linux/UNIX,      100,     100,     100,     100,     100,     100,     100,     100,     100,     100,     100,     100,
                      APN1-InstanceUsage:db.r3.xlarge, Aurora MySQL,    100,     100,     100,     100,     100,     100,     100,     100,     100,     100,     100,     100,
                      APN1-NodeUsage:cache.r3.4xlarge, Redis,           100,     100,     100,     100,     100,     100,     100,     100,     100,     100,     100,     100,

                      usage_type,                      os/engine,    ondemand_num_avg, reserved_num, full_ondemand_cost, reserved_applied_cost, subtraction, discount_rate,      reserved_quantity,
                      APN1-BoxUsage:c4.2xlarge,        Linux,        0,                100,          441504,             296200,                145304,      0.3291114010292092, 296200,
                      APN1-InstanceUsage:db.r3.xlarge, Aurora MySQL, 0,                100,          613200,             340800,                272400,      0.4442270058708415, 340800,
                      APN1-NodeUsage:cache.r3.4xlarge, Redis,        0,                100,          1.913184e+06,       1.245312e+06,          667872,      0.3490892668974861, 621600,

                      usage_type,                      os/engine,    instance_num,
                      APN1-BoxUsage:c4.large,          Linux,        400,
                      APN1-InstanceUsage:db.r3.large,  Aurora MySQL, 200,
                      APN1-NodeUsage:cache.r3.4xlarge, Redis,        100,
```

## API Example

```
# get current usage
date := []*costexp.Date{
  {
    Start: "2020-11-01",
    End:   "2020-12-01",
  },
}

repo := costexp.New(date)
for _, r := range repo.SelectAll() {
  fmt.Println(r)
}

{
  "account_id":"123456789012",
  "date":"2020-11",
  "usage_type":"APN1-BoxUsage:m4.4xlarge",
  "platform":"Linux/UNIX",
  "instance_hour":2264.238066,
  "instance_num":3.1447750916666664
}

# find aws pricing of current usage
repo := pricing.New([]string{"ap-northeast-1"})
rs := repo.FindByUsageType("APN1-BoxUsage:m4.4xlarge").
  OperatingSystem("Linux").
  Tenancy("Shared").
  LeaseContractLength("1yr").
  PurchaseOption("All Upfront").
  OfferingClass("standard")

# predict future usage (the method is various)
forecast := []pricing.Forecast{
  {Date: "2021-01", InstanceNum: 120.4},
  {Date: "2021-02", InstanceNum: 110.3},
  {Date: "2021-03", InstanceNum: 100.1},
  {Date: "2021-04", InstanceNum: 90.9},
  {Date: "2021-05", InstanceNum: 80.9},
  {Date: "2021-06", InstanceNum: 70.6},
  {Date: "2021-07", InstanceNum: 60.3},
  {Date: "2021-08", InstanceNum: 50.9},
  {Date: "2021-09", InstanceNum: 40.7},
  {Date: "2021-10", InstanceNum: 30.6},
  {Date: "2021-11", InstanceNum: 20.2},
  {Date: "2021-12", InstanceNum: 10.8},
}

# get recommended reserved instance
result, _ := repo.Recommend(rs[0], forecast)
fmt.Println(result)

{
  "record":{
    "sku":"XU2NYYPCRTK4T7CN",
    "offer_term_code":"6QCMYABX3D",
    "region":"ap-northeast-1",
    "instance_type":"m4.4xlarge",
    "usage_type":"APN1-BoxUsage:m4.4xlarge",
    "lease_contract_length":"1yr",
    "purchase_option":"All Upfront",
    "ondemand":1.032,
    "reserved_quantity":5700,
    "reserved_hrs":0,
    "tenancy":"Shared",
    "pre_installed":"NA",
    "operating_system":"Linux",
    "operation":"RunInstances",
    "offering_class":"standard",
    "normalization_size_factor":"32"
  },
  "breakevenpoint_in_month":8,
  "strategy":"breakevenpoint",
  "ondemand_instance_num_avg":23.7,
  "reserved_instance_num":50,
  "full_ondemand_cost":666271.584,
  "reserved_applied_cost":{
    "ondemand":214255.58399999997,
    "reserved":285000,
    "total":499255.584
  },
  "reserved_quantity":285000,
  "subtraction":167016.00000000006,
  "discount_rate":0.2506725545719808,
  "minimum_record":{
    "sku":"7MYWT7Y96UT3NJ2D",
    "offer_term_code":"6QCMYABX3D",
    "region":"ap-northeast-1",
    "instance_type":"m4.large",
    "usage_type":"APN1-BoxUsage:m4.large",
    "lease_contract_length":"1yr",
    "purchase_option":"All Upfront",
    "ondemand":0.129,
    "reserved_quantity":713,
    "reserved_hrs":0,
    "tenancy":"Shared",
    "pre_installed":"NA",
    "operating_system":"Linux",
    "operation":"RunInstances",
    "offering_class":"standard",
    "normalization_size_factor":"4"
  },
  "minimum_reserved_instance_num":400
}

# buy m4.large x400 instead of m4.4xlarge x50
# and

min := result.MinimumRecord

repo := reservation.New([]string{"ap-northeast-1"})
rs := repo.FindByInstanceType(min.InstanceType).
  Region(min.Region).
  Duration(func(length string) int64 {
    duration := 31536000
    if length == "3yr" {
      duration = 94608000
    }
    return int64(duration)
  }(min.LeaseContractLength)).
  OfferingClass(min.OfferingClass).
  OfferingType(min.PurchaseOption).
  ContainsProductDescription(min.OperatingSystem)

fmt.Println(rs[0])

{
  "region":"ap-northeast-1",
  "instance_type":"m4.large",
  "duration":31536000,
  "offering_type":"All Upfront",
  "offering_class":"standard",
  "product_description":"Linux/UNIX (Amazon VPC)",
  "instance_count":100,
  "start":"2020-12-01T12:00:00Z"
}

# already bought 100 instances
# finally, buy m4.large x300
```

```
repo := pricing.New([]string{"ap-northeast-1"})
rs := repo.FindByInstanceType("m4.large").
  OperatingSystem("Linux").
  Tenancy("Shared").
  PreInstalled("NA").
  OfferingClass("standard").
  LeaseContractLength("1yr").
  PurchaseOption("All Upfront")

for _, r := range rs {
  fmt.Printf("%s\n", r)
}

for _, r := range rs {
  fmt.Printf("%s\n", r.GetAnnualCost())
}

ondemand := 3
reserved := 10
for _, r := range rs {
  fmt.Printf("%s\n", r.GetCost(ondemand, reserved))
}

{
  "sku":"7MYWT7Y96UT3NJ2D",
  "offer_term_code":"6QCMYABX3D",
  "region":"ap-northeast-1",
  "instance_type":"m4.large",
  "usage_type":"APN1-BoxUsage:m4.large",
  "lease_contract_length":"1yr",
  "purchase_option":"All Upfront",
  "ondemand":0.129,
  "reserved_quantity":713,
  "reserved_hrs":0,
  "tenancy":"Shared",
  "pre_installed":"NA",
  "operating_system":"Linux",
  "operation":"RunInstances",
  "offering_class":"standard",
  "normalization_size_factor":"4"
}

{
  "lease_contract_length":"1yr",
  "purchase_option":"All Upfront",
  "ondemand":1130.04,
  "reserved":713,
  "reserved_quantity":713,
  "subtraction":417.03999999999996,
  "discount_rate":0.36904888322537255
}

{
  "lease_contract_length":"1yr",
  "purchase_option":"All Upfront",
  "full_ondemand":14690.52,
  "reserved_applied":
  {
    "ondemand":3390.12,
    "reserved":7130,
    "total":10520.119999999999
  },
  "reserved_quantity":7130,
  "subtraction":4170.4000000000015,
  "discount_rate":0.28388375632720975
}
```
