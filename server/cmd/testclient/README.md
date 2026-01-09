# QuickPic Test Client

A CLI tool to test messaging and encryption with the QuickPic server.

## Commands

## Full test flow

```bash
# 1. Setup test user
go run ./cmd/testclient setup testbot

# 2. Accept friend request in iOS app

# 3. Send message
go run ./cmd/testclient message testbot "Test message!"

# 4. Check iOS app for the decrypted message

go run ./cmd/testclient receive testbot
```
