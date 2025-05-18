# Paper.id Wallet Disbursement API

## Overview
This is a Go implementation of a wallet disbursement API for Paper.id. The API allows users to transfer funds from their Paper.id wallet to their bank account.

## Features
- Single API endpoint for wallet disbursement (`POST /api/disbursements`)
- Database-backed user accounts, wallets, and bank accounts
- Transaction-based operations to ensure data integrity
- Input validation
- Proper error handling and responses
- Test data seeding for easy demonstration

## Technical Implementation

### Data Models
- **User**: Basic user information
- **Wallet**: Represents a user's digital wallet with balance
- **BankAccount**: A user's registered bank account for disbursements
- **Disbursement**: Records of disbursement transactions

### API Endpoints

#### `GET /`
- Health check endpoint that confirms the API is running

#### `POST /api/disbursements`
- **Purpose**: Disburse funds from a user's wallet to their bank account
- **Request Body**:
  ```json
  {
    "user_id": 1,
    "bank_account_id": 1,
    "amount": 500000,
    "description": "Monthly withdrawal"
  }
  ```
- **Successful Response**:
  ```json
  {
    "success": true,
    "message": "Disbursement processed successfully",
    "disbursement": {
      "id": 1,
      "user_id": 1,
      "wallet_id": 1,
      "bank_account_id": 1,
      "amount": 500000,
      "status": "completed",
      "reference": "DSB16301234567890",
      "created_at": "2025-05-16T12:34:56Z",
      "updated_at": "2025-05-16T12:34:56Z",
      "completed_at": "2025-05-16T12:34:56Z"
    }
  }
  ```
- **Error Response**:
  ```json
  {
    "success": false,
    "message": "An error occurred",
    "error": "insufficient wallet balance"
  }
  ```

## Validation
The API performs the following validations:
1. JSON payload structure
2. User existence
3. Bank account existence and ownership
4. Wallet existence
5. Sufficient balance

## Database Schema
The application uses SQLite for simplicity but can be adapted to work with any SQL database.

Tables:
- users
- wallets
- bank_accounts
- disbursements

## Security Considerations
- Input validation to prevent injection attacks
- Proper error handling to avoid information leakage
- Transaction-based operations to maintain data integrity

## Deployment
To run the application:

```bash
go mod init github.com/yourusername/paper-id-disbursement
go mod tidy
go run main.go
```

The API will be available at http://localhost:8080

## Testing
You can test the API using curl or any API testing tool like Postman:

```bash
curl -X POST http://localhost:8080/api/disbursements \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "bank_account_id": 1, "amount": 500000, "description": "Test withdrawal"}'
```

## Future Improvements
In a production environment, consider implementing:
1. Authentication and authorization
2. Rate limiting
3. More sophisticated error handling
4. Event logging and monitoring
5. Integration with actual payment gateways for real disbursements
6. Idempotency keys to prevent duplicate disbursements

## Project Structure
```
paper-id-disbursement/
├── main.go          # Main application file
├── api.go           # API handler functions
├── constants.go     # Error constants (will contain other constants as the app grows)
├── init.go          # Database init functions
├── models.go        # Data structures and models
├── go.mod           # Go module definition
├── go.sum           # Go module checksums
├── paper_id.db      # SQLite database file
└── README.md        # Project documentation
```
