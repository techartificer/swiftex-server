package codes

// ErrorCode help to make more sense about individual http status
type ErrorCode string

const (
	InvalidRegisterData       ErrorCode = "400001"
	UserSignUpDataInvalid     ErrorCode = "400002"
	InvalidShopCreateData     ErrorCode = "400003"
	InvalidLoginCredential    ErrorCode = "401001"
	BearerTokenGiven          ErrorCode = "401002"
	InvalidAuthorizationToken ErrorCode = "401003"
	InvalidAccountType        ErrorCode = "401004"
	StatusNotActive           ErrorCode = "403001"
	NotSuperAdmin             ErrorCode = "403002"
	AccessDenied              ErrorCode = "403003"
	AdminNotFound             ErrorCode = "404001"
	RefreshTokenNotFound      ErrorCode = "404002"
	BearerTokenNotFound       ErrorCode = "404003"
	ShopNotFound              ErrorCode = "404004"
	MerchantNotFound          ErrorCode = "404005"
	AdminAlreadyExist         ErrorCode = "409001"
	MerchantAlreadyExist      ErrorCode = "409002"
	ShopAlreadyExist          ErrorCode = "409003"
	OrderAlreadyExist         ErrorCode = "409004"
	InvalidLimit              ErrorCode = "422001"
	DatabaseQueryFailed       ErrorCode = "500001"
	UserLoginFailed           ErrorCode = "500002"
	TokenRefreshFailed        ErrorCode = "500003"
	SomethingWentWrong        ErrorCode = "500004"
	PasswordHashFailed        ErrorCode = "500005"
)
