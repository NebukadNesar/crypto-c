package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CryptoDataRestApi interface {
	GetHistoricalDataByCoin(coin string) ([]HistoricalData, error)
	GetLatestSnapshotByCoin(coin string) (HistoricalData, error)
}

// CryptoDataAPI implements the CryptoDataRestApi interface.
type CryptoDataAPI struct {
	cache *DataCache // DataCache instance for storing historical data
}

// NewCryptoDataAPI create new instance of the crypto rest api
func NewCryptoDataAPI(cache *DataCache) *CryptoDataAPI {
	return &CryptoDataAPI{
		cache: cache,
	}
}

func (api *CryptoDataAPI) GetHistoricalDataByCoin(coin string) ([]HistoricalData, error) {
	history, err := api.cache.History(coin)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (api *CryptoDataAPI) GetLatestSnapshotByCoin(coin string) (HistoricalData, error) {
	snapshot, err := api.cache.Get(coin)
	if err != nil {
		return HistoricalData{}, err
	}
	return snapshot, nil
}

func (api *CryptoDataAPI) start() {
	r := setupRouter(api)

	if err := r.Run("localhost:9090"); err != nil {
		fmt.Println("Error starting REST server:", err)
	}
}

func setupRouter(api *CryptoDataAPI) *gin.Engine {
	r := gin.Default()
	r.GET("/crypto/history/:coin", func(c *gin.Context) {
		historicalDataByCoin, err := api.GetHistoricalDataByCoin(c.Param("coin"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, historicalDataByCoin)
	})
	r.GET("/crypto/snapshot/:coin", func(c *gin.Context) {
		snapshot, err := api.GetLatestSnapshotByCoin(c.Param("coin"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, snapshot)
	})
	return r
}
