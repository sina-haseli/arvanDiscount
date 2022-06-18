package models

type RedeemVoucherRequest struct {
	UserID string `json:"user_id"`
	Code   string `json:"code"`
}
