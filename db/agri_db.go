package db

import (
	"os"

	"github.com/anuragnitt/agri-backend/models"
	"github.com/anuragnitt/agri-backend/odm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgriDb struct {}

func (a *AgriDb) Shop() *ShopRepository {
	baseRepo := odm.AbstractRepository[models.Shop]{
		Database: os.Getenv("AGRI_DB"),
		CollectionName: os.Getenv("SHOP_COLLECTION"),
	}

	return &ShopRepository{baseRepo}
}

func (a *AgriDb) Product() *ProductRepository {
	baseRepo := odm.AbstractRepository[models.Product]{
		Database: os.Getenv("AGRI_DB"),
		CollectionName: os.Getenv("PRODUCT_COLLECTION"),
	}

	return &ProductRepository{baseRepo}
}

func (a *AgriDb) Catalogue() *CatalogueRepository {
	baseRepo := odm.AbstractRepository[models.Catalogue]{
		Database: os.Getenv("AGRI_DB"),
		CollectionName: os.Getenv("CTALOGUE_COLLECTION"),
	}

	return &CatalogueRepository{baseRepo}
}

func (a *AgriDb) FindShopByProductId(productId primitive.ObjectID, minQuantity uint32) (chan *models.Shop, chan error) {
	ch := make(chan *models.Shop)
	errorChan := make(chan error)

	go func() {
		catalogueChan, errChan_1 := a.Catalogue().FindCatalogueByProductId(productId, minQuantity)

		select {
		case catalogue := <- catalogueChan:
			shopChan, errChan_2 := a.Shop().FindShopByCatalogueId(catalogue.CatalogueId)

			select {
			case shop := <- shopChan:
				ch <- shop
			case err := <- errChan_2:
				errorChan <- err
				return
			}

		case err := <- errChan_1:
			errorChan <- err
			return
		}
	}()

	return ch, errorChan
}
