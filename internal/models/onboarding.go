package models

type OnboardingPayload struct {
	CompanyName string  `json:"companyName" example:"ABC Co."`
	CompanyType string  `json:"companyType" example:"SME"`
	Province    string  `json:"province" example:"Bangkok"`
	District    string  `json:"district" example:"Pathum Wan"`
	SubDistrict string  `json:"subDistrict" example:"Lumphini"`
	Address     string  `json:"address" example:"123 ถนนสุขุมวิท"`
	Alley       *string `json:"alley" example:"Soi 5"`
	PostalCode  string  `json:"postalCode" example:"10330"`
	Tel         string  `json:"tel" example:"0999999999"`
}
