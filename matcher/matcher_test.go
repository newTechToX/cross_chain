package matcher

import (
	"app/dao"
	"app/model"
	"app/utils"
	"fmt"
	"testing"
)

func TestSimpleMatcher(t *testing.T) {
	_dao := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := NewSimpleInMatcher(_dao)
	var results model.Results
	err := _dao.DB().Select(&results, "select * from common_cross_chain where match_tag = '0x0000000000000000000000000000000000000000000000000000000000000005' and direction = 'in'")
	if err != nil {
		fmt.Println(err)
		return
	}
	shouldUpdates, err := m.Match(results)
	fmt.Println(err)
	utils.PrintPretty(shouldUpdates)

	err = _dao.Update(shouldUpdates)
	fmt.Println(err)
}
