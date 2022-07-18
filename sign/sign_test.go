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

func TestSignInfo2(t *testing.T) {
	info := &pb.SignedInfo{
		Peer:        "16Uiu2HAm7rY7Ta5k9vHoNkSAMzsN7kunuPmw2ojrGCz9ob1gSQnu",
		CreatedTime: 1658141357,
		Version:     "2.1.3",
		Nonce:       1,
		BttcAddress: "0xf713ce4be8498782765abb35fe7b5785aaafe30a",
		SignedTime:  1658141357,
	}
	signature, hash, err := SignInfo(info)
	fmt.Println("signature, err = ", hexutil.Encode(signature), err)
	fmt.Println("hashBytes, err = ", hexutil.Encode(hash), err)
	fmt.Println()

	address, err := RecoverInfo(signature, hash)
	fmt.Println("parse signer = ", address.String(), err)
}

func TestSignInfo3(t *testing.T) {
	signature := []byte("0x440d717771ce7eeba45aa4517b23b2dfc874732413a571bda2a2ea908a4f34df50b58057c1966b1d66869a0dead3e147847c02ba6899f2fdc6edbab775929add01")
	hash := []byte("0x4fab5dc163ed4c85f1349a00c000c5c963bc73cfa10031a8f4e9fead4ede82a6")
	fmt.Println("signature, err = ", hexutil.Encode(signature), nil)
	fmt.Println("hashBytes, err = ", hexutil.Encode(hash), nil)
	fmt.Println()

	address, err := RecoverInfo(signature, hash)
	fmt.Println("parse signer = ", address.String(), err)
}
