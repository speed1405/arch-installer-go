package installer

import "github.com/speed1405/arch-installer-go/messages"

// Re-export message types from the shared messages package so that
// existing code within this package can reference them without a prefix.

// ErrorKind classifies the severity of an installation error.
type ErrorKind = messages.ErrorKind

const (
	ErrorCritical  = messages.ErrorCritical
	ErrorRetryable = messages.ErrorRetryable
	ErrorWarning   = messages.ErrorWarning
)

// ClassifyError inspects an error and returns the appropriate ErrorKind.
var ClassifyError = messages.ClassifyError

// CmdOutputMsg carries a single line of stdout/stderr from a running command.
type CmdOutputMsg = messages.CmdOutputMsg

// CmdDoneMsg signals that a background command completed successfully.
type CmdDoneMsg = messages.CmdDoneMsg

// CmdErrorMsg signals that a background command failed.
type CmdErrorMsg = messages.CmdErrorMsg

// StepForwardMsg advances the installer to the next step.
type StepForwardMsg = messages.StepForwardMsg

// StepBackMsg returns to the previous step.
type StepBackMsg = messages.StepBackMsg

// EnvDetectedMsg carries the results of environment detection.
type EnvDetectedMsg = messages.EnvDetectedMsg

// DiskInfo holds basic information about a block device.
type DiskInfo = messages.DiskInfo

// DisksLoadedMsg carries a list of detected disks.
type DisksLoadedMsg = messages.DisksLoadedMsg

// RetryMsg triggers a retry of the current step's background operation.
type RetryMsg = messages.RetryMsg
