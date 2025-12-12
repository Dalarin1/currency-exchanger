package currency

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

const cacheFile = "cache.json"

// структура одной записи кэша
type CacheEntry struct {
	Time   int64  `json:"time"`
	Result string `json:"result"`
}

// ищем запись в JSON-кэше
func FindCacheEntry(callArgs string) (string, int64, error) {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", -1, errors.New("cache file not found")
		}
		return "", -1, err
	}

	var cache map[string]CacheEntry
	if err := json.Unmarshal(data, &cache); err != nil {
		return "", -1, err
	}

	entry, ok := cache[callArgs]
	if !ok {
		return "", -1, errors.New("cache entry not found")
	}

	return entry.Result, entry.Time, nil
}

// сохраняем запись в JSON-кэш
func StoreCacheEntry(callArgs, result string) {
	var cache map[string]CacheEntry

	// читаем существующий кэш (если есть)
	data, err := os.ReadFile(cacheFile)
	if err == nil {
		_ = json.Unmarshal(data, &cache)
	}
	if cache == nil {
		cache = make(map[string]CacheEntry)
	}

	// добавляем/обновляем запись
	cache[callArgs] = CacheEntry{
		Time:   time.Now().Unix(),
		Result: result,
	}

	// сохраняем обратно в файл
	newData, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(cacheFile, newData, 0644)
}
