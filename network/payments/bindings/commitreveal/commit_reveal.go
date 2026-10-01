// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package commitreveal

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

// IEntropySourceRoundTicket is an auto generated low-level Go binding around an user-defined struct.
type IEntropySourceRoundTicket struct {
	Sender      common.Address
	Recipient   common.Address
	RoundId     *big.Int
	LocalIndex  *big.Int
	FaceValue   *big.Int
	WinProb     *big.Int
	SenderNonce *big.Int
}

// CommitRevealEntropyMetaData contains all meta data concerning the CommitRevealEntropy contract.
var CommitRevealEntropyMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_config\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ROUND_TICKET_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"computeWinnerIndex\",\"inputs\":[{\"name\":\"commitBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tau\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"winnerIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"targetBlockhash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"config\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractProtocolConfig\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoverSigner\",\"inputs\":[{\"name\":\"ticketHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"roundTicketHash\",\"inputs\":[{\"name\":\"ticket\",\"type\":\"tuple\",\"internalType\":\"structIEntropySource.RoundTicket\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"faceValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"winProb\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"senderNonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifySecretCommitment\",\"inputs\":[{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"committedHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"error\",\"name\":\"BlockhashUnavailable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidSecret\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidShortString\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundDrawExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundNotYetDrawable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[{\"name\":\"str\",\"type\":\"string\",\"internalType\":\"string\"}]}]",
	Bin: "0x61018080604052346101e457602081610b45803803809161002082856101e8565b8339810103126101e457516001600160a01b038116908190036101e45760405161004b6040826101e8565b6017815260208101907f434950484552205061796d656e742050726f746f636f6c0000000000000000008252604051916100866040846101e8565b600583526020830191640312e302e360dc1b83526100a38161021f565b610120526100b08461021f565b61014052519020918260e05251902080610100524660a0526040519060208201927f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f8452604083015260608201524660808201523060a082015260a0815261011960c0826101e8565b5190206080523060c052801561018e57610160526040516108ba908161028b82396080518161071f015260a051816107dc015260c051816106e9015260e0518161076e015261010051816107940152610120518161045e0152610140518161048701526101605181818160d001526105500152f35b60405162461bcd60e51b815260206004820152602860248201527f436f6d6d697452657665616c456e74726f70793a205a65726f20636f6e666967604482015267206164647265737360c01b6064820152608490fd5b5f80fd5b601f909101601f19168101906001600160401b0382119082101761020b57604052565b634e487b7160e01b5f52604160045260245ffd5b601f81511161024a57602081519101516020821061023b571790565b5f198260200360031b1b161790565b604460209160405192839163305a27a960e01b83528160048401528051918291826024860152018484015e5f828201840152601f01601f19168101030190fdfe6080806040526004361015610012575f80fd5b5f3560e01c90816379502c551461053e5750806384b0196e146104465780638dec32da1461040a57806397aba7f914610365578063bd5b99b014610288578063cbdd2539146100a75763f513a56514610069575f80fd5b346100a3575f3660031901126100a35760206040517fdb748eecae65a833ffc39ad4f2b845692691c9e96c79b52d9249207f04fc238a8152f35b5f80fd5b346100a35760803660031901126100a3576040516318bc00fb60e31b81526064359060048035917f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031691602090829081855afa801561023a575f90610254575b61011a9150836105d9565b9182431061024557602060049260405193848092636093c0e960e01b82525afa90811561023a575f91610204575b61015292506105d9565b43116101f557409081156101e65780156101a25760409182516020810190828252602435858201526044356060820152606081526101916080826105a3565b519020918351920682526020820152f35b606460405162461bcd60e51b815260206004820152602060248201527f526166666c654d6174683a20546175206d75737420626520706f7369746976656044820152fd5b63492aae0760e11b5f5260045ffd5b63107c3f0d60e21b5f5260045ffd5b90506020823d602011610232575b8161021f602093836105a3565b810103126100a357610152915190610148565b3d9150610212565b6040513d5f823e3d90fd5b6359ce759960e11b5f5260045ffd5b506020813d602011610280575b8161026e602093836105a3565b810103126100a35761011a905161010f565b3d9150610261565b346100a35760e03660031901126100a3576004356001600160a01b038116908181036100a357506024356001600160a01b038116918282036100a3576020926042925060405190848201927fdb748eecae65a833ffc39ad4f2b845692691c9e96c79b52d9249207f04fc238a845260408301526060820152604435608082015260643560a082015260843560c082015260a43560e082015260c435610100820152610100815261033a610120826105a3565b5190206103456106e6565b906040519161190160f01b83526002830152602282015220604051908152f35b346100a35760403660031901126100a35760243567ffffffffffffffff81116100a357366023820112156100a35780600401359067ffffffffffffffff82116100a35736602483830101116100a3576103ef6103f8915f6020809580602483601f19601f84011601956103db60405197886105a3565b828752018386013783010152600435610638565b90929192610672565b6040516001600160a01b039091168152f35b346100a35760403660031901126100a357602060405181810160043581528282526104366040836105a3565b6024359151902014604051908152f35b346100a3575f3660031901126100a3576104e26104827f00000000000000000000000000000000000000000000000000000000000000006105fa565b6104ab7f00000000000000000000000000000000000000000000000000000000000000006105fa565b60206104f0604051926104be83856105a3565b5f84525f368137604051958695600f60f81b875260e08588015260e087019061057f565b90858203604087015261057f565b4660608501523060808501525f60a085015283810360c08501528180845192838152019301915f5b82811061052757505050500390f35b835185528695509381019392810192600101610518565b346100a3575f3660031901126100a3577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b90601f8019910116810190811067ffffffffffffffff8211176105c557604052565b634e487b7160e01b5f52604160045260245ffd5b919082018092116105e657565b634e487b7160e01b5f52601160045260245ffd5b60ff811690601f821161062957604051916106166040846105a3565b6020808452838101919036833783525290565b632cd44ac360e21b5f5260045ffd5b8151919060418303610668576106619250602082015190606060408401519301515f1a90610802565b9192909190565b50505f9160029190565b60048110156106d25780610684575050565b6001810361069b5763f645eedf60e01b5f5260045ffd5b600281036106b6575063fce698f760e01b5f5260045260245ffd5b6003146106c05750565b6335e2f38360e21b5f5260045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614806107d9575b15610741577f000000000000000000000000000000000000000000000000000000000000000090565b60405160208101907f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f82527f000000000000000000000000000000000000000000000000000000000000000060408201527f000000000000000000000000000000000000000000000000000000000000000060608201524660808201523060a082015260a081526107d360c0826105a3565b51902090565b507f00000000000000000000000000000000000000000000000000000000000000004614610718565b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411610879579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa1561023a575f516001600160a01b0381161561086f57905f905f90565b505f906001905f90565b5050505f916003919056fea2646970667358221220306fa6aac57e9fd1c83340b8a7f03bb2b341fdf876ce62b398406aafb4c4d21e64736f6c63430008230033",
}

