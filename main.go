package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"engo.io/audio"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/sirupsen/logrus"
)

const (
	CoinMarketCapIOSOfBTCTUrl = "https://api.coinmarketcap.com/v2/ticker"
)

type AutoGenerated struct {
	Data struct {
		Num1 struct {
			ID                int     `json:"id"`
			Name              string  `json:"name"`
			Symbol            string  `json:"symbol"`
			WebsiteSlug       string  `json:"website_slug"`
			Rank              int     `json:"rank"`
			CirculatingSupply float64 `json:"circulating_supply"`
			TotalSupply       float64 `json:"total_supply"`
			MaxSupply         float64 `json:"max_supply"`
			Quotes            struct {
				USD struct {
					Price            float64 `json:"price"`
					Volume24H        float64 `json:"volume_24h"`
					MarketCap        float64 `json:"market_cap"`
					PercentChange1H  float64 `json:"percent_change_1h"`
					PercentChange24H float64 `json:"percent_change_24h"`
					PercentChange7D  float64 `json:"percent_change_7d"`
				} `json:"USD"`
			} `json:"quotes"`
			LastUpdated int `json:"last_updated"`
		} `json:"1"`
	}
}

var percentage float64

func main() {
	input := os.Args[1:]
	if len(input) > 0 {
		i, err := strconv.ParseFloat(input[0], 64)
		if err != nil {
			percentage = 1.0
		} else {
			percentage = i
		}
	} else {
		percentage = 1.0
	}
	getMarketInfo()
}

func alert() {
	player, err := audio.NewSimplePlayer("napalm-death-you-suffer.wav")
	if err != nil {
		panic(err)
	}
	player.SetVolume(0.9)
	player.Play()
	time.Sleep(time.Second * 5)
}

func getMarketInfo() {
	resp, err := http.Get(CoinMarketCapIOSOfBTCTUrl)
	if err != nil {
		logrus.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
	}

	var btc AutoGenerated
	err = ffjson.Unmarshal(body, &btc)
	if err != nil {
		logrus.Error(err)
	}

	if btc.Data.Num1.Quotes.USD.PercentChange1H > percentage {
		alert()
	}
}
