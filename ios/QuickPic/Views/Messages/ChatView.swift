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
    @State private var debugMessage: StoredMessage?
    @State private var fullscreenImage: UIImage?
    @FocusState private var isTextFieldFocused: Bool
    @Environment(\.dismiss) private var dismiss

    private let refreshTimer = Timer.publish(every: 3, on: .main, in: .common).autoconnect()

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
                                MessageBubble(message: message, onImageTap: { image in
                                    fullscreenImage = image
                                })
                                    .id(message.id)
                                    .onLongPressGesture {
                                        Haptics.medium()
                                        debugMessage = message
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
        .task {
            await viewModel.loadMessages()
            onMessagesViewed()
        }
        .onReceive(refreshTimer) { _ in
            Task {
                await viewModel.refresh()
            }
        }
        .sheet(item: $debugMessage) { message in
            MessageDebugSheet(message: message)
                .presentationDetents([.medium, .large])
                .presentationDragIndicator(.visible)
        }
        .fullScreenCover(item: Binding(
            get: { fullscreenImage.map { IdentifiableImage(image: $0) } },
            set: { fullscreenImage = $0?.image }
        )) { item in
            FullscreenImageView(image: item.image)
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
    var onImageTap: ((UIImage) -> Void)?

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
        if let uiImage = UIImage(data: message.decryptedContent) {
            Image(uiImage: uiImage)
                .resizable()
                .scaledToFit()
                .frame(maxWidth: 200, maxHeight: 200)
                .cornerRadius(AppRadius.md)
                .onTapGesture {
                    onImageTap?(uiImage)
                }
        }
    }

    @ViewBuilder
    private var textContent: some View {
        if let text = String(data: message.decryptedContent, encoding: .utf8) {
            Text(text)
                .font(.appBody)
                .foregroundColor(message.isFromMe ? .black : .textPrimary)
                .padding(.horizontal, AppSpacing.md)
                .padding(.vertical, AppSpacing.sm)
                .background(message.isFromMe ? Color.appPrimary : Color.cardBackground)
                .cornerRadius(AppRadius.lg)
        }
    }

    private func timeAgo(_ date: Date) -> String {
        let formatter = RelativeDateTimeFormatter()
        formatter.unitsStyle = .abbreviated
        return formatter.localizedString(for: date, relativeTo: Date())
    }
}

// MARK: - Message Debug Sheet

struct MessageDebugSheet: View {
    let message: StoredMessage
    @Environment(\.dismiss) private var dismiss

    private var isoTimestamp: String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        return formatter.string(from: message.createdAt)
    }

    private var encryptedBytesHex: String {
        guard let encrypted = message.encryptedContent else {
            return "Not available"
        }
        let previewBytes = encrypted.prefix(256)
        let hex = previewBytes.map { String(format: "%02x", $0) }.joined(separator: " ")
        if encrypted.count > 256 {
            return hex + " ... (\(encrypted.count) bytes total)"
        }
        return hex
    }

    var body: some View {
        NavigationStack {
            ZStack {
                Color.appBackground.ignoresSafeArea()

                ScrollView {
                    VStack(alignment: .leading, spacing: AppSpacing.lg) {
                        // Message UID
                        DebugInfoRow(title: "Message UID", value: message.id.uuidString)

                        // Exact timestamp in ISO format
                        DebugInfoRow(title: "Timestamp (ISO 8601)", value: isoTimestamp)

                        // Content type
                        DebugInfoRow(title: "Content Type", value: message.contentType.rawValue)

                        // Direction
                        DebugInfoRow(title: "Direction", value: message.isFromMe ? "Sent" : "Received")

                        // Encrypted content
                        VStack(alignment: .leading, spacing: AppSpacing.sm) {
                            Text("Encrypted Content (Hex)")
                                .font(.appCaption)
                                .foregroundColor(.textSecondary)

                            ScrollView(.horizontal, showsIndicators: true) {
                                Text(encryptedBytesHex)
                                    .font(.system(size: 10, design: .monospaced))
                                    .foregroundColor(.textPrimary)
                                    .textSelection(.enabled)
                            }
                            .padding(AppSpacing.md)
                            .frame(maxWidth: .infinity, alignment: .leading)
                            .frame(maxHeight: 150)
                            .background(Color.cardBackground)
                            .cornerRadius(AppRadius.md)
                        }

                        // Encrypted content size
                        if let encrypted = message.encryptedContent {
                            DebugInfoRow(title: "Encrypted Size", value: "\(encrypted.count) bytes")
                        }

                        // Decrypted content size
                        DebugInfoRow(title: "Decrypted Size", value: "\(message.decryptedContent.count) bytes")

                        Spacer()
                    }
                    .padding(AppSpacing.lg)
                }
            }
            .navigationTitle("Message Debug Info")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Done") {
                        dismiss()
                    }
                    .foregroundColor(.appPrimary)
                }
            }
        }
    }
}

