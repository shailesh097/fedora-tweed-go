package fedora

import (
	"fmt"
	"os"

	"github.com/shailesh097/fedora-tweed-go/internal/logger"
)

type Fedora struct{}

func New() *Fedora {
	return &Fedora{}
}

func (f *Fedora) Setup() error {
	log, cleanup, err := logger.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v", err)
	}

	defer cleanup()

	log.Error("Error while setting up fedora!")
	return nil
}
