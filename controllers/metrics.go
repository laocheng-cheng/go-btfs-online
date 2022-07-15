package controllers

import (
	"errors"
	"fmt"
	"time"

	nodepb "github.com/tron-us/go-btfs-common/protos/node"
	pb "github.com/tron-us/go-btfs-common/protos/online"
	"github.com/tron-us/go-common/v2/log"
	"github.com/tron-us/status-server/common"
	"github.com/tron-us/status-server/sign"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/libp2p/go-libp2p-core/peer"
	ic "github.com/libp2p/go-libp2p-crypto"
	"go.uber.org/zap"
)

func UpdateSignMetricsHandler(sm *pb.ReqSignMetrics) (resp *pb.RespSignMetrics, err error) {
	err = checkPeerKey(sm)
	if err != nil {
		return nil, err
	}

	// parse payLoad
	var payLoad pb.PayLoadInfo
	if err := payLoad.Unmarshal(sm.Payload); err != nil {
		// TODO: unify error messages
		log.Error(common.AuthenticationPayloadParseError, zap.String(common.StatusCode, common.AuthenticationPayloadParseError))
		return nil, err
	}
	nm := *payLoad.Node
	nm.TimeCreated = time.Now().UTC()

	// check decode node id
	id := nm.NodeId
	_, err = peer.Decode(id)
	if err != nil {
		return nil, nil
	}

	lastSignature := payLoad.LastSignature
	lastSignedInfo := payLoad.LastSignedInfo
	fmt.Println("----", len(lastSignature), string(lastSignature), lastSignedInfo)

	// 1.Verify last signature
	if len(lastSignature)>0 || (lastSignedInfo != nil && lastSignedInfo.Nonce>0) {
		lastSignedAddress, err := sign.RecoverInfoExt([]byte(lastSignature), lastSignedInfo)
		if err != nil {
			return nil, err
		}
		if !sign.VerifySignature(lastSignedAddress.String()) {
			return nil, errors.New("VerifySignature failed! ")
		}
	}

	// 2.check limit
	now := uint32(time.Now().Unix())
	if lastSignedInfo != nil && now-lastSignedInfo.SignedTime < 50*60 {
		return nil, errors.New("too many requests! ")
	}

	// 3.get signature
	signature, baseInfo, err := getSignedInfo(&nm, lastSignedInfo, now)
	log.Debug("online getSignedInfo",
		zap.String("signature", hexutil.Encode(signature)),
		zap.Any("baseInfo", baseInfo),
		zap.Error(err))
	if err != nil {
		return nil, err
	}

	// 4.return pack (signature, baseInfo) to host
	resp = &pb.RespSignMetrics{
		Code:        pb.ResponseCode_SUCCESS,
		SignedInfo: &pb.SignedInfo{
			Peer:        baseInfo.Peer,
			CreatedTime: baseInfo.CreatedTime,
			Version:     baseInfo.Version,
			Nonce:       baseInfo.Nonce,
			BttcAddress: baseInfo.BttcAddress,
			SignedTime:  baseInfo.SignedTime,
		},
		Signature:   hexutil.Encode(signature),
	}

	return resp, nil
}

func checkPeerKey(sm *pb.ReqSignMetrics) (err error) {
	// check empty
	if sm.Payload == nil || sm.PublicKey == nil || sm.Signature == nil {
		// TODO: unify error messages
		return fmt.Errorf("void authentication")
	}

	// verfity signature
	pubK, err := ic.UnmarshalPublicKey(sm.PublicKey)
	if err != nil {
		log.Error("cannot unmarshal public key", zap.String(string(sm.PublicKey), err.Error()))
		return err
	}
	valid, err := pubK.Verify(sm.Payload, sm.Signature)
	if err != nil {
		log.Error(common.AuthenticationVerifyError, zap.String("cannnot verify", err.Error()))
		return err
	}

	if !valid {
		// TODO: unify error messages
		log.Error(common.AuthenticationVerifyError, zap.String(common.StatusCode, common.AuthenticationVerifyError))
		return fmt.Errorf(common.AuthenticationVerifyError)
	}

	return nil
}


func getSignedInfo(nm *nodepb.Node, last *pb.SignedInfo, now uint32) (signature []byte, info *pb.SignedInfo, err error) {
	var nonce uint32
	var bttcAddr string
	var createTime uint32

	if last != nil {
		nonce = last.Nonce
		bttcAddr = last.BttcAddress
		createTime = last.CreatedTime
	} else {
		nonce = 0
		bttcAddr, err = common.ConvertPeerID2BttcAddr(nm.NodeId)
		if err != nil {
			return nil, nil, err
		}
		createTime = now
	}

	info = &pb.SignedInfo{
		Peer:        nm.NodeId,
		CreatedTime: createTime,
		Version:     nm.BtfsVersion,
		Nonce:       nonce + 1,
		BttcAddress: bttcAddr,
		SignedTime:  uint32(time.Now().Unix()),
	}
	signature, _, err = sign.SignInfo(info)
	if err != nil {
		return nil, nil, err
	}

	return signature, info, nil
}
