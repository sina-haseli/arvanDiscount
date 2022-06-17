package repositories

import (
	"context"
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
	validateVoucherQuery             = "SELECT id, code, usable, amount FROM vouchers WHERE code = $1 limit 1"
	findVoucherByCodeAndNotUsedQuery = "SELECT id, code, usable, amount FROM vouchers WHERE code = $1 AND usable > 0 limit 1"
	redeemVoucherQuery               = "INSERT INTO redeemed_voucher (user_id, voucher_id) VALUES ($1, $2)"
	validateFirstTimeRedeemQuery     = "SELECT 1 from redeemed_voucher WHERE voucher_id = $1 AND user_id = $2"
	createVoucherQuery               = "INSERT INTO vouchers (code, amount, usable) VALUES ($1, $2, $3)"
	usersRedeemedVoucherQuery        = "SELECT user_id, voucher_id FROM redeemed_voucher WHERE voucher_id = $1"
	updateVoucherQuery               = "UPDATE vouchers SET usable = usable - 1 WHERE id = $1"
)

var InvalidVoucherCode = errors.New("invalid voucher code")
var VoucherAlreadyUsed = errors.New("voucher already redeemed by user")
var VoucherSoldOut = errors.New("voucher sold out")

func NewVoucherRepository(db *sql.DB) *voucherRepository {
	return &voucherRepository{
		db:  db,
		dbq: db,
		tx:  nil,
	}
}

func (v *voucherRepository) Create(ctx context.Context, rq *models.VoucherRequestModel) (*models.VoucherModel, error) {
	var vm models.VoucherModel
	vm.Code = rq.Code
	vm.Amount = rq.Amount
	vm.Usable = rq.Usable

	_, err := v.dbq.Exec(createVoucherQuery, &vm.Code, &vm.Amount, &vm.Usable)
	if err != nil {
		return nil, err
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

	rows, err := v.dbq.Query(usersRedeemedVoucherQuery, voucher.ID)

	if err == sql.ErrNoRows {
		return nil, InvalidVoucherCode
	}

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// use logger
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

func (v *voucherRepository) FindVoucherByCodeAndNotUsed(ctx context.Context, code string) (models.VoucherModel, error) {
	var vm models.VoucherModel
	//check it
	//	v.db.Prepare(validateVoucherQuery)
	rows, err := v.dbq.Query(findVoucherByCodeAndNotUsedQuery, code)
	if err == sql.ErrNoRows {
		return vm, VoucherSoldOut
	}

	if err != nil {
		return vm, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// use logger
			fmt.Println("Error closing rows: ", err)
		}
	}(rows)
	//read docs
	//v.db.QueryRow(validateVoucherQuery, code)

	if rows.Next() == false {
		return vm, InvalidVoucherCode
	}

	err = rows.Scan(&vm.ID, &vm.Code, &vm.Usable, &vm.Amount)
	if err != nil {
		return vm, err
	}
	ctx.Done()

	return vm, nil
}

func (v *voucherRepository) FindVoucherByCode(ctx context.Context, code string) (models.VoucherModel, error) {
	var vm models.VoucherModel
	//check it
	//	v.db.Prepare(validateVoucherQuery)
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
			// use logger
			fmt.Println("Error closing rows: ", err)
		}
	}(rows)
	//read docs
	//v.db.QueryRow(validateVoucherQuery, code)

	if rows.Next() == false {
		return vm, InvalidVoucherCode
	}

	err = rows.Scan(&vm.ID, &vm.Code, &vm.Usable, &vm.Amount)
	if err != nil {
		return vm, err
	}

	ctx.Done()

	return vm, nil
}

func (v *voucherRepository) UpdateVoucherById(voucherID int) error {
	rows, err := v.dbq.Exec(updateVoucherQuery, voucherID)
	if err != nil {
		return err
	}
	ra, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return fmt.Errorf("failed to update into redeemed voucher, values : %d", voucherID)
	}
	return nil
}

func (v *voucherRepository) InsertIntoRedeemedVoucher(userID, voucherID int) error {
	rows, err := v.dbq.Exec(redeemVoucherQuery, userID, voucherID)
	if err != nil {
		return err
	}

	//use context
	//v.db.ExecContext(v.tx.Context(), redeemVoucherQuery, userID, voucherID, step)

	ra, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return fmt.Errorf("failed to insert into redeemed voucher, values : %d, %d", userID, voucherID)
	}

	return nil
}

func (v *voucherRepository) IsUserRedeemedVoucherBefore(userID, voucherID int) (bool, error) {
	rows, err := v.dbq.Query(validateFirstTimeRedeemQuery, voucherID, userID)
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

func (v *voucherRepository) RedeemVoucher(ctx context.Context, userID int, voucher models.VoucherModel) error {
	trx, err := v.beginTransaction()
	if err != nil {
		return err
	}

	ub, err := trx.IsUserRedeemedVoucherBefore(userID, voucher.ID)
	if ub {
		er := trx.rollbackTransaction()
		if er != nil {
			return er
		}

		return VoucherAlreadyUsed
	}

	err = trx.UpdateVoucherById(voucher.ID)
	if err != nil {
		er := trx.rollbackTransaction()
		if er != nil {
			return er
		}

		return VoucherSoldOut
	}

	err = trx.InsertIntoRedeemedVoucher(userID, voucher.ID)
	if err != nil {
		er := trx.rollbackTransaction()
		if er != nil {
			return er
		}

		return err
	}

	ctx.Done()

	return trx.commitTransaction()
}

func (v *voucherRepository) beginTransaction() (*voucherRepository, error) {
	tx, err := v.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return &voucherRepository{}, err
	}

	return &voucherRepository{tx: tx, dbq: tx}, nil
}

func (v *voucherRepository) commitTransaction() error {
	if v.tx == nil {
		return fmt.Errorf("you cant commit tansaction befor start it")
	}

	return v.tx.Commit()
}

func (v *voucherRepository) rollbackTransaction() error {
	if v.tx == nil {
		return fmt.Errorf("you cant rollback tansaction befor start it")
	}

	return v.tx.Rollback()
}
