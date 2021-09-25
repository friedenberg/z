package files_guard

import (
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type openFilesGuard struct {
	channel chan struct{}
}

var openFilesGuardInstance *openFilesGuard

func init() {
	limitCmd := exec.Command("ulimit", "-S", "-n")
	output, err := limitCmd.Output()

	if err != nil {
		panic(err)
	}

	limitStr := strings.TrimSpace(string(output))

	limit, err := strconv.ParseInt(string(limitStr), 10, 64)

	if err != nil {
		panic(err)
	}

	openFilesGuardInstance = &openFilesGuard{
		//TODO subtract from unknowable
		channel: make(chan struct{}, limit),
	}
}

func (g *openFilesGuard) Lock() {
	g.channel <- struct{}{}
}

func (g *openFilesGuard) LockN(n int) {
	for i := 0; i < n; i++ {
		g.channel <- struct{}{}
	}
}

func (g *openFilesGuard) Unlock() {
	<-g.channel
}

func (g *openFilesGuard) UnlockN(n int) {
	for i := 0; i < n; i++ {
		<-g.channel
	}
}

func Create(s string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()
	f, err = os.Create(s)

	if err != nil {
		openFilesGuardInstance.Unlock()
	}

	return
}

func OpenFile(name string, flag int, perm os.FileMode) (f *os.File, err error) {
	openFilesGuardInstance.Lock()
	return os.OpenFile(name, flag, perm)
}

func Open(s string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()
	return os.Open(s)
}

func Close(f io.Closer) error {
	defer openFilesGuardInstance.Unlock()
	return f.Close()
}

func CombinedOutput(c *exec.Cmd) ([]byte, error) {
	openFilesGuardInstance.LockN(3)
	defer openFilesGuardInstance.UnlockN(3)

	return c.CombinedOutput()
}
