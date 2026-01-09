# QuickPic E2E Tests

End-to-end tests for the QuickPic server API.

## Prerequisites

1. PostgreSQL running locally
2. Test database created

```bash
# Create test database
createdb quickpic_test
```

## Running Tests

### Using Make (recommended)

```bash
cd server

# Run all e2e tests
make test-e2e

# Run with verbose output
make test-e2e-verbose

# Run specific test
go test -v ./tests/e2e -run TestFullUserJourney
```

### Manual

```bash
cd server

# Set test database URL (optional, defaults to localhost)
export TEST_DATABASE_URL="postgres://localhost:5432/quickpic_test?sslmode=disable"

# Run tests
go test -v ./tests/e2e
```

## Test Coverage

### Auth Endpoints
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Token refresh
- `POST /auth/logout` - Logout (invalidate refresh token)

### Friend Endpoints
- `POST /friends/request` - Send friend request
- `GET /friends/requests` - Get pending requests
- `POST /friends/accept` - Accept request
- `POST /friends/reject` - Reject request
- `GET /friends` - List friends

### Message Endpoints
- `POST /messages` - Send encrypted message
- `GET /messages` - Get pending messages
- `POST /messages/:id/ack` - Acknowledge receipt

### Security Tests
- Protected routes require authentication
- Invalid tokens are rejected
- Users can only access their own data

### Full User Journey
- Complete flow: register → add friend → send message → acknowledge → logout

## Test Structure

```
tests/e2e/
├── setup_test.go   # Test configuration, helpers, types
├── api_test.go     # All API endpoint tests
└── README.md       # This file
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `TEST_DATABASE_URL` | `postgres://localhost:5432/quickpic_test?sslmode=disable` | PostgreSQL connection string |
