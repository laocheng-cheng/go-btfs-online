package sign

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	pb "github.com/tron-us/go-btfs-common/protos/online"
	"testing"
	"time"
)

func TestSignInfo(t *testing.T) {
	info := &pb.SignedInfo{
		Peer:        "1",
		CreatedTime: 1,
		Version:     "1",
		Nonce:       3,
		BttcAddress: "0x22df207EC3C8D18fEDeed87752C5a68E5b4f6FbD",
		SignedTime:  uint32(time.Now().Unix()),
	}
	signature, hash, err := SignInfo(info)
	fmt.Println("signature, err = ", hexutil.Encode(signature), err)
	fmt.Println("hashBytes, err = ", hexutil.Encode(hash), err)
	fmt.Println()

	address, err := RecoverInfo(signature, hash)
	fmt.Println("parse signer = ", address.String(), err)
}
