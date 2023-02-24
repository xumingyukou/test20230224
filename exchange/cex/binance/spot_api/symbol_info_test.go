package spot_api

import (
	"clients/config"
	"clients/exchange/cex/base"
	"clients/logger"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSymbolInfo(t *testing.T) {
	var timeOffset int64 = 30
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		AccessKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		SecretKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	}
	a := NewApiClient(conf)
	exchangeInfoRes, err := a.ExchangeInfo()
	if err != nil {
		t.Fatal(err)
	}
	configAll, err := a.CapitalConfigGetAll()
	if err != nil {
		t.Fatal(err)
	}
	trakeFeeRes, err := a.AssetTradeFee()
	if err != nil {
		t.Fatal(err)
	}

	symbols := []string{"BAKE/USDT", "AVAX/USDT", "BTT/USDT", "C98/USDT", "AXS/USDT", "SFP/USDT", "ALPHA/USDT", "ALICE/USDT", "DODO/USDT", "TLM/USDT", "BEL/USDT", "SXP/USDT", "BAND/USDT", "REEF/USDT", "IOTX/USDT", "FTM/USDT", "COMP/USDT", "MATIC/USDT", "CHR/USDT", "ZIL/USDT", "LINA/USDT", "MASK/USDT", "ONE/USDT", "ATA/USDT", "IOTA/USDT", "NEAR/USDT", "CELR/USDT", "ZEC/USDT", "BAT/USDT", "BLZ/USDT", "AAVE/USDT", "SNX/USDT", "YFI/USDT", "CAKE/USDT", "EOS/USDT", "ETH/USDT", "LTC/USDT", "ETC/USDT", "BTC/USDT", "BNB/USDT", "XRP/USDT", "FIL/USDT", "BCH/USDT", "TRX/USDT", "UNI/USDT", "SUSHI/USDT", "XLM/USDT", "DOT/USDT", "LINK/USDT", "ADA/USDT", "USDC/USDT", "BUSD/USDT", "DOGE/USDT", "MBOX/USDT", "ALPACA/USDT", "XVS/USDT", "RAY/USDT", "SOL/USDT", "SRM/USDT"}
	tokens := []string{"LTC", "BTC", "ETH", "ETC", "BCH", "XRP", "USDT", "BTG", "EOS", "XLM"}
	symbolMap := make(map[string]string)
	tokenMap := make(map[string]bool)
	for _, sym := range symbols {
		symbolMap[strings.Replace(sym, "/", "", -1)] = sym
	}
	for _, token := range tokens {
		tokenMap[token] = true
	}
	spotConf := config.ExchangeSpotConfig{}
	var (
		symbolInfoMap = make(map[string]*config.SymbolInfo)
		tokenInfoMap  = make(map[string]*config.TokenInfo)
	)
	for _, symbol := range exchangeInfoRes.Symbols {
		if _, ok := symbolMap[symbol.Symbol]; !ok {
			continue
		}
		precision := symbol.GetPrecision()

		symbolInfoMap[symbol.BaseAsset+"/"+symbol.QuoteAsset] = &config.SymbolInfo{
			Symbol: symbol.BaseAsset + "/" + symbol.QuoteAsset,
			Precision: config.SymbolPrecision{
				Amount:    int64(precision.AmountPrecision),
				Price:     int64(precision.PricePrecision),
				AmountMin: precision.MinAmount,
			},
		}
	}
	for _, tradeFee := range *trakeFeeRes {
		var (
			takerFee, makerFee float64
		)
		if _, ok := symbolMap[tradeFee.Symbol]; !ok {
			continue
		}
		takerFee, err = strconv.ParseFloat(tradeFee.TakerCommission, 64)
		if err != nil {
			t.Fatal(err)
		}
		makerFee, err = strconv.ParseFloat(tradeFee.MakerCommission, 64)
		if err != nil {
			t.Fatal(err)
		}
		symbolInfoMap[symbolMap[tradeFee.Symbol]].TradeFee = config.TradeFeeItem{
			Taker: takerFee,
			Maker: makerFee,
		}
	}

	for _, item := range *configAll {
		if _, ok := tokenMap[item.Coin]; !ok {
			continue
		}
		for _, network := range item.NetworkList {
			var (
				fee float64
			)
			fee, err = strconv.ParseFloat(network.WithdrawFee, 64)
			if err != nil {
				t.Fatal(err)
			}
			if _, ok := tokenInfoMap[item.Coin]; !ok {
				tokenInfoMap[item.Coin] = &config.TokenInfo{
					Token: item.Coin,
				}
			}
			//tokenInfoMap[item.Coin].WithdrawFee[network.Network] = fee
			tokenInfoMap[item.Coin].WithdrawFee = append(tokenInfoMap[item.Coin].WithdrawFee, &config.TransferFeeItem{
				NetWork: network.Network,
				Fee:     fee,
			})
		}
	}
	for symbol, info := range symbolInfoMap {
		spotConf.Symbols = append(spotConf.Symbols, &config.SymbolInfo{
			Symbol:    symbol,
			Precision: info.Precision,
			TradeFee:  info.TradeFee,
		})
	}
	for token, info := range tokenInfoMap {
		spotConf.Tokens = append(spotConf.Tokens, &config.TokenInfo{
			Token:       token,
			WithdrawFee: info.WithdrawFee,
		})
	}
	logger.SaveToFile("binance_spot.json", spotConf)
}

func TestGetSymbol(t *testing.T) {
	for _, i := range []string{
		"ETHBTC",
		"LTCBTC",
		"BNBBTC",
		"NEOBTC",
		"QTUMETH",
		"EOSETH",
		"SNTETH",
		"BNTETH",
		"BCCBTC",
		"GASBTC",
		"BNBETH",
		"BTCUSDT",
		"ETHUSDT",
		"HSRBTC",
		"OAXETH",
		"DNTETH",
		"MCOETH",
		"ICNETH",
		"MCOBTC",
		"WTCBTC",
		"LUNCBUSD",
		"USTCBUSD",
		"OPBTC",
		"OPBUSD",
		"OPUSDT",
		"OGBUSD",
		"KEYBUSD",
		"ASRBUSD",
		"FIROBUSD",
		"NKNBUSD",
		"OPBNB",
		"OPEUR",
		"GTOBUSD",
		"SNXETH",
		"WBTCBUSD",
	} {
		start := time.Now()
		ParseSymbolName(i)
		fmt.Println(time.Since(start).String())
	}

	m := make(map[string]string)

	for i := 0; i < 1000; i++ {
		m[strconv.Itoa(i)] = strconv.Itoa(i)
	}

	for i := 0; i < 1000; i++ {
		fmt.Println("------------------------")
		start := time.Now()
		_ = m[strconv.Itoa(i)]
		fmt.Println(time.Since(start).String())
	}

}

func TestSlice(t *testing.T) {
	a := []int{}
	for i := 0; i < 1000; i++ {
		a = append(a, i)
	}
	fmt.Println("start:", cap(a), len(a))
	fmt.Printf("%p", a)
	a = a[500:]
	fmt.Println("end", cap(a), len(a))
	fmt.Printf("%p", a)
}

func TestSymbol(t *testing.T) {
	for _, symbol := range []string{"BTCUSDT", "BTCUSD", "SDFSDTUSD", "ETHUSDP"} {
		fmt.Println(ParseSymbolName(symbol))
	}
}
