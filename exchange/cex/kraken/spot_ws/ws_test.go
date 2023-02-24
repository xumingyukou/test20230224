package spot_ws

import (
	"clients/exchange/cex/base"
	"context"
	"encoding/json"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"hash/crc32"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"testing"
)

var (
	b *KKWebsocket
)

func init() {
	//config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:7890",
		EndPoint: WS_PUBLIC_BASE_URL,
	}

	con := base.WsConf{
		APIConf: conf,
	}
	b = NewKKWebsocket(con)
}

func TestDepthIncre(t *testing.T) {
	b.ProxyUrl = "http://127.0.0.1:7890"
	ctx, _ := context.WithCancel(context.Background())
	symbol := client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   1,
	}
	symbol_ := client.SymbolInfo{
		Symbol: "LTC/USDT",
		Type:   1,
	}
	depthMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	ch1 := make(chan *client.WsDepthRsp)
	ch2 := make(chan *client.WsDepthRsp)
	depthMap[&symbol] = ch1
	depthMap[&symbol_] = ch2

	b.DepthIncrementGroup(ctx, 10, depthMap)
	//if err != nil {
	//	t.Fatal(err)
	//}
	for {
		select {
		case res, _ := <-ch1:
			fmt.Println("ch1 res", res)
		case res, _ := <-ch2:
			fmt.Println("ch2  res", res)
		}
	}
}

func TestTradeGroup(t *testing.T) {
	b.ProxyUrl = "http://127.0.0.1:7890"
	ctx, _ := context.WithCancel(context.Background())
	symbol := client.SymbolInfo{
		Symbol: "BTC/USD",
		Type:   1,
	}
	symbol_ := client.SymbolInfo{
		Symbol: "LTC/USD",
		Type:   1,
	}
	tradeMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	ch1 := make(chan *client.WsTradeRsp)
	ch2 := make(chan *client.WsTradeRsp)
	tradeMap[&symbol] = ch2
	tradeMap[&symbol_] = ch1

	err := b.TradeGroup(ctx, tradeMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, _ := <-ch1:
			fmt.Println("ch1 res", res)
		case res, _ := <-ch2:
			fmt.Println("ch2  res", res)
		}
	}
}

func TestTickerGroup(t *testing.T) {
	b.ProxyUrl = "http://127.0.0.1:7890"
	ctx, _ := context.WithCancel(context.Background())
	symbol := client.SymbolInfo{
		Symbol: "BTC/USD",
		Type:   1,
	}
	symbol_ := client.SymbolInfo{
		Symbol: "LTC/USD",
		Type:   1,
	}
	bookMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	ch1 := make(chan *client.WsBookTickerRsp)
	ch2 := make(chan *client.WsBookTickerRsp)
	bookMap[&symbol] = ch2
	bookMap[&symbol_] = ch1

	err := b.BookTickerGroup(ctx, bookMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, _ := <-ch1:
			fmt.Println("ch1 res", res)
		case res, _ := <-ch2:
			fmt.Println("ch2  res", res)
		}
	}
}

func TestSnapShotGroup(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	ch1 := make(chan *client.WsDepthRsp)
	ch2 := make(chan *depth.Depth)
	chDelatMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	symbol1 := client.SymbolInfo{
		//Symbol: "ZRX/BTC",

		Symbol: "BTC/USD",
		//Type:   1,
	}

	chDelatMap[&symbol1] = ch1
	chFullMap[&symbol1] = ch2

	symbol2 := client.SymbolInfo{
		Symbol: "BTC/EUR",
	}
	ch3 := make(chan *client.WsDepthRsp)
	ch4 := make(chan *depth.Depth)
	chDelatMap[&symbol2] = ch3
	chFullMap[&symbol2] = ch4
	ctx, _ := context.WithCancel(context.Background())

	err := b.DepthIncrementSnapshotGroup(ctx, 0, 1000, true, true, chDelatMap, chFullMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res := <-ch1:
			tmp, _ := json.Marshal(res)
			fmt.Println("增量ch1 res", string(tmp))
		case res := <-ch2:
			tmp, _ := json.Marshal(res)
			fmt.Println("全量i111212  res", string(tmp))
		case res := <-ch3:
			tmp, _ := json.Marshal(res)
			fmt.Println("ch3 res", string(tmp))
		case res := <-ch4:
			tmp, _ := json.Marshal(res)
			fmt.Println("ch4 i111212  res", string(tmp))

		}
	}
}

