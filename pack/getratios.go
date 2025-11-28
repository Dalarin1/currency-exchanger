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
		body, _ := io.ReadAll(resp.Body)
		var ans map[string]string
		json.Unmarshal(body, &ans)
		return zero, fmt.Errorf("API returned status: %d\nError type: %v", resp.StatusCode, ans["error-type"])
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
	api_key, err := GetApiKey()
	if err == nil {
		var website string = "https://v6.exchangerate-api.com/v6/" + api_key + "/latest/" + currency

		return GetData[STD_API_ANS](website)
	} else {
		return STD_API_ANS{}, err
	}

}

func GetHystoricalData(currency string, date time.Time, amount float64) (HYSTORICAL_API_ANS, error) {
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
			return HYSTORICAL_API_ANS{}, fmt.Errorf("ERROR: Too old date; 01/01/1990 is oldes supported date")
		}
	}
	api_key, err := GetApiKey()
	if err == nil {
		var website string = "https://v6.exchangerate-api.com/v6/" + api_key + "/history/" + currency + "/" + strconv.FormatInt((int64)(date.Year()), 10) + "/" + strconv.FormatInt((int64)(date.Month()), 10) + "/" + strconv.FormatInt((int64)(date.Day()), 10)
		if amount > 0 {
			website += "/" + strconv.FormatFloat(amount, 'f', 4, 64)
		}
		return GetData[HYSTORICAL_API_ANS](website)
	} else {
		return HYSTORICAL_API_ANS{}, err
	}

}

func GetEnrichedData(curr_a, curr_b string) (ENRICHED_API_ANS, error) {
	if !CheckCurrencyValid(curr_a) || !CheckCurrencyValid(curr_b) {
		return ENRICHED_API_ANS{}, fmt.Errorf("Wrong or unsupported currency")
	}
	api_key, err := GetApiKey()
	if err == nil {
		var website string = "https://v6.exchangerate-api.com/v6/" + api_key + "/enriched/" + curr_a + "/" + curr_b

		return GetData[ENRICHED_API_ANS](website)
	} else {
		return ENRICHED_API_ANS{}, err
	}

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
	api_key, err := GetApiKey()
	if err == nil {
		var website string = "https://v6.exchangerate-api.com/v6/" + api_key + "/pair/" + curr_a + "/" + curr_b
		if amount > 0 {
			website += "/" + strconv.FormatFloat(amount, 'f', 4, 64)
		}
		res, err := GetData[PAIR_API_ANS](website)
		return res, err
	} else {
		return PAIR_API_ANS{}, err
	}
}

func GetApiKey() (string, error) {
	godotenv.Load()

	apiKey := os.Getenv("EXCHANGERATE_API_KEY")
	if apiKey != "" {
		return apiKey, nil
	}
	return "", fmt.Errorf("Api key does not found")
}

func FormatStdData(stddata STD_API_ANS) string {
	var result string = ""
	for c, v := range stddata.Conversion_rates {
		result += fmt.Sprintf(" %s => %f \n", c, v)
	}
	return result
}

func FormatPairData(pairdata PAIR_API_ANS, amount float64) string {
	var result string = ""
	result += fmt.Sprintf("1 %s = %f %s\n", pairdata.Base_code, pairdata.Conversion_rate, pairdata.Target_code)
	if amount != 0 {
		result += fmt.Sprintf("%f %s = %f %s\n", amount, pairdata.Base_code, pairdata.Conversion_result, pairdata.Target_code)
	}
	return result
}

func FormatEnrichedData(endata ENRICHED_API_ANS) string {
	var result string = ""
	result += fmt.Sprintf("1 %s = %f %s", endata.Base_code, endata.Conversion_rate, endata.Target_code)
	/*"target_data": {
		"locale": "Japan",
		"two_letter_code": "JP",
		"currency_name": "Japanese Yen",
		"currency_name_short": "Yen",
		"display_symbol": "00A5",
		"flag_url": "https://www.exchangerate-api.com/img/docs/JP.gif"
	}*/
	result += "Target data: \n"
	result += "\tLocale: " + endata.Target_data["locale"] + "\n"
	result += "\tTwo letter code: " + endata.Target_data["two_letter_code"] + "\n"
	result += "\tCurrency name: " + endata.Target_data["currency_name"] + "\n"
	result += "\tCurrency short name: " + endata.Target_data["currency_short_name"] + "\n"
	result += "\tDisplay symbol: " + endata.Target_data["display_symbol"] + "\n"
	result += "\tFlag url: " + endata.Target_data["flag_url"] + "\n"
	return result
}

func FormatHystoricalData(hdata HYSTORICAL_API_ANS, use_worst_date_format_ever bool) string {
	var result string = ""
	if use_worst_date_format_ever {
		result += fmt.Sprintf("Date: %d/%d/%d \n", hdata.Month, hdata.Day, hdata.Year)
	} else {
		result += fmt.Sprintf("Date: %d/%d/%d \n", hdata.Day, hdata.Month, hdata.Year)
	}
	result += "Base code: " + hdata.Base_code + "\n"
	for c, v := range hdata.Conversion_rates {
		result += fmt.Sprintf(" %s => %f \n", c, v)
	}
	return result
}
