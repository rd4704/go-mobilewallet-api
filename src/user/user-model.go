package user

import (
	"database/sql"
	"fmt"
)

// User model
type User struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	WalletId    string  `json:"walletId"`
	Description string  `json:"description"`
	Balance     float32 `json:"balance"`
}

// GetUser get user by id
func (u *User) GetUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT u.id, u.name, u.email, w.id as walletId, w.description, w.balance FROM users u INNER JOIN wallets w on w.userId WHERE u.id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.ID, &u.Name, &u.Email, &u.WalletId, &u.Description, &u.Balance)
}

// GetUsers query all users
func GetUsers(db *sql.DB, start, count int) ([]User, error) {
	statement := fmt.Sprintf("SELECT u.id, u.name, u.email, w.id as walletId, w.description, w.balance FROM users u INNER JOIN wallets w on w.userId = u.id LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.WalletId, &u.Description, &u.Balance); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
