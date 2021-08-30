package lib

import (
	"encoding/json"
	"strings"
)

type AlfredItem struct {
	Title        string         `json:"title"`
	Arg          string         `json:"arg"`
	Subtitle     string         `json:"subtitle"`
	Match        string         `json:"match"`
	Icon         AlfredItemIcon `json:"icon"`
	Uid          string         `json:"uid"`
	ItemType     string         `json:"type"`
	QuicklookUrl string         `json:"quicklookurl"`
	Text         AlfredItemText `json:"text"`
}

type AlfredItemText struct {
	Copy string `json:"copy"`
}

type AlfredItemIcon struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

func GenerateAlfredItemsJson(i []AlfredItem) (j string, err error) {
	sb := &strings.Builder{}

	for idx, v := range i {
		if idx > 0 {
			sb.WriteString(",")
		}

		var alfredItemJson []byte
		alfredItemJson, err = json.Marshal(v)

		if err != nil {
			return
		}

		sb.WriteString(string(alfredItemJson))
	}

	j = sb.String()
	return
}
