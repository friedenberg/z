package printer

import (
	"encoding/json"
	"os/exec"
	"strings"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type ZettelPrinter interface {
	Begin()
	PrintZettel(int, *lib.Zettel, error)
	End()
}

//   _   _       _ _
//  | \ | |_   _| | |
//  |  \| | | | | | |
//  | |\  | |_| | | |
//  |_| \_|\__,_|_|_|
//

type NullZettelPrinter struct{}

func (p *NullZettelPrinter) Begin() {}
func (p *NullZettelPrinter) End()   {}

func (p *NullZettelPrinter) PrintZettel(_ int, _ *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}
}

type multiplexPrintLine struct {
	i int
	z *lib.Zettel
	e error
}

type MultiplexingZettelPrinter struct {
	Printer   ZettelPrinter
	channel   chan multiplexPrintLine
	waitGroup *sync.WaitGroup
}

func (p *MultiplexingZettelPrinter) Begin() {
	p.channel = make(chan multiplexPrintLine)
	p.Printer.Begin()
	p.waitGroup = &sync.WaitGroup{}
	p.waitGroup.Add(1)

	go func() {
		defer p.waitGroup.Done()

		for l := range p.channel {
			p.Printer.PrintZettel(l.i, l.z, l.e)
		}
	}()
}

func (p *MultiplexingZettelPrinter) PrintZettel(i int, z *lib.Zettel, e error) {
	p.channel <- multiplexPrintLine{i, z, e}
}

func (p *MultiplexingZettelPrinter) End() {
	close(p.channel)
	p.waitGroup.Wait()
	p.Printer.End()
}

//   _____ _ _
//  |  ___(_) | ___ _ __   __ _ _ __ ___   ___
//  | |_  | | |/ _ \ '_ \ / _` | '_ ` _ \ / _ \
//  |  _| | | |  __/ | | | (_| | | | | | |  __/
//  |_|   |_|_|\___|_| |_|\__,_|_| |_| |_|\___|
//

type FilenameZettelPrinter struct{}

func (p *FilenameZettelPrinter) Begin() {}
func (p *FilenameZettelPrinter) End()   {}

func (p *FilenameZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	util.StdPrinterOut(z.Path)
}

//       _ ____   ___  _   _
//      | / ___| / _ \| \ | |
//   _  | \___ \| | | |  \| |
//  | |_| |___) | |_| | |\  |
//   \___/|____/ \___/|_| \_|
//

type JsonZettelPrinter struct{}

func (p *JsonZettelPrinter) Begin() {}
func (p *JsonZettelPrinter) End()   {}

func (p *JsonZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	b, errOut := json.Marshal(z.IndexData)

	if errOut != nil {
		util.StdPrinterErr(errOut)
		return
	}

	util.StdPrinterOut(string(b))
}

//      _    _  __              _
//     / \  | |/ _|_ __ ___  __| |
//    / _ \ | | |_| '__/ _ \/ _` |
//   / ___ \| |  _| | |  __/ (_| |
//  /_/   \_\_|_| |_|  \___|\__,_|
//

type AlfredJsonZettelPrinter struct {
	afterFirstPrint bool
	sync.Mutex
}

func (p *AlfredJsonZettelPrinter) Begin() {
	util.StdPrinterOut(`{"items":[`)
}

func (p *AlfredJsonZettelPrinter) shouldPrintComma() bool {
	p.Lock()
	defer p.Unlock()

	return p.afterFirstPrint
}

func (p *AlfredJsonZettelPrinter) setShouldPrintComma() {
	p.Lock()
	defer p.Unlock()

	p.afterFirstPrint = true
}

func (p *AlfredJsonZettelPrinter) PrintZettel(_ int, z *lib.Zettel, errIn error) {
	defer p.setShouldPrintComma()

	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	sb := strings.Builder{}
	if p.shouldPrintComma() {
		sb.WriteString(",")
		sb.WriteString("\n")
	}

	sb.WriteString(z.AlfredData.Json)
	util.StdPrinterOut(sb.String())
}

func (p *AlfredJsonZettelPrinter) End() {
	util.StdPrinterOut(`]}`)
}

//   _____      _ _
//  |  ___|   _| | |
//  | |_ | | | | | |
//  |  _|| |_| | | |
//  |_|   \__,_|_|_|
//

type FullZettelPrinter struct{}

func (p *FullZettelPrinter) Begin() {}
func (p *FullZettelPrinter) End()   {}

func (p *FullZettelPrinter) PrintZettel(_ int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	sb := &strings.Builder{}
	sb.WriteString(lib.MetadataStartSequence)
	sb.WriteString(z.Data.MetadataYaml)
	sb.WriteString(lib.MetadataEndSequence)
	sb.WriteString(z.Data.Body)
	util.StdPrinterOutf(sb.String())
}

//   _____                          _
//  |  ___|__  _ __ _ __ ___   __ _| |_
//  | |_ / _ \| '__| '_ ` _ \ / _` | __|
//  |  _| (_) | |  | | | | | | (_| | |_
//  |_|  \___/|_|  |_| |_| |_|\__,_|\__|
//

type FormatZettelPrinter struct {
	Formatter lib.Formatter
}

func (p *FormatZettelPrinter) Begin() {}
func (p *FormatZettelPrinter) End()   {}

func (p *FormatZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	//TODO should empty strings be printed?
	util.StdPrinterOut(p.Formatter.Format(z))
}

//      _        _   _
//     / \   ___| |_(_) ___  _ __
//    / _ \ / __| __| |/ _ \| '_ \
//   / ___ \ (__| |_| | (_) | | | |
//  /_/   \_\___|\__|_|\___/|_| |_|
//

type ActionZettelPrinter struct {
	Env                    *lib.Env
	ShouldEdit, ShouldOpen bool
	zettels                []*lib.Zettel
	files                  []string
	urls                   []string
}

func (p *ActionZettelPrinter) Begin() {}

func (p *ActionZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	p.zettels = append(p.zettels, z)

	if z.HasFile() {
		p.files = append(p.files, z.FilePath())
	}

	if z.HasUrl() {
		p.urls = append(p.urls, z.IndexData.Url)
	}
}

func (p *ActionZettelPrinter) End() {
	if p.ShouldEdit {
		p.openZettels()
	}

	if p.ShouldOpen {
		p.openFiles()
		p.openUrls()
	}
}

func (p *ActionZettelPrinter) openZettels() {
	if len(p.zettels) == 0 {
		return
	}

	zettelFiles := make([]string, len(p.zettels))

	for i, z := range p.zettels {
		zettelFiles[i] = z.Path
	}

	args := []string{"-p"}

	cmd := exec.Command(
		"mvim",
		append(args, zettelFiles...)...,
	)

	cmd.Run()
}

func (p *ActionZettelPrinter) openFiles() {
	if len(p.files) == 0 {
		return
	}

	cmd := exec.Command(
		"open",
		p.files...,
	)

	cmd.Run()
}

func (p *ActionZettelPrinter) openUrls() {
	if len(p.urls) == 0 {
		return
	}

	args := []string{
		"-na",
		"Google Chrome",
		"--args",
		"--new-window",
	}

	cmd := exec.Command(
		"open",
		append(args, p.urls...)...,
	)

	cmd.Run()
	//TODO return errors
}
