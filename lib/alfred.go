package lib

import (
	"encoding/json"
	"net/url"
	"path"
	"strings"
	"time"
)

type ZettelAlfredItem struct {
	Title        string               `json:"title"`
	Arg          string               `json:"arg"`
	Subtitle     string               `json:"subtitle"`
	Match        string               `json:"match"`
	Icon         ZettelAlfredItemIcon `json:"icon"`
	Uid          string               `json:"uid"`
	ItemType     string               `json:"type"`
	QuicklookUrl string               `json:"quicklookurl"`
	Text         ZettelAlfredItemText `json:"text"`
}

type ZettelAlfredItemText struct {
	Copy string `json:"copy"`
}

type ZettelAlfredItemIcon struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

func MakeMatches(z *Zettel) string {
	m := z.IndexData

	join := func(s ...[]string) string {
		sb := strings.Builder{}

		for _, a := range s {
			for _, b := range a {
				sb.WriteString(b)
				sb.WriteString(" ")
			}
		}

		return sb.String()
	}

	t, err := TimeFromPath(z.Path)

	if err != nil {
		//TODO
	}

	day := t.Format("2006-01-02")

	base := []string{
		m.Description,
		"w:" + day,
	}

	if z.HasUrl() {
		url, err := url.Parse(m.Url)

		if err == nil {
			base = append(base, "d:"+url.Hostname())
		}

		base = append(base, "h:u")
	}

	if z.HasFile() {
		base = append(base, "h:f")
	}

	today := time.Now()

	if today.Format("2006-01-02") == day {
		base = append(base, "w:today")
	}

	return join(
		base,
		WithPrefix(Split(m.Areas, "-"), "a:"),
		WithPrefix(Split(m.Projects, "-"), "p:"),
		WithPrefix(Split(m.Tags, "-"), "t:"),
	)
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
	addMany(WithPrefix(z.IndexData.Areas, "a:"))
	addMany(WithPrefix(z.IndexData.Projects, "p:"))
	addMany(WithPrefix(z.IndexData.Tags, "t:"))

	return strings.Join(el, ", ")
}

func (z *Zettel) AddIcon() {
	getIconSuffix := func() string {
		if z.HasUrl() {
			return "pb"
		}

		return "note"
	}

	getIcon := func() (alfredIcon ZettelAlfredItemIcon) {
		if z.HasFile() {
			alfredIcon.Path = z.IndexData.File
			alfredIcon.Type = "fileicon"
			return
		}

		alfredIcon.Path = "icon-kind-" + getIconSuffix() + ".png"

		return
	}

	z.AlfredData.Item.Icon = getIcon()

	return
}

func (z *Zettel) AddAlfredItem() (err error) {
	z.AlfredData.Item.Arg = z.Path
	z.AlfredData.Item.ItemType = "file"
	z.AlfredData.Item.Uid = strings.TrimSuffix(
		path.Base(z.Path),
		path.Ext(z.Path),
	)

	z.AlfredData.Item.QuicklookUrl = z.Path

	if z.HasFile() {
		z.AlfredData.Item.QuicklookUrl = z.IndexData.File
	}

	z.AlfredData.Item.Subtitle = MakeSubtitle(z)
	z.AlfredData.Item.Title = z.IndexData.Description
	z.AlfredData.Item.Match = MakeMatches(z)
	z.AlfredData.Item.Text = ZettelAlfredItemText{
		Copy: path.Base(z.Path),
	}
	z.AddIcon()
	return nil
}

func (z *Zettel) GenerateAlfredJson() (err error) {
	alfredItemJson, err := json.Marshal(z.AlfredData.Item)
	z.AlfredData.Json = string(alfredItemJson)
	return
}
