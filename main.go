package main

import (
	"context"
	"log"
	"os"

	"github.com/GhandyP/adk-go-enterprise-cookbook/internal/adkhelper"
)

func main() {
	if err := adkhelper.Run(
		context.Background(),
		os.Args[1:],
		"case_42_sales_scripts_calls",
		"A text-only starting point for sales scripts and calls.",
		"Create a sales call script from the customer facts and approved objections supplied by the user. Personalize only supplied facts and list forbidden claims.",
	); err != nil {
		log.Fatal(err)
	}
}
