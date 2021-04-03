package charge

import (
	"strings"

	"github.com/techartificer/swiftex/constants"
)

const (
	Dhaka string = "Dhaka"
)

func Calculate(weight float32, deliverType string, city string) float64 {
	var charge float64 = 60.0
	if weight > 1.00 && weight <= 2.00 {
		charge += 20
	} else if weight > 2.00 {
		charge += 40
	}
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
