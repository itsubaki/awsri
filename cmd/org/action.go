package org

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/urfave/cli"
)

func Action(c *cli.Context) {
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
			fmt.Printf("list accounts: %v", err)
			os.Exit(1)
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

		return
	}

	if format == "csv" {
		fmt.Println("arn, email, id, joined_method, joined_timestamp, name, status")
		for _, o := range output {
			fmt.Printf("%s, %s, %s, %s, %s, %s, %s\n", *o.Arn, *o.Email, *o.Id, *o.JoinedMethod, *o.JoinedTimestamp, *o.Name, *o.Status)
		}
	}
}
