package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var API_KEY string = "YOUR_API_CODE"

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

var old_currencies map[string]bool = map[string]bool{
	"AUD": true, "ATS": true, "BEF": true, "BRL": true,
	"CAD": true, "CHF": true, "CNY": true, "DEM": true,
	"DKK": true, "ESP": true, "EUR": true, "FIM": true,
	"FRF": true, "GBP": true, "GRD": true, "HKD": true,
	"IEP": true, "INR": true, "IRR": true, "ITL": true,
	"JPY": true, "KRW": true, "LKR": true, "MXN": true,
	"MYR": true, "NOK": true, "NLG": true, "NZD": true,
	"PTE": true, "SEK": true, "SGD": true, "THB": true,
	"TWD": true, "USD": true, "ZAR": true}

var min_avaliable_date time.Time = time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)
var border_time time.Time = time.Date(2020, time.December, 31, 0, 0, 0, 0, time.UTC)

type STD_API_ANS struct {
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

func (ans STD_API_ANS) GetResult() string {
	return ans.Result
}

type PAIR_API_ANS struct {
	Result                string  `json:"result"`
	Documentation         string  `json:"documentation"`
	Terms_of_use          string  `json:"terms_of_use"`
	Time_last_update_unix uint    `json:"time_last_update_unix"`
	Time_last_update_utc  string  `json:"time_last_update_utc"`
	Time_next_update_unix uint    `json:"time_next_update_unix"`
	Time_next_update_utc  string  `json:"time_next_update_utc"`
	Base_code             string  `json:"base_code"`
	Target_code           string  `json:"target_code"`
	Conversion_rate       float64 `json:"conversion_rate"`
	Conversion_result     float64 `json:"conversion_result"`
}

func (ans PAIR_API_ANS) GetResult() string {
	return ans.Result
}

type ENRICHED_API_ANS struct {
	Result                string            `json:"result"`
	Documentation         string            `json:"documentation"`
	Terms_of_use          string            `json:"terms_of_use"`
	Time_last_update_unix uint              `json:"time_last_update_unix"`
	Time_last_update_utc  string            `json:"time_last_update_utc"`
	Time_next_update_unix uint              `json:"time_next_update_unix"`
	Time_next_update_utc  string            `json:"time_next_update_utc"`
	Base_code             string            `json:"base_code"`
	Target_code           string            `json:"target_code"`
	Conversion_rate       float64           `json:"conversion_rate"`
	Target_data           map[string]string `json:"target_data"`
}

func (ans ENRICHED_API_ANS) GetResult() string {
	return ans.Result
}

type HYSTORICAL_API_ANS struct {
	Result           string             `json:"result"`
	Documentation    string             `json:"documentation"`
	Terms_of_use     string             `json:"terms_of_use"`
	Year             uint               `json:"year"`
	Month            uint               `json:"month"`
	Day              uint8              `json:"day"`
	Base_code        string             `json:"base_code"`
	Conversion_rates map[string]float64 `json:"conversion_rates"`
}

func (ans HYSTORICAL_API_ANS) GetResult() string {
	return ans.Result
}
func CheckCurrencyValid(currency string) bool {
	_, ok := supported_currencies[currency]
	return len(currency) == 3 && ok
}

type API_ANS interface {
	STD_API_ANS | ENRICHED_API_ANS | PAIR_API_ANS | HYSTORICAL_API_ANS
	GetResult() string
}

func GetData[T API_ANS](website string) (T, error) {
	var zero T
	resp, err := http.Get(website)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, fmt.Errorf("failed to read response: %w", err)
	}

	var data T
	err = json.Unmarshal(body, &data)
	if err != nil {
		return zero, err
	}
	if data.GetResult() != "success" {
		return zero, fmt.Errorf("API returned error result: %s", data.GetResult())
	}
	return data, nil
}

func GetStdData(currency string) (STD_API_ANS, error) {
	if !CheckCurrencyValid(currency) {
		return STD_API_ANS{}, fmt.Errorf("Wrong or unsupported currency")
	}
	var website string = "https://v6.exchangerate-api.com/v6/" + GetApiKey() + "/latest/" + currency

	return GetData[STD_API_ANS](website)
}

func GetHystoricalData(currency string, year, month, day int) (HYSTORICAL_API_ANS, error) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	if date.After(border_time) {
		if !CheckCurrencyValid(currency) {
			return HYSTORICAL_API_ANS{}, fmt.Errorf("ERROR: Unsupported currency")
		}
	} else {
		_, ok := old_currencies[currency]
		if !ok {
			return HYSTORICAL_API_ANS{}, fmt.Errorf("ERROR: Unsupported currency")
		}
		if date.Before(min_avaliable_date) {
			return HYSTORICAL_API_ANS{}, fmt.Errorf("ERROR: To old date; 01/01/1990 is oldes supported date")
		}
	}
	var website string = "https://v6.exchangerate-api.com/v6/" + GetApiKey() + "/history/" + currency + "/" + strconv.FormatInt((int64)(year), 10) + "/" + strconv.FormatInt((int64)(month), 10) + "/" + strconv.FormatInt((int64)(day), 10)

	return GetData[HYSTORICAL_API_ANS](website)
}

func GetEnrichedData(curr_a, curr_b string) (ENRICHED_API_ANS, error) {
	if !CheckCurrencyValid(curr_a) || !CheckCurrencyValid(curr_b) {
		return ENRICHED_API_ANS{}, fmt.Errorf("Wrong or unsupported currency")
	}
	var website string = "https://v6.exchangerate-api.com/v6/" + GetApiKey() + "/enriched/" + curr_a + "/" + curr_b

	return GetData[ENRICHED_API_ANS](website)
}

func GetPairData(curr_a, curr_b string, amount float64) (PAIR_API_ANS, error) {
	if !CheckCurrencyValid(curr_a) || !CheckCurrencyValid(curr_b) {
		return PAIR_API_ANS{}, fmt.Errorf("Wrong or unsupported currency")
	}

	if curr_a == curr_b {
		if amount != 0 && amount > 0 {
			return PAIR_API_ANS{Conversion_rate: 1, Conversion_result: amount}, nil
		} else {
			return PAIR_API_ANS{Conversion_rate: 1}, nil
		}
	}

	var website string = "https://v6.exchangerate-api.com/v6/" + GetApiKey() + "/pair/" + curr_a + "/" + curr_b
	if amount > 0 {
		website += "/" + strconv.FormatFloat(amount, 'f', 4, 64)
	}

	return GetData[PAIR_API_ANS](website)
}

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
