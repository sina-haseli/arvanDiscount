package models

type VoucherRequestModel struct {
	Code   string `json:"code"`
	Usable int    `json:"usable"`
	Amount int    `json:"amount"`
}
