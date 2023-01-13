package main

import (
	"fmt"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

var kvStore storetypes.CommitKVStore

func init() {
	var err error
	kvStore, err = GetKVStore()
	if err != nil {
		panic(err.Error())
	}
}

var table = []struct {
	keysNum      uint64
	dataByteSize uint64
}{
	{
		keysNum:      uint64(1000),
		dataByteSize: uint64(800),
	},
}

func BenchmarkOrderedKeys(b *testing.B) {
	// set 1000 keys (8bytes) of value (8bytes) to kv store
	for _, v := range table {
		b.Run(fmt.Sprintf("running ordered keysNum = %v with dataByteSize = %v bytes", v.keysNum, v.dataByteSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// writing action
				PopulateStoreWithOrderedKeys(v.keysNum, v.dataByteSize, kvStore)

				// commit store to disk
				kvStore.Commit()

				// reading action
				iter := kvStore.Iterator(nil, nil)
				for ; iter.Valid(); iter.Next() {
					iter.Value()
				}
				iter.Close()
			}
		})
	}
}
