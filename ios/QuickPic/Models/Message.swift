import Foundation

enum ContentType: String, Codable {
    case text
    case image
}

struct Message: Codable, Identifiable {
    let id: UUID
    let fromUserID: UUID
    let toUserID: UUID
    let encryptedContent: Data
    let contentType: ContentType
    let signature: String
    let createdAt: Date
    let fromUsername: String
    let fromPublicKey: String

    enum CodingKeys: String, CodingKey {
        case id
        case fromUserID = "from_user_id"
        case toUserID = "to_user_id"
        case encryptedContent = "encrypted_content"
        case contentType = "content_type"
        case signature
        case createdAt = "created_at"
        case fromUsername = "from_username"
        case fromPublicKey = "from_public_key"
    }
}

struct SendMessageRequest: Codable {
    let toUsername: String
    let encryptedContent: Data
    let contentType: ContentType
    let signature: String

    enum CodingKeys: String, CodingKey {
        case toUsername = "to_username"
        case encryptedContent = "encrypted_content"
        case contentType = "content_type"
        case signature
    }
}

/// Decrypted message stored locally for 24-hour cache
struct CachedMessage: Codable, Identifiable {
    let id: UUID
    let fromUsername: String
    let contentType: ContentType
    let decryptedContent: Data
    let receivedAt: Date
    var hasBeenViewed: Bool

    var isExpired: Bool {
        Date().timeIntervalSince(receivedAt) > 24 * 60 * 60
    }
}
