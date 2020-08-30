package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	//"reflect"
	"strconv"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	stockPrice        map[string][]int64
	debug             = flag.Bool("debug", false, "Enable debug output")
	jsonout           = flag.Bool("json", false, "Enable REST API output")
	currentPrice      float64
	lastPrice         float64
	priceList         []float64
	trendList         []float64
	pricePollInterval int
	pricesToAverage   int
	averagesForTrend  int
	orderInterval     int
	config            Config
)

const (
	APIDomain = "https://api.kraken.com"
	tickerAPI = APIDomain + "/0/public/Ticker?pair="
	buyMode   = true
)

//{"error":[],"result":{"LINKEUR":{"a":["16.099230","789","789.000"],"b":["16.080190","301","301.000"],"c":["16.096100","26.00000000"],"v":["59933.65828008","371917.75198917"],"p":["16.093100","16.231972"],"t":[796,5286],"l":["15.700000","15.655200"],"h":["16.360670","16.930000"],"o":"15.865720"}}}
type Config struct {
	BotConfig struct {
		APIKey            string `yaml:"api_key"`
		PrivateKey        string `yaml:"private_key"`
		CurrencyPair      string `yaml:"currency_pair"`
		OrderInterval     string `yaml:"order_interval"`
		SellOrderSpread   string `yaml:"sell_order_spread"`
		BuyOrderSpread    string `yaml:"buy_order_spread"`
		BuySellSpread     string `yaml:"buy_sell_spread"`
		PricePollInterval string `yaml:"price_poll_interval"`
		PricesToAverage   string `yaml:"prices_to_average"`
		AveragesForTrend  string `yaml:"averages_for_trend"`
	} `yaml:"bot_config"`
}

// Since Go has no enum, create our own
type Trend int

const (
	TrendUp   Trend = 100
	TrendDown Trend = 1
)

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(config *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(config)
	if err != nil {
		processError(err)
	}
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func DoHTTPGet(uri string) (rawData []byte, err error) {
	response, err := http.Get(uri)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		rawData, _ = ioutil.ReadAll(response.Body)
	}
	return
}

func PrettyPrintJson(data []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(data), "", "    ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(out.String())
}

func CheckTrend() Trend {
	var marketTrend Trend
	groupTotal := 0.0
	groupAvg := 0.0
	counter := 0

	for i := range priceList {
		if counter < pricesToAverage-1 {
			groupTotal = groupTotal + priceList[i]
			if *debug {
				fmt.Printf("[DEBUG]  counter = %d | priceList[i] = %f | groupTotal = %f\n", counter, priceList[i], groupTotal)
			}
			counter = counter + 1
		} else {
			groupAvg = groupTotal / float64(pricesToAverage)
			trendList = append(trendList, groupAvg)
			if *debug {
				fmt.Printf("[DEBUG]  counter = %d | priceList[i] = %f | groupTotal = %f | groupAvg = %f \n", counter, priceList[i], groupTotal, groupAvg)
			}
			groupTotal = 0
			counter = 0
		}
	}

	if *debug {
		fmt.Printf("[DEBUG]  len(trendList) = %d | averagesForTrend = %d \n", len(trendList), averagesForTrend)
	}
	if len(trendList) > averagesForTrend {
		newTrendList := trendList[len(trendList)-averagesForTrend:]
		trendList = newTrendList
	}
	log.Println(trendList)

	marketTrend = TrendDown
	trendListLen := len(trendList)
	if trendListLen >= averagesForTrend {
		for i, v := range trendList {
			if *debug {
				fmt.Printf("[DEBUG]  %d - %f\n", i, v)
			}
			if i < trendListLen-1 {
				if v <= trendList[i+1] {
					if *debug {
						fmt.Println("[DEBUG]  CONTINUE")
					}
					continue
				} else {
					if *debug {
						fmt.Println("[DEBUG]  BREAK")
					}
					break
				}
			}
			marketTrend = TrendUp
		}
	}
	return marketTrend
}

func UpdateOrders() {

}

