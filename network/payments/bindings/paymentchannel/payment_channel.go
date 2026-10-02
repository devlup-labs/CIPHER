// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package paymentchannel

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PaymentChannelMetaData contains all meta data concerning the PaymentChannel contract.
var PaymentChannelMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_config\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizedPayoutEngine\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"channelBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"channels\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"totalDeposited\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"settledAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"withdrawnAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unlockBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"config\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractProtocolConfig\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositToChannel\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getChannelKey\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"initiateChannelUnlock\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isUnlockInitiated\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"openChannel\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"payoutProvider\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPayoutEngine\",\"inputs\":[{\"name\":\"engine\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authorized\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawChannel\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ChannelDeposited\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChannelOpened\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChannelPayoutExecuted\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChannelUnlockInitiated\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"unlockBlock\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChannelWithdrawn\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientChannelBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NothingToWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedCaller\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnlockAlreadyInitiated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnlockNotComplete\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnlockNotInitiated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroDeposit\",\"inputs\":[]}]",
	Bin: "0x60a03461010057601f610afe38819003918201601f19168301916001600160401b038311848410176101045780849260209460405283398101031261010057516001600160a01b038116908190036101005760017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580156100af57608052600380546001600160a01b031916331790556040516109e59081610119823960805181818161012101526104530152f35b60405162461bcd60e51b815260206004820152602360248201527f5061796d656e744368616e6e656c3a205a65726f20636f6e666967206164647260448201526265737360e81b6064820152608490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe6080806040526004361015610012575f80fd5b5f3560e01c90816302d5a975146107e9575080630f232621146107b257806328095f26146106a15780632abb5e9d146106215780633f4ba83a146105b35780635c975abb14610592578063756b7b891461048257806379502c551461043e5780637a7ebd7b146103f25780638456cb59146103915780638da5cb5b146103695780638db31f881461034b578063a6db5b1d146102bc578063a71c2df314610250578063a89e60c6146101f35763df087864146100cc575f80fd5b346101cd5760203660031901126101cd576100e5610822565b6100ed610975565b6100f781336108f1565b5f526001602052600360405f20019081546101e457604051633aa2c67960e11b81526020816004817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa80156101d9575f906101a1575b61016391504361086f565b918290556040519182526001600160a01b03169033907fcce983f8bacd60182acc7e585fbe219277117ceaf0871600204eae641c59838a90602090a3005b506020813d6020116101d1575b816101bb6020938361087c565b810103126101cd576101639051610158565b5f80fd5b3d91506101ae565b6040513d5f823e3d90fd5b6305d89a4360e11b5f5260045ffd5b346101cd5760403660031901126101cd5761021d61020f610822565b610217610838565b906108f1565b5f526001602052602061024860405f20600261023f825460018401549061084e565b9101549061084e565b604051908152f35b346101cd5760403660031901126101cd57610269610822565b602435908115158092036101cd576003546001600160a01b031633036102ad5760018060a01b03165f52600260205260405f209060ff801983541691161790555f80f35b635c427cd960e01b5f5260045ffd5b60203660031901126101cd576102d0610822565b6102d8610975565b341561033c576102e881336108f1565b5f52600160205260405f206102fe34825461086f565b90556040513481526001600160a01b039091169033907f05c20c1748de9db22bf7db8be0f4d51c217ee76839ad9ba6aabc1d5a82e953a190602090a3005b6356316e8760e01b5f5260045ffd5b346101cd5760403660031901126101cd57602061024861020f610822565b346101cd575f3660031901126101cd576003546040516001600160a01b039091168152602090f35b346101cd575f3660031901126101cd576003546001600160a01b031633036102ad576103bb610975565b600160ff195f5416175f557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a1005b346101cd5760203660031901126101cd576004355f526001602052608060405f208054906001810154906003600282015491015491604051938452602084015260408301526060820152f35b346101cd575f3660031901126101cd576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b346101cd5760603660031901126101cd5761049b610822565b6104a3610838565b60443590335f52600260205260ff60405f205416156102ad576104c461093d565b6104cc610975565b6104d681846108f1565b5f52600160205260405f2080546104f96001830192600261023f8554809561084e565b841161058357836105099161086f565b90556001600160a01b0316915f80808085875af16105256108b2565b5015610574576040519182526001600160a01b0316907fef4ac45c90db8ad23280036f26d8dc1584377d3a7fbcf6b3068baf7bd80b909090602090a360015f5160206109905f395f51905f5255005b6312171d8360e31b5f5260045ffd5b632c51d8db60e21b5f5260045ffd5b346101cd575f3660031901126101cd57602060ff5f54166040519015158152f35b346101cd575f3660031901126101cd576003546001600160a01b031633036102ad575f5460ff8116156106125760ff19165f557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a1005b638dfc202b60e01b5f5260045ffd5b60203660031901126101cd57610635610822565b61063d610975565b341561033c5761064d81336108f1565b5f52600160205260405f2061066334825461086f565b90556040513481526001600160a01b039091169033907f434d899467f051cde9a358523702f186412ffe4e86423150a08680ccdb1519ee90602090a3005b346101cd5760203660031901126101cd576106ba610822565b6106c261093d565b6106ca610975565b6106d481336108f1565b5f52600160205260405f2090600382015480156107a3574310610794576002610703835460018501549061084e565b92016107118154809461084e565b92831561078557836107229161086f565b90555f80808085335af16107346108b2565b5015610574576040519182526001600160a01b03169033907ffed549f521421ee5fb337d4e385525b833bbe3e8a060b62a0f25fb5bcb6f08c190602090a360015f5160206109905f395f51905f5255005b630686827b60e51b5f5260045ffd5b6334de422560e11b5f5260045ffd5b635b89f32f60e11b5f5260045ffd5b346101cd5760403660031901126101cd576107ce61020f610822565b5f5260016020526020600360405f2001541515604051908152f35b346101cd5760203660031901126101cd576020906001600160a01b0361080d610822565b165f526002825260ff60405f20541615158152f35b600435906001600160a01b03821682036101cd57565b602435906001600160a01b03821682036101cd57565b9190820391821161085b57565b634e487b7160e01b5f52601160045260245ffd5b9190820180921161085b57565b90601f8019910116810190811067ffffffffffffffff82111761089e57604052565b634e487b7160e01b5f52604160045260245ffd5b3d156108ec573d9067ffffffffffffffff821161089e57604051916108e1601f8201601f19166020018461087c565b82523d5f602084013e565b606090565b906040519060208201926bffffffffffffffffffffffff199060601b1683526bffffffffffffffffffffffff199060601b1660348201526028815261093760488261087c565b51902090565b60025f5160206109905f395f51905f5254146109665760025f5160206109905f395f51905f5255565b633ee5aeb560e01b5f5260045ffd5b60ff5f541661098057565b63d93c066560e01b5f5260045ffdfe9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00a264697066735822122063f6fe9546e80edd3cba2a1ad78f592042a712806398b378a5bdca5ba6bc003664736f6c63430008230033",
}

