package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// CurrencyRate представляет структуру данных для каждой валюты из API НБУ
type CurrencyRate struct {
	R030         int    `json:"r030"`
	Txt          string `json:"txt"`
	Rate         float64 `json:"rate"`
	Cc           string `json:"cc"`
	ExchangeDate string `json:"exchangedate"`
}

var (
	// Флаг для фильтрации по коду валюты
	currencyCode string
	// Флаг для указания даты
	date string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "exchange-rate",
		Short: "A CLI client that gives information about currency exchange rates",
	}

	// Команда для получения курсов валют
	getRatesCmd := &cobra.Command{
		Use:   "get",
		Short: "Fetch currency exchange rates from NBU API",
		Run: func(cmd *cobra.Command, args []string) {
			fetchExchangeRates()
		},
	}

	// Добавляем флаги для команды get
	getRatesCmd.Flags().StringVarP(&currencyCode, "currency", "c", "", "Currency code (e.g., USD, EUR)")
	getRatesCmd.Flags().StringVarP(&date, "date", "d", time.Now().Format("02.01.2006"), "Date for exchange rates (format: DD.MM.YYYY)")

	rootCmd.AddCommand(getRatesCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// fetchExchangeRates выполняет запрос к API 
func fetchExchangeRates() {
	url := "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?json"
	if date != "" {
		url = fmt.Sprintf("%s&date=%s", url, formatDateForAPI(date))
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: API returned status code %d\n", resp.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}

	var rates []CurrencyRate
	if err := json.Unmarshal(body, &rates); err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	// Фильтрация и вывод результатов
	for _, rate := range rates {
		if currencyCode == "" || rate.Cc == currencyCode {
			fmt.Printf("Currency: %s (%s), Rate: %.4f, Date: %s\n", rate.Txt, rate.Cc, rate.Rate, rate.ExchangeDate)
		}
	}
}

// formatDateForAPI преобразует дату из DD.MM.YYYY в YYYYMMDD для API
func formatDateForAPI(dateStr string) string {
	t, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		os.Exit(1)
	}
	return t.Format("20060102")
}