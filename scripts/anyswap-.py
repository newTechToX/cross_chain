from common import *
import time

LogAnySwapOut = "0x97116cf6cd4f6412bb47914d6db18da9e16ab2142f543b86e207c24fbd16b23a"
LogAnySwapIn = "0xaac9ce45fe3adf5143598c4f18a369591a20a3384aedaf1b525d29127e1fcd55"

AnySwapSigs = {LogAnySwapOut: "(uint256,uint256,uint256)", LogAnySwapIn: "(uint256,uint256,uint256)"}

AnySwapContracts = {
    "ethereum": [],
    "bsc": [],
    "polygon": [],
    "fantom": [],
    "arbitrum": [],
    "avalanche": [],
    "optimism": [],
}

Web3QueryStartBlock = {
    "ethereum": 12000000,
    "bsc": 5000000,
    "polygon": 15000000,
    "fantom": 2000000,
    "arbitrum": 900,
    "avalanche": 2400000,
    "optimism": 3400000,
}

Web3QueryBatchSize = {
    # alchemy
    "ethereum": 7000,
    "polygon": 7000,
    "optimism": 7000,
    "arbitrum": 7000,

    # getblock
    "bsc": 900,
    "fantom": 7000,
    "avalanche": 7000,
}

tableName = "anyswap"
Base = declarative_base()


class AnyswapSchema(Base):
    __tablename__ = tableName

    id = Column("id", Integer, primary_key=True)
    chain = Column("chain", String)
    blockNumber = Column("block_number", Integer)
    txIndex = Column("tx_index", Integer)
    txHash = Column("hash", String)
    logIndex = Column("log_index", Integer)
    direction = Column("direction", String)
    contract = Column("contract", String)
    fromChain = Column("from_chain", Numeric(256))
    token = Column("token", String)
    fromAddress = Column("from_address", String)
    amount = Column("amount", Numeric)
    toAddress = Column("to_address", String)
    toChain = Column("to_chain", Numeric(256))
    srcTxHash = Column("match_tag", String)
    project = Column("project", String)


def _query_and_parse_from_web3(chain, start_block, end_block):
    topics = [[LogAnySwapIn, LogAnySwapOut]]
    addresses = AnySwapContracts[chain]
    items = query_web3_auto(chain, start_block, end_block, topics, addresses)
    objs = []

    for item in items:
        topics = item["topics"]
        topic0 = topics[0]
        data = item["data"]
        direction = ""
        srcTxHash = ""
        _from = ""

        # 根据direction确定token地址
        if topic0 == LogAnySwapIn:
            srcTxHash = topics[1]
            token = "0x" + (topics[2])[26:]
            _to = "0x" + (topics[3])[26:]
            (amount, fromChainID, toChainID) = decode_single(AnySwapSigs[topic0], bytes.fromhex(data[2:]))
            direction = "in"
        elif topic0 == LogAnySwapOut:
            token = "0x" + (topics[1])[26:]
            _from = "0x" + (topics[2])[26:]
            _to = "0x" + (topics[3])[26:]
            (amount, fromChainID, toChainID) = decode_single(AnySwapSigs[topic0], bytes.fromhex(data[2:]))
            direction = "out"
            srcTxHash = item["transactionHash"]
        else:
            continue

        objs.append(
            AnyswapSchema(chain=chain,
                          blockNumber=item["blockNumber"],
                          txIndex=item["transactionIndex"],
                          txHash=item["transactionHash"],
                          logIndex=item["logIndex"],
                          direction=direction,
                          contract=item["address"],
                          fromChain=fromChainID,
                          token=token,
                          fromAddress=_from,
                          amount=amount,
                          toAddress=_to,
                          toChain=toChainID,
                          srcTxHash=srcTxHash,
                          project="anyswap"))
    return objs


def anyswapMain():
    if not inspect(engine).has_table(tableName):
        print(f"table {tableName} not exists")
        Base.metadata.create_all(engine, checkfirst=True)

    chains = ["bsc", "ethereum", "polygon", "avalanche", "arbitrum", "fantom", "optimism"]
    for chain in chains:
        session = Session()
        stmt = select(AnyswapSchema.blockNumber).where(AnyswapSchema.chain == chain).order_by(
            AnyswapSchema.blockNumber.desc()).limit(1)
        latest_block_in_db = session.scalar(stmt)
        print(f"latest block in database for {chain}: {latest_block_in_db}")

        latest_block_from_web3 = query_web3_latest_block(chain)
        print(f"latest block from web3 for {chain}: {latest_block_from_web3}")

        start_block = latest_block_in_db if latest_block_in_db is not None else Web3QueryStartBlock[chain]
        end_block = latest_block_from_web3
        step = Web3QueryBatchSize[chain]
        inserted = 0
        pbar = tqdm(range(start_block, end_block, step))
        for i in pbar:
            objs = _query_and_parse_from_web3(chain, i + 1, i + step)
            session.bulk_save_objects(objs)
            session.commit()
            inserted += len(objs)
            pbar.set_postfix({"last_inserted": len(objs)})
        print(f"objects inserted: {inserted}")


if __name__ == "__main__":
    while True:
        try:
            anyswapMain()
            time.sleep(1)
        except requests.exceptions.RequestException as e:   
            anyswapMain()
            time.sleep(1)
