from common import *
from Crypto.Hash import keccak

#Burn_2 (bytes32 burnId, address token, address account, uint256 amount, uint64 toChainId, address toAccount, uint64 nonce)

Send = "0x89d8051e597ab4178a863a5190407b98abfeff406aa8db90c59af76612e58f01"
Burn = "0x6298d7b58f235730b3b399dc5c282f15dae8b022e5fbbf89cee21fd83c8810a3"
Deposited_1 = "0x15d2eeefbe4963b5b2178f239ddcc730dda55f1c23c22efb79ded0eb854ac789"
Deposited_2 = "0x28d226819e371600e26624ebc4a9a3947117ee2760209f816c789d3a99bf481b"

Mint = "0x5bc84ecccfced5bb04bfc7f3efcdbe7f5cd21949ef146811b4d1967fe41f777a"
Relay = "0x79fa08de5149d912dce8e5e8da7a7c17ccdf23dd5d3bfe196802e6eb86347c7c"
Withdrawn = "0x296a629c5265cb4e5319803d016902eb70a9079b89655fe2b7737821ed88beeb"

sigs = {
    Relay: "(bytes32,address,address,address,uint256,uint64,bytes32)",
    Send: "(bytes32,address,address,address,uint256,uint64,uint64,uint32)",
	Burn: "(bytes32,address,address,uint256,uint64,address,uint64)",
  	Deposited_1: "(bytes32,address,address,uint256,uint64,address)",
	Deposited_2: "(bytes32,address,address,uint256,uint64,address,uint64)",
	Mint: "(bytes32,address,address,uint256,uint64,bytes32,address)",
	Withdrawn: "(bytes32,address,address,uint256,uint64,bytes32,address)",
}

CbridgeContracts = {
    "ethereum": [
        "0x5427FEFA711Eff984124bFBB1AB6fbf5E3DA1820", "0xB37D31b2A74029B5951a2778F959282E2D518595",
		"0x7510792A3B1969F9307F3845CE88e39578f2bAE1", "0x52E4f244f380f8fA51816c8a10A63105dd4De084",
		"0x16365b45EB269B5B5dACB34B4a15399Ec79b95eB",
    ],
    "bsc": [
		"0xdd90E5E87A2081Dcf0391920868eBc2FFB81a1aF", "0x78bc5Ee9F11d133A08b331C2e18fE81BE0Ed02DC",
		"0x11a0c9270D88C99e221360BCA50c2f6Fda44A980", "0x26c76F7FeF00e02a5DD4B5Cc8a0f717eB61e1E4b",
		"0xd443FE6bf23A4C9B78312391A30ff881a097580E",
    ],
    "polygon": [
		"0x88DCDC47D2f83a99CF0000FDF667A468bB958a78", "0xc1a2D967DfAa6A10f3461bc21864C23C1DD51EeA",
		"0x4C882ec256823eE773B25b414d36F92ef58a7c0C", "0xb51541df05DE07be38dcfc4a80c05389A54502BB",
		"0x4d58FDC7d0Ee9b674F49a0ADE11F26C3c9426F7A",
    ],
    "fantom": [
		"0x374B8a9f3eC5eB2D97ECA84Ea27aCa45aa1C57EF", "0x7D91603E79EA89149BAf73C9038c51669D8F03E9",
		"0x30F7Aa65d04d289cE319e88193A33A8eB1857fb9", "0x38D1e20B0039bFBEEf4096be00175227F8939E51",
    ],
    "arbitrum": [
		"0xb3833Ecd19D4Ff964fA7bc3f8aC070ad5e360E56", "0x1619DE6B6B20eD217a58d00f37B9d47C7663feca",
		"0xFe31bFc4f7C9b69246a6dc0087D91a91Cb040f76", "0xEA4B1b0aa3C110c55f650d28159Ce4AD43a4a58b",
		"0xbdd2739AE69A054895Be33A22b2D2ed71a1DE778",
    ],
    "avalanche": [
		"0xef3c714c9425a8F3697A9C969Dc1af30ba82e5d4", "0x5427FEFA711Eff984124bFBB1AB6fbf5E3DA1820",
		"0xb51541df05DE07be38dcfc4a80c05389A54502BB", "0xb774C6f82d1d5dBD36894762330809e512feD195",
		"0x88DCDC47D2f83a99CF0000FDF667A468bB958a78",
    ],
    "optimism": [
		"0x9D39Fc627A6d9d9F8C831c16995b209548cc3401", "0xbCfeF6Bb4597e724D720735d32A9249E0640aA11",
		"0x61f85fF2a2f4289Be4bb9B72Fc7010B3142B5f41",
    ],
    "moonbeam":[],
    "moonriver":[],
    "cronos":[],
    "boba":[],
    "klaytn":[],
    "dogechain":[],
    "harmony":[],
    "dfk":[],
    "metis":[],
    "canto":[],
}

