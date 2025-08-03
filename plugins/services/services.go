package services

import (
	"explorer/app/db"
	"explorer/app/types"
	"log"
	"net/http"
	"strconv"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

func HandleServices(kit *kit.Kit) error {
	var services []types.Service
	if err := db.Get().Find(&services).Error; err != nil {
		return err
	}
	return kit.Render(ServiceConfigModal(services))
}

func HandleServiceCreate(kit *kit.Kit) error {
	if err := kit.Request.ParseForm(); err != nil {
		return err
	}

	name := kit.Request.FormValue("name")
	priceStr := kit.Request.FormValue("price")

	if name == "" {
		http.Error(kit.Response, "Name is required", http.StatusBadRequest)
		return nil
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		http.Error(kit.Response, "Invalid price", http.StatusBadRequest)
		return nil
	}

	service := types.Service{Name: name, Price: price}
	if err := db.Get().Create(&service).Error; err != nil {
		return err
	}

	var services []types.Service
	if err := db.Get().Find(&services).Error; err != nil {
		return err
	}

	return kit.Render(ServiceConfigModal(services))
}

func HandleServiceDelete(kit *kit.Kit) error {
	idStr := chi.URLParam(kit.Request, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := db.Get().Delete(&types.Service{}, id).Error; err != nil {
		log.Println(err)
		return err
	}

	var services []types.Service
	if err := db.Get().Find(&services).Error; err != nil {
		log.Println(err)
		return err
	}

	return kit.Render(ServiceConfigModal(services))
}
