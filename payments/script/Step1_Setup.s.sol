// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script, console2} from "forge-std/Script.sol";
import {ProtocolConfig} from "../src/config/ProtocolConfig.sol";
import {ProviderRegistry} from "../src/core/ProviderRegistry.sol";
import {PaymentChannel} from "../src/core/PaymentChannel.sol";
import {CommitRevealEntropy} from "../src/randomness/CommitRevealEntropy.sol";
import {SettlementEngine} from "../src/core/SettlementEngine.sol";

contract Step1_Setup is Script {
    function run() external {
        uint256 deployerPk = 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80;
        uint256 providerPk = 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d;
        uint256 clientPk   = 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a;

        address deployer = vm.addr(deployerPk);
        address provider = vm.addr(providerPk);
        address client   = vm.addr(clientPk);

        console2.log("=== STEP 1: DEPLOYING CONTRACTS ON ANVIL ===");
        vm.startBroadcast(deployerPk);

        ProtocolConfig config = new ProtocolConfig(deployer);
        ProviderRegistry registry = new ProviderRegistry(address(config));
        PaymentChannel channel = new PaymentChannel(address(config));
        CommitRevealEntropy entropy = new CommitRevealEntropy(address(config));
        SettlementEngine engine = new SettlementEngine(
            address(config),
            address(entropy),
            address(channel),
            address(registry)
        );

        channel.setPayoutEngine(address(engine), true);
        vm.stopBroadcast();

        console2.log("  [+] ProtocolConfig:     ", address(config));
        console2.log("  [+] ProviderRegistry:   ", address(registry));
        console2.log("  [+] PaymentChannel:     ", address(channel));
        console2.log("  [+] CommitRevealEntropy:", address(entropy));
        console2.log("  [+] SettlementEngine:   ", address(engine));

        console2.log("\n=== STEP 2: PROVIDER STAKING COLLATERAL ===");
        vm.startBroadcast(providerPk);
        registry.registerProvider{value: 2 ether}();
        vm.stopBroadcast();
        console2.log("  [+] Provider Staked 2.0 ETH");
        console2.log("  [+] Provider Active:     ", registry.isProviderActive(provider));

        console2.log("\n=== STEP 3: CLIENT FUNDING ESCROW CHANNEL ===");
        vm.startBroadcast(clientPk);
        channel.openChannel{value: 5 ether}(provider);
        vm.stopBroadcast();
        console2.log("  [+] Client Deposited 5.0 ETH into PaymentChannel");
        console2.log("  [+] Channel Balance:     ", channel.channelBalance(client, provider) / 1e18, "ETH");

        console2.log("\n=== STEP 4: PROVIDER COMMITTING ROUND 1 ===");
        bytes32 secret = keccak256("CIPHER_PROVIDER_SECRET_DEMO");
        bytes32 secretHash = keccak256(abi.encodePacked(secret));
        uint256 roundId = 1;
        uint256 tau = 10;
        uint256 faceValue = 1 ether;

        vm.startBroadcast(providerPk);
        engine.commitRoundSecret(client, roundId, tau, faceValue, secretHash);
        vm.stopBroadcast();

        console2.log("  [+] Provider Committed Secret for Round 1");
        console2.log("  [+] Commit Block Number: ", block.number);
        console2.log("  [+] Confirmation Delay:  ", config.CONFIRMATION_DELAY(), "blocks");
    }
}
