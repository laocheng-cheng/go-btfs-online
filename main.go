package main

import (
	"context"

	pb "github.com/tron-us/go-btfs-common/protos/online"
	grpcUtils "github.com/tron-us/go-btfs-common/utils/grpc"
	"github.com/tron-us/go-common/v2/log"
	"github.com/tron-us/status-server/common"
	"github.com/tron-us/status-server/config"
	"github.com/tron-us/status-server/controllers"

	"go.uber.org/zap"
)

type OnlineServer struct {
	// for backward-compatibility
	pb.UnimplementedOnlineServiceServer
}

//implementation of the shared helper function
func (s *OnlineServer) UpdateSignMetrics(ctx context.Context, pbSM *pb.ReqSignMetrics) (*pb.RespSignMetrics, error) {
	resp, err := controllers.UpdateSignMetricsHandler(pbSM)
	if err != nil {
		log.Error(common.MetricsUpdateAndDiscoveryHandlerError, zap.Error(err))
		return new(pb.RespSignMetrics), err
	}
	if resp == nil {
		return new(pb.RespSignMetrics), nil
	} else {
		return resp, nil
	}
}

func main() {
	log.Info("Starting status server...")

	s := grpcUtils.GrpcServer{}
	s.GrpcServer(config.Host+config.Port, nil, "", &OnlineServer{})
	s.AcceptConnection()

	outerStop := make(chan bool, 1)
	<-outerStop

	log.Info("Status server is stopping...",
		zap.String("host", config.Host),
		zap.String("port", config.Port))
}
