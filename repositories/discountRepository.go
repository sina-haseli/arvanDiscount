package repositories

import (
	"context"
	"database/sql"
	"discount/models"
	"errors"
	"fmt"
)

type voucherRepository struct {
	db *sql.DB
}

const (
	validateVoucherQuery         = "SELECT id, code, usable, amount FROM vouchers WHERE code = $1 limit 1"
	redeemVoucherQuery           = "INSERT INTO redeemed_voucher (user_id, voucher_id) VALUES ($1, $2)"
	validateFirstTimeRedeemQuery = "SELECT 1 from redeemed_voucher WHERE voucher_id = $1 AND user_id = $2"
	createVoucherQuery           = "INSERT INTO vouchers (code, amount, usable) VALUES ($1, $2, $3)"
	usersRedeemedVoucherQuery    = "SELECT user_id, voucher_id FROM redeemed_voucher WHERE voucher_id = $1"
	updateVoucherQuery           = "UPDATE vouchers SET usable = usable - 1 WHERE id = $1 and usable > 0"
)

var InvalidVoucherCode = errors.New("invalid voucher code")
var VoucherExist = errors.New("voucher exist")
var VoucherAlreadyUsed = errors.New("voucher already redeemed by user")
var VoucherSoldOut = errors.New("voucher sold out")

func NewVoucherRepository(db *sql.DB) *voucherRepository {
	return &voucherRepository{
		db: db,
	}
}

func (v *voucherRepository) Create(ctx context.Context, rq *models.VoucherRequestModel) (*models.VoucherModel, error) {
	vm := models.VoucherModel{
		Code:   rq.Code,
		Amount: rq.Amount,
		Usable: rq.Usable,
	}

	_, err := v.db.ExecContext(ctx, createVoucherQuery, &vm.Code, &vm.Amount, &vm.Usable)
	if err != nil {
		return nil, VoucherExist
	}

	result, err := v.FindVoucherByCode(ctx, vm.Code)
	if err != nil {
		return nil, err
	}

	vm.ID = result.ID
	return &vm, nil
}

func (v *voucherRepository) GetVoucherCodeUsed(ctx context.Context, code string) (*[]models.RedeemVoucherRequest, error) {

	var rvr models.RedeemVoucherRequest

	voucher, err := v.FindVoucherByCode(ctx, code)
	if err != nil {
		return nil, InvalidVoucherCode
	}

	rows, err := v.db.QueryContext(ctx, usersRedeemedVoucherQuery, voucher.ID)

	if err == sql.ErrNoRows {
		return nil, InvalidVoucherCode
	}

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error closing rows: ", err)
		}
	}(rows)

	if rows.Next() == false {
		return nil, InvalidVoucherCode
	}

	var users []models.RedeemVoucherRequest
	for rows.Next() {
		err = rows.Scan(&rvr.UserID, &rvr.Code)
		if err != nil {
			return nil, err
		}
		users = append(users, rvr)
	}

	return &users, nil

}

func (v *voucherRepository) FindVoucherByCode(ctx context.Context, code string) (models.VoucherModel, error) {
	var vm models.VoucherModel

	rows, err := v.db.QueryContext(ctx, validateVoucherQuery, code)
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

func (v *voucherRepository) updateVoucherById(ctx context.Context, tx *sql.Tx, voucherID int) (bool, error) {
	rows, err := tx.ExecContext(ctx, updateVoucherQuery, voucherID)
	if err != nil {
		return false, err
	}
	ra, err := rows.RowsAffected()
	if err != nil {
		return false, err
	}

	if ra == 0 {
		return false, nil
	}
	return true, nil
}

func (v *voucherRepository) insertIntoRedeemedVoucher(ctx context.Context, tx *sql.Tx, userID string, voucherID int) error {
	rows, err := tx.ExecContext(ctx, redeemVoucherQuery, userID, voucherID)
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

func (v *voucherRepository) RedeemVoucher(ctx context.Context, userID string, voucherID int) error {
	trx, err := v.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	isRedeemed, err := v.isUserRedeemedVoucherBefore(ctx, trx, userID, voucherID)
	if err != nil {
		if er := trx.Rollback(); er != nil {
			return er
		}

		return err
	}

	if isRedeemed {
		if er := trx.Rollback(); er != nil {
			return er
		}

		return VoucherAlreadyUsed
	}

	vUpdated, err := v.updateVoucherById(ctx, trx, voucherID)
	if err != nil {
		if er := trx.Rollback(); er != nil {
			return er
		}

		return err
	}

	if !vUpdated {
		if er := trx.Rollback(); er != nil {
			return er
		}

		return VoucherSoldOut
	}

	err = v.insertIntoRedeemedVoucher(ctx, trx, userID, voucherID)
	if err != nil {
		if er := trx.Rollback(); er != nil {
			return er
		}

		return err
	}

	return trx.Commit()
}

func (v *voucherRepository) isUserRedeemedVoucherBefore(ctx context.Context, tx *sql.Tx, userID string, voucherID int) (bool, error) {
	rows, err := tx.QueryContext(ctx, validateFirstTimeRedeemQuery, voucherID, userID)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		fmt.Println("failed to check voucher redeemed before or not:", err)
		return true, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("failed to close rows: ", err)
		}
	}(rows)

	if rows.Next() {
		return true, nil
	}

	return false, nil
}
