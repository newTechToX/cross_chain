from common import *

LogAnySwapOut = "0x97116cf6cd4f6412bb47914d6db18da9e16ab2142f543b86e207c24fbd16b23a"
LogAnySwapIn = "0xaac9ce45fe3adf5143598c4f18a369591a20a3384aedaf1b525d29127e1fcd55"

AnySwapSigs = {LogAnySwapOut: "(uint256,uint256,uint256)", LogAnySwapIn: "(uint256,uint256,uint256)"}

AnySwapContracts = {
    "ethereum": [
        '0x6b7a87899490ece95443e979ca9485cbe7e71522', '0x7782046601e7b9b05ca55a3899780ce6ee6b8b2b',
        '0xf0457c4c99732b716e40d456acb3fc83c699b8ba', '0xe95fd76cf16008c12ff3b3a937cb16cd9cc20284',
        '0x765277eebeca2e31912c9946eae1021199b39c61', '0xba8da9dcf11b50b03fd5284f164ef5cdef910705',
    ],
    "bsc": [
        '0xabd380327fe66724ffda91a87c772fb8d00be488', '0x92c079d3155c2722dbf7e65017a5baf9cd15561c',
        '0x56a6c850cebe23f0c7891a004bef57265cda4d13', '0xd1c5966f9f5ee6881ff6b261bbeda45972b1b5f3',
        '0xe1d592c3322f1f714ca11f05b6bc0efef1907859', '0xd1a891e6eccb7471ebd6bc352f57150d4365db21',
        '0x58892974758a4013377a45fad698d2ff1f08d98e', '0xf9736ec3926703e85c843fc972bd89a7f8e827c0'
    ],
    "polygon": [
        '0x2ef4a574b72e1f555185afa8a09c6d1a8ac4025c', '0x0b23341fa1da0171f52aa8ef85f3946b44d35ac0',
        '0x6ff0609046a38d76bd40c5863b4d1a2dce687f73', '0xd50380e953603b37a74dc67c92fc5e19e0b65469',
        '0x4f3aff3a747fcade12598081e80c6605a8be192f', '0x72c290f3f13664b024ee611983aa2d5621ebe917',
        '0x1ccca1ce62c62f7be95d4a67722a8fdbed6eecb4', '0x84cebca6bd17fe11f7864f7003a1a30f2852b1dc',
        '0xafaace7138ab3c2bcb2db4264f8312e1bbb80653'
    ],
    "fantom": [
        '0xb576c9403f39829565bd6051695e2ac7ecf850e2', '0x85fd5f8dbd0c9ef1806e6c7d4b787d438621c1dc',
        '0x0b23341fa1da0171f52aa8ef85f3946b44d35ac0', '0xf98f70c265093a3b3adbef84ddc29eace900685b',
        '0x1ccca1ce62c62f7be95d4a67722a8fdbed6eecb4', '0x24e2a6f08e3cc2baba93bd9b89e19167a37d6694',
        '0xf3ce95ec61114a4b1bfc615c16e6726015913ccc'
    ],
    "arbitrum": [
        '0x0cae51e1032e8461f4806e26332c030e34de3adb', '0xcb9f441ffae898e7a2f32143fd79ac899517a9dc',
        '0x39fde572a18448f8139b7788099f0a0740f51205', '0xc931f61b1534eb21d8c11b24f3f5ab2471d4ab50',
        '0x2bf9b864cdc97b08b6d79ad4663e71b8ab65c45c', '0x650af55d5877f289837c30b94af91538a7504b76',
        '0xa71353bb71dda105d383b02fc2dd172c4d39ef8b'
    ],
    "avalanche": [
        '0x9b17baadf0f21f03e35249e0e59723f34994f806', '0xe5cf1558a1470cb5c166c2e8651ed0f3c5fb8f42',
        '0x05f024c6f5a94990d32191d6f36211e3ee33504e', '0xb0731d50c681c45856bfc3f7539d5f61d4be81d8',
        '0x34324e1598bf02ccd3dea93f4e332b5507097473', '0x833f307ac507d47309fd8cdd1f835bef8d702a93'
    ],
    "optimism": ['0xdc42728b0ea910349ed3c6e1c9dc06b5fb591f98', '0x80a16016cc4a2e6a2caca8a4a498b1699ff0f844'],
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
    "ethereum": 10000,
    "polygon": 10000,
    "optimism": 10000,
    "arbitrum": 10000,

    # getblock
    "bsc": 1000,
    "fantom": 10000,
    "avalanche": 1000,
}

tableName = "anyswap"
Base = declarative_base()


class AnyswapSchema(Base):
    __tablename__ = tableName

    id = Column("id", Integer, primary_key=True)
    chain = Column("chain", String)
    blockNumber = Column("block_number", Integer)
    txIndex = Column("tx_index", Integer)
    txHash = Column("tx_hash", String)
    logIndex = Column("log_index", Integer)
    direction = Column("direction", String)
    contract = Column("contract", String)
    fromChain = Column("from_chain", String)
    token = Column("token", String)
    fromAddress = Column("from_address", String)
    amount = Column("amount", Numeric)
    toAddress = Column("to_address", String)
    toChain = Column("to_chain", String)
    srcTxHash = Column("src_tx_hash", String)


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
                          srcTxHash=srcTxHash))
    return objs


def anyswapMain():
    if not inspect(engine).has_table(tableName):
        print(f"table {tableName} not exists")
        Base.metadata.create_all(engine, checkfirst=True)

    chains = ["bsc", "avalanche", "bsc","ethereum", "polygon", "optimism", "arbitrum", "fantom"]
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
    anyswapMain()
