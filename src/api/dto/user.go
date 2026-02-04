package dto

type GetOtpRequest struct {
	CountryCode string `json:"countryCode"`
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,min=1`
}

type RegisterLoginByMobileRequest struct {
	CountryCode string `json:"countryCode"`
	MobileNumber string `json:"mobileNumber" binding:"required,mobile"`
	Otp          string `json:"otp" binding:"required"`
}
