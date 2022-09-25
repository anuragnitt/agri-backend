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

type SearchByProductService struct {
	pb.UnimplementedSearchByProductServiceServer
	db *db.AgriDb
}

func NewSearchByProductService(db *db.AgriDb) *SearchByProductService {
	return &SearchByProductService{
		db: db,
	}
}

func (s *SearchByProductService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

func (s *SearchByProductService) SearchByProduct(req *pb.SearchByProductRequest, stream pb.SearchByProductService_SearchByProductServer) error {
	name := strings.TrimSpace(req.GetName())

	resChan, errChan := s.db.Product().GetProductsByName(name)

	select {
	case products := <- resChan:
		for _, product := range products {
			shopChan, errorChan := s.db.FindShopByProductId(product.ProductId, 1)

			select {
			case shop := <- shopChan:
				productRes := &pb.SearchByProductResponse{
					Product: &pb.Product{
						Id: product.ProductId.Hex(),
						ShopId: shop.ShopId.Hex(),
						Name: product.Name,
						Price: product.Price,
						Category: product.Category,
					},
				}

				stream.Send(productRes)

			case err := <- errorChan:
				logger.Error("Error fetching shop by productId", zap.Error(err))
				return status.Error(codes.Internal, "Unexpected error occured")
			}
		}

	case err := <- errChan:
		logger.Error("Error fetching products by name", zap.Error(err))
		return status.Error(codes.Internal, "Unexpected error occured")
	}

	return nil
}


func (s *SearchByProductService) SearchByProductAttr(req *pb.SearchByProductAttrRequest, stream pb.SearchByProductService_SearchByProductAttrServer) error {
	name := strings.TrimSpace(req.GetName())
	minPrice := req.GetMinPrice()
	maxPrice := req.GetMaxPrice()
	category := strings.TrimSpace(req.GetCategory())

	if maxPrice <= 0 {
		return status.Error(codes.InvalidArgument, "Invalid value for maxPrice")
	}

	resChan, errChan := s.db.Product().GetProductsByAttr(name, category, minPrice, maxPrice)

	select {
	case products := <- resChan:
		for _, product := range products {
			shopChan, errorChan := s.db.FindShopByProductId(product.ProductId, 1)

			select {
			case shop := <- shopChan:
				productRes := &pb.SearchByProductResponse{
					Product: &pb.Product{
						Id: product.ProductId.Hex(),
						ShopId: shop.ShopId.Hex(),
						Name: product.Name,
						Price: product.Price,
						Category: product.Category,
					},
				}

				stream.Send(productRes)

			case err := <- errorChan:
				logger.Error("Error fetching shop by productId", zap.Error(err))
				return status.Error(codes.Internal, "Unexpected error occured")
			}
		}

	case err := <- errChan:
		logger.Error("Error fetching products by attributes", zap.Error(err))
		return status.Error(codes.Internal, "Unexpected error occured")
	}

	return nil
}
