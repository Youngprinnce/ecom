package user

import (
	"database/sql"

	"github.com/youngprinnce/go-ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		user := &types.User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.CreatedAt); err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, nil
}

func (s *Store) GetUserById(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		user := &types.User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.CreatedAt); err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, nil
}

func (s *Store) CreateUser(user *types.User) error {
	_, err := s.db.Exec("INSERT INTO users (email, first_name, last_name, password) VALUES ($1, $2, $3, $4)", user.Email, user.FirstName, user.LastName, user.Password)
	return err
}
