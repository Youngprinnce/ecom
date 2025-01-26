package order

import (
	"context"
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

// CreateOrder creates a new order in the database.
func (s *Store) CreateOrder(order types.Order) (int, error) {
	ctx := context.Background()

	// Insert the order into the database
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO orders (userId, total, status, address)
		VALUES (?, ?, ?, ?)
	`, order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	// Get the ID of the newly created order
	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return int(orderID), nil
}

// CreateOrderItem creates a new order item in the database.
func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	ctx := context.Background()

	// Insert the order item into the database
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO order_items (orderId, productId, quantity, price)
		VALUES (?, ?, ?, ?)
	`, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	if err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}

	return nil
}

// GetOrdersByUserID retrieves all orders for a specific user.
func (s *Store) GetOrdersByUserID(userID int) ([]types.Order, error) {
	ctx := context.Background()

	// Query the database for orders by user ID
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, userId, total, status, address, createdAt
		FROM orders
		WHERE userId = ?
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	// Parse the rows into a slice of Order structs
	orders := make([]types.Order, 0)
	for rows.Next() {
		var order types.Order
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Total,
			&order.Status,
			&order.Address,
			&order.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// GetOrderByID retrieves an order by its ID.
func (s *Store) GetOrderByID(orderID int) (*types.Order, error) {
	ctx := context.Background()

	// Query the database for the order by ID
	row := s.db.QueryRowContext(ctx, `
		SELECT id, userId, total, status, address, createdAt
		FROM orders
		WHERE id = ?
	`, orderID)

	// Parse the row into an Order struct
	var order types.Order
	if err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Total,
		&order.Status,
		&order.Address,
		&order.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to scan order: %w", err)
	}

	return &order, nil
}

// UpdateOrderStatus updates the status of an order.
func (s *Store) UpdateOrderStatus(orderID int, status string) error {
	ctx := context.Background()

	// Update the order status in the database
	_, err := s.db.ExecContext(ctx, `
		UPDATE orders
		SET status = ?
		WHERE id = ?
	`, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

// CancelOrder cancels an order if it is still in the "pending" status.
func (s *Store) CancelOrder(orderID int, userID int) error {
	ctx := context.Background()

	// Cancel the order in the database
	_, err := s.db.ExecContext(ctx, `
		UPDATE orders
		SET status = 'cancelled'
		WHERE id = ? AND userId = ? AND status = 'pending'
	`, orderID, userID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}
