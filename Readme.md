# Transactions routine

This project is a simple web service that handles accounts and transactions.

### Features

- Create accounts
- Retrieve account details
- Process transactions

## Installation

To install and set up the project, follow these steps:

1. Clone the repository:
```
   bash
   git clone https://github.com/mohammadshabab/com.pismo.transaction.routine.git
   cd com.pismo.transaction.routine
   go mod download
   go run main.go
```

## Request response

The web service has 3 end points 
### 1. Create Account

**Endpoint:**
`POST /accounts`

**REQUEST**
```
{
    "document_number": "12345678900"
}
```
**RESPONSE**
```
{
    "account_id": 1,
    "document_number": "12345678900"
}
```

### 2. Get Account Details

**Endpoint:**
`GET /accounts/{id}` 

**RESPONSE**
```
{
    "account_id": 1,
    "document_number": "12345678900"
}
```

### 3. Create Transaction

**Endpoint:**
`POST /transactions` 

**REQUEST** 
```
{
  "account_id": 1,
  "operation_type_id":4,
  "amount": 250.45
}
```
**RESPONSE**
```
{
    "transaction_id": 10,
    "account_id": 1,
    "operation_type_id": 4,
    "amount": 250.45,
    "event_date": "2024-09-28T15:45:33.1492902+05:30"
}
```
