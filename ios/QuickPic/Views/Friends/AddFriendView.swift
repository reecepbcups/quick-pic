import SwiftUI

struct AddFriendView: View {
    @Environment(\.dismiss) var dismiss
    @State private var username = ""
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var successMessage: String?

    var onRequestSent: () -> Void

    var body: some View {
        NavigationStack {
            VStack(spacing: 24) {
                // Header
                VStack(spacing: 8) {
                    Image(systemName: "person.badge.plus")
                        .font(.system(size: 50))
                        .foregroundColor(.yellow)

                    Text("Add a Friend")
                        .font(.title2)
                        .fontWeight(.bold)

                    Text("Enter their username to send a friend request")
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                        .multilineTextAlignment(.center)
                }
                .padding(.top, 32)

                // Form
                VStack(spacing: 16) {
                    TextField("Username", text: $username)
                        .textFieldStyle(.roundedBorder)
                        .autocapitalization(.none)
                        .autocorrectionDisabled()

                    if let error = errorMessage {
                        Text(error)
                            .font(.caption)
                            .foregroundColor(.red)
                    }

                    if let success = successMessage {
                        HStack {
                            Image(systemName: "checkmark.circle.fill")
                                .foregroundColor(.green)
                            Text(success)
                                .foregroundColor(.green)
                        }
                        .font(.caption)
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
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(username.isEmpty ? Color.gray : Color.yellow)
                        .foregroundColor(.black)
                        .cornerRadius(10)
                    }
                    .disabled(isLoading || username.isEmpty)
                }
                .padding(.horizontal, 32)

                Spacer()
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
            }
        }
    }

    private func sendRequest() {
        isLoading = true
        errorMessage = nil
        successMessage = nil

        Task {
            do {
                _ = try await APIService.shared.sendFriendRequest(to: username)
                successMessage = "Friend request sent to @\(username)"
                onRequestSent()

                // Dismiss after a short delay
                try? await Task.sleep(nanoseconds: 1_500_000_000)
                dismiss()
            } catch APIError.notFound {
                errorMessage = "User not found"
            } catch APIError.httpError(409, _) {
                errorMessage = "Request already sent or already friends"
            } catch APIError.httpError(400, let message) {
                errorMessage = message ?? "Cannot add this user"
            } catch {
                errorMessage = "Failed to send request"
            }
            isLoading = false
        }
    }
}

#Preview {
    AddFriendView(onRequestSent: {})
}
