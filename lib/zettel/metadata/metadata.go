package metadata

import (
	"encoding/json"
	"strings"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

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
	newFile         *NewFile
	localFile       *File
	remoteFiles     TagSet
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
	m.newFile = nil
	m.localFile = nil
	m.remoteFiles = MakeTagSet()
	m.url = nil
}

func (m *Metadata) SetStringTags(tags []string) (err error) {
	m.init()

	for i, t := range tags {
		if i == 0 {
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
		if m.newFile != nil {
			err = xerrors.Errorf("already have new file, cannot add again")
			return
		}

		m.newFile = t2
	case *File:
		if t2.KastenName == "" {
			if m.localFile != nil {
				m.allTags.Del(m.localFile.Tag())
			}

			m.localFile = t2
		} else {
			m.remoteFiles.Add(t2)
		}

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
	ok = false

	_, ok = m.LocalFile()

	if ok {
		return
	}

	ok = len(m.RemoteFiles()) > 0

	return
}

func (m Metadata) NewFile() (fd NewFile, ok bool) {
	ok = m.newFile != nil

	if ok {
		fd = *m.newFile
	}

	return
}

func (m Metadata) LocalFile() (fd File, ok bool) {
	ok = m.localFile != nil

	if ok {
		fd = *m.localFile
	}

	return
}

func (m Metadata) RemoteFiles() (fds []File) {
	ts := m.remoteFiles.Tags()
	fds = make([]File, len(ts))

	for i, t := range ts {
		fds[i] = *(t.(*File))
	}

	return
}

func (m Metadata) Files() (fs []File) {
	fs = make([]File, 0, 1+m.remoteFiles.Len())

	if f, ok := m.LocalFile(); ok {
		fs = append(fs, f)
	}

	for _, f := range m.RemoteFiles() {
		fs = append(fs, f)
	}

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
	LocalFile   *File
	RemoteFiles []File
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
