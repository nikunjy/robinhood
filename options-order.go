package robinhood

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// OptionsOrderOpts encapsulates common Options order choices
type OptionsOrderOpts struct {
	Quantity    float64
	Price       float64
	Direction   OptionDirection
	TimeInForce TimeInForce
	Type        OrderType
	Side        OrderSide
}

// optionInput is the input object to the RobinHood API
type optionInput struct {
	Account                string          `json:"account"`
	Direction              OptionDirection `json:"direction"`
	Legs                   []Leg           `json:"legs"`
	OverrideDayTradeChecks bool            `json:"override_day_trade_checks"`
	OverrideDtbpChecks     bool            `json:"override_dtbp_checks"`
	Price                  float64         `json:"price,string"`
	Quantity               float64         `json:"quantity,string"`
	RefID                  string          `json:"ref_id"`
	TimeInForce            TimeInForce     `json:"time_in_force"`
	Trigger                string          `json:"trigger"`
	Type                   OrderType       `json:"type"`
}

// A Leg is a single option contract that will be purchased as part of a single
// order. Transactions! Lower Risk!
type Leg struct {
	Option         string    `json:"option"`
	PositionEffect string    `json:"position_effect"`
	RatioQuantity  float64   `json:"ratio_quantity,string"`
	Side           OrderSide `json:"side"`
}

// OrderOptions places a new order for options. Cancellation of the
// context.Context will cancel the _http request_, never the order itself if it
// has already been created.
func (c *Client) OrderOptions(ctx context.Context, q *OptionInstrument, o OptionsOrderOpts) (json.RawMessage, error) {
	b := optionInput{
		Account:     c.Account.URL,
		Direction:   o.Direction,
		TimeInForce: o.TimeInForce,
		Legs: []Leg{{
			Option:         q.URL,
			RatioQuantity:  1,
			Side:           o.Side,
			PositionEffect: "open",
		}},
		Trigger:  "immediate",
		Type:     o.Type,
		Quantity: o.Quantity,
		Price:    o.Price,
		RefID:    uuid.New().String(),
	}

	if o.Side != Buy {
		b.Legs[0].PositionEffect = "close"
	}

	bs, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", EPOptions+"orders/", bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var out json.RawMessage
	err = c.DoAndDecode(ctx, req, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type OptionOrderState string

const (
	ORDER_STATE_FILLED    OptionOrderState = "filled"
	ORDER_STATE_CANCELLED OptionOrderState = "cancelled"
)

type OptionOrder struct {
	CancelURL        string    `json:"cancel_url"`
	CanceledQuantity float64   `json:"canceled_quantity,string"`
	CreatedAt        time.Time `json:"created_at,string"`
	Direction        string    `json:"direction"`
	ID               string    `json:"id"`
	Legs             []struct {
		Executions []struct {
			ID             string
			Price          float64   `json:"price,string"`
			Quantity       float64   `json:"quantity,string"`
			SettlementDate string    `json:"settlement_date"`
			Timestamp      time.Time `json:"timestamp,string"`
		} `json:"executions"`
		ID             string  `json:"id"`
		Instrument     string  `json:"option"`
		PositionEffect string  `json:"position_effect"`
		RatioQuantity  float64 `json:"ratio_quantity"`
		Side           string  `json:"side"`
	} `json:"legs"`
	PendingQuantity   float64          `json:"pending_quantity,string"`
	Premium           float64          `json:"premium,string"`
	ProcessedPremium  float64          `json:"processed_premium,string"`
	Price             float64          `json:"price,string"`
	ProcessedQuantity float64          `json:"processed_quantity,string"`
	Quantity          float64          `json:"quantity,string"`
	RefID             string           `json:"ref_id"`
	State             OptionOrderState `json:"state"`
	Trigger           string           `json:"trigger"`
	Type              string           `json:"type"`
	UpdatedAt         time.Time        `json:"updated_at,string"`
	ChainID           string           `json:"chain_id"`
	ChainSymbol       string           `json:"chain_symbol"`
	ClosingStrategy   string           `json:"closing_strategy"`
}

func (o OptionOrder) IsClosingOrder() bool {
	return o.ClosingStrategy != ""
}

type OptionsOrdersIterator interface {
	HasNext() bool
	Next(ctx context.Context) ([]OptionOrder, error)
}

type optionsOrdersIterators struct {
	c    *Client
	next string
}

func (o *optionsOrdersIterators) HasNext() bool {
	return o.next != ""
}

func (o *optionsOrdersIterators) Next(ctx context.Context) ([]OptionOrder, error) {
	var tmp struct {
		Results []OptionOrder
		Next    string
	}
	err := o.c.GetAndDecode(ctx, o.next, &tmp)
	if err != nil {
		return nil, err
	}
	o.next = tmp.Next
	return tmp.Results, nil
}

// NewOptionsOrdersIterator returns an iterator which will return all
// the option orders ever seen.
func (c *Client) NewOptionsOrdersIterator() OptionsOrdersIterator {
	return &optionsOrdersIterators{
		next: EPOptions + "orders/",
		c:    c,
	}
}

// GetOptionsOrders returns all outstanding options orders
func (c *Client) GetOptionsOrders(ctx context.Context) ([]OptionOrder, error) {
	var tmp struct {
		Results []OptionOrder
		Next    string
	}
	err := c.GetAndDecode(ctx, EPOptions+"orders/", &tmp)
	if err != nil {
		return nil, err
	}
	return tmp.Results, nil
}
