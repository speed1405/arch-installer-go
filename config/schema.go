package config

// InstallConfig holds all user-selected installation parameters.
type InstallConfig struct {
	Disk            string      // e.g. "/dev/sda"
	PartitionScheme string      // "auto" | "manual"
	Partitions      []Partition
	Bootloader      string // "grub" | "systemd-boot"
	Hostname        string
	Locale          string // "en_US.UTF-8"
	Timezone        string // "America/New_York"
	Username        string
	Password        string // in-memory only, never persisted
	RootPassword    string // in-memory only, never persisted
	ExtraPackages   []string
	IsUEFI          bool
	Mirror          string
}

// Partition describes a single disk partition.
type Partition struct {
	Device     string // "/dev/sda1"
	MountPoint string // "/boot", "/", "/home"
	Filesystem string // "fat32", "ext4", "btrfs"
	Size       string // "512M", "rest"
	Encrypt    bool
}