func s(symbols []string) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:7890",
		EndPoint: WS_PUBLIC_BASE_URL,
	}

	con := base.WsConf{
		APIConf: conf,
	}
	b = NewKKWebsocket(con)
	ctx, _ := context.WithCancel(context.Background())
	chDelatMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	for i := 0; i < len(symbols); i += 10 {
		for j := i; j < i+10 && j < len(symbols); j++ {

			symbol := client.SymbolInfo{Symbol: strings.ReplaceAll(symbols[j], "XBT", "BTC")}
			chDelatMap[&symbol] = nil
			chFullMap[&symbol] = nil

		}

	}
	b.DepthIncrementSnapshotGroup(ctx, 0, 1000, false, false, chDelatMap, chFullMap)
}
func TestCheck(t *testing.T) {
	// var fields []string
	// asks := [][]float64{{0.05005, 0.00000500, 1582905487.684110},
	// 	{0.05010, 0.00000500, 1582905486.187983},
	// 	{0.05015, 0.00000500, 1582905484.480241},
	// 	{0.05020, 0.00000500, 1582905486.645658},
	// 	{0.05025, 0.00000500, 1582905486.859009},
	// 	{0.05030, 0.00000500, 1582905488.601486},
	// 	{0.05035, 0.00000500, 1582905488.357312},
	// 	{0.05040, 0.00000500, 1582905488.785484},
	// 	{0.05045, 0.00000500, 1582905485.302661},
	// 	{0.05050, 0.00000500, 1582905486.157467}}
	// bids := [][]float64{
	// 	{0.05000, 0.00000500, 1582905487.439814},
	// 	{0.04995, 0.00000500, 1582905485.119396},
	// 	{0.04990, 0.00000500, 1582905486.432052},
	// 	{0.04980, 0.00000500, 1582905480.609351},
	// 	{0.04975, 0.00000500, 1582905476.793880},
	// 	{0.04970, 0.00000500, 1582905486.767461},
	// 	{0.04965, 0.00000500, 1582905481.767528},
	// 	{0.04960, 0.00000500, 1582905487.378907},
	// 	{0.04955, 0.00000500, 1582905483.626664},
	// 	{0.04950, 0.00000500, 1582905488.509872},
	// }
	// cs_ := crc32.ChecksumIEEE([]byte("50055005010500501550050205005025500503050050355005040500504550050505005000500499550049905004980500497550049705004965500496050049555004950500"))
	// for i := 0; i < 10; i++ {
	// 	if i < 10 {
	// 		// fields = append(fields, ]), float2M(asks[i][1]))
	// 	}
	// 	fmt.Println(bids, asks)
	// 	fmt.Println(fields)
	// }
	// for i := 0; i < 10; i++ {
	// 	if i < 10 {
	// 		// fields = append(fields, ]), float2M(bids[i][1]))
	// 	}
	// 	fmt.Println(fields)
	// }
	// fmt.Println(len(fields))
	// raw := strings.Join(fields, "")
	// raw = strings.ReplaceAll(raw, ".", "")
	dep := int64(866411700)
	dep = 866411700
	fmt.Println(int32(dep))
	cs := crc32.ChecksumIEEE([]byte("238968000024786226723897300003000000023897800001044543523897900001850000023900300001850000023900400006274903742390400000185000002390450000145062552390650000104555180423908000008094232389670000321942388920000149600122388350000185000002388130000208993823881100001850000023878700001850000023876300001850000023876100007926732387590000628133110238753000030000000"))
	fmt.Println(cs)
	fmt.Println(int32(cs))
	// fmt.Println(int32(cs) == int32(cs_))
	// fmt.Println(fmt.Sprintf("%.5f", 12.20))
	//fmt.Println(float2M(0.00000500))
	//2405970000623348102406280000148599622406340000200000024066200001850000024066800001038227524067400006624000024067500006231478232406900000185000002407020000300000002407160000185000002405650000273000002405630000623434972405540000190251852405400000382002405310000200000024053000001152358240529000018500000240507000051240000240506000014485319240504000018500000
	//24059700006233481024062800001485996224063400002000000240662000018500000240668000010382275240674000066240000240675000062314782324069000001850000024070200003000000024071600001850000024056500002730000024056300006234349724055400001902518524054000003820024053000001152358240529000018500000240507000051240000240506000014485319240504000018500000240503000015000000

	//238968000024786226723897300003000000023897800001044543523897900001850000023900300001850000023900400006274903742390400000185000002390450000145062552390650000104555180423908000008094232389670000321942388920000149600122388350000185000002388130000208993823881100001850000023878700001850000023876300001850000023876100007926732387590000628133110238753000030000000
	//238968000024786226723897300003000000023897800001044543523897900001850000023900300001850000023900400006274903742390400000185000002390450000145062552390650000104555180423908000008094232389670000321942388920000149600122388370000627928102388350000185000002388130000208993823881100001850000023878700001850000023876300001850000023876100007926732387590000628133110
}

