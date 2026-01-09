//
//  ChatView.swift
//  QuickPic
//
//  Individual conversation view showing message history with a friend
//

import SwiftUI

struct ChatView: View {
    let conversation: Conversation
    let onMessagesViewed: () -> Void

    @StateObject private var viewModel: ChatViewModel
    @State private var messageText = ""
    @State private var selectedMessage: StoredMessage?
    @FocusState private var isTextFieldFocused: Bool

    init(conversation: Conversation, onMessagesViewed: @escaping () -> Void) {
        self.conversation = conversation
        self.onMessagesViewed = onMessagesViewed
        _viewModel = StateObject(wrappedValue: ChatViewModel(conversation: conversation))
    }

    var body: some View {
        VStack(spacing: 0) {
            // Messages list
            ScrollViewReader { proxy in
                ScrollView {
                    LazyVStack(spacing: 8) {
                        ForEach(viewModel.messages) { message in
                            MessageBubble(message: message)
                                .id(message.id)
                                .onTapGesture {
                                    if !message.isFromMe && !message.hasBeenViewed {
                                        selectedMessage = message
                                    }
                                }
                        }
                    }
                    .padding()
                }
                .onChange(of: viewModel.messages.count) { _, _ in
                    if let lastMessage = viewModel.messages.last {
                        withAnimation {
                            proxy.scrollTo(lastMessage.id, anchor: .bottom)
                        }
                    }
                }
            }

            Divider()

            // Message input
            HStack(spacing: 12) {
                TextField("Message", text: $messageText)
                    .textFieldStyle(.roundedBorder)
                    .focused($isTextFieldFocused)

                Button(action: sendMessage) {
                    Image(systemName: "paperplane.fill")
                        .foregroundColor(canSend ? .yellow : .gray)
                }
                .disabled(!canSend)
            }
            .padding()
            .background(Color(.systemBackground))
        }
        .navigationTitle(conversation.friendUsername)
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            ToolbarItem(placement: .navigationBarTrailing) {
                Button(action: { Task { await viewModel.refresh() } }) {
                    Image(systemName: "arrow.clockwise")
                }
            }
        }
        .fullScreenCover(item: $selectedMessage) { message in
            MessageContentView(message: message) {
                Task {
                    await viewModel.markAsViewed(message)
                }
                selectedMessage = nil
            }
        }
        .task {
            await viewModel.loadMessages()
            onMessagesViewed()
        }
    }

    private var canSend: Bool {
        !messageText.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty && !viewModel.isSending
    }

    private func sendMessage() {
        guard canSend else { return }

        let text = messageText
        messageText = ""
        isTextFieldFocused = false

        Task {
            await viewModel.sendTextMessage(text)
        }
    }
}

struct MessageBubble: View {
    let message: StoredMessage

