// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {QuickPicStorage} from "../src/QuickPicStorage.sol";

contract DeployScript is Script {
    function setUp() public {}

    function run() public {
        // Use the first Anvil default private key
        // Address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
        uint256 deployerPrivateKey = 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80;

        vm.startBroadcast(deployerPrivateKey);

        QuickPicStorage storage_ = new QuickPicStorage();
        console.log("QuickPicStorage deployed at:", address(storage_));

        vm.stopBroadcast();
    }
}
