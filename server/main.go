package main

import (
	"github.com/SaiNageswarS/go-api-boot/server"
	pb "github.com/anuragnitt/agri-backend/generated"
)

var grpcPort = ":8081"
var webPort = ":8080"

func main() {
	server.LoadSecretsIntoEnv(false)
	inject := NewInject()
	bootServer := server.NewGoApiBoot()

	pb.RegisterSearchByProductServiceServer(bootServer.GrpcServer, inject.SearchByProductService)
	pb.RegisterSearchNearbyShopServiceServer(bootServer.GrpcServer, inject.SearchNearbyShopService)

	bootServer.Start(grpcPort, webPort)
}
