package robinhood

import (
	"context"
	"net/url"
	"time"
)

type Position struct {
	Meta
	Account                 string  `json:"account"`
	AverageBuyPrice         float64 `json:"average_buy_price,string"`
	Instrument              string  `json:"instrument"`
	IntradayAverageBuyPrice float64 `json:"intraday_average_buy_price,string"`
	IntradayQuantity        float64 `json:"intraday_quantity,string"`
	Quantity                float64 `json:"quantity,string"`
	SharesHeldForBuys       float64 `json:"shares_held_for_buys,string"`
	SharesHeldForSells      float64 `json:"shares_held_for_sells,string"`
}

type OptionPostion struct {
	Chain                    string        `json:"chain"`
	AverageOpenPrice         float64       `json:"average_open_price,string"`
	Symbol                   string        `json:"symbol"`
	Quantity                 float64       `json:"quantity,string"`
	Direction                string        `json:"direction"`
	IntradayDirection        string        `json:"intraday_direction"`
	TradeValueMultiplier     string        `json:"trade_value_multiplier"`
	Account                  string        `json:"account"`
	Strategy                 string        `json:"strategy"`
	Legs                     []LegPosition `json:"legs"`
	IntradayQuantity         string        `json:"intraday_quantity"`
	UpdatedAt                time.Time     `json:"updated_at,string"`
	Id                       string        `json:"id"`
	IntradayAverageOpenPrice float64       `json:"intraday_average_open_price,string"`
	CreatedAt                time.Time     `json:"created_at,string"`
}

type PositionType string

const (
	Short PositionType = "short"
	Long               = "long"
)

type LegPosition struct {
	Id             string       `json:"id"`
	Position       string       `json:"position"`
	PositionType   PositionType `json:"position_type"`
	Option         string       `json:"option"`
	RatioQuantity  int          `json:"ratio_quantity"`
	ExpirationDate time.Time    `json:"expiration_date"`
	StrikePrice    float64      `json:"strike_price,string"`
	OptionType     string       `json:"option_type"`
}

// GetPositions returns all the positions associated with an account.
func (c *Client) GetOptionPositions(ctx context.Context) ([]OptionPostion, error) {
	return c.GetOptionPositionsParams(ctx, PositionParams{NonZero: true})
}

// GetPositions returns all the positions associated with an account.
func (c *Client) GetPositions(ctx context.Context) ([]Position, error) {
	return c.GetPositionsParams(ctx, PositionParams{NonZero: true})
}

// PositionParams encapsulates parameters known to the RobinHood positions API
// endpoint.
type PositionParams struct {
	NonZero bool
}

// Encode returns the query string associated with the requested parameters
func (p PositionParams) encode() string {
	v := url.Values{}
	if p.NonZero {
		v.Set("nonzero", "True")
	}
	return v.Encode()
}

// GetPositionsParams returns all the positions associated with a count, but
// passes the encoded PositionsParams object along to the RobinHood API as part
// of the query string.
func (c *Client) GetPositionsParams(ctx context.Context, p PositionParams) ([]Position, error) {
	u, err := url.Parse(EPPositions)
	if err != nil {
		return nil, err
	}
	u.RawQuery = p.encode()

	var r struct{ Results []Position }
	return r.Results, c.GetAndDecode(ctx, u.String(), &r)
}

// GetOptionPositionsParams returns all the positions associated with a count, but
// passes the encoded PositionsParams object along to the RobinHood API as part
// of the query string.
func (c *Client) GetOptionPositionsParams(ctx context.Context, p PositionParams) ([]OptionPostion, error) {
	u, err := url.Parse(EPOptions + "aggregate_positions/")
	if err != nil {
		return nil, err
	}
	u.RawQuery = p.encode()
	var r struct{ Results []OptionPostion }
	return r.Results, c.GetAndDecode(ctx, u.String(), &r)
}
