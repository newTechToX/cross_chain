package logic

import (
	"app/dao"
	"app/logic/replay"
	"app/model"
	"fmt"
	"testing"
)

func TestLogic_ReplayOutTxLogic(t *testing.T) {
	d := dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")

	a := &Logic{}
	re := &replay.Replayer{}
	a.replayer = re.NewReplayer()
	//hash := "0x4f2eb92a2a9a21bd0c19eab7b4dd3ff4cea4979b70ea4cf56fe20a6e14f73bbd"
	id := 5641147
	//stmt := fmt.Sprintf("select * from anyswap where hash = '%s'", hash)
	stmt := fmt.Sprintf("select * from anyswap where id = %d", id)

	var datas = []*model.Data{}
	err := d.DB().Select(&datas, stmt)
	if err != nil {
		fmt.Println(err)
	}

	err = a.replayOutTxLogic("anyswap", datas)
}
