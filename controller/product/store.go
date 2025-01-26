package product

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/youngprinnce/go-ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetProductByID retrieves a product by its ID
func (s *Store) GetProductByID(id int) (*types.Product, error) {
	var p types.Product
	query := "SELECT * FROM products WHERE id = ?"
	if err := s.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("could not get product: %w", err)
	}

	return &p, nil
}

// GetProducts retrieves all products
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

// CreateProduct creates a new product
func (s *Store) CreateProduct(p types.CreateProductPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)"
	_, err := s.db.ExecContext(ctx, query, p.Name, p.Description, p.Image, p.Price, p.Quantity)
	if err != nil {
		return fmt.Errorf("could not create product: %w", err)
	}

	return nil
}

// GetProductsByIDs retrieves products by their IDs
func (s *Store) GetProductsByIDs(productIds []int) ([]types.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dynamically generate placeholders for the IN clause
	placeholders := make([]string, len(productIds))
	args := make([]interface{}, len(productIds))
	for i, id := range productIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (%s)", strings.Join(placeholders, ","))
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not get products: %w", err)
	}
	defer rows.Close()

	products := make([]types.Product, 0)
	for rows.Next() {
		var p types.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("could not get product: %w", err)
		}
		products = append(products, p)
	}

	return products, nil
}

// UpdateProduct updates an existing product
func (s *Store) UpdateProduct(p types.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE products SET name = ?, description = ?, image = ?, price = ?, quantity = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, p.Name, p.Description, p.Image, p.Price, p.Quantity, p.ID)
	if err != nil {
		return fmt.Errorf("could not update product: %w", err)
	}

	return nil
}

// DeleteProduct deletes a product by its ID
func (s *Store) DeleteProduct(productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM products WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("could not delete product: %w", err)
	}

	return nil
}
