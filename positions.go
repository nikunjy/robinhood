package robinhood

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type CustomTime struct {
	time.Time
}

func (c *CustomTime) MarshalJSON() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c.Time)
}

func (c *CustomTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.Replace(str, `"`, "", -1)

	if val, err := time.Parse("2006-01-02", str); err == nil {
		c.Time = val
		return nil
	}
	return errors.New("no time formats matched")
}

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
	ExpirationDate CustomTime   `json:"expiration_date"`
	StrikePrice    float64      `json:"strike_price,string"`
	OptionType     string       `json:"option_type"`
}

type getPositionConfig struct {
	nonZero bool
}

func (c *getPositionConfig) params() PositionParams {
	params := PositionParams{}
	if c.nonZero {
		params.NonZero = true
	}
	return params
}

func newDefaultOptionsConfig() *getPositionConfig {
	return &getPositionConfig{
		nonZero: false,
	}
}

type GetPositionsParamsOptions func(*getPositionConfig)

func ExcludeZeroPositions() GetPositionsParamsOptions {
	return func(cfg *getPositionConfig) {
		cfg.nonZero = true
	}
}

// GetPositions returns all the positions associated with an account.
func (c *Client) GetOptionPositions(ctx context.Context, opts ...GetPositionsParamsOptions) ([]OptionPostion, error) {
	cfg := newDefaultOptionsConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return c.GetOptionPositionsParams(ctx, cfg.params())
}

// GetPositions returns all the positions associated with an account.
func (c *Client) GetPositions(ctx context.Context, opts ...GetPositionsParamsOptions) ([]Position, error) {
	cfg := newDefaultOptionsConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return c.GetPositionsParams(ctx, cfg.params())
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
	if err := c.GetAndDecode(ctx, u.String(), &r); err != nil {
		return nil, errors.Wrap(err, "error getting and decoding options")
	}
	return r.Results, nil
}
