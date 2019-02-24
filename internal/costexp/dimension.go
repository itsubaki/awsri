package costexp

import "strings"

type GetUsageQuantityInputFunc func(all []string) *GetUsageQuantityInput

func NewGetUsageQuantityInput() []GetUsageQuantityInputFunc {
	return []GetUsageQuantityInputFunc{
		NewComputeGetUsageQuantityInput,
		NewCacheGetUsageQuantityInput,
		NewDatabaseGetUsageQuantityInput,
	}
}

func NewComputeGetUsageQuantityInput(all []string) *GetUsageQuantityInput {
	usageType := []string{}
	for i := range all {
		if !strings.Contains(all[i], "BoxUsage") {
			continue
		}
		usageType = append(usageType, all[i])
	}

	return &GetUsageQuantityInput{
		Dimension: "PLATFORM",
		UsageType: usageType,
	}
}

func NewCacheGetUsageQuantityInput(all []string) *GetUsageQuantityInput {
	usageType := []string{}
	for i := range all {
		if !strings.Contains(all[i], "NodeUsage") {
			continue
		}
		usageType = append(usageType, all[i])
	}

	return &GetUsageQuantityInput{
		Dimension: "CACHE_ENGINE",
		UsageType: usageType,
	}
}

func NewDatabaseGetUsageQuantityInput(all []string) *GetUsageQuantityInput {
	usageType := []string{}
	for i := range all {
		if !strings.Contains(all[i], "InstanceUsage") &&
			!strings.Contains(all[i], "Multi-AZUsage") {
			continue
		}
		usageType = append(usageType, all[i])
	}

	return &GetUsageQuantityInput{
		Dimension: "DATABASE_ENGINE",
		UsageType: usageType,
	}
}
