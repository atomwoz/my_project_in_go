package model

var SwiftRow struct {
	SwiftCode   string `json:"swift_code"`
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	BankName    string `json:"bank_name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	TimeZone    string `json:"time_zone"`
}
