package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"atomwoz.com/remitly_task/internal/database"
	rtr "atomwoz.com/remitly_task/internal/router"

	"github.com/magiconair/properties/assert"
)

var router = rtr.CreateRouter("/v1")

func TestGetBySwiftCodeBranch(t *testing.T) {
	database.SetupTestDatabase("../../../config")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/ALBPPLP1BMW", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{
    "bankName": "ALIOR BANK SPOLKA AKCYJNA",
    "address": "WARSZAWA, MAZOWIECKIE",
    "countryISO2": "PL",
    "countryName": "POLAND",
    "isHeadquarter": false,
    "swiftCode": "ALBPPLP1BMW"
}`, w.Body.String())
}
func TestGetBySwiftHQ(t *testing.T) {
	database.SetupTestDatabase("../../../config")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/BKSACLRMXXX", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{
    "bankName": "SCOTIABANK CHILE",
    "address": "AVENIDA COSTANERA SUR 2710, FLOOR 10 EDIFICIO PARQUE TITANIUM SANTIAGO, PROVINCIA DE SANTIAGO",
    "countryISO2": "CL",
    "countryName": "CHILE",
    "isHeadquarter": true,
    "swiftCode": "BKSACLRMXXX",
    "branches": [
        {
            "bankName": "SCOTIABANK CHILE",
            "address": "",
            "countryISO2": "CL",
            "countryName": "CHILE",
            "isHeadquarter": false,
            "swiftCode": "BKSACLRM055"
        },
        {
            "bankName": "SCOTIABANK CHILE",
            "address": "",
            "countryISO2": "CL",
            "countryName": "CHILE",
            "isHeadquarter": false,
            "swiftCode": "BKSACLRM061"
        },
        {
            "bankName": "SCOTIABANK CHILE",
            "address": "21 DE MAYO 187  ARICA, PROVINCIA DE ARICA, 1000000",
            "countryISO2": "CL",
            "countryName": "CHILE",
            "isHeadquarter": false,
            "swiftCode": "BKSACLRM064"
        },
        {
            "bankName": "SCOTIABANK CHILE",
            "address": "",
            "countryISO2": "CL",
            "countryName": "CHILE",
            "isHeadquarter": false,
            "swiftCode": "BKSACLRM068"
        }
    ]
}`, w.Body.String())
}

func TestGetNoSwift(t *testing.T) {
	database.SetupTestDatabase("../../../config")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, w.Body.String(), "404 page not found")
}

func TestWrongSwift(t *testing.T) {
	database.SetupTestDatabase("../../../config")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/ALA_MA_KOTA", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, w.Body.String(), `{"error":2,"error_msg":"Swift code 'ALA_MA_KOTA' not found"}`)
}
