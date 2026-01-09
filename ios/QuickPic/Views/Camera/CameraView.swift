//
//  CameraView.swift
//  QuickPic
//
//  Camera view for capturing and sending photos to friends
//

import SwiftUI
import AVFoundation

struct CameraView: View {
    @StateObject private var viewModel = CameraViewModel()
    @State private var showFriendPicker = false
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        ZStack {
            Color.appBackground.ignoresSafeArea()

            if let capturedImage = viewModel.capturedImage {
                CapturedImageView(
                    image: capturedImage,
                    onRetake: {
                        Haptics.light()
                        viewModel.retake()
                    },
                    onSend: {
                        Haptics.medium()
                        showFriendPicker = true
                    }
                )
            } else {
                // Camera preview
                CameraPreviewView(session: viewModel.session)
                    .ignoresSafeArea()

                // Controls overlay
                VStack {
                    // Top bar with close button
                    HStack {
                        Button(action: {
                            Haptics.light()
                            dismiss()
                        }) {
                            Image(systemName: "xmark")
                                .font(.system(size: 18, weight: .semibold))
                                .foregroundColor(.white)
                                .frame(width: 44, height: 44)
                                .background(Color.black.opacity(0.4))
                                .clipShape(Circle())
                        }

                        Spacer()
                    }
                    .padding(AppSpacing.md)

                    Spacer()

                    // Bottom controls
                    HStack(spacing: 50) {
                        // Flash toggle
                        Button(action: {
                            Haptics.light()
                            viewModel.toggleFlash()
                        }) {
                            Image(systemName: viewModel.isFlashOn ? "bolt.fill" : "bolt.slash")
                                .font(.system(size: 22))
                                .foregroundColor(.white)
                                .frame(width: 50, height: 50)
                        }

                        // Capture button
                        Button(action: {
                            Haptics.heavy()
                            viewModel.capturePhoto()
                        }) {
                            ZStack {
                                Circle()
                                    .stroke(Color.white, lineWidth: 4)
                                    .frame(width: 76, height: 76)

                                Circle()
                                    .fill(Color.white)
                                    .frame(width: 62, height: 62)
                            }
                        }
                        .buttonStyle(ScaleButtonStyle())

                        // Switch camera
                        Button(action: {
                            Haptics.light()
                            viewModel.switchCamera()
                        }) {
                            Image(systemName: "camera.rotate")
                                .font(.system(size: 22))
                                .foregroundColor(.white)
                                .frame(width: 50, height: 50)
                        }
                    }
                    .padding(.bottom, 50)
                }
            }
        }
        .sheet(isPresented: $showFriendPicker) {
            FriendPickerSheet(
                imageData: viewModel.capturedImageData,
                onSent: {
                    viewModel.retake()
                    showFriendPicker = false
                    Haptics.success()
                    dismiss()
                }
            )
            .presentationDetents([.medium, .large])
            .presentationDragIndicator(.visible)
        }
        .onAppear {
            viewModel.startSession()
        }
        .onDisappear {
            viewModel.stopSession()
        }
    }
}

struct CapturedImageView: View {
    let image: UIImage
    let onRetake: () -> Void
    let onSend: () -> Void

    var body: some View {
        ZStack {
            Color.appBackground.ignoresSafeArea()

            Image(uiImage: image)
                .resizable()
                .scaledToFit()

            VStack {
                Spacer()

                HStack(spacing: 80) {
                    // Retake button
                    Button(action: onRetake) {
                        VStack(spacing: AppSpacing.sm) {
                            Image(systemName: "xmark")
                                .font(.system(size: 24, weight: .semibold))
                                .foregroundColor(.white)
                                .frame(width: 56, height: 56)
                                .background(Color.cardBackground)
                                .clipShape(Circle())

                            Text("Retake")
                                .font(.appCaption)
                                .foregroundColor(.textSecondary)
                        }
                    }
                    .buttonStyle(ScaleButtonStyle())

                    // Send button
                    Button(action: onSend) {
                        VStack(spacing: AppSpacing.sm) {
                            Image(systemName: "paperplane.fill")
                                .font(.system(size: 24, weight: .semibold))
                                .foregroundColor(.black)
                                .frame(width: 56, height: 56)
                                .background(Color.appPrimary)
                                .clipShape(Circle())

                            Text("Send")
                                .font(.appCaption)
                                .foregroundColor(.appPrimary)
                        }
                    }
                    .buttonStyle(ScaleButtonStyle())
                }
                .padding(.bottom, 50)
            }
        }
    }
}

