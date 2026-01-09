// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {QuickPicStorage} from "../src/QuickPicStorage.sol";

contract QueryMessagesScript is Script {
    function run() public view {
        // Replace with your deployed contract address
        QuickPicStorage store = QuickPicStorage(0x5FbDB2315678afecb367f032d93F642f64180aa3);

        // Get maddie's user info by username
        (
            bytes32 userId,
            uint256 userNumber,
            string memory username,
            ,  // passwordHash (skip)
            string memory publicKey,
            uint256 createdAt,
            uint256 updatedAt
        ) = store.getUserByUsername("maddie");

        console.log("=== User Info ===");
        console.log("Username:", username);
        console.log("User Number:", userNumber);
        console.logBytes32(userId);

        // Get messages sent BY maddie
        bytes32[] memory sentMessages = store.getMessagesSentByUser(userId);
        console.log("\n=== Messages Sent by Maddie ===");
        console.log("Total messages:", sentMessages.length);

        for (uint i = 0; i < sentMessages.length; i++) {
            (
                bytes32 msgId,
                bytes32 fromUserId,
                bytes32 toUserId,
                bytes memory encryptedContent,
                QuickPicStorage.ContentType contentType,
                string memory signature,
                uint256 msgCreatedAt
            ) = store.getMessage(sentMessages[i]);

            console.log("\n--- Message", i + 1, "---");
            console.log("Message ID:");
            console.logBytes32(msgId);
            console.log("Content Type:", uint(contentType)); // 0=Text, 1=Image
            console.log("Created At:", msgCreatedAt);
            console.log("Encrypted Content Length:", encryptedContent.length);
        }

        // Get messages sent TO maddie
        bytes32[] memory receivedMessages = store.getMessagesForUser(userId);
        console.log("\n=== Messages Received by Maddie ===");
        console.log("Total messages:", receivedMessages.length);
    }
}
