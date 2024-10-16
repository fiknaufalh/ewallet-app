# E-Wallet Application

This is a simple e-wallet application built with Go, using MySQL as the database. It allows users to create accounts, top up their balance, and make withdrawals.

## Features

- Create user accounts
- View user information and balance
- Top up user balance
- Withdraw from user balance

## Prerequisites

- Go 1.20 or higher
- Docker and Docker Compose
- Postman (for testing API endpoints)

## Project Structure

```
ewallet-app/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── controller/
│   │   └── wallet_controller.go
│   ├── db/
│   │   └── db.go
│   ├── models/
│   │   └── user.go
│   ├── repository/
│   │   └── user_repository.go
│   └── service/
│       └── wallet_service.go
├── migrations/
│   └── 001_create_users_table.sql
├── .env
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## Setup and Running the Application

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/ewallet-app.git
   cd ewallet-app
   ```

2. Copy the `.env.example` file to `.env` and update the values as needed:
   ```
   cp .env.example .env
   ```

3. Build and run the application using Docker Compose:
   ```
   docker-compose up --build
   ```

   The application will be available at `http://localhost:8080`.

## API Endpoints

- Create User: `POST /users`
- Get User: `GET /users/{id}`
- Top Up: `POST /topup`
- Withdraw: `POST /withdraw`

## Testing the Application

You can use Postman to test the API endpoints. Here are some example requests:

1. Create a new user:
   ```
   POST http://localhost:8080/users
   Content-Type: application/json

   {
     "name": "John Doe"
   }
   ```

2. Top up a user's balance:
   ```
   POST http://localhost:8080/topup
   Content-Type: application/json

   {
     "user_id": 1,
     "amount": 100.00
   }
   ```

3. Withdraw from a user's balance:
   ```
   POST http://localhost:8080/withdraw
   Content-Type: application/json

   {
     "user_id": 1,
     "amount": 50.00
   }
   ```

4. Get user information:
   ```
   GET http://localhost:8080/users/1
   ```

---

<p align="center">Developed by</p>
<p align="center">Fikri Naufal Hamdi</p>
<p align="center">Information System and Technology @ ITB</p>