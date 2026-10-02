package payments

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"cipher/network/payments/bindings/commitreveal"
	"cipher/network/payments/bindings/paymentchannel"
	"cipher/network/payments/bindings/providerregistry"
	"cipher/network/payments/bindings/settlementengine"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ContractAddresses holds the deployed contract addresses.
type ContractAddresses struct {
	PaymentChannel   common.Address
	SettlementEngine common.Address
	EntropySource    common.Address
	ProviderRegistry common.Address
}

// PaymentClient coordinates Ethereum RPC interactions, channel balances, and on-chain settlement.
type PaymentClient struct {
	Client     *ethclient.Client
	ChainID    *big.Int
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
	Signer     *TicketSigner

	PaymentChannel   *paymentchannel.PaymentChannel
	SettlementEngine *settlementengine.SettlementEngine
	EntropySource    *commitreveal.CommitRevealEntropy
	ProviderRegistry *providerregistry.ProviderRegistry
}

// NewPaymentClient dials an Ethereum JSON-RPC endpoint and binds to deployed contracts.
func NewPaymentClient(
	ctx context.Context,
	rpcURL string,
	privKeyHex string,
	addrs ContractAddresses,
) (*PaymentClient, error) {
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ethereum rpc: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chain id: %w", err)
	}

	var privKey *ecdsa.PrivateKey
	var addr common.Address
	if privKeyHex != "" {
		if len(privKeyHex) >= 2 && privKeyHex[:2] == "0x" {
			privKeyHex = privKeyHex[2:]
		}
		privKey, err = crypto.HexToECDSA(privKeyHex)
		if err != nil {
			return nil, fmt.Errorf("invalid private key hex: %w", err)
		}
		addr = crypto.PubkeyToAddress(privKey.PublicKey)
	}

	signer := NewTicketSigner(privKey, chainID, addrs.EntropySource)

	var pc *paymentchannel.PaymentChannel
	if addrs.PaymentChannel != (common.Address{}) {
		pc, err = paymentchannel.NewPaymentChannel(addrs.PaymentChannel, client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind PaymentChannel: %w", err)
		}
	}

	var se *settlementengine.SettlementEngine
	if addrs.SettlementEngine != (common.Address{}) {
		se, err = settlementengine.NewSettlementEngine(addrs.SettlementEngine, client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind SettlementEngine: %w", err)
		}
	}

	var es *commitreveal.CommitRevealEntropy
	if addrs.EntropySource != (common.Address{}) {
		es, err = commitreveal.NewCommitRevealEntropy(addrs.EntropySource, client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind CommitRevealEntropy: %w", err)
		}
	}

	var pr *providerregistry.ProviderRegistry
	if addrs.ProviderRegistry != (common.Address{}) {
		pr, err = providerregistry.NewProviderRegistry(addrs.ProviderRegistry, client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind ProviderRegistry: %w", err)
		}
	}

	return &PaymentClient{
		Client:           client,
		ChainID:          chainID,
		PrivateKey:       privKey,
		Address:          addr,
		Signer:           signer,
		PaymentChannel:   pc,
		SettlementEngine: se,
		EntropySource:    es,
		ProviderRegistry: pr,
	}, nil
}

// Close closes the underlying JSON-RPC client connection.
func (c *PaymentClient) Close() {
	if c.Client != nil {
		c.Client.Close()
	}
}

// GetTransactor creates a new keyed transaction transactor with latest suggested gas price.
func (c *PaymentClient) GetTransactor(ctx context.Context, value *big.Int) (*bind.TransactOpts, error) {
	if c.PrivateKey == nil {
		return nil, fmt.Errorf("cannot create transactor: private key is nil")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.PrivateKey, c.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	auth.Context = ctx
	if value != nil {
		auth.Value = value
	}

	gasPrice, err := c.Client.SuggestGasPrice(ctx)
	if err == nil {
		auth.GasPrice = gasPrice
	}

	return auth, nil
}

// GetChannelBalance queries the available balance in the payment channel between sender and recipient.
func (c *PaymentClient) GetChannelBalance(ctx context.Context, sender, recipient common.Address) (*big.Int, error) {
	if c.PaymentChannel == nil {
		return nil, fmt.Errorf("PaymentChannel contract not bound")
	}

	opts := &bind.CallOpts{Context: ctx}
	bal, err := c.PaymentChannel.ChannelBalance(opts, sender, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to query channel balance: %w", err)
	}

	return bal, nil
}

// GetProviderStake queries the active stake for a provider from ProviderRegistry.
func (c *PaymentClient) GetProviderStake(ctx context.Context, provider common.Address) (*big.Int, error) {
	if c.ProviderRegistry == nil {
		return nil, fmt.Errorf("ProviderRegistry contract not bound")
	}

	opts := &bind.CallOpts{Context: ctx}
	pInfo, err := c.ProviderRegistry.Providers(opts, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to query provider: %w", err)
	}

	return pInfo.Stake, nil
}

// SettleRound calls SettlementEngine.settleRound on-chain using a signed winning ticket and revealed secret.
func (c *PaymentClient) SettleRound(
	ctx context.Context,
	ticket RoundTicket,
	senderSig []byte,
	secret [32]byte,
) (*types.Transaction, error) {
	if c.SettlementEngine == nil {
		return nil, fmt.Errorf("SettlementEngine contract not bound")
	}

	auth, err := c.GetTransactor(ctx, nil)
	if err != nil {
		return nil, err
	}

	solTicket := settlementengine.IEntropySourceRoundTicket{
		Sender:      ticket.Sender,
		Recipient:   ticket.Recipient,
		RoundId:     ticket.RoundID,
		LocalIndex:  ticket.LocalIndex,
		FaceValue:   ticket.FaceValue,
		WinProb:     ticket.WinProb,
		SenderNonce: ticket.SenderNonce,
	}

	tx, err := c.SettlementEngine.SettleRound(auth, solTicket, senderSig, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to submit settleRound transaction: %w", err)
	}

	return tx, nil
}
