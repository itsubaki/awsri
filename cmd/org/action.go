package org

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	format := c.String("format")

	svc := organizations.New(session.Must(session.NewSession()))
	output := make([]organizations.Account, 0)
	token := ""
	for {
		input := &organizations.ListAccountsInput{}
		if token != "" {
			input.NextToken = &token
		}

		result, err := svc.ListAccounts(input)
		if err != nil {
			return fmt.Errorf("list accounts: %v", err)
		}

		for i := range result.Accounts {
			output = append(output, *result.Accounts[i])
		}

		if result.NextToken == nil {
			break
		}

		token = *result.NextToken
	}

	if format == "json" {
		for _, o := range output {
			b, err := json.Marshal(o)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(b))
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("arn, email, id, joined_method, joined_timestamp, name, status")
		for _, o := range output {
			fmt.Printf("%s, %s, %s, %s, %s, %s, %s\n", *o.Arn, *o.Email, *o.Id, *o.JoinedMethod, *o.JoinedTimestamp, *o.Name, *o.Status)
		}

		return nil
	}

	return nil
}
