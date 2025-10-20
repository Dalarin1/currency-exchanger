package main

import (
	currency "app/pack"
	"flag"
	"fmt"
)

func main() {

	var curr_from = flag.String("f", "RUB", "currency_from")
	var curr_to = flag.String("t", "USD", "currency_to")

	flag.Parse()

	result, err := currency.GetRatio(*curr_from, *curr_to)

	if err != nil {
		fmt.Printf("АХТУН ЕРРОР %s\n", err)
	} else {
		fmt.Printf("%s to %s is %f\n", *curr_from, *curr_to, result)
		fmt.Printf("%s to %s is %f\n", *curr_to, *curr_from, 1/result)
	}
}
