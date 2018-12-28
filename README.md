# awsri
aws reserved instance

## Memo

```
# awsprice/OperatingSystem
SUSE
Linux
RHEL
Windows
Memcached
Redis
Aurora PostgreSQL
Aurora MySQL
SQL Server
Oracle
PostgreSQL
MySQL
MariaDB

# costexp/Platform
Windows with SQL Server Web
Linux/UNIX
Windows (Amazon VPC)
Windows (BYOL)
NoOperatingSystem
Redis
Memcached
Aurora MySQL
Aurora PostgreSQL
PostgreSQL
MySQL
```

```
Linux/UNIX: Linux
Windows (Amazon VPC): Windows
```

## example

```
repo, _ := awsprice.NewRepository("/awsprice/ap-northeast-1.out")

rs := repo.FindByInstanceType("m4.large").
  OperatingSystem("Linux").
  Tenancy("Shared").
  PreInstalled("NA").
  OfferingClass("standard")

for _, r := range rs {
  fmt.Printf("%s\n", r)
}

for _, r := range rs {
  fmt.Printf("%s\n", r.GetAnnualCost())
}

onDemandInstanceNum := 3
reservedInstanceNum := 10
for _, r := range rs {
  fmt.Printf("%s\n", r.GetCost(onDemandInstanceNum, reservedInstanceNum))
}

{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"4NA7Y494T4","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"1yr","purchase_option":"No Upfront","ondemand":0.129,"reserved_quantity":0,"reserved_hrs":0.0871,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}
{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"HU7G6KETJZ","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"1yr","purchase_option":"Partial Upfront","ondemand":0.129,"reserved_quantity":364,"reserved_hrs":0.0415,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}
{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"6QCMYABX3D","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"1yr","purchase_option":"All Upfront","ondemand":0.129,"reserved_quantity":713,"reserved_hrs":0,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}
{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"BPH4J8HBKS","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"3yr","purchase_option":"No Upfront","ondemand":0.129,"reserved_quantity":0,"reserved_hrs":0.0637,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}
{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"38NPMPTW36","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"3yr","purchase_option":"Partial Upfront","ondemand":0.129,"reserved_quantity":775,"reserved_hrs":0.0295,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}
{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"NQ3QZPMQV9","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"3yr","purchase_option":"All Upfront","ondemand":0.129,"reserved_quantity":1457,"reserved_hrs":0,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}

{"lease_contract_length":"1yr","purchase_option":"No Upfront","ondemand":1130.04,"reserved":762.996,"reserved_quantity":0,"subtraction":367.044,"discount_rate":0.32480620155038764}
{"lease_contract_length":"1yr","purchase_option":"Partial Upfront","ondemand":1130.04,"reserved":727.54,"reserved_quantity":364,"subtraction":402.5,"discount_rate":0.35618208205019297}
{"lease_contract_length":"1yr","purchase_option":"All Upfront","ondemand":1130.04,"reserved":713,"reserved_quantity":713,"subtraction":417.03999999999996,"discount_rate":0.36904888322537255}
{"lease_contract_length":"3yr","purchase_option":"No Upfront","ondemand":3390.12,"reserved":1674.0360000000003,"reserved_quantity":0,"subtraction":1716.0839999999996,"discount_rate":0.5062015503875967}
{"lease_contract_length":"3yr","purchase_option":"Partial Upfront","ondemand":3390.12,"reserved":1550.26,"reserved_quantity":775,"subtraction":1839.86,"discount_rate":0.5427123523651081}
{"lease_contract_length":"3yr","purchase_option":"All Upfront","ondemand":3390.12,"reserved":1457,"reserved_quantity":1457,"subtraction":1933.12,"discount_rate":0.5702217030665582}

{"lease_contract_length":"1yr","purchase_option":"No Upfront","full_ondemand":{"ondemand":14690.52,"reserved":0,"total":14690.52},"reserved_applied":{"ondemand":3390.12,"reserved":7629.96,"total":11020.08},"reserved_quantity":0,"subtraction":3670.4400000000005,"discount_rate":0.24985092426952893}
{"lease_contract_length":"1yr","purchase_option":"Partial Upfront","full_ondemand":{"ondemand":14690.52,"reserved":0,"total":14690.52},"reserved_applied":{"ondemand":3390.12,"reserved":7275.4,"total":10665.52},"reserved_quantity":3640,"subtraction":4025,"discount_rate":0.2739862169616869}
{"lease_contract_length":"1yr","purchase_option":"All Upfront","full_ondemand":{"ondemand":14690.52,"reserved":0,"total":14690.52},"reserved_applied":{"ondemand":3390.12,"reserved":7130,"total":10520.119999999999},"reserved_quantity":7130,"subtraction":4170.4000000000015,"discount_rate":0.28388375632720975}
{"lease_contract_length":"3yr","purchase_option":"No Upfront","full_ondemand":{"ondemand":44071.56,"reserved":0,"total":44071.56},"reserved_applied":{"ondemand":10170.36,"reserved":16740.360000000004,"total":26910.720000000005},"reserved_quantity":0,"subtraction":17160.839999999993,"discount_rate":0.38938580799045897}
{"lease_contract_length":"3yr","purchase_option":"Partial Upfront","full_ondemand":{"ondemand":44071.56,"reserved":0,"total":44071.56},"reserved_applied":{"ondemand":10170.36,"reserved":15502.6,"total":25672.96},"reserved_quantity":7750,"subtraction":18398.6,"discount_rate":0.41747104028085236}
{"lease_contract_length":"3yr","purchase_option":"All Upfront","full_ondemand":{"ondemand":44071.56,"reserved":0,"total":44071.56},"reserved_applied":{"ondemand":10170.36,"reserved":14570,"total":24740.36},"reserved_quantity":14570,"subtraction":19331.199999999997,"discount_rate":0.4386320792819678}
```

