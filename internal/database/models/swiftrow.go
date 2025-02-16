package models

type SwiftRow struct {
	SwiftCode   string `json:"swift_code" gorm:"column:swift_code"`
	CountryCode string `json:"country_code" gorm:"column:country_code"`
	CountryName string `json:"country_name" gorm:"column:country_name"`
	BankName    string `json:"bank_name" gorm:"column:bank_name"`
	Address     string `json:"address" gorm:"column:address"`
	City        string `json:"city" gorm:"column:city"`
	TimeZone    string `json:"time_zone" gorm:"column:time_zone"`
}