// CommitRevealEntropyABI is the input ABI used to generate the binding from.
// Deprecated: Use CommitRevealEntropyMetaData.ABI instead.
var CommitRevealEntropyABI = CommitRevealEntropyMetaData.ABI

// CommitRevealEntropyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CommitRevealEntropyMetaData.Bin instead.
var CommitRevealEntropyBin = CommitRevealEntropyMetaData.Bin

// DeployCommitRevealEntropy deploys a new Ethereum contract, binding an instance of CommitRevealEntropy to it.
func DeployCommitRevealEntropy(auth *bind.TransactOpts, backend bind.ContractBackend, _config common.Address) (common.Address, *types.Transaction, *CommitRevealEntropy, error) {
	parsed, err := CommitRevealEntropyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CommitRevealEntropyBin), backend, _config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CommitRevealEntropy{CommitRevealEntropyCaller: CommitRevealEntropyCaller{contract: contract}, CommitRevealEntropyTransactor: CommitRevealEntropyTransactor{contract: contract}, CommitRevealEntropyFilterer: CommitRevealEntropyFilterer{contract: contract}}, nil
}

// CommitRevealEntropy is an auto generated Go binding around an Ethereum contract.
type CommitRevealEntropy struct {
	CommitRevealEntropyCaller     // Read-only binding to the contract
	CommitRevealEntropyTransactor // Write-only binding to the contract
	CommitRevealEntropyFilterer   // Log filterer for contract events
}

// CommitRevealEntropyCaller is an auto generated read-only Go binding around an Ethereum contract.
type CommitRevealEntropyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommitRevealEntropyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CommitRevealEntropyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommitRevealEntropyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CommitRevealEntropyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommitRevealEntropySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CommitRevealEntropySession struct {
	Contract     *CommitRevealEntropy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// CommitRevealEntropyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CommitRevealEntropyCallerSession struct {
	Contract *CommitRevealEntropyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// CommitRevealEntropyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CommitRevealEntropyTransactorSession struct {
	Contract     *CommitRevealEntropyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// CommitRevealEntropyRaw is an auto generated low-level Go binding around an Ethereum contract.
type CommitRevealEntropyRaw struct {
	Contract *CommitRevealEntropy // Generic contract binding to access the raw methods on
}

// CommitRevealEntropyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CommitRevealEntropyCallerRaw struct {
	Contract *CommitRevealEntropyCaller // Generic read-only contract binding to access the raw methods on
}

// CommitRevealEntropyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CommitRevealEntropyTransactorRaw struct {
	Contract *CommitRevealEntropyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCommitRevealEntropy creates a new instance of CommitRevealEntropy, bound to a specific deployed contract.
func NewCommitRevealEntropy(address common.Address, backend bind.ContractBackend) (*CommitRevealEntropy, error) {
	contract, err := bindCommitRevealEntropy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CommitRevealEntropy{CommitRevealEntropyCaller: CommitRevealEntropyCaller{contract: contract}, CommitRevealEntropyTransactor: CommitRevealEntropyTransactor{contract: contract}, CommitRevealEntropyFilterer: CommitRevealEntropyFilterer{contract: contract}}, nil
}

// NewCommitRevealEntropyCaller creates a new read-only instance of CommitRevealEntropy, bound to a specific deployed contract.
func NewCommitRevealEntropyCaller(address common.Address, caller bind.ContractCaller) (*CommitRevealEntropyCaller, error) {
	contract, err := bindCommitRevealEntropy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CommitRevealEntropyCaller{contract: contract}, nil
}

// NewCommitRevealEntropyTransactor creates a new write-only instance of CommitRevealEntropy, bound to a specific deployed contract.
func NewCommitRevealEntropyTransactor(address common.Address, transactor bind.ContractTransactor) (*CommitRevealEntropyTransactor, error) {
	contract, err := bindCommitRevealEntropy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CommitRevealEntropyTransactor{contract: contract}, nil
}

// NewCommitRevealEntropyFilterer creates a new log filterer instance of CommitRevealEntropy, bound to a specific deployed contract.
func NewCommitRevealEntropyFilterer(address common.Address, filterer bind.ContractFilterer) (*CommitRevealEntropyFilterer, error) {
	contract, err := bindCommitRevealEntropy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CommitRevealEntropyFilterer{contract: contract}, nil
}

// bindCommitRevealEntropy binds a generic wrapper to an already deployed contract.
func bindCommitRevealEntropy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CommitRevealEntropyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CommitRevealEntropy *CommitRevealEntropyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CommitRevealEntropy.Contract.CommitRevealEntropyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CommitRevealEntropy *CommitRevealEntropyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitRevealEntropy.Contract.CommitRevealEntropyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CommitRevealEntropy *CommitRevealEntropyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CommitRevealEntropy.Contract.CommitRevealEntropyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CommitRevealEntropy *CommitRevealEntropyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CommitRevealEntropy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CommitRevealEntropy *CommitRevealEntropyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CommitRevealEntropy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CommitRevealEntropy *CommitRevealEntropyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CommitRevealEntropy.Contract.contract.Transact(opts, method, params...)
}

