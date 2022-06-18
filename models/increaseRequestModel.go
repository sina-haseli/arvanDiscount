package models

type IncreaseRequestModel struct {
	UserID    string `json:"user_id"`
	VoucherID int    `json:"voucher_id"`
	Amount    int    `json:"amount"`
}
