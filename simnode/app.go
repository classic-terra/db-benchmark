package simnode

import (
	"encoding/json"
	"os"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/std"
	storemulti "github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/types"

	transfer "github.com/cosmos/ibc-go/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/modules/core"
	ibcclient "github.com/cosmos/ibc-go/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/modules/core/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/classic-terra/classic/x/market"
	marketkeeper "github.com/classic-terra/classic/x/market/keeper"
	markettypes "github.com/classic-terra/classic/x/market/types"
	"github.com/classic-terra/classic/x/oracle"
	oraclekeeper "github.com/classic-terra/classic/x/oracle/keeper"
	oracletypes "github.com/classic-terra/classic/x/oracle/types"
	"github.com/classic-terra/classic/x/treasury"
	treasurykeeper "github.com/classic-terra/classic/x/treasury/keeper"
	treasurytypes "github.com/classic-terra/classic/x/treasury/types"
	"github.com/classic-terra/classic/x/wasm"
	wasmconfig "github.com/classic-terra/classic/x/wasm/config"
	wasmkeeper "github.com/classic-terra/classic/x/wasm/keeper"
	wasmtypes "github.com/classic-terra/classic/x/wasm/types"

	bankwasm "github.com/classic-terra/classic/custom/bank/wasm"
	distrwasm "github.com/classic-terra/classic/custom/distribution/wasm"
	govwasm "github.com/classic-terra/classic/custom/gov/wasm"
	stakingwasm "github.com/classic-terra/classic/custom/staking/wasm"
	marketwasm "github.com/classic-terra/classic/x/market/wasm"
	oraclewasm "github.com/classic-terra/classic/x/oracle/wasm"
	treasurywasm "github.com/classic-terra/classic/x/treasury/wasm"
)

