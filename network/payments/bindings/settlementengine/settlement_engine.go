// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package settlementengine

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

// SettlementEngineMetaData contains all meta data concerning the SettlementEngine contract.
var SettlementEngineMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_config\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_entropySource\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_paymentChannel\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_providerRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimFallbackTicket\",\"inputs\":[{\"name\":\"ticket\",\"type\":\"tuple\",\"internalType\":\"structIEntropySource.RoundTicket\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"faceValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"winProb\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"senderNonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"senderSig\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitRoundSecret\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tau\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"faceValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipientRandHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"config\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractProtocolConfig\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"disputeResolver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDisputeResolver\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"entropySource\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEntropySource\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoundStatus\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumISettlementEngine.RoundStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"markRoundPartial\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paymentChannel\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPaymentChannel\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewWinnerIndex\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIProviderRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rounds\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"tau\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"faceValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipientRandHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumISettlementEngine.RoundStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDisputeResolver\",\"inputs\":[{\"name\":\"_disputeResolver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settleRound\",\"inputs\":[{\"name\":\"ticket\",\"type\":\"tuple\",\"internalType\":\"structIEntropySource.RoundTicket\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"faceValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"winProb\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"senderNonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"senderSig\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settleRoundBatch\",\"inputs\":[{\"name\":\"tickets\",\"type\":\"tuple[]\",\"internalType\":\"structIEntropySource.RoundTicket[]\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"localIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"faceValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"winProb\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"senderNonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"senderSigs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"secrets\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"usedFallbackTickets\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"usedTickets\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"voidRound\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"FallbackTicketClaimed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"localIndex\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"faceValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoundCommitted\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"tau\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"faceValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"commitBlock\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoundMarkedPartial\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoundSettled\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"winningLocalIndex\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"faceValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoundVoided\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BlockhashUnavailable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRoundParams\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSecret\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotWinningIndex\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundAlreadyCommitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundDisputePending\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundIsVoided\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundNotCommitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundNotOpen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundNotPartial\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RoundNotStaleEnough\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TicketAlreadyUsed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedCaller\",\"inputs\":[]}]",
	Bin: "0x6101003461029357601f611e0e38819003918201601f19168301916001600160401b038311848410176102975780849260809460405283398101031261029357610048816102ab565b610054602083016102ab565b9161006d6060610066604084016102ab565b92016102ab565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00556001600160a01b0390921692831561024e576001600160a01b0316908115610209576001600160a01b03169182156101c4576001600160a01b031692831561017f5760805260a05260c05260e052600480546001600160a01b03191633179055604051611b4e90816102c0823960805181818161035c0152818161093301528181610ab10152611036015260a0518181816102520152818161066801528181610d5d015261174b015260c0518181816101890152818161075801528181610977015281816109f001528181610cd501528181610efb015281816113020152611670015260e05181610c5b0152f35b60405162461bcd60e51b815260206004820152601f60248201527f536574746c656d656e74456e67696e653a205a65726f207265676973747279006044820152606490fd5b60405162461bcd60e51b815260206004820152601e60248201527f536574746c656d656e74456e67696e653a205a65726f206368616e6e656c00006044820152606490fd5b60405162461bcd60e51b815260206004820152601e60248201527f536574746c656d656e74456e67696e653a205a65726f20656e74726f707900006044820152606490fd5b60405162461bcd60e51b815260206004820152601d60248201527f536574746c656d656e74456e67696e653a205a65726f20636f6e6669670000006044820152606490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b03821682036102935756fe6080806040526004361015610012575f80fd5b5f905f3560e01c908163232ffabb146112a9575080633757a3da146111535780633e7678b4146111165780633f4ba83a146110a857806342bf1beb14610e8e57806351110b6514610c8a578063545921d914610c4657806359a515ba14610c175780635c975abb14610bf65780635da3d4f8146109a65780636df24cd91461096257806379502c551461091e5780638456cb59146108bd5780638da5cb5b14610895578063924e63f6146108445780639c390612146106f8578063a23ad88d14610697578063ae174b9f14610653578063bb71decb14610161578063d0a3155e146101325763f5a3f4af14610105575f80fd5b3461012f578060031936011261012f575460405160089190911c6001600160a01b03168152602090f35b80fd5b503461012f57602036600319011261012f5760ff60406020926004358152600384522054166040519015158152f35b50346104e5576101fb61017336611403565b9061017f939293611615565b61018761164d565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316916101bb84611527565b946020808601966101cb88611527565b6040516311b663f160e31b81526001600160a01b0392831660048201529116602482015297889081906044820190565b0381875afa9687156104da575f9761061f575b5061021e604086013580986115ef565b5f52600160205260405f209060ff600483015416600581101561060b576003036105fc57604051630bd5b99b60e41b8152937f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031691906020868061028d8b6004830161153b565b0381865afa9586156104da575f966105c8575b50855f52600360205260ff60405f2054166105b95760206102e1916102c48a611527565b9360405193849283926397aba7f960e01b84528b600485016115c2565b0381865afa9081156104da575f9161058a575b506001600160a01b0390811691160361057b57602060028301546044604051809481936346f6196d60e11b835288600484015260248301525afa9081156104da575f9161054c575b501561053d57600301546040516318bc00fb60e31b8152906020826004817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9081156104da575f91610507575b61039e9250611506565b4080156104f857604051916020830191825260c086013560408401526060830152606082526103ce6080836114a2565b60a08501359151902010156104e9575f52600360205260405f20600160ff198254161790556103fc82611527565b9061040684611527565b608084013592823b156104e55760405163756b7b8960e01b81526001600160a01b0391821660048201529116602482015260448101839052905f908290606490829084905af180156104da576104c4575b5060407f16c45b304a05c8734c1cf2f22e2bb402a1d810dbb2955e87dbcd2c2457457e279161048e61048885611527565b95611527565b825160609590950135855260208501919091526001600160a01b03908116941692a460015f516020611af95f395f51905f525580f35b6104d19195505f906114a2565b5f936040610457565b6040513d5f823e3d90fd5b5f80fd5b638401f37d60e01b5f5260045ffd5b63492aae0760e11b5f5260045ffd5b90506020823d602011610535575b81610522602093836114a2565b810103126104e55761039e915190610394565b3d9150610515565b63abab6bd760e01b5f5260045ffd5b61056e915060203d602011610574575b61056681836114a2565b8101906114d8565b5f61033c565b503d61055c565b638baa579f60e01b5f5260045ffd5b6105ac915060203d6020116105b2575b6105a481836114a2565b8101906115a3565b5f6102f4565b503d61059a565b630869773360e41b5f5260045ffd5b9095506020813d6020116105f4575b816105e4602093836114a2565b810103126104e55751945f6102a0565b3d91506105d7565b630522870960e11b5f5260045ffd5b634e487b7160e01b5f52602160045260245ffd5b9096506020813d60201161064b575b8161063b602093836114a2565b810103126104e55751955f61020e565b3d915061062e565b346104e5575f3660031901126104e5576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b346104e55760203660031901126104e5576004355f52600160205260a060405f206106f6815491600181015490600281015460ff60046003840154930154169260405195865260208601526040850152606084015260808301906113c5565bf35b346104e55761070636611464565b5f54919392909160081c6001600160a01b03168015159081610839575b5061082a576040516311b663f160e31b81526001600160a01b03808516600483015282166024820152602081806044810103817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa80156104da5785915f916107f3575b50916107c06020927fe15b7d92d30e2e0d903eccb550b6fe43024271cc0eca67f900b793307814ddc1946115ef565b5f908152600183526040908190206004908101805460ff19169091179055519384526001600160a01b03908116941692a4005b9150506020813d602011610822575b8161080f602093836114a2565b810103126104e5575184906107c0610791565b3d9150610802565b635c427cd960e01b5f5260045ffd5b905033141585610723565b346104e55760203660031901126104e55761085d61139b565b6004546001600160a01b0316330361082a575f8054610100600160a81b03191660089290921b610100600160a81b0316919091179055005b346104e5575f3660031901126104e5576004546040516001600160a01b039091168152602090f35b346104e5575f3660031901126104e5576004546001600160a01b0316330361082a576108e761164d565b600160ff195f5416175f557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a1005b346104e5575f3660031901126104e5576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b346104e5575f3660031901126104e5576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b346104e55760403660031901126104e5576109bf61139b565b602435906109cb61164d565b6040516311b663f160e31b81526001600160a01b0382811660048301523360248301527f00000000000000000000000000000000000000000000000000000000000000001690602081604481855afa80156104da5784905f90610bc0575b610a3392506115ef565b5f52600160205260405f2090600482019160ff835416600581101561060b57600103610bb157604051630f23262160e01b81526001600160a01b038516600482015233602482015291602090839060449082905afa9182156104da575f92610b8d575b506003015460405163d5beb3c760e01b8152906020826004817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9081156104da575f91610b57575b610af39250611506565b431190159081610b4e575b50610b3f57805460ff1916600317905533906001600160a01b03167f63f65de358a7452715376fc5a74b03fc4344129934500800a3814868042c35975f80a4005b63bc72195f60e01b5f5260045ffd5b90501584610afe565b90506020823d602011610b85575b81610b72602093836114a2565b810103126104e557610af3915190610ae9565b3d9150610b65565b6003919250610baa9060203d6020116105745761056681836114a2565b9190610a96565b63402bc00760e01b5f5260045ffd5b50506020813d602011610bee575b81610bdb602093836114a2565b810103126104e55783610a339151610a29565b3d9150610bce565b346104e5575f3660031901126104e557602060ff5f54166040519015158152f35b346104e55760203660031901126104e5576004355f526002602052602060ff60405f2054166040519015158152f35b346104e5575f3660031901126104e5576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b346104e5576020610cd1610c9d36611464565b6040516311b663f160e31b81526001600160a01b039485166004820152939092166024840152939092829081906044820190565b03817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa80156104da5783905f90610e58575b610d1892506115ef565b5f52600160205260405f209160ff600484015416600581101561060b5715610e495760028301546040516346f6196d60e11b81526004810184905260248101919091527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690602081604481855afa9081156104da575f91610e2a575b501561053d57600384015493546040805163cbdd253960e01b81526004810196909652602486019490945260448501929092526064840191909152829060849082905afa80156104da576020915f91610dfa575b50604051908152f35b610e1c915060403d604011610e23575b610e1481836114a2565b8101906114f0565b5082610df1565b503d610e0a565b610e43915060203d6020116105745761056681836114a2565b85610d9d565b635780553760e11b5f5260045ffd5b50506020813d602011610e86575b81610e73602093836114a2565b810103126104e55782610d189151610d0e565b3d9150610e66565b346104e55760a03660031901126104e557610ea761139b565b60243590604435606435610eb961164d565b81158015611020575b8015611018575b611009576040516311b663f160e31b81526001600160a01b0384166004820152336024820152602081806044810103817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa80156104da5785905f90610fd3575b610f3e92506115ef565b5f52600160205260405f20600481019060ff825416600581101561060b57610fc4578381558260018201556084356002820155600343910155600160ff1982541617905560405191825260208201524360408201527f37f4a67788af92647c831e67034db177a77fe8b0fe306d329422e139bc0c4acb6060339360018060a01b031692a4005b630cebf7ad60e41b5f5260045ffd5b50506020813d602011611001575b81610fee602093836114a2565b810103126104e55784610f3e9151610f34565b3d9150610fe1565b6331f4a25760e11b5f5260045ffd5b508015610ec9565b5060405163865ef51360e01b81526020816004817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa9081156104da575f91611076575b508211610ec2565b90506020813d6020116110a0575b81611091602093836114a2565b810103126104e557518561106e565b3d9150611084565b346104e5575f3660031901126104e5576004546001600160a01b0316330361082a575f5460ff8116156111075760ff19165f557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a1005b638dfc202b60e01b5f5260045ffd5b346104e55761114061112736611403565b92611133929192611615565b61113b61164d565b611667565b60015f516020611af95f395f51905f5255005b346104e55760603660031901126104e55760043567ffffffffffffffff81116104e557366023820112156104e55780600401359067ffffffffffffffff82116104e55736602460e08402830101116104e55760243567ffffffffffffffff81116104e5576111c59036906004016113d2565b9160443567ffffffffffffffff81116104e5576111e69036906004016113d2565b90946111f0611615565b6111f861164d565b84811480159061129f575b611290579336849003601e1901905f5b86811015611140578181101561127c578060051b9081870135848112156104e55787019182359267ffffffffffffffff84116104e55760200183360381136104e5578683101561127c57600193611276928c013591602460e086028b0101611667565b01611213565b634e487b7160e01b5f52603260045260245ffd5b63512509d360e11b5f5260045ffd5b5081811415611203565b346104e55760603660031901126104e5576112c261139b565b602435906001600160a01b03821682036104e5576311b663f160e31b83526001600160a01b039081166004840152166024820152602081806044810103817f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165afa80156104da575f90611367575b61134690604435906115ef565b5f526001602052602060ff600460405f200154166106f660405180926113c5565b506020813d602011611393575b81611381602093836114a2565b810103126104e5576113469051611339565b3d9150611374565b600435906001600160a01b03821682036104e557565b35906001600160a01b03821682036104e557565b90600582101561060b5752565b9181601f840112156104e55782359167ffffffffffffffff83116104e5576020808501948460051b0101116104e557565b600319810161012081126104e55760e0136104e55760049160e43567ffffffffffffffff81116104e557826023820112156104e55780600401359267ffffffffffffffff84116104e557602484830101116104e55760240191906101043590565b60809060031901126104e5576004356001600160a01b03811681036104e557906024356001600160a01b03811681036104e557906044359060643590565b90601f8019910116810190811067ffffffffffffffff8211176114c457604052565b634e487b7160e01b5f52604160045260245ffd5b908160209103126104e5575180151581036104e55790565b91908260409103126104e5576020825192015190565b9190820180921161151357565b634e487b7160e01b5f52601160045260245ffd5b356001600160a01b03811681036104e55790565b60e08101929160c09081906001600160a01b03611557826113b1565b1684526001600160a01b0361156e602083016113b1565b16602085015260408101356040850152606081013560608501526080810135608085015260a081013560a08501520135910152565b908160209103126104e557516001600160a01b03811681036104e55790565b91926060938192845260406020850152816040850152848401375f828201840152601f01601f1916010190565b90604051906020820192835260408201526040815261160f6060826114a2565b51902090565b60025f516020611af95f395f51905f52541461163e5760025f516020611af95f395f51905f5255565b633ee5aeb560e01b5f5260045ffd5b60ff5f541661165857565b63d93c066560e01b5f5260045ffd5b6116b3939192917f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031691906116a382611527565b946020808401966101cb88611527565b0381875afa9687156104da575f97611ac4575b506116d6604084013580986115ef565b5f52600160205260405f2094600486019160ff835416600581101561060b578015610e495760048114611ab5575f1901610bb1575f54899060081c6001600160a01b031680611a14575b505060808501359360018801548514801590611a05575b61100957604051630bd5b99b60e41b8152927f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169290602085806117868b6004830161153b565b0381875afa9485156104da575f956119d1575b50845f52600260205260ff60405f2054166105b95760206117da916117bd8a611527565b9360405193849283926397aba7f960e01b84528a600485016115c2565b0381875afa9081156104da575f916119b2575b506001600160a01b0390811691160361057b576002880154604051906346f6196d60e11b82528260048301526024820152602081604481865afa9081156104da575f91611993575b501561053d57600388015497546040805163cbdd253960e01b8152600481019a909a5260248a0192909252604489018b9052606489015290879060849082905afa9586156104da575f96611971575b5060608401359586036104e9575f52600260205260405f20600160ff19825416179055600260ff198254161790556118bb82611527565b6118c486611527565b90843b156104e55760405163756b7b8960e01b81526001600160a01b0391821660048201529116602482015260448101829052925f908490606490829084905af19283156104da57610488604093611943927f0b76e789e4341a8a74776cbf3af9cfc0df7246d30f52a25462d22762a0f346c596611961575b50611527565b825194855260208501919091526001600160a01b03908116941692a4565b5f61196b916114a2565b5f61193d565b61198b91965060403d604011610e2357610e1481836114a2565b50945f611884565b6119ac915060203d6020116105745761056681836114a2565b5f611835565b6119cb915060203d6020116105b2576105a481836114a2565b5f6117ed565b9094506020813d6020116119fd575b816119ed602093836114a2565b810103126104e55751935f611799565b3d91506119e0565b50875460608701351015611737565b611a6991602091611a2489611527565b611a2d8d611527565b6040516318387cc760e31b81526001600160a01b0392831660048201529116602482015260448101929092529092839190829081906064820190565b03915afa9081156104da575f91611a96575b50611a8757885f611720565b6307afee5f60e21b5f5260045ffd5b611aaf915060203d6020116105745761056681836114a2565b5f611a7b565b630e61107b60e11b5f5260045ffd5b9096506020813d602011611af0575b81611ae0602093836114a2565b810103126104e55751955f6116c6565b3d9150611ad356fe9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00a2646970667358221220ad81c101e3048bd9cf2d3ffdd9cd680ceb17c1aa272b354d537e8fc8ca565d3664736f6c63430008230033",
}

