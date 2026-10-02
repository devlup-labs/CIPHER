// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package providerregistry

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

// ProviderRegistryMetaData contains all meta data concerning the ProviderRegistry contract.
var ProviderRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_config\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizedSlashers\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"config\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractProtocolConfig\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProviderStake\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isProviderActive\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providers\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unstakeReleaseBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pendingUnstakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"registered\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerProvider\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"requestPartialUnstake\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestUnstake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlasher\",\"inputs\":[{\"name\":\"slasher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authorized\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slashProvider\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"reasonId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"slashedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawStake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ProviderRegistered\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"totalStake\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProviderSlashed\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"reasonId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProviderStakeWithdrawn\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProviderUnstakeRequested\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"releaseBlock\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InsufficientStake\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedCaller\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnbondingNotComplete\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnstakeAlreadyRequested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnstakeNotRequested\",\"inputs\":[]}]",
	Bin: "0x60a03461012457601f610c1838819003918201601f19168301916001600160401b038311848410176101285780849260209460405283398101031261012457516001600160a01b038116908190036101245760017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580156100d157608052600280546001600160a01b03191633179055604051610adb908161013d823960805181818160f20152818161022a0152818161033a015281816104e7015281816106dc0152818161078001526109af0152f35b60405162461bcd60e51b815260206004820152602560248201527f50726f766964657252656769737472793a205a65726f20636f6e666967206164604482015264647265737360d81b6064820152608490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe6080806040526004361015610012575f80fd5b5f3560e01c9081630787bc27146108ac57508063511fcb6a1461086357806352348080146107f7578063696724cb1461070b57806379502c55146106c75780638da5cb5b1461069f578063aed2d90814610662578063bed9d8611461042e578063bfebc370146103f7578063f769c6d014610302578063f8f9ba99146101cb5763fc63958e146100a0575f80fd5b34610196575f36600319011261019657335f525f60205260405f2060ff600382015416156101bc57600181019081546101ad578054600290910181905560405163366a44bb60e21b81526020816004817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa80156101a2575f9061016a575b610134915043610a98565b80925560405191825260208201527f2721af9124037857b9b71dab6c68562c42e1c4ee5aaadef494e3ff73f820c16460403392a2005b506020813d60201161019a575b816101846020938361093d565b81010312610196576101349051610129565b5f80fd5b3d9150610177565b6040513d5f823e3d90fd5b634b8cc7d360e11b5f5260045ffd5b63aba4733960e01b5f5260045ffd5b3461019657602036600319011261019657600435335f525f60205260405f209060ff600383015416156101bc57600182019182546101ad57811580156102f8575b6102e95761021b82825461091c565b60405163650190e760e01b81527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169291602082600481875afa9182156101a2575f926102b5575b50106102a65760049183600260209301556040519283809263366a44bb60e21b82525afa80156101a2575f9061016a57610134915043610a98565b6378de4a6960e11b5f5260045ffd5b9091506020813d6020116102e1575b816102d16020938361093d565b810103126101965751908661026b565b3d91506102c4565b63162908e360e11b5f5260045ffd5b508181541061020c565b5f36600319011261019657335f525f60205260405f2080546103243482610a98565b60405163650190e760e01b8152906020826004817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9182156101a2575f926103c3575b50106102a6576103856003913490610a98565b9182815501600160ff198254161790556040519081527f90c9734131c1e4fb36cde2d71e6feb93fb258f71be8a85411c173d25e1516e8060203392a2005b9091506020813d6020116103ef575b816103df6020938361093d565b8101031261019657519084610372565b3d91506103d2565b34610196576020366003190112610196576001600160a01b03610418610906565b165f525f602052602060405f2054604051908152f35b34610196575f3660031901126101965760027f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0054146106535760027f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055335f525f60205260405f2060018101908154801561064457431061063557600281018054928315801561062b575b6102a6575f90816104cb86865461091c565b80865593555560405163650190e760e01b8152906020826004817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9182156101a2575f926105f7575b50106105e6575b505f80808084335af13d156105e1573d67ffffffffffffffff81116105cd576040519061055d601f8201601f19166020018361093d565b81525f60203d92013e5b156105be576040519081527f8ad1823705a8f19ba96a6239ccb67ea0fc9739cc0c838cd924a52ed1a61e4e1560203392a260017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055005b6312171d8360e31b5f5260045ffd5b634e487b7160e01b5f52604160045260245ffd5b610567565b600301805460ff1916905581610526565b9091506020813d602011610623575b816106136020938361093d565b810103126101965751908461051f565b3d9150610606565b50838354106104b9565b6302d21cc960e61b5f5260045ffd5b63777c843f60e11b5f5260045ffd5b633ee5aeb560e01b5f5260045ffd5b34610196576020366003190112610196576001600160a01b03610683610906565b165f526001602052602060ff60405f2054166040519015158152f35b34610196575f366003190112610196576002546040516001600160a01b039091168152602090f35b34610196575f366003190112610196576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b34610196576020366003190112610196576001600160a01b0361072c610906565b165f525f60205260405f2060ff60038201541680610766575b60209181610759575b506040519015158152f35b600191500154158261074e565b50805460405163650190e760e01b815291906020836004817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9081156101a2575f916107c2575b9192501115610745565b90506020833d6020116107ef575b816107dd6020938361093d565b810103126101965760209251906107b8565b3d91506107d0565b3461019657604036600319011261019657610810610906565b60243590811515809203610196576002546001600160a01b031633036108545760018060a01b03165f52600160205260405f209060ff801983541691161790555f80f35b635c427cd960e01b5f5260045ffd5b346101965760603660031901126101965761087c610906565b335f52600160205260ff60405f20541615610854576108a4602091604435906024359061095f565b604051908152f35b34610196576020366003190112610196576080906001600160a01b036108d0610906565b165f525f60205260405f20805490600181015460ff60036002840154930154169284526020840152604083015215156060820152f35b600435906001600160a01b038216820361019657565b9190820391821161092957565b634e487b7160e01b5f52601160045260245ffd5b90601f8019910116810190811067ffffffffffffffff8211176105cd57604052565b6001600160a01b03165f81815260208190526040902080549493928515610a8f5785811115610a8a5750845b61099681809761091c565b80835560405163650190e760e01b8152906020826004817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9182156101a2575f92610a30575b5092602092917ff9b9fcac00d461bad4d9e376e00ea3cfc45261b9f23df684f6c5d949e3537e489410610a1f575b50604051908152a3565b600301805460ff191690555f610a15565b929150926020833d602011610a82575b81610a4d6020938361093d565b8101031261019657915191929091907ff9b9fcac00d461bad4d9e376e00ea3cfc45261b9f23df684f6c5d949e3537e486109e7565b3d9150610a40565b61098b565b505f9450505050565b919082018092116109295756fea2646970667358221220808f1a9032008d91020897cb9fa0aadbd72c53ed84832201e009d0e161142d5a64736f6c63430008230033",
}

// ProviderRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ProviderRegistryMetaData.ABI instead.
var ProviderRegistryABI = ProviderRegistryMetaData.ABI

// ProviderRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ProviderRegistryMetaData.Bin instead.
var ProviderRegistryBin = ProviderRegistryMetaData.Bin

// DeployProviderRegistry deploys a new Ethereum contract, binding an instance of ProviderRegistry to it.
func DeployProviderRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _config common.Address) (common.Address, *types.Transaction, *ProviderRegistry, error) {
	parsed, err := ProviderRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ProviderRegistryBin), backend, _config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ProviderRegistry{ProviderRegistryCaller: ProviderRegistryCaller{contract: contract}, ProviderRegistryTransactor: ProviderRegistryTransactor{contract: contract}, ProviderRegistryFilterer: ProviderRegistryFilterer{contract: contract}}, nil
}

// ProviderRegistry is an auto generated Go binding around an Ethereum contract.
type ProviderRegistry struct {
	ProviderRegistryCaller     // Read-only binding to the contract
	ProviderRegistryTransactor // Write-only binding to the contract
	ProviderRegistryFilterer   // Log filterer for contract events
}

// ProviderRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProviderRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProviderRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProviderRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProviderRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProviderRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProviderRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProviderRegistrySession struct {
	Contract     *ProviderRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProviderRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProviderRegistryCallerSession struct {
	Contract *ProviderRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ProviderRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProviderRegistryTransactorSession struct {
	Contract     *ProviderRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ProviderRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProviderRegistryRaw struct {
	Contract *ProviderRegistry // Generic contract binding to access the raw methods on
}

// ProviderRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProviderRegistryCallerRaw struct {
	Contract *ProviderRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ProviderRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProviderRegistryTransactorRaw struct {
	Contract *ProviderRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProviderRegistry creates a new instance of ProviderRegistry, bound to a specific deployed contract.
func NewProviderRegistry(address common.Address, backend bind.ContractBackend) (*ProviderRegistry, error) {
	contract, err := bindProviderRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ProviderRegistry{ProviderRegistryCaller: ProviderRegistryCaller{contract: contract}, ProviderRegistryTransactor: ProviderRegistryTransactor{contract: contract}, ProviderRegistryFilterer: ProviderRegistryFilterer{contract: contract}}, nil
}

// NewProviderRegistryCaller creates a new read-only instance of ProviderRegistry, bound to a specific deployed contract.
func NewProviderRegistryCaller(address common.Address, caller bind.ContractCaller) (*ProviderRegistryCaller, error) {
	contract, err := bindProviderRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProviderRegistryCaller{contract: contract}, nil
}

// NewProviderRegistryTransactor creates a new write-only instance of ProviderRegistry, bound to a specific deployed contract.
func NewProviderRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ProviderRegistryTransactor, error) {
	contract, err := bindProviderRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProviderRegistryTransactor{contract: contract}, nil
}

// NewProviderRegistryFilterer creates a new log filterer instance of ProviderRegistry, bound to a specific deployed contract.
func NewProviderRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ProviderRegistryFilterer, error) {
	contract, err := bindProviderRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProviderRegistryFilterer{contract: contract}, nil
}

// bindProviderRegistry binds a generic wrapper to an already deployed contract.
func bindProviderRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProviderRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProviderRegistry *ProviderRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProviderRegistry.Contract.ProviderRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProviderRegistry *ProviderRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProviderRegistry.Contract.ProviderRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProviderRegistry *ProviderRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProviderRegistry.Contract.ProviderRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProviderRegistry *ProviderRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProviderRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProviderRegistry *ProviderRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProviderRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProviderRegistry *ProviderRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProviderRegistry.Contract.contract.Transact(opts, method, params...)
}

