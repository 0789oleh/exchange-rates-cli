package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	currencyCode string
	date         string
	source       string
	target       string
	amount       float64
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "exchange-rate",
		Short: "A CLI client that gives information about currency exchange rates",
	}

	getRatesCmd := &cobra.Command{
		Use:   "get",
		Short: "Fetch currency exchange rates from NBU API",
		Run: func(cmd *cobra.Command, args []string) {
			// Передаем date в fetchExchangeRates
			rates, err := fetchExchangeRates(date)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			for _, rate := range rates {
				if currencyCode == "" || rate.Cc == currencyCode {
					fmt.Printf("Currency: %s (%s), Rate: %.4f, Date: %s\n", rate.Txt, rate.Cc, rate.Rate, rate.ExchangeDate)
				}
			}
		},
	}
	getRatesCmd.Flags().StringVarP(&currencyCode, "currency", "c", "", "Currency code (e.g., USD, EUR)")
	getRatesCmd.Flags().StringVarP(&date, "date", "d", "", "Date for exchange rates (format: DD.MM.YYYY)")

	convertCmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert amount from one currency to another",
		Run: func(cmd *cobra.Command, args []string) {
			if source == "" || target == "" {
				fmt.Println("Error: --source and --target are required")
				os.Exit(1)
			}

			rates, err := fetchExchangeRates(date)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			sourceRate := getRateForCurrency(rates, source)
			targetRate := getRateForCurrency(rates, target)

			if sourceRate == 0 || targetRate == 0 {
				fmt.Println("Error: Invalid source or target currency")
				os.Exit(1)
			}

			converted := amount * (sourceRate / targetRate)
			fmt.Printf("%.2f %s = %.4f %s (Date: %s)\n", amount, source, converted, target, rates[0].ExchangeDate)
		},
	}
	convertCmd.Flags().StringVarP(&source, "source", "s", "UAH", "Source currency code (e.g., USD, UAH)")
	convertCmd.Flags().StringVarP(&target, "target", "t", "", "Target currency code (e.g., EUR)")
	convertCmd.Flags().Float64VarP(&amount, "amount", "a", 1.0, "Amount to convert")
	convertCmd.Flags().StringVarP(&date, "date", "d", "", "Date for exchange rates (format: DD.MM.YYYY)")

	rootCmd.AddCommand(getRatesCmd, convertCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func getRateForCurrency(rates []CurrencyRate, cc string) float64 {
	if cc == "UAH" {
		return 1.0
	}
	for _, rate := range rates {
		if rate.Cc == cc {
			return rate.Rate
		}
	}
	return 0
}