struct FriendPickerSheet: View {
    let imageData: Data?
    let onSent: () -> Void

    @Environment(\.dismiss) private var dismiss
    @StateObject private var viewModel = FriendPickerViewModel()
    @State private var selectedFriend: Friend?
    @State private var isSending = false

    var body: some View {
        ZStack {
            Color.appBackground.ignoresSafeArea()

            VStack(spacing: 0) {
                // Header
                VStack(spacing: AppSpacing.sm) {
                    Text("Send to")
                        .font(.appTitle)
                        .foregroundColor(.textPrimary)
                }
                .padding(.top, AppSpacing.lg)
                .padding(.bottom, AppSpacing.md)

                if viewModel.isLoading {
                    Spacer()
                    ProgressView()
                        .tint(.appPrimary)
                    Spacer()
                } else if viewModel.friends.isEmpty {
                    Spacer()
                    VStack(spacing: AppSpacing.md) {
                        Image(systemName: "person.2.slash")
                            .font(.system(size: 50))
                            .foregroundColor(.textSecondary)

                        Text("No friends yet")
                            .font(.appHeadline)
                            .foregroundColor(.textPrimary)

                        Text("Add friends to send them photos")
                            .font(.appCaption)
                            .foregroundColor(.textSecondary)
                    }
                    Spacer()
                } else {
                    ScrollView {
                        LazyVStack(spacing: AppSpacing.sm) {
                            ForEach(viewModel.friends) { friend in
                                FriendPickerRow(
                                    friend: friend,
                                    isSelected: selectedFriend?.id == friend.id
                                ) {
                                    Haptics.light()
                                    selectedFriend = friend
                                }
                                .padding(.horizontal, AppSpacing.md)
                            }
                        }
                        .padding(.top, AppSpacing.sm)
                    }
                }

                // Send button
                Button(action: sendMessage) {
                    HStack {
                        if isSending {
                            ProgressView()
                                .tint(.black)
                        } else {
                            Text("Send")
                        }
                    }
                }
                .buttonStyle(PrimaryButtonStyle(isEnabled: selectedFriend != nil && !isSending))
                .disabled(selectedFriend == nil || isSending)
                .padding(AppSpacing.md)
            }
        }
        .task {
            await viewModel.loadFriends()
        }
    }

    private func sendMessage() {
        guard let friend = selectedFriend, let data = imageData else { return }
        Haptics.light()

        isSending = true
        Task {
            do {
                try await viewModel.sendImage(data, to: friend)
                onSent()
            } catch {
                Haptics.error()
                print("Failed to send: \(error)")
            }
            isSending = false
        }
    }
}

struct FriendPickerRow: View {
    let friend: Friend
    let isSelected: Bool
    let onTap: () -> Void

    var body: some View {
        Button(action: onTap) {
            HStack(spacing: AppSpacing.md) {
                StatusDot(
                    status: .read,
                    initial: String(friend.username.prefix(1)).uppercased()
                )

                Text(friend.username)
                    .font(.appHeadline)
                    .foregroundColor(.textPrimary)

                Spacer()

                if isSelected {
                    Image(systemName: "checkmark.circle.fill")
                        .font(.title2)
                        .foregroundColor(.appPrimary)
                } else {
                    Circle()
                        .stroke(Color.textSecondary, lineWidth: 2)
                        .frame(width: 24, height: 24)
                }
            }
            .padding(AppSpacing.md)
            .background(isSelected ? Color.appPrimary.opacity(0.1) : Color.cardBackground)
            .cornerRadius(AppRadius.lg)
        }
        .buttonStyle(FeedRowButtonStyle())
    }
}

// MARK: - View Models

@MainActor
class CameraViewModel: NSObject, ObservableObject {
    @Published var capturedImage: UIImage?
    @Published var capturedImageData: Data?
    @Published var isFlashOn = false

    let session = AVCaptureSession()
    private var photoOutput = AVCapturePhotoOutput()
    private var currentCameraPosition: AVCaptureDevice.Position = .back
    private let sessionQueue = DispatchQueue(label: "camera.session.queue")

    func startSession() {
        sessionQueue.async { [self] in
            guard !session.isRunning else { return }
            setupSession()
            session.startRunning()
        }
    }

