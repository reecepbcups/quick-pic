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
            } else {
                LoginView()
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
    @State private var friendInfoConversation: Conversation?

    private let refreshTimer = Timer.publish(every: 5, on: .main, in: .common).autoconnect()

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
            .sheet(item: $friendInfoConversation) { conversation in
                FriendInfoSheet(conversation: conversation)
                    .presentationDetents([.medium])
                    .presentationDragIndicator(.visible)
            }
            .task {
                await viewModel.loadData()
            }
            .onReceive(refreshTimer) { _ in
                Task {
                    await viewModel.refresh()
                }
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
                    } onLongPress: {
                        if case .conversation(let conversation) = item.type {
                            friendInfoConversation = conversation
                        }
                    }
                    .padding(.horizontal, AppSpacing.md)
                }
            }
            .padding(.top, AppSpacing.sm)
            .padding(.bottom, 100) // Space for FAB
        }
        .refreshable {
            await viewModel.refresh()
        }
    }

    // MARK: - Empty State

    private var emptyView: some View {
        ScrollView {
            VStack(spacing: AppSpacing.md) {
                Spacer()
                    .frame(height: 150)

                Image(systemName: "person.2")
                    .font(.system(size: 50))
                    .foregroundColor(.textSecondary)

                Text("No friends yet")
                    .font(.appHeadline)
                    .foregroundColor(.textPrimary)

                Text("Tap + to add friends\nPull down to refresh")
                    .font(.appCaption)
                    .foregroundColor(.textSecondary)
                    .multilineTextAlignment(.center)

                Spacer()
                    .frame(height: 150)
            }
            .frame(maxWidth: .infinity)
        }
        .refreshable {
            await viewModel.refresh()
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
    let onLongPress: () -> Void

    var body: some View {
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
        .contentShape(Rectangle())
        .onTapGesture {
            onTap()
        }
        .onLongPressGesture {
            Haptics.medium()
            onLongPress()
        }
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

// MARK: - Friend Info Sheet

struct FriendInfoSheet: View {
    let conversation: Conversation
    @Environment(\.dismiss) private var dismiss
    @State private var userNumber: Int64?
    @State private var isLoading = true

    private var joinDateFormatted: String {
        let formatter = DateFormatter()
        formatter.dateStyle = .medium
        return formatter.string(from: conversation.friendSince)
    }

    var body: some View {
        NavigationStack {
            VStack(spacing: AppSpacing.lg) {
                // Friend avatar and name
                VStack(spacing: AppSpacing.sm) {
                    ZStack {
                        Circle()
                            .fill(Color.appPrimary.opacity(0.2))
                            .frame(width: 70, height: 70)

                        Text(String(conversation.friendUsername.prefix(1)).uppercased())
                            .font(.system(size: 28, weight: .bold))
                            .foregroundColor(.appPrimary)
                    }

                    Text(conversation.friendUsername)
                        .font(.appTitle)
                        .foregroundColor(.textPrimary)
                }
                .padding(.top, AppSpacing.md)

                // Info rows
                VStack(spacing: AppSpacing.sm) {
                    // User number
                    HStack {
                        Text("User #")
                            .font(.appCaption)
                            .foregroundColor(.textSecondary)
                        Spacer()
                        if isLoading {
                            ProgressView()
                                .tint(.appPrimary)
                        } else if let number = userNumber {
                            Text("\(number)")
                                .font(.appBody)
                                .foregroundColor(.textPrimary)
                        }
                    }
                    .padding(AppSpacing.md)
                    .background(Color.cardBackground)
                    .cornerRadius(AppRadius.md)

                    // Friends since
                    HStack {
                        Text("Friends Since")
                            .font(.appCaption)
                            .foregroundColor(.textSecondary)
                        Spacer()
                        Text(joinDateFormatted)
                            .font(.appBody)
                            .foregroundColor(.textPrimary)
                    }
                    .padding(AppSpacing.md)
                    .background(Color.cardBackground)
                    .cornerRadius(AppRadius.md)

                    // Public Key
                    VStack(alignment: .leading, spacing: AppSpacing.xs) {
                        Text("Public Key")
                            .font(.appCaption)
                            .foregroundColor(.textSecondary)

                        Text(conversation.friendPublicKey)
                            .font(.system(size: 10, design: .monospaced))
                            .foregroundColor(.textPrimary)
                            .lineLimit(2)
                            .textSelection(.enabled)
                    }
                    .padding(AppSpacing.md)
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .background(Color.cardBackground)
                    .cornerRadius(AppRadius.md)
                }
                .padding(.horizontal, AppSpacing.lg)

                Spacer()
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .background(Color.appBackground)
            .navigationTitle("Friend Info")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Done") {
                        dismiss()
                    }
                    .foregroundColor(.appPrimary)
                }
            }
            .task {
                do {
                    let user = try await APIService.shared.getUser(username: conversation.friendUsername)
                    userNumber = user.userNumber
                } catch {
                    print("Failed to fetch user info: \(error)")
                }
                isLoading = false
            }
        }
    }
}

struct InfoRow: View {
    let title: String
    let value: String

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            Text(title)
                .font(.appCaption)
                .foregroundColor(.textSecondary)

            Text(value)
                .font(.appBody)
                .foregroundColor(.textPrimary)
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(AppSpacing.md)
                .background(Color.cardBackground)
                .cornerRadius(AppRadius.md)
        }
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

    func refresh(force: Bool = false) async {
        if !force, let lastLoad = lastLoadTime,
           Date().timeIntervalSince(lastLoad) < minRefreshInterval {
            print("[MainViewModel] Skipping refresh - too soon")
            return
        }

        print("[MainViewModel] Starting refresh\(force ? " (forced)" : "")...")
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

            print("[MainViewModel] Loaded \(friends.count) friends, \(pendingRequests.count) pending requests")
            for req in pendingRequests {
                print("[MainViewModel] Pending request from: \(req.fromUser.username)")
            }

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
            print("[MainViewModel] Built \(feedItems.count) feed items")
        } catch {
            print("[MainViewModel] Refresh failed: \(error)")
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
                encryptedContent: message.encryptedContent,
                isFromMe: false,
                hasBeenViewed: false,
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
                Haptics.success()

                // Remove the request from UI immediately
                pendingRequests.removeAll { $0.id == request.id }
                buildFeedItems()

                // Force refresh to get the new friend
                await refresh(force: true)
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
