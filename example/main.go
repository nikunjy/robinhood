package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nikunjy/robinhood"
	"golang.org/x/oauth2"
)

func iprintJSON(data interface{}) {}

func printJSON(data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error marshalling data", err)
	}
	fmt.Println(string(bytes))
}

func loginAttempt(
	email string,
	password string,
	mfa string,
) oauth2.TokenSource {
	source := &robinhood.OAuth{
		Username: email,
		Password: password,
	}
	if mfa != "" {
		fmt.Println("Using MFA", mfa)
		source.MFA = mfa
	}
	cacher := &robinhood.CredsCacher{
		Creds: source,
		Path:  "./creds",
	}
	return cacher
}

func doit() {
	username := ""
	password := ""
	if len(os.Args) > 2 {
		username = os.Args[1]
		password = os.Args[2]
	}

	cli, err := robinhood.Dial(context.Background(), loginAttempt(username, password, ""))
	switch err {
	case nil:
		break
	default:
		if robinhood.HasMissingMFA(err) {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter MFA from SMS: ")
			mfa, _ := reader.ReadString('\n')
			cli, err = robinhood.Dial(context.Background(), loginAttempt(username, password, mfa))
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	cli.Debug = true
	ctx := context.Background()

	portfolios, err := cli.GetPortfolios(ctx)
	if err != nil {
		panic(err)
	}
	iprintJSON(portfolios)

	options, err := cli.GetOptionPositions(ctx)
	if err != nil {
		log.Fatal("Error getting options ", err)
	}
	iprintJSON(options)

	it := cli.NewOptionsOrdersIterator()
	for it.HasNext() {
		val, err := it.Next(ctx)
		if err != nil {
			log.Fatal("Error getting option order from iterator ", err)
		}
		for _, option := range val {
			if option.State == robinhood.ORDER_STATE_CANCELLED {
				continue
			}
			printJSON(option)
		}

	}
}

func main() {
	doit()
}
