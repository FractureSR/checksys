package apratus

import (
	"math"
)

func SubString(mainString string, startIndex int, endIndex int) string {
	rs := []rune(mainString[startIndex:endIndex])
	return string(rs)
}

func HarvenSin(co1 Coordinate, co2 Coordinate) (float64, int) {
	radius := 6371000.0
	rad := math.Pi / 180.0

	co1.Lat *= rad
	co1.Lng *= rad
	co2.Lat *= rad
	co2.Lng *= rad

	Lngtheta := math.Abs(co1.Lng - co2.Lng)
	Lattheta := math.Abs(co1.Lat - co2.Lat)
	stepOne := math.Pow(math.Sin(Lattheta/2), 2) + math.Cos(co1.Lat)*math.Cos(co2.Lat)*math.Pow(math.Sin(Lngtheta/2), 2)
	dist := 2 * math.Asin(math.Min(1, math.Sqrt(stepOne)))

	return dist * radius, co1.Accuracy + co2.Accuracy
}
