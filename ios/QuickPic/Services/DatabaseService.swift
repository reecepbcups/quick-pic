//
//  DatabaseService.swift
//  QuickPic
//
//  SQLite-based local storage for messages and conversations
//

import Foundation
import SQLite3

final class DatabaseService: @unchecked Sendable {
    static let shared = DatabaseService()

    private var db: OpaquePointer?
    private let dbQueue = DispatchQueue(label: "com.quickpic.database", qos: .userInitiated)

    private init() {
        openDatabase()
        createTables()
        startCleanupTimer()
    }

    deinit {
        sqlite3_close(db)
    }

    // MARK: - Database Setup

    private func openDatabase() {
        let fileManager = FileManager.default
        let paths = fileManager.urls(for: .documentDirectory, in: .userDomainMask)
        let dbPath = paths[0].appendingPathComponent("quickpic.db").path

        if sqlite3_open(dbPath, &db) != SQLITE_OK {
            print("Failed to open database at \(dbPath)")
        }
    }

    private func createTables() {
        // Conversations table - stores friend conversations
        let createConversationsTable = """
            CREATE TABLE IF NOT EXISTS conversations (
                friend_user_id TEXT PRIMARY KEY,
                friend_username TEXT NOT NULL,
                friend_public_key TEXT NOT NULL,
                friend_since TEXT NOT NULL DEFAULT '',
                last_message_at TEXT,
                unread_count INTEGER DEFAULT 0,
                created_at TEXT NOT NULL
            );
        """

        // Messages table - stores all messages (sent and received)
        let createMessagesTable = """
            CREATE TABLE IF NOT EXISTS messages (
                id TEXT PRIMARY KEY,
                conversation_id TEXT NOT NULL,
                content_type TEXT NOT NULL,
                decrypted_content BLOB NOT NULL,
                encrypted_content BLOB,
                is_from_me INTEGER NOT NULL,
                has_been_viewed INTEGER DEFAULT 0,
                created_at TEXT NOT NULL,
                received_at TEXT NOT NULL,
                FOREIGN KEY (conversation_id) REFERENCES conversations(friend_user_id)
            );
        """

        // Create index for faster conversation lookups
        let createIndex = """
            CREATE INDEX IF NOT EXISTS idx_messages_conversation
            ON messages(conversation_id, created_at DESC);
        """

        executeSQL(createConversationsTable)
        executeSQL(createMessagesTable)
        executeSQL(createIndex)

        // Add new columns for existing databases (migrations) - silently ignore if already exists
        executeMigration("ALTER TABLE conversations ADD COLUMN friend_since TEXT NOT NULL DEFAULT '';")
        executeMigration("ALTER TABLE messages ADD COLUMN encrypted_content BLOB;")
    }

