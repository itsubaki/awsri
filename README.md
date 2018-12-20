# awsri
aws reserved instance

## example

```
repo, err := NewRepository("ap-northeast-1")
if err != nil {
  t.Errorf("%v", err)
}

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

{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"4NA7Y494T4","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"1yr","purchase_option":"No Upfront","on_demand":0.129,"reserved_quantity":0,"reserved_hrs":0.0871,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}
{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"HU7G6KETJZ","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"1yr","purchase_option":"Partial Upfront","on_demand":0.129,"reserved_quantity":364,"reserved_hrs":0.0415,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}
{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"6QCMYABX3D","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"1yr","purchase_option":"All Upfront","on_demand":0.129,"reserved_quantity":713,"reserved_hrs":0,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}
{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"BPH4J8HBKS","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"3yr","purchase_option":"No Upfront","on_demand":0.129,"reserved_quantity":0,"reserved_hrs":0.0637,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}
{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"38NPMPTW36","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"3yr","purchase_option":"Partial Upfront","on_demand":0.129,"reserved_quantity":775,"reserved_hrs":0.0295,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}
{"sku":"7MYWT7Y96UT3NJ2D","offer_term_code":"NQ3QZPMQV9","region":"ap-northeast-1","instance_type":"m4.large","usage_type":"APN1-BoxUsage:m4.large","lease_contract_length":"3yr","purchase_option":"All Upfront","on_demand":0.129,"reserved_quantity":1457,"reserved_hrs":0,"tenancy":"Shared","pre_installed":"NA","operating_system":"Linux","operation":"RunInstances","offering_class":"standard","normalization_size_factor":"4"}

{"lease_contract_length":"1yr","purchase_option":"No Upfront","on_demand":1130.04,"reserved":762.996,"reserved_quantity":0,"subtraction":367.044,"discount_rate":0.32480620155038764}
{"lease_contract_length":"1yr","purchase_option":"Partial Upfront","on_demand":1130.04,"reserved":727.54,"reserved_quantity":364,"subtraction":402.5,"discount_rate":0.35618208205019297}
{"lease_contract_length":"1yr","purchase_option":"All Upfront","on_demand":1130.04,"reserved":713,"reserved_quantity":713,"subtraction":417.03999999999996,"discount_rate":0.36904888322537255}
{"lease_contract_length":"3yr","purchase_option":"No Upfront","on_demand":3390.12,"reserved":1674.0360000000003,"reserved_quantity":0,"subtraction":1716.0839999999996,"discount_rate":0.5062015503875967}
{"lease_contract_length":"3yr","purchase_option":"Partial Upfront","on_demand":3390.12,"reserved":1550.26,"reserved_quantity":775,"subtraction":1839.86,"discount_rate":0.5427123523651081}
{"lease_contract_length":"3yr","purchase_option":"All Upfront","on_demand":3390.12,"reserved":1457,"reserved_quantity":1457,"subtraction":1933.12,"discount_rate":0.5702217030665582}
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

onDemandInstanceNum := 2
reservedInstanceNum := 3
fmt.Println(r.ExpectedCost(onDemandInstanceNum, reservedInstanceNum))

{"full_on_demand":{"on_demand":5650.2,"reserved":0,"total":5650.2},"reserved_applied":{"on_demand":2260.08,"reserved":2139,"total":4399.08},"reserved_quantity":2139,"subtraction":1251.12,"discount_rate":0.22142932993522357}
```
