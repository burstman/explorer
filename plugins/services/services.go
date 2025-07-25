package services

import (
	"explorer/app/db"
	"explorer/app/types"
	"fmt"
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
	fmt.Println("test")
	idStr := chi.URLParam(kit.Request, "id")
	if len(idStr) == 0 {
		log.Println("id is null")
		return fmt.Errorf("id is null")
	}

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

// func validateServices(values ServiceFromValues) v.Errors {
// 	errors := v.Errors{}
// 	for i, price := range values.Prices {
// 		key := fmt.Sprintf("prices[%d]", i)
// 		if price == "" {
// 			errors.Add(key, "Price is required.")
// 		} else if _, err := strconv.ParseFloat(price, 64); err != nil {
// 			errors.Add(key, "Price must be a valid number.")
// 		}
// 	}
// 	return errors
// }

// func HandleSaveServices(kit kit.Kit) error {
// 	var book types.Booking
// 	var values ServiceFromValues
// 	errors := v.Errors{}
// 	if err := kit.Request.ParseForm(); err != nil {
// 		kit.Text(http.StatusBadRequest, "Invalid form data")
// 	}
// 	v.Request(kit.Request, &values, nil)

// 	for i, strPrice := range values.Prices {
// 		price, err := strconv.ParseFloat(strPrice, 64)
// 		if err != nil {
// 			errors.Add("prices", "Invalid price")
// 		}
// 		values.Prices[i] = strconv.FormatFloat(price, 'f', 2, 64)
// 	}

// 	err := db.Get().Save(&values).Error
// 	if err != nil {
// 		return kit.Text(http.StatusBadRequest, "Failed to save services")
// 	}

// 	if errors.Any() {
// 		return kit.Render(ServicesForm(ServicesPageData{
// 			Values:  values,
// 			Errors:  errors,
// 			Success: "Services saved successfully",
// 		}))
// 	}
// }
