package services

import (
	"explorer/app/db"
	"explorer/app/types"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi"
)

type ServiceFromValues struct {
	Names  []string `form:"names"`
	Prices []string `form:"prices"`
}

type ServiceRow struct {
	Name  string
	Price string
	Error string
}

type ServicesPageData struct {
	Rows    []ServiceRow
	Success string
}

func HandleServices(kit *kit.Kit) error {
	var services []types.Service
	if err := db.Get().Find(&services).Error; err != nil {
		return err
	}
	return kit.Render(ServicesForm(services))
}

func HandleAddService(kit *kit.Kit) error {
	name := kit.Request.FormValue("name")
	priceStr := kit.Request.FormValue("price")

	if name == "" || priceStr == "" {
		if name == "" {
			return fmt.Errorf("name is required")
		}
		return fmt.Errorf("price is required")
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fmt.Errorf("invalid price format")
	}

	s := types.Service{Name: name, Price: price}
	if err := db.Get().Create(&s).Error; err != nil {
		return err
	}
	var allservices []types.Service
	if err := db.Get().Find(&allservices).Error; err != nil {
		return err
	}

	return kit.Render(ServicesForm(allservices))
}

func HandleDeleteService(kit *kit.Kit) error {
	log.Println("Request path:", kit.Request.URL.Path)

	// Extract id from URL
	idStr := chi.URLParam(kit.Request, "id")
	log.Println("ID from URL param:", idStr)
	for k, v := range chi.RouteContext(kit.Request.Context()).URLParams.Keys {
		log.Println("Route param key:", k, "=>", v)
	}

	if idStr == "" {
		return fmt.Errorf("id is required")
	}

	// Convert string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id, %s", err)
	}

	// Check if service exists
	var service types.Service
	if err := db.Get().First(&service, id).Error; err != nil {
		return kit.Text(http.StatusBadRequest, "Service not found")
	}

	// Delete the service
	if err := db.Get().Delete(&service).Error; err != nil {
		return kit.Text(http.StatusInternalServerError, "Failed to delete service")
	}

	// Return updated list
	var services []types.Service
	if err := db.Get().Find(&services).Error; err != nil {
		return err
	}

	return kit.Render(ServicesForm(services))
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
