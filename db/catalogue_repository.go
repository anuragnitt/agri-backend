package db

import (
	"github.com/anuragnitt/agri-backend/models"
	"github.com/anuragnitt/agri-backend/odm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CatalogueRepository struct {
	odm.AbstractRepository[models.Catalogue]
}

func (c *CatalogueRepository) FindCatalogueByProductId(productId primitive.ObjectID, minQuantity uint32) (chan *models.Catalogue, chan error) {
	ch := make(chan *models.Catalogue)
	errorChan := make(chan error)

	go func() {
		resChan, errChan := c.FindOne(bson.M{
			"inventory": bson.M{
				"$elemMatch": bson.M{
					"productId": productId,
					"quantity": bson.M{
						"$gte": minQuantity,
					},
				},
			},
		})

		select {
		case catalogue := <- resChan:
			ch <- catalogue
		case err := <- errChan:
			errorChan <- err
			return
		}
	}()

	return ch, errorChan
}
