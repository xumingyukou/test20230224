import requests,json,time,hmac,hashlib
from operator import itemgetter


with open("../../../.binance/key.json") as f:
    content = f.read()
mykey = json.loads(content)
url = "https://fapi.binance.com/fapi/v1/order"
data = {
    "quantity":5,
    "recvWindow":60000,
    "side":"BUY",
    "symbol": "BTCUSDT",
    "type": "MARKEY",
    "timestamp": int(time.time()*1000),
}
headers = {
    "X-Mbx-Apikey":mykey["access_key"]
}

def signature(data):
    params = []
    for key, value in data.items():
        params.append((key, value))
    # sort parameters by key
    params.sort(key=itemgetter(0))
    ordered_data = params
    query_string = '&'.join(["{}={}".format(d[0], d[1]) for d in ordered_data])
    m = hmac.new(mykey["secret_key"].encode('utf-8'), query_string.encode('utf-8'), hashlib.sha256)
    return m.hexdigest()

data["signature"] = signature(data)


data = {
    "quantity":5,
    "recvWindow":60000,
    "side":"BUY",
    "symbol": "BTCUSDT",
    "type": "MARKEY",
    "timestamp": 1657166407573,
    "signature":"41c2fb7b4813855e62df788d72d0dc814659e75528c25bcef96843d542be9219",
}
data = {"timestamp":1657166407573}
print(data, headers)
res = requests.post(url, data=data, headers=headers)
print(res)
print(res.content.decode())
