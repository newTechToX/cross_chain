package match

import (
	"app/dao"
	"app/model"
	"fmt"

	_ "github.com/lib/pq"
)

type Match struct {
}

const (
	MATCHED        = 0
	NO_MATCHED     = -1
	MULTI_MATCHED  = 2
	INFO_NOT_MATCH = 3 //to_addr或direction不匹配
)

// 输入的是一笔cross-in tx
func (a *Match) OneMatch(d *dao.Dao, data model.Result) (int, []model.MatchedId) {
	var res []model.MatchedId
	match_list, err := d.GetOneMatched(data)

	if len(match_list) == 0 || err != nil { //匹配不成功
		fmt.Println(err)
		return NO_MATCHED, res
	}

	//先将所有信息录入
	for _, e := range match_list {
		var ret = model.MatchedId{
			e.Id, data.Id,
		}
		res = append(res, ret)
	}

	//验证筛选出的条目
	if len(match_list) == 1 { // 只有一条数据，该数据的match_id一定为空，只需要匹配信息即可
		e := match_list[0]
		if e.ToAddress == data.ToAddress && e.Direction == model.OutDirection {
			d.UpdateMatchId(e.Id, data.Id)
			d.UpdateMatchId(data.Id, e.Id)
			return MATCHED, res

		} else {
			return INFO_NOT_MATCH, res
		}

	} else { //如果有多条数据
		return MULTI_MATCHED, res
	}
}

func (a *Match) MatchHistory(d *dao.Dao) {
	list_1 := []string{"Across", "Synapse", "CBridge", "Anyswap"}
	for _, n := range list_1 {
		matched_pair, _ := d.GetBatchMatched(n)
		a.BatchMatch(d, matched_pair, 0, len(matched_pair))
	}

	list_2 := []string{"Hop", "Stargate"}
	for _, n := range list_2 {
		matched_pair, _ := d.GetBatchMatched_Hop_Stargate(n)
		a.BatchMatch(d, matched_pair, 0, len(matched_pair))
	}
}

func (a *Match) BatchMatch(d *dao.Dao, matched_pair []*model.MatchedId, start, batch_size int) error {
	var err error
	for i := 0; i < batch_size; i++ {
		err = d.UpdateMatchId((matched_pair)[i+start].SrcID, (matched_pair)[i+start].DstID)
		if err != nil {
			fmt.Println("error to update ", (matched_pair)[i+start].SrcID)
			fmt.Println(err)
			return err
		}

		err = d.UpdateMatchId((matched_pair)[i+start].DstID, (matched_pair)[i+start].SrcID)
		if err != nil {
			fmt.Println("error to update ", (matched_pair)[i+start].DstID)
			fmt.Println(err)
			return err
		}
	}

	fmt.Println("update done: ", start+batch_size)
	return err
}
