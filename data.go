package main

import (
	"encoding/binary"
	"math/rand"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

func GetMockData(size uint64) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}

	return data
}

func PopulateStoreWithOrderedKeys(len, size uint64, store storetypes.CommitKVStore) {
	for i := uint64(0); i < len; i++ {
		store.Set(ConvertUint64ToBytes(i), GetMockData(size))
	}
}

func ConvertUint64ToBytes(num uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)

	return b
}
