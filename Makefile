run-main-badgerdb:
	go run -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb' -tags 'badgerdb' *.go

run-main-pebbledb:
	go run -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb' -tags 'pebbledb' *.go

run-main-rocksdb:
	go run -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb' -tags 'rocksdb' *.go

run-benchmark-badgerdb:
	mkdir -p benchmark-data 
	go test -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb' -tags 'badgerdb' -bench BenchmarkOrderedKeys -count=3 -benchmem github.com/nghuyenthevinh2000/db-benchmark -cpuprofile benchmark-data/badger-cpu.prof -memprofile benchmark-data/badger-mem.prof

run-benchmark-pebbledb:
	mkdir -p benchmark-data
	go test -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb' -tags 'pebbledb' -bench BenchmarkOrderedKeys -count=3 -benchmem github.com/nghuyenthevinh2000/db-benchmark -cpuprofile benchmark-data/pebble-cpu.prof -memprofile benchmark-data/pebble-mem.prof

run-benchmark-rocksdb:
	mkdir -p benchmark-data
	go test -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb' -tags 'rocksdb' -bench BenchmarkOrderedKeys -count=3 -benchmem github.com/nghuyenthevinh2000/db-benchmark -cpuprofile benchmark-data/pebble-cpu.prof -memprofile benchmark-data/pebble-mem.prof

run-benchmark:
	make run-benchmark-pebbledb
	make run-benchmark-badgerdb