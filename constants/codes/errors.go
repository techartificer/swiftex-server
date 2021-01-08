package codes

// ErrorCode help to make more sense about individual http status
type ErrorCode string

const (
	InvalidRegisterData       ErrorCode = "400001"
	UserSignUpDataInvalid     ErrorCode = "400002"
	InvalidLoginCredential    ErrorCode = "401001"
	BearerTokenGiven          ErrorCode = "401002"
	InvalidAuthorizationToken ErrorCode = "401003"
	StatusNotActive           ErrorCode = "403001"
	NotSuperAdmin             ErrorCode = "403002"
	AdminNotFound             ErrorCode = "404001"
	RefreshTokenNotFound      ErrorCode = "404002"
	BearerTokenNotFound       ErrorCode = "404003"
	AdminAlreadyExist         ErrorCode = "409001"
	DatabaseQueryFailed       ErrorCode = "500001"
	UserLoginFailed           ErrorCode = "500002"
	TokenRefreshFailed        ErrorCode = "500003"
	SomethingWentWrong        ErrorCode = "500004"
	PasswordHashFailed        ErrorCode = "500005"
)