// SettlementEngineABI is the input ABI used to generate the binding from.
// Deprecated: Use SettlementEngineMetaData.ABI instead.
var SettlementEngineABI = SettlementEngineMetaData.ABI

// SettlementEngineBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SettlementEngineMetaData.Bin instead.
var SettlementEngineBin = SettlementEngineMetaData.Bin

// DeploySettlementEngine deploys a new Ethereum contract, binding an instance of SettlementEngine to it.
func DeploySettlementEngine(auth *bind.TransactOpts, backend bind.ContractBackend, _config common.Address, _entropySource common.Address, _paymentChannel common.Address, _providerRegistry common.Address) (common.Address, *types.Transaction, *SettlementEngine, error) {
	parsed, err := SettlementEngineMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SettlementEngineBin), backend, _config, _entropySource, _paymentChannel, _providerRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SettlementEngine{SettlementEngineCaller: SettlementEngineCaller{contract: contract}, SettlementEngineTransactor: SettlementEngineTransactor{contract: contract}, SettlementEngineFilterer: SettlementEngineFilterer{contract: contract}}, nil
}

// SettlementEngine is an auto generated Go binding around an Ethereum contract.
type SettlementEngine struct {
	SettlementEngineCaller     // Read-only binding to the contract
	SettlementEngineTransactor // Write-only binding to the contract
	SettlementEngineFilterer   // Log filterer for contract events
}

