package printer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Rsync struct {
	Src, Dst   string
	rsyncCmd   *exec.Cmd
	rsyncStdin io.WriteCloser
}

func (p *Rsync) Begin() {
	p.rsyncCmd = exec.Command(
		"rsync",
		"--update",
		"--progress",
		"-0",
		"--files-from=/dev/stdin",
		p.Src,
		p.Dst,
	)

	p.rsyncCmd.Stdout = os.Stdout
	p.rsyncCmd.Stderr = os.Stderr

	r, w := io.Pipe()
	p.rsyncCmd.Stdin = r
	p.rsyncStdin = w

	p.rsyncCmd.Start()
}

func (p *Rsync) File(file string) {
	line := fmt.Sprintf("%s\x00", file)
	p.rsyncStdin.Write([]byte(line))
}

func (p *Rsync) End() {
	p.rsyncStdin.Close()
	p.rsyncCmd.Wait()
}
