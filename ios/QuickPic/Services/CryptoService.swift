import Foundation
import CryptoKit
import Compression

enum CryptoError: Error {
    case keyGenerationFailed
    case encryptionFailed
    case decryptionFailed
    case compressionFailed
    case decompressionFailed
    case invalidPublicKey
    case signingFailed
    case verificationFailed
}

/// Handles all encryption operations using X25519 + XChaCha20-Poly1305
final class CryptoService {
    static let shared = CryptoService()

    private init() {}

    // MARK: - Key Generation

    /// Generate a new X25519 key pair
    func generateKeyPair() -> (privateKey: Curve25519.KeyAgreement.PrivateKey, publicKey: Curve25519.KeyAgreement.PublicKey) {
        let privateKey = Curve25519.KeyAgreement.PrivateKey()
        return (privateKey, privateKey.publicKey)
    }

    /// Get public key as base64 string for storage/transmission
    func publicKeyToBase64(_ publicKey: Curve25519.KeyAgreement.PublicKey) -> String {
        publicKey.rawRepresentation.base64EncodedString()
    }

    /// Reconstruct public key from base64 string
    func publicKeyFromBase64(_ base64String: String) throws -> Curve25519.KeyAgreement.PublicKey {
        guard let data = Data(base64Encoded: base64String) else {
            throw CryptoError.invalidPublicKey
        }
        return try Curve25519.KeyAgreement.PublicKey(rawRepresentation: data)
    }

    /// Store private key securely
    func storePrivateKey(_ privateKey: Curve25519.KeyAgreement.PrivateKey) throws {
        try KeychainService.shared.storePrivateKey(privateKey.rawRepresentation)
    }

    /// Retrieve private key from keychain
    func getPrivateKey() throws -> Curve25519.KeyAgreement.PrivateKey {
        let data = try KeychainService.shared.getPrivateKey()
        return try Curve25519.KeyAgreement.PrivateKey(rawRepresentation: data)
    }

    // MARK: - Encryption

    /// Encrypt message content for a recipient
    /// Returns: (encryptedData, signature) - both base64 encoded where appropriate
    func encrypt(
        content: Data,
        recipientPublicKey: Curve25519.KeyAgreement.PublicKey,
        senderPrivateKey: Curve25519.KeyAgreement.PrivateKey
    ) throws -> (encryptedData: Data, signature: String) {
        // 1. Compress content with gzip
        let compressedData = try compress(content)

        // 2. Generate ephemeral symmetric key
        let symmetricKey = SymmetricKey(size: .bits256)

        // 3. Encrypt content with symmetric key using ChaChaPoly (similar to XChaCha20-Poly1305)
        let sealedBox = try ChaChaPoly.seal(compressedData, using: symmetricKey)
        let encryptedContent = sealedBox.combined

        // 4. Derive shared secret and encrypt symmetric key
        let sharedSecret = try senderPrivateKey.sharedSecretFromKeyAgreement(with: recipientPublicKey)
        let derivedKey = sharedSecret.hkdfDerivedSymmetricKey(
            using: SHA256.self,
            salt: Data(),
            sharedInfo: Data("QuickPic-Key-Encryption".utf8),
            outputByteCount: 32
        )

        let symmetricKeyData = symmetricKey.withUnsafeBytes { Data($0) }
        let encryptedSymmetricKey = try ChaChaPoly.seal(symmetricKeyData, using: derivedKey)

        // 5. Combine: [encrypted_symmetric_key_length (4 bytes)][encrypted_symmetric_key][encrypted_content]
        var result = Data()
        let keyLength = UInt32(encryptedSymmetricKey.combined.count)
        result.append(contentsOf: withUnsafeBytes(of: keyLength.bigEndian) { Data($0) })
        result.append(encryptedSymmetricKey.combined)
        result.append(encryptedContent)

        // 6. Sign the encrypted data
        let signingKey = Curve25519.Signing.PrivateKey(rawRepresentation: senderPrivateKey.rawRepresentation)
        let signature = try signingKey.signature(for: result)

        return (result, signature.base64EncodedString())
    }

    // MARK: - Decryption

    /// Decrypt message content from a sender
    func decrypt(
        encryptedData: Data,
        signature: String,
        senderPublicKey: Curve25519.KeyAgreement.PublicKey,
        recipientPrivateKey: Curve25519.KeyAgreement.PrivateKey
    ) throws -> Data {
        // 1. Verify signature
        guard let signatureData = Data(base64Encoded: signature) else {
            throw CryptoError.verificationFailed
        }

        let verifyingKey = try Curve25519.Signing.PublicKey(rawRepresentation: senderPublicKey.rawRepresentation)
        guard verifyingKey.isValidSignature(signatureData, for: encryptedData) else {
            throw CryptoError.verificationFailed
        }

        // 2. Extract encrypted symmetric key length
        guard encryptedData.count >= 4 else {
            throw CryptoError.decryptionFailed
        }

        let keyLengthData = encryptedData.prefix(4)
        let keyLength = keyLengthData.withUnsafeBytes { $0.load(as: UInt32.self).bigEndian }

        guard encryptedData.count >= 4 + Int(keyLength) else {
            throw CryptoError.decryptionFailed
        }

        // 3. Extract encrypted symmetric key and encrypted content
        let encryptedKeyData = encryptedData.subdata(in: 4..<(4 + Int(keyLength)))
        let encryptedContent = encryptedData.subdata(in: (4 + Int(keyLength))..<encryptedData.count)

        // 4. Derive shared secret and decrypt symmetric key
        let sharedSecret = try recipientPrivateKey.sharedSecretFromKeyAgreement(with: senderPublicKey)
        let derivedKey = sharedSecret.hkdfDerivedSymmetricKey(
            using: SHA256.self,
            salt: Data(),
            sharedInfo: Data("QuickPic-Key-Encryption".utf8),
            outputByteCount: 32
        )

        let encryptedKeyBox = try ChaChaPoly.SealedBox(combined: encryptedKeyData)
        let symmetricKeyData = try ChaChaPoly.open(encryptedKeyBox, using: derivedKey)
        let symmetricKey = SymmetricKey(data: symmetricKeyData)

        // 5. Decrypt content
        let contentBox = try ChaChaPoly.SealedBox(combined: encryptedContent)
        let compressedData = try ChaChaPoly.open(contentBox, using: symmetricKey)

        // 6. Decompress
        return try decompress(compressedData)
    }

    // MARK: - Compression

    private func compress(_ data: Data) throws -> Data {
        var compressedData = Data()

        let pageSize = 128
        var index = 0

        let outputFilter = try OutputFilter(.compress, using: .zlib) { (data: Data?) in
            if let data = data {
                compressedData.append(data)
            }
        }

        while index < data.count {
            let rangeLength = min(pageSize, data.count - index)
            let subdata = data.subdata(in: index..<(index + rangeLength))
            try outputFilter.write(subdata)
            index += rangeLength
        }

        try outputFilter.finalize()

        return compressedData
    }

    private func decompress(_ data: Data) throws -> Data {
        var decompressedData = Data()

        let pageSize = 128
        var index = 0

        let inputFilter = try InputFilter(.decompress, using: .zlib) { (length: Int) -> Data? in
            guard index < data.count else { return nil }
            let rangeLength = min(length, data.count - index)
            let subdata = data.subdata(in: index..<(index + rangeLength))
            index += rangeLength
            return subdata
        }

        while let page = try inputFilter.readData(ofLength: pageSize) {
            decompressedData.append(page)
        }

        return decompressedData
    }
}
