package model

import "strings"

type Form struct {
	Modules []map[string]interface{} `bson:"modules" json:"modules"`
}

func (form Form) HasVariables() bool {
	for _, v := range form.Modules {
		for _, n := range v {
			if str, ok := n.(string); ok {
				if strings.Contains(str, "$(") {
					return true
				}
			}
		}
	}
	return false
}
