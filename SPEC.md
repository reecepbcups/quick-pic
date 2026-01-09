# QuickPic Specification

## Overview

QuickPic is a privacy-focused ephemeral messaging app (Snapchat clone) with end-to-end encryption. Messages are encrypted client-side before transmission, ensuring the server never has access to plaintext content. Data persistence will eventually be handled by a Cosmos SDK blockchain.

> **Technical Decisions**: See [DECISIONS.md](./DECISIONS.md) for all finalized technical choices.

## Architecture

```
┌─────────────┐     E2E Encrypted      ┌─────────────┐
│   iOS App   │ ◄──────────────────────► │  Go Server  │
│  (SwiftUI)  │                          │ (Gin + PG)  │
└─────────────┘                          └──────┬──────┘
                                                │
                                                ▼ (Phase 2)
                                         ┌─────────────┐
                                         │ Cosmos SDK  │
                                         │ Blockchain  │
                                         └─────────────┘
```

## Core Principles

1. **Zero-knowledge server** - Server only sees encrypted blobs
2. **1:1 messaging only** - No group chats (simplifies encryption)
3. **Ephemeral by default** - View-once messages, 24hr local cache
4. **Client-side encryption** - All crypto happens on device
5. **Single device per user** - No multi-device sync complexity

---

## Phase 1: iOS App + Go Server

### iOS App Features

#### Authentication
- [ ] User registration with username/password
- [ ] Login/logout with JWT tokens
- [ ] X25519 key pair generation on first launch
- [ ] Private key storage in iOS Keychain (Secure Enclave)

#### Contacts/Friends
- [ ] Add friends by username
- [ ] Accept/reject friend requests
- [ ] View friends list
- [ ] Fetch friend's public key on connection

#### Camera & Media
- [ ] Take photo with front/back camera
- [ ] Photo preview before sending
- [ ] PNG conversion (lossless)
- [ ] Gzip compression before encryption
- [ ] 5MB max file size limit

#### Messaging
- [ ] Send encrypted photo to a friend
- [ ] Send encrypted text message
- [ ] View-once display (auto-hide after viewing)
- [ ] 24-hour local message cache
- [ ] Unread message indicators
- [ ] Generic push notifications ("You have a new message")

#### Encryption
- [ ] Generate X25519 key pair
- [ ] XChaCha20-Poly1305 symmetric encryption
- [ ] Encrypt messages with derived shared secret
- [ ] Sign messages with sender's private key
- [ ] Decrypt incoming messages

#### Local Storage
- [ ] Secure storage for decrypted message cache
- [ ] Auto-purge messages older than 24 hours
- [ ] Clear all data on logout

### Go Server

#### Technology Stack
- Go 1.21+
- Gin HTTP framework
- PostgreSQL database
- Argon2id password hashing
- JWT authentication

#### API Endpoints
```
POST   /auth/register     - Create account (username, password, public_key)
POST   /auth/login        - Authenticate, return JWT
POST   /auth/refresh      - Refresh access token
POST   /auth/logout       - Invalidate refresh token

GET    /users/:username   - Get user public info + public key
POST   /friends/request   - Send friend request
GET    /friends/requests  - List pending incoming requests
POST   /friends/accept    - Accept friend request
POST   /friends/reject    - Reject friend request
GET    /friends           - List friends with public keys

POST   /messages          - Send encrypted message blob
GET    /messages          - Fetch pending messages for user
POST   /messages/:id/ack  - Acknowledge receipt (triggers server delete)
```

#### Data Models
```
User {
  id: UUID
  username: String (unique, lowercase)
  password_hash: String (Argon2id)
  public_key: String (base64 X25519)
  created_at: Timestamp
  updated_at: Timestamp
}

FriendRequest {
  id: UUID
  from_user_id: UUID
  to_user_id: UUID
  status: pending | accepted | rejected
  created_at: Timestamp
}

Friendship {
  id: UUID
  user_a_id: UUID
  user_b_id: UUID
  created_at: Timestamp
}

Message {
  id: UUID
  from_user_id: UUID
  to_user_id: UUID
  encrypted_content: Blob
  content_type: text | image
  signature: String (base64)
  created_at: Timestamp
}
```

