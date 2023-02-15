package tests

import (
	"app/aggregator"
	"app/cross_chain/synapse"
	"app/matcher"
	"fmt"
	"testing"
)

func TestAgg(t *testing.T) {
	agg := aggregator.NewAggregator(srvCtx, "avalanche")
	agg.DoJob(synapse.NewSynapseCollector(srvCtx))
}

func TestMatcher(t *testing.T) {
	m := matcher.NewMatcher(srvCtx)
	fmt.Println(m.BeginMatch(15508103, 15518103, "Anyswap", matcher.NewSimpleInMatcher(srvCtx.Dao)))
}
