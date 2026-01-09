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
    @Environment(\.dismiss) private var dismiss

    init(conversation: Conversation, onMessagesViewed: @escaping () -> Void) {
        self.conversation = conversation
        self.onMessagesViewed = onMessagesViewed
        _viewModel = StateObject(wrappedValue: ChatViewModel(conversation: conversation))
    }

    var body: some View {
        ZStack {
            Color.appBackground.ignoresSafeArea()

            VStack(spacing: 0) {
                // Messages list
                ScrollViewReader { proxy in
                    ScrollView {
                        LazyVStack(spacing: AppSpacing.sm) {
                            ForEach(viewModel.messages) { message in
                                MessageBubble(message: message)
                                    .id(message.id)
                                    .onTapGesture {
                                        if !message.isFromMe && !message.hasBeenViewed {
                                            Haptics.light()
                                            selectedMessage = message
                                        }
                                    }
                            }
                        }
                        .padding(AppSpacing.md)
                    }
                    .onChange(of: viewModel.messages.count) { _, _ in
                        if let lastMessage = viewModel.messages.last {
                            withAnimation(.spring(response: 0.3)) {
                                proxy.scrollTo(lastMessage.id, anchor: .bottom)
                            }
                        }
                    }
                }

                // Message input
                messageInputBar
            }
        }
        .navigationTitle(conversation.friendUsername)
        .navigationBarTitleDisplayMode(.inline)
        .toolbarBackground(Color.appBackground, for: .navigationBar)
        .toolbarBackground(.visible, for: .navigationBar)
        .toolbar {
            ToolbarItem(placement: .navigationBarTrailing) {
                Button(action: { Task { await viewModel.refresh() } }) {
                    Image(systemName: "arrow.clockwise")
                        .foregroundColor(.appPrimary)
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

    private var messageInputBar: some View {
        HStack(spacing: AppSpacing.sm) {
            HStack {
                TextField("Message", text: $messageText)
                    .foregroundColor(.textPrimary)
                    .focused($isTextFieldFocused)
            }
            .padding(.horizontal, AppSpacing.md)
            .padding(.vertical, 12)
            .background(Color.cardBackground)
            .cornerRadius(AppRadius.xl)

            Button(action: sendMessage) {
                Image(systemName: "paperplane.fill")
                    .font(.system(size: 18))
                    .foregroundColor(canSend ? .appPrimary : .textSecondary)
                    .frame(width: 44, height: 44)
                    .background(canSend ? Color.appPrimary.opacity(0.15) : Color.cardBackground)
                    .clipShape(Circle())
            }
            .disabled(!canSend)
            .buttonStyle(ScaleButtonStyle())
        }
        .padding(AppSpacing.md)
        .background(Color.appBackground)
    }

    private var canSend: Bool {
        !messageText.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty && !viewModel.isSending
    }

    private func sendMessage() {
        guard canSend else { return }
        Haptics.light()

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
                Spacer(minLength: 60)
            }

            VStack(alignment: message.isFromMe ? .trailing : .leading, spacing: AppSpacing.xs) {
                if message.contentType == .image {
                    imageContent
                } else {
                    textContent
                }

                Text(timeAgo(message.createdAt))
                    .font(.appSmall)
                    .foregroundColor(.textSecondary)
            }

            if !message.isFromMe {
                Spacer(minLength: 60)
            }
        }
    }

    @ViewBuilder
    private var imageContent: some View {
        if message.isFromMe {
            if let uiImage = UIImage(data: message.decryptedContent) {
                Image(uiImage: uiImage)
                    .resizable()
                    .scaledToFit()
                    .frame(maxWidth: 200, maxHeight: 200)
                    .cornerRadius(AppRadius.md)
            }
        } else {
            HStack(spacing: AppSpacing.sm) {
                Image(systemName: message.hasBeenViewed ? "photo" : "photo.fill")
                Text(message.hasBeenViewed ? "Viewed" : "Tap to view")
            }
            .font(.appCaption)
            .foregroundColor(message.hasBeenViewed ? .textSecondary : .textPrimary)
            .padding(.horizontal, AppSpacing.md)
            .padding(.vertical, AppSpacing.sm)
            .background(message.hasBeenViewed ? Color.cardBackground : Color.appPrimary.opacity(0.2))
            .cornerRadius(AppRadius.lg)
        }
    }

    @ViewBuilder
    private var textContent: some View {
        if message.isFromMe {
            if let text = String(data: message.decryptedContent, encoding: .utf8) {
                Text(text)
                    .font(.appBody)
                    .foregroundColor(.black)
                    .padding(.horizontal, AppSpacing.md)
                    .padding(.vertical, AppSpacing.sm)
                    .background(Color.appPrimary)
                    .cornerRadius(AppRadius.lg)
            }
        } else {
            if message.hasBeenViewed {
                if let text = String(data: message.decryptedContent, encoding: .utf8) {
                    Text(text)
                        .font(.appBody)
                        .foregroundColor(.textPrimary)
                        .padding(.horizontal, AppSpacing.md)
                        .padding(.vertical, AppSpacing.sm)
                        .background(Color.cardBackground)
                        .cornerRadius(AppRadius.lg)
                }
            } else {
                HStack(spacing: AppSpacing.sm) {
                    Image(systemName: "text.bubble.fill")
                    Text("Tap to view")
                }
                .font(.appCaption)
                .foregroundColor(.textPrimary)
                .padding(.horizontal, AppSpacing.md)
                .padding(.vertical, AppSpacing.sm)
                .background(Color.appPrimary.opacity(0.2))
                .cornerRadius(AppRadius.lg)
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
            Color.appBackground.ignoresSafeArea()

            if showContent {
                contentView
            } else {
                ProgressView()
                    .tint(.appPrimary)
            }
        }
        .onAppear {
            Haptics.medium()
            DispatchQueue.main.asyncAfter(deadline: .now() + 0.3) {
                withAnimation(.spring(response: 0.4)) {
                    showContent = true
                }
            }
        }
        .gesture(
            DragGesture(minimumDistance: 0)
                .onEnded { _ in
                    Haptics.light()
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
                        .cornerRadius(AppRadius.md)
                        .padding(AppSpacing.md)
                }
            } else {
                if let text = String(data: message.decryptedContent, encoding: .utf8) {
                    Text(text)
                        .font(.title2)
                        .foregroundColor(.textPrimary)
                        .multilineTextAlignment(.center)
                        .padding(AppSpacing.xl)
                }
            }

            Spacer()

            Text("Release to close")
                .font(.appCaption)
                .foregroundColor(.textSecondary)
                .padding(.bottom, AppSpacing.xl)
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
        messages = db.getMessages(for: conversation.friendUserID)
        await refresh()
    }

    func refresh() async {
        do {
            let serverMessages = try await api.getMessages()
            let relevantMessages = serverMessages.filter { $0.fromUserID == conversation.friendUserID }

            for message in relevantMessages {
                await processIncomingMessage(message)
            }

            messages = db.getMessages(for: conversation.friendUserID)
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

        do {
            try await api.acknowledgeMessage(id: message.id)
            db.markMessageServerDeleted(messageID: message.id)
        } catch {
            // Will retry later
        }

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
            Haptics.success()
        } catch {
            Haptics.error()
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
            Haptics.success()
        } catch {
            Haptics.error()
            print("Failed to send image: \(error)")
        }
    }
}

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
