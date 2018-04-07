package testing

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/examples/basecoin/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/magiconair/properties/assert"
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/tmlibs/log"
)

func loggerAndDBs() (log.Logger, map[string]dbm.DB) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	dbs := map[string]dbm.DB{
		"main":    dbm.NewMemDB(),
		"acc":     dbm.NewMemDB(),
		"ibc":     dbm.NewMemDB(),
		"staking": dbm.NewMemDB(),
	}
	return logger, dbs
}

func newBasecoinApp() *BasecoinApp {
	logger, dbs := loggerAndDBs()
	return NewBasecoinApp(logger, dbs)
}

func setGenesisAccounts(bapp *BasecoinApp, accs ...auth.BaseAccount) error {
	genaccs := make([]*types.GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = types.NewGenesisAccount(&types.AppAccount{acc, accName})
	}

	genesisState := types.GenesisState{
		Accounts: genaccs,
	}

	stateBytes, err := json.MarshalIndent(genesisState, "", "\t")
	if err != nil {
		return err
	}

	// Initialize the chain
	vals := []abci.Validator{}
	bapp.InitChain(abci.RequestInitChain{vals, stateBytes})
	bapp.Commit()

	return nil
}

func SignCheckDeliver(t *testing.T, bapp *BasecoinApp, msg sdk.Msg, seq []int64, expPass bool, priv ...crypto.PrivKeyEd25519) {

	// Sign the tx
	tx := genTx(msg, seq, priv...)
	// Run a Check
	res := bapp.Check(tx)
	if expPass {
		require.Equal(t, sdk.CodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
	}

	// Simulate a Block
	bapp.BeginBlock(abci.RequestBeginBlock{})
	res = bapp.Deliver(tx)
	if expPass {
		require.Equal(t, sdk.CodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
	}
	bapp.EndBlock(abci.RequestEndBlock{})
	//bapp.Commit()
}

func CheckBalance(t *testing.T, bapp *BasecoinApp, addr sdk.Address, balExpected string) {
	ctxDeliver := bapp.BaseApp.NewContext(false, abci.Header{})
	res2 := bapp.accountMapper.GetAccount(ctxDeliver, addr)
	assert.Equal(t, balExpected, fmt.Sprintf("%v", res2.GetCoins()))
}
