#!/usr/bin/make -f

GOFILES := $(shell ls -1 | grep .go | grep -v _test.go)

run-main-badgerdb:
	rm -rf testdb*
	go run -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb' -tags 'badgerdb' $(GOFILES)

run-main-pebbledb:
	rm -rf testdb*
	go run -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb' -tags 'pebbledb' $(GOFILES)

run-main-rocksdb:
	rm -rf testdb*
	go run -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb' -tags 'rocksdb' $(GOFILES)

run-benchmark-badgerdb:
	rm -rf testdb*
	mkdir -p benchmark-data 
	go test -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb' -tags 'badgerdb' -bench BenchmarkOrderedKeys -count=3 -benchmem github.com/classic-terra/db-benchmark -cpuprofile benchmark-data/badger-cpu.prof -memprofile benchmark-data/badger-mem.prof

run-benchmark-pebbledb:
	rm -rf testdb*
	mkdir -p benchmark-data
	go test -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb' -tags 'pebbledb' -bench BenchmarkOrderedKeys -count=3 -benchmem github.com/classic-terra/db-benchmark -cpuprofile benchmark-data/pebble-cpu.prof -memprofile benchmark-data/pebble-mem.prof

run-benchmark-rocksdb:
	rm -rf testdb*
	mkdir -p benchmark-data
	go test -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb' -tags 'rocksdb' -bench BenchmarkOrderedKeys -count=3 -benchmem github.com/classic-terra/db-benchmark -cpuprofile benchmark-data/pebble-cpu.prof -memprofile benchmark-data/pebble-mem.prof

run-benchmark:
	make run-benchmark-pebbledb
	make run-benchmark-badgerdb