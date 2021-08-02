package lib

import "encoding/json"

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

func GenerateAlfredJson(i AlfredItem) (j string, err error) {
	alfredItemJson, err := json.Marshal(i)
	j = string(alfredItemJson)
	return
}