    private func executeSQL(_ sql: String) {
        var statement: OpaquePointer?
        if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
            if sqlite3_step(statement) != SQLITE_DONE {
                print("SQL execution failed: \(sql)")
            }
        } else {
            print("SQL preparation failed: \(sql)")
        }
        sqlite3_finalize(statement)
    }

    private func executeMigration(_ sql: String) {
        // Silently execute migration - ignore errors (e.g., duplicate column)
        var statement: OpaquePointer?
        if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
            sqlite3_step(statement)
        }
        sqlite3_finalize(statement)
    }

    // MARK: - Conversation Operations

    func getOrCreateConversation(for friend: Friend) -> Conversation {
        dbQueue.sync {
            // Store strings to ensure they stay alive during SQLite operations
            let friendIDString = friend.userID.uuidString
            let friendUsername = friend.username
            let friendPublicKey = friend.publicKey
            let friendSinceString = ISO8601DateFormatter().string(from: friend.since)
            let nowString = ISO8601DateFormatter().string(from: Date())

            // Check if conversation exists
            let selectSQL = "SELECT friend_user_id, friend_username, friend_public_key, friend_since, last_message_at, unread_count, created_at FROM conversations WHERE friend_user_id = ?;"
            var statement: OpaquePointer?

            if sqlite3_prepare_v2(db, selectSQL, -1, &statement, nil) == SQLITE_OK {
                friendIDString.withCString { cString in
                    sqlite3_bind_text(statement, 1, cString, -1, nil)
                }

                if sqlite3_step(statement) == SQLITE_ROW {
                    let conversation = conversationFromStatement(statement)
                    sqlite3_finalize(statement)
                    return conversation
                }
            }
            sqlite3_finalize(statement)

            // Create new conversation
            let insertSQL = """
                INSERT OR REPLACE INTO conversations (friend_user_id, friend_username, friend_public_key, friend_since, created_at, unread_count)
                VALUES (?, ?, ?, ?, ?, 0);
            """

            if sqlite3_prepare_v2(db, insertSQL, -1, &statement, nil) == SQLITE_OK {
                friendIDString.withCString { id in
                    friendUsername.withCString { name in
                        friendPublicKey.withCString { key in
                            friendSinceString.withCString { since in
                                nowString.withCString { now in
                                    sqlite3_bind_text(statement, 1, id, -1, nil)
                                    sqlite3_bind_text(statement, 2, name, -1, nil)
                                    sqlite3_bind_text(statement, 3, key, -1, nil)
                                    sqlite3_bind_text(statement, 4, since, -1, nil)
                                    sqlite3_bind_text(statement, 5, now, -1, nil)
                                    let result = sqlite3_step(statement)
                                    if result != SQLITE_DONE {
                                        print("Failed to insert conversation: \(result)")
                                    }
                                }
                            }
                        }
                    }
                }
            } else {
                print("Failed to prepare insert statement")
            }
            sqlite3_finalize(statement)

            return Conversation(
                friendUserID: friend.userID,
                friendUsername: friend.username,
                friendPublicKey: friend.publicKey,
                friendSince: friend.since,
                lastMessageAt: nil,
                unreadCount: 0,
                createdAt: Date()
            )
        }
    }

    func getAllConversations() -> [Conversation] {
        dbQueue.sync {
            var conversations: [Conversation] = []
            // Order by last_message_at with NULLs at end, then by created_at
            let sql = """
                SELECT friend_user_id, friend_username, friend_public_key, friend_since, last_message_at, unread_count, created_at
                FROM conversations
                ORDER BY
                    CASE WHEN last_message_at IS NULL THEN 1 ELSE 0 END,
                    last_message_at DESC,
                    created_at DESC;
            """
            var statement: OpaquePointer?

            if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
                while sqlite3_step(statement) == SQLITE_ROW {
                    conversations.append(conversationFromStatement(statement))
                }
            } else {
                print("Failed to prepare getAllConversations query")
            }
            sqlite3_finalize(statement)
            print("getAllConversations returning \(conversations.count) conversations")
            return conversations
        }
    }

    func updateConversationLastMessage(friendUserID: UUID, date: Date) {
        let dateString = ISO8601DateFormatter().string(from: date)
        let friendIDString = friendUserID.uuidString

        dbQueue.sync {
            let sql = "UPDATE conversations SET last_message_at = ? WHERE friend_user_id = ?;"
            var statement: OpaquePointer?

            if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
                dateString.withCString { dateStr in
                    friendIDString.withCString { idStr in
                        sqlite3_bind_text(statement, 1, dateStr, -1, nil)
                        sqlite3_bind_text(statement, 2, idStr, -1, nil)
                        sqlite3_step(statement)
                    }
                }
            }
            sqlite3_finalize(statement)
        }
    }

    func incrementUnreadCount(friendUserID: UUID) {
        let friendIDString = friendUserID.uuidString

        dbQueue.sync {
            let sql = "UPDATE conversations SET unread_count = unread_count + 1 WHERE friend_user_id = ?;"
            var statement: OpaquePointer?

            if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
                friendIDString.withCString { idStr in
                    sqlite3_bind_text(statement, 1, idStr, -1, nil)
                    sqlite3_step(statement)
                }
            }
            sqlite3_finalize(statement)
        }
    }

    func resetUnreadCount(friendUserID: UUID) {
        let friendIDString = friendUserID.uuidString

        dbQueue.sync {
            let sql = "UPDATE conversations SET unread_count = 0 WHERE friend_user_id = ?;"
            var statement: OpaquePointer?

            if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
                friendIDString.withCString { idStr in
                    sqlite3_bind_text(statement, 1, idStr, -1, nil)
                    sqlite3_step(statement)
                }
            }
            sqlite3_finalize(statement)
        }
    }

    private func conversationFromStatement(_ statement: OpaquePointer?) -> Conversation {
        let friendUserIDStr = String(cString: sqlite3_column_text(statement, 0))
        let friendUsername = String(cString: sqlite3_column_text(statement, 1))
        let friendPublicKey = String(cString: sqlite3_column_text(statement, 2))
        let friendSinceStr = sqlite3_column_text(statement, 3).map { String(cString: $0) } ?? ""
        let lastMessageAtStr = sqlite3_column_text(statement, 4).map { String(cString: $0) }
        let unreadCount = Int(sqlite3_column_int(statement, 5))
        let createdAtStr = String(cString: sqlite3_column_text(statement, 6))

        let formatter = ISO8601DateFormatter()

        return Conversation(
            friendUserID: UUID(uuidString: friendUserIDStr) ?? UUID(),
            friendUsername: friendUsername,
            friendPublicKey: friendPublicKey,
            friendSince: formatter.date(from: friendSinceStr) ?? Date(),
            lastMessageAt: lastMessageAtStr.flatMap { formatter.date(from: $0) },
            unreadCount: unreadCount,
            createdAt: formatter.date(from: createdAtStr) ?? Date()
        )
    }

    // MARK: - Message Operations

    func saveMessage(_ message: StoredMessage) {
        let formatter = ISO8601DateFormatter()
        let idString = message.id.uuidString
        let conversationIDString = message.conversationID.uuidString
        let contentTypeString = message.contentType.rawValue
        let createdAtString = formatter.string(from: message.createdAt)
        let receivedAtString = formatter.string(from: message.receivedAt)

        dbQueue.sync {
            let sql = """
                INSERT OR REPLACE INTO messages
                (id, conversation_id, content_type, decrypted_content, encrypted_content, is_from_me, has_been_viewed, created_at, received_at)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
            """
            var statement: OpaquePointer?

            if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
                idString.withCString { id in
                    conversationIDString.withCString { convID in
                        contentTypeString.withCString { contentType in
                            createdAtString.withCString { createdAt in
                                receivedAtString.withCString { receivedAt in
                                    sqlite3_bind_text(statement, 1, id, -1, nil)
                                    sqlite3_bind_text(statement, 2, convID, -1, nil)
                                    sqlite3_bind_text(statement, 3, contentType, -1, nil)
                                    sqlite3_bind_blob(statement, 4, (message.decryptedContent as NSData).bytes, Int32(message.decryptedContent.count), nil)
                                    if let encryptedContent = message.encryptedContent {
                                        sqlite3_bind_blob(statement, 5, (encryptedContent as NSData).bytes, Int32(encryptedContent.count), nil)
                                    } else {
                                        sqlite3_bind_null(statement, 5)
                                    }
                                    sqlite3_bind_int(statement, 6, message.isFromMe ? 1 : 0)
                                    sqlite3_bind_int(statement, 7, message.hasBeenViewed ? 1 : 0)
                                    sqlite3_bind_text(statement, 8, createdAt, -1, nil)
                                    sqlite3_bind_text(statement, 9, receivedAt, -1, nil)

                                    let result = sqlite3_step(statement)
                                    if result != SQLITE_DONE {
                                        print("Failed to save message: \(result)")
                                    } else {
                                        print("Saved message \(idString) to conversation \(conversationIDString)")
                                    }
                                }
                            }
                        }
                    }
                }
            } else {
                print("Failed to prepare saveMessage statement")
            }
            sqlite3_finalize(statement)
        }

        // Update conversation
        updateConversationLastMessage(friendUserID: message.conversationID, date: message.createdAt)
    }

    func getMessages(for conversationID: UUID, limit: Int = 50, offset: Int = 0) -> [StoredMessage] {
        dbQueue.sync {
            var messages: [StoredMessage] = []
            let conversationIDString = conversationID.uuidString
            let sql = """
                SELECT id, conversation_id, content_type, decrypted_content, encrypted_content, is_from_me, has_been_viewed, created_at, received_at
                FROM messages
                WHERE conversation_id = ?
                ORDER BY created_at ASC
                LIMIT ? OFFSET ?;
            """
            var statement: OpaquePointer?

            if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
                conversationIDString.withCString { cString in
                    sqlite3_bind_text(statement, 1, cString, -1, nil)
                    sqlite3_bind_int(statement, 2, Int32(limit))
                    sqlite3_bind_int(statement, 3, Int32(offset))

                    while sqlite3_step(statement) == SQLITE_ROW {
                        messages.append(messageFromStatement(statement))
                    }
                }
            } else {
                print("Failed to prepare getMessages query")
            }
            sqlite3_finalize(statement)
            print("getMessages for \(conversationIDString) returning \(messages.count) messages")
            return messages
        }
    }

    func getUnviewedMessages(for conversationID: UUID) -> [StoredMessage] {
        dbQueue.sync {
            var messages: [StoredMessage] = []
            let sql = """
                SELECT id, conversation_id, content_type, decrypted_content, encrypted_content, is_from_me, has_been_viewed, created_at, received_at
                FROM messages
                WHERE conversation_id = ? AND has_been_viewed = 0 AND is_from_me = 0
                ORDER BY created_at ASC;
            """
            var statement: OpaquePointer?

            if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
                sqlite3_bind_text(statement, 1, conversationID.uuidString, -1, nil)

                while sqlite3_step(statement) == SQLITE_ROW {
                    messages.append(messageFromStatement(statement))
                }
            }
            sqlite3_finalize(statement)
            return messages
        }
    }

    func markMessageAsViewed(messageID: UUID) {
        let messageIDString = messageID.uuidString

        dbQueue.sync {
            let sql = "UPDATE messages SET has_been_viewed = 1 WHERE id = ?;"
            var statement: OpaquePointer?

            if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
                messageIDString.withCString { idStr in
                    sqlite3_bind_text(statement, 1, idStr, -1, nil)
                    sqlite3_step(statement)
                }
            }
            sqlite3_finalize(statement)
        }
    }

    func messageExists(id: UUID) -> Bool {
        let idString = id.uuidString

        return dbQueue.sync {
            let sql = "SELECT 1 FROM messages WHERE id = ? LIMIT 1;"
            var statement: OpaquePointer?
            var exists = false

            if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
                idString.withCString { idStr in
                    sqlite3_bind_text(statement, 1, idStr, -1, nil)
                    exists = sqlite3_step(statement) == SQLITE_ROW
                }
            }
            sqlite3_finalize(statement)
            return exists
        }
    }

    func getTotalUnreadCount() -> Int {
        dbQueue.sync {
            let sql = "SELECT SUM(unread_count) FROM conversations;"
            var statement: OpaquePointer?
            var count = 0

            if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
                if sqlite3_step(statement) == SQLITE_ROW {
                    count = Int(sqlite3_column_int(statement, 0))
                }
            }
            sqlite3_finalize(statement)
            return count
        }
    }

    private func messageFromStatement(_ statement: OpaquePointer?) -> StoredMessage {
        let formatter = ISO8601DateFormatter()

        let idStr = String(cString: sqlite3_column_text(statement, 0))
        let conversationIDStr = String(cString: sqlite3_column_text(statement, 1))
        let contentTypeStr = String(cString: sqlite3_column_text(statement, 2))

        let contentBytes = sqlite3_column_blob(statement, 3)
        let contentLength = sqlite3_column_bytes(statement, 3)
        let content = Data(bytes: contentBytes!, count: Int(contentLength))

        // Read encrypted content (column 4) - may be NULL
        var encryptedContent: Data? = nil
        if let encryptedBytes = sqlite3_column_blob(statement, 4) {
            let encryptedLength = sqlite3_column_bytes(statement, 4)
            encryptedContent = Data(bytes: encryptedBytes, count: Int(encryptedLength))
        }

        let isFromMe = sqlite3_column_int(statement, 5) == 1
        let hasBeenViewed = sqlite3_column_int(statement, 6) == 1
        let createdAtStr = String(cString: sqlite3_column_text(statement, 7))
        let receivedAtStr = String(cString: sqlite3_column_text(statement, 8))

        return StoredMessage(
            id: UUID(uuidString: idStr) ?? UUID(),
            conversationID: UUID(uuidString: conversationIDStr) ?? UUID(),
            contentType: ContentType(rawValue: contentTypeStr) ?? .text,
            decryptedContent: content,
            encryptedContent: encryptedContent,
            isFromMe: isFromMe,
            hasBeenViewed: hasBeenViewed,
            createdAt: formatter.date(from: createdAtStr) ?? Date(),
            receivedAt: formatter.date(from: receivedAtStr) ?? Date()
        )
    }

    // MARK: - Cleanup

    func clearAll() {
        dbQueue.sync {
            executeSQL("DELETE FROM messages;")
            executeSQL("DELETE FROM conversations;")
        }
    }

    func purgeOldMessages(olderThan hours: Int = 24) {
        dbQueue.sync {
            let cutoffDate = Date().addingTimeInterval(-Double(hours) * 60 * 60)
            let formatter = ISO8601DateFormatter()

            let sql = "DELETE FROM messages WHERE received_at < ?;"
            var statement: OpaquePointer?

            if sqlite3_prepare_v2(db, sql, -1, &statement, nil) == SQLITE_OK {
                sqlite3_bind_text(statement, 1, formatter.string(from: cutoffDate), -1, nil)
                sqlite3_step(statement)
            }
            sqlite3_finalize(statement)
        }
    }

    private func startCleanupTimer() {
        Timer.scheduledTimer(withTimeInterval: 3600, repeats: true) { [weak self] _ in
            self?.purgeOldMessages()
        }
    }
}

// MARK: - Data Models

struct Conversation: Identifiable, Hashable {
    let friendUserID: UUID
    let friendUsername: String
    let friendPublicKey: String
    let friendSince: Date
    let lastMessageAt: Date?
    var unreadCount: Int
    let createdAt: Date

    var id: UUID { friendUserID }

    static func == (lhs: Conversation, rhs: Conversation) -> Bool {
        lhs.friendUserID == rhs.friendUserID
    }

    func hash(into hasher: inout Hasher) {
        hasher.combine(friendUserID)
    }
}

struct StoredMessage: Identifiable {
    let id: UUID
    let conversationID: UUID
    let contentType: ContentType
    let decryptedContent: Data
    let encryptedContent: Data?
    let isFromMe: Bool
    var hasBeenViewed: Bool
    let createdAt: Date
    let receivedAt: Date
}
