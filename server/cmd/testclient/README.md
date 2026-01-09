# QuickPic Test Client

A CLI tool to test messaging and encryption with the QuickPic server.

## Commands

### Setup (create account + send friend request)

```bash
cd server
go run ./cmd/testclient setup testbot
```

Creates user `testbot`, sends friend request to `reecepbcups`.

### Check friend status

```bash
go run ./cmd/testclient status testbot
```

### Send a message (after friend request is accepted)

```bash
go run ./cmd/testclient message testbot "Hello from the test client!"
```

## Full test flow

```bash
# 1. Setup test user
go run ./cmd/testclient setup testbot

# 2. Accept friend request in iOS app

# 3. Send message
go run ./cmd/testclient message testbot "Test message!"

# 4. Check iOS app for the decrypted message
```
