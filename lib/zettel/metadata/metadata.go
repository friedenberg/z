package metadata

import (
	"encoding/json"
	"regexp"
	"strings"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

var (
	RegexTag *regexp.Regexp
)

func init() {
	RegexTag = regexp.MustCompile(`^\w+-\w[^\s]+$`)
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
	localFile       *File
	remoteFiles     []File
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
	m.localFile = nil
	m.remoteFiles = nil
	m.url = nil
}

func (m *Metadata) SetStringTags(tags []string) (err error) {
	m.init()

	for i, t := range tags {
		if i == 0 && !RegexTag.MatchString(t) {
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
	if t == "" {
		return
	}

	var t1 ITag

	firstGroup := strings.Split(t, "-")[0]

	switch firstGroup {
	case "f":
		f := File{}
		err = f.Set(t)

		if err != nil {
			return
		}

		if f.KastenName == "" {
			m.localFile = &f
		} else {
			m.remoteFiles = append(m.remoteFiles, f)
		}

		t1 = &f
	case "u":
		u := Url{}
		err = u.Set(t)

		if err != nil {
			return
		}

		m.url = &u
		t1 = &u

	default:
		t2 := Tag(t)
		t1 = &t2
		m.stringTags.Add(t1)
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

func (m *Metadata) AddFile(fd File) {
	if fd.KastenName == "" {
		m.localFile = &fd
	} else {
		m.remoteFiles = append(m.remoteFiles, fd)
	}

	m.allTags.Add(&fd)
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

func (m Metadata) LocalFile() (fd File, ok bool) {
	ok = m.localFile != nil

	if ok {
		fd = *m.localFile
	}

	return
}

func (m Metadata) RemoteFiles() (fds []File) {
	fds = m.remoteFiles

	return
}

func (m Metadata) Files() (fs []File) {
	fs = make([]File, 0, 1+len(m.remoteFiles))

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
