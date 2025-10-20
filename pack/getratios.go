package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var supported_currencies map[string]bool = map[string]bool{"AED": true,
	"AFN": true, "ALL": true, "AMD": true, "ANG": true, "AOA": true, "ARS": true,
	"AUD": true, "AWG": true, "AZN": true, "BAM": true, "BBD": true, "BDT": true,
	"BGN": true, "BHD": true, "BIF": true, "BMD": true, "BND": true, "BOB": true,
	"BRL": true, "BSD": true, "BTN": true, "BWP": true, "BYN": true, "BZD": true,
	"CAD": true, "CDF": true, "CHF": true, "CLP": true, "CNY": true, "COP": true,
	"CRC": true, "CUP": true, "CVE": true, "CZK": true, "DJF": true, "DKK": true,
	"DOP": true, "DZD": true, "EGP": true, "ERN": true, "ETB": true, "EUR": true,
	"FJD": true, "FKP": true, "FOK": true, "GBP": true, "GEL": true, "GGP": true,
	"GHS": true, "GIP": true, "GMD": true, "GNF": true, "GTQ": true, "GYD": true,
	"HKD": true, "HNL": true, "HRK": true, "HTG": true, "HUF": true, "IDR": true,
	"ILS": true, "IMP": true, "INR": true, "IQD": true, "IRR": true, "ISK": true,
	"JEP": true, "JMD": true, "JOD": true, "JPY": true, "KES": true, "KGS": true,
	"KHR": true, "KID": true, "KMF": true, "KRW": true, "KWD": true, "KYD": true,
	"KZT": true, "LAK": true, "LBP": true, "LKR": true, "LRD": true, "LSL": true,
	"LYD": true, "MAD": true, "MDL": true, "MGA": true, "MKD": true, "MMK": true,
	"MNT": true, "MOP": true, "MRU": true, "MUR": true, "MVR": true, "MWK": true,
	"MXN": true, "MYR": true, "MZN": true, "NAD": true, "NGN": true, "NIO": true,
	"NOK": true, "NPR": true, "NZD": true, "OMR": true, "PAB": true, "PEN": true,
	"PGK": true, "PHP": true, "PKR": true, "PLN": true, "PYG": true, "QAR": true,
	"RON": true, "RSD": true, "RUB": true, "RWF": true, "SAR": true, "SBD": true,
	"SCR": true, "SDG": true, "SEK": true, "SGD": true, "SHP": true, "SLE": true,
	"SOS": true, "SRD": true, "SSP": true, "STN": true, "SYP": true, "SZL": true,
	"THB": true, "TJS": true, "TMT": true, "TND": true, "TOP": true, "TRY": true,
	"TTD": true, "TVD": true, "TWD": true, "TZS": true, "UAH": true, "UGX": true,
	"USD": true, "UYU": true, "UZS": true, "VES": true, "VND": true, "VUV": true,
	"WST": true, "XAF": true, "XCD": true, "XDR": true, "XOF": true, "XPF": true,
	"YER": true, "ZAR": true, "ZMW": true, "ZWL": true}

type API_ANS struct {
	Result                string             `json:"result"`
	Documentation         string             `json:"documentation"`
	Terms_of_use          string             `json:"terms_of_use"`
	Time_last_update_unix uint               `json:"time_last_update_unix"`
	Time_last_update_utc  string             `json:"time_last_update_utc"`
	Time_next_update_unix uint               `json:"time_next_update_unix"`
	Time_next_update_utc  string             `json:"time_next_update_utc"`
	Base_code             string             `json:"base_code"`
	Conversion_rates      map[string]float64 `json:"conversion_rates"`
}

var API_KEY string = "YOUR_API_CODE"

func GetApiKey() string {
	godotenv.Load()

	apiKey := os.Getenv("EXCHANGERATE_API_KEY")
	if apiKey != "" {
		return apiKey
	}

	if API_KEY == "YOUR_API_CODE" {
		os.Exit(1)
	}
	return API_KEY
}

func GetRatio(curr_a string, curr_b string) (float64, error) {

	if len(curr_a) < 3 || len(curr_b) < 3 {
		return 0, fmt.Errorf("Currency code must be 3 characters")
	}

	_, ok_1 := supported_currencies[curr_a]
	_, ok_2 := supported_currencies[curr_b]

	if !ok_1 || !ok_2 {
		return 0, fmt.Errorf("Unsupported currency detected")
	}

	if curr_a == curr_b {
		return 1, nil
	}

	var website string = "https://v6.exchangerate-api.com/v6/" + GetApiKey() + "/latest/" + curr_a
	resp, err := http.Get(website)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var data API_ANS
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}

	if data.Result != "success" {
		return 0, fmt.Errorf("API returned error result: %s", data.Result)
	}

	ratio, success := data.Conversion_rates[curr_b]
	if !success {
		return 0, fmt.Errorf("currency %s not found", curr_b)
	}

	return ratio, nil
}
