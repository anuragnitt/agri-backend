package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductQuantity struct {
	ProductId	primitive.ObjectID	`bson:"productId" json:"productId"`
	Quantity	uint32				`bson:"quantity" json:"quantity"`
}

type Catalogue struct {
	CatalogueId	primitive.ObjectID	`bson:"_id" json:"_id"`
	Inventory	[]ProductQuantity	`bson:"inventory" json:"inventory"`
}
