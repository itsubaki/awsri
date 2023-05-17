package cost

const (
	NetAmortizedCost = "NetAmortizedCost"
	NetUnblendedCost = "NetUnblendedCost"
	UnblendedCost    = "UnblendedCost"
	AmortizedCost    = "AmortizedCost"
	BlendedCost      = "BlendedCost"
)

func Metrics() []string {
	return []string{
		NetAmortizedCost,
		NetUnblendedCost,
		UnblendedCost,
		AmortizedCost,
		BlendedCost,
	}
}
