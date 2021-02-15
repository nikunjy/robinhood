[![Go Reference](https://pkg.go.dev/badge/github.com/nikunjy/robinhood.svg)](https://pkg.go.dev/github.com/nikunjy/robinhood)

This repo is just a fork of [this](https://github.com/andrewstuart/go-robinhood).
I found some errors in the original and I started adding code to this one.

## Usage
1. Please make sure you have MFA enabled on the robinhood account.
2. Cache the oauth token and you can keep on reusing it. Golang oauth2 will auto refresh
```go
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

	ctx := context.Background()
	portfolios, err := cli.GetPortfolios(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(portfolios)
}

func main() {
  	doit()
}
```
