package simnode

import (
	"encoding/json"
	"os"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	storemulti "github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type App struct {
	// store group
	store             *storemulti.Store
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry codectypes.InterfaceRegistry

	// logger
	logger log.Logger

	// the module manager
	mm *module.Manager

	// keepers
	AccountKeeper authkeeper.AccountKeeper
	ParamsKeeper  paramskeeper.Keeper
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	for _, name := range Modules {
		if name == paramstypes.ModuleName {
			continue
		}

		paramsKeeper.Subspace(name)
	}

	return paramsKeeper
}

func CosmosHandleGenesis(db dbm.DB, genDoc *tmtypes.GenesisDoc) (*App, error) {
	InitConfig()

	// create new Store
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	store := storemulti.NewStore(db, logger)

	// mount store of each keeper
	for _, key := range Keys {
		store.MountStoreWithDB(key, types.StoreTypeIAVL, nil)
	}

	for _, key := range Tkeys {
		store.MountStoreWithDB(key, types.StoreTypeTransient, nil)
	}

	// load store to create new empty iavl store
	if err := store.LoadLatestVersion(); err != nil {
		return nil, err
	}

	// genesis state handling
	var genesisState GenesisState
	if err := json.Unmarshal(genDoc.AppState, &genesisState); err != nil {
		panic(err)
	}

	// get a basic module manager and init genesis
	ctx := sdk.NewContext(store, tmproto.Header{}, false, logger)

	app := NewApp(store)
	app.logger = logger
	for _, name := range Modules {
		app.mm.Modules[name].InitGenesis(ctx, app.appCodec, genesisState[name])
	}

	return app, nil
}

func NewApp(store *storemulti.Store) *App {
	encoding := GetEncodingConfig()

	app := &App{
		store:             store,
		legacyAmino:       encoding.Amino,
		appCodec:          encoding.Marshaler,
		interfaceRegistry: encoding.InterfaceRegistry,
	}

	app.ParamsKeeper = initParamsKeeper(app.appCodec, app.legacyAmino, Keys[paramstypes.StoreKey], Tkeys[paramstypes.TStoreKey])

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		app.appCodec, Keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, nil,
	)

	app.mm = module.NewManager(
		auth.NewAppModule(app.appCodec, app.AccountKeeper, nil),
		params.NewAppModule(app.ParamsKeeper),
	)

	app.mm.SetOrderInitGenesis(Modules...)

	return app
}

func GetEncodingConfig() EncodingConfig {
	encodingConfig := MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}

func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}