var (
	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil, // just added to enable align fee
		treasurytypes.BurnModuleName:   {authtypes.Burner},
		minttypes.ModuleName:           {authtypes.Minter},
		markettypes.ModuleName:         {authtypes.Minter, authtypes.Burner},
		oracletypes.ModuleName:         nil,
		distrtypes.ModuleName:          nil,
		treasurytypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		oracletypes.ModuleName:       true,
		treasurytypes.BurnModuleName: true,
	}
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
	AccountKeeper    authkeeper.AccountKeeper
	AuthzKeeper      authzkeeper.Keeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper   evidencekeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	TransferKeeper   ibctransferkeeper.Keeper
	OracleKeeper     oraclekeeper.Keeper
	MarketKeeper     marketkeeper.Keeper
	TreasuryKeeper   treasurykeeper.Keeper
	WasmKeeper       wasmkeeper.Keeper
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	for _, name := range Modules {
		if name == paramstypes.ModuleName {
			continue
		}

		if name == govtypes.ModuleName {
			paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
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
	store := storemulti.NewStore(db)

	// mount store of each keeper
	for _, key := range Keys {
		store.MountStoreWithDB(key, types.StoreTypeIAVL, nil)
	}

	for _, key := range Tkeys {
		store.MountStoreWithDB(key, types.StoreTypeTransient, nil)
	}

	for _, key := range MemKeys {
		store.MountStoreWithDB(key, types.StoreTypeMemory, nil)
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
	appCodec := encoding.Marshaler

	app := &App{
		store:             store,
		legacyAmino:       encoding.Amino,
		appCodec:          appCodec,
		interfaceRegistry: encoding.InterfaceRegistry,
	}

	app.ParamsKeeper = initParamsKeeper(app.appCodec, app.legacyAmino, Keys[paramstypes.StoreKey], Tkeys[paramstypes.TStoreKey])

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, Keys[capabilitytypes.StoreKey], MemKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	// Applications that wish to enforce statically created ScopedKeepers should call `Seal` after creating
	// their scoped modules in `NewApp` with `ScopeToModule`
	app.CapabilityKeeper.Seal()

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, Keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, Keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlacklistedAccAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, Keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)
	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec, Keys[minttypes.StoreKey], app.GetSubspace(minttypes.ModuleName), &stakingKeeper,
		app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, Keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, app.BlacklistedAccAddrs(),
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, Keys[slashingtypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), 0, app.BankKeeper, authtypes.FeeCollectorName,
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, Keys[feegrant.StoreKey], app.AccountKeeper)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(nil, Keys[upgradetypes.StoreKey], appCodec, "", nil)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(Keys[authzkeeper.StoreKey], appCodec, nil)

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, Keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), app.StakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, Keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, Keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	// Initialize terra module keepers
	app.OracleKeeper = oraclekeeper.NewKeeper(
		appCodec, Keys[oracletypes.StoreKey], app.GetSubspace(oracletypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, app.DistrKeeper, &stakingKeeper, distrtypes.ModuleName,
	)
	app.MarketKeeper = marketkeeper.NewKeeper(
		appCodec, Keys[markettypes.StoreKey],
		app.GetSubspace(markettypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, app.OracleKeeper,
	)
	app.TreasuryKeeper = treasurykeeper.NewKeeper(
		appCodec, Keys[treasurytypes.StoreKey],
		app.GetSubspace(treasurytypes.ModuleName),
		app.AccountKeeper, app.BankKeeper,
		app.MarketKeeper, app.OracleKeeper,
		app.StakingKeeper, app.DistrKeeper,
		distrtypes.ModuleName)

	app.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec, Keys[wasmtypes.StoreKey],
		app.GetSubspace(wasmtypes.ModuleName),
		app.AccountKeeper, app.BankKeeper,
		app.TreasuryKeeper, nil,
		nil, wasmtypes.DefaultFeatures,
		"", wasmconfig.DefaultConfig(),
	)

	// register wasm msg parser & querier
	app.WasmKeeper.RegisterMsgParsers(map[string]wasmtypes.WasmMsgParserInterface{
		wasmtypes.WasmMsgParserRouteBank:         bankwasm.NewWasmMsgParser(),
		wasmtypes.WasmMsgParserRouteStaking:      stakingwasm.NewWasmMsgParser(),
		wasmtypes.WasmMsgParserRouteMarket:       marketwasm.NewWasmMsgParser(),
		wasmtypes.WasmMsgParserRouteWasm:         wasmkeeper.NewWasmMsgParser(),
		wasmtypes.WasmMsgParserRouteDistribution: distrwasm.NewWasmMsgParser(),
		wasmtypes.WasmMsgParserRouteGov:          govwasm.NewWasmMsgParser(),
	}, wasmkeeper.NewStargateWasmMsgParser(appCodec))
	app.WasmKeeper.RegisterQueriers(map[string]wasmtypes.WasmQuerierInterface{
		wasmtypes.WasmQueryRouteBank:     bankwasm.NewWasmQuerier(app.BankKeeper),
		wasmtypes.WasmQueryRouteStaking:  stakingwasm.NewWasmQuerier(app.StakingKeeper, app.DistrKeeper),
		wasmtypes.WasmQueryRouteMarket:   marketwasm.NewWasmQuerier(app.MarketKeeper),
		wasmtypes.WasmQueryRouteOracle:   oraclewasm.NewWasmQuerier(app.OracleKeeper),
		wasmtypes.WasmQueryRouteTreasury: treasurywasm.NewWasmQuerier(app.TreasuryKeeper),
		wasmtypes.WasmQueryRouteWasm:     wasmkeeper.NewWasmQuerier(app.WasmKeeper),
	}, wasmkeeper.NewStargateWasmQuerier(app.WasmKeeper))

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))
	app.GovKeeper = govkeeper.NewKeeper(
		appCodec, Keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govRouter,
	)

	/****  Module Options ****/
	skipGenesisInvariants := true

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		transferModule,
		market.NewAppModule(appCodec, app.MarketKeeper, app.AccountKeeper, app.BankKeeper, app.OracleKeeper),
		oracle.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
		treasury.NewAppModule(appCodec, app.TreasuryKeeper),
		wasm.NewAppModule(appCodec, app.WasmKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
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

// BlacklistedAccAddrs returns all the app's module account addresses black listed for receiving tokens.
func (app *App) BlacklistedAccAddrs() map[string]bool {
	blacklistedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blacklistedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blacklistedAddrs
}
