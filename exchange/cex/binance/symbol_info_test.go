package binance

import (
	"clients/config"
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance/spot_api"
	"clients/exchange/cex/binance/u_api"
	"clients/logger"
	"strconv"
	"strings"
	"testing"
)

func TestSpotSymbolInfo(t *testing.T) {
	var timeOffset int64 = 30
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		ProxyUrl:    ProxyUrl,
		AccessKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		SecretKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	}
	a := spot_api.NewApiClient(conf)
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

func TestUBaseSymbolInfo(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	var (
		timeOffset int64 = 30
		spotConf         = config.ExchangeSpotConfig{}
		conf             = base.APIConf{
			ReadTimeout: timeOffset,
			ProxyUrl:    ProxyUrl,
			AccessKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
			SecretKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
		}
		symbols              = []string{"BTC/USDT", "LTC/USDT", "ETH/USDT", "LINK/USDT", "BCH/USDT", "XRP/USDT", "EOS/USDT", "TRX/USDT", "ETC/USDT", "DOT/USDT", "ADA/USDT", "BNB/USDT", "FIL/USDT", "UNI/USDT", "XLM/USDT", "DOGE/USDT"}
		a                    = u_api.NewUApiClient(conf)
		exchangeInfoRes, err = a.ExchangeInfo()
		symbolMap            = make(map[string]string)
		symbolInfoMap        = make(map[string]*config.SymbolInfo)
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, sym := range symbols {
		symbolMap[strings.Replace(sym, "/", "", -1)] = sym
	}
	for _, symbol := range exchangeInfoRes.Symbols {
		if _, ok := symbolMap[symbol.Symbol]; !ok {
			continue
		}
		precision := symbol.GetPrecision()
		info := &config.SymbolInfo{
			Symbol: symbol.BaseAsset + "/" + symbol.QuoteAsset,
			Precision: config.SymbolPrecision{
				Amount:    int64(precision.AmountPrecision),
				Price:     int64(precision.PricePrecision),
				AmountMin: precision.MinAmount,
			},
		}
		symbolInfoMap[symbol.BaseAsset+"/"+symbol.QuoteAsset] = info
		spotConf.Symbols = append(spotConf.Symbols, &config.SymbolInfo{
			Symbol: symbol.BaseAsset + "/" + symbol.QuoteAsset,
			//Type:      u_api.GetFutureTypeFromExchange(u_api.ContractType(symbol.ContractType)),
			Precision: info.Precision,
			TradeFee:  info.TradeFee,
		})
		var (
			tradeFee           *u_api.RespCommissionRate
			takerFee, makerFee float64
		)
		tradeFee, err = a.CommissionRate(u_api.GetUBaseSymbol(symbolMap[symbol.Symbol], u_api.ContractType(symbol.ContractType)))
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := symbolMap[tradeFee.Symbol]; !ok {
			continue
		}
		takerFee, err = strconv.ParseFloat(tradeFee.TakerCommissionRate, 64)
		if err != nil {
			t.Fatal(err)
		}
		makerFee, err = strconv.ParseFloat(tradeFee.MakerCommissionRate, 64)
		if err != nil {
			t.Fatal(err)
		}
		symbolInfoMap[symbolMap[symbol.Symbol]].TradeFee = config.TradeFeeItem{
			Taker: takerFee,
			Maker: makerFee,
		}
	}

	logger.SaveToFile("binance_ubase.json", spotConf)
}
