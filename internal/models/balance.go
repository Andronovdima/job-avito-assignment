package models

type Balance struct {
	ID       int64 `json:"id" `
	UserID   int64 `json:"user_id" `
	Sum      int64 `json:"sum" `
}

