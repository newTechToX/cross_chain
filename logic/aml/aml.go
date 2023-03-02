package aml

import (
	"app/config"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type AML struct {
	url []string
	key string
}

func NewAML(config_path string) *AML {
	if config_path == "" {
		config_path = "../txt_config.yaml"
	}
	var cfg Config
	config.LoadCfg(&cfg, config_path)
	urls := []string{
		"https://aml.blocksec.com/api/aml/v2/address",
		"https://aml.blocksec.com/api/aml/v2/addresses",
	}
	return &AML{
		url: urls,
		key: cfg.AmlApiKey,
	}
	return nil
}

func (a *AML) AmlInfoContainWords(info []*AddressInfo, words []string) map[string]struct{} {
	if len(info) == 0 || len(words) == 0 {
		return nil
	}
	var res = make(map[string]struct{})
	for _, word := range words {
		for _, i := range info {
			description := i.Name
			for _, l := range i.Labels {
				description = description + "," + l
			}
			if strings.Contains(description, word) {
				res[word] = struct{}{}
				break
			}
		}
	}
	return res
}

//如果从aml没有查到信息，那么在map里没有信息

func (a *AML) QueryAml(chain string, address []string) (map[string][]*AddressInfo, error) {
	var raw_msg *RawMsg
	time.Sleep(500 * time.Millisecond)

	if chain == "ethereum" {
		chain = "eth"
	}

	buffer, err := a.queryAml(chain, address)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buffer, &raw_msg)

	if raw_msg == nil {
		err = fmt.Errorf("QueryAml(): null return data from aml, address: %s, error: %s", address[0], err)
		return nil, err
	}
	if raw_msg.Code != 200 || raw_msg.Msg != "" {
		err = fmt.Errorf("code: %d, address: %s, error: %s", raw_msg.Code, address[0], err)
		return nil, err
	}

	res := a.getInfo(raw_msg.Data)
	/*for _, addr := range address {
		if _, ok := res[addr]; !ok {
			res[addr] = nil
		}
	}*/
	return res, err
}

func (a *AML) queryAml(chain string, address []string) ([]byte, error) {
	var query_str string
	var err error
	if len(address) == 0 {
		err = fmt.Errorf("no adresses")
		return nil, err
	}

	if len(address) == 1 {
		query_str = fmt.Sprintf("%s/%s/%s", a.url[0], chain, address[0])
	} else {
		query_str = fmt.Sprintf("%s/%s/", a.url[1], chain)
		for _, e := range address {
			query_str = query_str + e + ","
		}
	}

	req, err := http.NewRequest("GET", query_str, strings.NewReader(""))
	req.Header.Set("API-KEY", a.key)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Error("failed request ", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (a *AML) getInfo(data []*ReturnData) map[string][]*AddressInfo {
	var info_map = make(map[string][]*AddressInfo)
	for _, d := range data {
		if !d.IsValid {
			continue
		}
		if d.Labels != nil && !reflect.DeepEqual(*d.Labels, Labels{}) {
			var labels []string
			for _, l := range d.Labels.Others {
				if l.Confidence >= 8 {
					labels = append(labels, l.Label)
				}
			}
			info_map[d.Address] = append(info_map[d.Address],
				&AddressInfo{
					Chain:   d.Chain,
					Address: d.Address,
					Name:    (d.Labels).name(),
					Risk:    d.Risk,
					Labels:  labels,
				})
		} else if len(d.CompatibleChainLabels) > 0 {
			for _, comp := range d.CompatibleChainLabels {
				info_map[d.Address] = append(info_map[d.Address],
					&AddressInfo{
						Chain:   comp.Chain,
						Address: d.Address,
						Name:    (comp.Labels()).name(),
						Risk:    d.Risk,
					})
			}
		}
	}

	return info_map
}
