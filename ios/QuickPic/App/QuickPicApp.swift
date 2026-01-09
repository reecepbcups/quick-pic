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
                MainTabView()
            } else {
                LoginView()
            }
        }
    }
}

struct LaunchView: View {
    var body: some View {
        ZStack {
            Color.black.ignoresSafeArea()

            VStack(spacing: 16) {
                Image(systemName: "camera.fill")
                    .font(.system(size: 60))
                    .foregroundColor(.yellow)

                Text("QuickPic")
                    .font(.largeTitle)
                    .fontWeight(.bold)
                    .foregroundColor(.white)

                ProgressView()
                    .tint(.white)
                    .padding(.top, 32)
            }
        }
    }
}

struct MainTabView: View {
    @State private var selectedTab = 1 // Start on camera
    @EnvironmentObject var authManager: AuthManager

    var body: some View {
        TabView(selection: $selectedTab) {
            InboxView()
                .tabItem {
                    Label("Inbox", systemImage: "tray.fill")
                }
                .tag(0)

            CameraView()
                .tabItem {
                    Label("Camera", systemImage: "camera.fill")
                }
                .tag(1)

            FriendsListView()
                .tabItem {
                    Label("Friends", systemImage: "person.2.fill")
                }
                .tag(2)

            ProfileView()
                .tabItem {
                    Label("Profile", systemImage: "person.circle.fill")
                }
                .tag(3)
        }
        .tint(.yellow)
    }
}

struct ProfileView: View {
    @EnvironmentObject var authManager: AuthManager

    var body: some View {
        NavigationStack {
            List {
                // User info section
                Section {
                    HStack(spacing: 16) {
                        Circle()
                            .fill(Color.yellow.opacity(0.3))
                            .frame(width: 60, height: 60)
                            .overlay(
                                Text(authManager.currentUser?.username.prefix(1).uppercased() ?? "?")
                                    .font(.title)
                                    .fontWeight(.bold)
                            )

                        VStack(alignment: .leading) {
                            Text(authManager.currentUser?.username ?? "Unknown")
                                .font(.headline)

                            Text("QuickPic User")
                                .font(.subheadline)
                                .foregroundColor(.secondary)
                        }
                    }
                    .padding(.vertical, 8)
                }

                // Security section
                Section("Security") {
                    HStack {
                        Image(systemName: "lock.shield.fill")
                            .foregroundColor(.green)
                        Text("End-to-end encrypted")
                    }

                    HStack {
                        Image(systemName: "key.fill")
                            .foregroundColor(.yellow)
                        Text("Private key stored on device")
                    }
                }

                // Actions section
                Section {
                    Button(action: { authManager.logout() }) {
                        HStack {
                            Image(systemName: "rectangle.portrait.and.arrow.right")
                            Text("Log Out")
                        }
                        .foregroundColor(.red)
                    }
                }

                // App info
                Section("About") {
                    HStack {
                        Text("Version")
                        Spacer()
                        Text("1.0.0")
                            .foregroundColor(.secondary)
                    }
                }
            }
            .navigationTitle("Profile")
        }
    }
}

#Preview {
    ContentView()
        .environmentObject(AuthManager())
}