func main() {
	flag.Parse()
	if *debug {
		log.Printf("*** DEBUG MODE ***\n")
	}
	fmt.Println("Starting the application...")

	SetupCloseHandler()

	readFile(&config)
	fmt.Println(config.BotConfig.CurrencyPair)

	pricePollInterval, _ = strconv.Atoi(config.BotConfig.PricePollInterval)
	pricesToAverage, _ = strconv.Atoi(config.BotConfig.PricesToAverage)
	averagesForTrend, _ = strconv.Atoi(config.BotConfig.AveragesForTrend)
	orderInterval, _ = strconv.Atoi(config.BotConfig.OrderInterval)

	if pricesToAverage < 3 {
		pricesToAverage = 3
	}
	if averagesForTrend < 3 {
		averagesForTrend = 3
	}

	// var lastPurchasePrice float64
	// lastPurchasePrice = 0
	pollTime := time.Now()
	//lastTickerTime := time.Now()
	//lastTrendTime := time.Now()
	lastOrderTime := time.Now()

	var marketTrend Trend
	marketTrend = TrendDown

	trendCounter := 0

	priceListMaxLength := averagesForTrend * pricesToAverage
	//trendListLength := averagesForTrend

	for {
		oldTrend := marketTrend
		pollTime = time.Now()
		rawdata, err := DoHTTPGet(tickerAPI + config.BotConfig.CurrencyPair)
		if err != nil {
			continue
		}
		var kvData map[string]interface{}
		json.Unmarshal(rawdata, &kvData)

		if *jsonout {
			PrettyPrintJson(rawdata)
		}

		pairData := kvData["result"].(map[string]interface{})
		for key, value := range pairData {
			// Each value is an interface{} type, that is type asserted as a string
			pd := value.(map[string]interface{})
			priceData := pd["c"].([]interface{})
			currentPrice, err = strconv.ParseFloat(priceData[0].(string), 64)

			// Update priceList
			{
				if *debug {
					fmt.Printf("[DEBUG]  len(priceList) = %d | ", len(priceList))
				}
				priceList = append(priceList, currentPrice)
				if *debug {
					fmt.Printf("[DEBUG]  len(priceList) = %d | ", len(priceList))
				}
				if len(priceList) > priceListMaxLength {
					newPriceList := priceList[len(priceList)-priceListMaxLength:]
					priceList = newPriceList
				}
				if *debug {
					fmt.Printf("[DEBUG]  len(priceList) = %d \n", len(priceList))
				}
			}

			// Update trendList
			if *debug {
				fmt.Printf("[DEBUG]  trendCounter = %d | pricesToAverage = %d \n", trendCounter, pricesToAverage)
			}
			if trendCounter == pricesToAverage {
				if *debug {
					fmt.Printf("[DEBUG]  len(priceList) = %d | (averagesForTrend * pricesToAverage) = %d \n", len(priceList), (averagesForTrend * pricesToAverage))
				}
				if len(priceList) >= (averagesForTrend * pricesToAverage) {
					if *debug {
						log.Printf("[DEBUG]  Length of priceList: %d\n", len(priceList))
					}
					marketTrend = CheckTrend()
				}
			}

			// Update orders
			{
				elapsed := pollTime.Sub(lastOrderTime)
				elapsedInSeconds := elapsed / 1000000000

				if elapsedInSeconds > time.Duration(orderInterval) {
					log.Println("Update orders")
					if marketTrend == TrendUp {
						log.Println("Update sell orders")
					}
					lastOrderTime = pollTime
				}
			}

			var out string
			if currentPrice == lastPrice {
				out = fmt.Sprintf("  =  %s : %s", key, pd["c"])
			} else if currentPrice > lastPrice {
				out = fmt.Sprintf("  ^  %s : %s", key, pd["c"])
			} else {
				out = fmt.Sprintf("  V  %s : %s", key, pd["c"])
			}
			log.Println(out)
			lastPrice = currentPrice
		}

		if trendCounter == pricesToAverage {
			trendCounter = 0
		}
		trendCounter = trendCounter + 1

		if oldTrend != marketTrend {
			trendDirection := "DOWN"
			if marketTrend == TrendUp {
				trendDirection = "UP"
			}
			log.Printf("New market trend detected: %s", trendDirection)
		}
		if *debug {
			fmt.Println("")
		}

		//fmt.Printf("len=%d cap=%d %v\n", len(priceList), cap(priceList), priceList)
		time.Sleep(time.Duration(pricePollInterval) * time.Second)
	}

	fmt.Println("Terminating the application...")
}
