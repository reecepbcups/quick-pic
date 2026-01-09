import Foundation

struct User: Codable, Identifiable {
    let id: UUID
    let userNumber: Int64
    let username: String
    let publicKey: String

    enum CodingKeys: String, CodingKey {
        case id
        case userNumber = "user_number"
        case username
        case publicKey = "public_key"
    }
}

struct AuthResponse: Codable {
    let accessToken: String
    let refreshToken: String
    let expiresIn: Int64
    let user: User

    enum CodingKeys: String, CodingKey {
        case accessToken = "access_token"
        case refreshToken = "refresh_token"
        case expiresIn = "expires_in"
        case user
    }
}

struct RegisterRequest: Codable {
    let username: String
    let password: String
    let publicKey: String

    enum CodingKeys: String, CodingKey {
        case username
        case password
        case publicKey = "public_key"
    }
}

struct LoginRequest: Codable {
    let username: String
    let password: String
}

struct RefreshRequest: Codable {
    let refreshToken: String

    enum CodingKeys: String, CodingKey {
        case refreshToken = "refresh_token"
    }
}
