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

    var body: some View {
        NavigationStack {
            ZStack {
                Color.black.ignoresSafeArea()

                if let capturedImage = viewModel.capturedImage {
                    // Preview captured image
                    CapturedImageView(
                        image: capturedImage,
                        onRetake: { viewModel.retake() },
                        onSend: { showFriendPicker = true }
                    )
                } else {
                    // Camera view
                    CameraPreviewView(session: viewModel.session)
                        .ignoresSafeArea()

                    VStack {
                        Spacer()

                        HStack(spacing: 60) {
                            // Flash toggle
                            Button(action: viewModel.toggleFlash) {
                                Image(systemName: viewModel.isFlashOn ? "bolt.fill" : "bolt.slash")
                                    .font(.title2)
                                    .foregroundColor(.white)
                            }

                            // Capture button
                            Button(action: viewModel.capturePhoto) {
                                ZStack {
                                    Circle()
                                        .stroke(Color.white, lineWidth: 4)
                                        .frame(width: 70, height: 70)

                                    Circle()
                                        .fill(Color.white)
                                        .frame(width: 58, height: 58)
                                }
                            }

                            // Switch camera
                            Button(action: viewModel.switchCamera) {
                                Image(systemName: "camera.rotate")
                                    .font(.title2)
                                    .foregroundColor(.white)
                            }
                        }
                        .padding(.bottom, 40)
                    }
                }
            }
            .navigationBarHidden(true)
            .sheet(isPresented: $showFriendPicker) {
                FriendPickerView(
                    imageData: viewModel.capturedImageData,
                    onSent: {
                        viewModel.retake()
                        showFriendPicker = false
                    }
                )
            }
            .onAppear {
                viewModel.startSession()
            }
            .onDisappear {
                viewModel.stopSession()
            }
        }
    }
}

struct CapturedImageView: View {
    let image: UIImage
    let onRetake: () -> Void
    let onSend: () -> Void

    var body: some View {
        ZStack {
            Image(uiImage: image)
                .resizable()
                .scaledToFit()

            VStack {
                Spacer()

                HStack(spacing: 60) {
                    Button(action: onRetake) {
                        VStack {
                            Image(systemName: "xmark")
                                .font(.title)
                            Text("Retake")
                                .font(.caption)
                        }
                        .foregroundColor(.white)
                    }

                    Button(action: onSend) {
                        VStack {
                            Image(systemName: "paperplane.fill")
                                .font(.title)
                            Text("Send")
                                .font(.caption)
                        }
                        .foregroundColor(.yellow)
                    }
                }
                .padding(.bottom, 40)
            }
        }
    }
}

struct FriendPickerView: View {
    let imageData: Data?
    let onSent: () -> Void

    @Environment(\.dismiss) private var dismiss
    @StateObject private var viewModel = FriendPickerViewModel()
    @State private var selectedFriend: Friend?
    @State private var isSending = false

    var body: some View {
        NavigationStack {
            Group {
                if viewModel.isLoading {
                    ProgressView("Loading friends...")
                } else if viewModel.friends.isEmpty {
                    VStack(spacing: 16) {
                        Image(systemName: "person.2.slash")
                            .font(.system(size: 50))
                            .foregroundColor(.secondary)
                        Text("No friends yet")
                            .font(.headline)
                        Text("Add friends to send them photos")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                } else {
                    List(viewModel.friends) { friend in
                        Button(action: { selectedFriend = friend }) {
                            HStack {
                                Circle()
                                    .fill(Color.yellow.opacity(0.3))
                                    .frame(width: 40, height: 40)
                                    .overlay(
                                        Text(friend.username.prefix(1).uppercased())
                                            .fontWeight(.semibold)
                                    )

                                Text(friend.username)

                                Spacer()

                                if selectedFriend?.id == friend.id {
                                    Image(systemName: "checkmark.circle.fill")
                                        .foregroundColor(.yellow)
                                }
                            }
                        }
                        .foregroundColor(.primary)
                    }
                }
            }
            .navigationTitle("Send to")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") { dismiss() }
                }
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Send") {
                        sendMessage()
                    }
                    .disabled(selectedFriend == nil || isSending)
                }
            }
            .task {
                await viewModel.loadFriends()
            }
        }
    }

    private func sendMessage() {
        guard let friend = selectedFriend, let data = imageData else { return }

        isSending = true
        Task {
            do {
                try await viewModel.sendImage(data, to: friend)
                onSent()
            } catch {
                print("Failed to send: \(error)")
            }
            isSending = false
        }
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

    func startSession() {
        guard !session.isRunning else { return }

        Task {
            await setupSession()
            session.startRunning()
        }
    }

    func stopSession() {
        session.stopRunning()
    }

    private func setupSession() async {
        session.beginConfiguration()
        session.sessionPreset = .photo

        // Add camera input
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
        currentCameraPosition = currentCameraPosition == .back ? .front : .back

        session.beginConfiguration()

        // Remove existing input
        session.inputs.forEach { session.removeInput($0) }

        // Add new input
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

extension CameraViewModel: AVCapturePhotoCaptureDelegate {
    nonisolated func photoOutput(_ output: AVCapturePhotoOutput, didFinishProcessingPhoto photo: AVCapturePhoto, error: Error?) {
        guard let data = photo.fileDataRepresentation(),
              let image = UIImage(data: data) else {
            return
        }

        // Convert to PNG for lossless quality
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

        // Save sent message locally to the conversation
        let storedMessage = StoredMessage(
            id: response.id,
            conversationID: friend.userID,
            contentType: .image,
            decryptedContent: imageData,
            isFromMe: true,
            hasBeenViewed: true,
            serverDeleted: true,
            createdAt: response.createdAt,
            receivedAt: Date()
        )

        db.saveMessage(storedMessage)

        // Ensure conversation exists
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
