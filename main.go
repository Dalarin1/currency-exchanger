package main

import (
	currency "app/pack"
	"flag"
	"fmt"
	"time"
)

func main() {

	var curr_base = flag.String("f", "RUB", "Base currency")
	var curr_target = flag.String("t", "", "Target currency")
	var amount = flag.Float64("a", 1.0, "Amount of base currency")
	var historical = flag.String("hist", "", "Show hystroical data about base currency")
	var enriched = flag.Bool("e", false, "Show enriched data about base currency")
	var help = flag.Bool("h", false, "Show help")
	flag.Parse()

	if *historical != "" && *enriched {
		fmt.Println("ERROR: Concurent flags -hist and -e")
		return
	}
	if *help {
		fmt.Println("HELP STRING")
		return
	}
	if len(*curr_target) == 0 {
		if *historical != "" {
			// HIST
			parsed_date, err := time.Parse("2006-01-02", *historical)
			if err != nil {
				fmt.Println(err)
			} else {
				hdata, err := currency.GetHystoricalData(*curr_base, parsed_date, *amount)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(currency.FormatHystoricalData(hdata, false))
				}
			}

		} else if *enriched {
			// ENRICH
			endata, err := currency.GetEnrichedData(*curr_base, *curr_target)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(currency.FormatEnrichedData(endata))
			}
		} else {
			// STD
			stddata, err := currency.GetStdData(*curr_base)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(currency.FormatStdData(stddata))
			}
		}
	} else {
		// PAIR
		pairdata, err := currency.GetPairData(*curr_base, *curr_target, *amount)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(currency.FormatPairData(pairdata, *amount))
		}
	}
}
