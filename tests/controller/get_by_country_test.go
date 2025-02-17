package controller_tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"atomwoz.com/remitly_task/internal/database"
	rtr "atomwoz.com/remitly_task/internal/router"

	"github.com/magiconair/properties/assert"
)

var router_country = rtr.CreateRouter("/v1")

func TestGetByCountryCode(t *testing.T) {
	database.SetupTestDatabase()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/country/AW", nil)
	router_country.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{
    "countryISO2": "AW",
    "countryName": "ARUBA",
    "Branches": [
        {
            "bankName": "AIB BANK NV",
            "address": "WILHELMINASTRAAT 36  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST",
            "countryISO2": "AW",
            "countryName": "ARUBA",
            "isHeadquarter": true,
            "swiftCode": "ANIBAWA1XXX"
        },
        {
            "bankName": "ARUBA BANK, LTD",
            "address": "CAMACURI 12  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST",
            "countryISO2": "AW",
            "countryName": "ARUBA",
            "isHeadquarter": true,
            "swiftCode": "ARUBAWAXXXX"
        },
        {
            "bankName": "CENTRALE BANK VAN ARUBA",
            "address": "J.E. IRAUSQUIN 8  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST",
            "countryISO2": "AW",
            "countryName": "ARUBA",
            "isHeadquarter": true,
            "swiftCode": "CBARAWAWXXX"
        },
        {
            "bankName": "CARIBBEAN MERCANTILE BANK N.V.",
            "address": "KAYA GILBERTO FRANCOIS CROES 53 ORANJESTAD, ORANJESTAD-WEST AND ORANJESTAD-EAST",
            "countryISO2": "AW",
            "countryName": "ARUBA",
            "isHeadquarter": true,
            "swiftCode": "CMBAAWAXXXX"
        },
        {
            "bankName": "INTERBANK ARUBA NV",
            "address": "CAYA G.F. CROES 38  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST",
            "countryISO2": "AW",
            "countryName": "ARUBA",
            "isHeadquarter": true,
            "swiftCode": "IARUAWA1XXX"
        },
        {
            "bankName": "IMTRADEX INTERNATIONAL N.V.",
            "address": "TANKI LENDEERT 143  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST",
            "countryISO2": "AW",
            "countryName": "ARUBA",
            "isHeadquarter": true,
            "swiftCode": "IMIEAWA1XXX"
        },
        {
            "bankName": "RBC ROYAL BANK (ARUBA) N.V. (FORMERLY RBTT BANK ARUBA N.V.)",
            "address": "ITALIESTRAAT 36  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST",
            "countryISO2": "AW",
            "countryName": "ARUBA",
            "isHeadquarter": true,
            "swiftCode": "RBTTAWAWXXX"
        },
        {
            "bankName": "BANCO DI CARIBE (ARUBA) N.V",
            "address": "VONDELLAAN 31  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST",
            "countryISO2": "AW",
            "countryName": "ARUBA",
            "isHeadquarter": true,
            "swiftCode": "BDCCAWAWXXX"
        }
    ]
}`, w.Body.String())
}

func TestGetNoCountry(t *testing.T) {
	database.SetupTestDatabase()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/country", nil)
	router_country.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, w.Body.String(), `{"error":2,"error_msg":"Swift code 'country' not found"}`)
}

func TestWrongCountry(t *testing.T) {
	database.SetupTestDatabase()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/country/XX", nil)
	router_country.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, w.Body.String(), `{"error":3,"error_msg":"Country ISO2 code 'XX' not found"}`)
}
