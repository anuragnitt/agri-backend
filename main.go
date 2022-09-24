package main

import (
	"github.com/SaiNageswarS/go-api-boot/server"
	pb "github.com/anuragnitt/agri-backend/generated"
)

var grpcPort = ":4000"
var webPort = ":4001"

func main() {
	server.LoadSecretsIntoEnv(false)
	inject := NewInject()
	bootServer := server.NewGoApiBoot()

	pb.RegisterSearchByProductServiceServer(bootServer.GrpcServer, inject.SearchByProductService)
	pb.RegisterSearchNearbyShopServiceServer(bootServer.GrpcServer, inject.SearchNearbyShopService)

	bootServer.Start(grpcPort, webPort)
}
