import Foundation

struct Friend: Codable, Identifiable {
    let userID: UUID
    let username: String
    let publicKey: String
    let since: Date

    var id: UUID { userID }

    enum CodingKeys: String, CodingKey {
        case userID = "user_id"
        case username
        case publicKey = "public_key"
        case since
    }
}

struct FriendRequest: Codable, Identifiable {
    let id: UUID
    let fromUserID: UUID
    let toUserID: UUID
    let status: String
    let createdAt: Date
    let fromUser: User

    enum CodingKeys: String, CodingKey {
        case id
        case fromUserID = "from_user_id"
        case toUserID = "to_user_id"
        case status
        case createdAt = "created_at"
        case fromUser = "from_user"
    }
}

struct SendFriendRequestRequest: Codable {
    let username: String
}

struct FriendRequestActionRequest: Codable {
    let requestID: UUID

    enum CodingKeys: String, CodingKey {
        case requestID = "request_id"
    }
}
