# QuickPic Test Client

A CLI tool to test messaging, encryption, and friend management with the QuickPic server.

## Prerequisites

- Server running at `http://localhost:8080`
- Go installed

## Commands

| Command | Description |
|---------|-------------|
| `register <username>` | Create a new account (no friend request) |
| `setup <username>` | Create account and send friend request to default target |
| `friend <username> <target>` | Send friend request to a specific user |
| `pending <username>` | List pending incoming friend requests |
| `accept <username> <request_id>` | Accept a friend request |
| `status <username>` | Check friend list and friendship status |
| `message --from <user> --to <user> <msg>` | Send an encrypted message (requires friendship) |
| `receive <username>` | Receive and decrypt incoming messages |
| `debug` | Test encryption/decryption locally |

## Usage Examples

### Create accounts and establish friendship

```bash
cd server

# Create two test accounts
go run ./cmd/testclient register maddie
# go run ./cmd/testclient register bob

go run ./cmd/testclient friend maddie reecepbcups

# reecepbcups checks pending requests (or accept in iOS app)
go run ./cmd/testclient pending maddie
# Output: ID: abc123... | From: maddie | At: 2024-01-08...

# reecepbcups accepts the request
go run ./cmd/testclient accept reecepbcups abc123-full-request-id

# Send a message from maddie to reecepbcups
go run ./cmd/testclient message --from maddie --to reecepbcups "Hello from test client!"

# get messages sent to maddie
go run ./cmd/testclient receive maddie
```

### Full test flow with iOS app

```bash
# 1. Create test user and send friend request to your iOS account
go run ./cmd/testclient setup testbot

# 2. Accept friend request in iOS app

# 3. Check friendship status
go run ./cmd/testclient status testbot

# 4. Send encrypted message (testbot -> reecepbcups)
go run ./cmd/testclient message --from testbot --to reecepbcups "Hello from test client!"

# 5. Check iOS app for the decrypted message

# 6. Receive messages sent from iOS
go run ./cmd/testclient receive testbot
```

### CLI-only testing (no iOS app)

```bash
# Create two users
go run ./cmd/testclient register user1
go run ./cmd/testclient register user2

# user1 sends friend request to user2
go run ./cmd/testclient friend user1 user2

# user2 accepts
go run ./cmd/testclient pending user2
go run ./cmd/testclient accept user2 <request_id>

# Now they can message each other
go run ./cmd/testclient message --from user1 --to user2 "Hello user2!"
go run ./cmd/testclient receive user2

# user2 replies
go run ./cmd/testclient message --from user2 --to user1 "Hey user1!"
go run ./cmd/testclient receive user1
```

## Credentials

Credentials are saved to `<username>.json` in the current directory after registration. This file contains:
- Username
- Private key (base64)
- Public key (base64)

The password for all test accounts is `testpassword123`.

## Configuration

Edit `main.go` to change:
- `baseURL` - Server address (default: `http://localhost:8080`)
- `targetUsername` - Default friend request target for `setup` command
