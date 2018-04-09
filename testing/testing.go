
import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/examples/basecoin/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	crypto "github.com/tendermint/go-crypto"
)

// Extended ABCI application
type TestApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the substores
	capKeyMainStore    *sdk.KVStoreKey
	capKeyAccountStore *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper sdk.AccountMapper
}

func NewTestApp() {

}

// generate a priv key and return it with its address
func generateKeys() (crypto.PrivKey, crypto.PubKey, sdk.Address) {
	priv := crypto.GenPrivKeyEd25519()
	pub := priv.PubKey()
	addr := pub.Address()
	return priv, pub, addr
}

// generate a priv key and return it with its address
func GenerateBaseAccounts(int32 numAccount) ([]sdk.Account, []crypto.PrivKey) {
	accounts := []BaseAccount
	privatekeys := []crypto.PrivKey

	for i := 0; i < numAccount; i++ {
		priv, pub, addr := generateKeys()
		baseAcc := BaseAcccount{
			Address: addr,
			PubKey: pub,
		}
		accounts = append(accounts, baseAcc)
		privatekeys = append(privatekeys, priv)
	}
	return accounts, privatekeys
}

func setupMultiStore(capkeys ...string) (sdk.MultiStore, []*sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	KVStoreKeys := []*sdk.KVStoreKey
	for _, capkey := range capkeys {
        kvstorekey := sdk.NewKVStoreKey(capkey)
        KVStoreKeys = append(KVStoreKeys, kvstorekey)
		ms.MountStoreWithDB(kvstorekey, sdk.StoreTypeIAVL, db)
    }
	ms.LoadLatestVersion()
	return ms, KVStoreKeys
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
