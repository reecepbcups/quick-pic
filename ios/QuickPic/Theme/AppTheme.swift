//
//  AppTheme.swift
//  QuickPic
//
//  Design system and theme constants
//

import SwiftUI
import UIKit

// MARK: - Colors

extension Color {
    // Backgrounds
    static let appBackground = Color(hex: "0D0D0D")
    static let cardBackground = Color(hex: "1A1A1C")
    static let inputBackground = Color(hex: "2A2A2C")

    // Primary
    static let appPrimary = Color(hex: "02A4D3")

    // Text
    static let textPrimary = Color.white
    static let textSecondary = Color(hex: "8E8E93")

    // Status
    static let pending = Color(hex: "FF9F0A")
    static let danger = Color(hex: "FF453A")
    static let success = Color(hex: "30D158")

    // Initialize from hex
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 3: // RGB (12-bit)
            (a, r, g, b) = (255, (int >> 8) * 17, (int >> 4 & 0xF) * 17, (int & 0xF) * 17)
        case 6: // RGB (24-bit)
            (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        case 8: // ARGB (32-bit)
            (a, r, g, b) = (int >> 24, int >> 16 & 0xFF, int >> 8 & 0xFF, int & 0xFF)
        default:
            (a, r, g, b) = (1, 1, 1, 0)
        }
        self.init(
            .sRGB,
            red: Double(r) / 255,
            green: Double(g) / 255,
            blue: Double(b) / 255,
            opacity: Double(a) / 255
        )
    }
}

// MARK: - Typography

extension Font {
    static let appTitle = Font.system(size: 28, weight: .bold)
    static let appHeadline = Font.system(size: 17, weight: .semibold)
    static let appBody = Font.system(size: 15, weight: .regular)
    static let appCaption = Font.system(size: 13, weight: .regular)
    static let appSmall = Font.system(size: 11, weight: .medium)
}

// MARK: - Spacing

struct AppSpacing {
    static let xs: CGFloat = 4
    static let sm: CGFloat = 8
    static let md: CGFloat = 16
    static let lg: CGFloat = 24
    static let xl: CGFloat = 32
}

// MARK: - Corner Radius

struct AppRadius {
    static let sm: CGFloat = 8
    static let md: CGFloat = 12
    static let lg: CGFloat = 16
    static let xl: CGFloat = 24
    static let full: CGFloat = 9999
}

// MARK: - Custom View Modifiers

struct CardStyle: ViewModifier {
    func body(content: Content) -> some View {
        content
            .background(Color.cardBackground)
            .cornerRadius(AppRadius.lg)
    }
}

struct PillInputStyle: ViewModifier {
    func body(content: Content) -> some View {
        content
            .padding(.horizontal, AppSpacing.md)
            .padding(.vertical, 14)
            .background(Color.inputBackground)
            .cornerRadius(AppRadius.md)
    }
}

struct PrimaryButtonStyle: ButtonStyle {
    let isEnabled: Bool

    init(isEnabled: Bool = true) {
        self.isEnabled = isEnabled
    }

    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .font(.appHeadline)
            .foregroundColor(.black)
            .frame(maxWidth: .infinity)
            .padding(.vertical, 16)
            .background(isEnabled ? Color.appPrimary : Color.inputBackground)
            .cornerRadius(AppRadius.md)
            .scaleEffect(configuration.isPressed ? 0.98 : 1.0)
            .animation(.easeInOut(duration: 0.1), value: configuration.isPressed)
    }
}

struct SecondaryButtonStyle: ButtonStyle {
    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .font(.appHeadline)
            .foregroundColor(.textPrimary)
            .frame(maxWidth: .infinity)
            .padding(.vertical, 16)
            .background(Color.inputBackground)
            .cornerRadius(AppRadius.md)
            .scaleEffect(configuration.isPressed ? 0.98 : 1.0)
            .animation(.easeInOut(duration: 0.1), value: configuration.isPressed)
    }
}

extension View {
    func cardStyle() -> some View {
        modifier(CardStyle())
    }

    func pillInputStyle() -> some View {
        modifier(PillInputStyle())
    }
}

// MARK: - Custom Components

struct AppTextField: View {
    let icon: String
    let placeholder: String
    @Binding var text: String
    var isSecure: Bool = false
    @State private var showPassword = false

