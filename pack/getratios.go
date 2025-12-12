package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

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

type HISTORICAL_API_ANS struct {
	Result           string             `json:"result"`
	Documentation    string             `json:"documentation"`
	Terms_of_use     string             `json:"terms_of_use"`
	Year             uint               `json:"year"`
	Month            uint               `json:"month"`
	Day              uint8              `json:"day"`
	Base_code        string             `json:"base_code"`
	Conversion_rates map[string]float64 `json:"conversion_rates"`
}

func (ans HISTORICAL_API_ANS) GetResult() string {
	return ans.Result
}
func CheckCurrencyValid(currency string) bool {
	_, ok := SupportedCurrencies[currency]
	return len(currency) == 3 && ok
}

type API_ANS interface {
	STD_API_ANS | ENRICHED_API_ANS | PAIR_API_ANS | HISTORICAL_API_ANS
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

func GetHistoricalData(currency string, date time.Time, amount float64) (HISTORICAL_API_ANS, error) {
	if date.After(border_time) {
		if !CheckCurrencyValid(currency) {
			return HISTORICAL_API_ANS{}, fmt.Errorf("ERROR: Unsupported currency")
		}
	} else {
		_, ok := OldCurrencies[currency]
		if !ok {
			return HISTORICAL_API_ANS{}, fmt.Errorf("ERROR: Unsupported currency")
		}
		if date.Before(min_avaliable_date) {
			return HISTORICAL_API_ANS{}, fmt.Errorf("ERROR: Too old date; 01/01/1990 is oldes supported date")
		}
	}
	api_key, err := GetApiKey()
	if err == nil {
		var website string = "https://v6.exchangerate-api.com/v6/" + api_key + "/history/" + currency + "/" + strconv.FormatInt((int64)(date.Year()), 10) + "/" + strconv.FormatInt((int64)(date.Month()), 10) + "/" + strconv.FormatInt((int64)(date.Day()), 10)
		if amount > 0 {
			website += "/" + strconv.FormatFloat(amount, 'f', 4, 64)
		}
		return GetData[HISTORICAL_API_ANS](website)
	} else {
		return HISTORICAL_API_ANS{}, err
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
	// Получаем путь к текущему файлу (getratios.go)
	_, filename, _, _ := runtime.Caller(0)
	// Определяем путь к каталогу пакета
	dir := filepath.Dir(filename)
	// Строим путь к .env в родительской директории
	envPath := filepath.Join(dir, "..", ".env")

	data, err := os.ReadFile(envPath)
	if err != nil {
		return "", fmt.Errorf("cannot read %s: %w", envPath, err)
	}

	text := string(data)
	parts := strings.SplitN(text, "=", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1]), nil
	}

	return "", fmt.Errorf("API key not found in %s", envPath)
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
	result += "\tCurrency short name: " + endata.Target_data["currency_name_short"] + "\n"
	result += "\tDisplay symbol: " + endata.Target_data["display_symbol"] + "\n"
	result += "\tFlag url: " + endata.Target_data["flag_url"] + "\n"
	return result
}

func FormatHistoricalData(hdata HISTORICAL_API_ANS, use_worst_date_format_ever bool) string {
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