    var body: some View {
        HStack {
            if message.isFromMe {
                Spacer()
            }

            VStack(alignment: message.isFromMe ? .trailing : .leading, spacing: 4) {
                if message.contentType == .image {
                    // Image message - show preview or tap indicator
                    if message.isFromMe {
                        // Sent image - show thumbnail
                        if let uiImage = UIImage(data: message.decryptedContent) {
                            Image(uiImage: uiImage)
                                .resizable()
                                .scaledToFit()
                                .frame(maxWidth: 200, maxHeight: 200)
                                .cornerRadius(12)
                        }
                    } else {
                        // Received image - show tap to view
                        HStack {
                            Image(systemName: message.hasBeenViewed ? "photo" : "photo.fill")
                            Text(message.hasBeenViewed ? "Viewed" : "Tap to view")
                        }
                        .padding(.horizontal, 16)
                        .padding(.vertical, 12)
                        .background(message.hasBeenViewed ? Color.gray.opacity(0.3) : Color.yellow.opacity(0.3))
                        .cornerRadius(16)
                    }
                } else {
                    // Text message
                    if message.isFromMe {
                        if let text = String(data: message.decryptedContent, encoding: .utf8) {
                            Text(text)
                                .padding(.horizontal, 16)
                                .padding(.vertical, 10)
                                .background(Color.yellow)
                                .foregroundColor(.black)
                                .cornerRadius(16)
                        }
                    } else {
                        // Received text - tap to view if not viewed
                        if message.hasBeenViewed {
                            if let text = String(data: message.decryptedContent, encoding: .utf8) {
                                Text(text)
                                    .padding(.horizontal, 16)
                                    .padding(.vertical, 10)
                                    .background(Color.gray.opacity(0.3))
                                    .cornerRadius(16)
                            }
                        } else {
                            HStack {
                                Image(systemName: "text.bubble.fill")
                                Text("Tap to view")
                            }
                            .padding(.horizontal, 16)
                            .padding(.vertical, 12)
                            .background(Color.yellow.opacity(0.3))
                            .cornerRadius(16)
                        }
                    }
                }

                Text(timeAgo(message.createdAt))
                    .font(.caption2)
                    .foregroundColor(.secondary)
            }

            if !message.isFromMe {
                Spacer()
            }
        }
    }

    private func timeAgo(_ date: Date) -> String {
        let formatter = RelativeDateTimeFormatter()
        formatter.unitsStyle = .abbreviated
        return formatter.localizedString(for: date, relativeTo: Date())
    }
}

struct MessageContentView: View {
    let message: StoredMessage
    let onDismiss: () -> Void

    @State private var showContent = false
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        ZStack {
            Color.black.ignoresSafeArea()

            if showContent {
                contentView
            } else {
                ProgressView()
                    .tint(.white)
            }
        }
        .onAppear {
            DispatchQueue.main.asyncAfter(deadline: .now() + 0.3) {
                withAnimation {
                    showContent = true
                }
            }
        }
        .gesture(
            DragGesture(minimumDistance: 0)
                .onEnded { _ in
                    onDismiss()
                    dismiss()
                }
        )
    }

    @ViewBuilder
    private var contentView: some View {
        VStack {
            Spacer()

            if message.contentType == .image {
                if let uiImage = UIImage(data: message.decryptedContent) {
                    Image(uiImage: uiImage)
                        .resizable()
                        .scaledToFit()
                        .cornerRadius(12)
                        .padding()
                }
            } else {
                if let text = String(data: message.decryptedContent, encoding: .utf8) {
                    Text(text)
                        .font(.title2)
                        .foregroundColor(.white)
                        .multilineTextAlignment(.center)
                        .padding(32)
                }
            }

            Spacer()

            Text("Release to close")
                .font(.caption)
                .foregroundColor(.gray)
                .padding(.bottom, 32)
        }
    }
}

@MainActor
class ChatViewModel: ObservableObject {
    @Published var messages: [StoredMessage] = []
    @Published var isSending = false

    private let conversation: Conversation
    private let db = DatabaseService.shared
    private let api = APIService.shared
    private let crypto = CryptoService.shared

    init(conversation: Conversation) {
        self.conversation = conversation
    }

    func loadMessages() async {
        print("ChatView loading messages for conversation: \(conversation.friendUserID)")
        messages = db.getMessages(for: conversation.friendUserID)
        print("Loaded \(messages.count) messages from database")
        await refresh()
    }

    func refresh() async {
        // Fetch new messages from server
        do {
            let serverMessages = try await api.getMessages()
            print("ChatView fetched \(serverMessages.count) total messages from server")

            let relevantMessages = serverMessages.filter { $0.fromUserID == conversation.friendUserID }
            print("Found \(relevantMessages.count) messages for this conversation")

            for message in relevantMessages {
                await processIncomingMessage(message)
            }

            messages = db.getMessages(for: conversation.friendUserID)
            print("After refresh, have \(messages.count) messages")
        } catch {
            print("ChatView refresh error: \(error)")
        }
    }

