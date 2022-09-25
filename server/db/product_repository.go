package db

import (
	"fmt"

	"github.com/anuragnitt/agri-backend/models"
	"github.com/anuragnitt/agri-backend/odm"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductRepository struct {
	odm.AbstractRepository[models.Product]
}

func (p *ProductRepository) GetProductsByName(name string) (chan []models.Product, chan error) {
	ch := make(chan []models.Product)
	errorChan := make(chan error)

	go func() {
		filters := bson.M{}
		if len(name) > 0 {
			filters["name"] = bson.M{
				"$regex": fmt.Sprintf(".*%v.*", name),
				"$options": "i",
			}
		}

		resChan, errChan := p.Find(filters, bson.D{}, 0, 0)

		select {
		case products := <- resChan:
			ch <- products
		case err := <- errChan:
			errorChan <- err
			return
		}
	}()

	return ch, errorChan
}

func (p *ProductRepository) GetProductsByAttr(name, category string, minPrice, maxPrice uint32) (chan []models.Product, chan error) {
	ch := make(chan []models.Product)
	errorChan := make(chan error)

	go func() {
		filters := bson.M{
			"price": bson.M{
				"$gte": minPrice,
				"$lte": maxPrice,
			},
		}
		if len(name) > 0 {
			filters["name"] = bson.M{
				"$regex": fmt.Sprintf(".*%v.*", name),
				"$options": "i",
			}
		}
		if len(category) > 0 {
			filters["category"] = bson.M{
				"$regex": fmt.Sprintf(".*%v.*", category),
				"$options": "i",
			}
		}

		resChan, errChan := p.Find(filters, bson.D{}, 0, 0)

		select {
		case products := <- resChan:
			ch <- products
		case err := <- errChan:
			errorChan <- err
			return
		}
	}()

	return ch, errorChan
}
