// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package bind

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ubiq/go-ubiq/v6/common"
)

var bindTests = []struct {
	name     string
	contract string
	bytecode []string
	abi      []string
	imports  string
	tester   string
	fsigs    []map[string]string
	libs     map[string]string
	aliases  map[string]string
	types    []string
}{
	// Test that the binding is available in combined and separate forms too
	{
		`Empty`,
		`contract NilContract {}`,
		[]string{`606060405260068060106000396000f3606060405200`},
		[]string{`[]`},
		`"github.com/ubiq/go-ubiq/v6/common"`,
		`
			if b, err := NewEmpty(common.Address{}, nil); b == nil || err != nil {
				t.Fatalf("combined binding (%v) nil or error (%v) not nil", b, nil)
			}
			if b, err := NewEmptyCaller(common.Address{}, nil); b == nil || err != nil {
				t.Fatalf("caller binding (%v) nil or error (%v) not nil", b, nil)
			}
			if b, err := NewEmptyTransactor(common.Address{}, nil); b == nil || err != nil {
				t.Fatalf("transactor binding (%v) nil or error (%v) not nil", b, nil)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Test that all the official sample contracts bind correctly
	{
		`Token`,
		`https://ethereum.org/token`,
		[]string{`60606040526040516107fd3803806107fd83398101604052805160805160a05160c051929391820192909101600160a060020a0333166000908152600360209081526040822086905581548551838052601f6002600019610100600186161502019093169290920482018390047f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e56390810193919290918801908390106100e857805160ff19168380011785555b506101189291505b8082111561017157600081556001016100b4565b50506002805460ff19168317905550505050610658806101a56000396000f35b828001600101855582156100ac579182015b828111156100ac5782518260005055916020019190600101906100fa565b50508060016000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061017557805160ff19168380011785555b506100c89291506100b4565b5090565b82800160010185558215610165579182015b8281111561016557825182600050559160200191906001019061018756606060405236156100775760e060020a600035046306fdde03811461007f57806323b872dd146100dc578063313ce5671461010e57806370a082311461011a57806395d89b4114610132578063a9059cbb1461018e578063cae9ca51146101bd578063dc3080f21461031c578063dd62ed3e14610341575b610365610002565b61036760008054602060026001831615610100026000190190921691909104601f810182900490910260809081016040526060828152929190828280156104eb5780601f106104c0576101008083540402835291602001916104eb565b6103d5600435602435604435600160a060020a038316600090815260036020526040812054829010156104f357610002565b6103e760025460ff1681565b6103d560043560036020526000908152604090205481565b610367600180546020600282841615610100026000190190921691909104601f810182900490910260809081016040526060828152929190828280156104eb5780601f106104c0576101008083540402835291602001916104eb565b610365600435602435600160a060020a033316600090815260036020526040902054819010156103f157610002565b60806020604435600481810135601f8101849004909302840160405260608381526103d5948235946024803595606494939101919081908382808284375094965050505050505060006000836004600050600033600160a060020a03168152602001908152602001600020600050600087600160a060020a031681526020019081526020016000206000508190555084905080600160a060020a0316638f4ffcb1338630876040518560e060020a0281526004018085600160a060020a0316815260200184815260200183600160a060020a03168152602001806020018281038252838181518152602001915080519060200190808383829060006004602084601f0104600f02600301f150905090810190601f1680156102f25780820380516001836020036101000a031916815260200191505b50955050505050506000604051808303816000876161da5a03f11561000257505050509392505050565b6005602090815260043560009081526040808220909252602435815220546103d59081565b60046020818152903560009081526040808220909252602435815220546103d59081565b005b60405180806020018281038252838181518152602001915080519060200190808383829060006004602084601f0104600f02600301f150905090810190601f1680156103c75780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60408051918252519081900360200190f35b6060908152602090f35b600160a060020a03821660009081526040902054808201101561041357610002565b806003600050600033600160a060020a03168152602001908152602001600020600082828250540392505081905550806003600050600084600160a060020a0316815260200190815260200160002060008282825054019250508190555081600160a060020a031633600160a060020a03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a35050565b820191906000526020600020905b8154815290600101906020018083116104ce57829003601f168201915b505050505081565b600160a060020a03831681526040812054808301101561051257610002565b600160a060020a0380851680835260046020908152604080852033949094168086529382528085205492855260058252808520938552929052908220548301111561055c57610002565b816003600050600086600160a060020a03168152602001908152602001600020600082828250540392505081905550816003600050600085600160a060020a03168152602001908152602001600020600082828250540192505081905550816005600050600086600160a060020a03168152602001908152602001600020600050600033600160a060020a0316815260200190815260200160002060008282825054019250508190555082600160a060020a031633600160a060020a03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a3939250505056`},
		[]string{`[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"},{"name":"_extraData","type":"bytes"}],"name":"approveAndCall","outputs":[{"name":"success","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"spentAllowance","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"inputs":[{"name":"initialSupply","type":"uint256"},{"name":"tokenName","type":"string"},{"name":"decimalUnits","type":"uint8"},{"name":"tokenSymbol","type":"string"}],"type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`},
		`"github.com/ubiq/go-ubiq/v6/common"`,
		`
			if b, err := NewToken(common.Address{}, nil); b == nil || err != nil {
				t.Fatalf("binding (%v) nil or error (%v) not nil", b, nil)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	{
		`Crowdsale`,
		`https://ethereum.org/crowdsale`,
		[]string{`606060408190526007805460ff1916905560a0806105a883396101006040529051608051915160c05160e05160008054600160a060020a03199081169095178155670de0b6b3a7640000958602600155603c9093024201600355930260045560058054909216909217905561052f90819061007990396000f36060604052361561006c5760e060020a600035046301cb3b20811461008257806329dcb0cf1461014457806338af3eed1461014d5780636e66f6e91461015f5780637a3a0e84146101715780637b3e5e7b1461017a578063a035b1fe14610183578063dc0d3dff1461018c575b61020060075460009060ff161561032357610002565b61020060035460009042106103205760025460015490106103cb576002548154600160a060020a0316908290606082818181858883f150915460025460408051600160a060020a039390931683526020830191909152818101869052517fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf6945090819003909201919050a15b60405160008054600160a060020a039081169230909116319082818181858883f150506007805460ff1916600117905550505050565b6103a160035481565b6103ab600054600160a060020a031681565b6103ab600554600160a060020a031681565b6103a160015481565b6103a160025481565b6103a160045481565b6103be60043560068054829081101561000257506000526002027ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f8101547ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d409190910154600160a060020a03919091169082565b005b505050815481101561000257906000526020600020906002020160005060008201518160000160006101000a815481600160a060020a030219169083021790555060208201518160010160005055905050806002600082828250540192505081905550600560009054906101000a9004600160a060020a0316600160a060020a031663a9059cbb3360046000505484046040518360e060020a0281526004018083600160a060020a03168152602001828152602001925050506000604051808303816000876161da5a03f11561000257505060408051600160a060020a03331681526020810184905260018183015290517fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf692509081900360600190a15b50565b5060a0604052336060908152346080819052600680546001810180835592939282908280158290116102025760020281600202836000526020600020918201910161020291905b8082111561039d57805473ffffffffffffffffffffffffffffffffffffffff19168155600060019190910190815561036a565b5090565b6060908152602090f35b600160a060020a03166060908152602090f35b6060918252608052604090f35b5b60065481101561010e576006805482908110156100025760009182526002027ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0190600680549254600160a060020a0316928490811015610002576002027ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d40015460405190915082818181858883f19350505050507fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf660066000508281548110156100025760008290526002027ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f01548154600160a060020a039190911691908490811015610002576002027ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d40015460408051600160a060020a0394909416845260208401919091526000838201525191829003606001919050a16001016103cc56`},
		[]string{`[{"constant":false,"inputs":[],"name":"checkGoalReached","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"deadline","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[],"name":"beneficiary","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":true,"inputs":[],"name":"tokenReward","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":true,"inputs":[],"name":"fundingGoal","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[],"name":"amountRaised","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[],"name":"price","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"funders","outputs":[{"name":"addr","type":"address"},{"name":"amount","type":"uint256"}],"type":"function"},{"inputs":[{"name":"ifSuccessfulSendTo","type":"address"},{"name":"fundingGoalInEthers","type":"uint256"},{"name":"durationInMinutes","type":"uint256"},{"name":"etherCostOfEachToken","type":"uint256"},{"name":"addressOfTokenUsedAsReward","type":"address"}],"type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"backer","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"isContribution","type":"bool"}],"name":"FundTransfer","type":"event"}]`},
		`"github.com/ubiq/go-ubiq/v6/common"`,
		`
			if b, err := NewCrowdsale(common.Address{}, nil); b == nil || err != nil {
				t.Fatalf("binding (%v) nil or error (%v) not nil", b, nil)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Test that named and anonymous inputs are handled correctly
	{
		`InputChecker`, ``, []string{``},
		[]string{`
			[
				{"type":"function","name":"noInput","constant":true,"inputs":[],"outputs":[]},
				{"type":"function","name":"namedInput","constant":true,"inputs":[{"name":"str","type":"string"}],"outputs":[]},
				{"type":"function","name":"anonInput","constant":true,"inputs":[{"name":"","type":"string"}],"outputs":[]},
				{"type":"function","name":"namedInputs","constant":true,"inputs":[{"name":"str1","type":"string"},{"name":"str2","type":"string"}],"outputs":[]},
				{"type":"function","name":"anonInputs","constant":true,"inputs":[{"name":"","type":"string"},{"name":"","type":"string"}],"outputs":[]},
				{"type":"function","name":"mixedInputs","constant":true,"inputs":[{"name":"","type":"string"},{"name":"str","type":"string"}],"outputs":[]}
			]
		`},
		`
			"fmt"

			"github.com/ubiq/go-ubiq/v6/common"
		`,
		`if b, err := NewInputChecker(common.Address{}, nil); b == nil || err != nil {
			 t.Fatalf("binding (%v) nil or error (%v) not nil", b, nil)
		 } else if false { // Don't run, just compile and test types
			 var err error

			 err = b.NoInput(nil)
			 err = b.NamedInput(nil, "")
			 err = b.AnonInput(nil, "")
			 err = b.NamedInputs(nil, "", "")
			 err = b.AnonInputs(nil, "", "")
			 err = b.MixedInputs(nil, "", "")

			 fmt.Println(err)
		 }`,
		nil,
		nil,
		nil,
		nil,
	},
	// Test that named and anonymous outputs are handled correctly
	{
		`OutputChecker`, ``, []string{``},
		[]string{`
			[
				{"type":"function","name":"noOutput","constant":true,"inputs":[],"outputs":[]},
				{"type":"function","name":"namedOutput","constant":true,"inputs":[],"outputs":[{"name":"str","type":"string"}]},
				{"type":"function","name":"anonOutput","constant":true,"inputs":[],"outputs":[{"name":"","type":"string"}]},
				{"type":"function","name":"namedOutputs","constant":true,"inputs":[],"outputs":[{"name":"str1","type":"string"},{"name":"str2","type":"string"}]},
				{"type":"function","name":"collidingOutputs","constant":true,"inputs":[],"outputs":[{"name":"str","type":"string"},{"name":"Str","type":"string"}]},
				{"type":"function","name":"anonOutputs","constant":true,"inputs":[],"outputs":[{"name":"","type":"string"},{"name":"","type":"string"}]},
				{"type":"function","name":"mixedOutputs","constant":true,"inputs":[],"outputs":[{"name":"","type":"string"},{"name":"str","type":"string"}]}
			]
		`},
		`
			"fmt"

			"github.com/ubiq/go-ubiq/v6/common"
		`,
		`if b, err := NewOutputChecker(common.Address{}, nil); b == nil || err != nil {
			 t.Fatalf("binding (%v) nil or error (%v) not nil", b, nil)
		 } else if false { // Don't run, just compile and test types
			 var str1, str2 string
			 var err error

			 err              = b.NoOutput(nil)
			 str1, err        = b.NamedOutput(nil)
			 str1, err        = b.AnonOutput(nil)
			 res, _          := b.NamedOutputs(nil)
			 str1, str2, err  = b.CollidingOutputs(nil)
			 str1, str2, err  = b.AnonOutputs(nil)
			 str1, str2, err  = b.MixedOutputs(nil)

			 fmt.Println(str1, str2, res.Str1, res.Str2, err)
		 }`,
		nil,
		nil,
		nil,
		nil,
	},
	// Tests that named, anonymous and indexed events are handled correctly
	{
		`EventChecker`, ``, []string{``},
		[]string{`
			[
				{"type":"event","name":"empty","inputs":[]},
				{"type":"event","name":"indexed","inputs":[{"name":"addr","type":"address","indexed":true},{"name":"num","type":"int256","indexed":true}]},
				{"type":"event","name":"mixed","inputs":[{"name":"addr","type":"address","indexed":true},{"name":"num","type":"int256"}]},
				{"type":"event","name":"anonymous","anonymous":true,"inputs":[]},
				{"type":"event","name":"dynamic","inputs":[{"name":"idxStr","type":"string","indexed":true},{"name":"idxDat","type":"bytes","indexed":true},{"name":"str","type":"string"},{"name":"dat","type":"bytes"}]},
				{"type":"event","name":"unnamed","inputs":[{"name":"","type":"uint256","indexed": true},{"name":"","type":"uint256","indexed":true}]}
			]
		`},
		`
			"fmt"
			"math/big"
			"reflect"

			"github.com/ubiq/go-ubiq/v6/common"
		`,
		`if e, err := NewEventChecker(common.Address{}, nil); e == nil || err != nil {
			 t.Fatalf("binding (%v) nil or error (%v) not nil", e, nil)
		 } else if false { // Don't run, just compile and test types
			 var (
				 err  error
			   res  bool
				 str  string
				 dat  []byte
				 hash common.Hash
			 )
			 _, err = e.FilterEmpty(nil)
			 _, err = e.FilterIndexed(nil, []common.Address{}, []*big.Int{})

			 mit, err := e.FilterMixed(nil, []common.Address{})

			 res = mit.Next()  // Make sure the iterator has a Next method
			 err = mit.Error() // Make sure the iterator has an Error method
			 err = mit.Close() // Make sure the iterator has a Close method

			 fmt.Println(mit.Event.Raw.BlockHash) // Make sure the raw log is contained within the results
			 fmt.Println(mit.Event.Num)           // Make sure the unpacked non-indexed fields are present
			 fmt.Println(mit.Event.Addr)          // Make sure the reconstructed indexed fields are present

			 dit, err := e.FilterDynamic(nil, []string{}, [][]byte{})

			 str  = dit.Event.Str    // Make sure non-indexed strings retain their type
			 dat  = dit.Event.Dat    // Make sure non-indexed bytes retain their type
			 hash = dit.Event.IdxStr // Make sure indexed strings turn into hashes
			 hash = dit.Event.IdxDat // Make sure indexed bytes turn into hashes

			 sink := make(chan *EventCheckerMixed)
			 sub, err := e.WatchMixed(nil, sink, []common.Address{})
			 defer sub.Unsubscribe()

			 event := <-sink
			 fmt.Println(event.Raw.BlockHash) // Make sure the raw log is contained within the results
			 fmt.Println(event.Num)           // Make sure the unpacked non-indexed fields are present
			 fmt.Println(event.Addr)          // Make sure the reconstructed indexed fields are present

			 fmt.Println(res, str, dat, hash, err)

			 oit, err := e.FilterUnnamed(nil, []*big.Int{}, []*big.Int{})

			 arg0  := oit.Event.Arg0    // Make sure unnamed arguments are handled correctly
			 arg1  := oit.Event.Arg1    // Make sure unnamed arguments are handled correctly
			 fmt.Println(arg0, arg1)
		 }
		 // Run a tiny reflection test to ensure disallowed methods don't appear
		 if _, ok := reflect.TypeOf(&EventChecker{}).MethodByName("FilterAnonymous"); ok {
		 	t.Errorf("binding has disallowed method (FilterAnonymous)")
		 }`,
		nil,
		nil,
		nil,
		nil,
	},
	// Test that contract interactions (deploy, transact and call) generate working code
	{
		`Interactor`,
		`
			contract Interactor {
				string public deployString;
				string public transactString;

				function Interactor(string str) {
				  deployString = str;
				}

				function transact(string str) {
				  transactString = str;
				}
			}
		`,
		[]string{`6060604052604051610328380380610328833981016040528051018060006000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10608d57805160ff19168380011785555b50607c9291505b8082111560ba57838155600101606b565b50505061026a806100be6000396000f35b828001600101855582156064579182015b828111156064578251826000505591602001919060010190609e565b509056606060405260e060020a60003504630d86a0e181146100315780636874e8091461008d578063d736c513146100ea575b005b610190600180546020600282841615610100026000190190921691909104601f810182900490910260809081016040526060828152929190828280156102295780601f106101fe57610100808354040283529160200191610229565b61019060008054602060026001831615610100026000190190921691909104601f810182900490910260809081016040526060828152929190828280156102295780601f106101fe57610100808354040283529160200191610229565b60206004803580820135601f81018490049093026080908101604052606084815261002f946024939192918401918190838280828437509496505050505050508060016000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061023157805160ff19168380011785555b506102619291505b808211156102665760008155830161017d565b60405180806020018281038252838181518152602001915080519060200190808383829060006004602084601f0104600f02600301f150905090810190601f1680156101f05780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b820191906000526020600020905b81548152906001019060200180831161020c57829003601f168201915b505050505081565b82800160010185558215610175579182015b82811115610175578251826000505591602001919060010190610243565b505050565b509056`},
		[]string{`[{"constant":true,"inputs":[],"name":"transactString","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"deployString","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":false,"inputs":[{"name":"str","type":"string"}],"name":"transact","outputs":[],"type":"function"},{"inputs":[{"name":"str","type":"string"}],"type":"constructor"}]`},
		`
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			// Deploy an interaction tester contract and call a transaction on it
			_, _, interactor, err := DeployInteractor(auth, sim, "Deploy string")
			if err != nil {
				t.Fatalf("Failed to deploy interactor contract: %v", err)
			}
			if _, err := interactor.Transact(auth, "Transact string"); err != nil {
				t.Fatalf("Failed to transact with interactor contract: %v", err)
			}
			// Commit all pending transactions in the simulator and check the contract state
			sim.Commit()

			if str, err := interactor.DeployString(nil); err != nil {
				t.Fatalf("Failed to retrieve deploy string: %v", err)
			} else if str != "Deploy string" {
				t.Fatalf("Deploy string mismatch: have '%s', want 'Deploy string'", str)
			}
			if str, err := interactor.TransactString(nil); err != nil {
				t.Fatalf("Failed to retrieve transact string: %v", err)
			} else if str != "Transact string" {
				t.Fatalf("Transact string mismatch: have '%s', want 'Transact string'", str)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Tests that plain values can be properly returned and deserialized
	{
		`Getter`,
		`
			contract Getter {
				function getter() constant returns (string, int, bytes32) {
					return ("Hi", 1, sha3(""));
				}
			}
		`,
		[]string{`606060405260dc8060106000396000f3606060405260e060020a6000350463993a04b78114601a575b005b600060605260c0604052600260809081527f486900000000000000000000000000000000000000000000000000000000000060a05260017fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47060e0829052610100819052606060c0908152600261012081905281906101409060a09080838184600060046012f1505081517fffff000000000000000000000000000000000000000000000000000000000000169091525050604051610160819003945092505050f3`},
		[]string{`[{"constant":true,"inputs":[],"name":"getter","outputs":[{"name":"","type":"string"},{"name":"","type":"int256"},{"name":"","type":"bytes32"}],"type":"function"}]`},
		`
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			// Deploy a tuple tester contract and execute a structured call on it
			_, _, getter, err := DeployGetter(auth, sim)
			if err != nil {
				t.Fatalf("Failed to deploy getter contract: %v", err)
			}
			sim.Commit()

			if str, num, _, err := getter.Getter(nil); err != nil {
				t.Fatalf("Failed to call anonymous field retriever: %v", err)
			} else if str != "Hi" || num.Cmp(big.NewInt(1)) != 0 {
				t.Fatalf("Retrieved value mismatch: have %v/%v, want %v/%v", str, num, "Hi", 1)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Tests that tuples can be properly returned and deserialized
	{
		`Tupler`,
		`
			contract Tupler {
				function tuple() constant returns (string a, int b, bytes32 c) {
					return ("Hi", 1, sha3(""));
				}
			}
		`,
		[]string{`606060405260dc8060106000396000f3606060405260e060020a60003504633175aae28114601a575b005b600060605260c0604052600260809081527f486900000000000000000000000000000000000000000000000000000000000060a05260017fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47060e0829052610100819052606060c0908152600261012081905281906101409060a09080838184600060046012f1505081517fffff000000000000000000000000000000000000000000000000000000000000169091525050604051610160819003945092505050f3`},
		[]string{`[{"constant":true,"inputs":[],"name":"tuple","outputs":[{"name":"a","type":"string"},{"name":"b","type":"int256"},{"name":"c","type":"bytes32"}],"type":"function"}]`},
		`
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			// Deploy a tuple tester contract and execute a structured call on it
			_, _, tupler, err := DeployTupler(auth, sim)
			if err != nil {
				t.Fatalf("Failed to deploy tupler contract: %v", err)
			}
			sim.Commit()

			if res, err := tupler.Tuple(nil); err != nil {
				t.Fatalf("Failed to call structure retriever: %v", err)
			} else if res.A != "Hi" || res.B.Cmp(big.NewInt(1)) != 0 {
				t.Fatalf("Retrieved value mismatch: have %v/%v, want %v/%v", res.A, res.B, "Hi", 1)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Tests that arrays/slices can be properly returned and deserialized.
	// Only addresses are tested, remainder just compiled to keep the test small.
	{
		`Slicer`,
		`
			contract Slicer {
				function echoAddresses(address[] input) constant returns (address[] output) {
					return input;
				}
				function echoInts(int[] input) constant returns (int[] output) {
					return input;
				}
				function echoFancyInts(uint24[23] input) constant returns (uint24[23] output) {
					return input;
				}
				function echoBools(bool[] input) constant returns (bool[] output) {
					return input;
				}
			}
		`,
		[]string{`606060405261015c806100126000396000f3606060405260e060020a6000350463be1127a3811461003c578063d88becc014610092578063e15a3db71461003c578063f637e5891461003c575b005b604080516020600480358082013583810285810185019096528085526100ee959294602494909392850192829185019084908082843750949650505050505050604080516020810190915260009052805b919050565b604080516102e0818101909252610138916004916102e491839060179083908390808284375090955050505050506102e0604051908101604052806017905b60008152602001906001900390816100d15790505081905061008d565b60405180806020018281038252838181518152602001915080519060200190602002808383829060006004602084601f0104600f02600301f1509050019250505060405180910390f35b60405180826102e0808381846000600461015cf15090500191505060405180910390f3`},
		[]string{`[{"constant":true,"inputs":[{"name":"input","type":"address[]"}],"name":"echoAddresses","outputs":[{"name":"output","type":"address[]"}],"type":"function"},{"constant":true,"inputs":[{"name":"input","type":"uint24[23]"}],"name":"echoFancyInts","outputs":[{"name":"output","type":"uint24[23]"}],"type":"function"},{"constant":true,"inputs":[{"name":"input","type":"int256[]"}],"name":"echoInts","outputs":[{"name":"output","type":"int256[]"}],"type":"function"},{"constant":true,"inputs":[{"name":"input","type":"bool[]"}],"name":"echoBools","outputs":[{"name":"output","type":"bool[]"}],"type":"function"}]`},
		`
			"math/big"
			"reflect"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/common"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			// Deploy a slice tester contract and execute a n array call on it
			_, _, slicer, err := DeploySlicer(auth, sim)
			if err != nil {
					t.Fatalf("Failed to deploy slicer contract: %v", err)
			}
			sim.Commit()

			if out, err := slicer.EchoAddresses(nil, []common.Address{auth.From, common.Address{}}); err != nil {
					t.Fatalf("Failed to call slice echoer: %v", err)
			} else if !reflect.DeepEqual(out, []common.Address{auth.From, common.Address{}}) {
					t.Fatalf("Slice return mismatch: have %v, want %v", out, []common.Address{auth.From, common.Address{}})
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Tests that anonymous default methods can be correctly invoked
	{
		`Defaulter`,
		`
			contract Defaulter {
				address public caller;

				function() {
					caller = msg.sender;
				}
			}
		`,
		[]string{`6060604052606a8060106000396000f360606040523615601d5760e060020a6000350463fc9c8d3981146040575b605e6000805473ffffffffffffffffffffffffffffffffffffffff191633179055565b606060005473ffffffffffffffffffffffffffffffffffffffff1681565b005b6060908152602090f3`},
		[]string{`[{"constant":true,"inputs":[],"name":"caller","outputs":[{"name":"","type":"address"}],"type":"function"}]`},
		`
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			// Deploy a default method invoker contract and execute its default method
			_, _, defaulter, err := DeployDefaulter(auth, sim)
			if err != nil {
				t.Fatalf("Failed to deploy defaulter contract: %v", err)
			}
			if _, err := (&DefaulterRaw{defaulter}).Transfer(auth); err != nil {
				t.Fatalf("Failed to invoke default method: %v", err)
			}
			sim.Commit()

			if caller, err := defaulter.Caller(nil); err != nil {
				t.Fatalf("Failed to call address retriever: %v", err)
			} else if (caller != auth.From) {
				t.Fatalf("Address mismatch: have %v, want %v", caller, auth.From)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Tests that structs are correctly unpacked
	{

		`Structs`,
		`
		pragma solidity ^0.6.5;
			pragma experimental ABIEncoderV2;
			contract Structs {
				struct A {
					bytes32 B;
				}
				
				function F() public view returns (A[] memory a, uint256[] memory c, bool[] memory d) {
					A[] memory a = new A[](2);
					a[0].B = bytes32(uint256(1234) << 96);
					uint256[] memory c;
					bool[] memory d;
					return (a, c, d);
				}
			
				function G() public view returns (A[] memory a) {
					A[] memory a = new A[](2);
					a[0].B = bytes32(uint256(1234) << 96);
					return a;
				}
			}
		`,
		[]string{`608060405234801561001057600080fd5b50610278806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806328811f591461003b5780636fecb6231461005b575b600080fd5b610043610070565b604051610052939291906101a0565b60405180910390f35b6100636100d6565b6040516100529190610186565b604080516002808252606082810190935282918291829190816020015b610095610131565b81526020019060019003908161008d575050805190915061026960611b9082906000906100be57fe5b60209081029190910101515293606093508392509050565b6040805160028082526060828101909352829190816020015b6100f7610131565b8152602001906001900390816100ef575050805190915061026960611b90829060009061012057fe5b602090810291909101015152905090565b60408051602081019091526000815290565b815260200190565b6000815180845260208085019450808401835b8381101561017b578151518752958201959082019060010161015e565b509495945050505050565b600060208252610199602083018461014b565b9392505050565b6000606082526101b3606083018661014b565b6020838203818501528186516101c98185610239565b91508288019350845b818110156101f3576101e5838651610143565b9484019492506001016101d2565b505084810360408601528551808252908201925081860190845b8181101561022b57825115158552938301939183019160010161020d565b509298975050505050505050565b9081526020019056fea2646970667358221220eb85327e285def14230424c52893aebecec1e387a50bb6b75fc4fdbed647f45f64736f6c63430006050033`},
		[]string{`[{"inputs":[],"name":"F","outputs":[{"components":[{"internalType":"bytes32","name":"B","type":"bytes32"}],"internalType":"structStructs.A[]","name":"a","type":"tuple[]"},{"internalType":"uint256[]","name":"c","type":"uint256[]"},{"internalType":"bool[]","name":"d","type":"bool[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"G","outputs":[{"components":[{"internalType":"bytes32","name":"B","type":"bytes32"}],"internalType":"structStructs.A[]","name":"a","type":"tuple[]"}],"stateMutability":"view","type":"function"}]`},
		`
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		
			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()
		
			// Deploy a structs method invoker contract and execute its default method
			_, _, structs, err := DeployStructs(auth, sim)
			if err != nil {
				t.Fatalf("Failed to deploy defaulter contract: %v", err)
			}
			sim.Commit()
			opts := bind.CallOpts{}
			if _, err := structs.F(&opts); err != nil {
				t.Fatalf("Failed to invoke F method: %v", err)
			}
			if _, err := structs.G(&opts); err != nil {
				t.Fatalf("Failed to invoke G method: %v", err)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Tests that non-existent contracts are reported as such (though only simulator test)
	{
		`NonExistent`,
		`
			contract NonExistent {
				function String() constant returns(string) {
					return "I don't exist";
				}
			}
		`,
		[]string{`6060604052609f8060106000396000f3606060405260e060020a6000350463f97a60058114601a575b005b600060605260c0604052600d60809081527f4920646f6e27742065786973740000000000000000000000000000000000000060a052602060c0908152600d60e081905281906101009060a09080838184600060046012f15050815172ffffffffffffffffffffffffffffffffffffff1916909152505060405161012081900392509050f3`},
		[]string{`[{"constant":true,"inputs":[],"name":"String","outputs":[{"name":"","type":"string"}],"type":"function"}]`},
		`
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/common"
			"github.com/ubiq/go-ubiq/v6/core"
		`,
		`
			// Create a simulator and wrap a non-deployed contract

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{}, uint64(10000000000))
			defer sim.Close()

			nonexistent, err := NewNonExistent(common.Address{}, sim)
			if err != nil {
				t.Fatalf("Failed to access non-existent contract: %v", err)
			}
			// Ensure that contract calls fail with the appropriate error
			if res, err := nonexistent.String(nil); err == nil {
				t.Fatalf("Call succeeded on non-existent contract: %v", res)
			} else if (err != bind.ErrNoCode) {
				t.Fatalf("Error mismatch: have %v, want %v", err, bind.ErrNoCode)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	{
		`NonExistentStruct`,
		`
			contract NonExistentStruct {
				function Struct() public view returns(uint256 a, uint256 b) {
					return (10, 10);
				}
			}
		`,
		[]string{`6080604052348015600f57600080fd5b5060888061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063d5f6622514602d575b600080fd5b6033604c565b6040805192835260208301919091528051918290030190f35b600a809156fea264697066735822beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef64736f6c6343decafe0033`},
		[]string{`[{"inputs":[],"name":"Struct","outputs":[{"internalType":"uint256","name":"a","type":"uint256"},{"internalType":"uint256","name":"b","type":"uint256"}],"stateMutability":"pure","type":"function"}]`},
		`
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/common"
			"github.com/ubiq/go-ubiq/v6/core"
		`,
		`
			// Create a simulator and wrap a non-deployed contract

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{}, uint64(10000000000))
			defer sim.Close()

			nonexistent, err := NewNonExistentStruct(common.Address{}, sim)
			if err != nil {
				t.Fatalf("Failed to access non-existent contract: %v", err)
			}
			// Ensure that contract calls fail with the appropriate error
			if res, err := nonexistent.Struct(nil); err == nil {
				t.Fatalf("Call succeeded on non-existent contract: %v", res)
			} else if (err != bind.ErrNoCode) {
				t.Fatalf("Error mismatch: have %v, want %v", err, bind.ErrNoCode)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Tests that gas estimation works for contracts with weird gas mechanics too.
	{
		`FunkyGasPattern`,
		`
			contract FunkyGasPattern {
				string public field;

				function SetField(string value) {
					// This check will screw gas estimation! Good, good!
					if (msg.gas < 100000) {
						throw;
					}
					field = value;
				}
			}
		`,
		[]string{`606060405261021c806100126000396000f3606060405260e060020a600035046323fcf32a81146100265780634f28bf0e1461007b575b005b6040805160206004803580820135601f8101849004840285018401909552848452610024949193602493909291840191908190840183828082843750949650505050505050620186a05a101561014e57610002565b6100db60008054604080516020601f600260001961010060018816150201909516949094049384018190048102820181019092528281529291908301828280156102145780601f106101e957610100808354040283529160200191610214565b60405180806020018281038252838181518152602001915080519060200190808383829060006004602084601f0104600302600f01f150905090810190601f16801561013b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b505050565b8060006000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106101b557805160ff19168380011785555b506101499291505b808211156101e557600081556001016101a1565b82800160010185558215610199579182015b828111156101995782518260005055916020019190600101906101c7565b5090565b820191906000526020600020905b8154815290600101906020018083116101f757829003601f168201915b50505050508156`},
		[]string{`[{"constant":false,"inputs":[{"name":"value","type":"string"}],"name":"SetField","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"field","outputs":[{"name":"","type":"string"}],"type":"function"}]`},
		`
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			// Deploy a funky gas pattern contract
			_, _, limiter, err := DeployFunkyGasPattern(auth, sim)
			if err != nil {
				t.Fatalf("Failed to deploy funky contract: %v", err)
			}
			sim.Commit()

			// Set the field with automatic estimation and check that it succeeds
			if _, err := limiter.SetField(auth, "automatic"); err != nil {
				t.Fatalf("Failed to call automatically gased transaction: %v", err)
			}
			sim.Commit()

			if field, _ := limiter.Field(nil); field != "automatic" {
				t.Fatalf("Field mismatch: have %v, want %v", field, "automatic")
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Test that constant functions can be called from an (optional) specified address
	{
		`CallFrom`,
		`
			contract CallFrom {
				function callFrom() constant returns(address) {
					return msg.sender;
				}
			}
		`, []string{`6060604052346000575b6086806100176000396000f300606060405263ffffffff60e060020a60003504166349f8e98281146022575b6000565b34600057602c6055565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b335b905600a165627a7a72305820aef6b7685c0fa24ba6027e4870404a57df701473fe4107741805c19f5138417c0029`},
		[]string{`[{"constant":true,"inputs":[],"name":"callFrom","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"}]`},
		`
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/common"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			// Deploy a sender tester contract and execute a structured call on it
			_, _, callfrom, err := DeployCallFrom(auth, sim)
			if err != nil {
				t.Fatalf("Failed to deploy sender contract: %v", err)
			}
			sim.Commit()

			if res, err := callfrom.CallFrom(nil); err != nil {
				t.Errorf("Failed to call constant function: %v", err)
			} else if res != (common.Address{}) {
				t.Errorf("Invalid address returned, want: %x, got: %x", (common.Address{}), res)
			}

			for _, addr := range []common.Address{common.Address{}, common.Address{1}, common.Address{2}} {
				if res, err := callfrom.CallFrom(&bind.CallOpts{From: addr}); err != nil {
					t.Fatalf("Failed to call constant function: %v", err)
				} else if res != addr {
					t.Fatalf("Invalid address returned, want: %x, got: %x", addr, res)
				}
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Tests that methods and returns with underscores inside work correctly.
	{
		`Underscorer`,
		`
		contract Underscorer {
			function UnderscoredOutput() constant returns (int _int, string _string) {
				return (314, "pi");
			}
			function LowerLowerCollision() constant returns (int _res, int res) {
				return (1, 2);
			}
			function LowerUpperCollision() constant returns (int _res, int Res) {
				return (1, 2);
			}
			function UpperLowerCollision() constant returns (int _Res, int res) {
				return (1, 2);
			}
			function UpperUpperCollision() constant returns (int _Res, int Res) {
				return (1, 2);
			}
			function PurelyUnderscoredOutput() constant returns (int _, int res) {
				return (1, 2);
			}
			function AllPurelyUnderscoredOutput() constant returns (int _, int __) {
				return (1, 2);
			}
			function _under_scored_func() constant returns (int _int) {
				return 0;
			}
		}
		`, []string{`6060604052341561000f57600080fd5b6103858061001e6000396000f30060606040526004361061008e576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806303a592131461009357806346546dbe146100c357806367e6633d146100ec5780639df4848514610181578063af7486ab146101b1578063b564b34d146101e1578063e02ab24d14610211578063e409ca4514610241575b600080fd5b341561009e57600080fd5b6100a6610271565b604051808381526020018281526020019250505060405180910390f35b34156100ce57600080fd5b6100d6610286565b6040518082815260200191505060405180910390f35b34156100f757600080fd5b6100ff61028e565b6040518083815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561014557808201518184015260208101905061012a565b50505050905090810190601f1680156101725780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b341561018c57600080fd5b6101946102dc565b604051808381526020018281526020019250505060405180910390f35b34156101bc57600080fd5b6101c46102f1565b604051808381526020018281526020019250505060405180910390f35b34156101ec57600080fd5b6101f4610306565b604051808381526020018281526020019250505060405180910390f35b341561021c57600080fd5b61022461031b565b604051808381526020018281526020019250505060405180910390f35b341561024c57600080fd5b610254610330565b604051808381526020018281526020019250505060405180910390f35b60008060016002819150809050915091509091565b600080905090565b6000610298610345565b61013a8090506040805190810160405280600281526020017f7069000000000000000000000000000000000000000000000000000000000000815250915091509091565b60008060016002819150809050915091509091565b60008060016002819150809050915091509091565b60008060016002819150809050915091509091565b60008060016002819150809050915091509091565b60008060016002819150809050915091509091565b6020604051908101604052806000815250905600a165627a7a72305820d1a53d9de9d1e3d55cb3dc591900b63c4f1ded79114f7b79b332684840e186a40029`},
		[]string{`[{"constant":true,"inputs":[],"name":"LowerUpperCollision","outputs":[{"name":"_res","type":"int256"},{"name":"Res","type":"int256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"_under_scored_func","outputs":[{"name":"_int","type":"int256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"UnderscoredOutput","outputs":[{"name":"_int","type":"int256"},{"name":"_string","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"PurelyUnderscoredOutput","outputs":[{"name":"_","type":"int256"},{"name":"res","type":"int256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"UpperLowerCollision","outputs":[{"name":"_Res","type":"int256"},{"name":"res","type":"int256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"AllPurelyUnderscoredOutput","outputs":[{"name":"_","type":"int256"},{"name":"__","type":"int256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"UpperUpperCollision","outputs":[{"name":"_Res","type":"int256"},{"name":"Res","type":"int256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"LowerLowerCollision","outputs":[{"name":"_res","type":"int256"},{"name":"res","type":"int256"}],"payable":false,"stateMutability":"view","type":"function"}]`},
		`
			"fmt"
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			// Deploy a underscorer tester contract and execute a structured call on it
			_, _, underscorer, err := DeployUnderscorer(auth, sim)
			if err != nil {
				t.Fatalf("Failed to deploy underscorer contract: %v", err)
			}
			sim.Commit()

			// Verify that underscored return values correctly parse into structs
			if res, err := underscorer.UnderscoredOutput(nil); err != nil {
				t.Errorf("Failed to call constant function: %v", err)
			} else if res.Int.Cmp(big.NewInt(314)) != 0 || res.String != "pi" {
				t.Errorf("Invalid result, want: {314, \"pi\"}, got: %+v", res)
			}
			// Verify that underscored and non-underscored name collisions force tuple outputs
			var a, b *big.Int

			a, b, _ = underscorer.LowerLowerCollision(nil)
			a, b, _ = underscorer.LowerUpperCollision(nil)
			a, b, _ = underscorer.UpperLowerCollision(nil)
			a, b, _ = underscorer.UpperUpperCollision(nil)
			a, b, _ = underscorer.PurelyUnderscoredOutput(nil)
			a, b, _ = underscorer.AllPurelyUnderscoredOutput(nil)
			a, _ = underscorer.UnderScoredFunc(nil)

			fmt.Println(a, b, err)
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Tests that logs can be successfully filtered and decoded.
	{
		`Eventer`,
		`
			contract Eventer {
				event SimpleEvent (
					address indexed Addr,
					bytes32 indexed Id,
					bool    indexed Flag,
					uint    Value
				);
				function raiseSimpleEvent(address addr, bytes32 id, bool flag, uint value) {
					SimpleEvent(addr, id, flag, value);
				}

				event NodataEvent (
					uint   indexed Number,
					int16  indexed Short,
					uint32 indexed Long
				);
				function raiseNodataEvent(uint number, int16 short, uint32 long) {
					NodataEvent(number, short, long);
				}

				event DynamicEvent (
					string indexed IndexedString,
					bytes  indexed IndexedBytes,
					string NonIndexedString,
					bytes  NonIndexedBytes
				);
				function raiseDynamicEvent(string str, bytes blob) {
					DynamicEvent(str, blob, str, blob);
				}

				event FixedBytesEvent (
					bytes24 indexed IndexedBytes,
					bytes24 NonIndexedBytes
				);
				function raiseFixedBytesEvent(bytes24 blob) {
					FixedBytesEvent(blob, blob);
				}
			}
		`,
		[]string{`608060405234801561001057600080fd5b5061043f806100206000396000f3006080604052600436106100615763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663528300ff8114610066578063630c31e2146100ff5780636cc6b94014610138578063c7d116dd1461015b575b600080fd5b34801561007257600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100fd94369492936024939284019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497506101829650505050505050565b005b34801561010b57600080fd5b506100fd73ffffffffffffffffffffffffffffffffffffffff60043516602435604435151560643561033c565b34801561014457600080fd5b506100fd67ffffffffffffffff1960043516610394565b34801561016757600080fd5b506100fd60043560243560010b63ffffffff604435166103d6565b806040518082805190602001908083835b602083106101b25780518252601f199092019160209182019101610193565b51815160209384036101000a6000190180199092169116179052604051919093018190038120875190955087945090928392508401908083835b6020831061020b5780518252601f1990920191602091820191016101ec565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390207f3281fd4f5e152dd3385df49104a3f633706e21c9e80672e88d3bcddf33101f008484604051808060200180602001838103835285818151815260200191508051906020019080838360005b8381101561029c578181015183820152602001610284565b50505050905090810190601f1680156102c95780820380516001836020036101000a031916815260200191505b50838103825284518152845160209182019186019080838360005b838110156102fc5781810151838201526020016102e4565b50505050905090810190601f1680156103295780820380516001836020036101000a031916815260200191505b5094505050505060405180910390a35050565b60408051828152905183151591859173ffffffffffffffffffffffffffffffffffffffff8816917f1f097de4289df643bd9c11011cc61367aa12983405c021056e706eb5ba1250c8919081900360200190a450505050565b6040805167ffffffffffffffff19831680825291517fcdc4c1b1aed5524ffb4198d7a5839a34712baef5fa06884fac7559f4a5854e0a9181900360200190a250565b8063ffffffff168260010b847f3ca7f3a77e5e6e15e781850bc82e32adfa378a2a609370db24b4d0fae10da2c960405160405180910390a45050505600a165627a7a72305820468b5843bf653145bd924b323c64ef035d3dd922c170644b44d61aa666ea6eee0029`},
		[]string{`[{"constant":false,"inputs":[{"name":"str","type":"string"},{"name":"blob","type":"bytes"}],"name":"raiseDynamicEvent","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"addr","type":"address"},{"name":"id","type":"bytes32"},{"name":"flag","type":"bool"},{"name":"value","type":"uint256"}],"name":"raiseSimpleEvent","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"blob","type":"bytes24"}],"name":"raiseFixedBytesEvent","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"number","type":"uint256"},{"name":"short","type":"int16"},{"name":"long","type":"uint32"}],"name":"raiseNodataEvent","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"Addr","type":"address"},{"indexed":true,"name":"Id","type":"bytes32"},{"indexed":true,"name":"Flag","type":"bool"},{"indexed":false,"name":"Value","type":"uint256"}],"name":"SimpleEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"Number","type":"uint256"},{"indexed":true,"name":"Short","type":"int16"},{"indexed":true,"name":"Long","type":"uint32"}],"name":"NodataEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"IndexedString","type":"string"},{"indexed":true,"name":"IndexedBytes","type":"bytes"},{"indexed":false,"name":"NonIndexedString","type":"string"},{"indexed":false,"name":"NonIndexedBytes","type":"bytes"}],"name":"DynamicEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"IndexedBytes","type":"bytes24"},{"indexed":false,"name":"NonIndexedBytes","type":"bytes24"}],"name":"FixedBytesEvent","type":"event"}]`},
		`
			"math/big"
			"time"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/common"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			// Deploy an eventer contract
			_, _, eventer, err := DeployEventer(auth, sim)
			if err != nil {
				t.Fatalf("Failed to deploy eventer contract: %v", err)
			}
			sim.Commit()

			// Inject a few events into the contract, gradually more in each block
			for i := 1; i <= 3; i++ {
				for j := 1; j <= i; j++ {
					if _, err := eventer.RaiseSimpleEvent(auth, common.Address{byte(j)}, [32]byte{byte(j)}, true, big.NewInt(int64(10*i+j))); err != nil {
						t.Fatalf("block %d, event %d: raise failed: %v", i, j, err)
					}
				}
				sim.Commit()
			}
			// Test filtering for certain events and ensure they can be found
			sit, err := eventer.FilterSimpleEvent(nil, []common.Address{common.Address{1}, common.Address{3}}, [][32]byte{{byte(1)}, {byte(2)}, {byte(3)}}, []bool{true})
			if err != nil {
				t.Fatalf("failed to filter for simple events: %v", err)
			}
			defer sit.Close()

			sit.Next()
			if sit.Event.Value.Uint64() != 11 || !sit.Event.Flag {
				t.Errorf("simple log content mismatch: have %v, want {11, true}", sit.Event)
			}
			sit.Next()
			if sit.Event.Value.Uint64() != 21 || !sit.Event.Flag {
				t.Errorf("simple log content mismatch: have %v, want {21, true}", sit.Event)
			}
			sit.Next()
			if sit.Event.Value.Uint64() != 31 || !sit.Event.Flag {
				t.Errorf("simple log content mismatch: have %v, want {31, true}", sit.Event)
			}
			sit.Next()
			if sit.Event.Value.Uint64() != 33 || !sit.Event.Flag {
				t.Errorf("simple log content mismatch: have %v, want {33, true}", sit.Event)
			}

			if sit.Next() {
				t.Errorf("unexpected simple event found: %+v", sit.Event)
			}
			if err = sit.Error(); err != nil {
				t.Fatalf("simple event iteration failed: %v", err)
			}
			// Test raising and filtering for an event with no data component
			if _, err := eventer.RaiseNodataEvent(auth, big.NewInt(314), 141, 271); err != nil {
				t.Fatalf("failed to raise nodata event: %v", err)
			}
			sim.Commit()

			nit, err := eventer.FilterNodataEvent(nil, []*big.Int{big.NewInt(314)}, []int16{140, 141, 142}, []uint32{271})
			if err != nil {
				t.Fatalf("failed to filter for nodata events: %v", err)
			}
			defer nit.Close()

			if !nit.Next() {
				t.Fatalf("nodata log not found: %v", nit.Error())
			}
			if nit.Event.Number.Uint64() != 314 {
				t.Errorf("nodata log content mismatch: have %v, want 314", nit.Event.Number)
			}
			if nit.Next() {
				t.Errorf("unexpected nodata event found: %+v", nit.Event)
			}
			if err = nit.Error(); err != nil {
				t.Fatalf("nodata event iteration failed: %v", err)
			}
			// Test raising and filtering for events with dynamic indexed components
			if _, err := eventer.RaiseDynamicEvent(auth, "Hello", []byte("World")); err != nil {
				t.Fatalf("failed to raise dynamic event: %v", err)
			}
			sim.Commit()

			dit, err := eventer.FilterDynamicEvent(nil, []string{"Hi", "Hello", "Bye"}, [][]byte{[]byte("World")})
			if err != nil {
				t.Fatalf("failed to filter for dynamic events: %v", err)
			}
			defer dit.Close()

			if !dit.Next() {
				t.Fatalf("dynamic log not found: %v", dit.Error())
			}
			if dit.Event.NonIndexedString != "Hello" || string(dit.Event.NonIndexedBytes) != "World" ||	dit.Event.IndexedString != common.HexToHash("0x06b3dfaec148fb1bb2b066f10ec285e7c9bf402ab32aa78a5d38e34566810cd2") || dit.Event.IndexedBytes != common.HexToHash("0xf2208c967df089f60420785795c0a9ba8896b0f6f1867fa7f1f12ad6f79c1a18") {
				t.Errorf("dynamic log content mismatch: have %v, want {'0x06b3dfaec148fb1bb2b066f10ec285e7c9bf402ab32aa78a5d38e34566810cd2, '0xf2208c967df089f60420785795c0a9ba8896b0f6f1867fa7f1f12ad6f79c1a18', 'Hello', 'World'}", dit.Event)
			}
			if dit.Next() {
				t.Errorf("unexpected dynamic event found: %+v", dit.Event)
			}
			if err = dit.Error(); err != nil {
				t.Fatalf("dynamic event iteration failed: %v", err)
			}
			// Test raising and filtering for events with fixed bytes components
			var fblob [24]byte
			copy(fblob[:], []byte("Fixed Bytes"))

			if _, err := eventer.RaiseFixedBytesEvent(auth, fblob); err != nil {
				t.Fatalf("failed to raise fixed bytes event: %v", err)
			}
			sim.Commit()

			fit, err := eventer.FilterFixedBytesEvent(nil, [][24]byte{fblob})
			if err != nil {
				t.Fatalf("failed to filter for fixed bytes events: %v", err)
			}
			defer fit.Close()

			if !fit.Next() {
				t.Fatalf("fixed bytes log not found: %v", fit.Error())
			}
			if fit.Event.NonIndexedBytes != fblob || fit.Event.IndexedBytes != fblob {
				t.Errorf("fixed bytes log content mismatch: have %v, want {'%x', '%x'}", fit.Event, fblob, fblob)
			}
			if fit.Next() {
				t.Errorf("unexpected fixed bytes event found: %+v", fit.Event)
			}
			if err = fit.Error(); err != nil {
				t.Fatalf("fixed bytes event iteration failed: %v", err)
			}
			// Test subscribing to an event and raising it afterwards
			ch := make(chan *EventerSimpleEvent, 16)
			sub, err := eventer.WatchSimpleEvent(nil, ch, nil, nil, nil)
			if err != nil {
				t.Fatalf("failed to subscribe to simple events: %v", err)
			}
			if _, err := eventer.RaiseSimpleEvent(auth, common.Address{255}, [32]byte{255}, true, big.NewInt(255)); err != nil {
				t.Fatalf("failed to raise subscribed simple event: %v", err)
			}
			sim.Commit()

			select {
			case event := <-ch:
				if event.Value.Uint64() != 255 {
					t.Errorf("simple log content mismatch: have %v, want 255", event)
				}
			case <-time.After(250 * time.Millisecond):
				t.Fatalf("subscribed simple event didn't arrive")
			}
			// Unsubscribe from the event and make sure we're not delivered more
			sub.Unsubscribe()

			if _, err := eventer.RaiseSimpleEvent(auth, common.Address{254}, [32]byte{254}, true, big.NewInt(254)); err != nil {
				t.Fatalf("failed to raise subscribed simple event: %v", err)
			}
			sim.Commit()

			select {
			case event := <-ch:
				t.Fatalf("unsubscribed simple event arrived: %v", event)
			case <-time.After(250 * time.Millisecond):
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	{
		`DeeplyNestedArray`,
		`
			contract DeeplyNestedArray {
				uint64[3][4][5] public deepUint64Array;
				function storeDeepUintArray(uint64[3][4][5] arr) public {
					deepUint64Array = arr;
				}
				function retrieveDeepArray() public view returns (uint64[3][4][5]) {
					return deepUint64Array;
				}
			}
		`,
		[]string{`6060604052341561000f57600080fd5b6106438061001e6000396000f300606060405260043610610057576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063344248551461005c5780638ed4573a1461011457806398ed1856146101ab575b600080fd5b341561006757600080fd5b610112600480806107800190600580602002604051908101604052809291906000905b828210156101055783826101800201600480602002604051908101604052809291906000905b828210156100f25783826060020160038060200260405190810160405280929190826003602002808284378201915050505050815260200190600101906100b0565b505050508152602001906001019061008a565b5050505091905050610208565b005b341561011f57600080fd5b61012761021d565b604051808260056000925b8184101561019b578284602002015160046000925b8184101561018d5782846020020151600360200280838360005b8381101561017c578082015181840152602081019050610161565b505050509050019260010192610147565b925050509260010192610132565b9250505091505060405180910390f35b34156101b657600080fd5b6101de6004808035906020019091908035906020019091908035906020019091905050610309565b604051808267ffffffffffffffff1667ffffffffffffffff16815260200191505060405180910390f35b80600090600561021992919061035f565b5050565b6102256103b0565b6000600580602002604051908101604052809291906000905b8282101561030057838260040201600480602002604051908101604052809291906000905b828210156102ed578382016003806020026040519081016040528092919082600380156102d9576020028201916000905b82829054906101000a900467ffffffffffffffff1667ffffffffffffffff16815260200190600801906020826007010492830192600103820291508084116102945790505b505050505081526020019060010190610263565b505050508152602001906001019061023e565b50505050905090565b60008360058110151561031857fe5b600402018260048110151561032957fe5b018160038110151561033757fe5b6004918282040191900660080292509250509054906101000a900467ffffffffffffffff1681565b826005600402810192821561039f579160200282015b8281111561039e5782518290600461038e9291906103df565b5091602001919060040190610375565b5b5090506103ac919061042d565b5090565b610780604051908101604052806005905b6103c9610459565b8152602001906001900390816103c15790505090565b826004810192821561041c579160200282015b8281111561041b5782518290600361040b929190610488565b50916020019190600101906103f2565b5b5090506104299190610536565b5090565b61045691905b8082111561045257600081816104499190610562565b50600401610433565b5090565b90565b610180604051908101604052806004905b6104726105a7565b81526020019060019003908161046a5790505090565b82600380016004900481019282156105255791602002820160005b838211156104ef57835183826101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555092602001926008016020816007010492830192600103026104a3565b80156105235782816101000a81549067ffffffffffffffff02191690556008016020816007010492830192600103026104ef565b505b50905061053291906105d9565b5090565b61055f91905b8082111561055b57600081816105529190610610565b5060010161053c565b5090565b90565b50600081816105719190610610565b50600101600081816105839190610610565b50600101600081816105959190610610565b5060010160006105a59190610610565b565b6060604051908101604052806003905b600067ffffffffffffffff168152602001906001900390816105b75790505090565b61060d91905b8082111561060957600081816101000a81549067ffffffffffffffff0219169055506001016105df565b5090565b90565b50600090555600a165627a7a7230582087e5a43f6965ab6ef7a4ff056ab80ed78fd8c15cff57715a1bf34ec76a93661c0029`},
		[]string{`[{"constant":false,"inputs":[{"name":"arr","type":"uint64[3][4][5]"}],"name":"storeDeepUintArray","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"retrieveDeepArray","outputs":[{"name":"","type":"uint64[3][4][5]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"},{"name":"","type":"uint256"},{"name":"","type":"uint256"}],"name":"deepUint64Array","outputs":[{"name":"","type":"uint64"}],"payable":false,"stateMutability":"view","type":"function"}]`},
		`
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			//deploy the test contract
			_, _, testContract, err := DeployDeeplyNestedArray(auth, sim)
			if err != nil {
				t.Fatalf("Failed to deploy test contract: %v", err)
			}

			// Finish deploy.
			sim.Commit()

			//Create coordinate-filled array, for testing purposes.
			testArr := [5][4][3]uint64{}
			for i := 0; i < 5; i++ {
				testArr[i] = [4][3]uint64{}
				for j := 0; j < 4; j++ {
					testArr[i][j] = [3]uint64{}
					for k := 0; k < 3; k++ {
						//pack the coordinates, each array value will be unique, and can be validated easily.
						testArr[i][j][k] = uint64(i) << 16 | uint64(j) << 8 | uint64(k)
					}
				}
			}

			if _, err := testContract.StoreDeepUintArray(&bind.TransactOpts{
				From: auth.From,
				Signer: auth.Signer,
			}, testArr); err != nil {
				t.Fatalf("Failed to store nested array in test contract: %v", err)
			}

			sim.Commit()

			retrievedArr, err := testContract.RetrieveDeepArray(&bind.CallOpts{
				From: auth.From,
				Pending: false,
			})
			if err != nil {
				t.Fatalf("Failed to retrieve nested array from test contract: %v", err)
			}

			//quick check to see if contents were copied
			// (See accounts/abi/unpack_test.go for more extensive testing)
			if retrievedArr[4][3][2] != testArr[4][3][2] {
				t.Fatalf("Retrieved value does not match expected value! got: %d, expected: %d. %v", retrievedArr[4][3][2], testArr[4][3][2], err)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	{
		`CallbackParam`,
		`
			contract FunctionPointerTest {
				function test(function(uint256) external callback) external {
					callback(1);
				}
			}
		`,
		[]string{`608060405234801561001057600080fd5b5061015e806100206000396000f3fe60806040526004361061003b576000357c010000000000000000000000000000000000000000000000000000000090048063d7a5aba214610040575b600080fd5b34801561004c57600080fd5b506100be6004803603602081101561006357600080fd5b810190808035806c0100000000000000000000000090049068010000000000000000900463ffffffff1677ffffffffffffffffffffffffffffffffffffffffffffffff169091602001919093929190939291905050506100c0565b005b818160016040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180828152602001915050600060405180830381600087803b15801561011657600080fd5b505af115801561012a573d6000803e3d6000fd5b50505050505056fea165627a7a7230582062f87455ff84be90896dbb0c4e4ddb505c600d23089f8e80a512548440d7e2580029`},
		[]string{`[
			{
				"constant": false,
				"inputs": [
					{
						"name": "callback",
						"type": "function"
					}
				],
				"name": "test",
				"outputs": [],
				"payable": false,
				"stateMutability": "nonpayable",
				"type": "function"
			}
		]`}, `
			"strings"
		`,
		`
			if strings.Compare("test(function)", CallbackParamFuncSigs["d7a5aba2"]) != 0 {
				t.Fatalf("")
			}
		`,
		[]map[string]string{
			{
				"test(function)": "d7a5aba2",
			},
		},
		nil,
		nil,
		nil,
	}, {
		`Tuple`,
		`
		pragma solidity >=0.4.19 <0.6.0;
		pragma experimental ABIEncoderV2;

		contract Tuple {
			struct S { uint a; uint[] b; T[] c; }
			struct T { uint x; uint y; }
			struct P { uint8 x; uint8 y; }
			struct Q { uint16 x; uint16 y; }
			event TupleEvent(S a, T[2][] b, T[][2] c, S[] d, uint[] e);
			event TupleEvent2(P[]);

			function func1(S memory a, T[2][] memory b, T[][2] memory c, S[] memory d, uint[] memory e) public pure returns (S memory, T[2][] memory, T[][2] memory, S[] memory, uint[] memory) {
				return (a, b, c, d, e);
			}
			function func2(S memory a, T[2][] memory b, T[][2] memory c, S[] memory d, uint[] memory e) public {
				emit TupleEvent(a, b, c, d, e);
			}
			function func3(Q[] memory) public pure {} // call function, nothing to return
		}
		`,
		[]string{`60806040523480156100115760006000fd5b50610017565b6110b2806100266000396000f3fe60806040523480156100115760006000fd5b50600436106100465760003560e01c8063443c79b41461004c578063d0062cdd14610080578063e4d9a43b1461009c57610046565b60006000fd5b610066600480360361006191908101906107b8565b6100b8565b604051610077959493929190610ccb565b60405180910390f35b61009a600480360361009591908101906107b8565b6100ef565b005b6100b660048036036100b19190810190610775565b610136565b005b6100c061013a565b60606100ca61015e565b606060608989898989945094509450945094506100e2565b9550955095509550959050565b7f18d6e66efa53739ca6d13626f35ebc700b31cced3eddb50c70bbe9c082c6cd008585858585604051610126959493929190610ccb565b60405180910390a15b5050505050565b5b50565b60405180606001604052806000815260200160608152602001606081526020015090565b60405180604001604052806002905b606081526020019060019003908161016d57905050905661106e565b600082601f830112151561019d5760006000fd5b81356101b06101ab82610d6f565b610d41565b915081818352602084019350602081019050838560808402820111156101d65760006000fd5b60005b8381101561020757816101ec888261037a565b8452602084019350608083019250505b6001810190506101d9565b5050505092915050565b600082601f83011215156102255760006000fd5b600261023861023382610d98565b610d41565b9150818360005b83811015610270578135860161025588826103f3565b8452602084019350602083019250505b60018101905061023f565b5050505092915050565b600082601f830112151561028e5760006000fd5b81356102a161029c82610dbb565b610d41565b915081818352602084019350602081019050838560408402820111156102c75760006000fd5b60005b838110156102f857816102dd888261058b565b8452602084019350604083019250505b6001810190506102ca565b5050505092915050565b600082601f83011215156103165760006000fd5b813561032961032482610de4565b610d41565b9150818183526020840193506020810190508360005b83811015610370578135860161035588826105d8565b8452602084019350602083019250505b60018101905061033f565b5050505092915050565b600082601f830112151561038e5760006000fd5b60026103a161039c82610e0d565b610d41565b915081838560408402820111156103b85760006000fd5b60005b838110156103e957816103ce88826106fe565b8452602084019350604083019250505b6001810190506103bb565b5050505092915050565b600082601f83011215156104075760006000fd5b813561041a61041582610e30565b610d41565b915081818352602084019350602081019050838560408402820111156104405760006000fd5b60005b83811015610471578161045688826106fe565b8452602084019350604083019250505b600181019050610443565b5050505092915050565b600082601f830112151561048f5760006000fd5b81356104a261049d82610e59565b610d41565b915081818352602084019350602081019050838560208402820111156104c85760006000fd5b60005b838110156104f957816104de8882610760565b8452602084019350602083019250505b6001810190506104cb565b5050505092915050565b600082601f83011215156105175760006000fd5b813561052a61052582610e82565b610d41565b915081818352602084019350602081019050838560208402820111156105505760006000fd5b60005b8381101561058157816105668882610760565b8452602084019350602083019250505b600181019050610553565b5050505092915050565b60006040828403121561059e5760006000fd5b6105a86040610d41565b905060006105b88482850161074b565b60008301525060206105cc8482850161074b565b60208301525092915050565b6000606082840312156105eb5760006000fd5b6105f56060610d41565b9050600061060584828501610760565b600083015250602082013567ffffffffffffffff8111156106265760006000fd5b6106328482850161047b565b602083015250604082013567ffffffffffffffff8111156106535760006000fd5b61065f848285016103f3565b60408301525092915050565b60006060828403121561067e5760006000fd5b6106886060610d41565b9050600061069884828501610760565b600083015250602082013567ffffffffffffffff8111156106b95760006000fd5b6106c58482850161047b565b602083015250604082013567ffffffffffffffff8111156106e65760006000fd5b6106f2848285016103f3565b60408301525092915050565b6000604082840312156107115760006000fd5b61071b6040610d41565b9050600061072b84828501610760565b600083015250602061073f84828501610760565b60208301525092915050565b60008135905061075a8161103a565b92915050565b60008135905061076f81611054565b92915050565b6000602082840312156107885760006000fd5b600082013567ffffffffffffffff8111156107a35760006000fd5b6107af8482850161027a565b91505092915050565b6000600060006000600060a086880312156107d35760006000fd5b600086013567ffffffffffffffff8111156107ee5760006000fd5b6107fa8882890161066b565b955050602086013567ffffffffffffffff8111156108185760006000fd5b61082488828901610189565b945050604086013567ffffffffffffffff8111156108425760006000fd5b61084e88828901610211565b935050606086013567ffffffffffffffff81111561086c5760006000fd5b61087888828901610302565b925050608086013567ffffffffffffffff8111156108965760006000fd5b6108a288828901610503565b9150509295509295909350565b60006108bb8383610a6a565b60808301905092915050565b60006108d38383610ac2565b905092915050565b60006108e78383610c36565b905092915050565b60006108fb8383610c8d565b60408301905092915050565b60006109138383610cbc565b60208301905092915050565b600061092a82610f0f565b6109348185610fb7565b935061093f83610eab565b8060005b8381101561097157815161095788826108af565b975061096283610f5c565b9250505b600181019050610943565b5085935050505092915050565b600061098982610f1a565b6109938185610fc8565b9350836020820285016109a585610ebb565b8060005b858110156109e257848403895281516109c285826108c7565b94506109cd83610f69565b925060208a019950505b6001810190506109a9565b50829750879550505050505092915050565b60006109ff82610f25565b610a098185610fd3565b935083602082028501610a1b85610ec5565b8060005b85811015610a585784840389528151610a3885826108db565b9450610a4383610f76565b925060208a019950505b600181019050610a1f565b50829750879550505050505092915050565b610a7381610f30565b610a7d8184610fe4565b9250610a8882610ed5565b8060005b83811015610aba578151610aa087826108ef565b9650610aab83610f83565b9250505b600181019050610a8c565b505050505050565b6000610acd82610f3b565b610ad78185610fef565b9350610ae283610edf565b8060005b83811015610b14578151610afa88826108ef565b9750610b0583610f90565b9250505b600181019050610ae6565b5085935050505092915050565b6000610b2c82610f51565b610b368185611011565b9350610b4183610eff565b8060005b83811015610b73578151610b598882610907565b9750610b6483610faa565b9250505b600181019050610b45565b5085935050505092915050565b6000610b8b82610f46565b610b958185611000565b9350610ba083610eef565b8060005b83811015610bd2578151610bb88882610907565b9750610bc383610f9d565b9250505b600181019050610ba4565b5085935050505092915050565b6000606083016000830151610bf76000860182610cbc565b5060208301518482036020860152610c0f8282610b80565b91505060408301518482036040860152610c298282610ac2565b9150508091505092915050565b6000606083016000830151610c4e6000860182610cbc565b5060208301518482036020860152610c668282610b80565b91505060408301518482036040860152610c808282610ac2565b9150508091505092915050565b604082016000820151610ca36000850182610cbc565b506020820151610cb66020850182610cbc565b50505050565b610cc581611030565b82525050565b600060a0820190508181036000830152610ce58188610bdf565b90508181036020830152610cf9818761091f565b90508181036040830152610d0d818661097e565b90508181036060830152610d2181856109f4565b90508181036080830152610d358184610b21565b90509695505050505050565b6000604051905081810181811067ffffffffffffffff82111715610d655760006000fd5b8060405250919050565b600067ffffffffffffffff821115610d875760006000fd5b602082029050602081019050919050565b600067ffffffffffffffff821115610db05760006000fd5b602082029050919050565b600067ffffffffffffffff821115610dd35760006000fd5b602082029050602081019050919050565b600067ffffffffffffffff821115610dfc5760006000fd5b602082029050602081019050919050565b600067ffffffffffffffff821115610e255760006000fd5b602082029050919050565b600067ffffffffffffffff821115610e485760006000fd5b602082029050602081019050919050565b600067ffffffffffffffff821115610e715760006000fd5b602082029050602081019050919050565b600067ffffffffffffffff821115610e9a5760006000fd5b602082029050602081019050919050565b6000819050602082019050919050565b6000819050919050565b6000819050602082019050919050565b6000819050919050565b6000819050602082019050919050565b6000819050602082019050919050565b6000819050602082019050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600081519050919050565b600081519050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b600082825260208201905092915050565b600082825260208201905092915050565b600061ffff82169050919050565b6000819050919050565b61104381611022565b811415156110515760006000fd5b50565b61105d81611030565b8114151561106b5760006000fd5b50565bfea365627a7a72315820d78c6ba7ee332581e6c4d9daa5fc07941841230f7ce49edf6e05b1b63853e8746c6578706572696d656e74616cf564736f6c634300050c0040`},
		[]string{`
[{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"a","type":"uint256"},{"internalType":"uint256[]","name":"b","type":"uint256[]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[]","name":"c","type":"tuple[]"}],"indexed":false,"internalType":"struct Tuple.S","name":"a","type":"tuple"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"indexed":false,"internalType":"struct Tuple.T[2][]","name":"b","type":"tuple[2][]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"indexed":false,"internalType":"struct Tuple.T[][2]","name":"c","type":"tuple[][2]"},{"components":[{"internalType":"uint256","name":"a","type":"uint256"},{"internalType":"uint256[]","name":"b","type":"uint256[]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[]","name":"c","type":"tuple[]"}],"indexed":false,"internalType":"struct Tuple.S[]","name":"d","type":"tuple[]"},{"indexed":false,"internalType":"uint256[]","name":"e","type":"uint256[]"}],"name":"TupleEvent","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"uint8","name":"x","type":"uint8"},{"internalType":"uint8","name":"y","type":"uint8"}],"indexed":false,"internalType":"struct Tuple.P[]","name":"","type":"tuple[]"}],"name":"TupleEvent2","type":"event"},{"constant":true,"inputs":[{"components":[{"internalType":"uint256","name":"a","type":"uint256"},{"internalType":"uint256[]","name":"b","type":"uint256[]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[]","name":"c","type":"tuple[]"}],"internalType":"struct Tuple.S","name":"a","type":"tuple"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[2][]","name":"b","type":"tuple[2][]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[][2]","name":"c","type":"tuple[][2]"},{"components":[{"internalType":"uint256","name":"a","type":"uint256"},{"internalType":"uint256[]","name":"b","type":"uint256[]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[]","name":"c","type":"tuple[]"}],"internalType":"struct Tuple.S[]","name":"d","type":"tuple[]"},{"internalType":"uint256[]","name":"e","type":"uint256[]"}],"name":"func1","outputs":[{"components":[{"internalType":"uint256","name":"a","type":"uint256"},{"internalType":"uint256[]","name":"b","type":"uint256[]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[]","name":"c","type":"tuple[]"}],"internalType":"struct Tuple.S","name":"","type":"tuple"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[2][]","name":"","type":"tuple[2][]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[][2]","name":"","type":"tuple[][2]"},{"components":[{"internalType":"uint256","name":"a","type":"uint256"},{"internalType":"uint256[]","name":"b","type":"uint256[]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[]","name":"c","type":"tuple[]"}],"internalType":"struct Tuple.S[]","name":"","type":"tuple[]"},{"internalType":"uint256[]","name":"","type":"uint256[]"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":false,"inputs":[{"components":[{"internalType":"uint256","name":"a","type":"uint256"},{"internalType":"uint256[]","name":"b","type":"uint256[]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[]","name":"c","type":"tuple[]"}],"internalType":"struct Tuple.S","name":"a","type":"tuple"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[2][]","name":"b","type":"tuple[2][]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[][2]","name":"c","type":"tuple[][2]"},{"components":[{"internalType":"uint256","name":"a","type":"uint256"},{"internalType":"uint256[]","name":"b","type":"uint256[]"},{"components":[{"internalType":"uint256","name":"x","type":"uint256"},{"internalType":"uint256","name":"y","type":"uint256"}],"internalType":"struct Tuple.T[]","name":"c","type":"tuple[]"}],"internalType":"struct Tuple.S[]","name":"d","type":"tuple[]"},{"internalType":"uint256[]","name":"e","type":"uint256[]"}],"name":"func2","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"components":[{"internalType":"uint16","name":"x","type":"uint16"},{"internalType":"uint16","name":"y","type":"uint16"}],"internalType":"struct Tuple.Q[]","name":"","type":"tuple[]"}],"name":"func3","outputs":[],"payable":false,"stateMutability":"pure","type":"function"}]
		`},
		`
			"math/big"
			"reflect"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,

		`
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			_, _, contract, err := DeployTuple(auth, sim)
			if err != nil {
				t.Fatalf("deploy contract failed %v", err)
			}
			sim.Commit()

			check := func(a, b interface{}, errMsg string) {
				if !reflect.DeepEqual(a, b) {
					t.Fatal(errMsg)
				}
			}

			a := TupleS{
				A: big.NewInt(1),
				B: []*big.Int{big.NewInt(2), big.NewInt(3)},
				C: []TupleT{
					{
						X: big.NewInt(4),
						Y: big.NewInt(5),
					},
					{
						X: big.NewInt(6),
						Y: big.NewInt(7),
					},
				},
			}

			b := [][2]TupleT{
				{
					{
						X: big.NewInt(8),
						Y: big.NewInt(9),
					},
					{
						X: big.NewInt(10),
						Y: big.NewInt(11),
					},
				},
			}

			c := [2][]TupleT{
				{
					{
						X: big.NewInt(12),
						Y: big.NewInt(13),
					},
					{
						X: big.NewInt(14),
						Y: big.NewInt(15),
					},
				},
				{
					{
						X: big.NewInt(16),
						Y: big.NewInt(17),
					},
				},
			}

			d := []TupleS{a}

			e := []*big.Int{big.NewInt(18), big.NewInt(19)}
			ret1, ret2, ret3, ret4, ret5, err := contract.Func1(nil, a, b, c, d, e)
			if err != nil {
				t.Fatalf("invoke contract failed, err %v", err)
			}
			check(ret1, a, "ret1 mismatch")
			check(ret2, b, "ret2 mismatch")
			check(ret3, c, "ret3 mismatch")
			check(ret4, d, "ret4 mismatch")
			check(ret5, e, "ret5 mismatch")

			_, err = contract.Func2(auth, a, b, c, d, e)
			if err != nil {
				t.Fatalf("invoke contract failed, err %v", err)
			}
			sim.Commit()

			iter, err := contract.FilterTupleEvent(nil)
			if err != nil {
				t.Fatalf("failed to create event filter, err %v", err)
			}
			defer iter.Close()

			iter.Next()
			check(iter.Event.A, a, "field1 mismatch")
			check(iter.Event.B, b, "field2 mismatch")
			check(iter.Event.C, c, "field3 mismatch")
			check(iter.Event.D, d, "field4 mismatch")
			check(iter.Event.E, e, "field5 mismatch")

			err = contract.Func3(nil, nil)
			if err != nil {
				t.Fatalf("failed to call function which has no return, err %v", err)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	{
		`UseLibrary`,
		`
		library Math {
    		function add(uint a, uint b) public view returns(uint) {
        		return a + b;
    		}
		}

		contract UseLibrary {
			function add (uint c, uint d) public view returns(uint) {
        		return Math.add(c,d);
    		}
		}
		`,
		[]string{
			// Bytecode for the UseLibrary contract
			`608060405234801561001057600080fd5b5061011d806100206000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063771602f714602d575b600080fd5b604d60048036036040811015604157600080fd5b5080359060200135605f565b60408051918252519081900360200190f35b600073__$b98c933f0a6ececcd167bd4f9d3299b1a0$__63771602f784846040518363ffffffff1660e01b8152600401808381526020018281526020019250505060206040518083038186803b15801560b757600080fd5b505af415801560ca573d6000803e3d6000fd5b505050506040513d602081101560df57600080fd5b5051939250505056fea265627a7a72305820eb5c38f42445604cfa43d85e3aa5ecc48b0a646456c902dd48420ae7241d06f664736f6c63430005090032`,
			// Bytecode for the Math contract
			`60a3610024600b82828239805160001a607314601757fe5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361060335760003560e01c8063771602f7146038575b600080fd5b605860048036036040811015604c57600080fd5b5080359060200135606a565b60408051918252519081900360200190f35b019056fea265627a7a723058206fc6c05f3078327f9c763edffdb5ab5f8bd212e293a1306c7d0ad05af3ad35f464736f6c63430005090032`,
		},
		[]string{
			`[{"constant":true,"inputs":[{"name":"c","type":"uint256"},{"name":"d","type":"uint256"}],"name":"add","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`,
			`[{"constant":true,"inputs":[{"name":"a","type":"uint256"},{"name":"b","type":"uint256"}],"name":"add","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`,
		},
		`
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			//deploy the test contract
			_, _, testContract, err := DeployUseLibrary(auth, sim)
			if err != nil {
				t.Fatalf("Failed to deploy test contract: %v", err)
			}

			// Finish deploy.
			sim.Commit()

			// Check that the library contract has been deployed
			// by calling the contract's add function.
			res, err := testContract.Add(&bind.CallOpts{
				From: auth.From,
				Pending: false,
			}, big.NewInt(1), big.NewInt(2))
			if err != nil {
				t.Fatalf("Failed to call linked contract: %v", err)
			}
			if res.Cmp(big.NewInt(3)) != 0 {
				t.Fatalf("Add did not return the correct result: %d != %d", res, 3)
			}
		`,
		nil,
		map[string]string{
			"b98c933f0a6ececcd167bd4f9d3299b1a0": "Math",
		},
		nil,
		[]string{"UseLibrary", "Math"},
	}, {
		"Overload",
		`
		pragma solidity ^0.5.10;

		contract overload {
		  mapping(address => uint256) balances;

		  event bar(uint256 i);
		  event bar(uint256 i, uint256 j);

		  function foo(uint256 i) public {
			  emit bar(i);
		  }
		  function foo(uint256 i, uint256 j) public {
			  emit bar(i, j);
		  }
		}
		`,
		[]string{`608060405234801561001057600080fd5b50610153806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806304bc52f81461003b5780632fbebd3814610073575b600080fd5b6100716004803603604081101561005157600080fd5b8101908080359060200190929190803590602001909291905050506100a1565b005b61009f6004803603602081101561008957600080fd5b81019080803590602001909291905050506100e4565b005b7fae42e9514233792a47a1e4554624e83fe852228e1503f63cd383e8a431f4f46d8282604051808381526020018281526020019250505060405180910390a15050565b7f0423a1321222a0a8716c22b92fac42d85a45a612b696a461784d9fa537c81e5c816040518082815260200191505060405180910390a15056fea265627a7a72305820e22b049858b33291cbe67eeaece0c5f64333e439d27032ea8337d08b1de18fe864736f6c634300050a0032`},
		[]string{`[{"constant":false,"inputs":[{"name":"i","type":"uint256"},{"name":"j","type":"uint256"}],"name":"foo","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"i","type":"uint256"}],"name":"foo","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"i","type":"uint256"}],"name":"bar","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"i","type":"uint256"},{"indexed":false,"name":"j","type":"uint256"}],"name":"bar","type":"event"}]`},
		`
		"math/big"
		"time"

		"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
		"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
		"github.com/ubiq/go-ubiq/v6/core"
		"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
		// Initialize test accounts
		key, _ := crypto.GenerateKey()
		auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
		defer sim.Close()

		// deploy the test contract
		_, _, contract, err := DeployOverload(auth, sim)
		if err != nil {
			t.Fatalf("Failed to deploy contract: %v", err)
		}
		// Finish deploy.
		sim.Commit()

		resCh, stopCh := make(chan uint64), make(chan struct{})

		go func() {
			barSink := make(chan *OverloadBar)
			sub, _ := contract.WatchBar(nil, barSink)
			defer sub.Unsubscribe()

			bar0Sink := make(chan *OverloadBar0)
			sub0, _ := contract.WatchBar0(nil, bar0Sink)
			defer sub0.Unsubscribe()

			for {
				select {
				case ev := <-barSink:
					resCh <- ev.I.Uint64()
				case ev := <-bar0Sink:
					resCh <- ev.I.Uint64() + ev.J.Uint64()
				case <-stopCh:
					return
				}
			}
		}()
		contract.Foo(auth, big.NewInt(1), big.NewInt(2))
		sim.Commit()
		select {
		case n := <-resCh:
			if n != 3 {
				t.Fatalf("Invalid bar0 event")
			}
		case <-time.NewTimer(3 * time.Second).C:
			t.Fatalf("Wait bar0 event timeout")
		}

		contract.Foo0(auth, big.NewInt(1))
		sim.Commit()
		select {
		case n := <-resCh:
			if n != 1 {
				t.Fatalf("Invalid bar event")
			}
		case <-time.NewTimer(3 * time.Second).C:
			t.Fatalf("Wait bar event timeout")
		}
		close(stopCh)
		`,
		nil,
		nil,
		nil,
		nil,
	},
	{
		"IdentifierCollision",
		`
		pragma solidity >=0.4.19 <0.6.0;

		contract IdentifierCollision {
			uint public _myVar;

			function MyVar() public view returns (uint) {
				return _myVar;
			}
		}
		`,
		[]string{"60806040523480156100115760006000fd5b50610017565b60c3806100256000396000f3fe608060405234801560105760006000fd5b506004361060365760003560e01c806301ad4d8714603c5780634ef1f0ad146058576036565b60006000fd5b60426074565b6040518082815260200191505060405180910390f35b605e607d565b6040518082815260200191505060405180910390f35b60006000505481565b60006000600050549050608b565b9056fea265627a7a7231582067c8d84688b01c4754ba40a2a871cede94ea1f28b5981593ab2a45b46ac43af664736f6c634300050c0032"},
		[]string{`[{"constant":true,"inputs":[],"name":"MyVar","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"_myVar","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`},
		`
		"math/big"

		"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
		"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
		"github.com/ubiq/go-ubiq/v6/crypto"
		"github.com/ubiq/go-ubiq/v6/core"
		`,
		`
		// Initialize test accounts
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)

		// Deploy registrar contract
		sim := backends.NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: big.NewInt(200000000000000000)}}, 10000000)
		defer sim.Close()

		transactOpts, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		_, _, _, err := DeployIdentifierCollision(transactOpts, sim)
		if err != nil {
			t.Fatalf("failed to deploy contract: %v", err)
		}
		`,
		nil,
		nil,
		map[string]string{"_myVar": "pubVar"}, // alias MyVar to PubVar
		nil,
	},
	{
		"MultiContracts",
		`
		pragma solidity ^0.5.11;
		pragma experimental ABIEncoderV2;

		library ExternalLib {
			struct SharedStruct{
				uint256 f1;
				bytes32 f2;
			}
		}

		contract ContractOne {
			function foo(ExternalLib.SharedStruct memory s) pure public {
				// Do stuff
			}
		}

		contract ContractTwo {
			function bar(ExternalLib.SharedStruct memory s) pure public {
				// Do stuff
			}
		}
        `,
		[]string{
			`60806040523480156100115760006000fd5b50610017565b6101b5806100266000396000f3fe60806040523480156100115760006000fd5b50600436106100305760003560e01c80639d8a8ba81461003657610030565b60006000fd5b610050600480360361004b91908101906100d1565b610052565b005b5b5056610171565b6000813590506100698161013d565b92915050565b6000604082840312156100825760006000fd5b61008c60406100fb565b9050600061009c848285016100bc565b60008301525060206100b08482850161005a565b60208301525092915050565b6000813590506100cb81610157565b92915050565b6000604082840312156100e45760006000fd5b60006100f28482850161006f565b91505092915050565b6000604051905081810181811067ffffffffffffffff8211171561011f5760006000fd5b8060405250919050565b6000819050919050565b6000819050919050565b61014681610129565b811415156101545760006000fd5b50565b61016081610133565b8114151561016e5760006000fd5b50565bfea365627a7a72315820749274eb7f6c01010d5322af4e1668b0a154409eb7968bd6cae5524c7ed669bb6c6578706572696d656e74616cf564736f6c634300050c0040`,
			`60806040523480156100115760006000fd5b50610017565b6101b5806100266000396000f3fe60806040523480156100115760006000fd5b50600436106100305760003560e01c8063db8ba08c1461003657610030565b60006000fd5b610050600480360361004b91908101906100d1565b610052565b005b5b5056610171565b6000813590506100698161013d565b92915050565b6000604082840312156100825760006000fd5b61008c60406100fb565b9050600061009c848285016100bc565b60008301525060206100b08482850161005a565b60208301525092915050565b6000813590506100cb81610157565b92915050565b6000604082840312156100e45760006000fd5b60006100f28482850161006f565b91505092915050565b6000604051905081810181811067ffffffffffffffff8211171561011f5760006000fd5b8060405250919050565b6000819050919050565b6000819050919050565b61014681610129565b811415156101545760006000fd5b50565b61016081610133565b8114151561016e5760006000fd5b50565bfea365627a7a723158209bc28ee7ea97c131a13330d77ec73b4493b5c59c648352da81dd288b021192596c6578706572696d656e74616cf564736f6c634300050c0040`,
			`606c6026600b82828239805160001a6073141515601857fe5b30600052607381538281f350fe73000000000000000000000000000000000000000030146080604052600436106023575b60006000fdfea365627a7a72315820518f0110144f5b3de95697d05e456a064656890d08e6f9cff47f3be710cc46a36c6578706572696d656e74616cf564736f6c634300050c0040`,
		},
		[]string{
			`[{"constant":true,"inputs":[{"components":[{"internalType":"uint256","name":"f1","type":"uint256"},{"internalType":"bytes32","name":"f2","type":"bytes32"}],"internalType":"struct ExternalLib.SharedStruct","name":"s","type":"tuple"}],"name":"foo","outputs":[],"payable":false,"stateMutability":"pure","type":"function"}]`,
			`[{"constant":true,"inputs":[{"components":[{"internalType":"uint256","name":"f1","type":"uint256"},{"internalType":"bytes32","name":"f2","type":"bytes32"}],"internalType":"struct ExternalLib.SharedStruct","name":"s","type":"tuple"}],"name":"bar","outputs":[],"payable":false,"stateMutability":"pure","type":"function"}]`,
			`[]`,
		},
		`
		"math/big"

		"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
		"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
		"github.com/ubiq/go-ubiq/v6/crypto"
		"github.com/ubiq/go-ubiq/v6/core"
        `,
		`
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)

		// Deploy registrar contract
		sim := backends.NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: big.NewInt(200000000000000000)}}, 10000000)
		defer sim.Close()

		transactOpts, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		_, _, c1, err := DeployContractOne(transactOpts, sim)
		if err != nil {
			t.Fatal("Failed to deploy contract")
		}
		sim.Commit()
		err = c1.Foo(nil, ExternalLibSharedStruct{
			F1: big.NewInt(100),
			F2: [32]byte{0x01, 0x02, 0x03},
		})
		if err != nil {
			t.Fatal("Failed to invoke function")
		}
		_, _, c2, err := DeployContractTwo(transactOpts, sim)
		if err != nil {
			t.Fatal("Failed to deploy contract")
		}
		sim.Commit()
		err = c2.Bar(nil, ExternalLibSharedStruct{
			F1: big.NewInt(100),
			F2: [32]byte{0x01, 0x02, 0x03},
		})
		if err != nil {
			t.Fatal("Failed to invoke function")
		}
        `,
		nil,
		nil,
		nil,
		[]string{"ContractOne", "ContractTwo", "ExternalLib"},
	},
	// Test the existence of the free retrieval calls
	{
		`PureAndView`,
		`pragma solidity >=0.6.0;
		contract PureAndView {
			function PureFunc() public pure returns (uint) {
				return 42;
			}
			function ViewFunc() public view returns (uint) {
				return block.number;
			}
		}
		`,
		[]string{`608060405234801561001057600080fd5b5060b68061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806376b5686a146037578063bb38c66c146053575b600080fd5b603d606f565b6040518082815260200191505060405180910390f35b60596077565b6040518082815260200191505060405180910390f35b600043905090565b6000602a90509056fea2646970667358221220d158c2ab7fdfce366a7998ec79ab84edd43b9815630bbaede2c760ea77f29f7f64736f6c63430006000033`},
		[]string{`[{"inputs": [],"name": "PureFunc","outputs": [{"internalType": "uint256","name": "","type": "uint256"}],"stateMutability": "pure","type": "function"},{"inputs": [],"name": "ViewFunc","outputs": [{"internalType": "uint256","name": "","type": "uint256"}],"stateMutability": "view","type": "function"}]`},
		`
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
		`,
		`
			// Generate a new random account and a funded simulator
			key, _ := crypto.GenerateKey()
			auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))

			sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(200000000000000000)}}, 10000000)
			defer sim.Close()

			// Deploy a tester contract and execute a structured call on it
			_, _, pav, err := DeployPureAndView(auth, sim)
			if err != nil {
				t.Fatalf("Failed to deploy PureAndView contract: %v", err)
			}
			sim.Commit()

			// This test the existence of the free retreiver call for view and pure functions
			if num, err := pav.PureFunc(nil); err != nil {
				t.Fatalf("Failed to call anonymous field retriever: %v", err)
			} else if num.Cmp(big.NewInt(42)) != 0 {
				t.Fatalf("Retrieved value mismatch: have %v, want %v", num, 42)
			}
			if num, err := pav.ViewFunc(nil); err != nil {
				t.Fatalf("Failed to call anonymous field retriever: %v", err)
			} else if num.Cmp(big.NewInt(1)) != 0 {
				t.Fatalf("Retrieved value mismatch: have %v, want %v", num, 1)
			}
		`,
		nil,
		nil,
		nil,
		nil,
	},
	// Test fallback separation introduced in v0.6.0
	{
		`NewFallbacks`,
		`
		pragma solidity >=0.6.0 <0.7.0;

		contract NewFallbacks {
			event Fallback(bytes data);
			fallback() external {
				emit Fallback(msg.data);
			}

			event Received(address addr, uint value);
			receive() external payable {
				emit Received(msg.sender, msg.value);
			}
		}
	   `,
		[]string{"6080604052348015600f57600080fd5b506101078061001f6000396000f3fe608060405236605f577f88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f885258743334604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a1005b348015606a57600080fd5b507f9043988963722edecc2099c75b0af0ff76af14ffca42ed6bce059a20a2a9f98660003660405180806020018281038252848482818152602001925080828437600081840152601f19601f820116905080830192505050935050505060405180910390a100fea26469706673582212201f994dcfbc53bf610b19176f9a361eafa77b447fd9c796fa2c615dfd0aaf3b8b64736f6c634300060c0033"},
		[]string{`[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"}],"name":"Fallback","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"addr","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Received","type":"event"},{"stateMutability":"nonpayable","type":"fallback"},{"stateMutability":"payable","type":"receive"}]`},
		`
			"bytes"
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
	   `,
		`
			key, _ := crypto.GenerateKey()
			addr := crypto.PubkeyToAddress(key.PublicKey)
	
			sim := backends.NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: big.NewInt(200000000000000000)}}, 1000000)
			defer sim.Close()
	
			opts, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
			_, _, c, err := DeployNewFallbacks(opts, sim)
			if err != nil {
				t.Fatalf("Failed to deploy contract: %v", err)
			}
			sim.Commit()

			// Test receive function
			opts.Value = big.NewInt(100)
			c.Receive(opts)
			sim.Commit()

			var gotEvent bool
			iter, _ := c.FilterReceived(nil)
			defer iter.Close()
			for iter.Next() {
				if iter.Event.Addr != addr {
					t.Fatal("Msg.sender mismatch")
				}
				if iter.Event.Value.Uint64() != 100 {
					t.Fatal("Msg.value mismatch")
				}
				gotEvent = true
				break
			}
			if !gotEvent {
				t.Fatal("Expect to receive event emitted by receive")
			}

			// Test fallback function
			gotEvent = false
			opts.Value = nil
			calldata := []byte{0x01, 0x02, 0x03}
			c.Fallback(opts, calldata)
			sim.Commit()

			iter2, _ := c.FilterFallback(nil)
			defer iter2.Close()
			for iter2.Next() {
				if !bytes.Equal(iter2.Event.Data, calldata) {
					t.Fatal("calldata mismatch")
				}
				gotEvent = true
				break
			}
			if !gotEvent {
				t.Fatal("Expect to receive event emitted by fallback")
			}
	   `,
		nil,
		nil,
		nil,
		nil,
	},
	// Test resolving single struct argument
	{
		`NewSingleStructArgument`,
		`
		 pragma solidity ^0.8.0;

		 contract NewSingleStructArgument {
			 struct MyStruct{
				 uint256 a;
				 uint256 b;
			 }
			 event StructEvent(MyStruct s);
			 function TestEvent() public {
				 emit StructEvent(MyStruct({a: 1, b: 2}));
			 }
		 }
	   `,
		[]string{"608060405234801561001057600080fd5b50610113806100206000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806324ec1d3f14602d575b600080fd5b60336035565b005b7fb4b2ff75e30cb4317eaae16dd8a187dd89978df17565104caa6c2797caae27d460405180604001604052806001815260200160028152506040516078919060ba565b60405180910390a1565b6040820160008201516096600085018260ad565b50602082015160a7602085018260ad565b50505050565b60b48160d3565b82525050565b600060408201905060cd60008301846082565b92915050565b600081905091905056fea26469706673582212208823628796125bf9941ce4eda18da1be3cf2931b231708ab848e1bd7151c0c9a64736f6c63430008070033"},
		[]string{`[{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"a","type":"uint256"},{"internalType":"uint256","name":"b","type":"uint256"}],"indexed":false,"internalType":"struct Test.MyStruct","name":"s","type":"tuple"}],"name":"StructEvent","type":"event"},{"inputs":[],"name":"TestEvent","outputs":[],"stateMutability":"nonpayable","type":"function"}]`},
		`
			"math/big"

			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
			"github.com/ubiq/go-ubiq/v6/eth/ethconfig"
	   `,
		`
			var (
				key, _  = crypto.GenerateKey()
				user, _ = bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
				sim     = backends.NewSimulatedBackend(core.GenesisAlloc{user.From: {Balance: big.NewInt(1000000000000000000)}}, ethconfig.Defaults.Miner.GasCeil)
			)
			defer sim.Close()

			_, _, d, err := DeployNewSingleStructArgument(user, sim)
			if err != nil {
				t.Fatalf("Failed to deploy contract %v", err)
			}
			sim.Commit()

			_, err = d.TestEvent(user)
			if err != nil {
				t.Fatalf("Failed to call contract %v", err)
			}
			sim.Commit()

			it, err := d.FilterStructEvent(nil)
			if err != nil {
				t.Fatalf("Failed to filter contract event %v", err)
			}
			var count int
			for it.Next() {
				if it.Event.S.A.Cmp(big.NewInt(1)) != 0 {
					t.Fatal("Unexpected contract event")
				}
				if it.Event.S.B.Cmp(big.NewInt(2)) != 0 {
					t.Fatal("Unexpected contract event")
				}
				count += 1
			}
			if count != 1 {
				t.Fatal("Unexpected contract event number")
			}
			`,
		nil,
		nil,
		nil,
		nil,
	},
	// Test errors introduced in v0.8.4
	{
		`NewErrors`,
		`
		pragma solidity >0.8.4;
	
		contract NewErrors {
			error MyError(uint256);
			error MyError1(uint256);
			error MyError2(uint256, uint256);
			error MyError3(uint256 a, uint256 b, uint256 c);
			function Error() public pure {
				revert MyError3(1,2,3);
			}
		}
	   `,
		[]string{"0x6080604052348015600f57600080fd5b5060998061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063726c638214602d575b600080fd5b60336035565b005b60405163024876cd60e61b815260016004820152600260248201526003604482015260640160405180910390fdfea264697066735822122093f786a1bc60216540cd999fbb4a6109e0fef20abcff6e9107fb2817ca968f3c64736f6c63430008070033"},
		[]string{`[{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"MyError","type":"error"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"MyError1","type":"error"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"MyError2","type":"error"},{"inputs":[{"internalType":"uint256","name":"a","type":"uint256"},{"internalType":"uint256","name":"b","type":"uint256"},{"internalType":"uint256","name":"c","type":"uint256"}],"name":"MyError3","type":"error"},{"inputs":[],"name":"Error","outputs":[],"stateMutability":"pure","type":"function"}]`},
		`
			"math/big"
	
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind"
			"github.com/ubiq/go-ubiq/v6/accounts/abi/bind/backends"
			"github.com/ubiq/go-ubiq/v6/core"
			"github.com/ubiq/go-ubiq/v6/crypto"
			"github.com/ubiq/go-ubiq/v6/eth/ethconfig"
	   `,
		`
			var (
				key, _  = crypto.GenerateKey()
				user, _ = bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
				sim     = backends.NewSimulatedBackend(core.GenesisAlloc{user.From: {Balance: big.NewInt(1000000000000000000)}}, ethconfig.Defaults.Miner.GasCeil)
			)
			defer sim.Close()
	
			_, tx, contract, err := DeployNewErrors(user, sim)
			if err != nil {
				t.Fatal(err)
			}
			sim.Commit()
			_, err = bind.WaitDeployed(nil, sim, tx)
			if err != nil {
				t.Error(err)
			}
			if err := contract.Error(new(bind.CallOpts)); err == nil {
				t.Fatalf("expected contract to throw error")
			}
			// TODO (MariusVanDerWijden unpack error using abigen
			// once that is implemented
	   `,
		nil,
		nil,
		nil,
		nil,
	},
}

// Tests that packages generated by the binder can be successfully compiled and
// the requested tester run against it.
func TestGolangBindings(t *testing.T) {
	// Skip the test if no Go command can be found
	gocmd := runtime.GOROOT() + "/bin/go"
	if !common.FileExist(gocmd) {
		t.Skip("go sdk not found for testing")
	}
	// Create a temporary workspace for the test suite
	ws, err := ioutil.TempDir("", "binding-test")
	if err != nil {
		t.Fatalf("failed to create temporary workspace: %v", err)
	}
	//defer os.RemoveAll(ws)

	pkg := filepath.Join(ws, "bindtest")
	if err = os.MkdirAll(pkg, 0700); err != nil {
		t.Fatalf("failed to create package: %v", err)
	}
	// Generate the test suite for all the contracts
	for i, tt := range bindTests {
		var types []string
		if tt.types != nil {
			types = tt.types
		} else {
			types = []string{tt.name}
		}
		// Generate the binding and create a Go source file in the workspace
		bind, err := Bind(types, tt.abi, tt.bytecode, tt.fsigs, "bindtest", LangGo, tt.libs, tt.aliases)
		if err != nil {
			t.Fatalf("test %d: failed to generate binding: %v", i, err)
		}
		if err = ioutil.WriteFile(filepath.Join(pkg, strings.ToLower(tt.name)+".go"), []byte(bind), 0600); err != nil {
			t.Fatalf("test %d: failed to write binding: %v", i, err)
		}
		// Generate the test file with the injected test code
		code := fmt.Sprintf(`
			package bindtest

			import (
				"testing"
				%s
			)

			func Test%s(t *testing.T) {
				%s
			}
		`, tt.imports, tt.name, tt.tester)
		if err := ioutil.WriteFile(filepath.Join(pkg, strings.ToLower(tt.name)+"_test.go"), []byte(code), 0600); err != nil {
			t.Fatalf("test %d: failed to write tests: %v", i, err)
		}
	}
	// Convert the package to go modules and use the current source for go-ethereum
	moder := exec.Command(gocmd, "mod", "init", "bindtest")
	moder.Dir = pkg
	if out, err := moder.CombinedOutput(); err != nil {
		t.Fatalf("failed to convert binding test to modules: %v\n%s", err, out)
	}
	pwd, _ := os.Getwd()
	replacer := exec.Command(gocmd, "mod", "edit", "-x", "-require", "github.com/ubiq/go-ubiq@v0.0.0", "-replace", "github.com/ubiq/go-ubiq="+filepath.Join(pwd, "..", "..", "..")) // Repo root
	replacer.Dir = pkg
	if out, err := replacer.CombinedOutput(); err != nil {
		t.Fatalf("failed to replace binding test dependency to current source tree: %v\n%s", err, out)
	}
	tidier := exec.Command(gocmd, "mod", "tidy")
	tidier.Dir = pkg
	if out, err := tidier.CombinedOutput(); err != nil {
		t.Fatalf("failed to tidy Go module file: %v\n%s", err, out)
	}
	// Test the entire package and report any failures
	cmd := exec.Command(gocmd, "test", "-v", "-count", "1")
	cmd.Dir = pkg
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to run binding test: %v\n%s", err, out)
	}
}

// Tests that java binding generated by the binder is exactly matched.
func TestJavaBindings(t *testing.T) {
	var cases = []struct {
		name     string
		contract string
		abi      string
		bytecode string
		expected string
	}{
		{
			"test",
			`
			pragma experimental ABIEncoderV2;
			pragma solidity ^0.5.2;

			contract test {
				function setAddress(address a) public returns(address){}
				function setAddressList(address[] memory a_l) public returns(address[] memory){}
				function setAddressArray(address[2] memory a_a) public returns(address[2] memory){}

				function setUint8(uint8 u8) public returns(uint8){}
				function setUint16(uint16 u16) public returns(uint16){}
				function setUint32(uint32 u32) public returns(uint32){}
				function setUint64(uint64 u64) public returns(uint64){}
				function setUint256(uint256 u256) public returns(uint256){}
				function setUint256List(uint256[] memory u256_l) public returns(uint256[] memory){}
				function setUint256Array(uint256[2] memory u256_a) public returns(uint256[2] memory){}

				function setInt8(int8 i8) public returns(int8){}
				function setInt16(int16 i16) public returns(int16){}
				function setInt32(int32 i32) public returns(int32){}
				function setInt64(int64 i64) public returns(int64){}
				function setInt256(int256 i256) public returns(int256){}
				function setInt256List(int256[] memory i256_l) public returns(int256[] memory){}
				function setInt256Array(int256[2] memory i256_a) public returns(int256[2] memory){}

				function setBytes1(bytes1 b1) public returns(bytes1) {}
				function setBytes32(bytes32 b32) public returns(bytes32) {}
				function setBytes(bytes memory bs) public returns(bytes memory) {}
				function setBytesList(bytes[] memory bs_l) public returns(bytes[] memory) {}
				function setBytesArray(bytes[2] memory bs_a) public returns(bytes[2] memory) {}

				function setString(string memory s) public returns(string memory) {}
				function setStringList(string[] memory s_l) public returns(string[] memory) {}
				function setStringArray(string[2] memory s_a) public returns(string[2] memory) {}

				function setBool(bool b) public returns(bool) {}
				function setBoolList(bool[] memory b_l) public returns(bool[] memory) {}
				function setBoolArray(bool[2] memory b_a) public returns(bool[2] memory) {}
			}`,
			`[{"constant":false,"inputs":[{"name":"u16","type":"uint16"}],"name":"setUint16","outputs":[{"name":"","type":"uint16"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"b_a","type":"bool[2]"}],"name":"setBoolArray","outputs":[{"name":"","type":"bool[2]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"a_a","type":"address[2]"}],"name":"setAddressArray","outputs":[{"name":"","type":"address[2]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"bs_l","type":"bytes[]"}],"name":"setBytesList","outputs":[{"name":"","type":"bytes[]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"u8","type":"uint8"}],"name":"setUint8","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"u32","type":"uint32"}],"name":"setUint32","outputs":[{"name":"","type":"uint32"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"b","type":"bool"}],"name":"setBool","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"i256_l","type":"int256[]"}],"name":"setInt256List","outputs":[{"name":"","type":"int256[]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"u256_a","type":"uint256[2]"}],"name":"setUint256Array","outputs":[{"name":"","type":"uint256[2]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"b_l","type":"bool[]"}],"name":"setBoolList","outputs":[{"name":"","type":"bool[]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"bs_a","type":"bytes[2]"}],"name":"setBytesArray","outputs":[{"name":"","type":"bytes[2]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"a_l","type":"address[]"}],"name":"setAddressList","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"i256_a","type":"int256[2]"}],"name":"setInt256Array","outputs":[{"name":"","type":"int256[2]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"s_a","type":"string[2]"}],"name":"setStringArray","outputs":[{"name":"","type":"string[2]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"s","type":"string"}],"name":"setString","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"u64","type":"uint64"}],"name":"setUint64","outputs":[{"name":"","type":"uint64"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"i16","type":"int16"}],"name":"setInt16","outputs":[{"name":"","type":"int16"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"i8","type":"int8"}],"name":"setInt8","outputs":[{"name":"","type":"int8"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"u256_l","type":"uint256[]"}],"name":"setUint256List","outputs":[{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"i256","type":"int256"}],"name":"setInt256","outputs":[{"name":"","type":"int256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"i32","type":"int32"}],"name":"setInt32","outputs":[{"name":"","type":"int32"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"b32","type":"bytes32"}],"name":"setBytes32","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"s_l","type":"string[]"}],"name":"setStringList","outputs":[{"name":"","type":"string[]"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"u256","type":"uint256"}],"name":"setUint256","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"bs","type":"bytes"}],"name":"setBytes","outputs":[{"name":"","type":"bytes"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"a","type":"address"}],"name":"setAddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"i64","type":"int64"}],"name":"setInt64","outputs":[{"name":"","type":"int64"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"b1","type":"bytes1"}],"name":"setBytes1","outputs":[{"name":"","type":"bytes1"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`,
			`608060405234801561001057600080fd5b5061265a806100206000396000f3fe608060405234801561001057600080fd5b50600436106101e1576000357c0100000000000000000000000000000000000000000000000000000000900480637fcaf66611610116578063c2b12a73116100b4578063da359dc81161008e578063da359dc814610666578063e30081a014610696578063e673eb32146106c6578063fba1a1c3146106f6576101e1565b8063c2b12a73146105d6578063c577796114610606578063d2282dc514610636576101e1565b80639a19a953116100f05780639a19a95314610516578063a0709e1914610546578063a53b1c1e14610576578063b7d5df31146105a6576101e1565b80637fcaf66614610486578063822cba69146104b657806386114cea146104e6576101e1565b806322722302116101835780635119655d1161015d5780635119655d146103c65780635be6b37e146103f65780636aa482fc146104265780637173b69514610456576101e1565b806322722302146103365780632766a755146103665780634d5ee6da14610396576101e1565b806316c105e2116101bf57806316c105e2146102765780631774e646146102a65780631c9352e2146102d65780631e26fd3314610306576101e1565b80630477988a146101e6578063118a971814610216578063151f547114610246575b600080fd5b61020060048036036101fb9190810190611599565b610726565b60405161020d9190611f01565b60405180910390f35b610230600480360361022b919081019061118d565b61072d565b60405161023d9190611ca6565b60405180910390f35b610260600480360361025b9190810190611123565b61073a565b60405161026d9190611c69565b60405180910390f35b610290600480360361028b9190810190611238565b610747565b60405161029d9190611d05565b60405180910390f35b6102c060048036036102bb919081019061163d565b61074e565b6040516102cd9190611f6d565b60405180910390f35b6102f060048036036102eb91908101906115eb565b610755565b6040516102fd9190611f37565b60405180910390f35b610320600480360361031b91908101906113cf565b61075c565b60405161032d9190611de5565b60405180910390f35b610350600480360361034b91908101906112a2565b610763565b60405161035d9190611d42565b60405180910390f35b610380600480360361037b9190810190611365565b61076a565b60405161038d9190611da8565b60405180910390f35b6103b060048036036103ab91908101906111b6565b610777565b6040516103bd9190611cc1565b60405180910390f35b6103e060048036036103db91908101906111f7565b61077e565b6040516103ed9190611ce3565b60405180910390f35b610410600480360361040b919081019061114c565b61078b565b60405161041d9190611c84565b60405180910390f35b610440600480360361043b9190810190611279565b610792565b60405161044d9190611d27565b60405180910390f35b610470600480360361046b91908101906112e3565b61079f565b60405161047d9190611d64565b60405180910390f35b6104a0600480360361049b9190810190611558565b6107ac565b6040516104ad9190611edf565b60405180910390f35b6104d060048036036104cb9190810190611614565b6107b3565b6040516104dd9190611f52565b60405180910390f35b61050060048036036104fb919081019061148b565b6107ba565b60405161050d9190611e58565b60405180910390f35b610530600480360361052b919081019061152f565b6107c1565b60405161053d9190611ec4565b60405180910390f35b610560600480360361055b919081019061138e565b6107c8565b60405161056d9190611dc3565b60405180910390f35b610590600480360361058b91908101906114b4565b6107cf565b60405161059d9190611e73565b60405180910390f35b6105c060048036036105bb91908101906114dd565b6107d6565b6040516105cd9190611e8e565b60405180910390f35b6105f060048036036105eb9190810190611421565b6107dd565b6040516105fd9190611e1b565b60405180910390f35b610620600480360361061b9190810190611324565b6107e4565b60405161062d9190611d86565b60405180910390f35b610650600480360361064b91908101906115c2565b6107eb565b60405161065d9190611f1c565b60405180910390f35b610680600480360361067b919081019061144a565b6107f2565b60405161068d9190611e36565b60405180910390f35b6106b060048036036106ab91908101906110fa565b6107f9565b6040516106bd9190611c4e565b60405180910390f35b6106e060048036036106db9190810190611506565b610800565b6040516106ed9190611ea9565b60405180910390f35b610710600480360361070b91908101906113f8565b610807565b60405161071d9190611e00565b60405180910390f35b6000919050565b61073561080e565b919050565b610742610830565b919050565b6060919050565b6000919050565b6000919050565b6000919050565b6060919050565b610772610852565b919050565b6060919050565b610786610874565b919050565b6060919050565b61079a61089b565b919050565b6107a76108bd565b919050565b6060919050565b6000919050565b6000919050565b6000919050565b6060919050565b6000919050565b6000919050565b6000919050565b6060919050565b6000919050565b6060919050565b6000919050565b6000919050565b6000919050565b6040805190810160405280600290602082028038833980820191505090505090565b6040805190810160405280600290602082028038833980820191505090505090565b6040805190810160405280600290602082028038833980820191505090505090565b60408051908101604052806002905b60608152602001906001900390816108835790505090565b6040805190810160405280600290602082028038833980820191505090505090565b60408051908101604052806002905b60608152602001906001900390816108cc5790505090565b60006108f082356124f2565b905092915050565b600082601f830112151561090b57600080fd5b600261091e61091982611fb5565b611f88565b9150818385602084028201111561093457600080fd5b60005b83811015610964578161094a88826108e4565b845260208401935060208301925050600181019050610937565b5050505092915050565b600082601f830112151561098157600080fd5b813561099461098f82611fd7565b611f88565b915081818352602084019350602081019050838560208402820111156109b957600080fd5b60005b838110156109e957816109cf88826108e4565b8452602084019350602083019250506001810190506109bc565b5050505092915050565b600082601f8301121515610a0657600080fd5b6002610a19610a1482611fff565b611f88565b91508183856020840282011115610a2f57600080fd5b60005b83811015610a5f5781610a458882610e9e565b845260208401935060208301925050600181019050610a32565b5050505092915050565b600082601f8301121515610a7c57600080fd5b8135610a8f610a8a82612021565b611f88565b91508181835260208401935060208101905083856020840282011115610ab457600080fd5b60005b83811015610ae45781610aca8882610e9e565b845260208401935060208301925050600181019050610ab7565b5050505092915050565b600082601f8301121515610b0157600080fd5b6002610b14610b0f82612049565b611f88565b9150818360005b83811015610b4b5781358601610b318882610eda565b845260208401935060208301925050600181019050610b1b565b5050505092915050565b600082601f8301121515610b6857600080fd5b8135610b7b610b768261206b565b611f88565b9150818183526020840193506020810190508360005b83811015610bc15781358601610ba78882610eda565b845260208401935060208301925050600181019050610b91565b5050505092915050565b600082601f8301121515610bde57600080fd5b6002610bf1610bec82612093565b611f88565b91508183856020840282011115610c0757600080fd5b60005b83811015610c375781610c1d8882610f9a565b845260208401935060208301925050600181019050610c0a565b5050505092915050565b600082601f8301121515610c5457600080fd5b8135610c67610c62826120b5565b611f88565b91508181835260208401935060208101905083856020840282011115610c8c57600080fd5b60005b83811015610cbc5781610ca28882610f9a565b845260208401935060208301925050600181019050610c8f565b5050505092915050565b600082601f8301121515610cd957600080fd5b6002610cec610ce7826120dd565b611f88565b9150818360005b83811015610d235781358601610d098882610fea565b845260208401935060208301925050600181019050610cf3565b5050505092915050565b600082601f8301121515610d4057600080fd5b8135610d53610d4e826120ff565b611f88565b9150818183526020840193506020810190508360005b83811015610d995781358601610d7f8882610fea565b845260208401935060208301925050600181019050610d69565b5050505092915050565b600082601f8301121515610db657600080fd5b6002610dc9610dc482612127565b611f88565b91508183856020840282011115610ddf57600080fd5b60005b83811015610e0f5781610df588826110aa565b845260208401935060208301925050600181019050610de2565b5050505092915050565b600082601f8301121515610e2c57600080fd5b8135610e3f610e3a82612149565b611f88565b91508181835260208401935060208101905083856020840282011115610e6457600080fd5b60005b83811015610e945781610e7a88826110aa565b845260208401935060208301925050600181019050610e67565b5050505092915050565b6000610eaa8235612504565b905092915050565b6000610ebe8235612510565b905092915050565b6000610ed2823561253c565b905092915050565b600082601f8301121515610eed57600080fd5b8135610f00610efb82612171565b611f88565b91508082526020830160208301858383011115610f1c57600080fd5b610f278382846125cd565b50505092915050565b600082601f8301121515610f4357600080fd5b8135610f56610f518261219d565b611f88565b91508082526020830160208301858383011115610f7257600080fd5b610f7d8382846125cd565b50505092915050565b6000610f928235612546565b905092915050565b6000610fa68235612553565b905092915050565b6000610fba823561255d565b905092915050565b6000610fce823561256a565b905092915050565b6000610fe28235612577565b905092915050565b600082601f8301121515610ffd57600080fd5b813561101061100b826121c9565b611f88565b9150808252602083016020830185838301111561102c57600080fd5b6110378382846125cd565b50505092915050565b600082601f830112151561105357600080fd5b8135611066611061826121f5565b611f88565b9150808252602083016020830185838301111561108257600080fd5b61108d8382846125cd565b50505092915050565b60006110a28235612584565b905092915050565b60006110b68235612592565b905092915050565b60006110ca823561259c565b905092915050565b60006110de82356125ac565b905092915050565b60006110f282356125c0565b905092915050565b60006020828403121561110c57600080fd5b600061111a848285016108e4565b91505092915050565b60006040828403121561113557600080fd5b6000611143848285016108f8565b91505092915050565b60006020828403121561115e57600080fd5b600082013567ffffffffffffffff81111561117857600080fd5b6111848482850161096e565b91505092915050565b60006040828403121561119f57600080fd5b60006111ad848285016109f3565b91505092915050565b6000602082840312156111c857600080fd5b600082013567ffffffffffffffff8111156111e257600080fd5b6111ee84828501610a69565b91505092915050565b60006020828403121561120957600080fd5b600082013567ffffffffffffffff81111561122357600080fd5b61122f84828501610aee565b91505092915050565b60006020828403121561124a57600080fd5b600082013567ffffffffffffffff81111561126457600080fd5b61127084828501610b55565b91505092915050565b60006040828403121561128b57600080fd5b600061129984828501610bcb565b91505092915050565b6000602082840312156112b457600080fd5b600082013567ffffffffffffffff8111156112ce57600080fd5b6112da84828501610c41565b91505092915050565b6000602082840312156112f557600080fd5b600082013567ffffffffffffffff81111561130f57600080fd5b61131b84828501610cc6565b91505092915050565b60006020828403121561133657600080fd5b600082013567ffffffffffffffff81111561135057600080fd5b61135c84828501610d2d565b91505092915050565b60006040828403121561137757600080fd5b600061138584828501610da3565b91505092915050565b6000602082840312156113a057600080fd5b600082013567ffffffffffffffff8111156113ba57600080fd5b6113c684828501610e19565b91505092915050565b6000602082840312156113e157600080fd5b60006113ef84828501610e9e565b91505092915050565b60006020828403121561140a57600080fd5b600061141884828501610eb2565b91505092915050565b60006020828403121561143357600080fd5b600061144184828501610ec6565b91505092915050565b60006020828403121561145c57600080fd5b600082013567ffffffffffffffff81111561147657600080fd5b61148284828501610f30565b91505092915050565b60006020828403121561149d57600080fd5b60006114ab84828501610f86565b91505092915050565b6000602082840312156114c657600080fd5b60006114d484828501610f9a565b91505092915050565b6000602082840312156114ef57600080fd5b60006114fd84828501610fae565b91505092915050565b60006020828403121561151857600080fd5b600061152684828501610fc2565b91505092915050565b60006020828403121561154157600080fd5b600061154f84828501610fd6565b91505092915050565b60006020828403121561156a57600080fd5b600082013567ffffffffffffffff81111561158457600080fd5b61159084828501611040565b91505092915050565b6000602082840312156115ab57600080fd5b60006115b984828501611096565b91505092915050565b6000602082840312156115d457600080fd5b60006115e2848285016110aa565b91505092915050565b6000602082840312156115fd57600080fd5b600061160b848285016110be565b91505092915050565b60006020828403121561162657600080fd5b6000611634848285016110d2565b91505092915050565b60006020828403121561164f57600080fd5b600061165d848285016110e6565b91505092915050565b61166f816123f7565b82525050565b61167e816122ab565b61168782612221565b60005b828110156116b95761169d858351611666565b6116a68261235b565b915060208501945060018101905061168a565b5050505050565b60006116cb826122b6565b8084526020840193506116dd8361222b565b60005b8281101561170f576116f3868351611666565b6116fc82612368565b91506020860195506001810190506116e0565b50849250505092915050565b611724816122c1565b61172d82612238565b60005b8281101561175f57611743858351611ab3565b61174c82612375565b9150602085019450600181019050611730565b5050505050565b6000611771826122cc565b80845260208401935061178383612242565b60005b828110156117b557611799868351611ab3565b6117a282612382565b9150602086019550600181019050611786565b50849250505092915050565b60006117cc826122d7565b836020820285016117dc8561224f565b60005b848110156118155783830388526117f7838351611b16565b92506118028261238f565b91506020880197506001810190506117df565b508196508694505050505092915050565b6000611831826122e2565b8084526020840193508360208202850161184a85612259565b60005b84811015611883578383038852611865838351611b16565b92506118708261239c565b915060208801975060018101905061184d565b508196508694505050505092915050565b61189d816122ed565b6118a682612266565b60005b828110156118d8576118bc858351611b5b565b6118c5826123a9565b91506020850194506001810190506118a9565b5050505050565b60006118ea826122f8565b8084526020840193506118fc83612270565b60005b8281101561192e57611912868351611b5b565b61191b826123b6565b91506020860195506001810190506118ff565b50849250505092915050565b600061194582612303565b836020820285016119558561227d565b60005b8481101561198e578383038852611970838351611bcd565b925061197b826123c3565b9150602088019750600181019050611958565b508196508694505050505092915050565b60006119aa8261230e565b808452602084019350836020820285016119c385612287565b60005b848110156119fc5783830388526119de838351611bcd565b92506119e9826123d0565b91506020880197506001810190506119c6565b508196508694505050505092915050565b611a1681612319565b611a1f82612294565b60005b82811015611a5157611a35858351611c12565b611a3e826123dd565b9150602085019450600181019050611a22565b5050505050565b6000611a6382612324565b808452602084019350611a758361229e565b60005b82811015611aa757611a8b868351611c12565b611a94826123ea565b9150602086019550600181019050611a78565b50849250505092915050565b611abc81612409565b82525050565b611acb81612415565b82525050565b611ada81612441565b82525050565b6000611aeb8261233a565b808452611aff8160208601602086016125dc565b611b088161260f565b602085010191505092915050565b6000611b218261232f565b808452611b358160208601602086016125dc565b611b3e8161260f565b602085010191505092915050565b611b558161244b565b82525050565b611b6481612458565b82525050565b611b7381612462565b82525050565b611b828161246f565b82525050565b611b918161247c565b82525050565b6000611ba282612350565b808452611bb68160208601602086016125dc565b611bbf8161260f565b602085010191505092915050565b6000611bd882612345565b808452611bec8160208601602086016125dc565b611bf58161260f565b602085010191505092915050565b611c0c81612489565b82525050565b611c1b816124b7565b82525050565b611c2a816124c1565b82525050565b611c39816124d1565b82525050565b611c48816124e5565b82525050565b6000602082019050611c636000830184611666565b92915050565b6000604082019050611c7e6000830184611675565b92915050565b60006020820190508181036000830152611c9e81846116c0565b905092915050565b6000604082019050611cbb600083018461171b565b92915050565b60006020820190508181036000830152611cdb8184611766565b905092915050565b60006020820190508181036000830152611cfd81846117c1565b905092915050565b60006020820190508181036000830152611d1f8184611826565b905092915050565b6000604082019050611d3c6000830184611894565b92915050565b60006020820190508181036000830152611d5c81846118df565b905092915050565b60006020820190508181036000830152611d7e818461193a565b905092915050565b60006020820190508181036000830152611da0818461199f565b905092915050565b6000604082019050611dbd6000830184611a0d565b92915050565b60006020820190508181036000830152611ddd8184611a58565b905092915050565b6000602082019050611dfa6000830184611ab3565b92915050565b6000602082019050611e156000830184611ac2565b92915050565b6000602082019050611e306000830184611ad1565b92915050565b60006020820190508181036000830152611e508184611ae0565b905092915050565b6000602082019050611e6d6000830184611b4c565b92915050565b6000602082019050611e886000830184611b5b565b92915050565b6000602082019050611ea36000830184611b6a565b92915050565b6000602082019050611ebe6000830184611b79565b92915050565b6000602082019050611ed96000830184611b88565b92915050565b60006020820190508181036000830152611ef98184611b97565b905092915050565b6000602082019050611f166000830184611c03565b92915050565b6000602082019050611f316000830184611c12565b92915050565b6000602082019050611f4c6000830184611c21565b92915050565b6000602082019050611f676000830184611c30565b92915050565b6000602082019050611f826000830184611c3f565b92915050565b6000604051905081810181811067ffffffffffffffff82111715611fab57600080fd5b8060405250919050565b600067ffffffffffffffff821115611fcc57600080fd5b602082029050919050565b600067ffffffffffffffff821115611fee57600080fd5b602082029050602081019050919050565b600067ffffffffffffffff82111561201657600080fd5b602082029050919050565b600067ffffffffffffffff82111561203857600080fd5b602082029050602081019050919050565b600067ffffffffffffffff82111561206057600080fd5b602082029050919050565b600067ffffffffffffffff82111561208257600080fd5b602082029050602081019050919050565b600067ffffffffffffffff8211156120aa57600080fd5b602082029050919050565b600067ffffffffffffffff8211156120cc57600080fd5b602082029050602081019050919050565b600067ffffffffffffffff8211156120f457600080fd5b602082029050919050565b600067ffffffffffffffff82111561211657600080fd5b602082029050602081019050919050565b600067ffffffffffffffff82111561213e57600080fd5b602082029050919050565b600067ffffffffffffffff82111561216057600080fd5b602082029050602081019050919050565b600067ffffffffffffffff82111561218857600080fd5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156121b457600080fd5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156121e057600080fd5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff82111561220c57600080fd5b601f19601f8301169050602081019050919050565b6000819050919050565b6000602082019050919050565b6000819050919050565b6000602082019050919050565b6000819050919050565b6000602082019050919050565b6000819050919050565b6000602082019050919050565b6000819050919050565b6000602082019050919050565b6000819050919050565b6000602082019050919050565b600060029050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600081519050919050565b600081519050919050565b600081519050919050565b600081519050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b600061240282612497565b9050919050565b60008115159050919050565b60007fff0000000000000000000000000000000000000000000000000000000000000082169050919050565b6000819050919050565b60008160010b9050919050565b6000819050919050565b60008160030b9050919050565b60008160070b9050919050565b60008160000b9050919050565b600061ffff82169050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600063ffffffff82169050919050565b600067ffffffffffffffff82169050919050565b600060ff82169050919050565b60006124fd82612497565b9050919050565b60008115159050919050565b60007fff0000000000000000000000000000000000000000000000000000000000000082169050919050565b6000819050919050565b60008160010b9050919050565b6000819050919050565b60008160030b9050919050565b60008160070b9050919050565b60008160000b9050919050565b600061ffff82169050919050565b6000819050919050565b600063ffffffff82169050919050565b600067ffffffffffffffff82169050919050565b600060ff82169050919050565b82818337600083830152505050565b60005b838110156125fa5780820151818401526020810190506125df565b83811115612609576000848401525b50505050565b6000601f19601f830116905091905056fea265627a7a723058206fe37171cf1b10ebd291cfdca61d67e7fc3c208795e999c833c42a14d86cf00d6c6578706572696d656e74616cf50037`,
			`
// This file is an automatically generated Java binding. Do not modify as any
// change will likely be lost upon the next re-generation!

package bindtest;

import org.ubiq.gubiq.*;
import java.util.*;

public class Test {
	// ABI is the input ABI used to generate the binding from.
	public final static String ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"u16\",\"type\":\"uint16\"}],\"name\":\"setUint16\",\"outputs\":[{\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"b_a\",\"type\":\"bool[2]\"}],\"name\":\"setBoolArray\",\"outputs\":[{\"name\":\"\",\"type\":\"bool[2]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"a_a\",\"type\":\"address[2]\"}],\"name\":\"setAddressArray\",\"outputs\":[{\"name\":\"\",\"type\":\"address[2]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"bs_l\",\"type\":\"bytes[]\"}],\"name\":\"setBytesList\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"u8\",\"type\":\"uint8\"}],\"name\":\"setUint8\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"u32\",\"type\":\"uint32\"}],\"name\":\"setUint32\",\"outputs\":[{\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"b\",\"type\":\"bool\"}],\"name\":\"setBool\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"i256_l\",\"type\":\"int256[]\"}],\"name\":\"setInt256List\",\"outputs\":[{\"name\":\"\",\"type\":\"int256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"u256_a\",\"type\":\"uint256[2]\"}],\"name\":\"setUint256Array\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[2]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"b_l\",\"type\":\"bool[]\"}],\"name\":\"setBoolList\",\"outputs\":[{\"name\":\"\",\"type\":\"bool[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"bs_a\",\"type\":\"bytes[2]\"}],\"name\":\"setBytesArray\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes[2]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"a_l\",\"type\":\"address[]\"}],\"name\":\"setAddressList\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"i256_a\",\"type\":\"int256[2]\"}],\"name\":\"setInt256Array\",\"outputs\":[{\"name\":\"\",\"type\":\"int256[2]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"s_a\",\"type\":\"string[2]\"}],\"name\":\"setStringArray\",\"outputs\":[{\"name\":\"\",\"type\":\"string[2]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"s\",\"type\":\"string\"}],\"name\":\"setString\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"u64\",\"type\":\"uint64\"}],\"name\":\"setUint64\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"i16\",\"type\":\"int16\"}],\"name\":\"setInt16\",\"outputs\":[{\"name\":\"\",\"type\":\"int16\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"i8\",\"type\":\"int8\"}],\"name\":\"setInt8\",\"outputs\":[{\"name\":\"\",\"type\":\"int8\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"u256_l\",\"type\":\"uint256[]\"}],\"name\":\"setUint256List\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"i256\",\"type\":\"int256\"}],\"name\":\"setInt256\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"i32\",\"type\":\"int32\"}],\"name\":\"setInt32\",\"outputs\":[{\"name\":\"\",\"type\":\"int32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"b32\",\"type\":\"bytes32\"}],\"name\":\"setBytes32\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"s_l\",\"type\":\"string[]\"}],\"name\":\"setStringList\",\"outputs\":[{\"name\":\"\",\"type\":\"string[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"u256\",\"type\":\"uint256\"}],\"name\":\"setUint256\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"bs\",\"type\":\"bytes\"}],\"name\":\"setBytes\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"a\",\"type\":\"address\"}],\"name\":\"setAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"i64\",\"type\":\"int64\"}],\"name\":\"setInt64\",\"outputs\":[{\"name\":\"\",\"type\":\"int64\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"b1\",\"type\":\"bytes1\"}],\"name\":\"setBytes1\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes1\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]";

	// BYTECODE is the compiled bytecode used for deploying new contracts.
	public final static String BYTECODE = "0x608060405234801561001057600080fd5b5061265a806100206000396000f3fe608060405234801561001057600080fd5b50600436106101e1576000357c0100000000000000000000000000000000000000000000000000000000900480637fcaf66611610116578063c2b12a73116100b4578063da359dc81161008e578063da359dc814610666578063e30081a014610696578063e673eb32146106c6578063fba1a1c3146106f6576101e1565b8063c2b12a73146105d6578063c577796114610606578063d2282dc514610636576101e1565b80639a19a953116100f05780639a19a95314610516578063a0709e1914610546578063a53b1c1e14610576578063b7d5df31146105a6576101e1565b80637fcaf66614610486578063822cba69146104b657806386114cea146104e6576101e1565b806322722302116101835780635119655d1161015d5780635119655d146103c65780635be6b37e146103f65780636aa482fc146104265780637173b69514610456576101e1565b806322722302146103365780632766a755146103665780634d5ee6da14610396576101e1565b806316c105e2116101bf57806316c105e2146102765780631774e646146102a65780631c9352e2146102d65780631e26fd3314610306576101e1565b80630477988a146101e6578063118a971814610216578063151f547114610246575b600080fd5b61020060048036036101fb9190810190611599565b610726565b60405161020d9190611f01565b60405180910390f35b610230600480360361022b919081019061118d565b61072d565b60405161023d9190611ca6565b60405180910390f35b610260600480360361025b9190810190611123565b61073a565b60405161026d9190611c69565b60405180910390f35b610290600480360361028b9190810190611238565b610747565b60405161029d9190611d05565b60405180910390f35b6102c060048036036102bb919081019061163d565b61074e565b6040516102cd9190611f6d565b60405180910390f35b6102f060048036036102eb91908101906115eb565b610755565b6040516102fd9190611f37565b60405180910390f35b610320600480360361031b91908101906113cf565b61075c565b60405161032d9190611de5565b60405180910390f35b610350600480360361034b91908101906112a2565b610763565b60405161035d9190611d42565b60405180910390f35b610380600480360361037b9190810190611365565b61076a565b60405161038d9190611da8565b60405180910390f35b6103b060048036036103ab91908101906111b6565b610777565b6040516103bd9190611cc1565b60405180910390f35b6103e060048036036103db91908101906111f7565b61077e565b6040516103ed9190611ce3565b60405180910390f35b610410600480360361040b919081019061114c565b61078b565b60405161041d9190611c84565b60405180910390f35b610440600480360361043b9190810190611279565b610792565b60405161044d9190611d27565b60405180910390f35b610470600480360361046b91908101906112e3565b61079f565b60405161047d9190611d64565b60405180910390f35b6104a0600480360361049b9190810190611558565b6107ac565b6040516104ad9190611edf565b60405180910390f35b6104d060048036036104cb9190810190611614565b6107b3565b6040516104dd9190611f52565b60405180910390f35b61050060048036036104fb919081019061148b565b6107ba565b60405161050d9190611e58565b60405180910390f35b610530600480360361052b919081019061152f565b6107c1565b60405161053d9190611ec4565b60405180910390f35b610560600480360361055b919081019061138e565b6107c8565b60405161056d9190611dc3565b60405180910390f35b610590600480360361058b91908101906114b4565b6107cf565b60405161059d9190611e73565b60405180910390f35b6105c060048036036105bb91908101906114dd565b6107d6565b6040516105cd9190611e8e565b60405180910390f35b6105f060048036036105eb9190810190611421565b6107dd565b6040516105fd9190611e1b565b60405180910390f35b610620600480360361061b9190810190611324565b6107e4565b60405161062d9190611d86565b60405180910390f35b610650600480360361064b91908101906115c2565b6107eb565b60405161065d9190611f1c565b60405180910390f35b610680600480360361067b919081019061144a565b6107f2565b60405161068d9190611e36565b60405180910390f35b6106b060048036036106ab91908101906110fa565b6107f9565b6040516106bd9190611c4e565b60405180910390f35b6106e060048036036106db9190810190611506565b610800565b6040516106ed9190611ea9565b60405180910390f35b610710600480360361070b91908101906113f8565b610807565b60405161071d9190611e00565b60405180910390f35b6000919050565b61073561080e565b919050565b610742610830565b919050565b6060919050565b6000919050565b6000919050565b6000919050565b6060919050565b610772610852565b919050565b6060919050565b610786610874565b919050565b6060919050565b61079a61089b565b919050565b6107a76108bd565b919050565b6060919050565b6000919050565b6000919050565b6000919050565b6060919050565b6000919050565b6000919050565b6000919050565b6060919050565b6000919050565b6060919050565b6000919050565b6000919050565b6000919050565b6040805190810160405280600290602082028038833980820191505090505090565b6040805190810160405280600290602082028038833980820191505090505090565b6040805190810160405280600290602082028038833980820191505090505090565b60408051908101604052806002905b60608152602001906001900390816108835790505090565b6040805190810160405280600290602082028038833980820191505090505090565b60408051908101604052806002905b60608152602001906001900390816108cc5790505090565b60006108f082356124f2565b905092915050565b600082601f830112151561090b57600080fd5b600261091e61091982611fb5565b611f88565b9150818385602084028201111561093457600080fd5b60005b83811015610964578161094a88826108e4565b845260208401935060208301925050600181019050610937565b5050505092915050565b600082601f830112151561098157600080fd5b813561099461098f82611fd7565b611f88565b915081818352602084019350602081019050838560208402820111156109b957600080fd5b60005b838110156109e957816109cf88826108e4565b8452602084019350602083019250506001810190506109bc565b5050505092915050565b600082601f8301121515610a0657600080fd5b6002610a19610a1482611fff565b611f88565b91508183856020840282011115610a2f57600080fd5b60005b83811015610a5f5781610a458882610e9e565b845260208401935060208301925050600181019050610a32565b5050505092915050565b600082601f8301121515610a7c57600080fd5b8135610a8f610a8a82612021565b611f88565b91508181835260208401935060208101905083856020840282011115610ab457600080fd5b60005b83811015610ae45781610aca8882610e9e565b845260208401935060208301925050600181019050610ab7565b5050505092915050565b600082601f8301121515610b0157600080fd5b6002610b14610b0f82612049565b611f88565b9150818360005b83811015610b4b5781358601610b318882610eda565b845260208401935060208301925050600181019050610b1b565b5050505092915050565b600082601f8301121515610b6857600080fd5b8135610b7b610b768261206b565b611f88565b9150818183526020840193506020810190508360005b83811015610bc15781358601610ba78882610eda565b845260208401935060208301925050600181019050610b91565b5050505092915050565b600082601f8301121515610bde57600080fd5b6002610bf1610bec82612093565b611f88565b91508183856020840282011115610c0757600080fd5b60005b83811015610c375781610c1d8882610f9a565b845260208401935060208301925050600181019050610c0a565b5050505092915050565b600082601f8301121515610c5457600080fd5b8135610c67610c62826120b5565b611f88565b91508181835260208401935060208101905083856020840282011115610c8c57600080fd5b60005b83811015610cbc5781610ca28882610f9a565b845260208401935060208301925050600181019050610c8f565b5050505092915050565b600082601f8301121515610cd957600080fd5b6002610cec610ce7826120dd565b611f88565b9150818360005b83811015610d235781358601610d098882610fea565b845260208401935060208301925050600181019050610cf3565b5050505092915050565b600082601f8301121515610d4057600080fd5b8135610d53610d4e826120ff565b611f88565b9150818183526020840193506020810190508360005b83811015610d995781358601610d7f8882610fea565b845260208401935060208301925050600181019050610d69565b5050505092915050565b600082601f8301121515610db657600080fd5b6002610dc9610dc482612127565b611f88565b91508183856020840282011115610ddf57600080fd5b60005b83811015610e0f5781610df588826110aa565b845260208401935060208301925050600181019050610de2565b5050505092915050565b600082601f8301121515610e2c57600080fd5b8135610e3f610e3a82612149565b611f88565b91508181835260208401935060208101905083856020840282011115610e6457600080fd5b60005b83811015610e945781610e7a88826110aa565b845260208401935060208301925050600181019050610e67565b5050505092915050565b6000610eaa8235612504565b905092915050565b6000610ebe8235612510565b905092915050565b6000610ed2823561253c565b905092915050565b600082601f8301121515610eed57600080fd5b8135610f00610efb82612171565b611f88565b91508082526020830160208301858383011115610f1c57600080fd5b610f278382846125cd565b50505092915050565b600082601f8301121515610f4357600080fd5b8135610f56610f518261219d565b611f88565b91508082526020830160208301858383011115610f7257600080fd5b610f7d8382846125cd565b50505092915050565b6000610f928235612546565b905092915050565b6000610fa68235612553565b905092915050565b6000610fba823561255d565b905092915050565b6000610fce823561256a565b905092915050565b6000610fe28235612577565b905092915050565b600082601f8301121515610ffd57600080fd5b813561101061100b826121c9565b611f88565b9150808252602083016020830185838301111561102c57600080fd5b6110378382846125cd565b50505092915050565b600082601f830112151561105357600080fd5b8135611066611061826121f5565b611f88565b9150808252602083016020830185838301111561108257600080fd5b61108d8382846125cd565b50505092915050565b60006110a28235612584565b905092915050565b60006110b68235612592565b905092915050565b60006110ca823561259c565b905092915050565b60006110de82356125ac565b905092915050565b60006110f282356125c0565b905092915050565b60006020828403121561110c57600080fd5b600061111a848285016108e4565b91505092915050565b60006040828403121561113557600080fd5b6000611143848285016108f8565b91505092915050565b60006020828403121561115e57600080fd5b600082013567ffffffffffffffff81111561117857600080fd5b6111848482850161096e565b91505092915050565b60006040828403121561119f57600080fd5b60006111ad848285016109f3565b91505092915050565b6000602082840312156111c857600080fd5b600082013567ffffffffffffffff8111156111e257600080fd5b6111ee84828501610a69565b91505092915050565b60006020828403121561120957600080fd5b600082013567ffffffffffffffff81111561122357600080fd5b61122f84828501610aee565b91505092915050565b60006020828403121561124a57600080fd5b600082013567ffffffffffffffff81111561126457600080fd5b61127084828501610b55565b91505092915050565b60006040828403121561128b57600080fd5b600061129984828501610bcb565b91505092915050565b6000602082840312156112b457600080fd5b600082013567ffffffffffffffff8111156112ce57600080fd5b6112da84828501610c41565b91505092915050565b6000602082840312156112f557600080fd5b600082013567ffffffffffffffff81111561130f57600080fd5b61131b84828501610cc6565b91505092915050565b60006020828403121561133657600080fd5b600082013567ffffffffffffffff81111561135057600080fd5b61135c84828501610d2d565b91505092915050565b60006040828403121561137757600080fd5b600061138584828501610da3565b91505092915050565b6000602082840312156113a057600080fd5b600082013567ffffffffffffffff8111156113ba57600080fd5b6113c684828501610e19565b91505092915050565b6000602082840312156113e157600080fd5b60006113ef84828501610e9e565b91505092915050565b60006020828403121561140a57600080fd5b600061141884828501610eb2565b91505092915050565b60006020828403121561143357600080fd5b600061144184828501610ec6565b91505092915050565b60006020828403121561145c57600080fd5b600082013567ffffffffffffffff81111561147657600080fd5b61148284828501610f30565b91505092915050565b60006020828403121561149d57600080fd5b60006114ab84828501610f86565b91505092915050565b6000602082840312156114c657600080fd5b60006114d484828501610f9a565b91505092915050565b6000602082840312156114ef57600080fd5b60006114fd84828501610fae565b91505092915050565b60006020828403121561151857600080fd5b600061152684828501610fc2565b91505092915050565b60006020828403121561154157600080fd5b600061154f84828501610fd6565b91505092915050565b60006020828403121561156a57600080fd5b600082013567ffffffffffffffff81111561158457600080fd5b61159084828501611040565b91505092915050565b6000602082840312156115ab57600080fd5b60006115b984828501611096565b91505092915050565b6000602082840312156115d457600080fd5b60006115e2848285016110aa565b91505092915050565b6000602082840312156115fd57600080fd5b600061160b848285016110be565b91505092915050565b60006020828403121561162657600080fd5b6000611634848285016110d2565b91505092915050565b60006020828403121561164f57600080fd5b600061165d848285016110e6565b91505092915050565b61166f816123f7565b82525050565b61167e816122ab565b61168782612221565b60005b828110156116b95761169d858351611666565b6116a68261235b565b915060208501945060018101905061168a565b5050505050565b60006116cb826122b6565b8084526020840193506116dd8361222b565b60005b8281101561170f576116f3868351611666565b6116fc82612368565b91506020860195506001810190506116e0565b50849250505092915050565b611724816122c1565b61172d82612238565b60005b8281101561175f57611743858351611ab3565b61174c82612375565b9150602085019450600181019050611730565b5050505050565b6000611771826122cc565b80845260208401935061178383612242565b60005b828110156117b557611799868351611ab3565b6117a282612382565b9150602086019550600181019050611786565b50849250505092915050565b60006117cc826122d7565b836020820285016117dc8561224f565b60005b848110156118155783830388526117f7838351611b16565b92506118028261238f565b91506020880197506001810190506117df565b508196508694505050505092915050565b6000611831826122e2565b8084526020840193508360208202850161184a85612259565b60005b84811015611883578383038852611865838351611b16565b92506118708261239c565b915060208801975060018101905061184d565b508196508694505050505092915050565b61189d816122ed565b6118a682612266565b60005b828110156118d8576118bc858351611b5b565b6118c5826123a9565b91506020850194506001810190506118a9565b5050505050565b60006118ea826122f8565b8084526020840193506118fc83612270565b60005b8281101561192e57611912868351611b5b565b61191b826123b6565b91506020860195506001810190506118ff565b50849250505092915050565b600061194582612303565b836020820285016119558561227d565b60005b8481101561198e578383038852611970838351611bcd565b925061197b826123c3565b9150602088019750600181019050611958565b508196508694505050505092915050565b60006119aa8261230e565b808452602084019350836020820285016119c385612287565b60005b848110156119fc5783830388526119de838351611bcd565b92506119e9826123d0565b91506020880197506001810190506119c6565b508196508694505050505092915050565b611a1681612319565b611a1f82612294565b60005b82811015611a5157611a35858351611c12565b611a3e826123dd565b9150602085019450600181019050611a22565b5050505050565b6000611a6382612324565b808452602084019350611a758361229e565b60005b82811015611aa757611a8b868351611c12565b611a94826123ea565b9150602086019550600181019050611a78565b50849250505092915050565b611abc81612409565b82525050565b611acb81612415565b82525050565b611ada81612441565b82525050565b6000611aeb8261233a565b808452611aff8160208601602086016125dc565b611b088161260f565b602085010191505092915050565b6000611b218261232f565b808452611b358160208601602086016125dc565b611b3e8161260f565b602085010191505092915050565b611b558161244b565b82525050565b611b6481612458565b82525050565b611b7381612462565b82525050565b611b828161246f565b82525050565b611b918161247c565b82525050565b6000611ba282612350565b808452611bb68160208601602086016125dc565b611bbf8161260f565b602085010191505092915050565b6000611bd882612345565b808452611bec8160208601602086016125dc565b611bf58161260f565b602085010191505092915050565b611c0c81612489565b82525050565b611c1b816124b7565b82525050565b611c2a816124c1565b82525050565b611c39816124d1565b82525050565b611c48816124e5565b82525050565b6000602082019050611c636000830184611666565b92915050565b6000604082019050611c7e6000830184611675565b92915050565b60006020820190508181036000830152611c9e81846116c0565b905092915050565b6000604082019050611cbb600083018461171b565b92915050565b60006020820190508181036000830152611cdb8184611766565b905092915050565b60006020820190508181036000830152611cfd81846117c1565b905092915050565b60006020820190508181036000830152611d1f8184611826565b905092915050565b6000604082019050611d3c6000830184611894565b92915050565b60006020820190508181036000830152611d5c81846118df565b905092915050565b60006020820190508181036000830152611d7e818461193a565b905092915050565b60006020820190508181036000830152611da0818461199f565b905092915050565b6000604082019050611dbd6000830184611a0d565b92915050565b60006020820190508181036000830152611ddd8184611a58565b905092915050565b6000602082019050611dfa6000830184611ab3565b92915050565b6000602082019050611e156000830184611ac2565b92915050565b6000602082019050611e306000830184611ad1565b92915050565b60006020820190508181036000830152611e508184611ae0565b905092915050565b6000602082019050611e6d6000830184611b4c565b92915050565b6000602082019050611e886000830184611b5b565b92915050565b6000602082019050611ea36000830184611b6a565b92915050565b6000602082019050611ebe6000830184611b79565b92915050565b6000602082019050611ed96000830184611b88565b92915050565b60006020820190508181036000830152611ef98184611b97565b905092915050565b6000602082019050611f166000830184611c03565b92915050565b6000602082019050611f316000830184611c12565b92915050565b6000602082019050611f4c6000830184611c21565b92915050565b6000602082019050611f676000830184611c30565b92915050565b6000602082019050611f826000830184611c3f565b92915050565b6000604051905081810181811067ffffffffffffffff82111715611fab57600080fd5b8060405250919050565b600067ffffffffffffffff821115611fcc57600080fd5b602082029050919050565b600067ffffffffffffffff821115611fee57600080fd5b602082029050602081019050919050565b600067ffffffffffffffff82111561201657600080fd5b602082029050919050565b600067ffffffffffffffff82111561203857600080fd5b602082029050602081019050919050565b600067ffffffffffffffff82111561206057600080fd5b602082029050919050565b600067ffffffffffffffff82111561208257600080fd5b602082029050602081019050919050565b600067ffffffffffffffff8211156120aa57600080fd5b602082029050919050565b600067ffffffffffffffff8211156120cc57600080fd5b602082029050602081019050919050565b600067ffffffffffffffff8211156120f457600080fd5b602082029050919050565b600067ffffffffffffffff82111561211657600080fd5b602082029050602081019050919050565b600067ffffffffffffffff82111561213e57600080fd5b602082029050919050565b600067ffffffffffffffff82111561216057600080fd5b602082029050602081019050919050565b600067ffffffffffffffff82111561218857600080fd5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156121b457600080fd5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156121e057600080fd5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff82111561220c57600080fd5b601f19601f8301169050602081019050919050565b6000819050919050565b6000602082019050919050565b6000819050919050565b6000602082019050919050565b6000819050919050565b6000602082019050919050565b6000819050919050565b6000602082019050919050565b6000819050919050565b6000602082019050919050565b6000819050919050565b6000602082019050919050565b600060029050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600060029050919050565b600081519050919050565b600081519050919050565b600081519050919050565b600081519050919050565b600081519050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b6000602082019050919050565b600061240282612497565b9050919050565b60008115159050919050565b60007fff0000000000000000000000000000000000000000000000000000000000000082169050919050565b6000819050919050565b60008160010b9050919050565b6000819050919050565b60008160030b9050919050565b60008160070b9050919050565b60008160000b9050919050565b600061ffff82169050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600063ffffffff82169050919050565b600067ffffffffffffffff82169050919050565b600060ff82169050919050565b60006124fd82612497565b9050919050565b60008115159050919050565b60007fff0000000000000000000000000000000000000000000000000000000000000082169050919050565b6000819050919050565b60008160010b9050919050565b6000819050919050565b60008160030b9050919050565b60008160070b9050919050565b60008160000b9050919050565b600061ffff82169050919050565b6000819050919050565b600063ffffffff82169050919050565b600067ffffffffffffffff82169050919050565b600060ff82169050919050565b82818337600083830152505050565b60005b838110156125fa5780820151818401526020810190506125df565b83811115612609576000848401525b50505050565b6000601f19601f830116905091905056fea265627a7a723058206fe37171cf1b10ebd291cfdca61d67e7fc3c208795e999c833c42a14d86cf00d6c6578706572696d656e74616cf50037";

	// deploy deploys a new Ethereum contract, binding an instance of Test to it.
	public static Test deploy(TransactOpts auth, EthereumClient client) throws Exception {
		Interfaces args = Geth.newInterfaces(0);
		String bytecode = BYTECODE;
		return new Test(Geth.deployContract(auth, ABI, Geth.decodeFromHex(bytecode), client, args));
	}

	// Internal constructor used by contract deployment.
	private Test(BoundContract deployment) {
		this.Address  = deployment.getAddress();
		this.Deployer = deployment.getDeployer();
		this.Contract = deployment;
	}

	// Ethereum address where this contract is located at.
	public final Address Address;

	// Ethereum transaction in which this contract was deployed (if known!).
	public final Transaction Deployer;

	// Contract instance bound to a blockchain address.
	private final BoundContract Contract;

	// Creates a new instance of Test, bound to a specific deployed contract.
	public Test(Address address, EthereumClient client) throws Exception {
		this(Geth.bindContract(address, ABI, client));
	}

	// setAddress is a paid mutator transaction binding the contract method 0xe30081a0.
	//
	// Solidity: function setAddress(address a) returns(address)
	public Transaction setAddress(TransactOpts opts, Address a) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setAddress(a);args.set(0,arg0);

		return this.Contract.transact(opts, "setAddress"	, args);
	}

	// setAddressArray is a paid mutator transaction binding the contract method 0x151f5471.
	//
	// Solidity: function setAddressArray(address[2] a_a) returns(address[2])
	public Transaction setAddressArray(TransactOpts opts, Addresses a_a) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setAddresses(a_a);args.set(0,arg0);

		return this.Contract.transact(opts, "setAddressArray"	, args);
	}

	// setAddressList is a paid mutator transaction binding the contract method 0x5be6b37e.
	//
	// Solidity: function setAddressList(address[] a_l) returns(address[])
	public Transaction setAddressList(TransactOpts opts, Addresses a_l) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setAddresses(a_l);args.set(0,arg0);

		return this.Contract.transact(opts, "setAddressList"	, args);
	}

	// setBool is a paid mutator transaction binding the contract method 0x1e26fd33.
	//
	// Solidity: function setBool(bool b) returns(bool)
	public Transaction setBool(TransactOpts opts, boolean b) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBool(b);args.set(0,arg0);

		return this.Contract.transact(opts, "setBool"	, args);
	}

	// setBoolArray is a paid mutator transaction binding the contract method 0x118a9718.
	//
	// Solidity: function setBoolArray(bool[2] b_a) returns(bool[2])
	public Transaction setBoolArray(TransactOpts opts, Bools b_a) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBools(b_a);args.set(0,arg0);

		return this.Contract.transact(opts, "setBoolArray"	, args);
	}

	// setBoolList is a paid mutator transaction binding the contract method 0x4d5ee6da.
	//
	// Solidity: function setBoolList(bool[] b_l) returns(bool[])
	public Transaction setBoolList(TransactOpts opts, Bools b_l) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBools(b_l);args.set(0,arg0);

		return this.Contract.transact(opts, "setBoolList"	, args);
	}

	// setBytes is a paid mutator transaction binding the contract method 0xda359dc8.
	//
	// Solidity: function setBytes(bytes bs) returns(bytes)
	public Transaction setBytes(TransactOpts opts, byte[] bs) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBinary(bs);args.set(0,arg0);

		return this.Contract.transact(opts, "setBytes"	, args);
	}

	// setBytes1 is a paid mutator transaction binding the contract method 0xfba1a1c3.
	//
	// Solidity: function setBytes1(bytes1 b1) returns(bytes1)
	public Transaction setBytes1(TransactOpts opts, byte[] b1) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBinary(b1);args.set(0,arg0);

		return this.Contract.transact(opts, "setBytes1"	, args);
	}

	// setBytes32 is a paid mutator transaction binding the contract method 0xc2b12a73.
	//
	// Solidity: function setBytes32(bytes32 b32) returns(bytes32)
	public Transaction setBytes32(TransactOpts opts, byte[] b32) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBinary(b32);args.set(0,arg0);

		return this.Contract.transact(opts, "setBytes32"	, args);
	}

	// setBytesArray is a paid mutator transaction binding the contract method 0x5119655d.
	//
	// Solidity: function setBytesArray(bytes[2] bs_a) returns(bytes[2])
	public Transaction setBytesArray(TransactOpts opts, Binaries bs_a) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBinaries(bs_a);args.set(0,arg0);

		return this.Contract.transact(opts, "setBytesArray"	, args);
	}

	// setBytesList is a paid mutator transaction binding the contract method 0x16c105e2.
	//
	// Solidity: function setBytesList(bytes[] bs_l) returns(bytes[])
	public Transaction setBytesList(TransactOpts opts, Binaries bs_l) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBinaries(bs_l);args.set(0,arg0);

		return this.Contract.transact(opts, "setBytesList"	, args);
	}

	// setInt16 is a paid mutator transaction binding the contract method 0x86114cea.
	//
	// Solidity: function setInt16(int16 i16) returns(int16)
	public Transaction setInt16(TransactOpts opts, short i16) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setInt16(i16);args.set(0,arg0);

		return this.Contract.transact(opts, "setInt16"	, args);
	}

	// setInt256 is a paid mutator transaction binding the contract method 0xa53b1c1e.
	//
	// Solidity: function setInt256(int256 i256) returns(int256)
	public Transaction setInt256(TransactOpts opts, BigInt i256) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBigInt(i256);args.set(0,arg0);

		return this.Contract.transact(opts, "setInt256"	, args);
	}

	// setInt256Array is a paid mutator transaction binding the contract method 0x6aa482fc.
	//
	// Solidity: function setInt256Array(int256[2] i256_a) returns(int256[2])
	public Transaction setInt256Array(TransactOpts opts, BigInts i256_a) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBigInts(i256_a);args.set(0,arg0);

		return this.Contract.transact(opts, "setInt256Array"	, args);
	}

	// setInt256List is a paid mutator transaction binding the contract method 0x22722302.
	//
	// Solidity: function setInt256List(int256[] i256_l) returns(int256[])
	public Transaction setInt256List(TransactOpts opts, BigInts i256_l) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBigInts(i256_l);args.set(0,arg0);

		return this.Contract.transact(opts, "setInt256List"	, args);
	}

	// setInt32 is a paid mutator transaction binding the contract method 0xb7d5df31.
	//
	// Solidity: function setInt32(int32 i32) returns(int32)
	public Transaction setInt32(TransactOpts opts, int i32) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setInt32(i32);args.set(0,arg0);

		return this.Contract.transact(opts, "setInt32"	, args);
	}

	// setInt64 is a paid mutator transaction binding the contract method 0xe673eb32.
	//
	// Solidity: function setInt64(int64 i64) returns(int64)
	public Transaction setInt64(TransactOpts opts, long i64) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setInt64(i64);args.set(0,arg0);

		return this.Contract.transact(opts, "setInt64"	, args);
	}

	// setInt8 is a paid mutator transaction binding the contract method 0x9a19a953.
	//
	// Solidity: function setInt8(int8 i8) returns(int8)
	public Transaction setInt8(TransactOpts opts, byte i8) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setInt8(i8);args.set(0,arg0);

		return this.Contract.transact(opts, "setInt8"	, args);
	}

	// setString is a paid mutator transaction binding the contract method 0x7fcaf666.
	//
	// Solidity: function setString(string s) returns(string)
	public Transaction setString(TransactOpts opts, String s) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setString(s);args.set(0,arg0);

		return this.Contract.transact(opts, "setString"	, args);
	}

	// setStringArray is a paid mutator transaction binding the contract method 0x7173b695.
	//
	// Solidity: function setStringArray(string[2] s_a) returns(string[2])
	public Transaction setStringArray(TransactOpts opts, Strings s_a) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setStrings(s_a);args.set(0,arg0);

		return this.Contract.transact(opts, "setStringArray"	, args);
	}

	// setStringList is a paid mutator transaction binding the contract method 0xc5777961.
	//
	// Solidity: function setStringList(string[] s_l) returns(string[])
	public Transaction setStringList(TransactOpts opts, Strings s_l) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setStrings(s_l);args.set(0,arg0);

		return this.Contract.transact(opts, "setStringList"	, args);
	}

	// setUint16 is a paid mutator transaction binding the contract method 0x0477988a.
	//
	// Solidity: function setUint16(uint16 u16) returns(uint16)
	public Transaction setUint16(TransactOpts opts, BigInt u16) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setUint16(u16);args.set(0,arg0);

		return this.Contract.transact(opts, "setUint16"	, args);
	}

	// setUint256 is a paid mutator transaction binding the contract method 0xd2282dc5.
	//
	// Solidity: function setUint256(uint256 u256) returns(uint256)
	public Transaction setUint256(TransactOpts opts, BigInt u256) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBigInt(u256);args.set(0,arg0);

		return this.Contract.transact(opts, "setUint256"	, args);
	}

	// setUint256Array is a paid mutator transaction binding the contract method 0x2766a755.
	//
	// Solidity: function setUint256Array(uint256[2] u256_a) returns(uint256[2])
	public Transaction setUint256Array(TransactOpts opts, BigInts u256_a) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBigInts(u256_a);args.set(0,arg0);

		return this.Contract.transact(opts, "setUint256Array"	, args);
	}

	// setUint256List is a paid mutator transaction binding the contract method 0xa0709e19.
	//
	// Solidity: function setUint256List(uint256[] u256_l) returns(uint256[])
	public Transaction setUint256List(TransactOpts opts, BigInts u256_l) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setBigInts(u256_l);args.set(0,arg0);

		return this.Contract.transact(opts, "setUint256List"	, args);
	}

	// setUint32 is a paid mutator transaction binding the contract method 0x1c9352e2.
	//
	// Solidity: function setUint32(uint32 u32) returns(uint32)
	public Transaction setUint32(TransactOpts opts, BigInt u32) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setUint32(u32);args.set(0,arg0);

		return this.Contract.transact(opts, "setUint32"	, args);
	}

	// setUint64 is a paid mutator transaction binding the contract method 0x822cba69.
	//
	// Solidity: function setUint64(uint64 u64) returns(uint64)
	public Transaction setUint64(TransactOpts opts, BigInt u64) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setUint64(u64);args.set(0,arg0);

		return this.Contract.transact(opts, "setUint64"	, args);
	}

	// setUint8 is a paid mutator transaction binding the contract method 0x1774e646.
	//
	// Solidity: function setUint8(uint8 u8) returns(uint8)
	public Transaction setUint8(TransactOpts opts, BigInt u8) throws Exception {
		Interfaces args = Geth.newInterfaces(1);
		Interface arg0 = Geth.newInterface();arg0.setUint8(u8);args.set(0,arg0);

		return this.Contract.transact(opts, "setUint8"	, args);
	}
}
`,
		},
	}
	for i, c := range cases {
		binding, err := Bind([]string{c.name}, []string{c.abi}, []string{c.bytecode}, nil, "bindtest", LangJava, nil, nil)
		if err != nil {
			t.Fatalf("test %d: failed to generate binding: %v", i, err)
		}
		// Remove empty lines
		removeEmptys := func(input string) string {
			lines := strings.Split(input, "\n")
			var index int
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					lines[index] = line
					index += 1
				}
			}
			lines = lines[:index]
			return strings.Join(lines, "\n")
		}
		binding = removeEmptys(binding)
		expect := removeEmptys(c.expected)
		if binding != expect {
			t.Fatalf("test %d: generated binding mismatch, has %s, want %s", i, binding, c.expected)
		}
	}
}
