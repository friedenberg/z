package metadata

import (
	"encoding/json"
	"regexp"
	"strings"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

var TagPrefixRegexp *regexp.Regexp

func init() {
	TagPrefixRegexp = regexp.MustCompile(`^\w{1,2}-\S`)
}

func MakeMetadata() (m Metadata) {
	m = Metadata{}
	m.init()
	return
}

type Metadata struct {
	description     string
	allTags         TagSet
	stringTags      TagSet
	searchMatchTags TagSet
	file            File
	url             *Url
}

func (m *Metadata) Set(s string) (err error) {
	var tags []string
	err = yaml.Unmarshal([]byte(s), &tags)

	if err != nil {
		err = xerrors.Errorf("parse metadata: %w", err)
		return
	}

	err = m.SetStringTags(tags)

	if err != nil {
		return
	}

	return
}

func (m *Metadata) init() {
	m.description = ""
	m.allTags = MakeTagSet()
	m.stringTags = MakeTagSet()
	m.searchMatchTags = MakeTagSet()
	m.file = nil
	m.url = nil
}

func (m *Metadata) SetStringTags(tags []string) (err error) {
	m.init()

	for i, t := range tags {
		if i == 0 && !TagPrefixRegexp.MatchString(t) {
			m.description = t
			continue
		}

		err = m.addStringTag(t)

		if err != nil {
			return
		}
	}

	return
}

func (m *Metadata) AddStringTags(t ...string) (err error) {
	for _, t1 := range t {
		err = m.addStringTag(t1)

		if err != nil {
			return
		}
	}

	return
}

func (m *Metadata) addStringTag(t string) (err error) {
	t1, err := makeTag(t)

	if err != nil || t1 == nil {
		return
	}

	switch t2 := t1.(type) {
	case *NewFile:
		if m.file != nil {
			m.allTags.Del(m.file.Tag())
		}

		m.file = t2
	case *LocalFile:
		if m.file != nil {
			m.allTags.Del(m.file.Tag())
		}

		m.file = t2

	case *Url:
		m.url = t2

	case *Tag:
		m.stringTags.Add(t2)

	default:
		panic(xerrors.Errorf("uncaught format for tag: '%s'", t))
	}

	m.allTags.Add(t1)
	m.searchMatchTags.Merge(t1.SearchMatchTags())

	return
}

func (m Metadata) Description() string {
	return m.description
}

func (m *Metadata) SetDescription(d string) {
	m.description = strings.TrimSpace(d)
}

func (m Metadata) StringTags() (s TagSet) {
	return m.stringTags
}

func (m Metadata) SearchMatchTags() (s TagSet) {
	s = m.searchMatchTags

	return
}

func (m Metadata) SearchMatchTagStrings() (t []string) {
	t = m.searchMatchTags.Strings()

	return
}

func (m Metadata) Url() (u Url, ok bool) {
	if ok = m.url != nil; ok {
		u = *m.url
	}

	return
}

func (m *Metadata) SetUrl(u Url) {
	old := m.url

	if old != nil {
		ot := old.Tag()
		m.allTags.Del(ot)
	}

	m.allTags.Add(&u)
	m.url = &u
}

func (m Metadata) HasFile() (ok bool) {
	ok = m.file != nil
	return
}

func (m *Metadata) SetFile(fd *NewFile) {
	m.file = fd
}

func (m Metadata) File() (f File) {
	f = m.file
	return
}

func (m Metadata) NewFile() (nf *NewFile, ok bool) {
	nf, ok = m.file.(*NewFile)
	return
}

func (m Metadata) LocalFile() (lf *LocalFile, ok bool) {
	lf, ok = m.file.(*LocalFile)
	return
}

func (m Metadata) TagSet() (s TagSet) {
	s = m.allTags
	return
}

func (m Metadata) Tags() (r []ITag) {
	r = m.allTags.Tags()
	return
}

func (m Metadata) TagStrings() (r []string) {
	r = m.allTags.Strings()
	return
}

func (m Metadata) Match(q string) bool {
	if m.TagSet().Match(q) {
		return true
	}

	if strings.Contains(m.description, q) {
		return true
	}

	return false
}

type jsonMetadata struct {
	Description string
	AllTags     []ITag
	StringTags  []ITag
	LocalFile   *LocalFile
	RemoteFiles []LocalFile
	Url         *Url
}

func (m *Metadata) UnmarshalJSON(b []byte) error {
	var t []string

	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	m.SetStringTags(t)

	return nil
}

func (m Metadata) MarshalJSON() (b []byte, err error) {
	b, err = json.Marshal(m.ToSortedTags())
	return
}
