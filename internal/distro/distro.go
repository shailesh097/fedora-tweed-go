package distro

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shailesh097/fedora-tweed-go/internal/distro/arch"
	"github.com/shailesh097/fedora-tweed-go/internal/distro/fedora"
)

type Distro interface {
	Setup() error
}

func DetectOS() (string, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "", fmt.Errorf("Could not open /etc/os-release: %w", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file) // reads a file line by line
	for scanner.Scan() {
		line := scanner.Text()
		if trimmedString, found := strings.CutPrefix(line, "ID="); found {
			id := strings.Trim(trimmedString, `"`)
			return strings.ToLower(id), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading /etc/os-release: %w", err)
	}

	return "", fmt.Errorf("could not find ID field in /etc/os-release")
}

func ResolveDistro(osID string) (Distro, error) {
	switch osID {
	case "fedora":
		return fedora.New(), nil
	case "arch":
		return arch.New(), nil
	default:
		return nil, fmt.Errorf("Unsupported operating system: %s", osID)
	}
}