// AuthorizedSlashers is a free data retrieval call binding the contract method 0xaed2d908.
//
// Solidity: function authorizedSlashers(address ) view returns(bool)
func (_ProviderRegistry *ProviderRegistryCaller) AuthorizedSlashers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ProviderRegistry.contract.Call(opts, &out, "authorizedSlashers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizedSlashers is a free data retrieval call binding the contract method 0xaed2d908.
//
// Solidity: function authorizedSlashers(address ) view returns(bool)
func (_ProviderRegistry *ProviderRegistrySession) AuthorizedSlashers(arg0 common.Address) (bool, error) {
	return _ProviderRegistry.Contract.AuthorizedSlashers(&_ProviderRegistry.CallOpts, arg0)
}

// AuthorizedSlashers is a free data retrieval call binding the contract method 0xaed2d908.
//
// Solidity: function authorizedSlashers(address ) view returns(bool)
func (_ProviderRegistry *ProviderRegistryCallerSession) AuthorizedSlashers(arg0 common.Address) (bool, error) {
	return _ProviderRegistry.Contract.AuthorizedSlashers(&_ProviderRegistry.CallOpts, arg0)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_ProviderRegistry *ProviderRegistryCaller) Config(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ProviderRegistry.contract.Call(opts, &out, "config")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_ProviderRegistry *ProviderRegistrySession) Config() (common.Address, error) {
	return _ProviderRegistry.Contract.Config(&_ProviderRegistry.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_ProviderRegistry *ProviderRegistryCallerSession) Config() (common.Address, error) {
	return _ProviderRegistry.Contract.Config(&_ProviderRegistry.CallOpts)
}

// GetProviderStake is a free data retrieval call binding the contract method 0xbfebc370.
//
// Solidity: function getProviderStake(address provider) view returns(uint256)
func (_ProviderRegistry *ProviderRegistryCaller) GetProviderStake(opts *bind.CallOpts, provider common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ProviderRegistry.contract.Call(opts, &out, "getProviderStake", provider)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetProviderStake is a free data retrieval call binding the contract method 0xbfebc370.
//
// Solidity: function getProviderStake(address provider) view returns(uint256)
func (_ProviderRegistry *ProviderRegistrySession) GetProviderStake(provider common.Address) (*big.Int, error) {
	return _ProviderRegistry.Contract.GetProviderStake(&_ProviderRegistry.CallOpts, provider)
}

// GetProviderStake is a free data retrieval call binding the contract method 0xbfebc370.
//
// Solidity: function getProviderStake(address provider) view returns(uint256)
func (_ProviderRegistry *ProviderRegistryCallerSession) GetProviderStake(provider common.Address) (*big.Int, error) {
	return _ProviderRegistry.Contract.GetProviderStake(&_ProviderRegistry.CallOpts, provider)
}

// IsProviderActive is a free data retrieval call binding the contract method 0x696724cb.
//
// Solidity: function isProviderActive(address provider) view returns(bool)
func (_ProviderRegistry *ProviderRegistryCaller) IsProviderActive(opts *bind.CallOpts, provider common.Address) (bool, error) {
	var out []interface{}
	err := _ProviderRegistry.contract.Call(opts, &out, "isProviderActive", provider)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsProviderActive is a free data retrieval call binding the contract method 0x696724cb.
//
// Solidity: function isProviderActive(address provider) view returns(bool)
func (_ProviderRegistry *ProviderRegistrySession) IsProviderActive(provider common.Address) (bool, error) {
	return _ProviderRegistry.Contract.IsProviderActive(&_ProviderRegistry.CallOpts, provider)
}

// IsProviderActive is a free data retrieval call binding the contract method 0x696724cb.
//
// Solidity: function isProviderActive(address provider) view returns(bool)
func (_ProviderRegistry *ProviderRegistryCallerSession) IsProviderActive(provider common.Address) (bool, error) {
	return _ProviderRegistry.Contract.IsProviderActive(&_ProviderRegistry.CallOpts, provider)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProviderRegistry *ProviderRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ProviderRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProviderRegistry *ProviderRegistrySession) Owner() (common.Address, error) {
	return _ProviderRegistry.Contract.Owner(&_ProviderRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProviderRegistry *ProviderRegistryCallerSession) Owner() (common.Address, error) {
	return _ProviderRegistry.Contract.Owner(&_ProviderRegistry.CallOpts)
}

// Providers is a free data retrieval call binding the contract method 0x0787bc27.
//
// Solidity: function providers(address ) view returns(uint256 stake, uint256 unstakeReleaseBlock, uint256 pendingUnstakeAmount, bool registered)
func (_ProviderRegistry *ProviderRegistryCaller) Providers(opts *bind.CallOpts, arg0 common.Address) (struct {
	Stake                *big.Int
	UnstakeReleaseBlock  *big.Int
	PendingUnstakeAmount *big.Int
	Registered           bool
}, error) {
	var out []interface{}
	err := _ProviderRegistry.contract.Call(opts, &out, "providers", arg0)

	outstruct := new(struct {
		Stake                *big.Int
		UnstakeReleaseBlock  *big.Int
		PendingUnstakeAmount *big.Int
		Registered           bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Stake = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.UnstakeReleaseBlock = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.PendingUnstakeAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Registered = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Providers is a free data retrieval call binding the contract method 0x0787bc27.
//
// Solidity: function providers(address ) view returns(uint256 stake, uint256 unstakeReleaseBlock, uint256 pendingUnstakeAmount, bool registered)
func (_ProviderRegistry *ProviderRegistrySession) Providers(arg0 common.Address) (struct {
	Stake                *big.Int
	UnstakeReleaseBlock  *big.Int
	PendingUnstakeAmount *big.Int
	Registered           bool
}, error) {
	return _ProviderRegistry.Contract.Providers(&_ProviderRegistry.CallOpts, arg0)
}

// Providers is a free data retrieval call binding the contract method 0x0787bc27.
//
// Solidity: function providers(address ) view returns(uint256 stake, uint256 unstakeReleaseBlock, uint256 pendingUnstakeAmount, bool registered)
func (_ProviderRegistry *ProviderRegistryCallerSession) Providers(arg0 common.Address) (struct {
	Stake                *big.Int
	UnstakeReleaseBlock  *big.Int
	PendingUnstakeAmount *big.Int
	Registered           bool
}, error) {
	return _ProviderRegistry.Contract.Providers(&_ProviderRegistry.CallOpts, arg0)
}

// RegisterProvider is a paid mutator transaction binding the contract method 0xf769c6d0.
//
// Solidity: function registerProvider() payable returns()
func (_ProviderRegistry *ProviderRegistryTransactor) RegisterProvider(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProviderRegistry.contract.Transact(opts, "registerProvider")
}

// RegisterProvider is a paid mutator transaction binding the contract method 0xf769c6d0.
//
// Solidity: function registerProvider() payable returns()
func (_ProviderRegistry *ProviderRegistrySession) RegisterProvider() (*types.Transaction, error) {
	return _ProviderRegistry.Contract.RegisterProvider(&_ProviderRegistry.TransactOpts)
}

// RegisterProvider is a paid mutator transaction binding the contract method 0xf769c6d0.
//
// Solidity: function registerProvider() payable returns()
func (_ProviderRegistry *ProviderRegistryTransactorSession) RegisterProvider() (*types.Transaction, error) {
	return _ProviderRegistry.Contract.RegisterProvider(&_ProviderRegistry.TransactOpts)
}

// RequestPartialUnstake is a paid mutator transaction binding the contract method 0xf8f9ba99.
//
// Solidity: function requestPartialUnstake(uint256 amount) returns()
func (_ProviderRegistry *ProviderRegistryTransactor) RequestPartialUnstake(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ProviderRegistry.contract.Transact(opts, "requestPartialUnstake", amount)
}

// RequestPartialUnstake is a paid mutator transaction binding the contract method 0xf8f9ba99.
//
// Solidity: function requestPartialUnstake(uint256 amount) returns()
func (_ProviderRegistry *ProviderRegistrySession) RequestPartialUnstake(amount *big.Int) (*types.Transaction, error) {
	return _ProviderRegistry.Contract.RequestPartialUnstake(&_ProviderRegistry.TransactOpts, amount)
}

// RequestPartialUnstake is a paid mutator transaction binding the contract method 0xf8f9ba99.
//
// Solidity: function requestPartialUnstake(uint256 amount) returns()
func (_ProviderRegistry *ProviderRegistryTransactorSession) RequestPartialUnstake(amount *big.Int) (*types.Transaction, error) {
	return _ProviderRegistry.Contract.RequestPartialUnstake(&_ProviderRegistry.TransactOpts, amount)
}

// RequestUnstake is a paid mutator transaction binding the contract method 0xfc63958e.
//
// Solidity: function requestUnstake() returns()
func (_ProviderRegistry *ProviderRegistryTransactor) RequestUnstake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProviderRegistry.contract.Transact(opts, "requestUnstake")
}

// RequestUnstake is a paid mutator transaction binding the contract method 0xfc63958e.
//
// Solidity: function requestUnstake() returns()
func (_ProviderRegistry *ProviderRegistrySession) RequestUnstake() (*types.Transaction, error) {
	return _ProviderRegistry.Contract.RequestUnstake(&_ProviderRegistry.TransactOpts)
}

// RequestUnstake is a paid mutator transaction binding the contract method 0xfc63958e.
//
// Solidity: function requestUnstake() returns()
func (_ProviderRegistry *ProviderRegistryTransactorSession) RequestUnstake() (*types.Transaction, error) {
	return _ProviderRegistry.Contract.RequestUnstake(&_ProviderRegistry.TransactOpts)
}

// SetSlasher is a paid mutator transaction binding the contract method 0x52348080.
//
// Solidity: function setSlasher(address slasher, bool authorized) returns()
func (_ProviderRegistry *ProviderRegistryTransactor) SetSlasher(opts *bind.TransactOpts, slasher common.Address, authorized bool) (*types.Transaction, error) {
	return _ProviderRegistry.contract.Transact(opts, "setSlasher", slasher, authorized)
}

// SetSlasher is a paid mutator transaction binding the contract method 0x52348080.
//
// Solidity: function setSlasher(address slasher, bool authorized) returns()
func (_ProviderRegistry *ProviderRegistrySession) SetSlasher(slasher common.Address, authorized bool) (*types.Transaction, error) {
	return _ProviderRegistry.Contract.SetSlasher(&_ProviderRegistry.TransactOpts, slasher, authorized)
}

// SetSlasher is a paid mutator transaction binding the contract method 0x52348080.
//
// Solidity: function setSlasher(address slasher, bool authorized) returns()
func (_ProviderRegistry *ProviderRegistryTransactorSession) SetSlasher(slasher common.Address, authorized bool) (*types.Transaction, error) {
	return _ProviderRegistry.Contract.SetSlasher(&_ProviderRegistry.TransactOpts, slasher, authorized)
}

// SlashProvider is a paid mutator transaction binding the contract method 0x511fcb6a.
//
// Solidity: function slashProvider(address provider, uint256 amount, bytes32 reasonId) returns(uint256 slashedAmount)
func (_ProviderRegistry *ProviderRegistryTransactor) SlashProvider(opts *bind.TransactOpts, provider common.Address, amount *big.Int, reasonId [32]byte) (*types.Transaction, error) {
	return _ProviderRegistry.contract.Transact(opts, "slashProvider", provider, amount, reasonId)
}

// SlashProvider is a paid mutator transaction binding the contract method 0x511fcb6a.
//
// Solidity: function slashProvider(address provider, uint256 amount, bytes32 reasonId) returns(uint256 slashedAmount)
func (_ProviderRegistry *ProviderRegistrySession) SlashProvider(provider common.Address, amount *big.Int, reasonId [32]byte) (*types.Transaction, error) {
	return _ProviderRegistry.Contract.SlashProvider(&_ProviderRegistry.TransactOpts, provider, amount, reasonId)
}

// SlashProvider is a paid mutator transaction binding the contract method 0x511fcb6a.
//
// Solidity: function slashProvider(address provider, uint256 amount, bytes32 reasonId) returns(uint256 slashedAmount)
func (_ProviderRegistry *ProviderRegistryTransactorSession) SlashProvider(provider common.Address, amount *big.Int, reasonId [32]byte) (*types.Transaction, error) {
	return _ProviderRegistry.Contract.SlashProvider(&_ProviderRegistry.TransactOpts, provider, amount, reasonId)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xbed9d861.
//
// Solidity: function withdrawStake() returns()
func (_ProviderRegistry *ProviderRegistryTransactor) WithdrawStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProviderRegistry.contract.Transact(opts, "withdrawStake")
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xbed9d861.
//
// Solidity: function withdrawStake() returns()
func (_ProviderRegistry *ProviderRegistrySession) WithdrawStake() (*types.Transaction, error) {
	return _ProviderRegistry.Contract.WithdrawStake(&_ProviderRegistry.TransactOpts)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xbed9d861.
//
// Solidity: function withdrawStake() returns()
func (_ProviderRegistry *ProviderRegistryTransactorSession) WithdrawStake() (*types.Transaction, error) {
	return _ProviderRegistry.Contract.WithdrawStake(&_ProviderRegistry.TransactOpts)
}

// ProviderRegistryProviderRegisteredIterator is returned from FilterProviderRegistered and is used to iterate over the raw logs and unpacked data for ProviderRegistered events raised by the ProviderRegistry contract.
type ProviderRegistryProviderRegisteredIterator struct {
	Event *ProviderRegistryProviderRegistered // Event containing the contract specifics and raw log

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
func (it *ProviderRegistryProviderRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderRegistryProviderRegistered)
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
		it.Event = new(ProviderRegistryProviderRegistered)
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
func (it *ProviderRegistryProviderRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderRegistryProviderRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderRegistryProviderRegistered represents a ProviderRegistered event raised by the ProviderRegistry contract.
type ProviderRegistryProviderRegistered struct {
	Provider   common.Address
	TotalStake *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProviderRegistered is a free log retrieval operation binding the contract event 0x90c9734131c1e4fb36cde2d71e6feb93fb258f71be8a85411c173d25e1516e80.
//
// Solidity: event ProviderRegistered(address indexed provider, uint256 totalStake)
func (_ProviderRegistry *ProviderRegistryFilterer) FilterProviderRegistered(opts *bind.FilterOpts, provider []common.Address) (*ProviderRegistryProviderRegisteredIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _ProviderRegistry.contract.FilterLogs(opts, "ProviderRegistered", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderRegistryProviderRegisteredIterator{contract: _ProviderRegistry.contract, event: "ProviderRegistered", logs: logs, sub: sub}, nil
}

// WatchProviderRegistered is a free log subscription operation binding the contract event 0x90c9734131c1e4fb36cde2d71e6feb93fb258f71be8a85411c173d25e1516e80.
//
// Solidity: event ProviderRegistered(address indexed provider, uint256 totalStake)
func (_ProviderRegistry *ProviderRegistryFilterer) WatchProviderRegistered(opts *bind.WatchOpts, sink chan<- *ProviderRegistryProviderRegistered, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _ProviderRegistry.contract.WatchLogs(opts, "ProviderRegistered", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderRegistryProviderRegistered)
				if err := _ProviderRegistry.contract.UnpackLog(event, "ProviderRegistered", log); err != nil {
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

// ParseProviderRegistered is a log parse operation binding the contract event 0x90c9734131c1e4fb36cde2d71e6feb93fb258f71be8a85411c173d25e1516e80.
//
// Solidity: event ProviderRegistered(address indexed provider, uint256 totalStake)
func (_ProviderRegistry *ProviderRegistryFilterer) ParseProviderRegistered(log types.Log) (*ProviderRegistryProviderRegistered, error) {
	event := new(ProviderRegistryProviderRegistered)
	if err := _ProviderRegistry.contract.UnpackLog(event, "ProviderRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderRegistryProviderSlashedIterator is returned from FilterProviderSlashed and is used to iterate over the raw logs and unpacked data for ProviderSlashed events raised by the ProviderRegistry contract.
type ProviderRegistryProviderSlashedIterator struct {
	Event *ProviderRegistryProviderSlashed // Event containing the contract specifics and raw log

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
func (it *ProviderRegistryProviderSlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderRegistryProviderSlashed)
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
		it.Event = new(ProviderRegistryProviderSlashed)
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
func (it *ProviderRegistryProviderSlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderRegistryProviderSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderRegistryProviderSlashed represents a ProviderSlashed event raised by the ProviderRegistry contract.
type ProviderRegistryProviderSlashed struct {
	Provider common.Address
	Amount   *big.Int
	ReasonId [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterProviderSlashed is a free log retrieval operation binding the contract event 0xf9b9fcac00d461bad4d9e376e00ea3cfc45261b9f23df684f6c5d949e3537e48.
//
// Solidity: event ProviderSlashed(address indexed provider, uint256 amount, bytes32 indexed reasonId)
func (_ProviderRegistry *ProviderRegistryFilterer) FilterProviderSlashed(opts *bind.FilterOpts, provider []common.Address, reasonId [][32]byte) (*ProviderRegistryProviderSlashedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	var reasonIdRule []interface{}
	for _, reasonIdItem := range reasonId {
		reasonIdRule = append(reasonIdRule, reasonIdItem)
	}

	logs, sub, err := _ProviderRegistry.contract.FilterLogs(opts, "ProviderSlashed", providerRule, reasonIdRule)
	if err != nil {
		return nil, err
	}
	return &ProviderRegistryProviderSlashedIterator{contract: _ProviderRegistry.contract, event: "ProviderSlashed", logs: logs, sub: sub}, nil
}

// WatchProviderSlashed is a free log subscription operation binding the contract event 0xf9b9fcac00d461bad4d9e376e00ea3cfc45261b9f23df684f6c5d949e3537e48.
//
// Solidity: event ProviderSlashed(address indexed provider, uint256 amount, bytes32 indexed reasonId)
func (_ProviderRegistry *ProviderRegistryFilterer) WatchProviderSlashed(opts *bind.WatchOpts, sink chan<- *ProviderRegistryProviderSlashed, provider []common.Address, reasonId [][32]byte) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	var reasonIdRule []interface{}
	for _, reasonIdItem := range reasonId {
		reasonIdRule = append(reasonIdRule, reasonIdItem)
	}

	logs, sub, err := _ProviderRegistry.contract.WatchLogs(opts, "ProviderSlashed", providerRule, reasonIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderRegistryProviderSlashed)
				if err := _ProviderRegistry.contract.UnpackLog(event, "ProviderSlashed", log); err != nil {
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

// ParseProviderSlashed is a log parse operation binding the contract event 0xf9b9fcac00d461bad4d9e376e00ea3cfc45261b9f23df684f6c5d949e3537e48.
//
// Solidity: event ProviderSlashed(address indexed provider, uint256 amount, bytes32 indexed reasonId)
func (_ProviderRegistry *ProviderRegistryFilterer) ParseProviderSlashed(log types.Log) (*ProviderRegistryProviderSlashed, error) {
	event := new(ProviderRegistryProviderSlashed)
	if err := _ProviderRegistry.contract.UnpackLog(event, "ProviderSlashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderRegistryProviderStakeWithdrawnIterator is returned from FilterProviderStakeWithdrawn and is used to iterate over the raw logs and unpacked data for ProviderStakeWithdrawn events raised by the ProviderRegistry contract.
type ProviderRegistryProviderStakeWithdrawnIterator struct {
	Event *ProviderRegistryProviderStakeWithdrawn // Event containing the contract specifics and raw log

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
func (it *ProviderRegistryProviderStakeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderRegistryProviderStakeWithdrawn)
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
		it.Event = new(ProviderRegistryProviderStakeWithdrawn)
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
func (it *ProviderRegistryProviderStakeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderRegistryProviderStakeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderRegistryProviderStakeWithdrawn represents a ProviderStakeWithdrawn event raised by the ProviderRegistry contract.
type ProviderRegistryProviderStakeWithdrawn struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterProviderStakeWithdrawn is a free log retrieval operation binding the contract event 0x8ad1823705a8f19ba96a6239ccb67ea0fc9739cc0c838cd924a52ed1a61e4e15.
//
// Solidity: event ProviderStakeWithdrawn(address indexed provider, uint256 amount)
func (_ProviderRegistry *ProviderRegistryFilterer) FilterProviderStakeWithdrawn(opts *bind.FilterOpts, provider []common.Address) (*ProviderRegistryProviderStakeWithdrawnIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _ProviderRegistry.contract.FilterLogs(opts, "ProviderStakeWithdrawn", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderRegistryProviderStakeWithdrawnIterator{contract: _ProviderRegistry.contract, event: "ProviderStakeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchProviderStakeWithdrawn is a free log subscription operation binding the contract event 0x8ad1823705a8f19ba96a6239ccb67ea0fc9739cc0c838cd924a52ed1a61e4e15.
//
// Solidity: event ProviderStakeWithdrawn(address indexed provider, uint256 amount)
func (_ProviderRegistry *ProviderRegistryFilterer) WatchProviderStakeWithdrawn(opts *bind.WatchOpts, sink chan<- *ProviderRegistryProviderStakeWithdrawn, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _ProviderRegistry.contract.WatchLogs(opts, "ProviderStakeWithdrawn", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderRegistryProviderStakeWithdrawn)
				if err := _ProviderRegistry.contract.UnpackLog(event, "ProviderStakeWithdrawn", log); err != nil {
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

// ParseProviderStakeWithdrawn is a log parse operation binding the contract event 0x8ad1823705a8f19ba96a6239ccb67ea0fc9739cc0c838cd924a52ed1a61e4e15.
//
// Solidity: event ProviderStakeWithdrawn(address indexed provider, uint256 amount)
func (_ProviderRegistry *ProviderRegistryFilterer) ParseProviderStakeWithdrawn(log types.Log) (*ProviderRegistryProviderStakeWithdrawn, error) {
	event := new(ProviderRegistryProviderStakeWithdrawn)
	if err := _ProviderRegistry.contract.UnpackLog(event, "ProviderStakeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderRegistryProviderUnstakeRequestedIterator is returned from FilterProviderUnstakeRequested and is used to iterate over the raw logs and unpacked data for ProviderUnstakeRequested events raised by the ProviderRegistry contract.
type ProviderRegistryProviderUnstakeRequestedIterator struct {
	Event *ProviderRegistryProviderUnstakeRequested // Event containing the contract specifics and raw log

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
func (it *ProviderRegistryProviderUnstakeRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderRegistryProviderUnstakeRequested)
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
		it.Event = new(ProviderRegistryProviderUnstakeRequested)
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
func (it *ProviderRegistryProviderUnstakeRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderRegistryProviderUnstakeRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderRegistryProviderUnstakeRequested represents a ProviderUnstakeRequested event raised by the ProviderRegistry contract.
type ProviderRegistryProviderUnstakeRequested struct {
	Provider     common.Address
	ReleaseBlock *big.Int
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterProviderUnstakeRequested is a free log retrieval operation binding the contract event 0x2721af9124037857b9b71dab6c68562c42e1c4ee5aaadef494e3ff73f820c164.
//
// Solidity: event ProviderUnstakeRequested(address indexed provider, uint256 releaseBlock, uint256 amount)
func (_ProviderRegistry *ProviderRegistryFilterer) FilterProviderUnstakeRequested(opts *bind.FilterOpts, provider []common.Address) (*ProviderRegistryProviderUnstakeRequestedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _ProviderRegistry.contract.FilterLogs(opts, "ProviderUnstakeRequested", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderRegistryProviderUnstakeRequestedIterator{contract: _ProviderRegistry.contract, event: "ProviderUnstakeRequested", logs: logs, sub: sub}, nil
}

// WatchProviderUnstakeRequested is a free log subscription operation binding the contract event 0x2721af9124037857b9b71dab6c68562c42e1c4ee5aaadef494e3ff73f820c164.
//
// Solidity: event ProviderUnstakeRequested(address indexed provider, uint256 releaseBlock, uint256 amount)
func (_ProviderRegistry *ProviderRegistryFilterer) WatchProviderUnstakeRequested(opts *bind.WatchOpts, sink chan<- *ProviderRegistryProviderUnstakeRequested, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _ProviderRegistry.contract.WatchLogs(opts, "ProviderUnstakeRequested", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderRegistryProviderUnstakeRequested)
				if err := _ProviderRegistry.contract.UnpackLog(event, "ProviderUnstakeRequested", log); err != nil {
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

// ParseProviderUnstakeRequested is a log parse operation binding the contract event 0x2721af9124037857b9b71dab6c68562c42e1c4ee5aaadef494e3ff73f820c164.
//
// Solidity: event ProviderUnstakeRequested(address indexed provider, uint256 releaseBlock, uint256 amount)
func (_ProviderRegistry *ProviderRegistryFilterer) ParseProviderUnstakeRequested(log types.Log) (*ProviderRegistryProviderUnstakeRequested, error) {
	event := new(ProviderRegistryProviderUnstakeRequested)
	if err := _ProviderRegistry.contract.UnpackLog(event, "ProviderUnstakeRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
