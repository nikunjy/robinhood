package robinhood

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptionsOrder(t *testing.T) {
	data := []byte(`{
        "cancel_url": null,
        "canceled_quantity": "0.00000",
        "created_at": "2018-09-06T14:48:20.305171Z",
        "direction": "credit",
        "id": "f43175ba-3191-4aa4-a024-780e1e2beecd",
        "legs": [
          {
            "executions": [
              {
                "id": "b854d39c-5554-47b9-b32b-a5352ab5955e",
                "price": "0.45000000",
                "quantity": "2.00000",
                "settlement_date": "2018-09-07",
                "timestamp": "2018-09-06T14:48:29.576000Z"
              }
            ],
            "id": "cadaa42-assdsb0-4sdxd-as4f-959soso7666aa24",
            "option": "https://api.robinhood.com/options/instruments/caadas-cb8xcx-4ooopo-lll-9bsdsds/",
            "position_effect": "close",
            "ratio_quantity": 1,
            "side": "sell"
          }
        ],
        "pending_quantity": "0.00000",
        "premium": "45.00000000",
        "processed_premium": "90.00000000000000000",
        "price": "0.45000000",
        "processed_quantity": "2.00000",
        "quantity": "2.00000",
        "ref_id": "CsSSDS22B5C9-7ALC-48FXC-BAS0-B26A9AA0009A4",
        "state": "filled",
        "time_in_force": "gfd",
        "trigger": "immediate",
        "type": "limit",
        "updated_at": "2019-01-01T14:48:29.835760Z",
        "chain_id": "e66adasz029-d2326-4572-87a0-b14232013c08bf",
        "chain_symbol": "AMD",
        "response_category": null,
        "opening_strategy": null,
        "closing_strategy": "long_put",
        "stop_price": null
      }`)
	var order OptionOrder
	require.NoError(t, json.Unmarshal(data, &order))
	require.EqualValues(t, order.CanceledQuantity, 0)
	require.Len(t, order.Legs, 1)
	require.Len(t, order.Legs[0].Executions, 1)
	execution := order.Legs[0].Executions[0]
	require.EqualValues(t, execution.Price, 0.45)
	require.EqualValues(t, execution.ID, "b854d39c-5554-47b9-b32b-a5352ab5955e")
	require.EqualValues(t, execution.SettlementDate, "2018-09-07")
	require.EqualValues(t, order.Legs[0].Instrument, "https://api.robinhood.com/options/instruments/caadas-cb8xcx-4ooopo-lll-9bsdsds/")
	require.EqualValues(t, order.Legs[0].ID, "cadaa42-assdsb0-4sdxd-as4f-959soso7666aa24")

}
