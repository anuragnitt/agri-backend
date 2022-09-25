package main

import (
	"github.com/anuragnitt/agri-backend/db"
	"github.com/anuragnitt/agri-backend/service"
)

type Inject struct {
	AgriDb						*db.AgriDb
	SearchByProductService		*service.SearchByProductService
	SearchNearbyShopService		*service.SearchNearbyShopService
}

func NewInject() *Inject {
	inj := &Inject{}
	inj.AgriDb = &db.AgriDb{}

	inj.SearchByProductService = service.NewSearchByProductService(inj.AgriDb)
	inj.SearchNearbyShopService = service.NewSearchNearbyShopService(inj.AgriDb)

	return inj
}