    var body: some View {
        HStack(spacing: AppSpacing.sm) {
            Image(systemName: icon)
                .foregroundColor(.textSecondary)
                .frame(width: 20)

            if isSecure && !showPassword {
                SecureField(placeholder, text: $text)
                    .foregroundColor(.textPrimary)
            } else {
                TextField(placeholder, text: $text)
                    .foregroundColor(.textPrimary)
            }

            if isSecure {
                Button(action: { showPassword.toggle() }) {
                    Image(systemName: showPassword ? "eye.slash" : "eye")
                        .foregroundColor(.textSecondary)
                }
            }
        }
        .pillInputStyle()
    }
}

struct StatusDot: View {
    enum Status {
        case unread
        case pending
        case read
    }

    let status: Status
    let initial: String
    var size: CGFloat = 44

    var backgroundColor: Color {
        switch status {
        case .unread: return .appPrimary
        case .pending: return .pending
        case .read: return Color(hex: "3A3A3C")
        }
    }

    var body: some View {
        Circle()
            .fill(backgroundColor)
            .frame(width: size, height: size)
            .overlay(
                Text(initial)
                    .font(.system(size: size * 0.4, weight: .semibold))
                    .foregroundColor(.white)
            )
    }
}

struct FloatingActionButton: View {
    let icon: String
    let action: () -> Void
    var size: CGFloat = 56

    var body: some View {
        Button(action: action) {
            ZStack {
                Circle()
                    .fill(Color.appPrimary)
                    .frame(width: size, height: size)
                    .shadow(color: Color.appPrimary.opacity(0.3), radius: 8, x: 0, y: 4)

                Image(systemName: icon)
                    .font(.system(size: size * 0.4, weight: .semibold))
                    .foregroundColor(.black)
            }
        }
        .buttonStyle(ScaleButtonStyle())
    }
}

struct ScaleButtonStyle: ButtonStyle {
    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .scaleEffect(configuration.isPressed ? 0.92 : 1.0)
            .animation(.spring(response: 0.3, dampingFraction: 0.6), value: configuration.isPressed)
    }
}

struct IconButton: View {
    let icon: String
    let action: () -> Void
    var size: CGFloat = 40

    var body: some View {
        Button(action: action) {
            ZStack {
                Circle()
                    .fill(Color.cardBackground)
                    .frame(width: size, height: size)

                Image(systemName: icon)
                    .font(.system(size: size * 0.4))
                    .foregroundColor(.textPrimary)
            }
        }
        .buttonStyle(ScaleButtonStyle())
    }
}

// MARK: - Haptics

@MainActor
struct Haptics {
    private static var lightGenerator: UIImpactFeedbackGenerator?
    private static var mediumGenerator: UIImpactFeedbackGenerator?
    private static var notificationGenerator: UINotificationFeedbackGenerator?

    static func prepare() {
        lightGenerator = UIImpactFeedbackGenerator(style: .light)
        mediumGenerator = UIImpactFeedbackGenerator(style: .medium)
        notificationGenerator = UINotificationFeedbackGenerator()
        lightGenerator?.prepare()
        mediumGenerator?.prepare()
        notificationGenerator?.prepare()
    }

    nonisolated static func light() {
        Task { @MainActor in
            if lightGenerator == nil {
                lightGenerator = UIImpactFeedbackGenerator(style: .light)
            }
            lightGenerator?.impactOccurred()
        }
    }

    nonisolated static func medium() {
        Task { @MainActor in
            if mediumGenerator == nil {
                mediumGenerator = UIImpactFeedbackGenerator(style: .medium)
            }
            mediumGenerator?.impactOccurred()
        }
    }

    nonisolated static func heavy() {
        Task { @MainActor in
            UIImpactFeedbackGenerator(style: .heavy).impactOccurred()
        }
    }

    nonisolated static func success() {
        Task { @MainActor in
            if notificationGenerator == nil {
                notificationGenerator = UINotificationFeedbackGenerator()
            }
            notificationGenerator?.notificationOccurred(.success)
        }
    }

    nonisolated static func error() {
        Task { @MainActor in
            if notificationGenerator == nil {
                notificationGenerator = UINotificationFeedbackGenerator()
            }
            notificationGenerator?.notificationOccurred(.error)
        }
    }
}
