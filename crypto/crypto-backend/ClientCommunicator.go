package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	allowedOrigins     = "http://localhost:9988"
	blockDataSourceUrl = "https://api.coincap.io/v2/assets"
	blockDataScrapper  BlockDataScrapperContract
	upgrade            = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins (modify as needed for security)
		},
	}
	_cryptoCache = GetCryptoCache()
)

func main() {
	initScrapper()

	//go initRestApi()
	go runAutomaticScrapper()

	http.Handle("/crypto", handleCORS(http.HandlerFunc(handleConnection)))
	log.Println("WebSocket server started on :9988")
	if err := http.ListenAndServe(":9988", nil); err != nil {
		log.Println("Error starting server:", err)
	}
}

func initRestApi() {
	cryptoDataRestApi := NewCryptoDataAPI(_cryptoCache)
	///cryptoDataRestApi
	cryptoDataRestApi.start()
}

func initScrapper() {
	blockDataScrapper = &BlockDataScrapper{
		APIURL: blockDataSourceUrl,
		sort:   true,
		limit:  10,
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error during connection upgrade:", err)
		return
	}

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Error during connection close:", err)
		}
	}(conn)

	go sendDataToClient(conn)

	messageType, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}
	log.Printf("Received: %s\n", msg)

	if _cryptoCache.GetLatestShapShot().Data != nil {
		cachedCryptoData, err := json.Marshal(_cryptoCache.GetLatestShapShot())
		if err != nil {
			log.Println("Error writing message:", err)
			return
		}
		err = conn.WriteMessage(messageType, cachedCryptoData)
		if err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}

// handleCORS is a middleware function to handle CORS headers
func handleCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler, chain of handlers
		h.ServeHTTP(w, r)
	})
}

func sendDataToClient(conn *websocket.Conn) {
	send(conn)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	log.Println("starting scrapper")
	go func() {
		for {
			select {
			case <-ticker.C:
				send(conn)
			}
		}
	}()
	select {}
}

func send(conn *websocket.Conn) {
	err, _ := sendSnapShotData(conn)
	if err != nil {
		//log
	}
	err2, _ := sendHistoricalData(conn)
	if err2 != nil {
		//log
	}
}

func sendSnapShotData(conn *websocket.Conn) (error, bool) {
	cryptoData := _cryptoCache.GetLatestShapShot().Data
	if len(cryptoData) == 0 {
		log.Println("No crypto data available.")
		return nil, true
	}

	snapshot := &ClientMessage{
		ID:   "snapshot",
		Data: cryptoData,
	}

	snapshotMessage, err := json.Marshal(snapshot)
	if err != nil {
		return nil, true
	}

	err = conn.WriteMessage(websocket.TextMessage, snapshotMessage)
	if err != nil {
		log.Println("Error writing crypto data:", err)
	}
	return err, false
}

func sendHistoricalData(conn *websocket.Conn) (error, bool) {
	historicalData, err := _cryptoCache.History("bitcoin")
	if err != nil || historicalData == nil {
		log.Println("Error marshaling historical data:", err)
		return err, false
	}

	historicalMessage := &ClientMessage{
		ID:   "historical",
		Data: historicalData,
	}

	byteData, err := json.Marshal(historicalMessage)
	if err != nil {
		log.Println("Error marshaling client message:", err)
		return err, false
	}

	log.Println("bitcoin chart data to send to clients len=", len(historicalData))
	err = conn.WriteMessage(websocket.TextMessage, byteData)
	if err != nil {
		log.Println("Error sending client message:", err)
		return err, false
	}
	return nil, true
}

func runAutomaticScrapper() {
	// Initial call to fetch crypto data
	getCryptoData()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("Starting scraper")

	go func() {
		for {
			select {
			case <-ticker.C:
				getCryptoData()
			}
		}
	}()
	select {}
}

func getCryptoData() {
	log.Println("New scraping started...")

	cryptoData, err := blockDataScrapper.FetchCryptoData()
	if err != nil {
		log.Println("Error fetching crypto data:", err.Error())
		return
	}

	if cryptoData.Data == nil {
		log.Println("Error: No crypto data received")
		return
	}

	_cryptoCache.UpdateLatestShapShot(cryptoData)
	for _, cryptoData := range cryptoData.Data {
		record := HistoricalData{
			Timestamp: time.Now(),
			Data:      cryptoData,
		}
		_cryptoCache.Put(cryptoData.ID, record)
	}
	log.Println("Crypto data cached successfully")
}
