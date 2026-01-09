//
//  AddFriendView.swift
//  QuickPic
//

import SwiftUI

struct AddFriendSheet: View {
    @Environment(\.dismiss) var dismiss
    @State private var username = ""
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var successMessage: String?

    var onRequestSent: () -> Void

    var body: some View {
        ZStack {
            Color.appBackground.ignoresSafeArea()

            VStack(spacing: AppSpacing.lg) {
                // Header
                VStack(spacing: AppSpacing.sm) {
                    Text("Add Friend")
                        .font(.appTitle)
                        .foregroundColor(.textPrimary)

                    Text("Enter their username to send a request")
                        .font(.appCaption)
                        .foregroundColor(.textSecondary)
                }
                .padding(.top, AppSpacing.lg)

                // Form
                VStack(spacing: AppSpacing.md) {
                    AppTextField(
                        icon: "at",
                        placeholder: "Username",
                        text: $username
                    )
                    .autocapitalization(.none)
                    .autocorrectionDisabled()

                    if let error = errorMessage {
                        HStack {
                            Image(systemName: "exclamationmark.circle.fill")
                                .foregroundColor(.danger)
                            Text(error)
                                .foregroundColor(.danger)
                        }
                        .font(.appCaption)
                    }

                    if let success = successMessage {
                        HStack {
                            Image(systemName: "checkmark.circle.fill")
                                .foregroundColor(.success)
                            Text(success)
                                .foregroundColor(.success)
                        }
                        .font(.appCaption)
                    }

                    Button(action: sendRequest) {
                        HStack {
                            if isLoading {
                                ProgressView()
                                    .tint(.black)
                            } else {
                                Text("Send Request")
                            }
                        }
                    }
                    .buttonStyle(PrimaryButtonStyle(isEnabled: !username.isEmpty && !isLoading))
                    .disabled(username.isEmpty || isLoading)
                    .padding(.top, AppSpacing.sm)
                }
                .padding(.horizontal, AppSpacing.lg)

                Spacer()
            }
        }
    }

    private func sendRequest() {
        Haptics.light()
        isLoading = true
        errorMessage = nil
        successMessage = nil

        Task {
            do {
                _ = try await APIService.shared.sendFriendRequest(to: username)
                successMessage = "Request sent to @\(username)"
                onRequestSent()
                Haptics.success()

                try? await Task.sleep(nanoseconds: 1_500_000_000)
                dismiss()
            } catch APIError.notFound {
                errorMessage = "User not found"
                Haptics.error()
            } catch APIError.httpError(409, _) {
                errorMessage = "Already sent or already friends"
                Haptics.error()
            } catch APIError.httpError(400, let message) {
                errorMessage = message ?? "Cannot add this user"
                Haptics.error()
            } catch {
                errorMessage = "Failed to send request"
                Haptics.error()
            }
            isLoading = false
        }
    }
}

#Preview {
    AddFriendSheet(onRequestSent: {})
}
