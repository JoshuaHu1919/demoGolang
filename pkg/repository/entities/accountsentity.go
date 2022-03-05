package entities

type AccountsEntity struct {
	Account       string `db:"account"`
	Agent         int    `db:"agent"`
	CreateDate    string `db:"createdate"`
	LastLoginTime string `db:"lastlogintime"`
}
