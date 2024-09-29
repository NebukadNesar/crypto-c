package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

// singleton instance and sync.Once
var (
	instance         *DataCache
	once             sync.Once
	MaxHistoryRecord = 100
)

type DataCache struct {
	latestSnapShot    CryptoResponse
	historicalRecords map[string][]HistoricalData
	mu                sync.RWMutex // Mutex for concurrent access
}

func GetCryptoCache() *DataCache {
	once.Do(func() {
		instance = &DataCache{
			historicalRecords: make(map[string][]HistoricalData),
		}
		go controlCache()
	})
	return instance
}

func controlCache() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				instance.mu.Lock()
				if instance.historicalRecords != nil {
					for _, historyRecords := range instance.historicalRecords {
						if len(historyRecords) > MaxHistoryRecord {
							historyRecords = historyRecords[:MaxHistoryRecord]
						}
					}
				}
				instance.mu.Unlock()
			}
		}
	}()
	select {}
}

type BlockDataCache interface {
	History(coin string) ([]HistoricalData, error)
	Get(coin string) (HistoricalData, error)
	Put(coin string, data HistoricalData)

	UpdateLatestShapShot([]HistoricalData)
	GetLatestShapShot() []HistoricalData
}

func (cache *DataCache) History(coin string) ([]HistoricalData, error) {
	cache.mu.RLock() // normally we do not need but just for play
	defer cache.mu.RUnlock()
	historicalData, ok := cache.historicalRecords[coin]
	if !ok {
		return nil, errors.New("cannot find the requested historical coin data")
	}
	return historicalData, nil
}

func (cache *DataCache) Get(coin string) (HistoricalData, error) {
	cache.mu.RLock() // normally we do not need ReadLock but just for play
	defer cache.mu.RUnlock()

	historicalData, ok := cache.historicalRecords[coin]
	if !ok || len(historicalData) == 0 {
		return HistoricalData{}, errors.New("cannot find the requested coin data")
	}
	return historicalData[len(historicalData)-1], nil
}

func (cache *DataCache) Put(coin string, data HistoricalData) {
	cache.mu.Lock() // here we need it, for keeping the order
	defer cache.mu.Unlock()
	log.Printf("Put %s historical record\n", coin)
	cache.historicalRecords[coin] = append(cache.historicalRecords[coin], data)
}

func (cache *DataCache) UpdateLatestShapShot(snapshot CryptoResponse) {
	cache.latestSnapShot = snapshot
}

func (cache *DataCache) GetLatestShapShot() CryptoResponse {
	return cache.latestSnapShot
}