struct DebugInfoRow: View {
    let title: String
    let value: String

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            Text(title)
                .font(.appCaption)
                .foregroundColor(.textSecondary)

            Text(value)
                .font(.system(size: 12, design: .monospaced))
                .foregroundColor(.textPrimary)
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(AppSpacing.md)
                .background(Color.cardBackground)
                .cornerRadius(AppRadius.md)
                .textSelection(.enabled)
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
                encryptedContent: message.encryptedContent,
                isFromMe: false,
                hasBeenViewed: false,
                                createdAt: message.createdAt,
                receivedAt: Date()
            )

            db.saveMessage(storedMessage)
        } catch {
            print("Failed to process message: \(error)")
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
                encryptedContent: encryptedData,
                isFromMe: true,
                hasBeenViewed: true,
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
                encryptedContent: encryptedData,
                isFromMe: true,
                hasBeenViewed: true,
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

// MARK: - Fullscreen Image Viewer

struct IdentifiableImage: Identifiable {
    let id = UUID()
    let image: UIImage
}

struct FullscreenImageView: View {
    let image: UIImage
    @Environment(\.dismiss) private var dismiss
    @State private var offset: CGSize = .zero
    @State private var opacity: Double = 1.0
    @State private var scale: CGFloat = 1.0
    @State private var lastScale: CGFloat = 1.0

    var body: some View {
        GeometryReader { geometry in
            ZStack {
                Color.black.opacity(opacity)
                    .ignoresSafeArea()
                    .onTapGesture {
                        dismiss()
                    }

                Image(uiImage: image)
                    .resizable()
                    .scaledToFit()
                    .scaleEffect(scale)
                    .offset(offset)
                    .gesture(
                        MagnificationGesture()
                            .onChanged { value in
                                scale = lastScale * value
                            }
                            .onEnded { value in
                                let finalScale = lastScale * value
                                if finalScale < 0.7 {
                                    withAnimation(.easeOut(duration: 0.2)) {
                                        scale = 0.1
                                        opacity = 0
                                    }
                                    DispatchQueue.main.asyncAfter(deadline: .now() + 0.2) {
                                        dismiss()
                                    }
                                } else if finalScale < 1.0 {
                                    withAnimation(.spring(response: 0.3)) {
                                        scale = 1.0
                                        lastScale = 1.0
                                    }
                                } else {
                                    lastScale = scale
                                }
                            }
                    )
                    .simultaneousGesture(
                        DragGesture()
                            .onChanged { value in
                                if scale <= 1.0 {
                                    offset = value.translation
                                    let progress = min(abs(value.translation.height) / 300, 1.0)
                                    opacity = 1.0 - progress * 0.5
                                } else {
                                    offset = value.translation
                                }
                            }
                            .onEnded { value in
                                if scale <= 1.0 && abs(value.translation.height) > 100 {
                                    let direction: CGFloat = value.translation.height > 0 ? 1 : -1
                                    withAnimation(.easeOut(duration: 0.2)) {
                                        offset = CGSize(width: 0, height: direction * geometry.size.height)
                                        opacity = 0
                                    }
                                    DispatchQueue.main.asyncAfter(deadline: .now() + 0.2) {
                                        dismiss()
                                    }
                                } else if scale <= 1.0 {
                                    withAnimation(.spring(response: 0.3)) {
                                        offset = .zero
                                        opacity = 1.0
                                    }
                                }
                            }
                    )
                    .onTapGesture(count: 2) {
                        withAnimation(.spring(response: 0.3)) {
                            if scale > 1.0 {
                                scale = 1.0
                                lastScale = 1.0
                                offset = .zero
                            } else {
                                scale = 2.5
                                lastScale = 2.5
                            }
                        }
                    }
            }
        }
        .statusBarHidden()
    }
}
