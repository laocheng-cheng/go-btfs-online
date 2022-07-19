// Package main implements a client for Status service.
package main

import (
	"context"
	"fmt"

	"github.com/tron-us/go-btfs-common/crypto"
	nodepb "github.com/tron-us/go-btfs-common/protos/node"
	onlinepb "github.com/tron-us/go-btfs-common/protos/online"
	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	cgrpc "github.com/tron-us/go-btfs-common/utils/grpc"
	"github.com/tron-us/go-common/v2/log"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	"go.uber.org/zap"
)

const (
	//address       = "https://status-dev.btfs.io"
	address         = "localhost:50051"
	dummyNodeId     = "16Uiu2HAm7rY7Ta5k9vHoNkSAMzsN7kunuPmw2ojrGCz9ob1gSQnu"
	dummyPrivKeyStr = "CAISIFi8bw3uGI/FhBe6sQWDygxKRmmaAc3OCgrhTBj1c1kc"
)

func main() {
	ctx := context.Background()
	var err error

	//err := cgrpc.RuntimeClient(address).WithContext(ctx, func(ctx context.Context,
	//	client sharedpb.RuntimeServiceClient) error {
	//	checkRuntimeCaller(ctx, client)
	//	return nil
	//})
	//if err != nil {
	//	log.Panic(err.Error())
	//}

	// node
	err = cgrpc.OnlineClient(address).WithContext(ctx, func(ctx context.Context,
		client onlinepb.OnlineServiceClient) error {
		updateSignMetricsCaller(ctx, client)
		return nil
	})
	if err != nil {
		log.Panic(err.Error())
	}
	return
}

func checkRuntimeCaller(ctx context.Context, c sharedpb.RuntimeServiceClient) {
	req := new(sharedpb.SignedRuntimeInfoRequest)
	res, err := c.CheckRuntime(ctx, req)
	if err != nil {
		log.Panic("client", zap.Error(err))
	}
	fmt.Println("res: ", res)
}

func updateSignMetricsCaller(ctx context.Context, c onlinepb.OnlineServiceClient) {
	// create dummy keys for authentication
	dummyPrivKey, err := crypto.ToPrivKey(dummyPrivKeyStr)
	if err != nil {
		log.Panic("PrivKey convert error", zap.Error(err))
	}
	dummyPubKey, err := ic.MarshalPublicKey(dummyPrivKey.GetPublic())
	if err != nil {
		log.Panic("PublicKey marshal error", zap.Error(err))
	}

	// construct dummy node data
	node := new(nodepb.Node)
	node.NodeId = dummyNodeId

	discovery := new(nodepb.DiscoveryNode)
	discovery.ToNodeId = dummyNodeId //+ rand.String(4)
	discovery.NodeConnectLatency = 100
	discovery.ErrCode = nodepb.DiscoveryErrorCode_SUCCESS

	payLoad := new(onlinepb.PayLoadInfo)
	payLoad.NodeId = dummyNodeId
	payLoad.Node = node
	//payLoad.LastSignature
	//payLoad.LastSignedInfo

	// sign
	sm := new(onlinepb.ReqSignMetrics)
	sm.PublicKey = dummyPubKey
	sm.Payload, err = payLoad.Marshal()
	if err != nil {
		log.Panic("node marshal error", zap.Error(err))
	}
	sm.Signature, err = crypto.Sign(dummyPrivKey, payLoad)
	if err != nil {
		log.Panic("sign error", zap.Error(err))
	}

	// call update rpc
	res, err := c.UpdateSignMetrics(ctx, sm)
	if err != nil {
		log.Panic("client", zap.Error(err))
	}
	fmt.Println("res: ", res)
}
