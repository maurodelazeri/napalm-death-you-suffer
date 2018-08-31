package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"engo.io/audio"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/sirupsen/logrus"
)

const (
	CoinMarketCapIOSOfBTCTUrl = "https://api.coinmarketcap.com/v2/ticker"
)

type Coinmarketcap struct {
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
	if _, err := os.Stat("/tmp/sound.wav"); os.IsNotExist(err) {
		// path/to/whatever does not exist
	}
	player, err := audio.NewSimplePlayer("/tmp/sound.wav")
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

	var btc Coinmarketcap
	err = ffjson.Unmarshal(body, &btc)
	if err != nil {
		logrus.Error(err)
	}

	if btc.Data.Num1.Quotes.USD.PercentChange1H > percentage {
		alert()
	}
}

func PrintDownloadPercent(done chan int64, path string, total int64) {
	var stop bool = false
	for {
		select {
		case <-done:
			stop = true
		default:
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			fi, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}
			size := fi.Size()
			if size == 0 {
				size = 1
			}
			var percent float64 = float64(size) / float64(total) * 100
			fmt.Printf("%.0f", percent)
			fmt.Println("%")
		}
		if stop {
			break
		}
		time.Sleep(time.Second)
	}
}

func DownloadFile(url string, dest string) {
	file := path.Base(url)
	log.Printf("Downloading file %s from %s\n", file, url)
	var path bytes.Buffer
	path.WriteString(dest)
	path.WriteString("/")
	path.WriteString(file)
	start := time.Now()
	out, err := os.Create(path.String())
	if err != nil {
		fmt.Println(path.String())
		panic(err)
	}
	defer out.Close()
	headResp, err := http.Head(url)
	if err != nil {
		panic(err)
	}
	defer headResp.Body.Close()
	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))
	if err != nil {
		panic(err)
	}
	done := make(chan int64)
	go PrintDownloadPercent(done, path.String(), int64(size))
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	done <- n
	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)
}
