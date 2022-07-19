package main

import (
	"fmt"

	"github.com/tron-us/status-server/common"
)

func main() {
	peer := "16Uiu2HAmEJEPHqAczriPjoe4Cb2dYMXAovNjz4T7GkrQRDbVUkMG"
	bttcAddr, err := common.ConvertPeerID2BttcAddr(peer)
	fmt.Println(bttcAddr, err)
}
