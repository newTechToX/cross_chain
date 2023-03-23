from common import *
from Crypto.Hash import keccak

TokenDeposit = "0xda5273705dbef4bf1b902a131c2eac086b7e1476a8ab0cb4da08af1fe1bd8e3b"
TokenDepositAndSwap = "0x79c15604b92ef54d3f61f0c40caab8857927ca3d5092367163b4562c1699eb5f"

TokenMint = "0xbf14b9fde87f6e1c29a7e0787ad1d0d64b4648d8ae63da21524d9fd0f283dd38"
TokenMintAndSwap = "0x4f56ec39e98539920503fd54ee56ae0cbebe9eb15aa778f18de67701eeae7c65"
TokenWithdraw = "0x8b0afdc777af6946e53045a4a75212769075d30455a212ac51c9b16f9c5c9b26"
TokenWithdrawAndRemove = "0xc1a608d0f8122d014d03cc915a91d98cef4ebaf31ea3552320430cba05211b6d"
TokenRedeem = "0xdc5bad4651c5fbe9977a696aadc65996c468cde1448dd468ec0d83bf61c4b57c"
TokenRedeemAndSwap = "0x91f25e9be0134ec851830e0e76dc71e06f9dade75a9b84e9524071dbbc319425"
TokenRedeemAndRemove = "0x9a7024cde1920aa50cdde09ca396229e8c4d530d5cfdc6233590def70a94408c"

chains = ["arbitrum",]

sigs = {
    TokenDeposit: "(uint256,address,uint256)",
    TokenDepositAndSwap: "(uint256,address,uint256,uint8,uint8,uint256,uint256)",
    TokenRedeem: "(uint256,address,uint256)",
    TokenRedeemAndSwap: "(uint256,address,uint256,uint8,uint8,uint256,uint256)",
    TokenRedeemAndRemove: "(uint256,address,uint256,uint8,uint256,uint256)",

    TokenMint: "(address,uint256,uint256)",
    TokenMintAndSwap: "(address,uint256,uint256,uint8,uint8,uint256,uint256,bool)",
    TokenWithdraw: "(address,uint256,uint256)",
    TokenWithdrawAndRemove: "(address,uint256,uint256,uint8,uint256,uint256,bool)",
}

SynapseContracts = {
    "ethereum": [
#        '0x2796317b0fF8538F253012862c06787Adfb8cEb6', '0x6571d6be3d8460CF5F7d6711Cd9961860029D85F'
    ],
    "bsc": [
 #       '0xd123f70ae324d34a9e76b67a27bf77593ba8749f', '0x749f37df06a99d6a8e065dd065f8cf947ca23697'
    ],
    "polygon": [
  #      '0x1c6ae197ff4bf7ba96c66c5fd64cb22450af9cc8', '0x8f5bbb2bb8c2ee94639e55d5f41de9b4839c1280'
    ],
    "fantom": [
   #     '0xb003e75f7e0b5365e814302192e99b4ee08c0ded', '0xaf41a65f786339e7911f4acdad6bd49426f2dc6b'
    ],
    "arbitrum": [
    #    '0x37f9ae2e0ea6742b9cad5abcfb6bbc3475b3862b', '0x6f4e8eba4d337f874ab57478acc2cb5bacdc19c9'
    ],
    "avalanche": [
     #   '0x0ef812f4c68dc84c22a4821ef30ba2ffab9c2f3a', '0xc05e61d0e7a63d27546389b7ad62fdff5a91aace'
    ],
    "optimism": [
      #  '0x470f9522ff620ee45df86c58e54e6a645fe3b4a7', '0xaf41a65f786339e7911f4acdad6bd49426f2dc6b'
        ],
    "moonbeam":[],
    "cronos":[],
    "boba":[],
    "klaytn":[],
    "dogechain":[],
    "harmony":[],
    "dfk":[],
    "metis":[],
    "canto":[],
    "moonriver":[],
}

Web3QueryStartBlock = {
    "ethereum": 13033669,
    "bsc": 10065475,
    "polygon": 18026806,
    "fantom": 18503502,
    "arbitrum": 657404,
    "avalanche": 3376709,
    "optimism": 30718,
}



tableName = "synapse"
Base = declarative_base()


class SynapseSchema(Base):
    __tablename__ = tableName

    id = Column("id", Integer, primary_key=True)
    chain = Column("chain", String)
    blockNumber = Column("block_number", Integer)
    txIndex = Column("tx_index", Integer)
    txHash = Column("hash", String)
    actionId = Column("log_index", Integer)
    direction = Column("direction", String)
    contract = Column("contract", String)
    fromChain = Column("from_chain", String)
    token = Column("token", String)
    fromAddress = Column("from_address", String)
    amount = Column("amount", Numeric)
    totalAmount = Column("total_amount", Numeric)
    fillAmount = Column("fill_amount", Numeric)
    toAddress = Column("to_address", String)
    toChain = Column("to_chain", String)
    kappa = Column("match_tag", String)
    minAmount = Column("min_amount", Numeric)
    srcTokenIdx = Column("src_token_idx", Integer)
    dstTokenIdx = Column("dst_token_idx", Integer)


