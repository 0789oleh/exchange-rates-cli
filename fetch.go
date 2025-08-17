package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func fetchExchangeRates(dateStr string) ([]CurrencyRate, error) {
	if dateStr == "" {
		dateStr = time.Now().Format("02.01.2006")
	}

	apiDate := formatDateForAPI(dateStr)
	cacheDir := filepath.Join(os.Getenv("HOME"), ".cache", "exchange-rate")
	cacheFile := filepath.Join(cacheDir, fmt.Sprintf("%s.json", apiDate))

	// Check cache
	if _, err := os.Stat(cacheFile); err == nil {
		data, err := ioutil.ReadFile(cacheFile)
		if err != nil {
			return nil, err
		}
		var rates []CurrencyRate
		if err := json.Unmarshal(data, &rates); err != nil {
			return nil, err
		}
		return rates, nil
	}

	// Fetch from API
	url := fmt.Sprintf("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?json&date=%s", apiDate)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rates []CurrencyRate
	if err := json.Unmarshal(body, &rates); err != nil {
		return nil, err
	}

	// Save to cache
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(cacheFile, body, 0644); err != nil {
		return nil, err
	}

	// Update ExchangeDate to user-friendly format
	for i := range rates {
		rates[i].ExchangeDate = dateStr
	}

	return rates, nil
}

func formatDateForAPI(dateStr string) string {
	t, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		os.Exit(1)
	}
	return t.Format("20060102")
}