#### Server Security
- Passwords: Argon2id hashing, never logged
- Transport: HTTPS/TLS 1.3 only
- Auth: Short-lived JWT (15min) + refresh tokens (7 days)
- Purge: Delete messages immediately after ACK
- Fallback purge: Cron job deletes undelivered messages after 7 days

---

## Phase 2: Cosmos SDK Blockchain

> Deferred - Server storage is sufficient while encryption is solid

The blockchain will:
- Store message metadata (not content - that stays encrypted)
- Provide immutable audit trail
- Enable decentralized identity
- Handle key registration/rotation

---

## Security Model

### Encryption Flow

```
Sender                                    Recipient
──────                                    ─────────
1. Compose message (text or photo)
2. If image: PNG → gzip compress
3. Generate ephemeral symmetric key
4. Encrypt content with XChaCha20-Poly1305
5. Encrypt symmetric key with recipient's X25519 public key
6. Sign encrypted blob with sender's private key
7. Send to server ─────────────────────► Server stores blob
                                         ◄─────────────────── 8. Fetch encrypted message
                                         9. Verify signature with sender's public key
                                         10. Decrypt symmetric key with private key
                                         11. Decrypt content with symmetric key
                                         12. If image: decompress → display
                                         13. ACK to server (triggers delete)
                                         14. Cache locally for 24hrs
```

### Key Management
- Private keys: iOS Keychain with Secure Enclave
- Public keys: Stored on server (they're public)
- No key rotation for MVP (future enhancement)

### Threat Model
| Threat | Mitigation |
|--------|------------|
| Server compromise | E2E encryption - server never sees plaintext |
| Man-in-the-middle | TLS 1.3 + signature verification |
| Device theft | Keychain + biometric protection |
| Credential stuffing | Rate limiting + Argon2id slow hash |
| Replay attacks | Message signatures + nonces |

---

## MVP Scope

### In Scope
1. Username/password registration and login
2. Add friends by username, accept/reject requests
3. Take photo with camera, preview, send
4. Simple text message composition
5. X25519 + XChaCha20-Poly1305 encryption
6. View-once message display
7. 24-hour local cache
8. Generic push notifications

### Out of Scope
- Stories/public posts
- Video messages
- Filters/lenses/AR
- Group chats
- Read receipts
- Typing indicators
- Bitmoji/avatars
- Multi-device sync
- Key rotation
- Account recovery

---

## File Structure

```
quick-pic/
├── SPEC.md
├── DECISIONS.md
├── ios/
│   └── QuickPic/
│       ├── QuickPic.xcodeproj
│       ├── App/
│       │   └── QuickPicApp.swift
│       ├── Views/
│       │   ├── Auth/
│       │   │   ├── LoginView.swift
│       │   │   └── RegisterView.swift
│       │   ├── Camera/
│       │   │   └── CameraView.swift
│       │   ├── Messages/
│       │   │   ├── InboxView.swift
│       │   │   └── MessageView.swift
│       │   └── Friends/
│       │       ├── FriendsListView.swift
│       │       └── AddFriendView.swift
│       ├── Models/
│       │   ├── User.swift
│       │   ├── Message.swift
│       │   └── Friend.swift
│       ├── Services/
│       │   ├── CryptoService.swift
│       │   ├── APIService.swift
│       │   ├── KeychainService.swift
│       │   └── MessageCacheService.swift
│       └── Resources/
│           └── Assets.xcassets
├── server/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── routes.go
│   │   ├── models/
│   │   ├── repository/
│   │   └── services/
│   ├── pkg/
│   │   └── crypto/
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
└── blockchain/ (Phase 2)
    └── cosmos-app/
```

---

## Implementation Order

1. **Server foundation**: Go project setup, DB schema, auth endpoints
2. **iOS foundation**: SwiftUI project, navigation, auth views
3. **Encryption layer**: Key generation, encrypt/decrypt on iOS
4. **Auth flow**: Register/login connected to server
5. **Friends system**: Add, accept, list friends
6. **Camera**: Capture, preview, compress images
7. **Messaging**: Send/receive encrypted messages
8. **View-once UI**: Display and auto-hide messages
9. **Local cache**: 24-hour message storage
10. **Push notifications**: APNs integration
11. **Polish**: Error handling, loading states, edge cases

---

## Related Documents

- [DECISIONS.md](./DECISIONS.md) - All technical decisions with rationale
