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
	AdminNotFound             ErrorCode = "404001"
	RefreshTokenNotFound      ErrorCode = "404002"
	BearerTokenNotFound       ErrorCode = "404003"
	DatabaseQueryFailed       ErrorCode = "500001"
	UserLoginFailed           ErrorCode = "500002"
	TokenRefreshFailed        ErrorCode = "500003"
)
