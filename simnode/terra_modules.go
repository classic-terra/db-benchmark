package simnode

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		params.AppModuleBasic{},
	)

	Modules = []string{
		authtypes.ModuleName,
		paramstypes.ModuleName,
	}

	Keys = sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		paramstypes.StoreKey,
	)

	Tkeys = sdk.NewTransientStoreKeys(
		paramstypes.TStoreKey,
	)
)
