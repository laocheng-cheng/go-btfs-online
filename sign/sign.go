package sign

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	pb "github.com/tron-us/go-btfs-common/protos/online"
	"log"
)

var (
	privateKey *ecdsa.PrivateKey
	stringTy   abi.Type
	uint32Ty   abi.Type
	addressTy  abi.Type
	bytesTy    abi.Type
	arguments  abi.Arguments
	arguments2 abi.Arguments
)

const constSignAddr = "0x22df207EC3C8D18fEDeed87752C5a68E5b4f6FbD"

func init() {
	var err error

	privateKey, err = crypto.HexToECDSA("744ba22387c27cf73dff283a37f0a7e63054a86be15965be97c807816d79da39")
	if err != nil {
		log.Panic("sign privateKey")
	}

	stringTy, err = abi.NewType("string", "string", nil)
	if err != nil {
		log.Panic("sign stringTy")
	}
	uint32Ty, err = abi.NewType("uint32", "uint32", nil)
	if err != nil {
		log.Panic("sign uint32Ty")
	}
	addressTy, err = abi.NewType("address", "address", nil)
	if err != nil {
		log.Panic("sign addressTy")
	}
	bytesTy, err = abi.NewType("bytes", "bytes", nil)
	if err != nil {
		log.Panic("sign bytesTy")
	}

	arguments = abi.Arguments{
		{
			Type: stringTy,
		},
		{
			Type: uint32Ty,
		},
		{
			Type: stringTy,
		},
		{
			Type: uint32Ty,
		},
		{
			Type: addressTy,
		},
		{
			Type: uint32Ty,
		},
	}
	arguments2 = abi.Arguments{
		{
			Type: stringTy,
		},
		{
			Type: uint32Ty,
		},
		{
			Type: bytesTy,
		},
	}
}

type BaseInfo struct {
	Peer        string `json:"peer"`
	CreatedTime uint32 `json:"created_time"`
	Version     string `json:"version"`
	Nonce       uint32 `json:"nonce"`
	BttcAddress string `json:"bttc_address"`
	SignedTime  uint32 `json:"signed_time"`
}

func SignInfo(info *pb.SignedInfo) ([]byte, []byte, error) {
	data, err := arguments.Pack(string(info.Peer), uint32(info.CreatedTime), string(info.Version), uint32(info.Nonce),
		common.HexToAddress(info.BttcAddress), uint32(info.SignedTime))
	if err != nil {
		return nil, nil, err
	}
	//fmt.Println("data, err ", hexutil.Encode(data), err)

	msg, err := arguments2.Pack(string("\x19Ethereum Signed Message:\n"), uint32(len(data)), data)
	if err != nil {
		return nil, nil, err
	}
	//fmt.Println("msg, err ", hexutil.Encode(msg), err)

	hash := crypto.Keccak256Hash(msg)
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, nil, err
	}

	return signature, hash.Bytes(), nil
}

func RecoverInfo(signature []byte, hashBytes []byte) (common.Address, error) {
	sigPublicKeyECDSA, err := crypto.SigToPub(hashBytes, signature)
	if err != nil {
		return [20]byte{}, err
	}
	address := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	return address, nil
}

func RecoverInfoExt(signature []byte, info *pb.SignedInfo) (common.Address, error) {
	_, hash, err := SignInfo(info); if err != nil {
		return [20]byte{}, err
	}
	return RecoverInfo(signature, hash)
}

func VerifySignature(addr string) bool {
	if addr == constSignAddr {
		return true
	}
	return false
}

