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
	CreatedAt   string  `json:"createdAt"`
}

// GetTransfer by id
func (t *Transfer) GetTransfer(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, description, fromWallet, toWallet, amount, createdAt from transfers WHERE id=%d", t.ID)
	return db.QueryRow(statement).Scan(&t.ID, &t.Description, &t.FromWallet, &t.ToWallet, &t.Amount, &t.CreatedAt)
}

// GetTransfers of a user
func getUserTransfers(db *sql.DB, userId, start, count int) ([]Transfer, error) {
	statement := fmt.Sprintf(`SELECT t.id as transferId, t.description, t.fromWallet, t.toWallet, CASE
	WHEN t.fromWallet = w.id THEN t.amount * -1
	ELSE  t.amount
	END as amount,
	t.createdAt
	from transfers t
	inner join wallets w on w.id = t.fromWallet or w.id = t.toWallet
	inner join transfers u on u.id = w.userId
	WHERE u.id = %d LIMIT %d OFFSET %d`, userId, count, start)

	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	transfers := []Transfer{}

	for rows.Next() {
		var t Transfer
		if err := rows.Scan(&t.ID, &t.Description, &t.FromWallet, &t.ToWallet, &t.Amount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transfers = append(transfers, t)
	}

	return transfers, nil
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
