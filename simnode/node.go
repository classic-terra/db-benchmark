package simnode

import (
	"fmt"
	"os"

	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
	sm "github.com/tendermint/tendermint/state"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func GetNode() (*App, error) {
	db, state, genDoc, err := TendermintHandleGenesis()
	if err != nil {
		return nil, err
	}

	fmt.Printf("getting genesis from chain = %v \n", state.ChainID)

	app, err := CosmosHandleGenesis(db, genDoc)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func TendermintHandleGenesis() (dbm.DB, sm.State, *tmtypes.GenesisDoc, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, sm.State{}, nil, err
	}

	db, err := GetDB(rootDir)
	if err != nil {
		return nil, sm.State{}, nil, err
	}

	config := &config.Config{}
	config.Genesis = "data/genesis.json"
	config.RootDir = rootDir
	genesisDocProvider := node.DefaultGenesisDocProviderFunc(config)

	fmt.Println("loading genesis state")

	state, genDoc, err := node.LoadStateFromDBOrGenesisDocProvider(db, genesisDocProvider)
	if err != nil {
		return nil, sm.State{}, nil, err
	}

	return db, state, genDoc, nil
}
