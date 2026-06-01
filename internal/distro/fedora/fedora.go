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
		{"Enabling Flatpak", f.enableFlatpak},
		{"Installing essential tools", f.installEssentialTools},
		{"Setting Catpuccin-Mocha theme for kitty", f.setKittyTheme},
		{"Configuring git", f.configureGit},
		// {"Installing Nvchad", f.insatallNvchad},
		{"Installing Extension Manager", f.installExtensionManager},
		// {"Setting up dotfiles", f.setupDotfiles} -- handled in dotfiles.go
		{"Installing VScode", f.installVscode},
		{"Installing Brave Browser", f.installBrave},
		{"Installing Obsidian", f.installObsidian},
		{"Installing Spotify", f.installSpotify},
		{"Enabling Maximize and Minimize buttons", f.enableTitlebarButtons},
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

func (f *Fedora) enableTitlebarButtons() error {
	return exec.Command("gsettings", "set", "org.gnome.desktop.wm.preferences",
		"button-layout", "appmenu:minimize,maximize,close").Run()
}
func (f *Fedora) installObsidian() error {
	return exec.Command("flatpak", "install", "-y", "flathub", "md.obsidian.Obsidian").Run()
}

func (f *Fedora) installSpotify() error {
	return exec.Command("flatpak", "install", "-y", "flathub", "com.spotify.Client").Run()
}
func (f *Fedora) installBrave() error {
	if err := exec.Command("sudo", "dnf", "install", "-y", "dnf-plugins-core").Run(); err != nil {
		return err
	}

	repoUrl := "https://brave-browser-rpm-release.s3.brave.com/brave-browser.repo"

	if err := exec.Command("sudo", "dnf", "config-manager", "addrepo", "--overwrite", "--from-repofile="+repoUrl).Run(); err != nil {
		return nil
	}

	return exec.Command("sudo", "dnf", "install", "brave-browser", "-y").Run()
}

func (f *Fedora) installVscode() error {
	steps := [][]string{
		{"sudo", "rpm", "--import", "https://packages.microsoft.com/keys/microsoft.asc"},
		{"sudo", "sh", "-c", `echo -e "[code]\nname=Visual Studio Code\nbaseurl=https://packages.microsoft.com/yumrepos/vscode\nenabled=1\ngpgcheck=1\ngpgkey=https://packages.microsoft.com/keys/microsoft.asc" > /etc/yum.repos.d/vscode.repo`},
		{"sudo", "dnf", "check-update"},
		{"sudo", "dnf", "install", "-y", "code"},
	}

	for _, step := range steps {
		if err := exec.Command(step[0], step[1:]...).Run(); err != nil {
			return err
		}
	}
	return nil
}

// func (f *Fedora) insatallNvchad() error {
// 	configPath := fmt.Sprintf("%s/.config/nvim", os.Getenv("HOME"))
// 	return exec.Command("git", "clone", "https://github.com/NvChad/starter", configPath).Run()
// }

func (f *Fedora) installExtensionManager() error {
	return exec.Command("flatpak", "install", "-y", "flathub", "com.mattjakeman.ExtensionManager").Run()
}

func (f *Fedora) configureGit() error {
	if err := exec.Command("git", "config", "--global", "user.name", "shailesh097").Run(); err != nil {
		return err
	}
	return exec.Command("git", "config", "--global", "user.email", "sailesh.pokharel.234@gmail.com").Run()
}

func (f *Fedora) setKittyTheme() error {
	return exec.Command("kitty", "+kitten", "themes", "Catppuccin-Mocha").Run()
}

func (f *Fedora) installEssentialTools() error {
	packages := []string{
		"curl", "wget", "git", "neovim", "fzf", "conky", "kitty",
		"fish", "nvtop", "btop", "fastfetch", "npm", "gnome-tweaks",
		"discord", "vlc", "gparted", "bash", "shc", "zoxide", "tldr",
		"xdotool", "golang",
	}

	args := append([]string{"dnf", "install", "-y", "--skip-unavailable"}, packages...)
	return exec.Command("sudo", args...).Run()
}

func (f *Fedora) enableFlatpak() error {
	if err := exec.Command("sudo", "dnf", "install", "-y", "flatpak").Run(); err != nil {
		return err
	}
	return exec.Command("flatpak", "remote-add", "--if-not-exists", "flathub",
		"https://dl.flathub.org/repo/flathub.flatpakrepo").Run()
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
		"defaultyes=true",
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
