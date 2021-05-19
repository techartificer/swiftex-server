package charge

import (
	"math"
	"strings"

	"github.com/techartificer/swiftex/constants"
)

const (
	Dhaka string = "Dhaka"
)

func Calculate(weight float32, deliverType string, city string, deliveryCharge float64) float64 {
	if deliveryCharge == 0 {
		deliveryCharge = constants.DeliveryCharge
	}
	cweight := math.Ceil(float64(weight)) - 1
	charge := deliveryCharge + (20 * cweight)

	isInsideDhaka := true
	if !strings.EqualFold(city, Dhaka) {
		charge += 70
		isInsideDhaka = false
	}
	if strings.EqualFold(deliverType, constants.Express) && isInsideDhaka {
		charge += 40
	}
	return charge
}
