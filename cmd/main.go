package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shailesh097/fedora-tweed-go/internal/distro"
	"github.com/shailesh097/fedora-tweed-go/internal/distro/arch"
	"github.com/shailesh097/fedora-tweed-go/internal/distro/fedora"
	"github.com/shailesh097/fedora-tweed-go/internal/logger"
)

func main() {
	log, cleanup, err := logger.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	defer cleanup()

	detectedOS, err := distro.DetectOS()
	if err != nil {
		log.Error("Unable to detect the linux distribution: %v", err)
		os.Exit(1)
	}

	log.Info("Detected %s as linux distribution!", detectedOS)

	if !promptContinue(detectedOS) {
		log.Warn("Setup process stopped!")
		os.Exit(0)
	}

	d, err := resolveDistro(detectedOS)
	if err != nil {
		log.Error("%v", err)
		os.Exit(1)
	}

	if err := d.Setup(); err != nil {
		log.Error("Setup failed: %v", err)
		os.Exit(1)
	}

}

func resolveDistro(osID string) (distro.Distro, error) {
	switch osID {
	case "fedora":
		return fedora.New(), nil
	case "arch":
		return arch.New(), nil
	default:
		return nil, fmt.Errorf("Unsupported operating system: %s", osID)
	}
}

// continue with asking prompt in case of invalid input
func promptContinue(osName string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Do you want to continue setting up %s? [y/n]: ", osName)

	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	switch strings.TrimSpace(strings.ToLower(input)) {
	case "y":
		return true
	case "n":
		return false
	default:
		fmt.Println("Invalid input! Please enter 'y' or 'n'.")
		return promptContinue(osName) // call itself again if invalid
	}
}
