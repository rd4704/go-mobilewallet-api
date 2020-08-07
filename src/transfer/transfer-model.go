package transfer

import (
	"database/sql"
	"fmt"
	"log"
)

// User model
type Transfer struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	FromWallet  int     `json:"fromWallet"`
	ToWallet    int     `json:"toWallet"`
	Amount      float32 `json:"amount"`
}

// GetUser get user by id
func (t *Transfer) GetTransfer(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, description, fromWallet, toWallet, amount from transfers WHERE id=%d", t.ID)
	return db.QueryRow(statement).Scan(&t.ID, &t.Description, &t.FromWallet, &t.ToWallet, &t.Amount)
}

/**
Make Transfer

- Deduce balance fromWallet
- Add amount toWallet
- Insert transfer transaction
*/
func (t *Transfer) MakeTransfer(db *sql.DB) error {

	tx, err := db.Begin()
	handleError(err)

	// Update the fromWallet
	fromWalletStmt := fmt.Sprintf("UPDATE wallets SET balance = (balance - %f) WHERE id=%d", t.Amount, t.FromWallet)
	res, err := tx.Exec(fromWalletStmt)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	rows, err := res.RowsAffected()
	handleError(err)

	if rows > 0 {
		log.Println("Updated fromWallet balance.")
	}

	// Update the fromWallet
	toWalletStmt := fmt.Sprintf("UPDATE wallets SET balance = (balance + %f) WHERE id=%d", t.Amount, t.ToWallet)
	res, err = tx.Exec(toWalletStmt)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	rows, err = res.RowsAffected()
	handleError(err)

	if rows > 0 {
		log.Println("Updated toWallet balance.")
	}

	// insert record into table2, referencing the first record from table1
	transferStmt := fmt.Sprintf("INSERT INTO transfers(description, fromWallet, toWallet, amount) VALUES('%s', %d, %d, %f)", t.Description, t.FromWallet, t.ToWallet, t.Amount)
	res, err = tx.Exec(transferStmt)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	rows, err = res.RowsAffected()
	handleError(err)

	if rows > 0 {
		log.Println("Inserted transfer record.")
	}

	// commit the transaction
	handleError(tx.Commit())

	log.Println("[MakeTransfer]: Transaction Completed.")

	return nil
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
