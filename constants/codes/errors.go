package codes

// ErrorCode help to make more sense about individual http status
type ErrorCode string

const (
	InvalidRegisterData    ErrorCode = "400001"
	InvalidLoginCredential ErrorCode = "401001"
	AdminNotFound          ErrorCode = "404001"
	DatabaseQueryFailed    ErrorCode = "500001"
	UserLoginFailed        ErrorCode = "500002"
)
