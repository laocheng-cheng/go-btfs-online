package common

import (
	"encoding/hex"

	eth "github.com/ethereum/go-ethereum/crypto"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	cp "github.com/tron-us/go-btfs-common/crypto"
)

func ConvertPeerID2BttcAddr(peerID string) (string, error) {
	tmp, err := peer.IDB58Decode(peerID)
	pppk, _ := tmp.ExtractPublicKey()

	pkBytes, _ := ic.RawFull(pppk)
	pk2, err := eth.UnmarshalPubkey(pkBytes)
	if err != nil {
		return "", err
	}

	addr, err := cp.PublicKeyToAddress(*pk2)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(addr.Bytes())[2:], nil
}
