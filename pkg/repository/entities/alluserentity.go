package entities

type GroupByAllUserEntity struct {
	Account   string `db:"Account"`
	CellScore int    `db:"CellScore"`
	WinNum    int    `db:"WinNum"`
	LostNum   int    `db:"LostNum"`
	WinGold   int    `db:"WinGold"`
	LostGold  int    `db:"LostGold"`
}
