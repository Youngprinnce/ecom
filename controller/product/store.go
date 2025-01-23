package product

import (
	"database/sql"
	"fmt"
	"context"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM products")
	if err != nil {
		return nil, fmt.Errorf("could not get products: %w", err)
	}
	defer rows.Close()

	products := make([]*types.Product, 0)
	for rows.Next() {
		var p types.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("could not get product: %w", err)
		}
		products = append(products, &p)
	}

	return products, nil
}

func (s *Store) CreateProduct(p types.CreateProductPayload) error {
	fmt.Println("CreateProductPayload: ", p)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	query := "INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)"
	if _, err := s.db.ExecContext(ctx, query, p.Name, p.Description, p.Image, p.Price, p.Quantity); err != nil {
		return fmt.Errorf("could not create product: %w", err)
	}

	return nil
}
