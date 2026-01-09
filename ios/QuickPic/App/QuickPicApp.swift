//
//  QuickPicApp.swift
//  QuickPic
//

import SwiftUI

@main
struct QuickPicApp: App {
    @StateObject private var authManager = AuthManager()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(authManager)
                .preferredColorScheme(.dark)
        }
    }
}

struct ContentView: View {
    @EnvironmentObject var authManager: AuthManager

    var body: some View {
        Group {
            if authManager.isLoading {
                LaunchView()
            } else if authManager.isAuthenticated {
                MainView()
                    .onAppear {
                        SyncService.shared.startSyncTimer()
                    }
            } else {
                LoginView()
                    .onAppear {
                        SyncService.shared.stopSyncTimer()
                    }
            }
        }
        .background(Color.appBackground)
    }
}

struct LaunchView: View {
    var body: some View {
        ZStack {
            Color.appBackground.ignoresSafeArea()

            VStack(spacing: 16) {
                Image(systemName: "camera.fill")
                    .font(.system(size: 60))
                    .foregroundColor(.appPrimary)

                Text("QuickPic")
                    .font(.appTitle)
                    .foregroundColor(.textPrimary)

                ProgressView()
                    .tint(.appPrimary)
                    .padding(.top, 32)
            }
        }
    }
}

// MARK: - Main View (Unified Feed)

struct MainView: View {
    @EnvironmentObject var authManager: AuthManager
    @StateObject private var viewModel = MainViewModel()
    @State private var showCamera = false
    @State private var showAddFriend = false
    @State private var showProfile = false
    @State private var selectedConversation: Conversation?

    var body: some View {
        NavigationStack {
            ZStack {
                Color.appBackground.ignoresSafeArea()

                VStack(spacing: 0) {
                    // Top bar
                    topBar
                        .padding(.horizontal, AppSpacing.md)
                        .padding(.top, AppSpacing.sm)

                    // Feed content
                    if viewModel.isLoading && viewModel.feedItems.isEmpty {
                        loadingView
                    } else if viewModel.feedItems.isEmpty {
                        emptyView
                    } else {
                        feedList
                    }
                }

                // Floating camera button
                VStack {
                    Spacer()
                    HStack {
                        Spacer()
                        FloatingActionButton(icon: "camera.fill") {
                            Haptics.medium()
                            showCamera = true
                        }
                        .padding(.trailing, AppSpacing.lg)
                        .padding(.bottom, AppSpacing.lg)
                    }
                }
            }
            .navigationDestination(item: $selectedConversation) { conversation in
                ChatView(conversation: conversation, onMessagesViewed: {
                    viewModel.markConversationRead(conversation)
                })
            }
            .fullScreenCover(isPresented: $showCamera) {
                CameraView()
            }
            .sheet(isPresented: $showAddFriend) {
                AddFriendSheet(onRequestSent: {
                    Task { await viewModel.refresh() }
                })
                .presentationDetents([.medium])
                .presentationDragIndicator(.visible)
            }
            .sheet(isPresented: $showProfile) {
                ProfileSheet()
                    .presentationDetents([.medium, .large])
                    .presentationDragIndicator(.visible)
            }
            .refreshable {
                await viewModel.refresh()
            }
            .task {
                await viewModel.loadData()
            }
        }
    }

    // MARK: - Top Bar

    private var topBar: some View {
        HStack {
            // Profile button
            Button(action: {
                Haptics.light()
                showProfile = true
            }) {
                ZStack {
                    Circle()
                        .fill(Color.cardBackground)
                        .frame(width: 40, height: 40)

                    Text(authManager.currentUser?.username.prefix(1).uppercased() ?? "?")
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundColor(.textPrimary)
                }
            }
            .buttonStyle(ScaleButtonStyle())

            Spacer()

            // Add friend button
            IconButton(icon: "plus") {
                Haptics.light()
                showAddFriend = true
            }
        }
        .padding(.bottom, AppSpacing.md)
    }

    // MARK: - Feed List

