// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title QuickPicStorage
 * @notice On-chain storage for QuickPic messaging app - replaces SQLite database
 * @dev All data is stored on-chain with the server submitting transactions on behalf of users
 */
contract QuickPicStorage {
    // ============ Enums ============

    enum FriendRequestStatus { Pending, Accepted, Rejected }
    enum ContentType { Text, Image }

    // ============ Structs ============

    struct User {
        bytes32 id;
        uint256 userNumber;
        string username;
        string passwordHash;  // Argon2id hash (hashing done off-chain)
        string publicKey;     // Base64-encoded X25519 public key
        uint256 createdAt;
        uint256 updatedAt;
        bool exists;
    }

    struct FriendRequest {
        bytes32 id;
        bytes32 fromUserId;
        bytes32 toUserId;
        FriendRequestStatus status;
        uint256 createdAt;
        bool exists;
    }

    struct Friendship {
        bytes32 id;
        bytes32 userAId;  // Lexicographically smaller
        bytes32 userBId;  // Lexicographically larger
        uint256 createdAt;
        bool exists;
    }

    struct Message {
        bytes32 id;
        bytes32 fromUserId;
        bytes32 toUserId;
        bytes encryptedContent;
        ContentType contentType;
        string signature;
        uint256 createdAt;
        bool exists;
    }

    // ============ State Variables ============

    address public owner;
    uint256 public nextUserNumber;

    // User storage
    mapping(bytes32 => User) public users;           // id => User
    mapping(string => bytes32) public usernameToId;  // username => id
    bytes32[] public userIds;

    // Friend request storage
    mapping(bytes32 => FriendRequest) public friendRequests;  // id => FriendRequest
    mapping(bytes32 => mapping(bytes32 => bytes32)) public friendRequestByUsers;  // fromUserId => toUserId => requestId
    mapping(bytes32 => bytes32[]) public pendingRequestsTo;   // toUserId => requestIds[]
    bytes32[] public friendRequestIds;

    // Friendship storage
    mapping(bytes32 => Friendship) public friendships;  // id => Friendship
    mapping(bytes32 => mapping(bytes32 => bytes32)) public friendshipByUsers;  // userAId => userBId => friendshipId
    mapping(bytes32 => bytes32[]) public userFriendships;  // userId => friendshipIds[]
    bytes32[] public friendshipIds;

    // Message storage
    mapping(bytes32 => Message) public messages;  // id => Message
    mapping(bytes32 => bytes32[]) public messagesToUser;  // toUserId => messageIds[]
    mapping(bytes32 => bytes32[]) public messagesFromUser;  // fromUserId => messageIds[]
    bytes32[] public messageIds;

    // ============ Events ============

    event UserCreated(bytes32 indexed id, uint256 userNumber, string username);
    event UserUpdated(bytes32 indexed id);
    event FriendRequestCreated(bytes32 indexed id, bytes32 indexed fromUserId, bytes32 indexed toUserId);
    event FriendRequestUpdated(bytes32 indexed id, FriendRequestStatus status);
    event FriendshipCreated(bytes32 indexed id, bytes32 indexed userAId, bytes32 indexed userBId);
    event MessageCreated(bytes32 indexed id, bytes32 indexed fromUserId, bytes32 indexed toUserId);
    event MessageDeleted(bytes32 indexed id);

    // ============ Modifiers ============

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    // ============ Constructor ============

    constructor() {
        owner = msg.sender;
        nextUserNumber = 1;
    }

    // ============ User Functions ============

    function createUser(
        bytes32 id,
        string calldata username,
        string calldata passwordHash,
        string calldata publicKey
    ) external onlyOwner returns (uint256 userNumber) {
        require(!users[id].exists, "User already exists");
        require(usernameToId[username] == bytes32(0), "Username already taken");
        require(bytes(username).length > 0, "Username cannot be empty");

        userNumber = nextUserNumber++;

        users[id] = User({
            id: id,
            userNumber: userNumber,
            username: username,
            passwordHash: passwordHash,
            publicKey: publicKey,
            createdAt: block.timestamp,
            updatedAt: block.timestamp,
            exists: true
        });

        usernameToId[username] = id;
        userIds.push(id);

        emit UserCreated(id, userNumber, username);
    }

    function getUser(bytes32 id) external view returns (
        bytes32 userId,
        uint256 userNumber,
        string memory username,
        string memory passwordHash,
        string memory publicKey,
        uint256 createdAt,
        uint256 updatedAt
    ) {
        User storage user = users[id];
        require(user.exists, "User not found");

        return (
            user.id,
            user.userNumber,
            user.username,
            user.passwordHash,
            user.publicKey,
            user.createdAt,
            user.updatedAt
        );
    }

    function getUserByUsername(string calldata username) external view returns (
        bytes32 userId,
        uint256 userNumber,
        string memory usernameOut,
        string memory passwordHash,
        string memory publicKey,
        uint256 createdAt,
        uint256 updatedAt
    ) {
        bytes32 id = usernameToId[username];
        require(id != bytes32(0), "User not found");

        User storage user = users[id];
        return (
            user.id,
            user.userNumber,
            user.username,
            user.passwordHash,
            user.publicKey,
            user.createdAt,
            user.updatedAt
        );
    }

    function updateUser(
        bytes32 id,
        string calldata passwordHash,
        string calldata publicKey
    ) external onlyOwner {
        require(users[id].exists, "User not found");

        users[id].passwordHash = passwordHash;
        users[id].publicKey = publicKey;
        users[id].updatedAt = block.timestamp;

        emit UserUpdated(id);
    }

    function userExists(bytes32 id) external view returns (bool) {
        return users[id].exists;
    }

    function usernameExists(string calldata username) external view returns (bool) {
        return usernameToId[username] != bytes32(0);
    }

    function getUserCount() external view returns (uint256) {
        return userIds.length;
    }

    // ============ Friend Request Functions ============

    function createFriendRequest(
        bytes32 id,
        bytes32 fromUserId,
        bytes32 toUserId
    ) external onlyOwner {
        require(!friendRequests[id].exists, "Friend request already exists");
        require(users[fromUserId].exists, "From user not found");
        require(users[toUserId].exists, "To user not found");
        require(fromUserId != toUserId, "Cannot send request to self");
        require(friendRequestByUsers[fromUserId][toUserId] == bytes32(0), "Request already exists");

        // Check if already friends
        (bytes32 userA, bytes32 userB) = _orderUserIds(fromUserId, toUserId);
        require(friendshipByUsers[userA][userB] == bytes32(0), "Already friends");

        friendRequests[id] = FriendRequest({
            id: id,
            fromUserId: fromUserId,
            toUserId: toUserId,
            status: FriendRequestStatus.Pending,
            createdAt: block.timestamp,
            exists: true
        });

        friendRequestByUsers[fromUserId][toUserId] = id;
        pendingRequestsTo[toUserId].push(id);
        friendRequestIds.push(id);

        emit FriendRequestCreated(id, fromUserId, toUserId);
    }

    function getFriendRequest(bytes32 id) external view returns (
        bytes32 requestId,
        bytes32 fromUserId,
        bytes32 toUserId,
        FriendRequestStatus status,
        uint256 createdAt
    ) {
        FriendRequest storage request = friendRequests[id];
        require(request.exists, "Friend request not found");

        return (
            request.id,
            request.fromUserId,
            request.toUserId,
            request.status,
            request.createdAt
        );
    }

    function updateFriendRequestStatus(bytes32 id, FriendRequestStatus status) external onlyOwner {
        require(friendRequests[id].exists, "Friend request not found");

        friendRequests[id].status = status;

        emit FriendRequestUpdated(id, status);
    }

    function getPendingRequestsForUser(bytes32 userId) external view returns (bytes32[] memory) {
        bytes32[] storage allRequests = pendingRequestsTo[userId];

        // Count pending requests
        uint256 pendingCount = 0;
        for (uint256 i = 0; i < allRequests.length; i++) {
            if (friendRequests[allRequests[i]].status == FriendRequestStatus.Pending) {
                pendingCount++;
            }
        }

        // Build result array
        bytes32[] memory result = new bytes32[](pendingCount);
        uint256 resultIndex = 0;
        for (uint256 i = 0; i < allRequests.length; i++) {
            if (friendRequests[allRequests[i]].status == FriendRequestStatus.Pending) {
                result[resultIndex++] = allRequests[i];
            }
        }

        return result;
    }

    // ============ Friendship Functions ============

    function createFriendship(
        bytes32 id,
        bytes32 userAId,
        bytes32 userBId
    ) external onlyOwner {
        require(!friendships[id].exists, "Friendship already exists");
        require(users[userAId].exists, "User A not found");
        require(users[userBId].exists, "User B not found");

        // Ensure consistent ordering
        (bytes32 orderedA, bytes32 orderedB) = _orderUserIds(userAId, userBId);
        require(friendshipByUsers[orderedA][orderedB] == bytes32(0), "Friendship already exists");

        friendships[id] = Friendship({
            id: id,
            userAId: orderedA,
            userBId: orderedB,
            createdAt: block.timestamp,
            exists: true
        });

        friendshipByUsers[orderedA][orderedB] = id;
        userFriendships[orderedA].push(id);
        userFriendships[orderedB].push(id);
        friendshipIds.push(id);

        emit FriendshipCreated(id, orderedA, orderedB);
    }

    function getFriendship(bytes32 id) external view returns (
        bytes32 friendshipId,
        bytes32 userAId,
        bytes32 userBId,
        uint256 createdAt
    ) {
        Friendship storage friendship = friendships[id];
        require(friendship.exists, "Friendship not found");

        return (
            friendship.id,
            friendship.userAId,
            friendship.userBId,
            friendship.createdAt
        );
    }

    function areFriends(bytes32 userId1, bytes32 userId2) external view returns (bool) {
        (bytes32 userA, bytes32 userB) = _orderUserIds(userId1, userId2);
        return friendshipByUsers[userA][userB] != bytes32(0);
    }

    function getUserFriendships(bytes32 userId) external view returns (bytes32[] memory) {
        return userFriendships[userId];
    }

    function getFriendsOfUser(bytes32 userId) external view returns (bytes32[] memory friendIds) {
        bytes32[] storage fships = userFriendships[userId];
        friendIds = new bytes32[](fships.length);

        for (uint256 i = 0; i < fships.length; i++) {
            Friendship storage f = friendships[fships[i]];
            if (f.userAId == userId) {
                friendIds[i] = f.userBId;
            } else {
                friendIds[i] = f.userAId;
            }
        }
    }

    // ============ Message Functions ============

    function createMessage(
        bytes32 id,
        bytes32 fromUserId,
        bytes32 toUserId,
        bytes calldata encryptedContent,
        ContentType contentType,
        string calldata signature
    ) external onlyOwner {
        require(!messages[id].exists, "Message already exists");
        require(users[fromUserId].exists, "From user not found");
        require(users[toUserId].exists, "To user not found");

        messages[id] = Message({
            id: id,
            fromUserId: fromUserId,
            toUserId: toUserId,
            encryptedContent: encryptedContent,
            contentType: contentType,
            signature: signature,
            createdAt: block.timestamp,
            exists: true
        });

        messagesToUser[toUserId].push(id);
        messagesFromUser[fromUserId].push(id);
        messageIds.push(id);

        emit MessageCreated(id, fromUserId, toUserId);
    }

    function getMessage(bytes32 id) external view returns (
        bytes32 messageId,
        bytes32 fromUserId,
        bytes32 toUserId,
        bytes memory encryptedContent,
        ContentType contentType,
        string memory signature,
        uint256 createdAt
    ) {
        Message storage message = messages[id];
        require(message.exists, "Message not found");

        return (
            message.id,
            message.fromUserId,
            message.toUserId,
            message.encryptedContent,
            message.contentType,
            message.signature,
            message.createdAt
        );
    }

    function deleteMessage(bytes32 id) external onlyOwner {
        require(messages[id].exists, "Message not found");

        messages[id].exists = false;

        emit MessageDeleted(id);
    }

    function getMessagesForUser(bytes32 userId) external view returns (bytes32[] memory) {
        bytes32[] storage allMessages = messagesToUser[userId];

        // Count existing messages
        uint256 count = 0;
        for (uint256 i = 0; i < allMessages.length; i++) {
            if (messages[allMessages[i]].exists) {
                count++;
            }
        }

        // Build result array
        bytes32[] memory result = new bytes32[](count);
        uint256 resultIndex = 0;
        for (uint256 i = 0; i < allMessages.length; i++) {
            if (messages[allMessages[i]].exists) {
                result[resultIndex++] = allMessages[i];
            }
        }

        return result;
    }

    function getMessagesSentByUser(bytes32 userId) external view returns (bytes32[] memory) {
        bytes32[] storage allMessages = messagesFromUser[userId];

        // Count existing messages
        uint256 count = 0;
        for (uint256 i = 0; i < allMessages.length; i++) {
            if (messages[allMessages[i]].exists) {
                count++;
            }
        }

        // Build result array
        bytes32[] memory result = new bytes32[](count);
        uint256 resultIndex = 0;
        for (uint256 i = 0; i < allMessages.length; i++) {
            if (messages[allMessages[i]].exists) {
                result[resultIndex++] = allMessages[i];
            }
        }

        return result;
    }

    // ============ Internal Functions ============

    function _orderUserIds(bytes32 id1, bytes32 id2) internal pure returns (bytes32, bytes32) {
        if (id1 < id2) {
            return (id1, id2);
        }
        return (id2, id1);
    }

    // ============ Admin Functions ============

    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "Invalid address");
        owner = newOwner;
    }
}
