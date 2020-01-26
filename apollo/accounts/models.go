package accounts

import "time"

type Account struct {
	Id        string    `json:"account_id"`
	UserId    string    `json:"user_id"`
	AssetId   string    `json:"asset_id"`
	Balance   int64     `json:"balance"`
	Holds     int64     `json:"holds"`
	CreatedAt time.Time `json:"created_at"`
}
