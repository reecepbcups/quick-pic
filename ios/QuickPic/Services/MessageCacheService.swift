import Foundation

/// Manages local storage of decrypted messages with 24-hour expiration
final class MessageCacheService {
    static let shared = MessageCacheService()

    private let cacheKey = "cached_messages"
    private let cacheExpirationHours: TimeInterval = 24
    private let fileManager = FileManager.default

    private var cacheURL: URL {
        let paths = fileManager.urls(for: .documentDirectory, in: .userDomainMask)
        return paths[0].appendingPathComponent("message_cache.json")
    }

    private init() {
        // Start cleanup timer
        startCleanupTimer()
    }

    // MARK: - Public API

    /// Store a decrypted message in cache
    func cache(message: CachedMessage) {
        var messages = loadMessages()
        messages.append(message)
        saveMessages(messages)
    }

    /// Get all cached messages (excludes expired)
    func getCachedMessages() -> [CachedMessage] {
        let messages = loadMessages()
        return messages.filter { !$0.isExpired }
    }

    /// Get unviewed messages count
    func getUnviewedCount() -> Int {
        getCachedMessages().filter { !$0.hasBeenViewed }.count
    }

    /// Mark a message as viewed
    func markAsViewed(messageID: UUID) {
        var messages = loadMessages()
        if let index = messages.firstIndex(where: { $0.id == messageID }) {
            messages[index].hasBeenViewed = true
            saveMessages(messages)
        }
    }

    /// Delete a specific message
    func delete(messageID: UUID) {
        var messages = loadMessages()
        messages.removeAll { $0.id == messageID }
        saveMessages(messages)
    }

    /// Clear all cached messages (called on logout)
    func clearAll() {
        try? fileManager.removeItem(at: cacheURL)
    }

    /// Remove expired messages
    func purgeExpired() {
        var messages = loadMessages()
        let beforeCount = messages.count
        messages.removeAll { $0.isExpired }
        let removedCount = beforeCount - messages.count

        if removedCount > 0 {
            saveMessages(messages)
            print("Purged \(removedCount) expired messages")
        }
    }

    // MARK: - Private Helpers

    private func loadMessages() -> [CachedMessage] {
        guard fileManager.fileExists(atPath: cacheURL.path) else {
            return []
        }

        do {
            let data = try Data(contentsOf: cacheURL)
            let decoder = JSONDecoder()
            decoder.dateDecodingStrategy = .iso8601
            return try decoder.decode([CachedMessage].self, from: data)
        } catch {
            print("Failed to load message cache: \(error)")
            return []
        }
    }

    private func saveMessages(_ messages: [CachedMessage]) {
        do {
            let encoder = JSONEncoder()
            encoder.dateEncodingStrategy = .iso8601
            let data = try encoder.encode(messages)
            try data.write(to: cacheURL, options: .atomic)
        } catch {
            print("Failed to save message cache: \(error)")
        }
    }

    private func startCleanupTimer() {
        // Run cleanup every hour
        Timer.scheduledTimer(withTimeInterval: 3600, repeats: true) { [weak self] _ in
            self?.purgeExpired()
        }
    }
}
