package entities

type GameRecordEntity struct {
	ChannelId        int    `db:"ChannelId"`
	Accounts         string `db:"Accounts"`
	KindId           int    `db:"KindId"`
	ServerId         int    `db:"ServerId"`
	CellScore        int    `db:"CellScore"`
	Profit           int    `db:"Profit"`
	CreateTime       string `db:"CreateTime"`
	IsNew            bool   `db:"IsNew"`
	LastLoginTime    string `db:"LastLoginTime"`
	RegisterTime     string `db:"RegisterTime"`
	HistoryCellScore int    `db:"HistoryCellScore"`
	HistoryProfit    int    `db:"HistoryProfit"`
	HistoryGameNum   int    `db:"HistoryGameNum"`
}
