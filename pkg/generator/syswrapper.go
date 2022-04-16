package generator

import (
	"os"
	"os/exec"
)

type sysWrapper interface {
	execCommand(name string, arg string) ([]byte, error)
	osGetwd() (string, error)
	osReadFile(name string) ([]byte, error)
}

type sysWrap struct{}

func (s *sysWrap) execCommand(name string, arg string) ([]byte, error) {
	return exec.Command(name, arg).Output()
}

func (s *sysWrap) osGetwd() (string, error) {
	return os.Getwd()
}

func (s *sysWrap) osReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