// SettlementEngineCaller is an auto generated read-only Go binding around an Ethereum contract.
type SettlementEngineCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementEngineTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SettlementEngineTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementEngineFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SettlementEngineFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementEngineSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SettlementEngineSession struct {
	Contract     *SettlementEngine // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SettlementEngineCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SettlementEngineCallerSession struct {
	Contract *SettlementEngineCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// SettlementEngineTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SettlementEngineTransactorSession struct {
	Contract     *SettlementEngineTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// SettlementEngineRaw is an auto generated low-level Go binding around an Ethereum contract.
type SettlementEngineRaw struct {
	Contract *SettlementEngine // Generic contract binding to access the raw methods on
}

// SettlementEngineCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SettlementEngineCallerRaw struct {
	Contract *SettlementEngineCaller // Generic read-only contract binding to access the raw methods on
}

// SettlementEngineTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SettlementEngineTransactorRaw struct {
	Contract *SettlementEngineTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSettlementEngine creates a new instance of SettlementEngine, bound to a specific deployed contract.
func NewSettlementEngine(address common.Address, backend bind.ContractBackend) (*SettlementEngine, error) {
	contract, err := bindSettlementEngine(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SettlementEngine{SettlementEngineCaller: SettlementEngineCaller{contract: contract}, SettlementEngineTransactor: SettlementEngineTransactor{contract: contract}, SettlementEngineFilterer: SettlementEngineFilterer{contract: contract}}, nil
}

// NewSettlementEngineCaller creates a new read-only instance of SettlementEngine, bound to a specific deployed contract.
func NewSettlementEngineCaller(address common.Address, caller bind.ContractCaller) (*SettlementEngineCaller, error) {
	contract, err := bindSettlementEngine(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementEngineCaller{contract: contract}, nil
}

// NewSettlementEngineTransactor creates a new write-only instance of SettlementEngine, bound to a specific deployed contract.
func NewSettlementEngineTransactor(address common.Address, transactor bind.ContractTransactor) (*SettlementEngineTransactor, error) {
	contract, err := bindSettlementEngine(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementEngineTransactor{contract: contract}, nil
}

// NewSettlementEngineFilterer creates a new log filterer instance of SettlementEngine, bound to a specific deployed contract.
func NewSettlementEngineFilterer(address common.Address, filterer bind.ContractFilterer) (*SettlementEngineFilterer, error) {
	contract, err := bindSettlementEngine(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SettlementEngineFilterer{contract: contract}, nil
}

// bindSettlementEngine binds a generic wrapper to an already deployed contract.
func bindSettlementEngine(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SettlementEngineMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementEngine *SettlementEngineRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SettlementEngine.Contract.SettlementEngineCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementEngine *SettlementEngineRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementEngine.Contract.SettlementEngineTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementEngine *SettlementEngineRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementEngine.Contract.SettlementEngineTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SettlementEngine *SettlementEngineCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SettlementEngine.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SettlementEngine *SettlementEngineTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementEngine.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SettlementEngine *SettlementEngineTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SettlementEngine.Contract.contract.Transact(opts, method, params...)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_SettlementEngine *SettlementEngineCaller) Config(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "config")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_SettlementEngine *SettlementEngineSession) Config() (common.Address, error) {
	return _SettlementEngine.Contract.Config(&_SettlementEngine.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(address)
func (_SettlementEngine *SettlementEngineCallerSession) Config() (common.Address, error) {
	return _SettlementEngine.Contract.Config(&_SettlementEngine.CallOpts)
}

// DisputeResolver is a free data retrieval call binding the contract method 0xf5a3f4af.
//
// Solidity: function disputeResolver() view returns(address)
func (_SettlementEngine *SettlementEngineCaller) DisputeResolver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "disputeResolver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DisputeResolver is a free data retrieval call binding the contract method 0xf5a3f4af.
//
// Solidity: function disputeResolver() view returns(address)
func (_SettlementEngine *SettlementEngineSession) DisputeResolver() (common.Address, error) {
	return _SettlementEngine.Contract.DisputeResolver(&_SettlementEngine.CallOpts)
}

// DisputeResolver is a free data retrieval call binding the contract method 0xf5a3f4af.
//
// Solidity: function disputeResolver() view returns(address)
func (_SettlementEngine *SettlementEngineCallerSession) DisputeResolver() (common.Address, error) {
	return _SettlementEngine.Contract.DisputeResolver(&_SettlementEngine.CallOpts)
}

// EntropySource is a free data retrieval call binding the contract method 0xae174b9f.
//
// Solidity: function entropySource() view returns(address)
func (_SettlementEngine *SettlementEngineCaller) EntropySource(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "entropySource")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EntropySource is a free data retrieval call binding the contract method 0xae174b9f.
//
// Solidity: function entropySource() view returns(address)
func (_SettlementEngine *SettlementEngineSession) EntropySource() (common.Address, error) {
	return _SettlementEngine.Contract.EntropySource(&_SettlementEngine.CallOpts)
}

// EntropySource is a free data retrieval call binding the contract method 0xae174b9f.
//
// Solidity: function entropySource() view returns(address)
func (_SettlementEngine *SettlementEngineCallerSession) EntropySource() (common.Address, error) {
	return _SettlementEngine.Contract.EntropySource(&_SettlementEngine.CallOpts)
}

// GetRoundStatus is a free data retrieval call binding the contract method 0x232ffabb.
//
// Solidity: function getRoundStatus(address sender, address recipient, uint256 roundId) view returns(uint8)
func (_SettlementEngine *SettlementEngineCaller) GetRoundStatus(opts *bind.CallOpts, sender common.Address, recipient common.Address, roundId *big.Int) (uint8, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "getRoundStatus", sender, recipient, roundId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetRoundStatus is a free data retrieval call binding the contract method 0x232ffabb.
//
// Solidity: function getRoundStatus(address sender, address recipient, uint256 roundId) view returns(uint8)
func (_SettlementEngine *SettlementEngineSession) GetRoundStatus(sender common.Address, recipient common.Address, roundId *big.Int) (uint8, error) {
	return _SettlementEngine.Contract.GetRoundStatus(&_SettlementEngine.CallOpts, sender, recipient, roundId)
}

// GetRoundStatus is a free data retrieval call binding the contract method 0x232ffabb.
//
// Solidity: function getRoundStatus(address sender, address recipient, uint256 roundId) view returns(uint8)
func (_SettlementEngine *SettlementEngineCallerSession) GetRoundStatus(sender common.Address, recipient common.Address, roundId *big.Int) (uint8, error) {
	return _SettlementEngine.Contract.GetRoundStatus(&_SettlementEngine.CallOpts, sender, recipient, roundId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SettlementEngine *SettlementEngineCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SettlementEngine *SettlementEngineSession) Owner() (common.Address, error) {
	return _SettlementEngine.Contract.Owner(&_SettlementEngine.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SettlementEngine *SettlementEngineCallerSession) Owner() (common.Address, error) {
	return _SettlementEngine.Contract.Owner(&_SettlementEngine.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SettlementEngine *SettlementEngineCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SettlementEngine *SettlementEngineSession) Paused() (bool, error) {
	return _SettlementEngine.Contract.Paused(&_SettlementEngine.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SettlementEngine *SettlementEngineCallerSession) Paused() (bool, error) {
	return _SettlementEngine.Contract.Paused(&_SettlementEngine.CallOpts)
}

// PaymentChannel is a free data retrieval call binding the contract method 0x6df24cd9.
//
// Solidity: function paymentChannel() view returns(address)
func (_SettlementEngine *SettlementEngineCaller) PaymentChannel(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "paymentChannel")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PaymentChannel is a free data retrieval call binding the contract method 0x6df24cd9.
//
// Solidity: function paymentChannel() view returns(address)
func (_SettlementEngine *SettlementEngineSession) PaymentChannel() (common.Address, error) {
	return _SettlementEngine.Contract.PaymentChannel(&_SettlementEngine.CallOpts)
}

// PaymentChannel is a free data retrieval call binding the contract method 0x6df24cd9.
//
// Solidity: function paymentChannel() view returns(address)
func (_SettlementEngine *SettlementEngineCallerSession) PaymentChannel() (common.Address, error) {
	return _SettlementEngine.Contract.PaymentChannel(&_SettlementEngine.CallOpts)
}

// PreviewWinnerIndex is a free data retrieval call binding the contract method 0x51110b65.
//
// Solidity: function previewWinnerIndex(address sender, address recipient, uint256 roundId, bytes32 secret) view returns(uint256)
func (_SettlementEngine *SettlementEngineCaller) PreviewWinnerIndex(opts *bind.CallOpts, sender common.Address, recipient common.Address, roundId *big.Int, secret [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "previewWinnerIndex", sender, recipient, roundId, secret)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewWinnerIndex is a free data retrieval call binding the contract method 0x51110b65.
//
// Solidity: function previewWinnerIndex(address sender, address recipient, uint256 roundId, bytes32 secret) view returns(uint256)
func (_SettlementEngine *SettlementEngineSession) PreviewWinnerIndex(sender common.Address, recipient common.Address, roundId *big.Int, secret [32]byte) (*big.Int, error) {
	return _SettlementEngine.Contract.PreviewWinnerIndex(&_SettlementEngine.CallOpts, sender, recipient, roundId, secret)
}

// PreviewWinnerIndex is a free data retrieval call binding the contract method 0x51110b65.
//
// Solidity: function previewWinnerIndex(address sender, address recipient, uint256 roundId, bytes32 secret) view returns(uint256)
func (_SettlementEngine *SettlementEngineCallerSession) PreviewWinnerIndex(sender common.Address, recipient common.Address, roundId *big.Int, secret [32]byte) (*big.Int, error) {
	return _SettlementEngine.Contract.PreviewWinnerIndex(&_SettlementEngine.CallOpts, sender, recipient, roundId, secret)
}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_SettlementEngine *SettlementEngineCaller) ProviderRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "providerRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_SettlementEngine *SettlementEngineSession) ProviderRegistry() (common.Address, error) {
	return _SettlementEngine.Contract.ProviderRegistry(&_SettlementEngine.CallOpts)
}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_SettlementEngine *SettlementEngineCallerSession) ProviderRegistry() (common.Address, error) {
	return _SettlementEngine.Contract.ProviderRegistry(&_SettlementEngine.CallOpts)
}

// Rounds is a free data retrieval call binding the contract method 0xa23ad88d.
//
// Solidity: function rounds(bytes32 ) view returns(uint256 tau, uint256 faceValue, bytes32 recipientRandHash, uint256 commitBlock, uint8 status)
func (_SettlementEngine *SettlementEngineCaller) Rounds(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Tau               *big.Int
	FaceValue         *big.Int
	RecipientRandHash [32]byte
	CommitBlock       *big.Int
	Status            uint8
}, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "rounds", arg0)

	outstruct := new(struct {
		Tau               *big.Int
		FaceValue         *big.Int
		RecipientRandHash [32]byte
		CommitBlock       *big.Int
		Status            uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Tau = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FaceValue = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.RecipientRandHash = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.CommitBlock = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Status = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

// Rounds is a free data retrieval call binding the contract method 0xa23ad88d.
//
// Solidity: function rounds(bytes32 ) view returns(uint256 tau, uint256 faceValue, bytes32 recipientRandHash, uint256 commitBlock, uint8 status)
func (_SettlementEngine *SettlementEngineSession) Rounds(arg0 [32]byte) (struct {
	Tau               *big.Int
	FaceValue         *big.Int
	RecipientRandHash [32]byte
	CommitBlock       *big.Int
	Status            uint8
}, error) {
	return _SettlementEngine.Contract.Rounds(&_SettlementEngine.CallOpts, arg0)
}

// Rounds is a free data retrieval call binding the contract method 0xa23ad88d.
//
// Solidity: function rounds(bytes32 ) view returns(uint256 tau, uint256 faceValue, bytes32 recipientRandHash, uint256 commitBlock, uint8 status)
func (_SettlementEngine *SettlementEngineCallerSession) Rounds(arg0 [32]byte) (struct {
	Tau               *big.Int
	FaceValue         *big.Int
	RecipientRandHash [32]byte
	CommitBlock       *big.Int
	Status            uint8
}, error) {
	return _SettlementEngine.Contract.Rounds(&_SettlementEngine.CallOpts, arg0)
}

// UsedFallbackTickets is a free data retrieval call binding the contract method 0xd0a3155e.
//
// Solidity: function usedFallbackTickets(bytes32 ) view returns(bool)
func (_SettlementEngine *SettlementEngineCaller) UsedFallbackTickets(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "usedFallbackTickets", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedFallbackTickets is a free data retrieval call binding the contract method 0xd0a3155e.
//
// Solidity: function usedFallbackTickets(bytes32 ) view returns(bool)
func (_SettlementEngine *SettlementEngineSession) UsedFallbackTickets(arg0 [32]byte) (bool, error) {
	return _SettlementEngine.Contract.UsedFallbackTickets(&_SettlementEngine.CallOpts, arg0)
}

// UsedFallbackTickets is a free data retrieval call binding the contract method 0xd0a3155e.
//
// Solidity: function usedFallbackTickets(bytes32 ) view returns(bool)
func (_SettlementEngine *SettlementEngineCallerSession) UsedFallbackTickets(arg0 [32]byte) (bool, error) {
	return _SettlementEngine.Contract.UsedFallbackTickets(&_SettlementEngine.CallOpts, arg0)
}

// UsedTickets is a free data retrieval call binding the contract method 0x59a515ba.
//
// Solidity: function usedTickets(bytes32 ) view returns(bool)
func (_SettlementEngine *SettlementEngineCaller) UsedTickets(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _SettlementEngine.contract.Call(opts, &out, "usedTickets", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedTickets is a free data retrieval call binding the contract method 0x59a515ba.
//
// Solidity: function usedTickets(bytes32 ) view returns(bool)
func (_SettlementEngine *SettlementEngineSession) UsedTickets(arg0 [32]byte) (bool, error) {
	return _SettlementEngine.Contract.UsedTickets(&_SettlementEngine.CallOpts, arg0)
}

// UsedTickets is a free data retrieval call binding the contract method 0x59a515ba.
//
// Solidity: function usedTickets(bytes32 ) view returns(bool)
func (_SettlementEngine *SettlementEngineCallerSession) UsedTickets(arg0 [32]byte) (bool, error) {
	return _SettlementEngine.Contract.UsedTickets(&_SettlementEngine.CallOpts, arg0)
}

// ClaimFallbackTicket is a paid mutator transaction binding the contract method 0xbb71decb.
//
// Solidity: function claimFallbackTicket((address,address,uint256,uint256,uint256,uint256,uint256) ticket, bytes senderSig, bytes32 secret) returns()
func (_SettlementEngine *SettlementEngineTransactor) ClaimFallbackTicket(opts *bind.TransactOpts, ticket IEntropySourceRoundTicket, senderSig []byte, secret [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.contract.Transact(opts, "claimFallbackTicket", ticket, senderSig, secret)
}

// ClaimFallbackTicket is a paid mutator transaction binding the contract method 0xbb71decb.
//
// Solidity: function claimFallbackTicket((address,address,uint256,uint256,uint256,uint256,uint256) ticket, bytes senderSig, bytes32 secret) returns()
func (_SettlementEngine *SettlementEngineSession) ClaimFallbackTicket(ticket IEntropySourceRoundTicket, senderSig []byte, secret [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.Contract.ClaimFallbackTicket(&_SettlementEngine.TransactOpts, ticket, senderSig, secret)
}

// ClaimFallbackTicket is a paid mutator transaction binding the contract method 0xbb71decb.
//
// Solidity: function claimFallbackTicket((address,address,uint256,uint256,uint256,uint256,uint256) ticket, bytes senderSig, bytes32 secret) returns()
func (_SettlementEngine *SettlementEngineTransactorSession) ClaimFallbackTicket(ticket IEntropySourceRoundTicket, senderSig []byte, secret [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.Contract.ClaimFallbackTicket(&_SettlementEngine.TransactOpts, ticket, senderSig, secret)
}

// CommitRoundSecret is a paid mutator transaction binding the contract method 0x42bf1beb.
//
// Solidity: function commitRoundSecret(address sender, uint256 roundId, uint256 tau, uint256 faceValue, bytes32 recipientRandHash) returns()
func (_SettlementEngine *SettlementEngineTransactor) CommitRoundSecret(opts *bind.TransactOpts, sender common.Address, roundId *big.Int, tau *big.Int, faceValue *big.Int, recipientRandHash [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.contract.Transact(opts, "commitRoundSecret", sender, roundId, tau, faceValue, recipientRandHash)
}

// CommitRoundSecret is a paid mutator transaction binding the contract method 0x42bf1beb.
//
// Solidity: function commitRoundSecret(address sender, uint256 roundId, uint256 tau, uint256 faceValue, bytes32 recipientRandHash) returns()
func (_SettlementEngine *SettlementEngineSession) CommitRoundSecret(sender common.Address, roundId *big.Int, tau *big.Int, faceValue *big.Int, recipientRandHash [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.Contract.CommitRoundSecret(&_SettlementEngine.TransactOpts, sender, roundId, tau, faceValue, recipientRandHash)
}

// CommitRoundSecret is a paid mutator transaction binding the contract method 0x42bf1beb.
//
// Solidity: function commitRoundSecret(address sender, uint256 roundId, uint256 tau, uint256 faceValue, bytes32 recipientRandHash) returns()
func (_SettlementEngine *SettlementEngineTransactorSession) CommitRoundSecret(sender common.Address, roundId *big.Int, tau *big.Int, faceValue *big.Int, recipientRandHash [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.Contract.CommitRoundSecret(&_SettlementEngine.TransactOpts, sender, roundId, tau, faceValue, recipientRandHash)
}

// MarkRoundPartial is a paid mutator transaction binding the contract method 0x5da3d4f8.
//
// Solidity: function markRoundPartial(address sender, uint256 roundId) returns()
func (_SettlementEngine *SettlementEngineTransactor) MarkRoundPartial(opts *bind.TransactOpts, sender common.Address, roundId *big.Int) (*types.Transaction, error) {
	return _SettlementEngine.contract.Transact(opts, "markRoundPartial", sender, roundId)
}

// MarkRoundPartial is a paid mutator transaction binding the contract method 0x5da3d4f8.
//
// Solidity: function markRoundPartial(address sender, uint256 roundId) returns()
func (_SettlementEngine *SettlementEngineSession) MarkRoundPartial(sender common.Address, roundId *big.Int) (*types.Transaction, error) {
	return _SettlementEngine.Contract.MarkRoundPartial(&_SettlementEngine.TransactOpts, sender, roundId)
}

// MarkRoundPartial is a paid mutator transaction binding the contract method 0x5da3d4f8.
//
// Solidity: function markRoundPartial(address sender, uint256 roundId) returns()
func (_SettlementEngine *SettlementEngineTransactorSession) MarkRoundPartial(sender common.Address, roundId *big.Int) (*types.Transaction, error) {
	return _SettlementEngine.Contract.MarkRoundPartial(&_SettlementEngine.TransactOpts, sender, roundId)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SettlementEngine *SettlementEngineTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementEngine.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SettlementEngine *SettlementEngineSession) Pause() (*types.Transaction, error) {
	return _SettlementEngine.Contract.Pause(&_SettlementEngine.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SettlementEngine *SettlementEngineTransactorSession) Pause() (*types.Transaction, error) {
	return _SettlementEngine.Contract.Pause(&_SettlementEngine.TransactOpts)
}

// SetDisputeResolver is a paid mutator transaction binding the contract method 0x924e63f6.
//
// Solidity: function setDisputeResolver(address _disputeResolver) returns()
func (_SettlementEngine *SettlementEngineTransactor) SetDisputeResolver(opts *bind.TransactOpts, _disputeResolver common.Address) (*types.Transaction, error) {
	return _SettlementEngine.contract.Transact(opts, "setDisputeResolver", _disputeResolver)
}

// SetDisputeResolver is a paid mutator transaction binding the contract method 0x924e63f6.
//
// Solidity: function setDisputeResolver(address _disputeResolver) returns()
func (_SettlementEngine *SettlementEngineSession) SetDisputeResolver(_disputeResolver common.Address) (*types.Transaction, error) {
	return _SettlementEngine.Contract.SetDisputeResolver(&_SettlementEngine.TransactOpts, _disputeResolver)
}

// SetDisputeResolver is a paid mutator transaction binding the contract method 0x924e63f6.
//
// Solidity: function setDisputeResolver(address _disputeResolver) returns()
func (_SettlementEngine *SettlementEngineTransactorSession) SetDisputeResolver(_disputeResolver common.Address) (*types.Transaction, error) {
	return _SettlementEngine.Contract.SetDisputeResolver(&_SettlementEngine.TransactOpts, _disputeResolver)
}

// SettleRound is a paid mutator transaction binding the contract method 0x3e7678b4.
//
// Solidity: function settleRound((address,address,uint256,uint256,uint256,uint256,uint256) ticket, bytes senderSig, bytes32 secret) returns()
func (_SettlementEngine *SettlementEngineTransactor) SettleRound(opts *bind.TransactOpts, ticket IEntropySourceRoundTicket, senderSig []byte, secret [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.contract.Transact(opts, "settleRound", ticket, senderSig, secret)
}

// SettleRound is a paid mutator transaction binding the contract method 0x3e7678b4.
//
// Solidity: function settleRound((address,address,uint256,uint256,uint256,uint256,uint256) ticket, bytes senderSig, bytes32 secret) returns()
func (_SettlementEngine *SettlementEngineSession) SettleRound(ticket IEntropySourceRoundTicket, senderSig []byte, secret [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.Contract.SettleRound(&_SettlementEngine.TransactOpts, ticket, senderSig, secret)
}

// SettleRound is a paid mutator transaction binding the contract method 0x3e7678b4.
//
// Solidity: function settleRound((address,address,uint256,uint256,uint256,uint256,uint256) ticket, bytes senderSig, bytes32 secret) returns()
func (_SettlementEngine *SettlementEngineTransactorSession) SettleRound(ticket IEntropySourceRoundTicket, senderSig []byte, secret [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.Contract.SettleRound(&_SettlementEngine.TransactOpts, ticket, senderSig, secret)
}

// SettleRoundBatch is a paid mutator transaction binding the contract method 0x3757a3da.
//
// Solidity: function settleRoundBatch((address,address,uint256,uint256,uint256,uint256,uint256)[] tickets, bytes[] senderSigs, bytes32[] secrets) returns()
func (_SettlementEngine *SettlementEngineTransactor) SettleRoundBatch(opts *bind.TransactOpts, tickets []IEntropySourceRoundTicket, senderSigs [][]byte, secrets [][32]byte) (*types.Transaction, error) {
	return _SettlementEngine.contract.Transact(opts, "settleRoundBatch", tickets, senderSigs, secrets)
}

// SettleRoundBatch is a paid mutator transaction binding the contract method 0x3757a3da.
//
// Solidity: function settleRoundBatch((address,address,uint256,uint256,uint256,uint256,uint256)[] tickets, bytes[] senderSigs, bytes32[] secrets) returns()
func (_SettlementEngine *SettlementEngineSession) SettleRoundBatch(tickets []IEntropySourceRoundTicket, senderSigs [][]byte, secrets [][32]byte) (*types.Transaction, error) {
	return _SettlementEngine.Contract.SettleRoundBatch(&_SettlementEngine.TransactOpts, tickets, senderSigs, secrets)
}

// SettleRoundBatch is a paid mutator transaction binding the contract method 0x3757a3da.
//
// Solidity: function settleRoundBatch((address,address,uint256,uint256,uint256,uint256,uint256)[] tickets, bytes[] senderSigs, bytes32[] secrets) returns()
func (_SettlementEngine *SettlementEngineTransactorSession) SettleRoundBatch(tickets []IEntropySourceRoundTicket, senderSigs [][]byte, secrets [][32]byte) (*types.Transaction, error) {
	return _SettlementEngine.Contract.SettleRoundBatch(&_SettlementEngine.TransactOpts, tickets, senderSigs, secrets)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SettlementEngine *SettlementEngineTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SettlementEngine.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SettlementEngine *SettlementEngineSession) Unpause() (*types.Transaction, error) {
	return _SettlementEngine.Contract.Unpause(&_SettlementEngine.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SettlementEngine *SettlementEngineTransactorSession) Unpause() (*types.Transaction, error) {
	return _SettlementEngine.Contract.Unpause(&_SettlementEngine.TransactOpts)
}

// VoidRound is a paid mutator transaction binding the contract method 0x9c390612.
//
// Solidity: function voidRound(address sender, address recipient, uint256 roundId, bytes32 reason) returns()
func (_SettlementEngine *SettlementEngineTransactor) VoidRound(opts *bind.TransactOpts, sender common.Address, recipient common.Address, roundId *big.Int, reason [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.contract.Transact(opts, "voidRound", sender, recipient, roundId, reason)
}

// VoidRound is a paid mutator transaction binding the contract method 0x9c390612.
//
// Solidity: function voidRound(address sender, address recipient, uint256 roundId, bytes32 reason) returns()
func (_SettlementEngine *SettlementEngineSession) VoidRound(sender common.Address, recipient common.Address, roundId *big.Int, reason [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.Contract.VoidRound(&_SettlementEngine.TransactOpts, sender, recipient, roundId, reason)
}

// VoidRound is a paid mutator transaction binding the contract method 0x9c390612.
//
// Solidity: function voidRound(address sender, address recipient, uint256 roundId, bytes32 reason) returns()
func (_SettlementEngine *SettlementEngineTransactorSession) VoidRound(sender common.Address, recipient common.Address, roundId *big.Int, reason [32]byte) (*types.Transaction, error) {
	return _SettlementEngine.Contract.VoidRound(&_SettlementEngine.TransactOpts, sender, recipient, roundId, reason)
}

// SettlementEngineFallbackTicketClaimedIterator is returned from FilterFallbackTicketClaimed and is used to iterate over the raw logs and unpacked data for FallbackTicketClaimed events raised by the SettlementEngine contract.
type SettlementEngineFallbackTicketClaimedIterator struct {
	Event *SettlementEngineFallbackTicketClaimed // Event containing the contract specifics and raw log

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
func (it *SettlementEngineFallbackTicketClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementEngineFallbackTicketClaimed)
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
		it.Event = new(SettlementEngineFallbackTicketClaimed)
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
func (it *SettlementEngineFallbackTicketClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementEngineFallbackTicketClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementEngineFallbackTicketClaimed represents a FallbackTicketClaimed event raised by the SettlementEngine contract.
type SettlementEngineFallbackTicketClaimed struct {
	Sender     common.Address
	Recipient  common.Address
	RoundId    *big.Int
	LocalIndex *big.Int
	FaceValue  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterFallbackTicketClaimed is a free log retrieval operation binding the contract event 0x16c45b304a05c8734c1cf2f22e2bb402a1d810dbb2955e87dbcd2c2457457e27.
//
// Solidity: event FallbackTicketClaimed(address indexed sender, address indexed recipient, uint256 indexed roundId, uint256 localIndex, uint256 faceValue)
func (_SettlementEngine *SettlementEngineFilterer) FilterFallbackTicketClaimed(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, roundId []*big.Int) (*SettlementEngineFallbackTicketClaimedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _SettlementEngine.contract.FilterLogs(opts, "FallbackTicketClaimed", senderRule, recipientRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &SettlementEngineFallbackTicketClaimedIterator{contract: _SettlementEngine.contract, event: "FallbackTicketClaimed", logs: logs, sub: sub}, nil
}

// WatchFallbackTicketClaimed is a free log subscription operation binding the contract event 0x16c45b304a05c8734c1cf2f22e2bb402a1d810dbb2955e87dbcd2c2457457e27.
//
// Solidity: event FallbackTicketClaimed(address indexed sender, address indexed recipient, uint256 indexed roundId, uint256 localIndex, uint256 faceValue)
func (_SettlementEngine *SettlementEngineFilterer) WatchFallbackTicketClaimed(opts *bind.WatchOpts, sink chan<- *SettlementEngineFallbackTicketClaimed, sender []common.Address, recipient []common.Address, roundId []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _SettlementEngine.contract.WatchLogs(opts, "FallbackTicketClaimed", senderRule, recipientRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementEngineFallbackTicketClaimed)
				if err := _SettlementEngine.contract.UnpackLog(event, "FallbackTicketClaimed", log); err != nil {
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

// ParseFallbackTicketClaimed is a log parse operation binding the contract event 0x16c45b304a05c8734c1cf2f22e2bb402a1d810dbb2955e87dbcd2c2457457e27.
//
// Solidity: event FallbackTicketClaimed(address indexed sender, address indexed recipient, uint256 indexed roundId, uint256 localIndex, uint256 faceValue)
func (_SettlementEngine *SettlementEngineFilterer) ParseFallbackTicketClaimed(log types.Log) (*SettlementEngineFallbackTicketClaimed, error) {
	event := new(SettlementEngineFallbackTicketClaimed)
	if err := _SettlementEngine.contract.UnpackLog(event, "FallbackTicketClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementEnginePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the SettlementEngine contract.
type SettlementEnginePausedIterator struct {
	Event *SettlementEnginePaused // Event containing the contract specifics and raw log

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
func (it *SettlementEnginePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementEnginePaused)
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
		it.Event = new(SettlementEnginePaused)
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
func (it *SettlementEnginePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementEnginePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementEnginePaused represents a Paused event raised by the SettlementEngine contract.
type SettlementEnginePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SettlementEngine *SettlementEngineFilterer) FilterPaused(opts *bind.FilterOpts) (*SettlementEnginePausedIterator, error) {

	logs, sub, err := _SettlementEngine.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &SettlementEnginePausedIterator{contract: _SettlementEngine.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SettlementEngine *SettlementEngineFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *SettlementEnginePaused) (event.Subscription, error) {

	logs, sub, err := _SettlementEngine.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementEnginePaused)
				if err := _SettlementEngine.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_SettlementEngine *SettlementEngineFilterer) ParsePaused(log types.Log) (*SettlementEnginePaused, error) {
	event := new(SettlementEnginePaused)
	if err := _SettlementEngine.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementEngineRoundCommittedIterator is returned from FilterRoundCommitted and is used to iterate over the raw logs and unpacked data for RoundCommitted events raised by the SettlementEngine contract.
type SettlementEngineRoundCommittedIterator struct {
	Event *SettlementEngineRoundCommitted // Event containing the contract specifics and raw log

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
func (it *SettlementEngineRoundCommittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementEngineRoundCommitted)
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
		it.Event = new(SettlementEngineRoundCommitted)
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
func (it *SettlementEngineRoundCommittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementEngineRoundCommittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementEngineRoundCommitted represents a RoundCommitted event raised by the SettlementEngine contract.
type SettlementEngineRoundCommitted struct {
	Sender      common.Address
	Recipient   common.Address
	RoundId     *big.Int
	Tau         *big.Int
	FaceValue   *big.Int
	CommitBlock *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRoundCommitted is a free log retrieval operation binding the contract event 0x37f4a67788af92647c831e67034db177a77fe8b0fe306d329422e139bc0c4acb.
//
// Solidity: event RoundCommitted(address indexed sender, address indexed recipient, uint256 indexed roundId, uint256 tau, uint256 faceValue, uint256 commitBlock)
func (_SettlementEngine *SettlementEngineFilterer) FilterRoundCommitted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, roundId []*big.Int) (*SettlementEngineRoundCommittedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _SettlementEngine.contract.FilterLogs(opts, "RoundCommitted", senderRule, recipientRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &SettlementEngineRoundCommittedIterator{contract: _SettlementEngine.contract, event: "RoundCommitted", logs: logs, sub: sub}, nil
}

// WatchRoundCommitted is a free log subscription operation binding the contract event 0x37f4a67788af92647c831e67034db177a77fe8b0fe306d329422e139bc0c4acb.
//
// Solidity: event RoundCommitted(address indexed sender, address indexed recipient, uint256 indexed roundId, uint256 tau, uint256 faceValue, uint256 commitBlock)
func (_SettlementEngine *SettlementEngineFilterer) WatchRoundCommitted(opts *bind.WatchOpts, sink chan<- *SettlementEngineRoundCommitted, sender []common.Address, recipient []common.Address, roundId []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _SettlementEngine.contract.WatchLogs(opts, "RoundCommitted", senderRule, recipientRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementEngineRoundCommitted)
				if err := _SettlementEngine.contract.UnpackLog(event, "RoundCommitted", log); err != nil {
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

// ParseRoundCommitted is a log parse operation binding the contract event 0x37f4a67788af92647c831e67034db177a77fe8b0fe306d329422e139bc0c4acb.
//
// Solidity: event RoundCommitted(address indexed sender, address indexed recipient, uint256 indexed roundId, uint256 tau, uint256 faceValue, uint256 commitBlock)
func (_SettlementEngine *SettlementEngineFilterer) ParseRoundCommitted(log types.Log) (*SettlementEngineRoundCommitted, error) {
	event := new(SettlementEngineRoundCommitted)
	if err := _SettlementEngine.contract.UnpackLog(event, "RoundCommitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementEngineRoundMarkedPartialIterator is returned from FilterRoundMarkedPartial and is used to iterate over the raw logs and unpacked data for RoundMarkedPartial events raised by the SettlementEngine contract.
type SettlementEngineRoundMarkedPartialIterator struct {
	Event *SettlementEngineRoundMarkedPartial // Event containing the contract specifics and raw log

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
func (it *SettlementEngineRoundMarkedPartialIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementEngineRoundMarkedPartial)
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
		it.Event = new(SettlementEngineRoundMarkedPartial)
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
func (it *SettlementEngineRoundMarkedPartialIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementEngineRoundMarkedPartialIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementEngineRoundMarkedPartial represents a RoundMarkedPartial event raised by the SettlementEngine contract.
type SettlementEngineRoundMarkedPartial struct {
	Sender    common.Address
	Recipient common.Address
	RoundId   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRoundMarkedPartial is a free log retrieval operation binding the contract event 0x63f65de358a7452715376fc5a74b03fc4344129934500800a3814868042c3597.
//
// Solidity: event RoundMarkedPartial(address indexed sender, address indexed recipient, uint256 indexed roundId)
func (_SettlementEngine *SettlementEngineFilterer) FilterRoundMarkedPartial(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, roundId []*big.Int) (*SettlementEngineRoundMarkedPartialIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _SettlementEngine.contract.FilterLogs(opts, "RoundMarkedPartial", senderRule, recipientRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &SettlementEngineRoundMarkedPartialIterator{contract: _SettlementEngine.contract, event: "RoundMarkedPartial", logs: logs, sub: sub}, nil
}

// WatchRoundMarkedPartial is a free log subscription operation binding the contract event 0x63f65de358a7452715376fc5a74b03fc4344129934500800a3814868042c3597.
//
// Solidity: event RoundMarkedPartial(address indexed sender, address indexed recipient, uint256 indexed roundId)
func (_SettlementEngine *SettlementEngineFilterer) WatchRoundMarkedPartial(opts *bind.WatchOpts, sink chan<- *SettlementEngineRoundMarkedPartial, sender []common.Address, recipient []common.Address, roundId []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _SettlementEngine.contract.WatchLogs(opts, "RoundMarkedPartial", senderRule, recipientRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementEngineRoundMarkedPartial)
				if err := _SettlementEngine.contract.UnpackLog(event, "RoundMarkedPartial", log); err != nil {
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

// ParseRoundMarkedPartial is a log parse operation binding the contract event 0x63f65de358a7452715376fc5a74b03fc4344129934500800a3814868042c3597.
//
// Solidity: event RoundMarkedPartial(address indexed sender, address indexed recipient, uint256 indexed roundId)
func (_SettlementEngine *SettlementEngineFilterer) ParseRoundMarkedPartial(log types.Log) (*SettlementEngineRoundMarkedPartial, error) {
	event := new(SettlementEngineRoundMarkedPartial)
	if err := _SettlementEngine.contract.UnpackLog(event, "RoundMarkedPartial", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementEngineRoundSettledIterator is returned from FilterRoundSettled and is used to iterate over the raw logs and unpacked data for RoundSettled events raised by the SettlementEngine contract.
type SettlementEngineRoundSettledIterator struct {
	Event *SettlementEngineRoundSettled // Event containing the contract specifics and raw log

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
func (it *SettlementEngineRoundSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementEngineRoundSettled)
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
		it.Event = new(SettlementEngineRoundSettled)
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
func (it *SettlementEngineRoundSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementEngineRoundSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementEngineRoundSettled represents a RoundSettled event raised by the SettlementEngine contract.
type SettlementEngineRoundSettled struct {
	Sender            common.Address
	Recipient         common.Address
	RoundId           *big.Int
	WinningLocalIndex *big.Int
	FaceValue         *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoundSettled is a free log retrieval operation binding the contract event 0x0b76e789e4341a8a74776cbf3af9cfc0df7246d30f52a25462d22762a0f346c5.
//
// Solidity: event RoundSettled(address indexed sender, address indexed recipient, uint256 indexed roundId, uint256 winningLocalIndex, uint256 faceValue)
func (_SettlementEngine *SettlementEngineFilterer) FilterRoundSettled(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, roundId []*big.Int) (*SettlementEngineRoundSettledIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _SettlementEngine.contract.FilterLogs(opts, "RoundSettled", senderRule, recipientRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &SettlementEngineRoundSettledIterator{contract: _SettlementEngine.contract, event: "RoundSettled", logs: logs, sub: sub}, nil
}

// WatchRoundSettled is a free log subscription operation binding the contract event 0x0b76e789e4341a8a74776cbf3af9cfc0df7246d30f52a25462d22762a0f346c5.
//
// Solidity: event RoundSettled(address indexed sender, address indexed recipient, uint256 indexed roundId, uint256 winningLocalIndex, uint256 faceValue)
func (_SettlementEngine *SettlementEngineFilterer) WatchRoundSettled(opts *bind.WatchOpts, sink chan<- *SettlementEngineRoundSettled, sender []common.Address, recipient []common.Address, roundId []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _SettlementEngine.contract.WatchLogs(opts, "RoundSettled", senderRule, recipientRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementEngineRoundSettled)
				if err := _SettlementEngine.contract.UnpackLog(event, "RoundSettled", log); err != nil {
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

// ParseRoundSettled is a log parse operation binding the contract event 0x0b76e789e4341a8a74776cbf3af9cfc0df7246d30f52a25462d22762a0f346c5.
//
// Solidity: event RoundSettled(address indexed sender, address indexed recipient, uint256 indexed roundId, uint256 winningLocalIndex, uint256 faceValue)
func (_SettlementEngine *SettlementEngineFilterer) ParseRoundSettled(log types.Log) (*SettlementEngineRoundSettled, error) {
	event := new(SettlementEngineRoundSettled)
	if err := _SettlementEngine.contract.UnpackLog(event, "RoundSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementEngineRoundVoidedIterator is returned from FilterRoundVoided and is used to iterate over the raw logs and unpacked data for RoundVoided events raised by the SettlementEngine contract.
type SettlementEngineRoundVoidedIterator struct {
	Event *SettlementEngineRoundVoided // Event containing the contract specifics and raw log

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
func (it *SettlementEngineRoundVoidedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementEngineRoundVoided)
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
		it.Event = new(SettlementEngineRoundVoided)
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
func (it *SettlementEngineRoundVoidedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementEngineRoundVoidedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementEngineRoundVoided represents a RoundVoided event raised by the SettlementEngine contract.
type SettlementEngineRoundVoided struct {
	Sender    common.Address
	Recipient common.Address
	RoundId   *big.Int
	Reason    [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRoundVoided is a free log retrieval operation binding the contract event 0xe15b7d92d30e2e0d903eccb550b6fe43024271cc0eca67f900b793307814ddc1.
//
// Solidity: event RoundVoided(address indexed sender, address indexed recipient, uint256 indexed roundId, bytes32 reason)
func (_SettlementEngine *SettlementEngineFilterer) FilterRoundVoided(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, roundId []*big.Int) (*SettlementEngineRoundVoidedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _SettlementEngine.contract.FilterLogs(opts, "RoundVoided", senderRule, recipientRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &SettlementEngineRoundVoidedIterator{contract: _SettlementEngine.contract, event: "RoundVoided", logs: logs, sub: sub}, nil
}

// WatchRoundVoided is a free log subscription operation binding the contract event 0xe15b7d92d30e2e0d903eccb550b6fe43024271cc0eca67f900b793307814ddc1.
//
// Solidity: event RoundVoided(address indexed sender, address indexed recipient, uint256 indexed roundId, bytes32 reason)
func (_SettlementEngine *SettlementEngineFilterer) WatchRoundVoided(opts *bind.WatchOpts, sink chan<- *SettlementEngineRoundVoided, sender []common.Address, recipient []common.Address, roundId []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _SettlementEngine.contract.WatchLogs(opts, "RoundVoided", senderRule, recipientRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementEngineRoundVoided)
				if err := _SettlementEngine.contract.UnpackLog(event, "RoundVoided", log); err != nil {
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

// ParseRoundVoided is a log parse operation binding the contract event 0xe15b7d92d30e2e0d903eccb550b6fe43024271cc0eca67f900b793307814ddc1.
//
// Solidity: event RoundVoided(address indexed sender, address indexed recipient, uint256 indexed roundId, bytes32 reason)
func (_SettlementEngine *SettlementEngineFilterer) ParseRoundVoided(log types.Log) (*SettlementEngineRoundVoided, error) {
	event := new(SettlementEngineRoundVoided)
	if err := _SettlementEngine.contract.UnpackLog(event, "RoundVoided", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementEngineUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the SettlementEngine contract.
type SettlementEngineUnpausedIterator struct {
	Event *SettlementEngineUnpaused // Event containing the contract specifics and raw log

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
func (it *SettlementEngineUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementEngineUnpaused)
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
		it.Event = new(SettlementEngineUnpaused)
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
func (it *SettlementEngineUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementEngineUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementEngineUnpaused represents a Unpaused event raised by the SettlementEngine contract.
type SettlementEngineUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SettlementEngine *SettlementEngineFilterer) FilterUnpaused(opts *bind.FilterOpts) (*SettlementEngineUnpausedIterator, error) {

	logs, sub, err := _SettlementEngine.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &SettlementEngineUnpausedIterator{contract: _SettlementEngine.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SettlementEngine *SettlementEngineFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *SettlementEngineUnpaused) (event.Subscription, error) {

	logs, sub, err := _SettlementEngine.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementEngineUnpaused)
				if err := _SettlementEngine.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_SettlementEngine *SettlementEngineFilterer) ParseUnpaused(log types.Log) (*SettlementEngineUnpaused, error) {
	event := new(SettlementEngineUnpaused)
	if err := _SettlementEngine.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
