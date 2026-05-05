package arch

import "fmt"

type Arch struct{}

func New() *Arch {
	return &Arch{}
}

func (a *Arch) Setup() error {
	return fmt.Errorf("Support for arch linux setup comming soon!!!")
}