    private var feedList: some View {
        ScrollView {
            LazyVStack(spacing: AppSpacing.sm) {
                ForEach(viewModel.feedItems) { item in
                    FeedRow(item: item) {
                        Haptics.light()
                        if case .request(let request) = item.type {
                            viewModel.acceptRequest(request)
                        }
                    } onReject: {
                        Haptics.light()
                        if case .request(let request) = item.type {
                            viewModel.rejectRequest(request)
                        }
                    } onTap: {
                        Haptics.light()
                        if case .conversation(let conversation) = item.type {
                            selectedConversation = conversation
                        }
                    }
                    .padding(.horizontal, AppSpacing.md)
                }
            }
            .padding(.top, AppSpacing.sm)
            .padding(.bottom, 100) // Space for FAB
        }
    }

    // MARK: - Empty State

    private var emptyView: some View {
        VStack(spacing: AppSpacing.md) {
            Spacer()

            Image(systemName: "person.2")
                .font(.system(size: 50))
                .foregroundColor(.textSecondary)

            Text("No friends yet")
                .font(.appHeadline)
                .foregroundColor(.textPrimary)

            Text("Tap + to add friends")
                .font(.appCaption)
                .foregroundColor(.textSecondary)

            Spacer()
        }
    }

    // MARK: - Loading State

    private var loadingView: some View {
        VStack {
            Spacer()
            ProgressView()
                .tint(.appPrimary)
            Spacer()
        }
    }
}

// MARK: - Feed Row

struct FeedRow: View {
    let item: FeedItem
    let onAccept: () -> Void
    let onReject: () -> Void
    let onTap: () -> Void

    var body: some View {
        Button(action: onTap) {
            HStack(spacing: AppSpacing.md) {
                // Status dot with initial
                StatusDot(
                    status: item.dotStatus,
                    initial: String(item.username.prefix(1)).uppercased()
                )

                // Content
                VStack(alignment: .leading, spacing: AppSpacing.xs) {
                    Text(item.username)
                        .font(.appHeadline)
                        .foregroundColor(.textPrimary)

                    Text(item.subtitle)
                        .font(.appCaption)
                        .foregroundColor(.textSecondary)
                }

                Spacer()

                // Right content (badge or actions)
                rightContent
            }
            .padding(AppSpacing.md)
            .cardStyle()
        }
        .buttonStyle(FeedRowButtonStyle())
    }

    @ViewBuilder
    private var rightContent: some View {
        switch item.type {
        case .request:
            HStack(spacing: AppSpacing.sm) {
                Button(action: onReject) {
                    Image(systemName: "xmark")
                        .font(.system(size: 14, weight: .semibold))
                        .foregroundColor(.danger)
                        .frame(width: 32, height: 32)
                        .background(Color.danger.opacity(0.15))
                        .clipShape(Circle())
                }

                Button(action: onAccept) {
                    Image(systemName: "checkmark")
                        .font(.system(size: 14, weight: .semibold))
                        .foregroundColor(.success)
                        .frame(width: 32, height: 32)
                        .background(Color.success.opacity(0.15))
                        .clipShape(Circle())
                }
            }

        case .conversation(let conversation):
            if conversation.unreadCount > 0 {
                Text("\(conversation.unreadCount)")
                    .font(.appSmall)
                    .foregroundColor(.white)
                    .padding(.horizontal, 8)
                    .padding(.vertical, 4)
                    .background(Color.appPrimary)
                    .clipShape(Capsule())
            }
        }
    }
}

struct FeedRowButtonStyle: ButtonStyle {
    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .scaleEffect(configuration.isPressed ? 0.98 : 1.0)
            .opacity(configuration.isPressed ? 0.9 : 1.0)
            .animation(.easeInOut(duration: 0.15), value: configuration.isPressed)
    }
}

// MARK: - Feed Item Model

struct FeedItem: Identifiable {
    let id: String
    let type: FeedItemType
    let username: String
    let subtitle: String
    let sortDate: Date

    var dotStatus: StatusDot.Status {
        switch type {
        case .request:
            return .pending
        case .conversation(let conv):
            return conv.unreadCount > 0 ? .unread : .read
        }
    }
}

enum FeedItemType {
    case request(FriendRequest)
    case conversation(Conversation)
}

// MARK: - Main View Model

@MainActor
class MainViewModel: ObservableObject {
    @Published var feedItems: [FeedItem] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private var friends: [Friend] = []
    private var pendingRequests: [FriendRequest] = []
    private var conversations: [Conversation] = []

