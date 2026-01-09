//
//  RegisterView.swift
//  QuickPic
//

import SwiftUI

struct RegisterView: View {
    @EnvironmentObject var authManager: AuthManager
    @Environment(\.dismiss) var dismiss
    @State private var username = ""
    @State private var password = ""
    @State private var confirmPassword = ""
    @State private var isLoading = false
    @State private var errorMessage: String?

    private var passwordsMatch: Bool {
        password == confirmPassword && !password.isEmpty
    }

    private var isValidForm: Bool {
        username.count >= 3 && password.count >= 8 && passwordsMatch
    }

    var body: some View {
        ZStack {
            Color.appBackground.ignoresSafeArea()

            VStack(spacing: 0) {
                Spacer()

                // Header
                VStack(spacing: AppSpacing.sm) {
                    Text("Create")
                        .font(.system(size: 34, weight: .bold))
                        .foregroundColor(.textPrimary)

                    Text("Account")
                        .font(.system(size: 34, weight: .bold))
                        .foregroundColor(.textPrimary)

                    Text("Your messages are end-to-end encrypted")
                        .font(.appCaption)
                        .foregroundColor(.textSecondary)
                        .padding(.top, AppSpacing.xs)
                }

                Spacer()

                // Form
                VStack(spacing: AppSpacing.md) {
                    VStack(alignment: .leading, spacing: AppSpacing.xs) {
                        AppTextField(
                            icon: "person",
                            placeholder: "Username",
                            text: $username
                        )
                        .textContentType(.username)
                        .autocapitalization(.none)
                        .autocorrectionDisabled()

                        if !username.isEmpty && username.count < 3 {
                            Text("Username must be at least 3 characters")
                                .font(.appSmall)
                                .foregroundColor(.pending)
                                .padding(.leading, AppSpacing.xs)
                        }
                    }

                    VStack(alignment: .leading, spacing: AppSpacing.xs) {
                        AppTextField(
                            icon: "lock",
                            placeholder: "Password",
                            text: $password,
                            isSecure: true
                        )
                        .textContentType(.newPassword)

                        if !password.isEmpty && password.count < 8 {
                            Text("Password must be at least 8 characters")
                                .font(.appSmall)
                                .foregroundColor(.pending)
                                .padding(.leading, AppSpacing.xs)
                        }
                    }

                    VStack(alignment: .leading, spacing: AppSpacing.xs) {
                        AppTextField(
                            icon: "lock",
                            placeholder: "Confirm Password",
                            text: $confirmPassword,
                            isSecure: true
                        )
                        .textContentType(.newPassword)

                        if !confirmPassword.isEmpty && !passwordsMatch {
                            Text("Passwords do not match")
                                .font(.appSmall)
                                .foregroundColor(.danger)
                                .padding(.leading, AppSpacing.xs)
                        }
                    }

                    if let error = errorMessage {
                        Text(error)
                            .font(.appCaption)
                            .foregroundColor(.danger)
                            .multilineTextAlignment(.center)
                    }

                    Button(action: register) {
                        HStack {
                            if isLoading {
                                ProgressView()
                                    .tint(.black)
                            } else {
                                Text("Create Account")
                            }
                        }
                    }
                    .buttonStyle(PrimaryButtonStyle(isEnabled: isValidForm && !isLoading))
                    .disabled(!isValidForm || isLoading)
                    .padding(.top, AppSpacing.sm)
                }
                .padding(.horizontal, AppSpacing.xl)

                Spacer()

                // Security note
                VStack(spacing: AppSpacing.sm) {
                    Image(systemName: "lock.shield.fill")
                        .font(.title2)
                        .foregroundColor(.success)

                    Text("A unique encryption key will be generated on your device. Your private key never leaves this device.")
                        .font(.appSmall)
                        .foregroundColor(.textSecondary)
                        .multilineTextAlignment(.center)
                }
                .padding(.horizontal, AppSpacing.xl)
                .padding(.bottom, AppSpacing.xl)
            }
        }
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            ToolbarItem(placement: .navigationBarLeading) {
                Button(action: { dismiss() }) {
                    Image(systemName: "arrow.left")
                        .foregroundColor(.textPrimary)
                }
            }
        }
    }

    private func register() {
        Haptics.light()
        isLoading = true
        errorMessage = nil

        Task {
            do {
                try await authManager.register(username: username, password: password)
                Haptics.success()
            } catch APIError.httpError(409, _) {
                errorMessage = "Username already exists"
                Haptics.error()
            } catch APIError.httpError(_, let message) {
                errorMessage = message ?? "Registration failed"
                Haptics.error()
            } catch {
                errorMessage = "Registration failed. Please try again."
                Haptics.error()
            }
            isLoading = false
        }
    }
}

#Preview {
    NavigationStack {
        RegisterView()
            .environmentObject(AuthManager())
    }
}
