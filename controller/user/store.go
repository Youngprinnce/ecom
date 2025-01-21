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
	return nil, nil
}
