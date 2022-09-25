package main

import (
	"fmt"
	"io"
	"context"

	pb "github.com/anuragnitt/agri-client/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

var grpcServer string = "localhost:8081"
var ctx context.Context
var cancel context.CancelFunc

func FindNearbyShops(c pb.SearchNearbyShopServiceClient, req *pb.SearchNearbyShopRequest) {
	shopStream, err := c.SearchNearbyShop(ctx, req)
	if err != nil {
		logger.Error("Error fetching nearby shops", zap.Error(err))
		return
	}

	for {
		shop, err := shopStream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			logger.Error("Error reading from shopStream", zap.Error(err))
			return
		}

		fmt.Printf("Shop name: %v, distance: %.6f, shopId: %v\n", shop.Shop.Name, shop.Shop.Distance, shop.Shop.Id)
	}

	fmt.Println("-----------------\n")
}

func FindProductsByName(c pb.SearchByProductServiceClient, req *pb.SearchByProductRequest) {
	productStream, err := c.SearchByProduct(ctx, req)
	if err != nil {
		logger.Error("Error fetching products by name", zap.Error(err))
		return
	}

	for {
		product, err := productStream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			logger.Error("Error reading from productStream", zap.Error(err))
			return
		}

		fmt.Printf("Product name: %v, price: %d, category: %v, shopId: %v\n", product.Product.Name, product.Product.Price, product.Product.Category, product.Product.ShopId)
	}

	fmt.Println("-----------------\n")
}

func FindProductsByAttr(c pb.SearchByProductServiceClient, req *pb.SearchByProductAttrRequest) {
	productStream, err := c.SearchByProductAttr(ctx, req)
	if err != nil {
		logger.Error("Error fetching products by attributes", zap.Error(err))
		return
	}

	for {
		product, err := productStream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			logger.Error("Error reading from productStream", zap.Error(err))
			return
		}

		fmt.Printf("Product name: %v, price: %d, category: %v, shopId: %v\n", product.Product.Name, product.Product.Price, product.Product.Category, product.Product.ShopId)
	}

	fmt.Println("-----------------\n")
}

func main() {
	conn, err := grpc.Dial(grpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Cannot connect to gRPC server", zap.Error(err))
	}

	shopClient := pb.NewSearchNearbyShopServiceClient(conn)
	productClient := pb.NewSearchByProductServiceClient(conn)

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	FindNearbyShops(shopClient, &pb.SearchNearbyShopRequest{
		Name: "Shop",
		Longitude: -1.20631,
		Latitude: 51.73213,
		Distance: 6000,
	})

	FindProductsByName(productClient, &pb.SearchByProductRequest{
		Name: "NEX",
	})

	FindProductsByAttr(productClient, &pb.SearchByProductAttrRequest{
		Name: "",
		MinPrice: 150,
		MaxPrice: 12000,
		Category: "clothing",
	})
}
