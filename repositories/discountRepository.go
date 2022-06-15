package repositories

import (
	"database/sql"
	"discount/models"
	"errors"
	"fmt"
)

type voucherRepository struct {
	db  *sql.DB
	dbq dbQE
	tx  *sql.Tx
}

type dbQE interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

const (
	validateVoucherQuery = "SELECT id, code, usable, amount FROM vouchers WHERE code = $1 limit 1"
	redeemVoucherQuery   = "INSERT INTO redeemed_voucher (user_id, voucher_id, step) VALUES ($1, $2, $3)"
)

var InvalidVoucherCode = errors.New("invalid voucher code")

func NewVoucherRepository(db *sql.DB) *voucherRepository {
	return &voucherRepository{
		db:  db,
		dbq: db,
		tx:  nil,
	}
}

func (v *voucherRepository) FindVoucherByCode(code string) (models.VoucherModel, error) {
	var vm models.VoucherModel
	rows, err := v.dbq.Query(validateVoucherQuery, code)
	if err == sql.ErrNoRows {
		return vm, InvalidVoucherCode
	}

	if err != nil {
		return vm, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error closing rows: ", err)
		}
	}(rows)

	if rows.Next() == false {
		return vm, InvalidVoucherCode
	}

	err = rows.Scan(&vm.ID, &vm.Code, &vm.Usable, &vm.Amount)
	if err != nil {
		return vm, err
	}

	return vm, nil
}

func (v *voucherRepository) InsertIntoRedeemedVoucher(userID, voucherID, step int) error {
	rows, err := v.dbq.Exec(redeemVoucherQuery, userID, voucherID, step)
	if err != nil {
		return err
	}

	ra, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return fmt.Errorf("failed to insert into redeemed voucher, values : %d, %d", userID, voucherID)
	}

	return nil
}

func (v *voucherRepository) RedeemVoucher(userID int, voucher models.VoucherModel, getStep func(voucher models.VoucherModel) (int, error)) error {
	// use transaction to ensure atomicity
	// validate if voucher is already redeemed by user
	// if yes, return error
	// if no, insert into redeemed_voucher
	// if insert failed, return error
	// if insert success, return nil
	return nil
}
