package db

import (
	"fmt"

	"github.com/anuragnitt/agri-backend/models"
	"github.com/anuragnitt/agri-backend/odm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopRepository struct {
	odm.AbstractRepository[models.Shop]
}

func (s *ShopRepository) GetNearbyShops(name string, longitude, latitude, maxDistance float64) (chan []models.Shop, chan error) {
	ch := make(chan []models.Shop)
	errorChan := make(chan error)

	go func() {
		filters := bson.M{
			"location": bson.M{
				"$nearSphere": bson.M{
					"$geometry": bson.M{
						"type": "Point",
						"coordinates": []float64{longitude, latitude},
					},
					"$maxDistance": maxDistance,
				},
			},
		}
		if len(name) > 0 {
			filters["name"] = bson.M{
				"$regex": fmt.Sprintf(".*%v.*", name),
				"$options": "i",
			}
		}

		resChan, errChan := s.Find(filters, bson.D{}, 0, 0)

		select {
		case shops := <- resChan:
			ch <- shops
		case err := <- errChan:
			errorChan <- err
			return
		}
	}()

	return ch, errorChan
}

func (s *ShopRepository) FindShopByCatalogueId(catalogueId primitive.ObjectID) (chan *models.Shop, chan error) {
	ch := make(chan *models.Shop)
	errorChan := make(chan error)

	go func() {
		resChan, errChan := s.FindOne(bson.M{
			"catalogueId": catalogueId,
		})

		select {
		case shop := <- resChan:
			ch <- shop
		case err := <- errChan:
			errorChan <- err
			return
		}
	}()

	return ch, errorChan
}
