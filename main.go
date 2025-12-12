package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	currency "app/pack"
)

func main() {
	currBase := flag.String("f", "RUB", "Base currency (e.g. USD)")
	currTarget := flag.String("t", "", "Target currency (e.g. EUR)")
	amount := flag.Float64("a", 1.0, "Amount of base currency")
	historical := flag.String("hist", "", "Show historical data (format: YYYY-MM-DD)")
	enriched := flag.Bool("e", false, "Show enriched data")
	help := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *historical != "" && *enriched {
		log.Fatal("ERROR: flags -hist and -e cannot be used together")
	}

	switch {
	case *historical != "":
		key := fmt.Sprintf("-f %s -hist %s", *currBase, *historical)
		withCache(key, func() string {
			return showHistorical(*currBase, *historical, *amount)
		})

	case *enriched:
		key := fmt.Sprintf("-f %s -t %s -e", *currBase, *currTarget)
		withCache(key, func() string {
			return showEnriched(*currBase, *currTarget)
		})

	case *currTarget == "":
		key := fmt.Sprintf("-f %s", *currBase)
		withCache(key, func() string {
			return showStandard(*currBase)
		})

	default:
		key := fmt.Sprintf("-f %s -t %s -a %.2f", *currBase, *currTarget, *amount)
		withCache(key, func() string {
			return showPair(*currBase, *currTarget, *amount)
		})
	}
}

func withCache(key string, fetch func() string) {
	if cached, timestamp, err := currency.FindCacheEntry(key); err == nil {
		age := time.Since(time.Unix(timestamp, 0))
		if age < 6*time.Hour {
			fmt.Println("[cache hit]")
			fmt.Println(cached)
			return
		}
	}

	fmt.Println("[cache miss]")
	result := fetch()
	currency.StoreCacheEntry(key, result)
	fmt.Println(result)
}

func printHelp() {
	fmt.Println(`Usage:
  app [flags]

Flags:
  -f string       Base currency (default "RUB")
  -t string       Target currency (optional)
  -a float        Amount (default 1.0)
  -hist date      Show historical data (format: YYYY-MM-DD)
  -e              Show enriched data
  -h              Show this help`)
}

func showHistorical(base, dateStr string, amount float64) string {
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatalf("Invalid date format: %v", err)
	}

	data, err := currency.GetHistoricalData(base, parsed, amount)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return currency.FormatHistoricalData(data, false)
}

func showEnriched(base, target string) string {
	if target == "" {
		log.Fatal("ERROR: target currency (-t) required for enriched data")
	}

	data, err := currency.GetEnrichedData(base, target)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return currency.FormatEnrichedData(data)
}

func showStandard(base string) string {
	data, err := currency.GetStdData(base)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return currency.FormatStdData(data)
}

func showPair(base, target string, amount float64) string {
	data, err := currency.GetPairData(base, target, amount)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return currency.FormatPairData(data, amount)
}
