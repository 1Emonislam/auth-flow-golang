OwnPayNet Service
A cost-free, non-custodial Bitcoin payment gateway built with Golang, Gin, PostgreSQL, and Bitcoin Core.
Features

User signup, signin, and password reset.
JWT-based authentication for protected routes.
Create Bitcoin payment requests with unique addresses.
Transaction monitoring via Bitcoin Core.
Non-custodial: funds go directly to merchant wallets.

Prerequisites

Go 1.21+
PostgreSQL
Bitcoin Core (testnet for development, mainnet for production)
ngrok (for local webhook testing)

Setup

Clone the repository:git clone <repository-url>
cd own-paynet


Install Bitcoin Core and configure bitcoin.conf:testnet=1
rpcuser=your_rpc_user
rpcpassword=your_rpc_password
rpcallowip=127.0.0.1
rpcbind=127.0.0.1
server=1

Run Bitcoin Core:bitcoind -testnet -rpcuser=your_rpc_user -rpcpassword=your_rpc_password


Create a .env file based on the example.
Install dependencies:go mod tidy


Run the application:go run main.go



API Endpoints

POST /api/v1/signup: Register a new user.
POST /api/v1/signin: Login and get JWT token.
POST /api/v1/reset-password: Reset user password.
POST /api/v1/payments: Create a payment request (protected).
POST /api/v1/webhook: Receive transaction updates.

Testing

Start PostgreSQL and Bitcoin Core.
Use Postman to test endpoints:
Signup: POST http://localhost:8080/api/v1/signup{"email": "user@example.com", "password": "password123"}


Signin: POST http://localhost:8080/api/v1/signin{"email": "user@example.com", "password": "password123"}


Create Payment: POST http://localhost:8080/api/v1/payments (with Authorization header){"amount": 0.001, "merchant_wallet": "tb1q...", "currency": "BTC"}


Webhook: POST http://localhost:8080/api/v1/webhook (with X-Webhook-Signature){"payment_id": "generated_payment_id", "status": "confirmed", "address": "btc_address"}





Deployment

Dockerize the application:FROM golang:1.21
WORKDIR /app
COPY . .
RUN go build -o main
CMD ["./main"]


Deploy to a cloud provider with Bitcoin Core running.

