// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script, console2} from "forge-std/Script.sol";
import {PaymentChannel} from "../src/core/PaymentChannel.sol";
import {CommitRevealEntropy} from "../src/randomness/CommitRevealEntropy.sol";
import {SettlementEngine} from "../src/core/SettlementEngine.sol";
import {IEntropySource} from "../src/interfaces/IEntropySource.sol";

contract Step2_Settle is Script {
    function run() external {
        uint256 providerPk = 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d;
        uint256 clientPk   = 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a;

        address provider = vm.addr(providerPk);
        address client   = vm.addr(clientPk);

        PaymentChannel channel = PaymentChannel(0x5FbDB2315678afecb367f032d93F642f64180aa3); // Will be read or set
        // In local deterministic deployment, the addresses are:
        // ProtocolConfig: 0x5FbDB2315678afecb367f032d93F642f64180aa3
        // ProviderRegistry: 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
        // PaymentChannel: 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0
        // CommitRevealEntropy: 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9
        // SettlementEngine: 0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9
        channel = PaymentChannel(0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0);
        CommitRevealEntropy entropy = CommitRevealEntropy(0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9);
        SettlementEngine engine = SettlementEngine(0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9);

        bytes32 secret = keccak256("CIPHER_PROVIDER_SECRET_DEMO");
        uint256 roundId = 1;

        bytes32 channelKey = channel.getChannelKey(client, provider);
        bytes32 roundKey = keccak256(abi.encodePacked(channelKey, roundId));
        (uint256 tau, uint256 faceValue, , uint256 commitBlock, ) = engine.rounds(roundKey);

        console2.log("=== STEP 5: ON-CHAIN RAFFLE WINNER CALCULATION ===");
        (uint256 winnerIndex, bytes32 targetBh) = entropy.computeWinnerIndex(commitBlock, secret, roundId, tau);

        console2.log("  [+] Target Blockhash:   ", vm.toString(targetBh));
        console2.log("  [+] Winning Chunk Index:", winnerIndex);

        // Client signs the winning ticket
        IEntropySource.RoundTicket memory ticket = IEntropySource.RoundTicket({
            sender: client,
            recipient: provider,
            roundId: roundId,
            localIndex: winnerIndex,
            faceValue: faceValue,
            winProb: type(uint256).max,
            senderNonce: 1
        });

        bytes32 tHash = entropy.roundTicketHash(ticket);
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(clientPk, tHash);
        bytes memory sig = abi.encodePacked(r, s, v);

        console2.log("\n=== STEP 6: EXECUTING ON-CHAIN SETTLEMENT ===");
        console2.log("  [+] Provider Bal Before:", provider.balance / 1e18, "ETH");
        console2.log("  [+] Channel Bal Before: ", channel.channelBalance(client, provider) / 1e18, "ETH");

        uint256 balBefore = provider.balance;

        vm.startBroadcast(providerPk);
        engine.settleRound(ticket, sig, secret);
        vm.stopBroadcast();

        uint256 balAfter = provider.balance;

        console2.log("==================================================");
        console2.log("SUCCESS! ROUND SETTLED ON LIVE ANVIL EVM");
        console2.log("==================================================");
        console2.log("  [+] Provider Bal After: ", balAfter / 1e18, "ETH");
        console2.log("  [+] Channel Bal After:  ", channel.channelBalance(client, provider) / 1e18, "ETH");
        console2.log("  [+] NET PAYOUT EARNED:  ", (balAfter - balBefore) / 1e18, "ETH");
        console2.log("==================================================");
    }
}
