package currencyhelper

func ConvertIntValueToFloat(priceInt int) float64 {
	return float64(priceInt) / 100
}
