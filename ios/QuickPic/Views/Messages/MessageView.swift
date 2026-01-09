//
//  MessageView.swift
//  QuickPic
//
//  View-once message display with auto-hide after viewing
//  Note: This is now primarily used through ChatView's MessageContentView
//

import SwiftUI

struct MessageView: View {
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
                instructionsView
            }
        }
        .onAppear {
            // Auto-show after brief delay
            DispatchQueue.main.asyncAfter(deadline: .now() + 0.3) {
                withAnimation {
                    showContent = true
                }
            }
        }
        .gesture(
            DragGesture(minimumDistance: 0)
                .onEnded { _ in
                    // Dismiss when released
                    onDismiss()
                    dismiss()
                }
        )
    }

    @ViewBuilder
    private var contentView: some View {
        VStack {
            Spacer()

            // Content
            if message.contentType == .image {
                if let uiImage = UIImage(data: message.decryptedContent) {
                    Image(uiImage: uiImage)
                        .resizable()
                        .scaledToFit()
                        .cornerRadius(12)
                        .padding()
                } else {
                    Text("Unable to load image")
                        .foregroundColor(.red)
                }
            } else {
                if let text = String(data: message.decryptedContent, encoding: .utf8) {
                    Text(text)
                        .font(.title2)
                        .foregroundColor(.white)
                        .multilineTextAlignment(.center)
                        .padding(32)
                } else {
                    Text("Unable to decode message")
                        .foregroundColor(.red)
                }
            }

            Spacer()

            // Footer
            VStack(spacing: 8) {
                Text("Release to close")
                    .font(.caption)
                    .foregroundColor(.gray)

                if !message.hasBeenViewed {
                    Text("This message will be marked as viewed")
                        .font(.caption2)
                        .foregroundColor(.yellow)
                }
            }
            .padding(.bottom, 32)
        }
        .transition(.opacity)
    }

    private var instructionsView: some View {
        VStack(spacing: 16) {
            ProgressView()
                .tint(.white)

            Text("Loading message...")
                .foregroundColor(.gray)
        }
    }
}

#Preview {
    MessageView(
        message: StoredMessage(
            id: UUID(),
            conversationID: UUID(),
            contentType: .text,
            decryptedContent: "Hello, this is a test message!".data(using: .utf8)!,
            isFromMe: false,
            hasBeenViewed: false,
            serverDeleted: false,
            createdAt: Date(),
            receivedAt: Date()
        ),
        onDismiss: {}
    )
}
