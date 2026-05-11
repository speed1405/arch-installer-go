package config

// DefaultConfig returns a config pre-populated with sensible defaults.
func DefaultConfig() InstallConfig {
	return InstallConfig{
		PartitionScheme: "auto",
		Bootloader:      "grub",
		Locale:          "en_US.UTF-8",
		Timezone:        "UTC",
		Hostname:        "archlinux",
		ExtraPackages:   []string{},
	}
}
