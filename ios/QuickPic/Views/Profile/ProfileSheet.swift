//
//  ProfileSheet.swift
//  QuickPic
//

import SwiftUI

struct ProfileSheet: View {
    @EnvironmentObject var authManager: AuthManager
    @Environment(\.dismiss) var dismiss

    var body: some View {
        ZStack {
            Color.appBackground.ignoresSafeArea()

            VStack(spacing: AppSpacing.lg) {
                // Header
                VStack(spacing: AppSpacing.md) {
                    // Avatar
                    ZStack {
                        Circle()
                            .fill(Color.appPrimary.opacity(0.2))
                            .frame(width: 80, height: 80)

                        Text(authManager.currentUser?.username.prefix(1).uppercased() ?? "?")
                            .font(.system(size: 32, weight: .bold))
                            .foregroundColor(.appPrimary)
                    }

                    // Username
                    Text(authManager.currentUser?.username ?? "Unknown")
                        .font(.appTitle)
                        .foregroundColor(.textPrimary)

                    Text("QuickPic User #\(authManager.currentUser?.userNumber ?? 0)")
                        .font(.appCaption)
                        .foregroundColor(.textSecondary)
                }
                .padding(.top, AppSpacing.lg)

                Divider()
                    .background(Color.inputBackground)

                // Scrollable settings section
                ScrollView {
                    VStack(spacing: AppSpacing.lg) {
                        // Security info
                        VStack(spacing: AppSpacing.md) {
                            SecurityRow(
                                icon: "lock.shield.fill",
                                iconColor: .success,
                                title: "End-to-end encrypted",
                                subtitle: "Messages are encrypted on your device"
                            )

                            SecurityRow(
                                icon: "key.fill",
                                iconColor: .appPrimary,
                                title: "Private key on device",
                                subtitle: "Your key never leaves this device"
                            )
                        }
                        .padding(.horizontal, AppSpacing.md)

                        // Logout button
                        Button(action: logout) {
                            HStack {
                                Image(systemName: "rectangle.portrait.and.arrow.right")
                                Text("Log Out")
                            }
                            .font(.appHeadline)
                            .foregroundColor(.danger)
                            .frame(maxWidth: .infinity)
                            .padding(.vertical, 16)
                            .background(Color.danger.opacity(0.1))
                            .cornerRadius(AppRadius.md)
                        }
                        .padding(.horizontal, AppSpacing.md)

                        // Version
                        Text("Version 1.0.0")
                            .font(.appSmall)
                            .foregroundColor(.textSecondary)
                            .padding(.bottom, AppSpacing.md)
                    }
                }
            }
        }
    }

    private func logout() {
        Haptics.medium()
        authManager.logout()
        dismiss()
    }
}

struct SecurityRow: View {
    let icon: String
    let iconColor: Color
    let title: String
    let subtitle: String

    var body: some View {
        HStack(spacing: AppSpacing.md) {
            Image(systemName: icon)
                .font(.title3)
                .foregroundColor(iconColor)
                .frame(width: 32)

            VStack(alignment: .leading, spacing: 2) {
                Text(title)
                    .font(.appBody)
                    .foregroundColor(.textPrimary)

                Text(subtitle)
                    .font(.appSmall)
                    .foregroundColor(.textSecondary)
            }

            Spacer()
        }
        .padding(AppSpacing.md)
        .background(Color.cardBackground)
        .cornerRadius(AppRadius.md)
    }
}

#Preview {
    ProfileSheet()
        .environmentObject(AuthManager())
}
