package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

type BlockDataScrapperContract interface {
	FetchCryptoData() (CryptoResponse, error)
}

type BlockDataScrapper struct {
	APIURL string // api address to scrap blockchain data
	sort   bool   // sort the incoming data according to natural order of the blockchain names
	limit  int    // limit the number of element that we get from the api based on the volume
}

func (scrapper *BlockDataScrapper) FetchCryptoData() (CryptoResponse, error) {
	fmt.Println("....Fetching crypto data....")
	response, err := http.Get(scrapper.APIURL)

	if err != nil {
		fmt.Println("Error:", err)
		return CryptoResponse{}, err
	}

	fmt.Println(response.Status)
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error:", err)
		return CryptoResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println(err.Error())
		}
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return CryptoResponse{}, err
	}
	var cryptoResponse CryptoResponse

	err = json.Unmarshal(body, &cryptoResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return CryptoResponse{}, err
	}

	if scrapper.sort {
		sort.Slice(cryptoResponse.Data, func(i, j int) bool {
			capI, errI := strconv.ParseFloat(cryptoResponse.Data[i].MarketCapUsd, 64)
			capJ, errJ := strconv.ParseFloat(cryptoResponse.Data[j].MarketCapUsd, 64)
			if errI != nil || errJ != nil {
				return false
			}
			return capI > capJ
		})
	}

	limit := scrapper.limit
	if limit > 0 && len(cryptoResponse.Data) > limit {
		cryptoResponse.Data = cryptoResponse.Data[:limit]
	}

	return cryptoResponse, nil
}
