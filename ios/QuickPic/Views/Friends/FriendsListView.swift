import SwiftUI

struct FriendsListView: View {
    @StateObject private var viewModel = FriendsViewModel()
    @State private var showAddFriend = false
    @State private var showInbox = false

    var body: some View {
        NavigationStack {
            List {
                // Pending Requests Section
                if !viewModel.pendingRequests.isEmpty {
                    Section("Friend Requests") {
                        ForEach(viewModel.pendingRequests) { request in
                            FriendRequestRow(
                                request: request,
                                onAccept: { viewModel.acceptRequest(request) },
                                onReject: { viewModel.rejectRequest(request) }
                            )
                        }
                    }
                }

                // Friends Section
                Section("Friends") {
                    if viewModel.friends.isEmpty && !viewModel.isLoading {
                        Text("No friends yet")
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(viewModel.friends) { friend in
                            NavigationLink(destination: friendChatDestination(for: friend)) {
                                FriendRow(friend: friend)
                            }
                        }
                    }
                }
            }
            .navigationTitle("Friends")
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button(action: { showInbox = true }) {
                        Image(systemName: "tray.fill")
                            .foregroundColor(.yellow)
                    }
                }
                ToolbarItem(placement: .navigationBarTrailing) {
                    HStack(spacing: 16) {
                        Button(action: { Task { await viewModel.refresh() } }) {
                            if viewModel.isLoading {
                                ProgressView()
                                    .scaleEffect(0.8)
                            } else {
                                Image(systemName: "arrow.clockwise")
                            }
                        }
                        .disabled(viewModel.isLoading)

                        Button(action: { showAddFriend = true }) {
                            Image(systemName: "person.badge.plus")
                        }
                    }
                }
            }
            .refreshable {
                await viewModel.refresh()
            }
            .sheet(isPresented: $showAddFriend) {
                AddFriendView(onRequestSent: {
                    Task { await viewModel.refresh() }
                })
            }
            .navigationDestination(isPresented: $showInbox) {
                InboxView()
            }
            .alert("Error", isPresented: $viewModel.showError) {
                Button("Retry") {
                    Task { await viewModel.refresh() }
                }
                Button("OK", role: .cancel) {}
            } message: {
                Text(viewModel.errorMessage ?? "Failed to load friends")
            }
            .task {
                await viewModel.loadData()
            }
            .onAppear {
                // Auto-refresh when view appears
                if viewModel.friends.isEmpty && !viewModel.isLoading {
                    Task { await viewModel.refresh() }
                }
            }
        }
    }

    @ViewBuilder
    private func friendChatDestination(for friend: Friend) -> some View {
        let conversation = Conversation(
            friendUserID: friend.userID,
            friendUsername: friend.username,
            friendPublicKey: friend.publicKey,
            lastMessageAt: nil,
            unreadCount: 0,
            createdAt: friend.since
        )
        ChatView(conversation: conversation, onMessagesViewed: {})
    }
}

struct FriendRow: View {
    let friend: Friend

    var body: some View {
        HStack {
            Circle()
                .fill(Color.yellow.opacity(0.3))
                .frame(width: 44, height: 44)
                .overlay(
                    Text(friend.username.prefix(1).uppercased())
                        .fontWeight(.semibold)
                )

            VStack(alignment: .leading) {
                Text(friend.username)
                    .fontWeight(.medium)

                Text("Friends since \(friend.since.formatted(date: .abbreviated, time: .omitted))")
                    .font(.caption)
                    .foregroundColor(.secondary)
            }

            Spacer()

            Image(systemName: "bubble.left.fill")
                .foregroundColor(.yellow.opacity(0.7))
                .font(.caption)
        }
    }
}

struct FriendRequestRow: View {
    let request: FriendRequest
    let onAccept: () -> Void
    let onReject: () -> Void

    var body: some View {
        HStack {
            Circle()
                .fill(Color.blue.opacity(0.3))
                .frame(width: 44, height: 44)
                .overlay(
                    Text(request.fromUser.username.prefix(1).uppercased())
                        .fontWeight(.semibold)
                )

            VStack(alignment: .leading) {
                Text(request.fromUser.username)
                    .fontWeight(.medium)

                Text("Wants to be friends")
                    .font(.caption)
                    .foregroundColor(.secondary)
            }

            Spacer()

            HStack(spacing: 8) {
                Button(action: onReject) {
                    Image(systemName: "xmark")
                        .foregroundColor(.red)
                        .padding(8)
                        .background(Color.red.opacity(0.1))
                        .clipShape(Circle())
                }

                Button(action: onAccept) {
                    Image(systemName: "checkmark")
                        .foregroundColor(.green)
                        .padding(8)
                        .background(Color.green.opacity(0.1))
                        .clipShape(Circle())
                }
            }
        }
    }
}

@MainActor
class FriendsViewModel: ObservableObject {
    @Published var friends: [Friend] = []
    @Published var pendingRequests: [FriendRequest] = []
    @Published var errorMessage: String?
    @Published var showError = false
    @Published var isLoading = false

    private var lastLoadTime: Date?
    private let minRefreshInterval: TimeInterval = 5 // Minimum 5 seconds between refreshes

    func loadData() async {
        await refresh()
    }

    func refresh() async {
        // Prevent too frequent refreshes
        if let lastLoad = lastLoadTime,
           Date().timeIntervalSince(lastLoad) < minRefreshInterval {
            return
        }

        isLoading = true
        defer {
            isLoading = false
            lastLoadTime = Date()
        }

        // Use a timeout for the API calls
        do {
            async let friendsTask = APIService.shared.getFriends()
            async let requestsTask = APIService.shared.getPendingFriendRequests()

            let (loadedFriends, loadedRequests) = try await (friendsTask, requestsTask)
            friends = loadedFriends
            pendingRequests = loadedRequests

            // Sync friends to conversations database
            let db = DatabaseService.shared
            for friend in loadedFriends {
                _ = db.getOrCreateConversation(for: friend)
            }
        } catch {
            // Only show error if we have no data
            if friends.isEmpty {
                errorMessage = getErrorMessage(for: error)
                showError = true
            }
        }
    }

    func acceptRequest(_ request: FriendRequest) {
        Task {
            do {
                try await APIService.shared.acceptFriendRequest(requestID: request.id)
                // Remove from local list immediately for responsiveness
                pendingRequests.removeAll { $0.id == request.id }
                await refresh()
            } catch {
                errorMessage = "Failed to accept request"
                showError = true
            }
        }
    }

    func rejectRequest(_ request: FriendRequest) {
        Task {
            do {
                try await APIService.shared.rejectFriendRequest(requestID: request.id)
                // Remove from local list immediately for responsiveness
                pendingRequests.removeAll { $0.id == request.id }
            } catch {
                errorMessage = "Failed to reject request"
                showError = true
            }
        }
    }

    private func getErrorMessage(for error: Error) -> String {
        if let apiError = error as? APIError {
            switch apiError {
            case .unauthorized:
                return "Session expired. Please log in again."
            case .requestFailed:
                return "Network error. Check your connection."
            default:
                return "Failed to load friends. Pull to retry."
            }
        }
        return "Failed to load friends. Pull to retry."
    }
}

#Preview {
    FriendsListView()
}
