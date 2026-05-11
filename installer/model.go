package installer

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/speed1405/arch-installer-go/config"
	"github.com/speed1405/arch-installer-go/system"
	"github.com/speed1405/arch-installer-go/ui"
	"github.com/speed1405/arch-installer-go/ui/components"
)

// Model is the root Bubble Tea model for the installer.
type Model struct {
	step      Step
	config    config.InstallConfig
	prevSteps []Step
	blocked   map[Step]bool

	// Screen-specific sub-models
	diskList components.DiskListModel
	partForm components.PartFormModel
	userForm components.UserFormModel
	confirm  components.ConfirmModel
	progress components.ProgressModel

	// Shared state
	logLines []string
	err      error
	errKind  ErrorKind
	width    int
	height   int
}

// New creates the root installer model with default configuration.
func New() Model {
	return Model{
		step:      StepWelcome,
		config:    config.DefaultConfig(),
		prevSteps: []Step{},
		blocked:   make(map[Step]bool),
		diskList:  components.NewDiskList(),
		partForm:  components.NewPartForm(),
		userForm:  components.NewUserForm(),
		confirm:   components.NewConfirm("Are you sure you want to partition the disk? This cannot be undone."),
		progress:  components.NewProgress(),
	}
}

// Init starts the first async command.
func (m Model) Init() tea.Cmd {
	return system.DetectEnvironment()
}

// Update handles all incoming messages and state transitions.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.diskList.SetSize(msg.Width, msg.Height)
		m.progress.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.step != StepError && m.step != StepComplete {
				return m, func() tea.Msg { return StepBackMsg{} }
			}
		}

	case EnvDetectedMsg:
		m.config.IsUEFI = msg.IsUEFI
		if m.step == StepDetectEnv {
			return m, func() tea.Msg { return StepForwardMsg{} }
		}
		return m, nil

	case DisksLoadedMsg:
		disks := make([]components.DiskDisplay, len(msg.Disks))
		for i, d := range msg.Disks {
			disks[i] = components.DiskDisplay{
				Name:  d.Name,
				Path:  d.Path,
				Size:  d.Size,
				Model: d.Model,
			}
		}
		m.diskList.SetDisks(disks)
		return m, nil

	case StepForwardMsg:
		next := nextStep(m.step)
		m.prevSteps = append(m.prevSteps, m.step)
		m.step = next
		if destructiveSteps[next] {
			m.blocked[next] = true
		}
		return m, m.initStep(next)

	case StepBackMsg:
		n := len(m.prevSteps)
		if n > 0 && !m.blocked[m.step] {
			prev := m.prevSteps[n-1]
			m.prevSteps = m.prevSteps[:n-1]
			m.step = prev
		}
		return m, nil

	case CmdOutputMsg:
		m.logLines = append(m.logLines, msg.Line)
		var cmd tea.Cmd
		m.progress, cmd = m.progress.Update(msg)
		return m, cmd

	case CmdDoneMsg:
		m.logLines = append(m.logLines, "✓ Done.")
		return m, func() tea.Msg { return StepForwardMsg{} }

	case CmdErrorMsg:
		m.err = msg.Err
		m.errKind = msg.Kind
		if msg.Kind == ErrorCritical {
			m.step = StepError
		}
		return m, nil

	case RetryMsg:
		return m, m.initStep(m.step)
	}

	return m.updateActiveComponent(msg)
}

// initStep returns the initial tea.Cmd for a newly entered step.
func (m *Model) initStep(s Step) tea.Cmd {
	switch s {
	case StepDetectEnv:
		return system.DetectEnvironment()
	case StepSelectDisk:
		return system.ListDisks()
	case StepPartitioning:
		return system.PartitionDisk(m.config)
	case StepInstallingBase:
		return system.RunPacstrap("/mnt", system.BasePackages)
	case StepFstabGen:
		return system.GenerateFstab()
	case StepTimezone:
		return system.SetTimezone(m.config.Timezone)
	case StepLocale:
		return system.SetLocale(m.config.Locale)
	case StepHostname:
		return system.SetHostname(m.config.Hostname)
	case StepUserConfig:
		return system.CreateUser(m.config.Username, m.config.Password, m.config.RootPassword)
	case StepBootloader:
		return system.InstallBootloader(m.config)
	}
	return nil
}

// updateActiveComponent delegates to whichever sub-model is currently active.
func (m Model) updateActiveComponent(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.step {
	case StepSelectDisk:
		var selected string
		m.diskList, selected, cmd = m.diskList.Update(msg)
		if selected != "" {
			m.config.Disk = selected
			return m, func() tea.Msg { return StepForwardMsg{} }
		}
	case StepPartitionConfig:
		var done bool
		m.partForm, done, cmd = m.partForm.Update(msg)
		if done {
			m.config.Partitions = m.partForm.Result()
			return m, func() tea.Msg { return StepForwardMsg{} }
		}
	case StepConfirmPartition:
		var confirmed bool
		m.confirm, confirmed, cmd = m.confirm.Update(msg)
		if confirmed {
			return m, func() tea.Msg { return StepForwardMsg{} }
		}
	case StepUserConfig:
		var done bool
		m.userForm, done, cmd = m.userForm.Update(msg)
		if done {
			res := m.userForm.Result()
			m.config.Hostname = res.Hostname
			m.config.Username = res.Username
			m.config.Password = res.Password
			m.config.RootPassword = res.RootPassword
			return m, func() tea.Msg { return StepForwardMsg{} }
		}
	case StepInstallingBase, StepPartitioning, StepFstabGen,
		StepTimezone, StepLocale, StepHostname, StepBootloader:
		m.progress, cmd = m.progress.Update(msg)
	}
	return m, cmd
}

// View renders the complete TUI for the current step.
func (m Model) View() string {
	header := ui.RenderHeader(StepName(m.step), m.step == StepError)
	footer := ui.RenderFooter(m.step == StepComplete || m.step == StepError)
	body := m.renderBody()

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m Model) renderBody() string {
	switch m.step {
	case StepWelcome:
		return ui.RenderWelcome()
	case StepDetectEnv:
		return ui.RenderSpinner("Detecting environment…")
	case StepSelectDisk:
		return m.diskList.View()
	case StepPartitionStrategy:
		return ui.RenderPartitionStrategy(m.config.PartitionScheme)
	case StepPartitionConfig:
		return m.partForm.View()
	case StepConfirmPartition:
		return m.confirm.View()
	case StepPartitioning:
		return m.progress.View()
	case StepSelectMirrors:
		return ui.RenderMirrorInput(m.config.Mirror)
	case StepInstallingBase:
		return m.progress.View()
	case StepFstabGen:
		return m.progress.View()
	case StepTimezone:
		return ui.RenderTextInput("Timezone", m.config.Timezone)
	case StepLocale:
		return ui.RenderTextInput("Locale", m.config.Locale)
	case StepHostname:
		return ui.RenderTextInput("Hostname", m.config.Hostname)
	case StepUserConfig:
		return m.userForm.View()
	case StepBootloader:
		return m.progress.View()
	case StepPostInstall:
		return ui.RenderPostInstall()
	case StepComplete:
		return ui.RenderComplete()
	case StepError:
		return ui.RenderError(m.err, m.errKind == ErrorRetryable)
	default:
		return fmt.Sprintf("Step: %s", StepName(m.step))
	}
}