func TestCheckSum(t *testing.T) {
	//bid := "(24090.200000 ,0.622481)(24091.600000 ,0.185000)(24094.700000 ,0.149793)(24094.800000 ,6.223668)(24096.000000 ,0.020000)(24098.200000 ,0.150000)(24099.500000 ,0.515600)(24100.800000 ,0.300000)(24102.400000 ,0.043335)(24103.400000 ,0.086671)(24107.800000 ,0.234660)(24108.200000 ,0.008094)(24109.000000 ,0.516500)(24109.100000 ,10.366598)(24111.600000 ,0.100000)(24111.700000 ,10.365493)(24112.400000 ,0.103000)"
	//ask := "(24086.800000 ,0.622566)(24086.200000 ,0.088000)(24085.300000 ,0.020000)(24085.000000 ,0.185000)(24082.900000 ,0.103787)(24082.800000 ,0.150000)(24082.500000 ,0.323947)(24081.600000 ,0.515600)(24081.500000 ,0.300000)(24079.900000 ,0.220670)(24079.700000 ,0.220670)(24079.600000 ,0.220670)(24079.500000 ,0.185000)(24078.500000 ,0.220670)(24078.400000 ,0.300000)(24078.200000 ,0.515600)(24078.000000 ,6.227977)(24077.600000 ,0.515600)(24077.300000 ,0.515600)(24076.900000 ,0.185000)(24074.500000 ,0.234660)(24074.200000 ,0.185000)(24073.800000 ,0.043345)(24072.700000 ,10.382211)(24072.500000 ,0.363870)(24071.700000 ,0.185000)"

	ask := "24087.0000000000 ,0.6225699900)(24089.7000000000 ,0.1360409800)(24092.9000000000 ,0.0080942300)(24095.0000000000 ,6.2236153200)(24100.5000000000 ,0.0433390900)(24101.6000000000 ,0.0866781800)(24103.6000000000 ,10.3690203300)(24105.7000000000 ,0.0930000000)(24106.0000000000 ,10.3680277900)(24108.3000000000 ,0.2346600000)(24113.9000000000 ,0.2500000000"
	bid := "24081.7000000000 ,0.1850000000)(24080.5000000000 ,0.1850000000)(24080.1000000000 ,0.1850000000)(24077.8000000000 ,0.5356000000)(24077.6000000000 ,0.1850000000)(24076.2000000000 ,0.1038434500)(24074.9000000000 ,0.1500000000)(24073.3000000000 ,0.1850000000)(24072.8000000000 ,0.3000000000)(24070.6000000000 ,0.2206699000)(24070.5000000000 ,0.1850000000)(24070.4000000000 ,0.0077766700)(24069.0000000000 ,0.5346600000)(24068.6000000000 ,6.2304332800)(24068.0000000000 ,0.1850000000)(24066.5000000000 ,0.0433536100)(24065.5000000000 ,0.1951054200"

	content := ""
	for _, i := range strings.Split(ask, ")(")[:10] {
		j := strings.Split(i, ",")
		price := strings.TrimLeft(strings.ReplaceAll(j[0][:len(j[0])-6], ".", ""), "0")
		amount := strings.TrimLeft(strings.ReplaceAll(j[1][:len(j[1])-2], ".", ""), "0")
		content += fmt.Sprint(price, amount)
	}
	for _, i := range strings.Split(bid, ")(")[:10] {
		j := strings.Split(i, ",")
		price := strings.TrimLeft(strings.ReplaceAll(j[0][:len(j[0])-6], ".", ""), "0")
		amount := strings.TrimLeft(strings.ReplaceAll(j[1][:len(j[1])-2], ".", ""), "0")
		content += fmt.Sprint(price, amount)
	}
	fmt.Println(content == "240870000062256999240897000013604098240929000080942324095000006223615322410050000433390924101600008667818241036000010369020332410570000930000024106000001036802779241083000023466000240817000018500000240805000018500000240801000018500000240778000053560000240776000018500000240762000010384345240749000015000000240733000018500000240728000030000000240706000022066990")
}

func TestChecksum(t *testing.T) {
	//asks := [][]float64{{23961.5, 1.03519746}, {23962.2, 3.1299329}, {23962.6, 0.73000456}, {23964.4, 0.0001234}, {23964.5, 0.41143185}, {23965, 0.0001}, {23965.1, 0.7}, {23966, 0.36}, {23966.1, 0.5546846}, {23966.3, 0.03058149}, {23966.5, 1.73845542}, {23966.6, 6.25871993}, {23967.3, 2.17902017}, {23967.8, 6.25839787}, {23968.1, 0.25}, {23968.3, 0.7}, {23968.4, 2.5411}, {23968.8, 0.22889068}, {23969.9, 0.7583068}, {23970, 6.2578287}, {23971, 1.254}, {23973.4, 0.33}, {23973.5, 0.5}, {23973.6, 0.14417226}, {23973.8, 3.05451855}, {23974.1, 1.668}, {23974.5, 0.7}}
	//bisd := [][]float64{}
}
