package match

import (
	"app/dao"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

// since multimatched data are too many in Across
// TestGetMultiMatched_Across gets those multi-matched ones
func TestGetMultiMatched_Across(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	file, err := os.OpenFile("./t.json", os.O_RDWR, 0666)
	if err != nil {
		println("failed to open file")
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	n := 0
	for {
		str, err := reader.ReadString('\n')
		n++
		if err == io.EOF {
			break
		}
		if n%4 != 3 {
			continue
		}
		oriId, _ := strconv.ParseInt(strings.Fields(strings.TrimSpace(str))[1], 10, 64)
		println(oriId)
		data, err := d.GetOne(0, "", oriId)
		fmt.Println(data.Id, data.Chain, data.Hash, data.MatchId, data.MatchTag)
	}
}

// Followings are functions to complete basic matches