    private func processIncomingMessage(_ message: Message) async {
        guard !db.messageExists(id: message.id) else { return }

        do {
            let senderPublicKey = try crypto.publicKeyFromBase64(message.fromPublicKey)
            let privateKey = try crypto.getPrivateKey()

            let decryptedData = try crypto.decrypt(
                encryptedData: message.encryptedContent,
                signature: message.signature,
                senderPublicKey: senderPublicKey,
                recipientPrivateKey: privateKey
            )

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
        } catch {
            print("Failed to process message: \(error)")
        }
    }

    func markAsViewed(_ message: StoredMessage) async {
        db.markMessageAsViewed(messageID: message.id)

        // Delete from server now that it's been viewed
        do {
            try await api.acknowledgeMessage(id: message.id)
            db.markMessageServerDeleted(messageID: message.id)
        } catch {
            // Will retry later
        }

        // Update local state
        if let index = messages.firstIndex(where: { $0.id == message.id }) {
            messages[index].hasBeenViewed = true
        }
    }

    func sendTextMessage(_ text: String) async {
        guard let textData = text.data(using: .utf8) else { return }

        isSending = true
        defer { isSending = false }

        do {
            let privateKey = try crypto.getPrivateKey()
            let recipientPublicKey = try crypto.publicKeyFromBase64(conversation.friendPublicKey)

            let (encryptedData, signature) = try crypto.encrypt(
                content: textData,
                recipientPublicKey: recipientPublicKey,
                senderPrivateKey: privateKey
            )

            let response = try await api.sendMessage(
                to: conversation.friendUsername,
                encryptedContent: encryptedData,
                contentType: .text,
                signature: signature
            )

            // Save sent message locally
            let storedMessage = StoredMessage(
                id: response.id,
                conversationID: conversation.friendUserID,
                contentType: .text,
                decryptedContent: textData,
                isFromMe: true,
                hasBeenViewed: true,
                serverDeleted: true,
                createdAt: response.createdAt,
                receivedAt: Date()
            )

            db.saveMessage(storedMessage)
            messages = db.getMessages(for: conversation.friendUserID)
        } catch {
            print("Failed to send message: \(error)")
        }
    }

    func sendImage(_ imageData: Data) async {
        isSending = true
        defer { isSending = false }

        do {
            let privateKey = try crypto.getPrivateKey()
            let recipientPublicKey = try crypto.publicKeyFromBase64(conversation.friendPublicKey)

            let (encryptedData, signature) = try crypto.encrypt(
                content: imageData,
                recipientPublicKey: recipientPublicKey,
                senderPrivateKey: privateKey
            )

            let response = try await api.sendMessage(
                to: conversation.friendUsername,
                encryptedContent: encryptedData,
                contentType: .image,
                signature: signature
            )

            // Save sent message locally
            let storedMessage = StoredMessage(
                id: response.id,
                conversationID: conversation.friendUserID,
                contentType: .image,
                decryptedContent: imageData,
                isFromMe: true,
                hasBeenViewed: true,
                serverDeleted: true,
                createdAt: response.createdAt,
                receivedAt: Date()
            )

            db.saveMessage(storedMessage)
            messages = db.getMessages(for: conversation.friendUserID)
        } catch {
            print("Failed to send image: \(error)")
        }
    }
}

// Make StoredMessage conform to Identifiable for ForEach
extension StoredMessage: Hashable {
    static func == (lhs: StoredMessage, rhs: StoredMessage) -> Bool {
        lhs.id == rhs.id
    }

    func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }
}

#Preview {
    NavigationStack {
        ChatView(
            conversation: Conversation(
                friendUserID: UUID(),
                friendUsername: "testuser",
                friendPublicKey: "",
                lastMessageAt: Date(),
                unreadCount: 2,
                createdAt: Date()
            ),
            onMessagesViewed: {}
        )
    }
}
