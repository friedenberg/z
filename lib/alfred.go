package lib

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/friedenberg/z/lib/alfred"
	"github.com/friedenberg/z/util"
)

type AlfredItem struct {
	Title        string               `json:"title"`
	Arg          string               `json:"arg"`
	Subtitle     string               `json:"subtitle"`
	Match        string               `json:"match"`
	Icon         AlfredItemIcon `json:"icon"`
	Uid          string               `json:"uid"`
	ItemType     string               `json:"type"`
	QuicklookUrl string               `json:"quicklookurl"`
	Text         AlfredItemText `json:"text"`
}

type AlfredItemText struct {
	Copy string `json:"copy"`
}

type AlfredItemIcon struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type AlfredItemFormat struct {
	Title    FormatFunc
	Arg      FormatFunc
	Subtitle FormatFunc
	// Match        FormatFunc
	Text         FormatFunc
	QuicklookUrl FormatFunc
	IconType     FormatFunc
	IconPath     FormatFunc
}

func MakeMatches(z *Zettel) string {
	//TODO add more variations and match against item format
	//e.g., Project: 2020-zettel -> p:2020-zettel, p:2020, 2020, zettel
	m := z.IndexData
	sb := &strings.Builder{}

	addMatch := func(s string) {
		sb.WriteString(s)
		sb.WriteString(" ")
	}

	t, err := TimeFromPath(z.Path)

	if err != nil {
		panic(fmt.Errorf("make alfred match field: %w", err))
	}

	addMatch(m.Description)

	day := t.Format("2006-01-02")

	addMatch("w-" + day)

	if z.HasUrl() {
		url, err := url.Parse(m.Url)

		if err == nil {
			addMatch("d-" + url.Hostname())
		}

		addMatch("h-u")
	}

	if z.HasFile() {
		addMatch("h-f")
	}

	today := time.Now()

	if today.Format("2006-01-02") == day {
		addMatch("w-today")
	}

	for _, t := range m.Tags {
		for _, m := range util.ExpandTags(t) {
			addMatch(m)
		}
	}

	return sb.String()
}

func MakeSubtitle(z *Zettel) string {
	el := make([]string, 0)

	add := func(s string) {
		el = append(el, s)
	}

	addMany := func(s []string) {
		for _, v := range s {
			add(v)
		}
	}

	add(z.IndexData.Date)
	addMany(alfred.WithPrefix(z.IndexData.Tags, "t:"))

	return strings.Join(el, ", ")
}

func (z *Zettel) AddIcon() {
	getIconSuffix := func() string {
		if z.HasUrl() {
			return "pb"
		}

		return "note"
	}

	getIcon := func() (alfredIcon AlfredItemIcon) {
		if z.HasFile() {
			alfredIcon.Path = z.FilePath()
			alfredIcon.Type = "fileicon"
			return
		}

		alfredIcon.Path = "icon-kind-" + getIconSuffix() + ".png"

		return
	}

	z.AlfredData.Item.Icon = getIcon()

	return
}

func (z *Zettel) AddAlfredItem(f AlfredItemFormat) (err error) {
	z.AlfredData.Item.ItemType = "file"

	z.AlfredData.Item.Uid = strings.TrimSuffix(
		path.Base(z.Path),
		path.Ext(z.Path),
	)

	if f.QuicklookUrl != nil {
		z.AlfredData.Item.QuicklookUrl = f.QuicklookUrl(z)
	}

	if f.Arg != nil {
		z.AlfredData.Item.Arg = f.Arg(z)
	}

	if f.Subtitle != nil {
		z.AlfredData.Item.Subtitle = f.Subtitle(z)
	}

	if f.Title != nil {
		z.AlfredData.Item.Title = f.Title(z)
	}

	if f.Text != nil {
		z.AlfredData.Item.Text = AlfredItemText{
			Copy: f.Text(z),
		}
	}

	z.AlfredData.Item.Icon = AlfredItemIcon{}

	if f.IconType != nil {
		z.AlfredData.Item.Icon.Type = f.IconType(z)
	}

	if f.IconPath != nil {
		z.AlfredData.Item.Icon.Path = f.IconPath(z)
	}

	z.AlfredData.Item.Match = MakeMatches(z)
	// z.AddIcon()
	return nil
}

func (z *Zettel) GenerateAlfredJson() (err error) {
	alfredItemJson, err := json.Marshal(z.AlfredData.Item)
	z.AlfredData.Json = string(alfredItemJson)
	return
}
