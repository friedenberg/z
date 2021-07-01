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
	m := z.Metadata

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

	base := []string{
		m.Description,
		"w:" + m.Date,
		"k:" + m.Kind,
	}

	if m.Kind == "pb" {
		url, err := url.Parse(m.Url)

		if err == nil {
			base = append(base, "d:"+url.Hostname())
		}
	}

	t := time.Now()

	if t.Format("2006-01-02") == m.Date {
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

	add(z.Metadata.Date)
	addMany(WithPrefix(z.Metadata.Areas, "a:"))
	addMany(WithPrefix(z.Metadata.Projects, "p:"))
	addMany(WithPrefix(z.Metadata.Tags, "t:"))

	return strings.Join(el, ", ")
}

func (z *Zettel) AddIcon() {
	getIconSuffix := func() string {
		switch z.Metadata.Kind {
		default:
			return ""
		case "bs", "brainstorm":
			return "bs"
		case "mn", "meeting notes":
			return "mn"
		case "pb", "link", "url":
			return "pb"
		case "em", "email":
			return "email"
		}
	}

	getIcon := func() ZettelAlfredItemIcon {
		iconPath := ""
		if iconSuffix := getIconSuffix(); iconSuffix != "" {
			iconPath = "icon-kind-" + getIconSuffix() + ".png"
		}

		obj := ZettelAlfredItemIcon{
			Path: iconPath,
		}

		return obj
	}

	z.AlfredItem.Icon = getIcon()

	return
}

func (z *Zettel) AddAlfredItem() (err error) {
	z.AlfredItem.Arg = z.Path
	z.AlfredItem.ItemType = "file"
	z.AlfredItem.Uid = strings.TrimSuffix(
		path.Base(z.Path),
		path.Ext(z.Path),
	)
	z.AlfredItem.QuicklookUrl = z.Path
	z.AlfredItem.Subtitle = MakeSubtitle(z)
	z.AlfredItem.Title = z.Metadata.Description
	z.AlfredItem.Match = MakeMatches(z)
	z.AlfredItem.Text = ZettelAlfredItemText{
		Copy: path.Base(z.Path),
	}
	z.AddIcon()
	return nil
}

func (z *Zettel) GenerateAlfredJson() (err error) {
	alfredItemJson, err := json.Marshal(z.AlfredItem)
	z.AlfredItemJson = string(alfredItemJson)
	return
}
