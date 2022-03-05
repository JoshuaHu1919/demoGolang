package errors

const (
	// Success 指成功
	Success = "200001"

	// SuccessNoContent 代表成功同時不需要回傳的內容
	SuccessNoContent = "204001"

	// InvalidAuthenticationInfo 代表認證資訊錯誤
	InvalidAuthenticationInfo = "400003"

	// InvalidHeaderValue - The value provided for one of the HTTP headers was not in the correct format.")]
	InvalidHeaderValue = "400004"

	// InvalidInput - One of the request inputs is not valid.")]
	InvalidInput = "400006"

	// InvalidQueryParameterValue - An invalid value was specified for one of the query parameters in the request URI.")]
	InvalidQueryParameterValue = "400009"

	// OutOfRangeInput - One of the request inputs is out of range.")]
	OutOfRangeInput = "400020"

	// Unauthorized 指未授權
	Unauthorized = "401001"

	// AccountIsDisabled - The specified account is disabled." )]
	AccountIsDisabled = "403001"

	// NotAllowed - The request is understood, but it has been refused or access is not allowed.")]
	NotAllowed = "403003"

	// UsernameOrPasswordIncorrect - Username or Password is incorrect
	UsernameOrPasswordIncorrect = "403006"

	// OtpRequired - OTP Binding is required.
	OtpRequired = "403007"

	// OtpAuthorizationRequired - Two-factor authorization is required
	OtpAuthorizationRequired = "403008"

	// OtpIncorrect - OTP is incorrect
	OtpIncorrect = "403009"

	// ResetPasswordRequired - Reset Password Required
	ResetPasswordRequired = "403010"

	// ResourceNotFound - The specified resource does not exist.
	ResourceNotFound = "404001"

	// ResourceDependencyNotFound - The specified resource dependency does not exist
	ResourceDependencyNotFound = "404002"

	// AccountAlreadyExists - The specified account already exists.
	AccountAlreadyExists = "409001"

	// ResourceAlreadyExists - Conflict (409) - The specified resource already exists.")]
	ResourceAlreadyExists = "409004"

	// InternalError - "Internal Server Error (500) - The server encountered an internal error. Please retry the request.
	InternalError = "500001"
)
