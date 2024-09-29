package main

import "time"

var BlockDataChannel = make(chan CryptoResponse)

type CryptoData struct {
	ID                string `json:"id"`
	Rank              string `json:"rank"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Supply            string `json:"supply"`
	MaxSupply         string `json:"maxSupply"`
	MarketCapUsd      string `json:"marketCapUsd"`
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
	PriceUsd          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
	Vwap24Hr          string `json:"vwap24Hr"`
	Explorer          string `json:"explorer"`
}

type CryptoResponse struct {
	Data []CryptoData `json:"data"`
}

type HelloMessage struct {
	Text string `json:"text"`
}

type Cache struct {
	blockData CryptoResponse
}

type HistoricalData struct {
	Timestamp time.Time
	Data      CryptoData
}

type ClientMessage struct {
	ID   string `json:"id"`
	Data any    `json:"data"`
}
