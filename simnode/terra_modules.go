package simnode

import (
	transfer "github.com/cosmos/ibc-go/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/modules/core"
	ibcclientclient "github.com/cosmos/ibc-go/modules/core/02-client/client"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	customauth "github.com/classic-terra/classic/custom/auth"
	customauthz "github.com/classic-terra/classic/custom/authz"
	custombank "github.com/classic-terra/classic/custom/bank"
	customcrisis "github.com/classic-terra/classic/custom/crisis"
	customdistr "github.com/classic-terra/classic/custom/distribution"
	customevidence "github.com/classic-terra/classic/custom/evidence"
	customfeegrant "github.com/classic-terra/classic/custom/feegrant"
	customgov "github.com/classic-terra/classic/custom/gov"
	custommint "github.com/classic-terra/classic/custom/mint"
	customparams "github.com/classic-terra/classic/custom/params"
	customslashing "github.com/classic-terra/classic/custom/slashing"
	customstaking "github.com/classic-terra/classic/custom/staking"
	customupgrade "github.com/classic-terra/classic/custom/upgrade"

	"github.com/classic-terra/classic/x/market"
	markettypes "github.com/classic-terra/classic/x/market/types"
	"github.com/classic-terra/classic/x/oracle"
	oracletypes "github.com/classic-terra/classic/x/oracle/types"
	"github.com/classic-terra/classic/x/treasury"
	treasurytypes "github.com/classic-terra/classic/x/treasury/types"
	"github.com/classic-terra/classic/x/vesting"
	"github.com/classic-terra/classic/x/wasm"
	wasmtypes "github.com/classic-terra/classic/x/wasm/types"
)

var (
	// ModuleBasics = The ModuleBasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		customauth.AppModuleBasic{},
		customauthz.AppModuleBasic{},
		custombank.AppModuleBasic{},
		capability.AppModuleBasic{},
		customstaking.AppModuleBasic{},
		custommint.AppModuleBasic{},
		customdistr.AppModuleBasic{},
		customgov.NewAppModuleBasic(
			paramsclient.ProposalHandler,
			distrclient.ProposalHandler,
			upgradeclient.ProposalHandler,
			upgradeclient.CancelProposalHandler,
			ibcclientclient.UpdateClientProposalHandler,
			ibcclientclient.UpgradeProposalHandler,
		),
		customparams.AppModuleBasic{},
		customcrisis.AppModuleBasic{},
		customslashing.AppModuleBasic{},
		customfeegrant.AppModuleBasic{},
		ibc.AppModuleBasic{},
		customupgrade.AppModuleBasic{},
		customevidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		oracle.AppModuleBasic{},
		market.AppModuleBasic{},
		treasury.AppModuleBasic{},
		wasm.AppModuleBasic{},
	)

	Modules = []string{
		capabilitytypes.ModuleName, authtypes.ModuleName,
		banktypes.ModuleName, distrtypes.ModuleName,
		stakingtypes.ModuleName, slashingtypes.ModuleName,
		govtypes.ModuleName, markettypes.ModuleName,
		oracletypes.ModuleName, treasurytypes.ModuleName,
		wasmtypes.ModuleName, authz.ModuleName,
		minttypes.ModuleName, crisistypes.ModuleName,
		ibchost.ModuleName,
		evidencetypes.ModuleName, ibctransfertypes.ModuleName,
		feegrant.ModuleName,
	}

	Keys = sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		oracletypes.StoreKey, markettypes.StoreKey, treasurytypes.StoreKey,
		wasmtypes.StoreKey, authzkeeper.StoreKey, feegrant.StoreKey,
	)

	Tkeys = sdk.NewTransientStoreKeys(
		paramstypes.TStoreKey,
	)

	MemKeys = sdk.NewMemoryStoreKeys(
		capabilitytypes.MemStoreKey,
	)
)
