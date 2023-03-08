package services

import (
	"errors"
	"fmt"
	"sort"

	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/models"
)

func GetWeight(places []Place, packageType string) (float64, bool, error) {
	weights, usedVolumeWeights, err := GetWeights(places, packageType)
	if err != nil {
		return 0, false, err
	}

	var sum float64
	for _, weight := range weights {
		sum += weight
	}

	fmt.Println(weights)
	return sum, usedVolumeWeights, nil
}

// Get array of weights
func GetWeights(places []Place, packageType string) ([]float64, bool, error) {
	var err error
	var weights = make([]float64, 0, len(places))

	//If user uses documents package type, will be checked only real weights
	if packageType == "documents" {
		for _, place := range places {
			weights = append(weights, place.Weight)
		}
		return weights, false, nil
	}

	//If user uses standart or B2B package type, will be checked volume weights too
	weights, usedVolumeWeights, err := GetMaxWeights(places)
	if err != nil {
		return nil, false, err
	}

	return weights, usedVolumeWeights, nil
}

// Checking volume weights
func GetMaxWeights(places []Place) ([]float64, bool, error) {
	maxWeights := make([]float64, 0, len(places))
	usedVolumeWeights := false

	for _, place := range places {
		if place.Length == 0 || place.Weight == 0 || place.Height == 0 {
			return nil, false, errors.New("invalid sizes")
		}

		volumeWeight := CalculateVolumeWeight(place.Length, place.Width, place.Height)

		if volumeWeight > place.Weight {
			maxWeights = append(maxWeights, volumeWeight)
			usedVolumeWeights = true
			continue
		}

		maxWeights = append(maxWeights, place.Weight)
	}

	return maxWeights, usedVolumeWeights, nil
}

// rounding weight
func RoundUpWeight(weight *float64, weights []float64, maxWeight float64) {
	sort.Float64s(weights)

	numOfMaxWeights := int(*weight / maxWeight)
	reminder := *weight - float64(numOfMaxWeights)*maxWeight

	fmt.Println("reminder in rounding: ", reminder)
	if reminder == 0 {
		return
	}

	for i := range weights {
		if weights[i] >= reminder {
			reminder = weights[i]
			*weight = float64(numOfMaxWeights)*maxWeight + reminder
			fmt.Println("Rounded weight: ", *weight)
			return
		}
	}
}

// Finding in DB all weights that can be userd for package
func GetAllWeightsOfPackage(regionID int, packageTypeID int) ([]float64, error) {
	database := db.GetDB()
	var prices []models.Price
	res := database.Preload("Weight").Find(&prices, "region_id = ? AND package_type_id = ?", regionID, packageTypeID)

	if res.Error != nil {
		return nil, res.Error
	}

	weights := make([]float64, len(prices))
	for i := range prices {
		weights[i] = prices[i].Weight.Weight
	}

	return weights, nil
}

func GetMaxWeightFromArray(weights []float64) float64 {
	max := weights[0]
	for _, weight := range weights {
		if weight > max {
			max = weight
		}
	}
	fmt.Println("Max:", max)
	return max
}

func CalculateVolumeWeight(length float64, width float64, height float64) float64 {
	return length * width * height / 5000
}
