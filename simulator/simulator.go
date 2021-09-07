package simulator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/env"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-utils/log"
	sim "github.com/bitrise-io/go-xcode/simulator"
)

// Simulator ...
type Simulator sim.InfoModel

// Manager ...
type Manager interface {
	GetLatestSimulatorAndVersion(osName, deviceName string) (Simulator, string, error)
	GetSimulator(osNameAndVersion, deviceName string) (Simulator, error)
	LaunchSimulator(simulatorID string, xcodebuildMajorVersion int) error

	ResetLaunchServices() error
	SimulatorBoot(id string) error
	SimulatorEnableVerboseLog(id string) error
	SimulatorCollectDiagnostics() (string, error)
	SimulatorShutdown(id string) error
	SimulatorDiagnosticsName() (string, error)
}

type manager struct {
}

// NewManager ...
func NewManager() Manager {
	return manager{}
}

func (m manager) GetLatestSimulatorAndVersion(osName, deviceName string) (Simulator, string, error) {
	info, ver, err := sim.GetLatestSimulatorInfoAndVersion(osName, deviceName)
	return Simulator(info), ver, err
}

func (m manager) GetSimulator(osNameAndVersion, deviceName string) (Simulator, error) {
	info, err := sim.GetSimulatorInfo(osNameAndVersion, deviceName)
	return Simulator(info), err
}

func (m manager) LaunchSimulator(simulatorID string, xcodebuildMajorVersion int) error {
	return sim.BootSimulator(simulatorID, xcodebuildMajorVersion)
}

// Reset launch services database to avoid Big Sur's sporadic failure to find the Simulator App
// The following error is printed when this happens: "kLSNoExecutableErr: The executable is missing"
// Details:
// - https://stackoverflow.com/questions/2182040/the-application-cannot-be-opened-because-its-executable-is-missing/16546673#16546673
// - https://ss64.com/osx/lsregister.html
func (m manager) ResetLaunchServices() error {
	f := command.NewFactory(env.NewRepository())
	cmd := f.Create("sw_vers", []string{"-productVersion"}, nil)

	macOSVersion, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return err
	}

	if strings.HasPrefix(macOSVersion, "11.") { // It's Big Sur
		cmd := f.Create("xcode-select", []string{"--print-path"}, nil)
		xcodeDevDirPath, err := cmd.RunAndReturnTrimmedCombinedOutput()
		if err != nil {
			return err
		}

		simulatorAppPath := filepath.Join(xcodeDevDirPath, "Applications", "Simulator.app")

		cmdString := "/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister"
		cmd = f.Create(cmdString, []string{"-f", simulatorAppPath}, nil)

		log.Infof("Applying launch services reset workaround before booting simulator")
		_, err = cmd.RunAndReturnTrimmedCombinedOutput()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m manager) SimulatorBoot(id string) error {
	f := command.NewFactory(env.NewRepository())
	cmd := f.Create("xcrun", []string{"simctl", "boot", id}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	log.Donef("$ %s", cmd.PrintableCommandArgs())
	exitCode, err := cmd.RunAndReturnExitCode()
	if err != nil {
		if errorutil.IsExitStatusError(err) {
			if exitCode == 149 { // Simulator already booted
				return nil
			}
			log.Warnf("Failed to boot Simulator, command exited with code %d", exitCode)
			return nil
		}
		return fmt.Errorf("failed to boot Simulator, command execution failed: %v", err)
	}

	return nil
}

// Simulator needs to be booted to enable verbose log
func (m manager) SimulatorEnableVerboseLog(id string) error {
	f := command.NewFactory(env.NewRepository())
	cmd := f.Create("xcrun", []string{"simctl", "logverbose", id, "enable"}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	log.Donef("$ %s", cmd.PrintableCommandArgs())
	if err := cmd.Run(); err != nil {
		if errorutil.IsExitStatusError(err) {
			log.Warnf("Failed to enable Simulator verbose logging, command exited with code %d", err)
			return nil
		}

		return fmt.Errorf("failed to enable Simulator verbose logging, command execution failed: %v", err)
	}

	return nil
}

func (m manager) SimulatorCollectDiagnostics() (string, error) {
	diagnosticsName, err := m.SimulatorDiagnosticsName()
	if err != nil {
		return "", err
	}
	diagnosticsOutDir, err := ioutil.TempDir("", diagnosticsName)
	if err != nil {
		return "", fmt.Errorf("failed to collect Simulator diagnostics, could not create temporary directory: %v", err)
	}

	f := command.NewFactory(env.NewRepository())
	cmd := f.Create("xcrun", []string{"simctl", "diagnose", "-b", "--no-archive", fmt.Sprintf("--output=%s", diagnosticsOutDir)}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  bytes.NewReader([]byte("\n")),
	})

	log.Donef("$ %s", cmd.PrintableCommandArgs())
	if err := cmd.Run(); err != nil {
		if errorutil.IsExitStatusError(err) {
			return "", fmt.Errorf("failed to collect Simulator diagnostics: %v", err)

		}
		return "", fmt.Errorf("failed to collect Simulator diagnostics, command execution failed: %v", err)
	}

	return diagnosticsOutDir, nil
}

func (m manager) SimulatorShutdown(id string) error {
	f := command.NewFactory(env.NewRepository())
	cmd := f.Create("xcrun", []string{"simctl", "shutdown", id}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	log.Donef("$ %s", cmd.PrintableCommandArgs())
	exitCode, err := cmd.RunAndReturnExitCode()
	if err != nil {
		if errorutil.IsExitStatusError(err) {
			if exitCode == 149 { // Simulator already shut down
				return nil
			}
			log.Warnf("Failed to shutdown Simulator, command exited with code %d", exitCode)
			return nil
		}
		return fmt.Errorf("failed to shutdown Simulator, command execution failed: %v", err)
	}

	return nil
}

func (m manager) SimulatorDiagnosticsName() (string, error) {
	timestamp, err := time.Now().MarshalText()
	if err != nil {
		return "", fmt.Errorf("failed to collect Simulator diagnostics, failed to marshal timestamp: %v", err)
	}

	return fmt.Sprintf("simctl_diagnose_%s.zip", strings.ReplaceAll(string(timestamp), ":", "-")), nil
}
