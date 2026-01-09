//
//  InboxView.swift
//  QuickPic
//
//  Conversation-based inbox showing friends with message threads
//

import SwiftUI

struct InboxView: View {
    @StateObject private var viewModel = InboxViewModel()
    @State private var selectedConversation: Conversation?

    var body: some View {
        NavigationStack {
            Group {
                if viewModel.isLoading && viewModel.conversations.isEmpty {
                    ProgressView("Loading conversations...")
                } else if viewModel.conversations.isEmpty {
                    EmptyInboxView()
                } else {
                    conversationsList
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
            .navigationDestination(item: $selectedConversation) { conversation in
                ChatView(conversation: conversation, onMessagesViewed: {
                    viewModel.markConversationRead(conversation)
                })
            }
            .task {
                await viewModel.loadConversations()
            }
        }
    }

    private var conversationsList: some View {
        List {
            ForEach(viewModel.conversations) { conversation in
                ConversationRow(conversation: conversation)
                    .onTapGesture {
                        selectedConversation = conversation
                    }
            }
        }
        .listStyle(.plain)
    }
}

struct ConversationRow: View {
    let conversation: Conversation

    var body: some View {
        HStack(spacing: 12) {
            // Avatar
            ZStack {
                Circle()
                    .fill(conversation.unreadCount > 0 ? Color.yellow.opacity(0.3) : Color.gray.opacity(0.3))
                    .frame(width: 50, height: 50)

                Text(conversation.friendUsername.prefix(1).uppercased())
                    .fontWeight(.semibold)
                    .foregroundColor(conversation.unreadCount > 0 ? .yellow : .primary)
            }

            VStack(alignment: .leading, spacing: 4) {
                Text(conversation.friendUsername)
                    .fontWeight(conversation.unreadCount > 0 ? .semibold : .regular)

                if let lastMessage = conversation.lastMessageAt {
                    Text(timeAgo(lastMessage))
                        .font(.caption)
                        .foregroundColor(.secondary)
                } else {
                    Text("Start a conversation")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
            }

            Spacer()

            if conversation.unreadCount > 0 {
                ZStack {
                    Circle()
                        .fill(Color.yellow)
                        .frame(width: 24, height: 24)

                    Text("\(conversation.unreadCount)")
                        .font(.caption2)
                        .fontWeight(.bold)
                        .foregroundColor(.black)
                }
            }

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
            Image(systemName: "bubble.left.and.bubble.right")
                .font(.system(size: 50))
                .foregroundColor(.secondary)

            Text("No conversations yet")
                .font(.headline)

            Text("Add friends to start chatting")
                .font(.subheadline)
                .foregroundColor(.secondary)
        }
    }
}

@MainActor
class InboxViewModel: ObservableObject {
    @Published var conversations: [Conversation] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let api = APIService.shared
    private let db = DatabaseService.shared
    private let crypto = CryptoService.shared

    func loadConversations() async {
        isLoading = true

        // First sync friends to create conversations
        await syncFriendsToConversations()

        // Fetch new messages from server
        await fetchNewMessages()

        // Load all conversations from database
        conversations = db.getAllConversations()
        print("Loaded \(conversations.count) conversations")

        isLoading = false
    }

    func refresh() async {
        await syncFriendsToConversations()
        await fetchNewMessages()
        conversations = db.getAllConversations()
    }

    private func fetchNewMessages() async {
        do {
            let serverMessages = try await api.getMessages()
            print("Fetched \(serverMessages.count) messages from server")

            for message in serverMessages {
                await processMessage(message)
            }
        } catch {
            // Don't show error for empty messages or 404
            if case APIError.httpError(let code, _) = error, code == 404 {
                return
            }
            if case APIError.notFound = error {
                return
            }
            print("Failed to fetch messages: \(error)")
        }
    }

    private func syncFriendsToConversations() async {
        do {
            let friends = try await api.getFriends()
            print("Syncing \(friends.count) friends to conversations")

            for friend in friends {
                let conv = db.getOrCreateConversation(for: friend)
                print("Created/found conversation for \(conv.friendUsername)")
            }
        } catch {
            print("Failed to sync friends: \(error)")
            errorMessage = "Failed to load friends"
        }
    }

    private func processMessage(_ message: Message) async {
        // Skip if we already have this message
        guard !db.messageExists(id: message.id) else {
            print("Message \(message.id) already exists, skipping")
            return
        }

        do {
            let senderPublicKey = try crypto.publicKeyFromBase64(message.fromPublicKey)
            let privateKey = try crypto.getPrivateKey()

            let decryptedData = try crypto.decrypt(
                encryptedData: message.encryptedContent,
                signature: message.signature,
                senderPublicKey: senderPublicKey,
                recipientPrivateKey: privateKey
            )

            // Create stored message
            let storedMessage = StoredMessage(
                id: message.id,
                conversationID: message.fromUserID,
                contentType: message.contentType,
                decryptedContent: decryptedData,
                isFromMe: false,
                hasBeenViewed: false,
                serverDeleted: false,
                createdAt: message.createdAt,
                receivedAt: Date()
            )

            db.saveMessage(storedMessage)
            db.incrementUnreadCount(friendUserID: message.fromUserID)
            print("Saved message \(message.id) from \(message.fromUsername)")

            // Note: We DON'T acknowledge to server here anymore
            // We wait until the message is actually viewed
        } catch {
            print("Failed to process message \(message.id): \(error)")
        }
    }

    func markConversationRead(_ conversation: Conversation) {
        db.resetUnreadCount(friendUserID: conversation.friendUserID)
        if let index = conversations.firstIndex(where: { $0.id == conversation.id }) {
            conversations[index].unreadCount = 0
        }
    }
}

// Make Conversation conform to Hashable for navigation
extension Conversation: Hashable {
    static func == (lhs: Conversation, rhs: Conversation) -> Bool {
        lhs.friendUserID == rhs.friendUserID
    }

    func hash(into hasher: inout Hasher) {
        hasher.combine(friendUserID)
    }
}

#Preview {
    InboxView()
}
