# QuickPic Technical Decisions

This document records all finalized technical decisions for the QuickPic project.

---

## 1. Encryption Scheme

**Decision: X25519 + XChaCha20-Poly1305**

- Industry standard used by Signal, WhatsApp, and other secure messengers
- X25519 for key exchange (Curve25519 Diffie-Hellman)
- XChaCha20-Poly1305 for symmetric encryption of message content
- 32-byte keys, fast performance, strong security guarantees
- iOS: Use `libsodium` via CryptoKit or Swift wrapper

---

## 2. Authentication Method

**Decision: Username + Password**

Security requirements:
- Passwords hashed with Argon2id (preferred) or bcrypt
- **Never store plaintext passwords**
- **Never log passwords or password-adjacent data**
- Salt stored per-user
- Minimum password requirements TBD (suggest: 8+ chars)

Implementation notes:
- Server receives password over TLS
- Immediately hash on receipt
- Clear password from memory after hashing
- Use constant-time comparison for verification

---

## 3. Message Expiration

**Decision: View-once with 24-hour local cache**

Behavior:
- Messages auto-delete from UI after viewing
- Device caches last 24 hours of messages locally (for replay/review)
- Server purges messages aggressively after delivery confirmation
- Single-device assumption per user (no multi-device sync)

Server purge strategy:
- Delete message from DB once recipient confirms receipt
- Fallback: Delete undelivered messages after 7 days
- No server-side message history

Local cache:
- Store decrypted messages in secure local storage
- Auto-purge messages older than 24 hours
- Clear cache on logout

---

## 4. Server Technology Stack

**Decision: Go**

Stack:
- **Language**: Go 1.21+
- **Framework**: Gin or Chi (lightweight HTTP router)
- **Database**: PostgreSQL
- **ORM**: sqlc or GORM
- **Auth**: JWT with short-lived access tokens + refresh tokens

Rationale:
- Go aligns with Cosmos SDK (Phase 2)
- Excellent performance
- Strong standard library
- Easy deployment

---

## 5. Key Exchange Protocol

**Decision: Simple Public Key Exchange**

Flow:
1. User generates X25519 key pair on first app launch
2. Public key uploaded to server during registration
3. When adding a friend, client fetches friend's public key from server
4. Messages encrypted using recipient's public key
5. No forward secrecy (acceptable for MVP, single-device model)

Key storage:
- Private key: iOS Keychain (Secure Enclave if available)
- Public key: Server database (plaintext, it's public)

Key rotation:
- Not implemented for MVP
- Future: Allow key regeneration with friend notification

---

## 6. Image Handling

**Decision: 5MB max with lossless compression**

Pipeline:
1. Capture image from camera
2. Convert to PNG (lossless)
3. Compress using zlib/gzip before encryption
4. Encrypt compressed data
5. Transmit encrypted blob
6. Recipient decrypts → decompresses → displays

Specifications:
- **Max file size**: 5MB (post-compression, pre-encryption)
- **Format**: PNG (lossless quality)
- **Compression**: gzip/zlib
- **Resolution**: Preserve original (no downscaling)

Rationale:
- Lossless ensures image quality
- Compression reduces bandwidth without quality loss
- 5MB limit prevents abuse while allowing high-res photos

---

## 7. Push Notifications

**Decision: Generic notifications only**

Format:
```
Title: "QuickPic"
Body: "You have a new message"
```

Privacy measures:
- No sender name in notification
- No message preview
- No image thumbnail
- Badge count only shows unread count (no details)

Implementation:
- Apple Push Notification Service (APNs)
- Server sends push when message stored
- Notification contains no sensitive data
- App fetches actual message on open

---

## Summary Table

| Decision | Choice |
|----------|--------|
| Encryption | X25519 + XChaCha20-Poly1305 |
| Auth | Username/password (Argon2id hash) |
| Message expiry | View-once + 24hr local cache |
| Server | Go + Gin + PostgreSQL |
| Key exchange | Simple public key |
| Images | 5MB max, PNG, gzip compressed |
| Notifications | Generic "new message" only |

---

## Additional Constraints

- **Single device per user**: No multi-device sync
- **1:1 messages only**: No group chats
- **Zero-knowledge server**: Server never sees plaintext
- **No message history on server**: Purge after delivery
