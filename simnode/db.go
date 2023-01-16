package simnode

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/store/iavl"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbm "github.com/tendermint/tm-db"
)

func GetDB(dir string) (dbm.DB, error) {
	name := fmt.Sprintf("testdb")
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("you are getting backend = %v \n", sdk.DBBackend)
	db, err := dbm.NewDB(name, dbm.BackendType(sdk.DBBackend), dir)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetIAVLKVStore() (storetypes.CommitKVStore, error) {
	db, err := GetDB("")
	if err != nil {
		return nil, err
	}

	lazyLoad := false
	kvStore, err := iavl.LoadStore(db, storetypes.CommitID{}, lazyLoad, iavl.DefaultIAVLCacheSize)
	if err != nil {
		return nil, err
	}

	return kvStore, nil
}
