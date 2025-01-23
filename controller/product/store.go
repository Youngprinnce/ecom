package product

import (
	"database/sql"
	"fmt"

	"github.com/youngprinnce/go-ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProductByID(id int) (*types.Product, error) {
	var p types.Product
	query := "SELECT * FROM products WHERE id = $1"
	if err := s.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt); err != nil {
		return nil, fmt.Errorf("could not get product: %w", err)
	}

	return &p, nil
}

func (s *Store) GetProducts() ([]*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, fmt.Errorf("could not get products: %w", err)
	}
	defer rows.Close()

	var products []*types.Product
	for rows.Next() {
		var p types.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("could not get product: %w", err)
		}
		products = append(products, &p)
	}

	return products, nil
}
