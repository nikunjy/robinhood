package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	path string,
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
		Path:  path,
	}
	return cacher
}

func doit() {
	userName := flag.String("user", "", "username")
	password := flag.String("password", "", "Password for the account")
	var path string
	flag.StringVar(&path, "path", "", "Path for creds to be stored")
	flag.Parse()

	if path == "" {
		if *userName == "" || *password == "" {
			log.Fatal("Need to provide either creds file or username/password")
		}
		f, err := ioutil.TempFile(os.TempDir(), "creds")
		if err != nil {
			log.Fatalf("Error creating tmp file %v", err)
		}
		path = f.Name()
		f.Close()
		log.Printf("Picking creds path %s\n", path)
	}
	cli, err := robinhood.Dial(context.Background(), loginAttempt(*userName, *password, "", path))
	switch err {
	case nil:
		break
	default:
		if robinhood.HasMissingMFA(err) {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter MFA from SMS: ")
			mfa, _ := reader.ReadString('\n')
			cli, err = robinhood.Dial(context.Background(), loginAttempt(*userName, *password, mfa, path))
			if err != nil {
				log.Fatalf("Error dialing robinhood %v", err)
			}
		} else {
			log.Fatalf("Error getting creds from robinhood %v", err)
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
				//continue
			}
			printJSON(option)
		}

	}
}

func main() {
	doit()
}
