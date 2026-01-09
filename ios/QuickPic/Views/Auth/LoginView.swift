//
//  LoginView.swift
//  QuickPic
//

import SwiftUI

struct LoginView: View {
    @EnvironmentObject var authManager: AuthManager
    @State private var username = ""
    @State private var password = ""
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var showRegister = false

    var body: some View {
        NavigationStack {
            ZStack {
                Color.appBackground.ignoresSafeArea()

                VStack(spacing: 0) {
                    Spacer()

                    // Header
                    VStack(spacing: AppSpacing.sm) {
                        Text("Welcome")
                            .font(.system(size: 34, weight: .bold))
                            .foregroundColor(.textPrimary)

                        Text("Back !")
                            .font(.system(size: 34, weight: .bold))
                            .foregroundColor(.textPrimary)

                        Text("Sign in to continue")
                            .font(.appCaption)
                            .foregroundColor(.textSecondary)
                            .padding(.top, AppSpacing.xs)
                    }

                    Spacer()

                    // Form
                    VStack(spacing: AppSpacing.md) {
                        AppTextField(
                            icon: "person",
                            placeholder: "Username",
                            text: $username
                        )
                        .textContentType(.username)
                        .autocapitalization(.none)
                        .autocorrectionDisabled()

                        AppTextField(
                            icon: "lock",
                            placeholder: "Password",
                            text: $password,
                            isSecure: true
                        )
                        .textContentType(.password)

                        if let error = errorMessage {
                            Text(error)
                                .font(.appCaption)
                                .foregroundColor(.danger)
                                .frame(maxWidth: .infinity, alignment: .leading)
                                .padding(.horizontal, AppSpacing.xs)
                        }

                        Button(action: login) {
                            HStack {
                                if isLoading {
                                    ProgressView()
                                        .tint(.black)
                                } else {
                                    Text("Sign in")
                                }
                            }
                        }
                        .buttonStyle(PrimaryButtonStyle(isEnabled: canLogin))
                        .disabled(!canLogin)
                        .padding(.top, AppSpacing.sm)
                    }
                    .padding(.horizontal, AppSpacing.xl)

                    Spacer()

                    // Register link
                    HStack(spacing: AppSpacing.xs) {
                        Text("Don't have account?")
                            .font(.appCaption)
                            .foregroundColor(.textSecondary)

                        Button("Create now") {
                            Haptics.light()
                            showRegister = true
                        }
                        .font(.appCaption)
                        .fontWeight(.semibold)
                        .foregroundColor(.appPrimary)
                    }
                    .padding(.bottom, AppSpacing.md)

                    Button(action: debugCreateReece) {
                        Text("Debug Create REECE")
                            .font(.appCaption)
                            .fontWeight(.medium)
                    }
                    .buttonStyle(.bordered)
                    .tint(.orange)
                    .padding(.bottom, AppSpacing.xl)
                }
            }
            .navigationDestination(isPresented: $showRegister) {
                RegisterView()
            }
        }
    }

    private var canLogin: Bool {
        !username.isEmpty && !password.isEmpty && !isLoading
    }

    private func login() {
        Haptics.light()
        isLoading = true
        errorMessage = nil

        Task {
            do {
                try await authManager.login(username: username, password: password)
                Haptics.success()
            } catch APIError.unauthorized {
                errorMessage = "Invalid username or password"
                Haptics.error()
            } catch APIError.httpError(_, let message) {
                errorMessage = message ?? "Login failed"
                Haptics.error()
            } catch {
                errorMessage = "Login failed. Please try again."
                Haptics.error()
            }
            isLoading = false
        }
    }

    private func debugCreateReece() {
        Haptics.light()
        isLoading = true
        errorMessage = nil

        Task {
            do {
                try await authManager.register(username: "reecepbcups", password: "password")
                Haptics.success()
            } catch APIError.httpError(409, _) {
                errorMessage = "Username already exists"
                Haptics.error()
            } catch APIError.httpError(_, let message) {
                errorMessage = message ?? "Debug registration failed"
                Haptics.error()
            } catch {
                errorMessage = "Debug registration failed. Please try again."
                Haptics.error()
            }
            isLoading = false
        }
    }
}

#Preview {
    LoginView()
        .environmentObject(AuthManager())
}
