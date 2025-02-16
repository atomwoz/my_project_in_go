package models

import "github.com/spf13/viper"

type SwiftModel struct {
	SwiftCode      string `json:"swiftCode" gorm:"column:swift_code;primaryKey"`
	CountryCode    string `json:"countryISO2" gorm:"column:country_code"`
	CountryName    string `json:"countryName" gorm:"column:country_name"`
	BankName       string `json:"bankName" gorm:"column:bank_name"`
	Address        string `json:"address" gorm:"column:address"`
	City           string `json:"city" gorm:"column:city"`
	TimeZone       string `json:"timeZone" gorm:"column:time_zone"`
	BankSymbol     string `json:"bankSymbol" gorm:"column:bank_symbol"`
	IsHeadquarters bool   `json:"isHeadquarter" gorm:"column:is_headquarters"`
}

func (SwiftModel) TableName() string {
	return viper.GetString("db.table")
}

type SwiftHeadquarterRow struct {
	BankName      string           `json:"bankName" gorm:"column:bank_name"`
	Address       string           `json:"address" gorm:"column:address"`
	CountryCode   string           `json:"countryISO2" gorm:"column:country_code"`
	CountryName   string           `json:"countryName" gorm:"column:country_name"`
	IsHeadquarter bool             `json:"isHeadquarter" gorm:"column:is_headquarter"`
	SwiftCode     string           `json:"swiftCode" gorm:"column:swift_code"`
	Branches      []SwiftBranchRow `json:"branches" gorm:"foreignKey:BankName;references:BankName"`
	bankSymbol    string           `json:"bankSymbol" gorm:"column:bank_symbol"`
}

type SwiftBranchRow struct {
	BankName      string `json:"bankName" gorm:"column:bank_name"`
	Address       string `json:"address" gorm:"column:address"`
	CountryCode   string `json:"countryISO2" gorm:"column:country_code"`
	CountryName   string `json:"countryName" gorm:"column:country_name"`
	IsHeadquarter bool   `json:"isHeadquarter" gorm:"column:is_headquarter"`
	SwiftCode     string `json:"swiftCode" gorm:"column:swift_code"`
}
