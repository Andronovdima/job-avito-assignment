package models

type Transaction struct {
	ID       int64 `json:"id" `
	Sender   int64 `json:"sender" `
	Receiver int64 `json:"receiver" `
	Sum      int64 `json:"sum" `
}
