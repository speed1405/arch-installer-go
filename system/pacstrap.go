package system

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

// BasePackages is the default package set installed by pacstrap.
var BasePackages = []string{
	"base",
	"base-devel",
	"linux",
	"linux-firmware",
	"networkmanager",
	"vim",
	"sudo",
}

// RunPacstrap installs the base system onto mountpoint.
func RunPacstrap(mountpoint string, packages []string) tea.Cmd {
	args := append([]string{"pacstrap", "-K", mountpoint}, packages...)
	return Execute(context.Background(), true, args...)
}
