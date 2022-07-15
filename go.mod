module github.com/tron-us/status-server

go 1.13

require (
	github.com/ethereum/go-ethereum v1.9.24
	github.com/libp2p/go-libp2p-core v0.0.6
	github.com/libp2p/go-libp2p-crypto v0.1.0
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/tron-us/go-btfs-common v0.8.9-pre11
	github.com/tron-us/go-common/v2 v2.3.0
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
)

replace github.com/libp2p/go-libp2p-core => github.com/TRON-US/go-libp2p-core v0.5.0
