// Package messages defines shared message types used by both the installer
// state machine and the system command layer.
package messages

import "strings"

// ErrorKind classifies the severity of an installation error.
type ErrorKind int

const (
	// ErrorCritical means the installation cannot continue.
	ErrorCritical ErrorKind = iota
	// ErrorRetryable means the operation can be attempted again.
	ErrorRetryable
	// ErrorWarning is non-fatal; log it and continue.
	ErrorWarning
)

// ClassifyError inspects an error message and returns the appropriate ErrorKind.
func ClassifyError(err error) ErrorKind {
	if err == nil {
		return ErrorWarning
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "no space left"):
		return ErrorCritical
	case strings.Contains(msg, "permission denied"):
		return ErrorCritical
	case strings.Contains(msg, "read-only file system"):
		return ErrorCritical
	case strings.Contains(msg, "network"):
		return ErrorRetryable
	case strings.Contains(msg, "timeout"):
		return ErrorRetryable
	case strings.Contains(msg, "connection"):
		return ErrorRetryable
	case strings.Contains(msg, "temporary failure"):
		return ErrorRetryable
	default:
		return ErrorCritical
	}
}

// CmdOutputMsg carries a single line of stdout/stderr from a running command.
type CmdOutputMsg struct {
	Line string
}

// CmdDoneMsg signals that a background command completed successfully.
type CmdDoneMsg struct {
	Output string
}

// CmdErrorMsg signals that a background command failed.
type CmdErrorMsg struct {
	Err  error
	Kind ErrorKind
}

// StepForwardMsg advances the installer to the next step.
type StepForwardMsg struct{}

// StepBackMsg returns to the previous step.
type StepBackMsg struct{}

// EnvDetectedMsg carries the results of environment detection.
type EnvDetectedMsg struct {
	IsUEFI   bool
	Internet bool
}

// DiskInfo holds basic information about a block device.
type DiskInfo struct {
	Name  string // "sda"
	Path  string // "/dev/sda"
	Size  string // "256G"
	Model string // "Samsung SSD 860"
}

// DisksLoadedMsg carries a list of detected disks.
type DisksLoadedMsg struct {
	Disks []DiskInfo
}

// RetryMsg triggers a retry of the current step's background operation.
type RetryMsg struct{}
