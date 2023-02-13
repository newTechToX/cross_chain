package aml

import (
	"app/utils"
	"sort"
)

type CompatibleChainLabels struct {
	Chain        string          `json:"chain"`
	ContractInfo *ContractInfo   `json:"contract_info"`
	NameTag      *string         `json:"name_tag"`
	EntityInfo   []*EntityInfo   `json:"entity_info"`
	PropertyInfo []*PropertyInfo `json:"property_info"`
	Others       []*Others       `json:"others"`
}

type Labels struct {
	ContractInfo *ContractInfo   `json:"contract_info"`
	NameTag      *string         `json:"name_tag"`
	EntityInfo   []*EntityInfo   `json:"entity_info"`
	PropertyInfo []*PropertyInfo `json:"property_info"`
	Others       []*Others       `json:"others"`
}

type ContractInfo struct {
	ContractName string `json:"contract_name"`
	TokenName    string `json:"token_name"`
}

type EntityInfo struct {
	EntityType     string `json:"entity_type"`
	Entity         string `json:"entity"`
	Category       string `json:"category"`
	EntityProperty string `json:"entity_property"`
	Confidence     int    `json:"confidence"`
}

type PropertyInfo struct {
	AddressProperty string `json:"address_property"`
	Category        string `json:"category"`
	Confidence      int    `json:"confidence"`
}

type Others struct {
	Label      string `json:"label"`
	Confidence int    `json:"confidence"`
}

type NameAndConfidence struct {
	Name       string
	Confidence int
}

type NamesAndConfidences []*NameAndConfidence

func (l CompatibleChainLabels) Labels() *Labels {
	return &Labels{
		EntityInfo:   l.EntityInfo,
		PropertyInfo: l.PropertyInfo,
		ContractInfo: l.ContractInfo,
		NameTag:      l.NameTag,
		Others:       l.Others,
	}
}

func (l Labels) namesAndConfidences() NamesAndConfidences {
	var res NamesAndConfidences
	for _, e := range l.EntityInfo {
		res = append(res,
			&NameAndConfidence{e.Entity,
				e.Confidence},
		)
	}
	for _, e := range l.PropertyInfo {
		res = append(res,
			&NameAndConfidence{e.AddressProperty,
				e.Confidence},
		)
	}
	for _, e := range l.Others {
		res = append(res,
			&NameAndConfidence{e.Label,
				e.Confidence},
		)
	}
	return res
}

func (n NamesAndConfidences) Len() int {
	return len(n)
}

func (n NamesAndConfidences) Less(i, j int) bool {
	return n[i].Confidence < n[j].Confidence
}

func (n NamesAndConfidences) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (l Labels) name() (name string) {
	if l.NameTag != nil {
		name = *l.NameTag
	} else if l.ContractInfo != nil {
		if name = l.ContractInfo.TokenName; name != "" {
		} else {
			name = l.ContractInfo.TokenName
		}
	} else {
		names := l.namesAndConfidences()

		sort.Sort(names)
		name = names[len(names)-1].Name
	}
	return
}

// 没什么用
func (l Labels) minConfidence() int {
	var confidences = []int{}
	for _, e := range l.EntityInfo {
		confidences = append(confidences, e.Confidence)
	}
	for _, e := range l.PropertyInfo {
		confidences = append(confidences, e.Confidence)
	}
	for _, e := range l.Others {
		confidences = append(confidences, e.Confidence)
	}
	return utils.Min(confidences...)
}