```
repo, _ := costexp.NewRepository("/costexp/example_2018-11.out")

for _, r := range repo.SelectAll() {
  fmt.Println(r)
}

{"account_id":"123456789012","date":"2018-11","usage_type":"APN1-BoxUsage:c4.2xlarge","platform":"Linux/UNIX","instance_hour":175.600833,"instance_num":0.24389004583333332}
{"account_id":"123456789012","date":"2018-11","usage_type":"APN1-BoxUsage:c4.large","platform":"Linux/UNIX","instance_hour":720,"instance_num":1}
{"account_id":"123456789012","date":"2018-11","usage_type":"APN1-BoxUsage:t2.micro","platform":"Linux/UNIX","instance_hour":2264.238066,"instance_num":3.1447750916666664}
{"account_id":"123456789012","date":"2018-11","usage_type":"APN1-BoxUsage:t2.nano","platform":"Linux/UNIX","instance_hour":720,"instance_num":1}
{"account_id":"123456789012","date":"2018-11","usage_type":"APN1-BoxUsage:t2.small","platform":"Linux/UNIX","instance_hour":1440,"instance_num":2}
{"account_id":"123456789012","date":"2018-11","usage_type":"APN1-NodeUsage:cache.r5.large","engine":"Redis","instance_hour":2,"instance_num":0.002777777777777778}
{"account_id":"123456789012","date":"2018-11","usage_type":"APN1-NodeUsage:cache.t2.micro","engine":"Redis","instance_hour":344,"instance_num":0.4777777777777778}
{"account_id":"123456789012","date":"2018-11","usage_type":"APN1-NodeUsage:cache.t2.small","engine":"Redis","instance_hour":72,"instance_num":0.1}
{"account_id":"123456789012","date":"2018-11","usage_type":"APN1-InstanceUsage:db.r3.large","engine":"Aurora MySQL","instance_hour":1,"instance_num":0.001388888888888889}
{"account_id":"123456789012","date":"2018-11","usage_type":"APN1-InstanceUsage:db.r4.large","engine":"Aurora MySQL","instance_hour":2,"instance_num":0.002777777777777778}
{"account_id":"123456789012","date":"2018-11","usage_type":"APN1-InstanceUsage:db.t2.small","engine":"Aurora MySQL","instance_hour":237,"instance_num":0.32916666666666666}
```

```
r := &Record{
  SKU:                     "7MYWT7Y96UT3NJ2D",
  OfferTermCode:           "4NA7Y494T4",
  Region:                  "ap-northeast-1",
  InstanceType:            "m4.large",
  UsageType:               "APN1-BoxUsage:m4.large",
  LeaseContractLength:     "1yr",
  PurchaseOption:          "All Upfront",
  OnDemand:                0.129,
  ReservedHrs:             0,
  ReservedQuantity:        713,
  Tenancy:                 "Shared",
  PreInstalled:            "NA",
  OperatingSystem:         "Linux",
  Operation:               "RunInstances",
  OfferingClass:           "standard",
  NormalizationSizeFactor: "4",
}

forecast := []Forecast{
  {Month: "2018-01", InstanceNum: 120.4},
  {Month: "2018-02", InstanceNum: 110.3},
  {Month: "2018-03", InstanceNum: 100.1},
  {Month: "2018-04", InstanceNum: 90.9},
  {Month: "2018-05", InstanceNum: 80.9},
  {Month: "2018-06", InstanceNum: 70.6},
  {Month: "2018-07", InstanceNum: 60.3},
  {Month: "2018-08", InstanceNum: 50.9},
  {Month: "2018-09", InstanceNum: 40.7},
  {Month: "2018-10", InstanceNum: 30.6},
  {Month: "2018-11", InstanceNum: 20.2},
  {Month: "2018-12", InstanceNum: 10.8},
}

fmt.Println(r.Recommend(forecast))
{
  "record":
  {
    "sku":"7MYWT7Y96UT3NJ2D",
    "offer_term_code":"4NA7Y494T4",
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
  "breakevenpoint_in_month":8,
  "ondemand_instance_num_avg":23.7,
  "reserved_instance_num":50,
  "full_ondemand_cost":83283.948,
  "reserved_applied_cost":
  {
    "ondemand":26781.947999999997,
    "reserved":35650,
    "total":62431.948
  },
  "reserved_quantity":35650,
  "subtraction":20852.000000000007,
  "discount_rate":0.2503723766793573
}
```
