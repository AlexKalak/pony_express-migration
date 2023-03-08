package currencyhelper

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

type Currency struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

type Data struct {
	TRY Currency `json:"TRY"`
}

type ApiResponse struct {
	Data Data `json:"data"`
}

const (
	apiCurrencyFloating = 1_000_000
)

var TRYPrice int

func StartGettingCurrencies(delay time.Duration) {
	for {
		client := &http.Client{}

		access_key := "3UvsTADx0A39OXkbiDTx5aigVSxrE0w9xGaVOfYK"

		url := fmt.Sprintf("https://api.currencyapi.com/v3/latest/?apikey=%s&currencies=TRY", access_key)

		req, _ := http.NewRequest("GET", url, nil)

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(bytes))

		var currencyInf ApiResponse
		err = json.Unmarshal(bytes, &currencyInf)
		if err != nil {
			panic(err)
		}

		TRYPrice = int(currencyInf.Data.TRY.Value * apiCurrencyFloating)
		fmt.Println(TRYPrice)
		time.Sleep(delay)
	}
}

func ConvertTRYtoUSD(tryCount int) int {
	f_price := float64(tryCount) / (float64(TRYPrice) / apiCurrencyFloating)
	fmt.Println("f_price:	", f_price)
	i_price := int(math.Round(f_price))
	fmt.Println("i_price:	", i_price)
	return i_price
}

func ConvertIntValueToFloat(priceInt int) float64 {
	return float64(priceInt) / 100
}
