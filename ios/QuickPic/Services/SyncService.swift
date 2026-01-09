//
//  SyncService.swift
//  QuickPic
//
//  Background service for syncing viewed messages with the server
//  Ensures messages are deleted from server after being viewed
//

import Foundation

final class SyncService: @unchecked Sendable {
    static let shared = SyncService()

    private let db = DatabaseService.shared
    private let api = APIService.shared
    private var syncTimer: Timer?
    private let syncInterval: TimeInterval = 60 // Sync every 60 seconds

    private init() {
        startSyncTimer()
    }

    // MARK: - Public API

    /// Start the background sync timer
    func startSyncTimer() {
        stopSyncTimer()

        syncTimer = Timer.scheduledTimer(withTimeInterval: syncInterval, repeats: true) { [weak self] _ in
            Task {
                await self?.syncViewedMessages()
            }
        }

        // Also run immediately
        Task {
            await syncViewedMessages()
        }
    }

    /// Stop the background sync timer
    func stopSyncTimer() {
        syncTimer?.invalidate()
        syncTimer = nil
    }

    /// Manually trigger a sync
    func syncNow() async {
        await syncViewedMessages()
    }

    // MARK: - Private

    private func syncViewedMessages() async {
        let messagesNeedingDeletion = db.getMessagesNeedingServerDeletion()

        for message in messagesNeedingDeletion {
            do {
                try await api.acknowledgeMessage(id: message.id)
                db.markMessageServerDeleted(messageID: message.id)
                print("Successfully deleted message \(message.id) from server")
            } catch {
                // Will retry on next sync cycle
                print("Failed to delete message \(message.id) from server: \(error)")
            }
        }
    }
}
