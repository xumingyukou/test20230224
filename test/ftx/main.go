package main

import (
	"clients/crypto"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ENDPOINT               = "https://ftx.com/api/"
	access_key, secret_key = getKey()
	proxyUrl, _            = url.Parse("http://127.0.0.1:7890")
	transport              = http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	httpClient = &http.Client{
		//Transport: &transport,
		Timeout: 10 * time.Second,
	}
	subaccount string
)

func getKey() (string, string) {
	filePath := "../.ftx/key.json"
	data := make(map[string]string)
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}
	return data["access_key"], data["secret_key"]
}

func DoRequest(method, uri string, body []byte, params *map[string]string, resp interface{}) error {
	// avoid Pointer's butting
	u, _ := url.ParseRequestURI(ENDPOINT)
	fmt.Println("U:", u)
	u.Path = u.Path + uri
	fmt.Println("U2:", u)
	if params != nil {
		q := u.Query()
		for k, v := range *params {
			q.Set(k, v)
		}
		fmt.Println("U3:", q)
		u.RawQuery = q.Encode()
	}

	nonce := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	fmt.Println("Time:", strconv.FormatInt(time.Now().UnixMilli(), 10))
	fmt.Println("Nonce:", nonce)
	var q string
	if u.RawQuery != "" {
		q = "?" + u.Query().Encode()
	}
	payload := nonce + method + "/api/" + uri + q
	fmt.Println("U4:", payload)
	if body != nil {
		payload += string(body)
	}
	fmt.Println("U5:", body)
	signature, err := crypto.GetParamHmacSHA256Sign(secret_key, payload)
	fmt.Println("Signature: ", signature)
	if err != nil {
		return err
	}
	fmt.Println("request:", method, u.Path, payload, signature)
	req, err := http.NewRequest(method, u.String(), strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTX-KEY", access_key)
	req.Header.Set("FTX-SIGN", signature)
	req.Header.Set("FTX-TS", nonce)
	if subaccount != "" {
		req.Header.Set("FTX-SUBACCOUNT", subaccount)
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			return
		}
	}()
	resBytes, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(resBytes))
	if res.StatusCode != 200 {
		//c.Logger.Printf("status: %s", res.Status)
		return errors.New(res.Status)
	}
	err = json.Unmarshal(resBytes, resp)
	return err
}

type RespAccount struct {
	Success bool `json:"success"`
	Result  struct {
		BackstopProvider             bool    `json:"backstopProvider"`
		Collateral                   float64 `json:"collateral"`
		FreeCollateral               float64 `json:"freeCollateral"`
		InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
		Leverage                     float64 `json:"leverage"`
		Liquidating                  bool    `json:"liquidating"`
		MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
		MakerFee                     float64 `json:"makerFee"`
		MarginFraction               float64 `json:"marginFraction"`
		OpenMarginFraction           float64 `json:"openMarginFraction"`
		TakerFee                     float64 `json:"takerFee"`
		TotalAccountValue            float64 `json:"totalAccountValue"`
		TotalPositionSize            float64 `json:"totalPositionSize"`
		Username                     string  `json:"username"`
		Positions                    []struct {
			Cost                         float64 `json:"cost"`
			EntryPrice                   float64 `json:"entryPrice"`
			Future                       string  `json:"future"`
			InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
			LongOrderSize                float64 `json:"longOrderSize"`
			MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
			NetSize                      float64 `json:"netSize"`
			OpenSize                     float64 `json:"openSize"`
			RealizedPnl                  float64 `json:"realizedPnl"`
			ShortOrderSize               float64 `json:"shortOrderSize"`
			Side                         string  `json:"side"`
			Size                         float64 `json:"size"`
			UnrealizedPnl                float64 `json:"unrealizedPnl"`
		} `json:"positions"`
	} `json:"result"`
}

func main() {
	var (
		url    = "account"
		params = make(map[string]string)
		data   = []byte("")
		resp   RespAccount
	)
	err := DoRequest("GET", url, data, &params, &resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