// PaymentChannelABI is the input ABI used to generate the binding from.
// Deprecated: Use PaymentChannelMetaData.ABI instead.
var PaymentChannelABI = PaymentChannelMetaData.ABI

// PaymentChannelBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PaymentChannelMetaData.Bin instead.
var PaymentChannelBin = PaymentChannelMetaData.Bin

// DeployPaymentChannel deploys a new Ethereum contract, binding an instance of PaymentChannel to it.
func DeployPaymentChannel(auth *bind.TransactOpts, backend bind.ContractBackend, _config common.Address) (common.Address, *types.Transaction, *PaymentChannel, error) {
	parsed, err := PaymentChannelMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PaymentChannelBin), backend, _config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PaymentChannel{PaymentChannelCaller: PaymentChannelCaller{contract: contract}, PaymentChannelTransactor: PaymentChannelTransactor{contract: contract}, PaymentChannelFilterer: PaymentChannelFilterer{contract: contract}}, nil
}

// PaymentChannel is an auto generated Go binding around an Ethereum contract.
type PaymentChannel struct {
	PaymentChannelCaller     // Read-only binding to the contract
	PaymentChannelTransactor // Write-only binding to the contract
	PaymentChannelFilterer   // Log filterer for contract events
}

// PaymentChannelCaller is an auto generated read-only Go binding around an Ethereum contract.
type PaymentChannelCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentChannelTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PaymentChannelTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentChannelFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PaymentChannelFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentChannelSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PaymentChannelSession struct {
	Contract     *PaymentChannel   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PaymentChannelCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PaymentChannelCallerSession struct {
	Contract *PaymentChannelCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// PaymentChannelTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PaymentChannelTransactorSession struct {
	Contract     *PaymentChannelTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// PaymentChannelRaw is an auto generated low-level Go binding around an Ethereum contract.
type PaymentChannelRaw struct {
	Contract *PaymentChannel // Generic contract binding to access the raw methods on
}

// PaymentChannelCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PaymentChannelCallerRaw struct {
	Contract *PaymentChannelCaller // Generic read-only contract binding to access the raw methods on
}

// PaymentChannelTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PaymentChannelTransactorRaw struct {
	Contract *PaymentChannelTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPaymentChannel creates a new instance of PaymentChannel, bound to a specific deployed contract.
func NewPaymentChannel(address common.Address, backend bind.ContractBackend) (*PaymentChannel, error) {
	contract, err := bindPaymentChannel(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PaymentChannel{PaymentChannelCaller: PaymentChannelCaller{contract: contract}, PaymentChannelTransactor: PaymentChannelTransactor{contract: contract}, PaymentChannelFilterer: PaymentChannelFilterer{contract: contract}}, nil
}

// NewPaymentChannelCaller creates a new read-only instance of PaymentChannel, bound to a specific deployed contract.
func NewPaymentChannelCaller(address common.Address, caller bind.ContractCaller) (*PaymentChannelCaller, error) {
	contract, err := bindPaymentChannel(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PaymentChannelCaller{contract: contract}, nil
}

// NewPaymentChannelTransactor creates a new write-only instance of PaymentChannel, bound to a specific deployed contract.
func NewPaymentChannelTransactor(address common.Address, transactor bind.ContractTransactor) (*PaymentChannelTransactor, error) {
	contract, err := bindPaymentChannel(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PaymentChannelTransactor{contract: contract}, nil
}

// NewPaymentChannelFilterer creates a new log filterer instance of PaymentChannel, bound to a specific deployed contract.
func NewPaymentChannelFilterer(address common.Address, filterer bind.ContractFilterer) (*PaymentChannelFilterer, error) {
	contract, err := bindPaymentChannel(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PaymentChannelFilterer{contract: contract}, nil
}

// bindPaymentChannel binds a generic wrapper to an already deployed contract.
func bindPaymentChannel(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PaymentChannelMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PaymentChannel *PaymentChannelRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PaymentChannel.Contract.PaymentChannelCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PaymentChannel *PaymentChannelRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PaymentChannel.Contract.PaymentChannelTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PaymentChannel *PaymentChannelRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PaymentChannel.Contract.PaymentChannelTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PaymentChannel *PaymentChannelCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PaymentChannel.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PaymentChannel *PaymentChannelTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PaymentChannel.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PaymentChannel *PaymentChannelTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PaymentChannel.Contract.contract.Transact(opts, method, params...)
}

// AuthorizedPayoutEngine is a free data retrieval call binding the contract method 0x02d5a975.
//
// Solidity: function authorizedPayoutEngine(address ) view returns(bool)
func (_PaymentChannel *PaymentChannelCaller) AuthorizedPayoutEngine(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _PaymentChannel.contract.Call(opts, &out, "authorizedPayoutEngine", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizedPayoutEngine is a free data retrieval call binding the contract method 0x02d5a975.
//
// Solidity: function authorizedPayoutEngine(address ) view returns(bool)
func (_PaymentChannel *PaymentChannelSession) AuthorizedPayoutEngine(arg0 common.Address) (bool, error) {
	return _PaymentChannel.Contract.AuthorizedPayoutEngine(&_PaymentChannel.CallOpts, arg0)
}

// AuthorizedPayoutEngine is a free data retrieval call binding the contract method 0x02d5a975.
//
// Solidity: function authorizedPayoutEngine(address ) view returns(bool)
func (_PaymentChannel *PaymentChannelCallerSession) AuthorizedPayoutEngine(arg0 common.Address) (bool, error) {
	return _PaymentChannel.Contract.AuthorizedPayoutEngine(&_PaymentChannel.CallOpts, arg0)
}

// ChannelBalance is a free data retrieval call binding the contract method 0xa89e60c6.
//
// Solidity: function channelBalance(address sender, address recipient) view returns(uint256)
func (_PaymentChannel *PaymentChannelCaller) ChannelBalance(opts *bind.CallOpts, sender common.Address, recipient common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PaymentChannel.contract.Call(opts, &out, "channelBalance", sender, recipient)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChannelBalance is a free data retrieval call binding the contract method 0xa89e60c6.
//
// Solidity: function channelBalance(address sender, address recipient) view returns(uint256)
func (_PaymentChannel *PaymentChannelSession) ChannelBalance(sender common.Address, recipient common.Address) (*big.Int, error) {
	return _PaymentChannel.Contract.ChannelBalance(&_PaymentChannel.CallOpts, sender, recipient)
}

// ChannelBalance is a free data retrieval call binding the contract method 0xa89e60c6.
//
// Solidity: function channelBalance(address sender, address recipient) view returns(uint256)
func (_PaymentChannel *PaymentChannelCallerSession) ChannelBalance(sender common.Address, recipient common.Address) (*big.Int, error) {
	return _PaymentChannel.Contract.ChannelBalance(&_PaymentChannel.CallOpts, sender, recipient)
}

// Channels is a free data retrieval call binding the contract method 0x7a7ebd7b.
//
// Solidity: function channels(bytes32 ) view returns(uint256 totalDeposited, uint256 settledAmount, uint256 withdrawnAmount, uint256 unlockBlock)
func (_PaymentChannel *PaymentChannelCaller) Channels(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TotalDeposited  *big.Int
	SettledAmount   *big.Int
	WithdrawnAmount *big.Int
	UnlockBlock     *big.Int
}, error) {
	var out []interface{}
	err := _PaymentChannel.contract.Call(opts, &out, "channels", arg0)

	outstruct := new(struct {
		TotalDeposited  *big.Int
		SettledAmount   *big.Int
		WithdrawnAmount *big.Int
		UnlockBlock     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalDeposited = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SettledAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.WithdrawnAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UnlockBlock = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Channels is a free data retrieval call binding the contract method 0x7a7ebd7b.
//
// Solidity: function channels(bytes32 ) view returns(uint256 totalDeposited, uint256 settledAmount, uint256 withdrawnAmount, uint256 unlockBlock)
func (_PaymentChannel *PaymentChannelSession) Channels(arg0 [32]byte) (struct {
	TotalDeposited  *big.Int
	SettledAmount   *big.Int
	WithdrawnAmount *big.Int
	UnlockBlock     *big.Int
}, error) {
	return _PaymentChannel.Contract.Channels(&_PaymentChannel.CallOpts, arg0)
}

// Channels is a free data retrieval call binding the contract method 0x7a7ebd7b.
//
// Solidity: function channels(bytes32 ) view returns(uint256 totalDeposited, uint256 settledAmount, uint256 withdrawnAmount, uint256 unlockBlock)
func (_PaymentChannel *PaymentChannelCallerSession) Channels(arg0 [32]byte) (struct {
	TotalDeposited  *big.Int
	SettledAmount   *big.Int
	WithdrawnAmount *big.Int
	UnlockBlock     *big.Int
}, error) {
	return _PaymentChannel.Contract.Channels(&_PaymentChannel.CallOpts, arg0)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_PaymentChannel *PaymentChannelCaller) Config(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PaymentChannel.contract.Call(opts, &out, "config")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_PaymentChannel *PaymentChannelSession) Config() (common.Address, error) {
	return _PaymentChannel.Contract.Config(&_PaymentChannel.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_PaymentChannel *PaymentChannelCallerSession) Config() (common.Address, error) {
	return _PaymentChannel.Contract.Config(&_PaymentChannel.CallOpts)
}

// GetChannelKey is a free data retrieval call binding the contract method 0x8db31f88.
//
// Solidity: function getChannelKey(address sender, address recipient) pure returns(bytes32)
func (_PaymentChannel *PaymentChannelCaller) GetChannelKey(opts *bind.CallOpts, sender common.Address, recipient common.Address) ([32]byte, error) {
	var out []interface{}
	err := _PaymentChannel.contract.Call(opts, &out, "getChannelKey", sender, recipient)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetChannelKey is a free data retrieval call binding the contract method 0x8db31f88.
//
// Solidity: function getChannelKey(address sender, address recipient) pure returns(bytes32)
func (_PaymentChannel *PaymentChannelSession) GetChannelKey(sender common.Address, recipient common.Address) ([32]byte, error) {
	return _PaymentChannel.Contract.GetChannelKey(&_PaymentChannel.CallOpts, sender, recipient)
}

// GetChannelKey is a free data retrieval call binding the contract method 0x8db31f88.
//
// Solidity: function getChannelKey(address sender, address recipient) pure returns(bytes32)
func (_PaymentChannel *PaymentChannelCallerSession) GetChannelKey(sender common.Address, recipient common.Address) ([32]byte, error) {
	return _PaymentChannel.Contract.GetChannelKey(&_PaymentChannel.CallOpts, sender, recipient)
}

// IsUnlockInitiated is a free data retrieval call binding the contract method 0x0f232621.
//
// Solidity: function isUnlockInitiated(address sender, address recipient) view returns(bool)
func (_PaymentChannel *PaymentChannelCaller) IsUnlockInitiated(opts *bind.CallOpts, sender common.Address, recipient common.Address) (bool, error) {
	var out []interface{}
	err := _PaymentChannel.contract.Call(opts, &out, "isUnlockInitiated", sender, recipient)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsUnlockInitiated is a free data retrieval call binding the contract method 0x0f232621.
//
// Solidity: function isUnlockInitiated(address sender, address recipient) view returns(bool)
func (_PaymentChannel *PaymentChannelSession) IsUnlockInitiated(sender common.Address, recipient common.Address) (bool, error) {
	return _PaymentChannel.Contract.IsUnlockInitiated(&_PaymentChannel.CallOpts, sender, recipient)
}

// IsUnlockInitiated is a free data retrieval call binding the contract method 0x0f232621.
//
// Solidity: function isUnlockInitiated(address sender, address recipient) view returns(bool)
func (_PaymentChannel *PaymentChannelCallerSession) IsUnlockInitiated(sender common.Address, recipient common.Address) (bool, error) {
	return _PaymentChannel.Contract.IsUnlockInitiated(&_PaymentChannel.CallOpts, sender, recipient)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PaymentChannel *PaymentChannelCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PaymentChannel.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PaymentChannel *PaymentChannelSession) Owner() (common.Address, error) {
	return _PaymentChannel.Contract.Owner(&_PaymentChannel.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PaymentChannel *PaymentChannelCallerSession) Owner() (common.Address, error) {
	return _PaymentChannel.Contract.Owner(&_PaymentChannel.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_PaymentChannel *PaymentChannelCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PaymentChannel.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_PaymentChannel *PaymentChannelSession) Paused() (bool, error) {
	return _PaymentChannel.Contract.Paused(&_PaymentChannel.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_PaymentChannel *PaymentChannelCallerSession) Paused() (bool, error) {
	return _PaymentChannel.Contract.Paused(&_PaymentChannel.CallOpts)
}

// DepositToChannel is a paid mutator transaction binding the contract method 0xa6db5b1d.
//
// Solidity: function depositToChannel(address recipient) payable returns()
func (_PaymentChannel *PaymentChannelTransactor) DepositToChannel(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.contract.Transact(opts, "depositToChannel", recipient)
}

// DepositToChannel is a paid mutator transaction binding the contract method 0xa6db5b1d.
//
// Solidity: function depositToChannel(address recipient) payable returns()
func (_PaymentChannel *PaymentChannelSession) DepositToChannel(recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.Contract.DepositToChannel(&_PaymentChannel.TransactOpts, recipient)
}

// DepositToChannel is a paid mutator transaction binding the contract method 0xa6db5b1d.
//
// Solidity: function depositToChannel(address recipient) payable returns()
func (_PaymentChannel *PaymentChannelTransactorSession) DepositToChannel(recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.Contract.DepositToChannel(&_PaymentChannel.TransactOpts, recipient)
}

// InitiateChannelUnlock is a paid mutator transaction binding the contract method 0xdf087864.
//
// Solidity: function initiateChannelUnlock(address recipient) returns()
func (_PaymentChannel *PaymentChannelTransactor) InitiateChannelUnlock(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.contract.Transact(opts, "initiateChannelUnlock", recipient)
}

// InitiateChannelUnlock is a paid mutator transaction binding the contract method 0xdf087864.
//
// Solidity: function initiateChannelUnlock(address recipient) returns()
func (_PaymentChannel *PaymentChannelSession) InitiateChannelUnlock(recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.Contract.InitiateChannelUnlock(&_PaymentChannel.TransactOpts, recipient)
}

// InitiateChannelUnlock is a paid mutator transaction binding the contract method 0xdf087864.
//
// Solidity: function initiateChannelUnlock(address recipient) returns()
func (_PaymentChannel *PaymentChannelTransactorSession) InitiateChannelUnlock(recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.Contract.InitiateChannelUnlock(&_PaymentChannel.TransactOpts, recipient)
}

// OpenChannel is a paid mutator transaction binding the contract method 0x2abb5e9d.
//
// Solidity: function openChannel(address recipient) payable returns()
func (_PaymentChannel *PaymentChannelTransactor) OpenChannel(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.contract.Transact(opts, "openChannel", recipient)
}

// OpenChannel is a paid mutator transaction binding the contract method 0x2abb5e9d.
//
// Solidity: function openChannel(address recipient) payable returns()
func (_PaymentChannel *PaymentChannelSession) OpenChannel(recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.Contract.OpenChannel(&_PaymentChannel.TransactOpts, recipient)
}

// OpenChannel is a paid mutator transaction binding the contract method 0x2abb5e9d.
//
// Solidity: function openChannel(address recipient) payable returns()
func (_PaymentChannel *PaymentChannelTransactorSession) OpenChannel(recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.Contract.OpenChannel(&_PaymentChannel.TransactOpts, recipient)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PaymentChannel *PaymentChannelTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PaymentChannel.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PaymentChannel *PaymentChannelSession) Pause() (*types.Transaction, error) {
	return _PaymentChannel.Contract.Pause(&_PaymentChannel.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PaymentChannel *PaymentChannelTransactorSession) Pause() (*types.Transaction, error) {
	return _PaymentChannel.Contract.Pause(&_PaymentChannel.TransactOpts)
}

// PayoutProvider is a paid mutator transaction binding the contract method 0x756b7b89.
//
// Solidity: function payoutProvider(address sender, address recipient, uint256 amount) returns()
func (_PaymentChannel *PaymentChannelTransactor) PayoutProvider(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PaymentChannel.contract.Transact(opts, "payoutProvider", sender, recipient, amount)
}

// PayoutProvider is a paid mutator transaction binding the contract method 0x756b7b89.
//
// Solidity: function payoutProvider(address sender, address recipient, uint256 amount) returns()
func (_PaymentChannel *PaymentChannelSession) PayoutProvider(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PaymentChannel.Contract.PayoutProvider(&_PaymentChannel.TransactOpts, sender, recipient, amount)
}

// PayoutProvider is a paid mutator transaction binding the contract method 0x756b7b89.
//
// Solidity: function payoutProvider(address sender, address recipient, uint256 amount) returns()
func (_PaymentChannel *PaymentChannelTransactorSession) PayoutProvider(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PaymentChannel.Contract.PayoutProvider(&_PaymentChannel.TransactOpts, sender, recipient, amount)
}

// SetPayoutEngine is a paid mutator transaction binding the contract method 0xa71c2df3.
//
// Solidity: function setPayoutEngine(address engine, bool authorized) returns()
func (_PaymentChannel *PaymentChannelTransactor) SetPayoutEngine(opts *bind.TransactOpts, engine common.Address, authorized bool) (*types.Transaction, error) {
	return _PaymentChannel.contract.Transact(opts, "setPayoutEngine", engine, authorized)
}

// SetPayoutEngine is a paid mutator transaction binding the contract method 0xa71c2df3.
//
// Solidity: function setPayoutEngine(address engine, bool authorized) returns()
func (_PaymentChannel *PaymentChannelSession) SetPayoutEngine(engine common.Address, authorized bool) (*types.Transaction, error) {
	return _PaymentChannel.Contract.SetPayoutEngine(&_PaymentChannel.TransactOpts, engine, authorized)
}

// SetPayoutEngine is a paid mutator transaction binding the contract method 0xa71c2df3.
//
// Solidity: function setPayoutEngine(address engine, bool authorized) returns()
func (_PaymentChannel *PaymentChannelTransactorSession) SetPayoutEngine(engine common.Address, authorized bool) (*types.Transaction, error) {
	return _PaymentChannel.Contract.SetPayoutEngine(&_PaymentChannel.TransactOpts, engine, authorized)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PaymentChannel *PaymentChannelTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PaymentChannel.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PaymentChannel *PaymentChannelSession) Unpause() (*types.Transaction, error) {
	return _PaymentChannel.Contract.Unpause(&_PaymentChannel.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PaymentChannel *PaymentChannelTransactorSession) Unpause() (*types.Transaction, error) {
	return _PaymentChannel.Contract.Unpause(&_PaymentChannel.TransactOpts)
}

// WithdrawChannel is a paid mutator transaction binding the contract method 0x28095f26.
//
// Solidity: function withdrawChannel(address recipient) returns()
func (_PaymentChannel *PaymentChannelTransactor) WithdrawChannel(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.contract.Transact(opts, "withdrawChannel", recipient)
}

// WithdrawChannel is a paid mutator transaction binding the contract method 0x28095f26.
//
// Solidity: function withdrawChannel(address recipient) returns()
func (_PaymentChannel *PaymentChannelSession) WithdrawChannel(recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.Contract.WithdrawChannel(&_PaymentChannel.TransactOpts, recipient)
}

// WithdrawChannel is a paid mutator transaction binding the contract method 0x28095f26.
//
// Solidity: function withdrawChannel(address recipient) returns()
func (_PaymentChannel *PaymentChannelTransactorSession) WithdrawChannel(recipient common.Address) (*types.Transaction, error) {
	return _PaymentChannel.Contract.WithdrawChannel(&_PaymentChannel.TransactOpts, recipient)
}

// PaymentChannelChannelDepositedIterator is returned from FilterChannelDeposited and is used to iterate over the raw logs and unpacked data for ChannelDeposited events raised by the PaymentChannel contract.
type PaymentChannelChannelDepositedIterator struct {
	Event *PaymentChannelChannelDeposited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PaymentChannelChannelDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentChannelChannelDeposited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PaymentChannelChannelDeposited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PaymentChannelChannelDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentChannelChannelDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentChannelChannelDeposited represents a ChannelDeposited event raised by the PaymentChannel contract.
type PaymentChannelChannelDeposited struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChannelDeposited is a free log retrieval operation binding the contract event 0x05c20c1748de9db22bf7db8be0f4d51c217ee76839ad9ba6aabc1d5a82e953a1.
//
// Solidity: event ChannelDeposited(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) FilterChannelDeposited(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*PaymentChannelChannelDepositedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentChannel.contract.FilterLogs(opts, "ChannelDeposited", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PaymentChannelChannelDepositedIterator{contract: _PaymentChannel.contract, event: "ChannelDeposited", logs: logs, sub: sub}, nil
}

// WatchChannelDeposited is a free log subscription operation binding the contract event 0x05c20c1748de9db22bf7db8be0f4d51c217ee76839ad9ba6aabc1d5a82e953a1.
//
// Solidity: event ChannelDeposited(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) WatchChannelDeposited(opts *bind.WatchOpts, sink chan<- *PaymentChannelChannelDeposited, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentChannel.contract.WatchLogs(opts, "ChannelDeposited", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentChannelChannelDeposited)
				if err := _PaymentChannel.contract.UnpackLog(event, "ChannelDeposited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelDeposited is a log parse operation binding the contract event 0x05c20c1748de9db22bf7db8be0f4d51c217ee76839ad9ba6aabc1d5a82e953a1.
//
// Solidity: event ChannelDeposited(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) ParseChannelDeposited(log types.Log) (*PaymentChannelChannelDeposited, error) {
	event := new(PaymentChannelChannelDeposited)
	if err := _PaymentChannel.contract.UnpackLog(event, "ChannelDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentChannelChannelOpenedIterator is returned from FilterChannelOpened and is used to iterate over the raw logs and unpacked data for ChannelOpened events raised by the PaymentChannel contract.
type PaymentChannelChannelOpenedIterator struct {
	Event *PaymentChannelChannelOpened // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PaymentChannelChannelOpenedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentChannelChannelOpened)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PaymentChannelChannelOpened)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PaymentChannelChannelOpenedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentChannelChannelOpenedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentChannelChannelOpened represents a ChannelOpened event raised by the PaymentChannel contract.
type PaymentChannelChannelOpened struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChannelOpened is a free log retrieval operation binding the contract event 0x434d899467f051cde9a358523702f186412ffe4e86423150a08680ccdb1519ee.
//
// Solidity: event ChannelOpened(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) FilterChannelOpened(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*PaymentChannelChannelOpenedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentChannel.contract.FilterLogs(opts, "ChannelOpened", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PaymentChannelChannelOpenedIterator{contract: _PaymentChannel.contract, event: "ChannelOpened", logs: logs, sub: sub}, nil
}

// WatchChannelOpened is a free log subscription operation binding the contract event 0x434d899467f051cde9a358523702f186412ffe4e86423150a08680ccdb1519ee.
//
// Solidity: event ChannelOpened(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) WatchChannelOpened(opts *bind.WatchOpts, sink chan<- *PaymentChannelChannelOpened, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentChannel.contract.WatchLogs(opts, "ChannelOpened", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentChannelChannelOpened)
				if err := _PaymentChannel.contract.UnpackLog(event, "ChannelOpened", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelOpened is a log parse operation binding the contract event 0x434d899467f051cde9a358523702f186412ffe4e86423150a08680ccdb1519ee.
//
// Solidity: event ChannelOpened(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) ParseChannelOpened(log types.Log) (*PaymentChannelChannelOpened, error) {
	event := new(PaymentChannelChannelOpened)
	if err := _PaymentChannel.contract.UnpackLog(event, "ChannelOpened", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentChannelChannelPayoutExecutedIterator is returned from FilterChannelPayoutExecuted and is used to iterate over the raw logs and unpacked data for ChannelPayoutExecuted events raised by the PaymentChannel contract.
type PaymentChannelChannelPayoutExecutedIterator struct {
	Event *PaymentChannelChannelPayoutExecuted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PaymentChannelChannelPayoutExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentChannelChannelPayoutExecuted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PaymentChannelChannelPayoutExecuted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PaymentChannelChannelPayoutExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentChannelChannelPayoutExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentChannelChannelPayoutExecuted represents a ChannelPayoutExecuted event raised by the PaymentChannel contract.
type PaymentChannelChannelPayoutExecuted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChannelPayoutExecuted is a free log retrieval operation binding the contract event 0xef4ac45c90db8ad23280036f26d8dc1584377d3a7fbcf6b3068baf7bd80b9090.
//
// Solidity: event ChannelPayoutExecuted(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) FilterChannelPayoutExecuted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*PaymentChannelChannelPayoutExecutedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentChannel.contract.FilterLogs(opts, "ChannelPayoutExecuted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PaymentChannelChannelPayoutExecutedIterator{contract: _PaymentChannel.contract, event: "ChannelPayoutExecuted", logs: logs, sub: sub}, nil
}

// WatchChannelPayoutExecuted is a free log subscription operation binding the contract event 0xef4ac45c90db8ad23280036f26d8dc1584377d3a7fbcf6b3068baf7bd80b9090.
//
// Solidity: event ChannelPayoutExecuted(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) WatchChannelPayoutExecuted(opts *bind.WatchOpts, sink chan<- *PaymentChannelChannelPayoutExecuted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentChannel.contract.WatchLogs(opts, "ChannelPayoutExecuted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentChannelChannelPayoutExecuted)
				if err := _PaymentChannel.contract.UnpackLog(event, "ChannelPayoutExecuted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelPayoutExecuted is a log parse operation binding the contract event 0xef4ac45c90db8ad23280036f26d8dc1584377d3a7fbcf6b3068baf7bd80b9090.
//
// Solidity: event ChannelPayoutExecuted(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) ParseChannelPayoutExecuted(log types.Log) (*PaymentChannelChannelPayoutExecuted, error) {
	event := new(PaymentChannelChannelPayoutExecuted)
	if err := _PaymentChannel.contract.UnpackLog(event, "ChannelPayoutExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentChannelChannelUnlockInitiatedIterator is returned from FilterChannelUnlockInitiated and is used to iterate over the raw logs and unpacked data for ChannelUnlockInitiated events raised by the PaymentChannel contract.
type PaymentChannelChannelUnlockInitiatedIterator struct {
	Event *PaymentChannelChannelUnlockInitiated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PaymentChannelChannelUnlockInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentChannelChannelUnlockInitiated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PaymentChannelChannelUnlockInitiated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PaymentChannelChannelUnlockInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentChannelChannelUnlockInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentChannelChannelUnlockInitiated represents a ChannelUnlockInitiated event raised by the PaymentChannel contract.
type PaymentChannelChannelUnlockInitiated struct {
	Sender      common.Address
	Recipient   common.Address
	UnlockBlock *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterChannelUnlockInitiated is a free log retrieval operation binding the contract event 0xcce983f8bacd60182acc7e585fbe219277117ceaf0871600204eae641c59838a.
//
// Solidity: event ChannelUnlockInitiated(address indexed sender, address indexed recipient, uint256 unlockBlock)
func (_PaymentChannel *PaymentChannelFilterer) FilterChannelUnlockInitiated(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*PaymentChannelChannelUnlockInitiatedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentChannel.contract.FilterLogs(opts, "ChannelUnlockInitiated", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PaymentChannelChannelUnlockInitiatedIterator{contract: _PaymentChannel.contract, event: "ChannelUnlockInitiated", logs: logs, sub: sub}, nil
}

// WatchChannelUnlockInitiated is a free log subscription operation binding the contract event 0xcce983f8bacd60182acc7e585fbe219277117ceaf0871600204eae641c59838a.
//
// Solidity: event ChannelUnlockInitiated(address indexed sender, address indexed recipient, uint256 unlockBlock)
func (_PaymentChannel *PaymentChannelFilterer) WatchChannelUnlockInitiated(opts *bind.WatchOpts, sink chan<- *PaymentChannelChannelUnlockInitiated, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentChannel.contract.WatchLogs(opts, "ChannelUnlockInitiated", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentChannelChannelUnlockInitiated)
				if err := _PaymentChannel.contract.UnpackLog(event, "ChannelUnlockInitiated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelUnlockInitiated is a log parse operation binding the contract event 0xcce983f8bacd60182acc7e585fbe219277117ceaf0871600204eae641c59838a.
//
// Solidity: event ChannelUnlockInitiated(address indexed sender, address indexed recipient, uint256 unlockBlock)
func (_PaymentChannel *PaymentChannelFilterer) ParseChannelUnlockInitiated(log types.Log) (*PaymentChannelChannelUnlockInitiated, error) {
	event := new(PaymentChannelChannelUnlockInitiated)
	if err := _PaymentChannel.contract.UnpackLog(event, "ChannelUnlockInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentChannelChannelWithdrawnIterator is returned from FilterChannelWithdrawn and is used to iterate over the raw logs and unpacked data for ChannelWithdrawn events raised by the PaymentChannel contract.
type PaymentChannelChannelWithdrawnIterator struct {
	Event *PaymentChannelChannelWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PaymentChannelChannelWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentChannelChannelWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PaymentChannelChannelWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PaymentChannelChannelWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentChannelChannelWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentChannelChannelWithdrawn represents a ChannelWithdrawn event raised by the PaymentChannel contract.
type PaymentChannelChannelWithdrawn struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChannelWithdrawn is a free log retrieval operation binding the contract event 0xfed549f521421ee5fb337d4e385525b833bbe3e8a060b62a0f25fb5bcb6f08c1.
//
// Solidity: event ChannelWithdrawn(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) FilterChannelWithdrawn(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*PaymentChannelChannelWithdrawnIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentChannel.contract.FilterLogs(opts, "ChannelWithdrawn", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PaymentChannelChannelWithdrawnIterator{contract: _PaymentChannel.contract, event: "ChannelWithdrawn", logs: logs, sub: sub}, nil
}

// WatchChannelWithdrawn is a free log subscription operation binding the contract event 0xfed549f521421ee5fb337d4e385525b833bbe3e8a060b62a0f25fb5bcb6f08c1.
//
// Solidity: event ChannelWithdrawn(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) WatchChannelWithdrawn(opts *bind.WatchOpts, sink chan<- *PaymentChannelChannelWithdrawn, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentChannel.contract.WatchLogs(opts, "ChannelWithdrawn", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentChannelChannelWithdrawn)
				if err := _PaymentChannel.contract.UnpackLog(event, "ChannelWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelWithdrawn is a log parse operation binding the contract event 0xfed549f521421ee5fb337d4e385525b833bbe3e8a060b62a0f25fb5bcb6f08c1.
//
// Solidity: event ChannelWithdrawn(address indexed sender, address indexed recipient, uint256 amount)
func (_PaymentChannel *PaymentChannelFilterer) ParseChannelWithdrawn(log types.Log) (*PaymentChannelChannelWithdrawn, error) {
	event := new(PaymentChannelChannelWithdrawn)
	if err := _PaymentChannel.contract.UnpackLog(event, "ChannelWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentChannelPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the PaymentChannel contract.
type PaymentChannelPausedIterator struct {
	Event *PaymentChannelPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PaymentChannelPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentChannelPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PaymentChannelPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PaymentChannelPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentChannelPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentChannelPaused represents a Paused event raised by the PaymentChannel contract.
type PaymentChannelPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_PaymentChannel *PaymentChannelFilterer) FilterPaused(opts *bind.FilterOpts) (*PaymentChannelPausedIterator, error) {

	logs, sub, err := _PaymentChannel.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &PaymentChannelPausedIterator{contract: _PaymentChannel.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_PaymentChannel *PaymentChannelFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *PaymentChannelPaused) (event.Subscription, error) {

	logs, sub, err := _PaymentChannel.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentChannelPaused)
				if err := _PaymentChannel.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_PaymentChannel *PaymentChannelFilterer) ParsePaused(log types.Log) (*PaymentChannelPaused, error) {
	event := new(PaymentChannelPaused)
	if err := _PaymentChannel.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentChannelUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the PaymentChannel contract.
type PaymentChannelUnpausedIterator struct {
	Event *PaymentChannelUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PaymentChannelUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentChannelUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PaymentChannelUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PaymentChannelUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentChannelUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentChannelUnpaused represents a Unpaused event raised by the PaymentChannel contract.
type PaymentChannelUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_PaymentChannel *PaymentChannelFilterer) FilterUnpaused(opts *bind.FilterOpts) (*PaymentChannelUnpausedIterator, error) {

	logs, sub, err := _PaymentChannel.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &PaymentChannelUnpausedIterator{contract: _PaymentChannel.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_PaymentChannel *PaymentChannelFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *PaymentChannelUnpaused) (event.Subscription, error) {

	logs, sub, err := _PaymentChannel.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentChannelUnpaused)
				if err := _PaymentChannel.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_PaymentChannel *PaymentChannelFilterer) ParseUnpaused(log types.Log) (*PaymentChannelUnpaused, error) {
	event := new(PaymentChannelUnpaused)
	if err := _PaymentChannel.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