    func stopSession() {
        sessionQueue.async { [self] in
            session.stopRunning()
        }
    }

    private func setupSession() {
        session.beginConfiguration()
        session.sessionPreset = .photo

        guard let camera = AVCaptureDevice.default(.builtInWideAngleCamera, for: .video, position: currentCameraPosition),
              let input = try? AVCaptureDeviceInput(device: camera) else {
            session.commitConfiguration()
            return
        }

        if session.canAddInput(input) {
            session.addInput(input)
        }

        if session.canAddOutput(photoOutput) {
            session.addOutput(photoOutput)
        }

        session.commitConfiguration()
    }

    func capturePhoto() {
        let settings = AVCapturePhotoSettings()
        settings.flashMode = isFlashOn ? .on : .off
        photoOutput.capturePhoto(with: settings, delegate: self)
    }

    func retake() {
        capturedImage = nil
        capturedImageData = nil
    }

    func toggleFlash() {
        isFlashOn.toggle()
    }

    func switchCamera() {
        sessionQueue.async { [self] in
            currentCameraPosition = currentCameraPosition == .back ? .front : .back

            session.beginConfiguration()
            session.inputs.forEach { session.removeInput($0) }

            guard let camera = AVCaptureDevice.default(.builtInWideAngleCamera, for: .video, position: currentCameraPosition),
                  let input = try? AVCaptureDeviceInput(device: camera) else {
                session.commitConfiguration()
                return
            }

            if session.canAddInput(input) {
                session.addInput(input)
            }

            session.commitConfiguration()
        }
    }
}

extension CameraViewModel: AVCapturePhotoCaptureDelegate {
    nonisolated func photoOutput(_ output: AVCapturePhotoOutput, didFinishProcessingPhoto photo: AVCapturePhoto, error: Error?) {
        guard let data = photo.fileDataRepresentation(),
              let image = UIImage(data: data) else {
            return
        }

        if let pngData = image.pngData() {
            Task { @MainActor in
                self.capturedImage = image
                self.capturedImageData = pngData
            }
        }
    }
}

@MainActor
class FriendPickerViewModel: ObservableObject {
    @Published var friends: [Friend] = []
    @Published var isLoading = false

    private let api = APIService.shared
    private let crypto = CryptoService.shared
    private let db = DatabaseService.shared

    func loadFriends() async {
        isLoading = true
        defer { isLoading = false }

        do {
            friends = try await api.getFriends()
        } catch {
            print("Failed to load friends: \(error)")
        }
    }

    func sendImage(_ imageData: Data, to friend: Friend) async throws {
        let privateKey = try crypto.getPrivateKey()
        let recipientPublicKey = try crypto.publicKeyFromBase64(friend.publicKey)

        let (encryptedData, signature) = try crypto.encrypt(
            content: imageData,
            recipientPublicKey: recipientPublicKey,
            senderPrivateKey: privateKey
        )

        let response = try await api.sendMessage(
            to: friend.username,
            encryptedContent: encryptedData,
            contentType: .image,
            signature: signature
        )

        let storedMessage = StoredMessage(
            id: response.id,
            conversationID: friend.userID,
            contentType: .image,
            decryptedContent: imageData,
            encryptedContent: encryptedData,
            isFromMe: true,
            hasBeenViewed: true,
            createdAt: response.createdAt,
            receivedAt: Date()
        )

        db.saveMessage(storedMessage)
        _ = db.getOrCreateConversation(for: friend)
    }
}

// MARK: - Camera Preview

struct CameraPreviewView: UIViewRepresentable {
    let session: AVCaptureSession

    func makeUIView(context: Context) -> UIView {
        let view = UIView(frame: .zero)

        let previewLayer = AVCaptureVideoPreviewLayer(session: session)
        previewLayer.videoGravity = .resizeAspectFill
        view.layer.addSublayer(previewLayer)

        context.coordinator.previewLayer = previewLayer

        return view
    }

    func updateUIView(_ uiView: UIView, context: Context) {
        context.coordinator.previewLayer?.frame = uiView.bounds
    }

    func makeCoordinator() -> Coordinator {
        Coordinator()
    }

    class Coordinator {
        var previewLayer: AVCaptureVideoPreviewLayer?
    }
}

#Preview {
    CameraView()
}
