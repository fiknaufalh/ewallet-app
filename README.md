# E-Wallet Service

## Description
A robust e-wallet service built with Go that handles user wallet operations with a focus on reliability, consistency, and safety. This service allows users to create an account with a wallet, perform top-ups, withdrawals, and check their balance.

## Features
### Core Features
- User account creation and management
- Wallet balance management
- Top-up functionality
- Withdrawal to bank account
- Balance inquiry

### Technical Features
1. **ACID Compliance**
   - Transaction atomicity for all operations
   - Consistent balance updates
   - Isolated concurrent operations
   - Durable transaction records

2. **Security & Safety**
   - Idempotency handling to prevent duplicate transactions
   - Transaction limits enforcement
   - Optimistic locking for concurrent updates
   - Input validation and sanitization

3. **Race Condition Handling**
   - Database-level locks
   - Transaction isolation level: Serializable
   - Version control for wallet updates

## Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL 14
- Postman (for testing)

## Setup and Installation

### 1. Clone the Repository
```bash
git clone <repository-url>
cd ewallet-app
```

### 2. Environment Configuration
Create `.env` file in the config directory:

```bash
cp config/.env.example config/.env
```

Example `.env` configuration:
```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ewallet
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_CONN_MAX_LIFETIME=300s

# Transaction Configuration
MAX_WITHDRAWAL_AMOUNT=10000000
MIN_WITHDRAWAL_AMOUNT=10000
MAX_TOPUP_AMOUNT=50000000
MIN_TOPUP_AMOUNT=10000

# Security
IDEMPOTENCY_KEY_EXPIRATION=24h
```

### 3. Running the Application
Using Docker Compose (recommended):
```bash
# Build and start all services
docker-compose up --build
```

The application will be available at `http://localhost:3000`.

## API Documentation and Testing Guide

### 1. Create User
Creates a new user with an associated wallet.

**Endpoint:** `POST /api/v1/users`

**Request:**
```json
{
    "username": "john_doe",
    "email": "john@example.com"
}
```

**Response (201 Created):**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "john_doe",
    "email": "john@example.com"
}
```

### 2. Top Up Wallet
Adds money to user's wallet.

**Endpoint:** `POST /api/v1/topup`

**Headers:**
- Content-Type: application/json
- X-Idempotency-Key: unique-key-123 (must be unique for each request)

**Request:**
```json
{
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "amount": 100000,
    "reference_id": "TOP-123456"
}
```

**Response (200 OK):**
```json
{
    "transaction_id": "123e4567-e89b-12d3-a456-426614174001",
    "status": "completed",
    "balance": 100000
}
```

### 3. Withdraw from Wallet
Withdraws money from user's wallet to a bank account.

**Endpoint:** `POST /api/v1/withdraw`

**Headers:**
- Content-Type: application/json
- X-Idempotency-Key: unique-key-456 (must be unique for each request)

**Request:**
```json
{
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "amount": 50000,
    "reference_id": "WD-123456",
    "bank_account": "1234567890"
}
```

**Response (200 OK):**
```json
{
    "transaction_id": "123e4567-e89b-12d3-a456-426614174002",
    "status": "completed",
    "balance": 50000
}
```

### 4. Get Balance
Retrieves current wallet balance.

**Endpoint:** `GET /api/v1/balance/{user_id}`

**Response (200 OK):**
```json
{
    "balance": 50000
}
```

## Testing Scenarios

### Complete Flow Test
1. Create a new user
```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"john_doe","email":"john@example.com"}'
```

2. Top up the wallet (save the user_id from step 1)
```bash
curl -X POST http://localhost:3000/api/v1/topup \
  -H "Content-Type: application/json" \
  -H "X-Idempotency-Key: unique-key-123" \
  -d '{"user_id":"YOUR-USER-ID","amount":100000,"reference_id":"TOP-123456"}'
```

3. Check balance
```bash
curl http://localhost:3000/api/v1/balance/YOUR-USER-ID
```

4. Withdraw some money
```bash
curl -X POST http://localhost:3000/api/v1/withdraw \
  -H "Content-Type: application/json" \
  -H "X-Idempotency-Key: unique-key-456" \
  -d '{"user_id":"YOUR-USER-ID","amount":50000,"reference_id":"WD-123456","bank_account":"1234567890"}'
```

5. Check balance again
```bash
curl http://localhost:3000/api/v1/balance/YOUR-USER-ID
```

### Testing Idempotency
Try sending the same top-up request twice with the same idempotency key:
```bash
# First request - should succeed
curl -X POST http://localhost:3000/api/v1/topup \
  -H "Content-Type: application/json" \
  -H "X-Idempotency-Key: test-key-1" \
  -d '{"user_id":"YOUR-USER-ID","amount":100000,"reference_id":"TOP-789"}'

# Second request with same key - should return the same response without processing again
curl -X POST http://localhost:3000/api/v1/topup \
  -H "Content-Type: application/json" \
  -H "X-Idempotency-Key: test-key-1" \
  -d '{"user_id":"YOUR-USER-ID","amount":100000,"reference_id":"TOP-789"}'
```

### Testing Validation
1. Try withdrawing more than available balance
2. Try top-up with amount outside limits
3. Try withdrawal with amount outside limits
4. Try requests without idempotency key
5. Try invalid user IDs

## Error Handling
The service handles various error cases:
- Invalid input validation
- Insufficient balance
- Duplicate transactions
- Concurrent updates
- Database errors

## Common Issues and Solutions

### Cannot connect to database
1. Check if PostgreSQL container is running:
```bash
docker ps
```
2. Check logs:
```bash
docker-compose logs db
```

### API returns 500 error
1. Check application logs:
```bash
docker-compose logs app
```

### Database migrations not applied
1. The migrations are automatically applied when the container starts
2. To manually apply migrations:
```bash
docker-compose exec app ./main migrate
```

## Development Notes
- All monetary values are handled in the smallest currency unit (e.g., cents)
- Idempotency keys expire after 24 hours
- Transaction isolation level is set to Serializable for maximum consistency

## Monitoring and Logging
- Application logs are available through Docker:
```bash
docker-compose logs -f app
```

---

<p align="center">Developed by</p>
<p align="center">Fikri Naufal Hamdi</p>
<p align="center">Information System and Technology @ ITB</p>