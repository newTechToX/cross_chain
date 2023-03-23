from common import *
from Crypto.Hash import keccak

tableName = "synapse"
Base = declarative_base()


class SynapseSchema(Base):
    __tablename__ = tableName

    id = Column("id", Integer, primary_key=True)
    chain = Column("chain", String)
    blockNumber = Column("block_number", Integer)
    txIndex = Column("index", Integer)
    txHash = Column("hash", String)
    actionId = Column("action_id", Integer)
    direction = Column("direction", String)
    token = Column("token", String)
    profit = Column("profit", JSON)
    fromAddressError = Column("from_address_error", Integer)
    toAddressProfit = Column("to_address_profit", Integer)
    tokenProfitError = Column("token_profit_error", Integer)
    isFakeToken  = Column("isfaketoken", Integer)
    project = Column("project",String)


def _query_and_parse_from_web3(chain, start_block, end_block):
    topic0s = [[TokenDeposit, TokenDepositAndSwap, TokenMint, TokenMintAndSwap, TokenWithdraw, TokenWithdrawAndRemove,
                TokenRedeem, TokenRedeemAndSwap, TokenRedeemAndRemove]]
    addresses = SynapseContracts[chain]
    items = query_web3_auto(chain, start_block, end_block, topic0s, addresses)
    objs = []

    for item in items:
        topics = item["topics"]
        topic0 = topics[0]
        data = item["data"]
        direction = ""
        _from = ""
        minDy = -1
        toChainID = -1
        fromChainID = -1
        tokenIdxTo = -1
        tokenIdxFrom = -1

        # 根据direction确定token地址
        _to = "0x" + (topics[1])[26:] if len(topics) > 1 else ""
        kappa = topics[2] if len(topics) == 3 else ""

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
        elif topic0 == TokenWithdraw:
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
            kappa = "0x" + keccak_hash.hexdigest()
            fromChainID = ChainIds[chain]
        if direction == "in":
            toChainID = ChainIds[chain]

        objs.append(
            SynapseSchema(chain=chain,
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
                          kappa=kappa,
                          minAmount=minDy,
                          srcTokenIdx=tokenIdxFrom,
                          dstTokenIdx=tokenIdxTo))
    return objs


def synapseMain():
    if not inspect(engine).has_table(tableName):
        print(f"table {tableName} not exists")
        Base.metadata.create_all(engine, checkfirst=True)


if __name__ == "__main__":
    synapseMain()
