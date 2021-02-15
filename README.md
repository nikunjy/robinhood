[![Go Reference](https://pkg.go.dev/badge/github.com/nikunjy/robinhood.svg)](https://pkg.go.dev/github.com/nikunjy/robinhood)

This repo is just a fork [of](https://github.com/andrewstuart/go-robinhood).
I found some errors in the original and I started adding code to this one.

# Robinhood the rich and feeding the poor, now automated

> Even though robinhood makes me poor

## Notice

If you have used this library before, and use credential caching, you will need
to remove any credential cache and rebuild if you experience errors.

## General usage

```go
cli, err := robinhood.Dial(&robinhood.OAuth{
  Username: "andrewstuart",
  Password: "mypasswordissecure",
})

// err

i, err := cli.GetInstrumentForSymbol("SPY")

// err

o, err := cli.Order(i, robinhood.OrderOpts{
  Price: 100.0,
  Side: robinhood.Buy,
  Quantity: 1,
})

// err

time.Sleep(5*time.Second) //Let me think about it some more...

// Ah crap, I need to buy groceries.

err := o.Cancel()

if err != nil {
  // Oh well
}
```
