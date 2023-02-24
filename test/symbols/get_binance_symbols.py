import requests,json

info = requests.get("https://api.binance.com/api/v1/exchangeInfo",proxies={"http":"http://127.0.0.1:7890"},timeout=5).json()
#content = "switch symbol{\n"
#for symbol in info["symbols"]:
#    content += f"case \"{symbol['symbol']}\":return \"{symbol['baseAsset']+'/'+symbol['quoteAsset']}\"\n"
#content += "default: return \"\"\n}"
#with open("binance_symbol.txt", "w") as f:
#    f.write(content)
symbol_list = {}
quote_list = ["USD"]
for symbol in info["symbols"]:
    symbol_list[symbol['symbol']] = symbol['baseAsset'] +'/'+ symbol['quoteAsset']
    quote_list.append(symbol['quoteAsset'])
print(json.dumps(list(set(quote_list))))

def get_symbol_name(symbol:str):
    symbol = symbol.upper()
    if symbol.endswith("USD"):
        if symbol.endswith("BUSD"):
            return symbol[:len(symbol)-4]+"/"+"BUSD"
        elif symbol.endswith("TUSD"):
            return symbol[:len(symbol)-4]+"/"+"TUSD"
        return symbol[:len(symbol)-3]+"/"+"USD"

    for quote in ["BUSD", "UST", "TUSD", "BKRW", "AUD", "DOT", "NGN", "BVND", "XRP", "BTC", "TRX", "UAH", "BIDR", "VAI", "DAI", "DOGE", "GBP", "BRL", "USDP", "USDS", "USDT", "IDRT", "EUR", "ETH", "USD", "PAX", "TRY", "BNB", "RUB", "ZAR", "USDC"]:
        if symbol.endswith(quote):
            return symbol[:len(symbol)-len(quote)]+"/"+quote
    return symbol


for e_sym, c_sym in symbol_list.items():
    if c_sym != get_symbol_name(e_sym):
        print(c_sym,e_sym,get_symbol_name(e_sym))

