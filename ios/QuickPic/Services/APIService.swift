import Foundation

enum APIError: Error {
    case invalidURL
    case requestFailed(Error)
    case invalidResponse
    case httpError(statusCode: Int, message: String?)
    case decodingError(Error)
    case unauthorized
    case notFound
}

final class APIService: @unchecked Sendable {
    static let shared = APIService()

    private let baseURL: String
    private let session: URLSession
    private let decoder: JSONDecoder
    private let encoder: JSONEncoder

    private init() {
        #if DEBUG
        // Use your computer's local IP for testing on physical device
        // Make sure phone and computer are on the same WiFi network
        self.baseURL = "http://172.21.11.82:8080"
        #else
        self.baseURL = "https://api.quickpic.app"  // Production URL
        #endif

        let config = URLSessionConfiguration.default
        config.timeoutIntervalForRequest = 30
        self.session = URLSession(configuration: config)

        self.decoder = JSONDecoder()
        self.decoder.dateDecodingStrategy = .custom { decoder in
            let container = try decoder.singleValueContainer()
            let dateString = try container.decode(String.self)

            // Try ISO8601 with fractional seconds first (Go's default format)
            let formatterWithFractional = ISO8601DateFormatter()
            formatterWithFractional.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
            if let date = formatterWithFractional.date(from: dateString) {
                return date
            }

            // Fall back to ISO8601 without fractional seconds
            let formatter = ISO8601DateFormatter()
            formatter.formatOptions = [.withInternetDateTime]
            if let date = formatter.date(from: dateString) {
                return date
            }

            throw DecodingError.dataCorruptedError(in: container, debugDescription: "Cannot decode date: \(dateString)")
        }

        self.encoder = JSONEncoder()
        self.encoder.dateEncodingStrategy = .iso8601
    }

    // MARK: - Auth Endpoints

    func register(username: String, password: String, publicKey: String) async throws -> AuthResponse {
        let request = RegisterRequest(username: username, password: password, publicKey: publicKey)
        return try await post("/auth/register", body: request)
    }

    func login(username: String, password: String) async throws -> AuthResponse {
        let request = LoginRequest(username: username, password: password)
        return try await post("/auth/login", body: request)
    }

    func refreshToken(_ refreshToken: String) async throws -> AuthResponse {
        let request = RefreshRequest(refreshToken: refreshToken)
        return try await post("/auth/refresh", body: request, allowRetry: false)
    }

    func logout(refreshToken: String) async throws {
        let request = RefreshRequest(refreshToken: refreshToken)
        let _: EmptyResponse = try await post("/auth/logout", body: request)
    }

    // MARK: - User Endpoints

    func getUser(username: String) async throws -> User {
        try await get("/users/\(username)", authenticated: true)
    }

    // MARK: - Friend Endpoints

    func sendFriendRequest(to username: String) async throws -> FriendRequest {
        let request = SendFriendRequestRequest(username: username)
        return try await post("/friends/request", body: request, authenticated: true)
    }

    func getPendingFriendRequests() async throws -> [FriendRequest] {
        try await get("/friends/requests", authenticated: true)
    }

    func acceptFriendRequest(requestID: UUID) async throws {
        let request = FriendRequestActionRequest(requestID: requestID)
        let _: EmptyResponse = try await post("/friends/accept", body: request, authenticated: true)
    }

    func rejectFriendRequest(requestID: UUID) async throws {
        let request = FriendRequestActionRequest(requestID: requestID)
        let _: EmptyResponse = try await post("/friends/reject", body: request, authenticated: true)
    }

    func getFriends() async throws -> [Friend] {
        try await get("/friends", authenticated: true)
    }

    // MARK: - Message Endpoints

    func sendMessage(to username: String, encryptedContent: Data, contentType: ContentType, signature: String) async throws -> MessageSendResponse {
        let request = SendMessageRequest(
            toUsername: username,
            encryptedContent: encryptedContent,
            contentType: contentType,
            signature: signature
        )
        return try await post("/messages", body: request, authenticated: true)
    }

    func getMessages() async throws -> [Message] {
        try await get("/messages", authenticated: true)
    }

    func acknowledgeMessage(id: UUID) async throws {
        let _: EmptyResponse = try await post("/messages/\(id.uuidString)/ack", body: EmptyRequest(), authenticated: true)
    }

    // MARK: - Private Helpers

    private func get<T: Decodable>(_ path: String, authenticated: Bool = false) async throws -> T {
        var request = try makeRequest(path: path, method: "GET")

        if authenticated {
            try addAuthHeader(to: &request)
        }

        return try await execute(request)
    }

    private func post<T: Decodable, B: Encodable>(_ path: String, body: B, authenticated: Bool = false, allowRetry: Bool = true) async throws -> T {
        var request = try makeRequest(path: path, method: "POST")
        request.httpBody = try encoder.encode(body)
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        if authenticated {
            try addAuthHeader(to: &request)
        }

        return try await execute(request, allowRetry: allowRetry)
    }

    private func makeRequest(path: String, method: String) throws -> URLRequest {
        guard let url = URL(string: baseURL + path) else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = method
        return request
    }

    private func addAuthHeader(to request: inout URLRequest) throws {
        let accessToken = try KeychainService.shared.getAccessToken()
        request.setValue("Bearer \(accessToken)", forHTTPHeaderField: "Authorization")
    }

    private func execute<T: Decodable>(_ request: URLRequest, allowRetry: Bool = true) async throws -> T {
        let (data, response) = try await session.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        switch httpResponse.statusCode {
        case 200...299:
            do {
                return try decoder.decode(T.self, from: data)
            } catch {
                throw APIError.decodingError(error)
            }

        case 401:
            // Try to refresh token (but not if this is already a refresh/auth request)
            if allowRetry, let refreshedResponse = try? await refreshAndRetry(request) as T {
                return refreshedResponse
            }
            throw APIError.unauthorized

        case 404:
            throw APIError.notFound

        default:
            let message = try? decoder.decode(ErrorResponse.self, from: data).error
            throw APIError.httpError(statusCode: httpResponse.statusCode, message: message)
        }
    }

    private func refreshAndRetry<T: Decodable>(_ originalRequest: URLRequest) async throws -> T {
        let storedRefreshToken = try KeychainService.shared.getRefreshToken()
        let authResponse = try await refreshToken(storedRefreshToken)

        // Store new tokens
        try KeychainService.shared.storeTokens(
            accessToken: authResponse.accessToken,
            refreshToken: authResponse.refreshToken
        )

        // Retry original request with new token
        var newRequest = originalRequest
        newRequest.setValue("Bearer \(authResponse.accessToken)", forHTTPHeaderField: "Authorization")

        let (data, response) = try await session.data(for: newRequest)

        guard let httpResponse = response as? HTTPURLResponse,
              (200...299).contains(httpResponse.statusCode) else {
            throw APIError.unauthorized
        }

        return try decoder.decode(T.self, from: data)
    }
}

// MARK: - Helper Types

private struct EmptyRequest: Codable {}

private struct EmptyResponse: Codable {
    let message: String?

    init(from decoder: Decoder) throws {
        let container = try? decoder.container(keyedBy: CodingKeys.self)
        message = try? container?.decode(String.self, forKey: .message)
    }

    enum CodingKeys: String, CodingKey {
        case message
    }
}

private struct ErrorResponse: Codable {
    let error: String
}

struct MessageSendResponse: Codable {
    let id: UUID
    let createdAt: Date

    enum CodingKeys: String, CodingKey {
        case id
        case createdAt = "created_at"
    }
}