Web3QueryStartBlock = {
    "ethereum": 13719989,
    "bsc": 13099216,
    "polygon": 22006535,
    "fantom": 23658021,
    "arbitrum": 3483123,
    "avalanche": 7660096,
    "optimism": 646444,
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
    "cronos":5000,
    "boba":5000,
    "klaytn":5000,
    "dogechain":5000,
    "harmony":5000,
    "dfk":5000,
    "metis":5000,
    "canto":5000,
}

tableName = "cbridge_v2"
Base = declarative_base()


class CbridgeSchema(Base):
    __tablename__ = tableName

    id = Column("id", Integer, primary_key=True)
    chain = Column("chain", String)
    blockNumber = Column("block_number", Integer)
    txIndex = Column("tx_index", Integer)
    txHash = Column("hash", String)
    actionId = Column("action_id", Integer)
    direction = Column("direction", String)
    contract = Column("contract", String)
    fromChain = Column("from_chain", String)
    token = Column("token", String)
    fromAddress = Column("from_address", String)
    amount = Column("amount", Numeric)
    toAddress = Column("to_address", String)
    toChain = Column("to_chain", String)
    matchTag = Column("match_tag", String)
    detail = Column("detail", String)
    slippage = Column("slippage", Integer)



def _query_and_parse_from_web3(chain, start_block, end_block):
    topics = [[Burn, Send, Deposited_1, Deposited_2, Mint, Relay, Withdrawn]]
    addresses = CbridgeContracts[chain]
    items = query_web3_auto(chain, start_block, end_block, topics, addresses)
    objs = []

    for item in items:
        topics = item["topics"]
        topic0 = topics[0]
        data = item["data"]
        direction = ""
        _from = ""
        detail = ""
        toChainID = -1
        fromChainID = -1
        matchTag = ""
        slippage = -1
       

        # 根据direction确定token地址
        if topic0 == Burn:
            (matchTag, token, _from, amount, toChainID, _to, detail) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "out"
        elif topic0 == Deposited_1:
            (matchTag, _from, token, amount, toChainID, _to) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "out"
        elif topic0 == Deposited_2:
            (matchTag, _from, token, amount, toChainID, _to, detail) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "out"
        elif topic0 == Send:
            (matchTag, _from, _to, token, amount, toChainID, detail, slippage) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "out"

        elif topic0 == Mint:
            (detail, token, _to, amount, fromChainID, matchTag, _from) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "in"
        elif topic0 == Relay:
            (detail, _from, _to, token, amount, fromChainID, matchTag) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "in"
        elif topic0 ==Withdrawn:
            (detail, _to, token, amount, fromChainID, matchTag, _from) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "in"
        else:
            continue

        if direction == "in":
            toChainID = ChainIds[chain]
        else:
            fromChainID = ChainIds[chain]

        objs.append(
            CbridgeSchema(chain=chain,
                          blockNumber=item["blockNumber"],
                          txIndex=item["transactionIndex"],
                          txHash=item["transactionHash"],
                          actionId=item["logIndex"],
                          direction=direction,
                          contract=item["address"],
                          fromChain=fromChainID,
                          token=token,
                          fromAddress=_from,
                          amount=amount,
                          toAddress=_to,
                          toChain=toChainID,
                          matchTag=matchTag,
                          detail=detail,
                          slippage=slippage                          
                          ))
    return objs


def cbridgeMain():
    if not inspect(engine).has_table(tableName):
        print(f"table {tableName} not exists")
        Base.metadata.create_all(engine, checkfirst=True)

    chains = ["ethereum", "polygon", "optimism", "arbitrum", "fantom", "avalanche", "bsc", "cronos", 
    "boba", "moonriver", "moonbeam", "harmony", "metis", "astar"]
    for chain in chains:
        session = Session()
        stmt = select(CbridgeSchema.blockNumber).where(CbridgeSchema.chain == chain).order_by(
            CbridgeSchema.blockNumber.desc()).limit(1)
        latest_block_in_db = session.scalar(stmt)
        print(f"latest block in database for {chain}: {latest_block_in_db}")

        latest_block_from_web3 = query_web3_latest_block(chain)
        print(f"latest block from web3 for {chain}: {latest_block_from_web3}")

        if latest_block_in_db is not None:
            start_block = latest_block_in_db
        elif chain in Web3QueryStartBlock:
            start_block = Web3QueryStartBlock[chain]
        else:
            start_block = 1

        start_block = latest_block_in_db if latest_block_in_db is not None else Web3QueryStartBlock[chain]
        end_block = latest_block_from_web3
        step = Web3QueryBatchSize[chain]
        inserted = 0
        pbar = tqdm(range(start_block, end_block, step))
        for i in pbar:
            end = min(i + step, latest_block_from_web3)
            objs = _query_and_parse_from_web3(chain, i + 1, end)
            session.bulk_save_objects(objs)
            session.commit()
            inserted += len(objs)
            pbar.set_postfix({"last_inserted": len(objs)})
        print(f"objects inserted: {inserted}")


if __name__ == "__main__":
    cbridgeMain()
