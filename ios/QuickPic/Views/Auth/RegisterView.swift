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
        VStack(spacing: 24) {
            // Header
            VStack(spacing: 8) {
                Text("Create Account")
                    .font(.largeTitle)
                    .fontWeight(.bold)

                Text("Your messages are end-to-end encrypted")
                    .font(.subheadline)
                    .foregroundColor(.secondary)
            }
            .padding(.top, 32)

            // Form
            VStack(spacing: 16) {
                VStack(alignment: .leading, spacing: 4) {
                    TextField("Username", text: $username)
                        .textFieldStyle(.roundedBorder)
                        .textContentType(.username)
                        .autocapitalization(.none)
                        .autocorrectionDisabled()

                    if !username.isEmpty && username.count < 3 {
                        Text("Username must be at least 3 characters")
                            .font(.caption)
                            .foregroundColor(.orange)
                    }
                }

                VStack(alignment: .leading, spacing: 4) {
                    SecureField("Password", text: $password)
                        .textFieldStyle(.roundedBorder)
                        .textContentType(.newPassword)

                    if !password.isEmpty && password.count < 8 {
                        Text("Password must be at least 8 characters")
                            .font(.caption)
                            .foregroundColor(.orange)
                    }
                }

                VStack(alignment: .leading, spacing: 4) {
                    SecureField("Confirm Password", text: $confirmPassword)
                        .textFieldStyle(.roundedBorder)
                        .textContentType(.newPassword)

                    if !confirmPassword.isEmpty && !passwordsMatch {
                        Text("Passwords do not match")
                            .font(.caption)
                            .foregroundColor(.red)
                    }
                }

                if let error = errorMessage {
                    Text(error)
                        .font(.caption)
                        .foregroundColor(.red)
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
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(isValidForm ? Color.yellow : Color.gray)
                    .foregroundColor(.black)
                    .cornerRadius(10)
                }
                .disabled(isLoading || !isValidForm)
            }
            .padding(.horizontal, 32)

            Spacer()

            // Security Note
            VStack(spacing: 8) {
                Image(systemName: "lock.shield.fill")
                    .font(.title2)
                    .foregroundColor(.green)

                Text("A unique encryption key will be generated on your device. Your private key never leaves this device.")
                    .font(.caption)
                    .foregroundColor(.secondary)
                    .multilineTextAlignment(.center)
            }
            .padding(.horizontal, 32)
            .padding(.bottom, 32)
        }
        .navigationBarTitleDisplayMode(.inline)
    }

    private func register() {
        isLoading = true
        errorMessage = nil

        Task {
            do {
                try await authManager.register(username: username, password: password)
            } catch APIError.httpError(409, _) {
                errorMessage = "Username already exists"
            } catch APIError.httpError(_, let message) {
                errorMessage = message ?? "Registration failed"
            } catch {
                errorMessage = "Registration failed. Please try again."
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