    private let api = APIService.shared
    private let db = DatabaseService.shared
    private let crypto = CryptoService.shared

    private var lastLoadTime: Date?
    private let minRefreshInterval: TimeInterval = 3

    func loadData() async {
        await refresh()
    }

    func refresh() async {
        if let lastLoad = lastLoadTime,
           Date().timeIntervalSince(lastLoad) < minRefreshInterval {
            return
        }

        isLoading = true
        defer {
            isLoading = false
            lastLoadTime = Date()
        }

        // Fetch all data in parallel
        do {
            async let friendsTask = api.getFriends()
            async let requestsTask = api.getPendingFriendRequests()

            let (loadedFriends, loadedRequests) = try await (friendsTask, requestsTask)
            friends = loadedFriends
            pendingRequests = loadedRequests

            // Sync friends to conversations
            for friend in friends {
                _ = db.getOrCreateConversation(for: friend)
            }

            // Fetch new messages
            await fetchNewMessages()

            // Load conversations
            conversations = db.getAllConversations()

            // Build feed items
            buildFeedItems()
        } catch {
            if feedItems.isEmpty {
                errorMessage = "Failed to load data"
            }
        }
    }

    private func fetchNewMessages() async {
        do {
            let serverMessages = try await api.getMessages()

            for message in serverMessages {
                await processMessage(message)
            }
        } catch {
            // Silent fail for messages
        }
    }

    private func processMessage(_ message: Message) async {
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
            db.incrementUnreadCount(friendUserID: message.fromUserID)
        } catch {
            print("Failed to process message: \(error)")
        }
    }

    private func buildFeedItems() {
        var items: [FeedItem] = []

        // Add pending requests
        for request in pendingRequests {
            items.append(FeedItem(
                id: "request-\(request.id)",
                type: .request(request),
                username: request.fromUser.username,
                subtitle: "wants to add you",
                sortDate: request.createdAt
            ))
        }

        // Add conversations
        for conversation in conversations {
            let subtitle: String
            if conversation.unreadCount > 0 {
                subtitle = conversation.unreadCount == 1 ? "1 new message" : "\(conversation.unreadCount) new messages"
            } else if let lastMessage = conversation.lastMessageAt {
                subtitle = timeAgo(lastMessage)
            } else {
                subtitle = "Start chatting"
            }

            items.append(FeedItem(
                id: "conv-\(conversation.friendUserID)",
                type: .conversation(conversation),
                username: conversation.friendUsername,
                subtitle: subtitle,
                sortDate: conversation.lastMessageAt ?? conversation.createdAt
            ))
        }

        // Sort: requests first, then by date (newest first)
        items.sort { a, b in
            if case .request = a.type, case .conversation = b.type {
                return true
            }
            if case .conversation = a.type, case .request = b.type {
                return false
            }
            return a.sortDate > b.sortDate
        }

        feedItems = items
    }

    private func timeAgo(_ date: Date) -> String {
        let formatter = RelativeDateTimeFormatter()
        formatter.unitsStyle = .abbreviated
        return formatter.localizedString(for: date, relativeTo: Date())
    }

    func acceptRequest(_ request: FriendRequest) {
        Task {
            do {
                try await api.acceptFriendRequest(requestID: request.id)
                pendingRequests.removeAll { $0.id == request.id }
                buildFeedItems()
                await refresh()
                Haptics.success()
            } catch {
                Haptics.error()
            }
        }
    }

    func rejectRequest(_ request: FriendRequest) {
        Task {
            do {
                try await api.rejectFriendRequest(requestID: request.id)
                pendingRequests.removeAll { $0.id == request.id }
                buildFeedItems()
                Haptics.success()
            } catch {
                Haptics.error()
            }
        }
    }

    func markConversationRead(_ conversation: Conversation) {
        db.resetUnreadCount(friendUserID: conversation.friendUserID)
        if let index = conversations.firstIndex(where: { $0.id == conversation.id }) {
            conversations[index].unreadCount = 0
        }
        buildFeedItems()
    }
}

#Preview {
    ContentView()
        .environmentObject(AuthManager())
}
