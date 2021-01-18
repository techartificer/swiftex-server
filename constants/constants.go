package constants

// AdminRole is the role type of admin
type AdminRole string

const (
	AdminType    string = "Admin"
	MerchantType string = "Merchant"
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
)

var Roles = []AdminRole{SuperAdmin, Admin, Moderator, ZoneManager}

const (
	// Active status: user can access
	Active string = "Active"
	// Deactive status: user can not access
	Deactive  string = "Deactive"
	Pending   string = "Pending"
	Delevered string = "Delevered"
)

var AllStatus = []string{Active, Deactive}

const (
	Role   string = "role"
	Phone  string = "phone"
	UserID string = "userId"
)

const TrackIDSize = 8
