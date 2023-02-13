package rebalanceaccount

type RebalanceAccountRequest struct {
	AccountsID []int64 `json:"accounts"`
	Limit      int     `json:"limit"`
}

type AccountsToRebalanceModel struct {
	AccountID  int64   `db:"account_id"`
	Profit     float32 `db:"profit"`
	ProfitRate float32 `db:"profit_rate"`
	Volume     float32 `db:"volume"`
	Deals      int     `db:"deals"`
	WeekNum    int     `db:"week_num"`
	PondProfit float32 `db:"pond_profit"`
}
