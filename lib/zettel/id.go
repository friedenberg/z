package zettel

import (
	"strconv"

	"github.com/friedenberg/z/util"
)

type Id int64

func IdFromString(s string) (id Id, err error) {
	base := util.BaseNameNoSuffix(s)
	in, err := strconv.ParseInt(base, 10, 64)
	return Id(in), err
}

func (id Id) Int() int64 {
	return int64(id)
}

func (id Id) String() string {
	return strconv.FormatInt(id.Int(), 10)
}