// ROUNDTICKETTYPEHASH is a free data retrieval call binding the contract method 0xf513a565.
//
// Solidity: function ROUND_TICKET_TYPEHASH() view returns(bytes32)
func (_CommitRevealEntropy *CommitRevealEntropyCaller) ROUNDTICKETTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _CommitRevealEntropy.contract.Call(opts, &out, "ROUND_TICKET_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ROUNDTICKETTYPEHASH is a free data retrieval call binding the contract method 0xf513a565.
//
// Solidity: function ROUND_TICKET_TYPEHASH() view returns(bytes32)
func (_CommitRevealEntropy *CommitRevealEntropySession) ROUNDTICKETTYPEHASH() ([32]byte, error) {
	return _CommitRevealEntropy.Contract.ROUNDTICKETTYPEHASH(&_CommitRevealEntropy.CallOpts)
}

// ROUNDTICKETTYPEHASH is a free data retrieval call binding the contract method 0xf513a565.
//
// Solidity: function ROUND_TICKET_TYPEHASH() view returns(bytes32)
func (_CommitRevealEntropy *CommitRevealEntropyCallerSession) ROUNDTICKETTYPEHASH() ([32]byte, error) {
	return _CommitRevealEntropy.Contract.ROUNDTICKETTYPEHASH(&_CommitRevealEntropy.CallOpts)
}

// ComputeWinnerIndex is a free data retrieval call binding the contract method 0xcbdd2539.
//
// Solidity: function computeWinnerIndex(uint256 commitBlock, bytes32 secret, uint256 roundId, uint256 tau) view returns(uint256 winnerIndex, bytes32 targetBlockhash)
func (_CommitRevealEntropy *CommitRevealEntropyCaller) ComputeWinnerIndex(opts *bind.CallOpts, commitBlock *big.Int, secret [32]byte, roundId *big.Int, tau *big.Int) (struct {
	WinnerIndex     *big.Int
	TargetBlockhash [32]byte
}, error) {
	var out []interface{}
	err := _CommitRevealEntropy.contract.Call(opts, &out, "computeWinnerIndex", commitBlock, secret, roundId, tau)

	outstruct := new(struct {
		WinnerIndex     *big.Int
		TargetBlockhash [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.WinnerIndex = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TargetBlockhash = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// ComputeWinnerIndex is a free data retrieval call binding the contract method 0xcbdd2539.
//
// Solidity: function computeWinnerIndex(uint256 commitBlock, bytes32 secret, uint256 roundId, uint256 tau) view returns(uint256 winnerIndex, bytes32 targetBlockhash)
func (_CommitRevealEntropy *CommitRevealEntropySession) ComputeWinnerIndex(commitBlock *big.Int, secret [32]byte, roundId *big.Int, tau *big.Int) (struct {
	WinnerIndex     *big.Int
	TargetBlockhash [32]byte
}, error) {
	return _CommitRevealEntropy.Contract.ComputeWinnerIndex(&_CommitRevealEntropy.CallOpts, commitBlock, secret, roundId, tau)
}

// ComputeWinnerIndex is a free data retrieval call binding the contract method 0xcbdd2539.
//
// Solidity: function computeWinnerIndex(uint256 commitBlock, bytes32 secret, uint256 roundId, uint256 tau) view returns(uint256 winnerIndex, bytes32 targetBlockhash)
func (_CommitRevealEntropy *CommitRevealEntropyCallerSession) ComputeWinnerIndex(commitBlock *big.Int, secret [32]byte, roundId *big.Int, tau *big.Int) (struct {
	WinnerIndex     *big.Int
	TargetBlockhash [32]byte
}, error) {
	return _CommitRevealEntropy.Contract.ComputeWinnerIndex(&_CommitRevealEntropy.CallOpts, commitBlock, secret, roundId, tau)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_CommitRevealEntropy *CommitRevealEntropyCaller) Config(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CommitRevealEntropy.contract.Call(opts, &out, "config")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_CommitRevealEntropy *CommitRevealEntropySession) Config() (common.Address, error) {
	return _CommitRevealEntropy.Contract.Config(&_CommitRevealEntropy.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_CommitRevealEntropy *CommitRevealEntropyCallerSession) Config() (common.Address, error) {
	return _CommitRevealEntropy.Contract.Config(&_CommitRevealEntropy.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_CommitRevealEntropy *CommitRevealEntropyCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _CommitRevealEntropy.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_CommitRevealEntropy *CommitRevealEntropySession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _CommitRevealEntropy.Contract.Eip712Domain(&_CommitRevealEntropy.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_CommitRevealEntropy *CommitRevealEntropyCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _CommitRevealEntropy.Contract.Eip712Domain(&_CommitRevealEntropy.CallOpts)
}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 ticketHash, bytes sig) pure returns(address)
func (_CommitRevealEntropy *CommitRevealEntropyCaller) RecoverSigner(opts *bind.CallOpts, ticketHash [32]byte, sig []byte) (common.Address, error) {
	var out []interface{}
	err := _CommitRevealEntropy.contract.Call(opts, &out, "recoverSigner", ticketHash, sig)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 ticketHash, bytes sig) pure returns(address)
func (_CommitRevealEntropy *CommitRevealEntropySession) RecoverSigner(ticketHash [32]byte, sig []byte) (common.Address, error) {
	return _CommitRevealEntropy.Contract.RecoverSigner(&_CommitRevealEntropy.CallOpts, ticketHash, sig)
}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 ticketHash, bytes sig) pure returns(address)
func (_CommitRevealEntropy *CommitRevealEntropyCallerSession) RecoverSigner(ticketHash [32]byte, sig []byte) (common.Address, error) {
	return _CommitRevealEntropy.Contract.RecoverSigner(&_CommitRevealEntropy.CallOpts, ticketHash, sig)
}

// RoundTicketHash is a free data retrieval call binding the contract method 0xbd5b99b0.
//
// Solidity: function roundTicketHash((address,address,uint256,uint256,uint256,uint256,uint256) ticket) view returns(bytes32)
func (_CommitRevealEntropy *CommitRevealEntropyCaller) RoundTicketHash(opts *bind.CallOpts, ticket IEntropySourceRoundTicket) ([32]byte, error) {
	var out []interface{}
	err := _CommitRevealEntropy.contract.Call(opts, &out, "roundTicketHash", ticket)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RoundTicketHash is a free data retrieval call binding the contract method 0xbd5b99b0.
//
// Solidity: function roundTicketHash((address,address,uint256,uint256,uint256,uint256,uint256) ticket) view returns(bytes32)
func (_CommitRevealEntropy *CommitRevealEntropySession) RoundTicketHash(ticket IEntropySourceRoundTicket) ([32]byte, error) {
	return _CommitRevealEntropy.Contract.RoundTicketHash(&_CommitRevealEntropy.CallOpts, ticket)
}

// RoundTicketHash is a free data retrieval call binding the contract method 0xbd5b99b0.
//
// Solidity: function roundTicketHash((address,address,uint256,uint256,uint256,uint256,uint256) ticket) view returns(bytes32)
func (_CommitRevealEntropy *CommitRevealEntropyCallerSession) RoundTicketHash(ticket IEntropySourceRoundTicket) ([32]byte, error) {
	return _CommitRevealEntropy.Contract.RoundTicketHash(&_CommitRevealEntropy.CallOpts, ticket)
}

// VerifySecretCommitment is a free data retrieval call binding the contract method 0x8dec32da.
//
// Solidity: function verifySecretCommitment(bytes32 secret, bytes32 committedHash) pure returns(bool)
func (_CommitRevealEntropy *CommitRevealEntropyCaller) VerifySecretCommitment(opts *bind.CallOpts, secret [32]byte, committedHash [32]byte) (bool, error) {
	var out []interface{}
	err := _CommitRevealEntropy.contract.Call(opts, &out, "verifySecretCommitment", secret, committedHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifySecretCommitment is a free data retrieval call binding the contract method 0x8dec32da.
//
// Solidity: function verifySecretCommitment(bytes32 secret, bytes32 committedHash) pure returns(bool)
func (_CommitRevealEntropy *CommitRevealEntropySession) VerifySecretCommitment(secret [32]byte, committedHash [32]byte) (bool, error) {
	return _CommitRevealEntropy.Contract.VerifySecretCommitment(&_CommitRevealEntropy.CallOpts, secret, committedHash)
}

// VerifySecretCommitment is a free data retrieval call binding the contract method 0x8dec32da.
//
// Solidity: function verifySecretCommitment(bytes32 secret, bytes32 committedHash) pure returns(bool)
func (_CommitRevealEntropy *CommitRevealEntropyCallerSession) VerifySecretCommitment(secret [32]byte, committedHash [32]byte) (bool, error) {
	return _CommitRevealEntropy.Contract.VerifySecretCommitment(&_CommitRevealEntropy.CallOpts, secret, committedHash)
}

// CommitRevealEntropyEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the CommitRevealEntropy contract.
type CommitRevealEntropyEIP712DomainChangedIterator struct {
	Event *CommitRevealEntropyEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *CommitRevealEntropyEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommitRevealEntropyEIP712DomainChanged)
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
		it.Event = new(CommitRevealEntropyEIP712DomainChanged)
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
func (it *CommitRevealEntropyEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommitRevealEntropyEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommitRevealEntropyEIP712DomainChanged represents a EIP712DomainChanged event raised by the CommitRevealEntropy contract.
type CommitRevealEntropyEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_CommitRevealEntropy *CommitRevealEntropyFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*CommitRevealEntropyEIP712DomainChangedIterator, error) {

	logs, sub, err := _CommitRevealEntropy.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &CommitRevealEntropyEIP712DomainChangedIterator{contract: _CommitRevealEntropy.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_CommitRevealEntropy *CommitRevealEntropyFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *CommitRevealEntropyEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _CommitRevealEntropy.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommitRevealEntropyEIP712DomainChanged)
				if err := _CommitRevealEntropy.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_CommitRevealEntropy *CommitRevealEntropyFilterer) ParseEIP712DomainChanged(log types.Log) (*CommitRevealEntropyEIP712DomainChanged, error) {
	event := new(CommitRevealEntropyEIP712DomainChanged)
	if err := _CommitRevealEntropy.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
