package repository

import "time"

type Account struct {
	AccountID  int    `db:"account_id"`
	DocumentNo string `db:"document_number"`
}

type Transaction struct {
	AccountID       int       `db:"account_id"`
	OperationTypeID int       `db:"operation_type_id"`
	Amount          float64   `db:"amount"`
	EventDate       time.Time `db:"event_date"`
}
