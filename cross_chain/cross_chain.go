package crosschain

import (
	"app/cross_chain/across"
	"app/cross_chain/anyswap"
	"app/cross_chain/celer_bridge"
	"app/cross_chain/hop"
	renbridge "app/cross_chain/ren_bridge"
	"app/cross_chain/stargate"
	"app/cross_chain/synapse"
	"app/cross_chain/wormhole"
	"app/model"
	"app/svc"
	"app/utils"
)

//有的项目的corss-outtx 全部转到一个contract 里面，有的则是不同的地址
//只有当全部转到同一个contract的时候，才在下面的OutTxReceiver里面有key

var OutTxReceiver = map[string]map[string]map[string]struct{}{
	"across": getOutTxreceiver(across.AcrossContracts),
}

var TokenTransferDirectly = map[string]struct{}{
	"anyswap":    struct{}{},
	"multichain": struct{}{},
}

func getOutTxreceiver(contracts map[string][]string) map[string]map[string]struct{} {
	ret := make(map[string]map[string]struct{}, len(contracts))
	contracts = utils.LowerStringMap(contracts)
	for k, v := range contracts {
		ret[k] = utils.ConvertSlice2Map(v)
	}
	return ret
}

func GetCollectors(svc *svc.ServiceContext) []model.Colletcor {
	return []model.Colletcor{
		anyswap.NewAnyswapCollector(svc),
		across.NewAcrossCollector(),
		synapse.NewSynapseCollector(svc),
		celer_bridge.NewCBridgeCollector(),
		hop.NewHopCollector(),
		renbridge.NewRenbridgeCollector(),
		stargate.NewStargateCollector(svc),
		wormhole.NewWormHoleCollector(svc),
	}
}
