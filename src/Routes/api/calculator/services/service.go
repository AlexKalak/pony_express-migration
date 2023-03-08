package services

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/alexkalak/pony_express/src/Routes/validation"
	currencyhelper "github.com/alexkalak/pony_express/src/currencyHelper"
	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/helpers/city_helper"
	"github.com/alexkalak/pony_express/src/helpers/package_types_helper"
	"github.com/alexkalak/pony_express/src/models"
	"github.com/alexkalak/pony_express/src/types"
	"github.com/gofiber/fiber/v2"
)

type CalculatorService interface {
	Calculate(c *fiber.Ctx) (*ResponsePrices, bool, []*types.ErrorResponse, error)
}

type calculatorService struct {
}

func New() CalculatorService {
	return &calculatorService{}
}

type Place struct {
	Weight float64 `json:"weight" validate:"required"`
	Width  float64 `json:"width" validate:""`
	Length float64 `json:"length" validate:""`
	Height float64 `json:"height" validate:""`
}

type CalculateRequestBody struct {
	SenderRegion string  `json:"sender-city" validate:"required"`
	Places       []Place `json:"places" validate:"required,dive"`
	PackageType  string  `json:"package-type" validate:"required"`
}

type ResponsePrices struct {
	DoorDoor     float64 `json:"door-door"`
	DoorOffice   float64 `json:"door-office"`
	OfficeDoor   float64 `json:"office-door"`
	OfficeOffice float64 `json:"office-office"`
}

func (cs *calculatorService) Calculate(c *fiber.Ctx) (*ResponsePrices, bool, []*types.ErrorResponse, error) {
	database := db.GetDB()
	reqBody, validationErrors, err := handleRequest(c)
	if err != nil {
		return nil, false, nil, err
	}
	if validationErrors != nil {
		return nil, false, validationErrors, nil
	}

	//Finding sender region and city
	var senderRegionFromDB models.SenderRegion
	res := database.Preload("SenderCity").First(&senderRegionFromDB, "name = ?", reqBody.SenderRegion)
	if res.Error != nil {
		return nil, false, nil, fmt.Errorf("sender-city not found")
	}
	senderCityFromDB := senderRegionFromDB.SenderCity

	//Getting region_id
	receiverRegionID, err := getRegionID(c)
	if err != nil {
		return nil, false, nil, err
	}

	//Getitng package type from db
	packageTypeFromDB, err := package_types_helper.GetPackageTypeFromDB(reqBody.PackageType)
	if err != nil {
		return nil, false, nil, err
	}

	//Getting total weight
	weight, usedVolumeWeights, err := GetWeight(reqBody.Places, packageTypeFromDB.Name)
	if err != nil {
		return nil, false, nil, err
	}

	price, err := GetPrice(weight, receiverRegionID, packageTypeFromDB.ID, packageTypeFromDB.Name, senderCityFromDB.ID)
	if err != nil {
		return nil, false, nil, err
	}

	floatedOfficePrice := currencyhelper.ConvertIntValueToFloat(price)
	floatedDoorPrice := currencyhelper.ConvertIntValueToFloat(senderRegionFromDB.PriceForDoor)

	fmt.Println("floated door price", floatedDoorPrice)

	ResponsePrices := ResponsePrices{
		OfficeOffice: floatedOfficePrice,
		OfficeDoor:   floatedOfficePrice,
		DoorDoor:     roundFloat(floatedOfficePrice+floatedDoorPrice, 2),
		DoorOffice:   roundFloat(floatedOfficePrice+floatedDoorPrice, 2),
	}

	return &ResponsePrices, usedVolumeWeights, nil, nil
}

func handleRequest(c *fiber.Ctx) (*CalculateRequestBody, []*types.ErrorResponse, error) {
	var reqBody CalculateRequestBody
	//parsing user request
	err := c.BodyParser(&reqBody)
	if err != nil {
		return nil, nil, err
	}

	//printing req body in console in readable format
	fmt.Println("REQ BODY")
	str, _ := json.MarshalIndent(&reqBody, "", "\t")
	fmt.Println(string(str))

	//Validating user request
	validationErrors := validation.Validate(&reqBody)
	if len(reqBody.Places) < 1 {
		//If places length less than 1 returning validation error
		validationErrors = append(validationErrors, &types.ErrorResponse{
			FailedField: "places",
			Tag:         "min-length=1",
		})
	}

	if len(validationErrors) > 0 {
		return nil, validationErrors, nil
	}

	return &reqBody, nil, nil
}

func getRegionID(c *fiber.Ctx) (int, error) {
	database := db.GetDB()
	var usrInput struct {
		ReceiverCountry  string `json:"receiver-country"`
		ReceiverCity     string `json:"receiver-city"`
		ReceiverDistrict string `json:"receiver-district"`
		ReceiverArea     string `json:"receiver-area"`
	}

	err := c.BodyParser(&usrInput)
	if err != nil {
		return 0, err
	}

	str, _ := json.MarshalIndent(usrInput, "", "\t")
	fmt.Println(string(str))

	if usrInput.ReceiverCountry == "" {
		return 0, fmt.Errorf("receiver country is required")
	}

	var country *models.Country
	res := database.First(&country, "name = ?", usrInput.ReceiverCountry)
	if res.Error != nil {
		return 0, fmt.Errorf("country not found")
	}

	if country.RegionID != 0 {
		return country.RegionID, nil
	}

	if usrInput.ReceiverCity == "" {
		return 0, fmt.Errorf("country doesn't have price itself, so city name is required")
	}

	var receiverCityFromDB *models.City
	receiverCityFromDB, err = city_helper.
		GetCityByCityNameCountryDistrictAndArea(usrInput.ReceiverCity, usrInput.ReceiverCountry, usrInput.ReceiverDistrict, usrInput.ReceiverArea)

	if err != nil {
		return 0, fmt.Errorf("receiver-city not found")
	}

	//Checking if receiver area is Moscow of Saint_Petersburg
	switch receiverCityFromDB.Area.Name {
	case "Московская":
		receiverCityFromDB, err = city_helper.GetCityByName("Москва")
	case "Ленинградская":
		receiverCityFromDB, err = city_helper.GetCityByName("Санкт-Петербург")
	}
	if err != nil {
		return 0, fmt.Errorf("error occured")
	}

	return receiverCityFromDB.RegionID, nil
}

func roundFloat(number float64, precision int) float64 {
	ratio := math.Pow10(precision)
	return math.Round(number*ratio) / ratio
}
