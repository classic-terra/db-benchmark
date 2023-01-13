run-main-badgerdb:
	go run -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb' -tags 'badgerdb' *.go

run-main-pebbledb:
	go run -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb' -tags 'pebbledb' *.go

run-main-rocksdb:
	go run -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb' -tags 'rocksdb' *.go

run-benchmark-badgerdb:
	go test -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb' -tags 'badgerdb' -bench BenchmarkOrderedKeys -count=3 -benchmem github.com/nghuyenthevinh2000/db-benchmark -cpuprofile badger-cpu.prof -memprofile badger-mem.prof

run-benchmark-pebbledb:
	go test -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb' -tags 'pebbledb' -bench BenchmarkOrderedKeys -count=3 -benchmem github.com/nghuyenthevinh2000/db-benchmark -cpuprofile pebble-cpu.prof -memprofile pebble-mem.prof

run-benchmark-rocksdb:
	go test -ldflags '-X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb' -tags 'rocksdb' -bench BenchmarkOrderedKeys -count=3 -benchmem github.com/nghuyenthevinh2000/db-benchmark

run-benchmark:
	make run-benchmark-pebbledb
	make run-benchmark-badgerdb