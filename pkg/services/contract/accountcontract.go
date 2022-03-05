package contract

import "time"

type AccountsContract struct {
	Platform      string    `json:"platform"`
	Account       string    `json:"account"`
	Agent         int       `json:"agent"`
	RegisterTime  string    `json:"registertime"`
	LastLoginTime string    `json:"lastlogintime"`
	StatisticDate string    `json:"statisticdate"`
	ExpireTime    time.Time `json:"expiretime"`
}
