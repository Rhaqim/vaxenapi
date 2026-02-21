# Vaxen API (Golang)

This is the Golang implementation of the Vaxen API, rewritten from the original NestJS project. It provides a comprehensive API for global treasury and cross-border payment operations.

## 🚀 Features

- **Authentication & Authorization**: JWT-based authentication with role-based access control
- **Organizations**: Multi-organization support with KYB verification
- **Wallets & Accounts**: Multi-currency wallet and account management
- **Orders & Conversions**: FX trading and currency conversion
- **Payouts**: Cross-border payment processing
- **Beneficiaries**: Beneficiary management system
- **Crypto**: Cryptocurrency address generation and withdrawals
- **Reports**: FX P&L, transaction, and balance reporting
- **Webhooks**: Provider webhook integration
- **Admin**: Administrative controls and user management
- **Audit**: Comprehensive audit logging

## 🛠 Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin) - High-performance HTTP web framework
- **Database**: [PostgreSQL](https://www.postgresql.org/) with [GORM](https://gorm.io/)
- **Authentication**: [JWT](https://github.com/golang-jwt/jwt) - JSON Web Tokens
- **Documentation**: [Swaggo](https://github.com/swaggo/swag) - Automated Swagger documentation
- **Environment**: [godotenv](https://github.com/joho/godotenv) - Environment variable management
- **Caching**: Redis (optional)

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Redis (optional, for caching)
- Make (optional, for using Makefile commands)

## 🏁 Getting Started

### 1. Clone and Setup

```bash
cd api
cp .env.example .env
```

### 2. Configure Environment

Edit `.env` file with your configuration:

```env
PORT=8080
ENVIRONMENT=development
DATABASE_URL=postgresql://user:password@localhost:5432/vaxen?sslmode=disable
JWT_SECRET=your-secret-key
REDIS_HOST=localhost
REDIS_PORT=6379
```

### 3. Install Dependencies

```bash
make install
# or
go mod download
```

### 4. Run Database Migrations

```bash
make migrate
# or manually start the app (migrations run automatically)
```

### 5. Start the Server

```bash
make run
# or
go run main.go
```

The API will be available at:

- **API**: <http://localhost:8080>
- **Swagger Docs**: <http://localhost:8080/docs>
- **Health Check**: <http://localhost:8080/health>

## 🐳 Docker Setup

### Using Docker Compose (Recommended)

```bash
make docker-compose-up
```

This will start:

- API server on port 8080
- PostgreSQL on port 5432
- Redis on port 6379

### Using Docker

```bash
make docker-build
make docker-run
```

## 📚 API Endpoints

### Authentication

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Refresh token

### Organizations

- `GET /api/v1/organizations` - List organizations
- `GET /api/v1/organizations/:id` - Get organization
- `POST /api/v1/organizations` - Create organization
- `PUT /api/v1/organizations/:id` - Update organization

### KYB

- `POST /api/v1/kyb/submit` - Submit KYB information
- `GET /api/v1/kyb/status` - Get KYB status

### Wallets

- `GET /api/v1/wallets` - List wallets
- `GET /api/v1/wallets/:id` - Get wallet
- `POST /api/v1/wallets` - Create wallet
- `GET /api/v1/wallets/:id/balance` - Get wallet balance

### Accounts

- `GET /api/v1/accounts` - List accounts
- `GET /api/v1/accounts/:id` - Get account
- `POST /api/v1/accounts` - Create account

### Orders

- `GET /api/v1/orders` - List orders
- `GET /api/v1/orders/open` - List open orders
- `GET /api/v1/orders/:id` - Get order
- `POST /api/v1/orders` - Create order
- `PUT /api/v1/orders/:id/cancel` - Cancel order

### Quotes & Conversions

- `POST /api/v1/quotes` - Create quote
- `GET /api/v1/quotes/:id` - Get quote
- `GET /api/v1/conversions` - List conversions
- `POST /api/v1/conversions` - Create conversion
- `GET /api/v1/conversions/:id` - Get conversion

### Beneficiaries

- `GET /api/v1/beneficiaries` - List beneficiaries
- `POST /api/v1/beneficiaries` - Create beneficiary
- `GET /api/v1/beneficiaries/:id` - Get beneficiary
- `PUT /api/v1/beneficiaries/:id` - Update beneficiary
- `DELETE /api/v1/beneficiaries/:id` - Delete beneficiary

### Payouts

- `GET /api/v1/payouts` - List payouts
- `POST /api/v1/payouts` - Create payout
- `GET /api/v1/payouts/:id` - Get payout

### Crypto

- `GET /api/v1/crypto/addresses` - List crypto addresses
- `POST /api/v1/crypto/addresses` - Create crypto address
- `POST /api/v1/crypto/withdraw` - Withdraw crypto

### Reports

- `GET /api/v1/reports/fx-pnl` - FX P&L report
- `GET /api/v1/reports/transactions` - Transaction report
- `GET /api/v1/reports/balances` - Balance report

### Statements

- `GET /api/v1/statements` - List statements
- `GET /api/v1/statements/:id` - Get statement

### Webhooks

- `POST /api/v1/webhooks/provider/:provider` - Process provider webhook

### Admin (Requires Admin Role)

- `GET /api/v1/admin/users` - List all users
- `GET /api/v1/admin/users/:id` - Get user
- `PUT /api/v1/admin/users/:id` - Update user
- `DELETE /api/v1/admin/users/:id` - Delete user
- `GET /api/v1/admin/organizations` - List all organizations
- `POST /api/v1/admin/organizations/:id/approve` - Approve organization
- `POST /api/v1/admin/organizations/:id/reject` - Reject organization

### Audit

- `GET /api/v1/audit/logs` - Get audit logs

## 🔐 Authentication

All protected endpoints require a JWT token in the Authorization header:

```shell
Authorization: Bearer <your-jwt-token>
```

## 🧪 Testing

```bash
make test
# or with coverage
make test-coverage
```

## 📖 Generate Swagger Documentation

```bash
make swagger
```

## 🔨 Available Make Commands

```bash
make help              # Display available commands
make install           # Install dependencies
make build             # Build the application
make run               # Run the application
make dev               # Run with auto-reload (requires air)
make test              # Run tests
make test-coverage     # Run tests with coverage
make clean             # Clean build artifacts
make swagger           # Generate Swagger documentation
make migrate           # Run database migrations
make docker-build      # Build Docker image
make docker-run        # Run Docker container
make docker-compose-up # Start services with docker-compose
make lint              # Run linter
make fmt               # Format code
make vet               # Run go vet
```

## 📁 Project Structure

```shell
api/
├── main.go                    # Application entry point
├── go.mod                     # Go module dependencies
├── Makefile                   # Build and development commands
├── Dockerfile                 # Docker configuration
├── docker-compose.yml         # Docker Compose configuration
├── .env.example               # Environment variables template
├── internal/
│   ├── config/                # Configuration management
│   │   └── config.go
│   ├── middleware/            # HTTP middleware
│   │   ├── auth.go
│   │   └── middleware.go
│   ├── routes/                # Route definitions
│   │   └── routes.go
│   ├── handlers/              # HTTP handlers
│   │   ├── auth.go
│   │   ├── organizations.go
│   │   ├── wallets.go
│   │   ├── orders.go
│   │   ├── payouts.go
│   │   ├── conversions.go
│   │   ├── reports.go
│   │   ├── kyb.go
│   │   ├── crypto.go
│   │   ├── webhooks.go
│   │   ├── admin.go
│   │   └── health.go
│   ├── models/                # Database models
│   │   └── models.go
│   ├── database/              # Database connection
│   │   └── database.go
│   └── utils/                 # Utility functions
│       ├── jwt.go
│       └── response.go
└── README.md
```

## 🚧 Development

### Hot Reload

Install [Air](https://github.com/cosmtrek/air) for hot reload during development:

```bash
go install github.com/cosmtrek/air@latest
make dev
```

### Code Formatting

```bash
make fmt
make vet
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License.

## 🔗 Related Projects

- [Backend (Original Go Implementation)](../backend)
- [Frontend](../frontend)

## 📞 Support

For support, email [support@vaxen.io](mailto:support@vaxen.io) or open an issue in the repository.
