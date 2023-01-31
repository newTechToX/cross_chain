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
)

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