def _query_and_parse_from_web3(chain, start_block, end_block):
    topics = [[TokenDeposit, TokenDepositAndSwap, TokenMint, TokenMintAndSwap, TokenWithdraw, TokenWithdrawAndRemove,
    TokenRedeem, TokenRedeemAndSwap, TokenRedeemAndRemove]]
    addresses = SynapseContracts[chain]
    items = query_web3_auto(chain, start_block, end_block, topics, addresses)
    objs = []

    for item in items:
        topics = item["topics"]
        topic0 = topics[0]
        data = item["data"]
        direction = ""
        _from = ""
        minDy = -1
        toChainID = -1
        tokenIdxTo = -1
        tokenIdxFrom = -1
       

        # 根据direction确定token地址
        _to = "0x" + (topics[1])[26:]
        kappa = topics[2] if len(topics)==3 else ""

        if topic0 == TokenDeposit:
            (toChainID, token, amount) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "out"
        elif topic0 == TokenDepositAndSwap:
            (toChainID, token, amount, tokenIdxFrom,
             tokenIdxTo, minDy, ddl) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "out"
        elif topic0 == TokenRedeem:
            (toChainID, token, amount) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "out"
        elif topic0 == TokenRedeemAndSwap:
            (toChainID, token, amount, tokenIdxFrom,
             tokenIdxTo, minDy, ddl) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "out"
        elif topic0 == TokenRedeemAndRemove:
            (toChainID, token, amount, swapTokenIdx, minDy,
             ddl) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "out"

        elif topic0 == TokenMint:
            (token, amount, fee) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "in"
        elif topic0 == TokenMintAndSwap:
            (token, amount, fee, tokenIdxFrom, tokenIdxTo,
             minDy, ddl, swapSuccess) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "in"
        elif topic0 ==TokenWithdraw:
            (token, amount, fee) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "in"
        elif topic0 == TokenWithdrawAndRemove:
            (token, amount, fee, swapTokenIdx, minDy,
             ddl, swapSuccess) = decode_single(sigs[topic0], bytes.fromhex(data[2:]))
            direction = "in"
        else:
            continue

        if direction == "out":
            keccak_hash = keccak.new(digest_bits=256)
            keccak_hash.update(bytes(item["transactionHash"], encoding="utf8"))
            kappa = "0x"+keccak_hash.hexdigest()

        objs.append(
            SynapseSchema(chain=chain,
                          blockNumber=item["blockNumber"],
                          txIndex=item["transactionIndex"],
                          txHash=item["transactionHash"],
                          actionId=item["logIndex"],
                          direction=direction,
                          contract=item["address"],
                          fromChain=_from,
                          token=token,
                          fromAddress=_from,
                          amount=amount,
                          toAddress=_to,
                          toChain=toChainID,
                          kappa=kappa,
                          minAmount=minDy,
                          srcTokenIdx=tokenIdxFrom,
                          dstTokenIdx=tokenIdxTo))
    return objs


def synapseMain(Web3QueryStartBlock):
    if not inspect(engine).has_table(tableName):
        print(f"table {tableName} not exists")
        Base.metadata.create_all(engine, checkfirst=True)

    for chain in chains:
        session = Session()
        latest_block_from_web3 = query_web3_latest_block(chain)

        start_block = Web3QueryStartBlock[chain]
        print(f"latest block in database for {chain}: {start_block}")
        end_block = latest_block_from_web3
        print(f"latest block from web3 for {chain}: {latest_block_from_web3}")

        step = Web3QueryBatchSize[chain]
        inserted = 0
        pbar = tqdm(range(start_block, end_block, step))
        for i in pbar:
            j = min(i+step, end_block)
            objs = _query_and_parse_from_web3(chain, i + 1, j)
            session.bulk_save_objects(objs)
            session.commit()
            inserted += len(objs)
            pbar.set_postfix({"last_inserted": len(objs)})

        dt = (datetime.datetime.now().strftime('%Y-%m-%d  %H:%M:%S'))
        print(f"{dt}\tobjects inserted: {inserted}")
        Web3QueryStartBlock[chain] = objs[-1].blockNumber if len(objs)>0 else end_block



if __name__ == "__main__":
    # Web3QueryStartBlock = {}
     for chain in chains:
        session = Session()
        #stmt = select(max(AnyswapSchema.blockNumber)).where(AnyswapSchema.chain == chain)
        stmt = text("select max(block_number) from synapse where chain = '%s'" %(chain))
        last_block_in_db = session.scalar(stmt)

        if last_block_in_db is not None:
            Web3QueryStartBlock[chain] = session.scalar(stmt)
        elif chain not in Web3QueryStartBlock:
            Web3QueryStartBlock[chain] = 1

     while True:
        try:
            synapseMain(Web3QueryStartBlock)
            time.sleep(46)
        except requests.exceptions.SSLError as e:
            print("SSLError: ", e)
            time.sleep(600)
        except requests.exceptions.ConnectionError as e:
            print("ConnectionError: ", e)
            time.sleep(600)
        except requests.exceptions.RequestException as e:
            print("RequestException: ", e)
            time.sleep(600)
        except:
            time.sleep(600)
