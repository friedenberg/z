package commands

import (
	"path/filepath"
	"testing"
)

func TestMain(t *testing.T) {
	glob := filepath.Join("/Users/sasha/Zettelkasten", "*.md")
	processor, _ := MakeProcessor(glob, nil, &NullPutter{Channel: make(PutterChannel)})
	processor.Run()
}
