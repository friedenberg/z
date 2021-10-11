package util

import (
	"bytes"
	"io"
	"net/url"
	"os/exec"

	"github.com/friedenberg/z/util/files_guard"
)

func PandocSummaryFromUrl(u *url.URL) (summary string, err error) {
	//TODO-P4 description from title

	chromeCommand := exec.Command(
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"--headless",
		"--dump-dom",
		u.String(),
	)

	//TODO-P4 add filters
	pandocCommand := exec.Command(
		"pandoc",
		"-f",
		"html-native_divs-native_spans",
		"-t",
		"markdown-header_attributes-fenced_code_attributes-link_attributes",
	)

	r, w := io.Pipe()
	chromeCommand.Stdout = w
	pandocCommand.Stdin = r

	var md bytes.Buffer
	pandocCommand.Stdout = &md

	chromeCommand.Start()
	pandocCommand.Start()
	chromeCommand.Wait()
	files_guard.Close(w)
	pandocCommand.Wait()

	summary = md.String()

	return
}
