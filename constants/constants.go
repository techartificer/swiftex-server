package constants

// AdminRole is the role type of admin
type AdminRole string

const Version = "v1.0.6 beta"

const (
	AdminType    string = "Admin"
	MerchantType string = "Merchant"
	RiderType    string = "Rider"
)

const (
	ShopOwner     string = "Owner"
	ShopModerator string = "Moderator"
)

const (
	// SuperAdmin has all sorts of access
	SuperAdmin AdminRole = "Super Admin"
	// Admin has all access exepet admin creation
	Admin AdminRole = "Admin"
	// Moderator has all access except admin creation and payment
	Moderator AdminRole = "Moderator"
	// ZoneManager has zone centric access
	ZoneManager AdminRole = "Zone Manager"
	Rider       string    = "Rider"
)

var Roles = []AdminRole{SuperAdmin, Admin, Moderator, ZoneManager}

const (
	Active   string = "Active"
	Deactive string = "Deactive"
)

const (
	Cancelled   string = "Cancelled"
	Created     string = "Created"
	Delivered   string = "Delivered"
	Accepted    string = "Accepted"
	Assigned    string = "Assigned"
	Apporved    string = "Approved"
	Pending     string = "Pending"
	Declined    string = "Declined"
	InTransit   string = "In Transit"
	Returned    string = "Returned"
	Rescheduled string = "Rescheduled"
	Picked      string = "Picked"
)

var AllStatus = []string{Active, Deactive}

const (
	COD  string = "COD"
	PAID string = "PAID"
)

const (
	Role   string = "role"
	Phone  string = "phone"
	UserID string = "userId"
)

const TrackIDSize = 8

const (
	Express string = "Express"
	Regular string = "Regular"
)

const (
	CreatedMsg    string = "Your parcel has been placed"
	AcceptedMsg   string = "Parcel has been accepted"
	PickedMsg     string = "Parcel has been picked up"
	InTransitMsg  string = "Out for delivery"
	CancelledMsg  string = "Parcel has been cancelled"
	ReturnedMsg   string = "Parcel has been returned"
	RescheduleMsg string = "Parcel has been rescheduled"
	DeleveredMsg  string = "Successfully delevered at your door"
)

const (
	DeliveryCharge float64 = 60
	CodCharge      float64 = 1
)
