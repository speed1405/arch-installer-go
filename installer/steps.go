package installer

// Step represents a single installation phase in the state machine.
type Step int

const (
	StepWelcome Step = iota
	StepDetectEnv
	StepSelectDisk
	StepPartitionStrategy
	StepPartitionConfig
	StepConfirmPartition
	StepPartitioning
	StepSelectMirrors
	StepInstallingBase
	StepFstabGen
	StepTimezone
	StepLocale
	StepHostname
	StepUserConfig
	StepBootloader
	StepPostInstall
	StepComplete
	StepError
)

// stepNames maps each Step to a human-readable title shown in the header.
var stepNames = map[Step]string{
	StepWelcome:           "Welcome",
	StepDetectEnv:         "Environment Detection",
	StepSelectDisk:        "Select Disk",
	StepPartitionStrategy: "Partition Strategy",
	StepPartitionConfig:   "Partition Configuration",
	StepConfirmPartition:  "Confirm Partitioning",
	StepPartitioning:      "Partitioning Disk",
	StepSelectMirrors:     "Select Mirrors",
	StepInstallingBase:    "Installing Base System",
	StepFstabGen:          "Generating Fstab",
	StepTimezone:          "Timezone",
	StepLocale:            "Locale",
	StepHostname:          "Hostname",
	StepUserConfig:        "User Configuration",
	StepBootloader:        "Bootloader",
	StepPostInstall:       "Post-Install Extras",
	StepComplete:          "Installation Complete",
	StepError:             "Error",
}

// StepName returns the display name for a step.
func StepName(s Step) string {
	if name, ok := stepNames[s]; ok {
		return name
	}
	return "Unknown"
}

// nextStep returns the logically next step in the installation sequence.
func nextStep(s Step) Step {
	switch s {
	case StepWelcome:
		return StepDetectEnv
	case StepDetectEnv:
		return StepSelectDisk
	case StepSelectDisk:
		return StepPartitionStrategy
	case StepPartitionStrategy:
		return StepPartitionConfig
	case StepPartitionConfig:
		return StepConfirmPartition
	case StepConfirmPartition:
		return StepPartitioning
	case StepPartitioning:
		return StepSelectMirrors
	case StepSelectMirrors:
		return StepInstallingBase
	case StepInstallingBase:
		return StepFstabGen
	case StepFstabGen:
		return StepTimezone
	case StepTimezone:
		return StepLocale
	case StepLocale:
		return StepHostname
	case StepHostname:
		return StepUserConfig
	case StepUserConfig:
		return StepBootloader
	case StepBootloader:
		return StepPostInstall
	case StepPostInstall:
		return StepComplete
	default:
		return s
	}
}

// destructiveSteps are steps after which backward navigation is blocked.
var destructiveSteps = map[Step]bool{
	StepPartitioning:   true,
	StepInstallingBase: true,
	StepFstabGen:       true,
}
