package zettel

type Tag interface {
	Set(string) error
	Tag() string
}
