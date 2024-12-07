//go:generate mockery --name=PismoRepo --output=../mocks
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	pismoRepo struct {
		db *sqlx.DB
	}

	PismoRepo interface {
		GetAccountByDocumentNo(ctx context.Context, document_number string) (isExists bool, err error)
		CreateAccount(ctx context.Context, document_number string) (err error)
		GetAccountByAccountID(ctx context.Context, account_id int) (account *Account, err error)
		CreateTransaction(ctx context.Context, txn Transaction) (err error)
	}
)

// NewPismoRepo configures and returns the object for PismoRepo
func NewPismoRepo(db *sqlx.DB) PismoRepo {
	return &pismoRepo{
		db,
	}
}

// GetAccountByDocumentNo retrives the account for given document_number, if account exists it will return true
func (p *pismoRepo) GetAccountByDocumentNo(ctx context.Context, docNo string) (isExists bool, err error) {
	err = p.db.GetContext(
		ctx,
		&isExists,
		"SELECT EXISTS ( SELECT 1 FROM accounts WHERE document_number = $1 limit 1)",
		docNo,
	)
	if err != nil {
		return false, fmt.Errorf("failed to query account: %w", err)
	}

	return
}

// CreateAccount creates new account record in accounts table
func (p *pismoRepo) CreateAccount(ctx context.Context, docNo string) (err error) {
	_, err = p.db.ExecContext(
		ctx,
		"INSERT INTO accounts (document_number) VALUES ($1)",
		docNo,
	)

	if err != nil {
		return fmt.Errorf("failed to insert account: %w", err)
	}

	return
}

// GetAccountByAccountID retrives account for given account_id
func (p *pismoRepo) GetAccountByAccountID(ctx context.Context, accID int) (*Account, error) {
	var acc Account
	err := p.db.GetContext(
		ctx,
		&acc,
		"SELECT account_id, document_number FROM accounts WHERE account_id = $1",
		accID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return &acc, fmt.Errorf("failed to query account: %w", err)
	}

	return &acc, nil
}

// CreateTransaction creates new record for in transactions table
func (p *pismoRepo) CreateTransaction(ctx context.Context, txn Transaction) (err error) {
	_, err = p.db.NamedExecContext(ctx,
		`INSERT INTO transactions 
			(account_id, operation_type_id, amount) 
		VALUES 
			(:account_id, :operation_type_id, :amount)
		`,
		txn,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction record: %w", err)
	}

	return
}
