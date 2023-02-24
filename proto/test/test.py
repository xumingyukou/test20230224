import asyncio
import sys
import time

sys.path.append('./py')

import nats.errors
from sdk import config
from sdk import logger
from sdk.broker import msg
from sdk.broker import subjects

from common import constant_pb2 as constant 
from order import order_pb2 as order
from depth import depth_pb2 as depth

async def testPubSub():
    c = config.LoadConfigFile('./conf/nats.toml')
    def handler(subject, data):
        print(subject, data.decode())
    subjects = ["dd.1.2.3.*", "dd.1.2.4.*"]
    sub = await msg.Sub.new(c, handler, True)
    p = await msg.Pub.new(c)

    for subject in subjects:
        await sub.subscribe(subject)

    i = 0
    while True:
        await asyncio.sleep(5)
        if i % 2 == 0:
            await p.publish("dd.1.2.3.ETH/USDT", bytes(str(i),encoding='utf8'))
            await p.publish("dd.1.2.4.ETH/USDT", bytes(str(i),encoding='utf8'))
        else:
            await sub.subscribe("dd.1.2.3.ETH/USDT")
            await sub.subscribe("dd.1.2.4.*")
        i = i + 1

async def testClientServer():
    c = config.LoadConfigFile('./conf/nats.toml')
    def handler(subject, data):
        input = int(data.decode())
        input = input + 1
        return bytes(str(input),encoding='utf8')

    subjects = ["dd.1.2.3.*", "dd.1.2.4.*"]
    server = await msg.Server.new(c, handler, True, subjects)
    p = await msg.Client.new(c)

    i = 0
    while True:
        await asyncio.sleep(2)
        r1 = None
        r2 = None
        try:
            r1 = await p.request("dd.1.2.3.ETH/USDT", bytes(str(i),encoding='utf8'), 1.0)
            r2 = await p.request("dd.1.2.4.ETH/USDT", bytes(str(i),encoding='utf8'), 1.0)
        except nats.errors.TimeoutError:
            print('Timeout error')

        if r1 != None:
            print("recv from dd.1.2.3.ETH/USDT " + r1.decode())

        if r2 != None:
            print("recv from dd.1.2.4.ETH/USDT " + r2.decode())

        i = i + 1

if __name__ == "__main__":
    a = constant.Exchange
    dd = depth.Depth()
    dd.timestamp = int(time.time())
    print(dd)

    #asyncio.run(testPubSub())
    asyncio.run(testClientServer())