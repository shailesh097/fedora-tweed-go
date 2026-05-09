package fedora

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/shailesh097/fedora-tweed-go/internal/logger"
)

type Fedora struct {
	log *logger.Logger
}

func New() *Fedora {
	return &Fedora{}
}

func (f *Fedora) Setup() error {
	log, cleanup, err := logger.New()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	defer cleanup()
	f.log = log

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Configuring DNF", f.configureDNF},
		{"Creating Workspace Directory", f.createWorkspace},
		{"Setting up Power Profile", f.setupPowerProfile},
		{"Updating System", f.updateSystem},
		{"Enabling RPM Fusion Repositories", f.enableRPMFusion},
	}

	for _, step := range steps {
		f.log.Info(step.name + "...")
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s failed: %w", step.name, err)
		}
	}

	f.log.Message("Fedora Setup Completed!")
	return nil
}

func (f *Fedora) enableRPMFusion() error {
	cmd := exec.Command("bash", "-c", `sudo dnf install -y \
	https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm \
		https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm`)
	return cmd.Run()
}

func (f *Fedora) updateSystem() error {
	return exec.Command("sudo", "dnf", "update", "-y").Run()
}

func (f *Fedora) setupPowerProfile() error {
	settings := [][]string{
		{"org.gnome.desktop.session", "idle-delay", "960"},
		{"org.gnome.settings-daemon.plugins.power", "sleep-inactive-ac-type", "nothing"},
		{"org.gnome.settings-daemon.plugins.power", "power-button-action", "nothing"},
	}

	for _, s := range settings {
		cmd := exec.Command("gsettings", "set", s[0], s[1], s[2])
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (f *Fedora) configureDNF() error {
	configs := []string{
		"defaultyes=True",
		"max_parallel_downloads=10",
	}

	for _, config := range configs {
		cmd := exec.Command("bash", "-c",
			fmt.Sprintf(`grep -q "^%s" /etc/dnf/dnf.conf ||  echo "%s" | sudo tee -a /etc/dnf/dnf.conf > /dev/null`, config, config))
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (f *Fedora) createWorkspace() error {
	workspace := fmt.Sprintf("%s/Workspace", os.Getenv("HOME"))
	return os.MkdirAll(workspace, 0755)
}
