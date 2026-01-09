//
//  InboxView.swift
//  QuickPic
//

import SwiftUI

struct InboxView: View {
    @StateObject private var viewModel = InboxViewModel()
    @State private var selectedMessage: CachedMessage?

    var body: some View {
        NavigationStack {
            Group {
                if viewModel.isLoading && viewModel.messages.isEmpty {
                    ProgressView("Loading messages...")
                } else if viewModel.messages.isEmpty {
                    EmptyInboxView()
                } else {
                    messagesList
                }
            }
            .navigationTitle("Inbox")
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: { Task { await viewModel.refresh() } }) {
                        Image(systemName: "arrow.clockwise")
                    }
                }
            }
            .refreshable {
                await viewModel.refresh()
            }
            .alert("Error", isPresented: .constant(viewModel.errorMessage != nil)) {
                Button("OK") { viewModel.errorMessage = nil }
            } message: {
                Text(viewModel.errorMessage ?? "")
            }
            .fullScreenCover(item: $selectedMessage) { message in
                MessageView(message: message) {
                    viewModel.markAsViewed(message)
                    selectedMessage = nil
                }
            }
            .task {
                await viewModel.loadMessages()
            }
        }
    }

    private var messagesList: some View {
        List {
            ForEach(viewModel.messages) { message in
                MessageRow(message: message)
                    .onTapGesture {
                        selectedMessage = message
                    }
            }
        }
        .listStyle(.plain)
    }
}

struct MessageRow: View {
    let message: CachedMessage

    var body: some View {
        HStack(spacing: 12) {
            // Avatar
            ZStack {
                Circle()
                    .fill(message.hasBeenViewed ? Color.gray.opacity(0.3) : Color.yellow.opacity(0.3))
                    .frame(width: 50, height: 50)

                if message.contentType == .image {
                    Image(systemName: "photo.fill")
                        .foregroundColor(message.hasBeenViewed ? .gray : .yellow)
                } else {
                    Image(systemName: "text.bubble.fill")
                        .foregroundColor(message.hasBeenViewed ? .gray : .yellow)
                }
            }

            VStack(alignment: .leading, spacing: 4) {
                Text(message.fromUsername)
                    .fontWeight(message.hasBeenViewed ? .regular : .semibold)

                HStack {
                    Text(message.contentType == .image ? "Sent a photo" : "Sent a message")
                        .font(.subheadline)
                        .foregroundColor(.secondary)

                    if !message.hasBeenViewed {
                        Circle()
                            .fill(Color.yellow)
                            .frame(width: 8, height: 8)
                    }
                }

                Text(timeAgo(message.receivedAt))
                    .font(.caption)
                    .foregroundColor(.secondary)
            }

            Spacer()

            Image(systemName: "chevron.right")
                .foregroundColor(.secondary)
                .font(.caption)
        }
        .padding(.vertical, 4)
    }

    private func timeAgo(_ date: Date) -> String {
        let formatter = RelativeDateTimeFormatter()
        formatter.unitsStyle = .abbreviated
        return formatter.localizedString(for: date, relativeTo: Date())
    }
}

struct EmptyInboxView: View {
    var body: some View {
        VStack(spacing: 16) {
            Image(systemName: "tray")
                .font(.system(size: 50))
                .foregroundColor(.secondary)

            Text("No messages yet")
                .font(.headline)

            Text("Messages from friends will appear here")
                .font(.subheadline)
                .foregroundColor(.secondary)
        }
    }
}

@MainActor
class InboxViewModel: ObservableObject {
    @Published var messages: [CachedMessage] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let api = APIService.shared
    private let cache = MessageCacheService.shared
    private let crypto = CryptoService.shared
    private let keychain = KeychainService.shared

    func loadMessages() async {
        isLoading = true
        defer { isLoading = false }

        // Load cached messages first
        messages = cache.getCachedMessages().sorted { $0.receivedAt > $1.receivedAt }

        // Fetch new messages from server
        await refresh()
    }

    func refresh() async {
        do {
            let serverMessages = try await api.getMessages()

            for message in serverMessages {
                await processMessage(message)
            }

            // Reload from cache
            messages = cache.getCachedMessages().sorted { $0.receivedAt > $1.receivedAt }
        } catch {
            errorMessage = "Failed to fetch messages"
        }
    }

    private func processMessage(_ message: Message) async {
        do {
            // Get sender's public key and decrypt
            let senderPublicKey = try crypto.publicKeyFromBase64(message.fromPublicKey)
            let privateKey = try crypto.getPrivateKey()

            let decryptedData = try crypto.decrypt(
                encryptedData: message.encryptedContent,
                signature: message.signature,
                senderPublicKey: senderPublicKey,
                recipientPrivateKey: privateKey
            )

            // Create cached message
            let cachedMessage = CachedMessage(
                id: message.id,
                fromUsername: message.fromUsername,
                contentType: message.contentType,
                decryptedContent: decryptedData,
                receivedAt: Date(),
                hasBeenViewed: false
            )

            cache.cache(message: cachedMessage)

            // Acknowledge receipt to server
            try await api.acknowledgeMessage(id: message.id)
        } catch {
            print("Failed to process message \(message.id): \(error)")
        }
    }

    func markAsViewed(_ message: CachedMessage) {
        cache.markAsViewed(messageID: message.id)
        if let index = messages.firstIndex(where: { $0.id == message.id }) {
            messages[index].hasBeenViewed = true
        }
    }
}

#Preview {
    InboxView()
}
