//
//  AuthManager.swift
//  QuickPic
//

import Foundation
import SwiftUI

@MainActor
final class AuthManager: ObservableObject {
    @Published private(set) var isAuthenticated = false
    @Published private(set) var currentUser: User?
    @Published private(set) var isLoading = true

    private let api = APIService.shared
    private let keychain = KeychainService.shared
    private let crypto = CryptoService.shared

    init() {
        Task {
            await checkAuthState()
        }
    }

    private func checkAuthState() async {
        defer { isLoading = false }

        // Check if we have stored user and tokens
        guard let _ = try? keychain.getAccessToken(),
              let user = try? keychain.getCurrentUser() else {
            isAuthenticated = false
            return
        }

        // Try to refresh token to validate session
        do {
            let refreshToken = try keychain.getRefreshToken()
            let response = try await api.refreshToken(refreshToken)
            try saveAuthResponse(response)
            isAuthenticated = true
            currentUser = response.user
        } catch {
            // Token invalid, clear and require re-login
            logout()
        }
    }

    func register(username: String, password: String) async throws {
        // Generate key pair first
        let (privateKey, publicKey) = crypto.generateKeyPair()
        let publicKeyBase64 = crypto.publicKeyToBase64(publicKey)

        // Store private key
        try crypto.storePrivateKey(privateKey)

        // Register with server
        let response = try await api.register(
            username: username,
            password: password,
            publicKey: publicKeyBase64
        )

        try saveAuthResponse(response)
        isAuthenticated = true
        currentUser = response.user
    }

    func login(username: String, password: String) async throws {
        let response = try await api.login(username: username, password: password)

        try saveAuthResponse(response)
        isAuthenticated = true
        currentUser = response.user
    }

    func logout() {
        // Try to notify server (fire and forget)
        if let refreshToken = try? keychain.getRefreshToken() {
            Task {
                try? await api.logout(refreshToken: refreshToken)
            }
        }

        // Clear local data
        keychain.clearAll()
        MessageCacheService.shared.clearAll()

        isAuthenticated = false
        currentUser = nil
    }

    private func saveAuthResponse(_ response: AuthResponse) throws {
        try keychain.storeTokens(
            accessToken: response.accessToken,
            refreshToken: response.refreshToken
        )
        try keychain.storeCurrentUser(response.user)
    }
}
