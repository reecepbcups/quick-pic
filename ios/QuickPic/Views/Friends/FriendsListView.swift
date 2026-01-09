import SwiftUI

struct FriendsListView: View {
    @StateObject private var viewModel = FriendsViewModel()
    @State private var showAddFriend = false

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
                    if viewModel.friends.isEmpty {
                        Text("No friends yet")
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(viewModel.friends) { friend in
                            FriendRow(friend: friend)
                        }
                    }
                }
            }
            .navigationTitle("Friends")
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: { showAddFriend = true }) {
                        Image(systemName: "person.badge.plus")
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
            .alert("Error", isPresented: .constant(viewModel.errorMessage != nil)) {
                Button("OK") { viewModel.errorMessage = nil }
            } message: {
                Text(viewModel.errorMessage ?? "")
            }
            .task {
                await viewModel.loadData()
            }
        }
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
    @Published var isLoading = false

    func loadData() async {
        await refresh()
    }

    func refresh() async {
        isLoading = true
        defer { isLoading = false }

        do {
            async let friendsTask = APIService.shared.getFriends()
            async let requestsTask = APIService.shared.getPendingFriendRequests()

            let (loadedFriends, loadedRequests) = try await (friendsTask, requestsTask)
            friends = loadedFriends
            pendingRequests = loadedRequests
        } catch {
            errorMessage = "Failed to load friends"
        }
    }

    func acceptRequest(_ request: FriendRequest) {
        Task {
            do {
                try await APIService.shared.acceptFriendRequest(requestID: request.id)
                await refresh()
            } catch {
                errorMessage = "Failed to accept request"
            }
        }
    }

    func rejectRequest(_ request: FriendRequest) {
        Task {
            do {
                try await APIService.shared.rejectFriendRequest(requestID: request.id)
                await refresh()
            } catch {
                errorMessage = "Failed to reject request"
            }
        }
    }
}

#Preview {
    FriendsListView()
}
