package charge

import (
	"math"
	"strings"

	"github.com/techartificer/swiftex/constants"
)

const (
	Dhaka string = "Dhaka"
)

func Calculate(weight float32, deliverType string, city string) float64 {
	var charge float64 = 60.0

	cweight := math.Ceil(float64(weight)) - 1
	charge += (20 * cweight)

	isInsideDhaka := true
	if strings.ToLower(city) != strings.ToLower(Dhaka) {
		charge += 70
		isInsideDhaka = false
	}
	if strings.ToLower(deliverType) == strings.ToLower(constants.Express) && isInsideDhaka {
		charge += 40
	}
	return charge
}
