package tests

import (
	"app/aggregator"
	"app/cross_chain/synapse"
	"testing"
)

func TestAgg(t *testing.T) {
	agg := aggregator.NewAggregator(srvCtx, "avalanche")
	agg.DoJob(synapse.NewSynapseCollector(srvCtx))
}
