package ok_api

import (
	"clients/transform"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/order"

	"github.com/warmplanet/proto/go/common"
)

func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sha)
}

func ComputeHmacSha256Ex(secret string, items ...string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	for _, v := range items {
		h.Write([]byte(v))
	}
	sha := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sha)
}

func S2M(i interface{}, n int) []map[string]string {
	var m = make([]map[string]string, n)
	j, _ := json.Marshal(i)
	_ = json.Unmarshal(j, &m)
	//for _, j := range m {
	//	fmt.Println("after", j)
	//}
	return m
}

func GetNetWorkFromChain(chain common.Chain) string {
	switch chain {
	case common.Chain_ETH:
		return "ERC20"
	case common.Chain_TRON:
		return "TRC20"
	case common.Chain_AVALANCHE:
		return "Avalanche C-Chain"
	case common.Chain_SOLANA:
		return "Solana"
	case common.Chain_POLYGON:
		return "Polygon"
	case common.Chain_OKC:
		return "OKC"
	case common.Chain_ARBITRUM:
		return "Arbitrum one"
	case common.Chain_OPTIMISM:
		return "Optimism"
	case common.Chain_BTC:
		return "Bitcoin"
	default:
		return ""
	}
}
func ParseSymbolName(symbol string) string {
	s := strings.ReplaceAll(symbol, "-", "/")
	return s
}

func GetOrderStatusFromExchange(status OrderStatus) order.OrderStatusCode {
	switch status {
	case ORDER_STATUS_LIVE:
		return order.OrderStatusCode_OPENED
	case ORDER_STATUS_NEW:
		return order.OrderStatusCode_OPENED
	case ORDER_STATUS_PARTIALLY_FILLED:
		return order.OrderStatusCode_PARTFILLED
	case ORDER_STATUS_FILLED:
		return order.OrderStatusCode_FILLED
	case ORDER_STATUS_CANCELED:
		return order.OrderStatusCode_CANCELED
	case ORDER_STATUS_PENDING_CANCEL:
		return order.OrderStatusCode_CANCELED
	case ORDER_STATUS_REJECTED:
		return order.OrderStatusCode_FAILED
	case ORDER_STATUS_EXPIRED:
		return order.OrderStatusCode_EXPIRED
	default:
		return order.OrderStatusCode_OrderStatusInvalid
	}
}

func TransformMarket(instType string) common.Market {
	switch instType {
	case "FUTURES":
		return common.Market_FUTURE
	case "MARGIN":
		return common.Market_MARGIN
	default:
		return common.Market_INVALID_MARKET
	}
}

func GetSymbol(symbol string) string {
	symbol = strings.Replace(symbol, "/", "", 1)
	symbol = strings.Replace(symbol, "-", "", 1)
	symbol = strings.Replace(symbol, "_", "", 1)
	return strings.ToUpper(symbol)
}

func GetChainFromNetWork(network string) common.Chain {
	i := strings.Index(network, "-")
	chain := network[i+1:]

	switch chain {
	case "ERC20":
		return common.Chain_ETH
	case "TRC20":
		return common.Chain_TRON
	case "Avalanche C-Chain":
		return common.Chain_AVALANCHE
	case "Solana":
		return common.Chain_SOLANA
	case "Polygon":
		return common.Chain_POLYGON
	case "Arbitrum one":
		return common.Chain_ARBITRUM
	case "Optimism":
		return common.Chain_OPTIMISM
	case "Bitcoin":
		return common.Chain_BTC
	default:
		return common.Chain_INVALID_CAHIN
	}
}

func GetSideTypeToExchange(side order.TradeSide) SideType {
	switch side {
	case order.TradeSide_BUY:
		return SIDE_TYPE_BUY
	case order.TradeSide_SELL:
		return SIDE_TYPE_SELL
	default:
		return ""
	}
}

func GetTransferTypeToExchange(status client.TransferStatus) TransferStatus {
	switch status {
	case client.TransferStatus_TRANSFERSTATUS_CREATED:
		return TRANSFER_TYPE_CREATED
	case client.TransferStatus_TRANSFERSTATUS_CANCELLED:
		return TRANSFER_TYPE_CANCELLED
	case client.TransferStatus_TRANSFERSTATUS_CONFORMING:
		return TRANSFER_TYPE_CONFIRMING
	case client.TransferStatus_TRANSFERSTATUS_REJECTED:
		return TRANSFER_TYPE_REJECTED
	case client.TransferStatus_TRANSFERSTATUS_PROCESSING:
		return TRANSFER_TYPE_PROCESSING
	case client.TransferStatus_TRANSFERSTATUS_FAILED:
		return TRANSFER_TYPE_FAILED
	case client.TransferStatus_TRANSFERSTATUS_COMPLETE:
		return TRANSFER_TYPE_SUCCESS
	default:
		return -1
	}
}

func GetDepositTypeFromExchange(status DepositStatus) client.DepositStatus {
	switch status {
	case DOPSIT_TYPE_PENDING:
		return client.DepositStatus_DEPOSITSTATUS_PENDING
	case DOPSIT_TYPE_PAUSE:
		return client.DepositStatus_DEPOSITSTATUS_PENDING
	case DOPSIT_TYPE_CONFIRMED:
		return client.DepositStatus_DEPOSITSTATUS_CONFIRMED
	case DOPSIT_TYPE_SUCCESS:
		return client.DepositStatus_DEPOSITSTATUS_SUCCESS
	case DOPSIT_TYPE_LOCK:
		return client.DepositStatus_DEPOSITSTATUS_FAILED
	case DOPSIT_TYPE_INTERCEPT:
		return client.DepositStatus_DEPOSITSTATUS_FAILED
	default:
		return client.DepositStatus_DEPOSITSTATUS_INVALID
	}
}

func GetDate(contract transform.CONTRACTTYPE) string {
	// 当周、次周、当季度、次季度
	now := time.Now().In(transform.BeiJingLoc)

	switch contract {
	case transform.THISWEEK:
		return transform.GetThisWeek(now, 5, 16).Format("060102")
	case transform.NEXTWEEK:
		return transform.GetNextWeek(now, 5, 16).Format("060102")
	case transform.THISMONTH:
		return transform.GetThisMonth(now, 5, 16).Format("060102")
	case transform.NEXTMONTH:
		return transform.GetNextMonth(now, 5, 16).Format("060102")
	case transform.THISQUARTER:
		thisWeek := transform.GetThisWeek(now, 5, 16).Format("060102")
		nextWeek := transform.GetNextWeek(now, 5, 16).Format("060102")
		thisQuarter := transform.GetThisQuarter(now, 5, 16).Format("060102")
		if thisQuarter == thisWeek || thisQuarter == nextWeek {
			now = now.AddDate(0, 0, 20)
			return transform.GetThisQuarter(now, 5, 16).Format("060102")
		}
		return thisQuarter
	case transform.NEXTQUARTER:
		thisWeek := transform.GetThisWeek(now, 5, 16).Format("060102")
		nextWeek := transform.GetNextWeek(now, 5, 16).Format("060102")
		thisQuarter := transform.GetThisQuarter(now, 5, 16).Format("060102")
		if thisQuarter == thisWeek || thisQuarter == nextWeek {
			now = now.AddDate(0, 0, 20)
		}
		return transform.GetNextQuarter(now, 5, 16).Format("060102")
	default:
		return ""
	}
}
