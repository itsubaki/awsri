package usage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type Account struct {
	ID          string
	Description string
}

func FetchLinkedAccount(start, end string) ([]Account, error) {
	input := costexplorer.GetDimensionValuesInput{
		Dimension: aws.String("LINKED_ACCOUNT"),
		TimePeriod: &costexplorer.DateInterval{
			Start: &start,
			End:   &end,
		},
	}

	c := costexplorer.New(session.Must(session.NewSession()))
	val, err := c.GetDimensionValues(&input)
	if err != nil {
		return []Account{}, fmt.Errorf("get dimension values: %v", err)
	}

	out := make([]Account, 0)
	for _, v := range val.DimensionValues {
		out = append(out, Account{
			ID:          *v.Value,
			Description: *v.Attributes["description"],
		})
	}

	return out, nil
}
