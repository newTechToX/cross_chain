from collections import defaultdict
import json
import logging
import requests
from sqlalchemy import *
from sqlalchemy.orm import *
from tqdm import tqdm
from eth_abi import decode_single
from time import sleep
import traceback

engine = create_engine("postgresql://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain", future=True)
Session = sessionmaker(bind=engine, future=True)
metadata = MetaData()


def os_search(body, index="eth_block"):
    TIMEOUT = 10000
    AUTH_USER = "yufeng_hu"
    AUTH_PASSWORD = "qO2MyAEInTHAODP3oJ"

    from opensearchpy import OpenSearch

    try:
        elastic = OpenSearch(hosts=[{
            'host': '192.168.3.146',
            'port': 9200
        }],
                             timeout=TIMEOUT,
                             http_auth=(AUTH_USER, AUTH_PASSWORD))
        response = elastic.search(index=index, body=body)
        return response
    except Exception as esx:
        import pprint

        print("---------------------------------")
        pprint.pprint(esx.info)
        raise esx


logging.basicConfig(format="%(asctime)s - %(name)s - %(levelname)s - %(message)s")
logger = logging.getLogger('chainbase.request')
logger.setLevel(logging.INFO)


def query_chainbase(sql, task_id=None, page=1):
    result = []
    # prepare for request
    url = "https://api.chainbase.online/v1/dw/query"
    headers = {"x-api-key": "2FtLTBTxc9h7CX3YwBeEkrMlnhc"}
    payload = {"query": sql} if not task_id else {"query": sql, "task_id": task_id, "page": page}
    if task_id:
        logger.debug(f"Requesting task: {task_id}-{page})")
    else:
        logger.debug(f"Requesting SQL: '{sql}'")
    # request
    response = requests.post(url, json=payload, headers=headers)
    # deal with response
    if response.status_code == 200:
        try:
            response_text = response.json()
            if response_text["message"] == "ok":
                result = response_text["data"]["result"]
                if "next_page" in response_text["data"]:
                    result.extend(
                        query_chainbase(sql,
                                        task_id=response_text["data"]["task_id"],
                                        page=response_text["data"]["next_page"]))
            else:
                logger.error(f"Query failed. Error message: {response_text['data']['err_msg']}")
        except Exception as e:
            logger.error(f"Failed to decode response: {str(e)}")
    else:
        logger.error('Unexpected status code: {}. Response: {}'.format(response.status_code, response))
    return result




def query_web3(baseUrl, fromBlock, toBlock, topics, addresses=[], extraHeaders={}):
    import copy
    payload = {
        "id": 1,
        "jsonrpc": "2.0",
        "method": "eth_getLogs",
        "params": [{
            "fromBlock": hex(fromBlock),
            "toBlock": hex(toBlock),
            "address": addresses,
            "topics": topics,
        }]
    }
    resp = requests.post(baseUrl, json=payload, headers=extraHeaders, timeout=8)
    resp.raise_for_status()
    resp = resp.json()
    if "error" in resp.keys():
        raise Exception(f"request failed: {resp['error']}")

    entries = []
    for raw in resp["result"]:
        if raw["removed"]:
            continue
        entry = copy.deepcopy(raw)
        entry["blockNumber"] = int(entry["blockNumber"], 16)
        entry["transactionIndex"] = int(entry["transactionIndex"], 16)
        entry["logIndex"] = int(entry["logIndex"], 16)
        entries.append(entry)
    return entries


def query_web3_latest_block(chain, extraHeaders={}):
    url_or_urls = Web3Endpoints[chain]
    url = None
    if isinstance(url_or_urls, str):
        # single url
        url = url_or_urls
    else:
        url = url_or_urls[0]
    payload = {"id": 1, "jsonrpc": "2.0", "method": "eth_blockNumber", "params": []}
    resp = requests.post(url, json=payload, headers=extraHeaders)
    resp.raise_for_status()
    resp = resp.json()
    return int(resp["result"], 16)


_Web3EndpointStatus = defaultdict(bool)


def query_web3_auto(chain, fromBlock, toBlock, topics, addresses=[], extraHeaders={}, retryCount=16):

    def select_url(chain):
        url_or_urls = Web3Endpoints[chain]
        if isinstance(url_or_urls, str):
            # single url
            return url_or_urls
        # urls
        for url in url_or_urls:
            if not _Web3EndpointStatus[url]:
                return url
        # no urls available, retry everything
        for url in url_or_urls:
            _Web3EndpointStatus[url] = False
        return url_or_urls[0]

    def step(retryCounter):
        url = select_url(chain)
        try:
            return query_web3(url, fromBlock, toBlock, topics, addresses, extraHeaders)
        except Exception as e:
            print(f"request count = {retryCounter} failed (url = {url}): {e}")
            _Web3EndpointStatus[url] = True
            

    # the retry loop
    for i in range(1, retryCount + 1):
        r = step(i)
        if r is not None:
            return r


