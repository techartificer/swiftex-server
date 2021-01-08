package constants

// AdminRole is the role type of admin
type AdminRole string

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

const (
	// Active status: user can access
	Active string = "Active"
	// Deactive status: user can not access
	Deactive string = "Deactive"
)
const (
	UserScope string = "user"
	Phone     string = "phone"
	UserID    string = "userId"
)
