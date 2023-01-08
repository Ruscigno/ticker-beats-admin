package rebalanceaccount

type RebalanceAccountRequest struct {
	AccoundID int64 `json:"account_id"`
	Limit     int   `json:"limit"`
}

type AccountsToRebalanceModel struct {
	AccountID  int64   `json:"account_id"`
	Profit     float32 `json:"profit"`
	ProfitRate float32 `json:"profit_rate"`
	Volume     float32 `json:"volume"`
	Deals      int     `json:"deals"`
	WeekNum    int     `json:"week_num"`
	PondProfit float32 `json:"pond_profit"`
}
