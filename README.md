# Go E-Commerce API

A lightweight, scalable, and modular e-commerce API built with **Golang**. This API provides endpoints for user management, product management, and order management, with support for JWT-based authentication and role-based access control.

---

## Features

- **User Management**:
  - Register a new user.
  - Login and issue JWT tokens.
  - Role-based access control (`admin` and `user` roles).

- **Product Management**:
  - Create, read, update, and delete products.
  - Only accessible by authenticated users with `admin` privileges.

- **Order Management**:
  - Place an order for one or more products.
  - List all orders for a specific user.
  - Cancel an order if it is still in the `pending` status.
  - Update the status of an order (admin only).

- **Authentication**:
  - JWT-based authentication for secure access to protected endpoints.

- **Database**:
  - MySQL database for storing users, products, orders, and order items.

- **Migrations**:
  - Database schema migrations using `golang-migrate`.

---

## Table of Contents

1. [Getting Started](#getting-started)
   - [Prerequisites](#prerequisites)
   - [Installation](#installation)
   - [Configuration](#configuration)
   - [Running the Application](#running-the-application)

2. [API Documentation](#api-documentation)
   - [User Management](#user-management)
   - [Product Management](#product-management)
   - [Order Management](#order-management)

3. [Database Schema](#database-schema)

4. [Migrations](#migrations)

5. [Contributing](#contributing)

6. [License](#license)

---

## Getting Started

### Prerequisites

- **Go** (version 1.20 or higher)
- **MySQL** (version 8.0 or higher)
- **golang-migrate** (for database migrations)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/youngprinnce/ecom.git
   cd ecom
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up the MySQL database:
   - Create a new database named `ecom`.
   - Update the database configuration in the `.env` file (see [Configuration](#configuration)).

4. Run database migrations:
   ```bash
   make migrate-up
   ```

### Configuration

Create a `.env` file in the root directory with the following environment variables:

```env
DB_USER=your_db_user
DB_PASSWD=your_db_password
DB_NET=tcp
DB_ADDR=localhost:3306
DB_NAME=ecom
DB_ALLOW_NATIVE_PASSWORDS=true
DB_PARSE_TIME=true
PORT=8080
JWT_SECRET=your_jwt_secret
JWT_EXPIRE_IN_SECONDS=604800 # 7 days
```

### Running the Application

1. Start the API server:
   ```bash
   make run
   ```

2. The API will be available at `http://localhost:8080`.

---

## API Documentation

### User Management

#### Register a New User
- **Endpoint**: `POST /register`
- **Request Body**:
  ```json
  {
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "password": "password123",
    "role": "user"
  }
  ```
- **Response**:
  ```json
  {
    "message": "User registered successfully"
  }
  ```

#### Login
- **Endpoint**: `POST /login`
- **Request Body**:
  ```json
  {
    "email": "john.doe@example.com",
    "password": "password123"
  }
  ```
- **Response**:
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

### Product Management

#### Create a Product (Admin Only)
- **Endpoint**: `POST /products`
- **Request Body**:
  ```json
  {
    "name": "Product A",
    "description": "A great product",
    "image": "https://example.com/product-a.jpg",
    "price": 19.99,
    "quantity": 100
  }
  ```
- **Response**:
  ```json
  {
    "id": 1,
    "name": "Product A",
    "description": "A great product",
    "image": "https://example.com/product-a.jpg",
    "price": 19.99,
    "quantity": 100,
    "createdAt": "2023-10-01T12:00:00Z"
  }
  ```

#### Get All Products
- **Endpoint**: `GET /products`
- **Response**:
  ```json
  [
    {
      "id": 1,
      "name": "Product A",
      "description": "A great product",
      "image": "https://example.com/product-a.jpg",
      "price": 19.99,
      "quantity": 100,
      "createdAt": "2023-10-01T12:00:00Z"
    }
  ]
  ```

### Order Management

#### Place an Order
- **Endpoint**: `POST /orders`
- **Request Body**:
  ```json
  {
    "items": [
      {
        "productId": 1,
        "quantity": 2
      }
    ]
  }
  ```
- **Response**:
  ```json
  {
    "orderID": 1,
    "totalPrice": 39.98
  }
  ```

#### List Orders for a User
- **Endpoint**: `GET /orders`
- **Response**:
  ```json
  [
    {
      "id": 1,
      "userId": 1,
      "total": 39.98,
      "status": "pending",
      "address": "123 Main St",
      "createdAt": "2023-10-01T12:00:00Z"
    }
  ]
  ```

---

## Database Schema

### Users Table
```sql
CREATE TABLE users (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  firstName VARCHAR(255) NOT NULL,
  lastName VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  role ENUM('admin', 'user') NOT NULL DEFAULT 'user',
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
```

### Products Table
```sql
CREATE TABLE products (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  image VARCHAR(255) NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  quantity INT UNSIGNED NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
```

### Orders Table
```sql
CREATE TABLE orders (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  userId INT UNSIGNED NOT NULL,
  total DECIMAL(10, 2) NOT NULL,
  status VARCHAR(50) NOT NULL,
  address TEXT NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (userId) REFERENCES users(id)
);
```

### Order Items Table
```sql
CREATE TABLE order_items (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  orderId INT UNSIGNED NOT NULL,
  productId INT UNSIGNED NOT NULL,
  quantity INT UNSIGNED NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (orderId) REFERENCES orders(id),
  FOREIGN KEY (productId) REFERENCES products(id)
);
```

---

## Migrations

Database migrations are managed using `golang-migrate`. To apply migrations, run:

```bash
make migrate-up
```

To rollback migrations, run:

```bash
make migrate-down
```

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Let me know if you need further assistance! 😊
