package service

import (
	"context"
	"strings"

	pb "github.com/anuragnitt/agri-backend/generated"
	"github.com/anuragnitt/agri-backend/db"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

type SearchNearbyShopService struct {
	pb.UnimplementedSearchNearbyShopServiceServer
	db *db.AgriDb
}

func NewSearchNearbyShopService(db *db.AgriDb) *SearchNearbyShopService {
	return &SearchNearbyShopService{
		db: db,
	}
}

func (s *SearchNearbyShopService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

func (s *SearchNearbyShopService) SearchNearbyShop(req *pb.SearchNearbyShopRequest, stream pb.SearchNearbyShopService_SearchNearbyShopServer) error {
	name := strings.TrimSpace(req.GetName())
	long := req.GetLongitude()
	lat := req.GetLatitude()
	dist := req.GetDistance()

	if (lat < -90 || lat > 90) {
		return status.Error(codes.InvalidArgument, "Invalid value for latitude")
	} else if (long < -180 || long > 180) {
		return status.Error(codes.InvalidArgument, "Invalid value for longitude")
	} else if (dist <= 0) {
		return status.Error(codes.InvalidArgument, "Invalid value for distance")
	}

	resChan, errChan := s.db.Shop().GetNearbyShops(name, long, lat, dist)

	select {
	case shops := <- resChan:
		for _, shop := range shops {
			shopRes := &pb.SearchNearbyShopResponse{
				Shop: &pb.Shop{
					Id: shop.ShopId.Hex(),
					CatalogueId: shop.CatalogueId.Hex(),
					Name: shop.Name,
					Owner: shop.Owner,
					Distance: shop.Location.Distance(long, lat),
				},
			}

			stream.Send(shopRes)
		}

	case err := <- errChan:
		logger.Error("Error fetching nearby shops", zap.Error(err))
		return status.Error(codes.Internal, "Unexpected error occured")
	}

	return nil
}