Web3Endpoints = {
    # alchemy
    "ethereum":
        "https://eth-mainnet.g.alchemy.com/v2/OckUvkRX1aiG0hWaCJiK0cV4dcr2xf5t",
    "polygon":
        "https://polygon-mainnet.g.alchemy.com/v2/Qxm0mhaxxs8HkXWdC8CAl24tNBSTSgk5",
    "optimism":
        "https://opt-mainnet.g.alchemy.com/v2/mpeRTW6iPVW4b7eYYxYWZq0FNQOfSXnX",
    "arbitrum":
        "https://arb-mainnet.g.alchemy.com/v2/YQ9AU6YZxrJ4pmhybu0eE2JX1m5A4He9",

    # getblock
    "fantom":
        "https://ftm.getblock.io/0b86cec7-52c8-4a2a-858a-683d28ddccce/mainnet/",
    "avalanche": [
        "https://ava-mainnet.public.blastapi.io/ext/bc/C/rpc",
        "https://avax.getblock.io/mainnet/0b86cec7-52c8-4a2a-858a-683d28ddccce/ext/bc/C/rpc",
        "https://rpc.ankr.com/avalanche",
        "https://api.avax.network/ext/bc/C/rpc",
        "https://1rpc.io/avax/c",
        "https://1rpc.io/avax/c",
    ],
    # others
    "bsc": [
        "https://bsc-mainnet.nodereal.io/v1/64a9df0874fb4a93b9d0a3849de012d3", "https://rpc.ankr.com/bsc",
        "https://rpc-bsc.bnb48.club", "https://bsc-dataseed1.defibit.io", "https://bsc-dataseed2.defibit.io",
        "https://bsc-dataseed3.defibit.io", "https://bsc-dataseed4.defibit.io", "https://bsc-dataseed1.ninicoin.io",
        "https://bsc-dataseed2.ninicoin.io", "https://bsc-dataseed3.ninicoin.io", "https://bsc-dataseed4.ninicoin.io",
        "https://bscrpc.com", "https://bsc-mainnet.public.blastapi.io", "https://binance.nodereal.io",
        "https://bsc.mytokenpocket.vip"
    ],

    "cronos":
        "https://node.croswap.com/rpc",
    "boba":
        "https://mainnet.boba.network",
    "moonriver": [
        "https://moonriver.api.onfinality.io/public",
        "https://rpc.api.moonriver.moonbeam.network",
    ],
    "moonbeam":
        "https://rpc.ankr.com/moonbeam",
    "klaytn":
        "https://klaytn02.fandom.finance",
    "dogechain":
        "https://rpc-us.dogechain.dog",
    "harmony":
        "https://api.s0.t.hmny.io",
    "dfk":
        "https://avax-dfk.gateway.pokt.network/v1/lb/6244818c00b9f0003ad1b619/ext/bc/q2aTwKuyzgs8pynF7UXBZCU7DejbZbZ6EUyHr3JQzYgwNPUPi/rpc",
    "metis":
        "https://andromeda.metis.io/?owner=1088",
    "canto":
        "https://canto.slingshot.finance",
    "astar":
        "https://1rpc.io/astr",
}

Web3QueryBatchSize = {
    # alchemy
    "ethereum": 10000,
    "polygon": 10000,
    "optimism": 10000,
    "arbitrum": 10000,

    # getblock
    "bsc": 1000,
    "fantom": 10000,
    "avalanche": 10000,

    "moonbeam":5000,
    "moonriver":5000,
    "cronos":1000,
    "boba":5000,
    "klaytn":5000,
    "dogechain":5000,
    "harmony":5000,
    "dfk":5000,
    "metis":5000,
    "canto":5000,
}

ChainIds = {
    "ethereum": 1,
    "polygon": 137,
    "optimism": 10,
    "arbitrum": 42161,
    "bsc": 56,
    "fantom": 250,
    "avalanche": 43114,
    "moonbeam": 1284,
    "moonriver": 1285,
    "cronos": 25,
    "boba": 288,
    "klaytn": 8217,
    "dogechain": 2000,
    "harmony": 1666600000,
    "dfk": 53935,
    "metis": 1088,
    "canto": 7700,
}