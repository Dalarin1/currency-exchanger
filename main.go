package main

import (
	currency "app/pack"
	"flag"
	"fmt"
)

func main() {

	var curr_base = flag.String("f", "RUB", "Base currency")
	var curr_target = flag.String("t", "", "Target currency")
	var amount = flag.Float64("a", 1.0, "Amount of base currency")
	var historical = flag.Bool("hist", false, "Show hystroical data about base currency")
	var enriched = flag.Bool("e", false, "Show enriched data about base currency")
	var help = flag.Bool("h", false, "Show help")
	flag.Parse()

	if *historical && *enriched {
		fmt.Println("ERROR: Concurent flags -hist and -e")
		return
	}
	if *help {
		fmt.Println("HELP STRING")
		return
	}
	if len(*curr_target) == 0 {
		if *historical {
			// TODO
			fmt.Println("HISTORICAL")
		} else if *enriched {
			a, err := currency.GetEnrichedData(*curr_base, *curr_target)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%s to %s : %f", *curr_base, *curr_target, a.Conversion_rate)
				fmt.Println(a.Target_data)
			}

		} else {
			a, err := currency.GetStdData(*curr_base)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%s to other currencies is:\n", a.Base_code)
				fmt.Println(a.Conversion_rates)
			}
		}
	} else {
		a, err := currency.GetPairData(*curr_base, *curr_target, *amount)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%f %s = %f %s", *amount, *curr_base, a.Conversion_result, *curr_target)
		}
	}
}
