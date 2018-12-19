package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	cmd "github.com/bitrise-io/steps-xcode-test/command"
	"github.com/bitrise-io/steps-xcode-test/models"
	"github.com/bitrise-io/steps-xcode-test/pretty"
	"github.com/bitrise-io/steps-xcode-test/xcodeutil"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/progress"
	"github.com/bitrise-io/go-utils/stringutil"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/bitrise-tools/go-xcode/utility"
	shellquote "github.com/kballard/go-shellquote"
)

// On performance limited OS X hosts (ex: VMs) the iPhone/iOS Simulator might time out
//  while booting. So far it seems that a simple retry solves these issues.

const (
	minSupportedXcodeMajorVersion = 6
	// This boot timeout can happen when running Unit Tests with Xcode Command Line `xcodebuild`.
	timeOutMessageIPhoneSimulator = "iPhoneSimulator: Timed out waiting"
	// This boot timeout can happen when running Xcode (7+) UI tests with Xcode Command Line `xcodebuild`.
	timeOutMessageUITest                     = "Terminating app due to uncaught exception '_XCTestCaseInterruptionException'"
	earlyUnexpectedExit                      = "Early unexpected exit, operation never finished bootstrapping - no restart will be attempted"
	failureAttemptingToLaunch                = "Assertion Failure: <unknown>:0: UI Testing Failure - Failure attempting to launch <XCUIApplicationImpl:"
	failedToBackgroundTestRunner             = `Error Domain=IDETestOperationsObserverErrorDomain Code=12 "Failed to background test runner.`
	appStateIsStillNotRunning                = `App state is still not running active, state = XCApplicationStateNotRunning`
	appAccessibilityIsNotLoaded              = `UI Testing Failure - App accessibility isn't loaded`
	testRunnerFailedToInitializeForUITesting = `Test runner failed to initialize for UI testing`
	timedOutRegisteringForTestingEvent       = `Timed out registering for testing event accessibility notifications`
)

var automaticRetryReasonPatterns = []string{
	timeOutMessageIPhoneSimulator,
	timeOutMessageUITest,
	earlyUnexpectedExit,
	failureAttemptingToLaunch,
	failedToBackgroundTestRunner,
	appStateIsStillNotRunning,
	appAccessibilityIsNotLoaded,
	testRunnerFailedToInitializeForUITesting,
	timedOutRegisteringForTestingEvent,
}

var xcodeCommandEnvs = []string{"NSUnbufferedIO=YES"}

// -----------------------
// --- Models
// -----------------------

// Configs ...
type Configs struct {
	// Project Parameters
	ProjectPath string `env:"project_path,required"`
	Scheme      string `env:"scheme,required"`

	// Simulator Configs
	SimulatorPlatform  string `env:"simulator_platform,required"`
	SimulatorDevice    string `env:"simulator_device,required"`
	SimulatorOsVersion string `env:"simulator_os_version,required"`

	// Test Run Configs
	OutputTool    string `env:"output_tool,opt[xcpretty,xcodebuild]"`
	IsCleanBuild  bool   `env:"is_clean_build,opt[yes,no]"`
	IsSingleBuild bool   `env:"single_build,opt[true,false]"`

	ShouldBuildBeforeTest bool `env:"should_build_before_test,opt[yes,no]"`
	ShouldRetryTestOnFail bool `env:"should_retry_test_on_fail,opt[yes,no]"`

	GenerateCodeCoverageFiles bool `env:"generate_code_coverage_files,opt[yes,no]"`
	ExportUITestArtifacts     bool `env:"export_uitest_artifacts,opt[true,false]"`

	// Not required parameters
	TestOptions         string `env:"xcodebuild_test_options"`
	XcprettyTestOptions string `env:"xcpretty_test_options"`

	// Debug
	Verbose      bool `env:"verbose,opt[yes,no]"`
	HeadlessMode bool `env:"headless_mode,opt[yes,no]"`
}

const testSummaryFileName = "TestSummaries.plist"

func isStringFoundInOutput(searchStr, outputToSearchIn string) bool {
	r, err := regexp.Compile("(?i)" + searchStr)
	if err != nil {
		log.Warnf("Failed to compile regexp: %s", err)
		return false
	}
	return r.MatchString(outputToSearchIn)
}

func runXcodeBuildCmd(useStdOut bool, args ...string) (string, int, error) {
	// command
	buildCmd := cmd.CreateXcodebuildCmd(args...)
	// output buffer
	var outBuffer bytes.Buffer
	// additional output writers, like StdOut
	outWritters := []io.Writer{}
	if useStdOut {
		outWritters = append(outWritters, os.Stdout)
	}
	// unify as a single writer
	outWritter := cmd.CreateBufferedWriter(&outBuffer, outWritters...)
	// and set the writer
	buildCmd.Stdin = nil
	buildCmd.Stdout = outWritter
	buildCmd.Stderr = outWritter
	buildCmd.Env = append(os.Environ(), xcodeCommandEnvs...)

	cmdArgsForPrint := cmd.PrintableCommandArgsWithEnvs(buildCmd.Args, xcodeCommandEnvs)

	log.Printf("$ %s", cmdArgsForPrint)

	err := buildCmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus, ok := exitError.Sys().(syscall.WaitStatus)
			if !ok {
				return outBuffer.String(), 1, errors.New("Failed to cast exit status")
			}
			return outBuffer.String(), waitStatus.ExitStatus(), err
		}
		return outBuffer.String(), 1, err
	}
	return outBuffer.String(), 0, nil
}

func runPrettyXcodeBuildCmd(useStdOut bool, xcprettyArgs []string, xcodebuildArgs []string) (string, int, error) {
	//
	buildCmd := cmd.CreateXcodebuildCmd(xcodebuildArgs...)
	prettyCmd := cmd.CreateXcprettyCmd(xcprettyArgs...)
	//
	var buildOutBuffer bytes.Buffer
	//
	pipeReader, pipeWriter := io.Pipe()
	//
	// build outputs:
	// - write it into a buffer
	// - write it into the pipe, which will be fed into xcpretty
	buildOutWriters := []io.Writer{pipeWriter}
	buildOutWriter := cmd.CreateBufferedWriter(&buildOutBuffer, buildOutWriters...)
	//
	var prettyOutWriter io.Writer
	if useStdOut {
		prettyOutWriter = os.Stdout
	}

	// and set the writers
	buildCmd.Stdin = nil
	buildCmd.Stdout = buildOutWriter
	buildCmd.Stderr = buildOutWriter
	//
	prettyCmd.Stdin = pipeReader
	prettyCmd.Stdout = prettyOutWriter
	prettyCmd.Stderr = prettyOutWriter
	//
	buildCmd.Env = append(os.Environ(), xcodeCommandEnvs...)

	log.Printf("$ set -o pipefail && %s | %v",
		cmd.PrintableCommandArgsWithEnvs(buildCmd.Args, xcodeCommandEnvs),
		cmd.PrintableCommandArgs(prettyCmd.Args))

	fmt.Println()

	if err := buildCmd.Start(); err != nil {
		return buildOutBuffer.String(), 1, err
	}
	if err := prettyCmd.Start(); err != nil {
		return buildOutBuffer.String(), 1, err
	}

	defer func() {
		if err := pipeWriter.Close(); err != nil {
			log.Warnf("Failed to close xcodebuild-xcpretty pipe, error: %s", err)
		}

		if err := prettyCmd.Wait(); err != nil {
			log.Warnf("xcpretty command failed, error: %s", err)
		}
	}()

	if err := buildCmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus, ok := exitError.Sys().(syscall.WaitStatus)
			if !ok {
				return buildOutBuffer.String(), 1, errors.New("Failed to cast exit status")
			}
			return buildOutBuffer.String(), waitStatus.ExitStatus(), err
		}
		return buildOutBuffer.String(), 1, err
	}

	return buildOutBuffer.String(), 0, nil
}

func runBuild(buildParams models.XcodeBuildParamsModel, outputTool string) (string, int, error) {
	xcodebuildArgs := []string{buildParams.Action, buildParams.ProjectPath, "-scheme", buildParams.Scheme}
	if buildParams.CleanBuild {
		xcodebuildArgs = append(xcodebuildArgs, "clean")
	}
	xcodebuildArgs = append(xcodebuildArgs, "build", "-destination", buildParams.DeviceDestination)

	log.Infof("Building the project...")

	if outputTool == "xcpretty" {
		return runPrettyXcodeBuildCmd(false, []string{}, xcodebuildArgs)
	}
	return runXcodeBuildCmd(false, xcodebuildArgs...)
}

func runTest(buildTestParams models.XcodeBuildTestParamsModel, outputTool, xcprettyOptions string, isAutomaticRetryOnReason, isRetryOnFail bool) (string, int, error) {
	handleTestError := func(fullOutputStr string, exitCode int, testError error) (string, int, error) {
		//
		// Automatic retry
		for _, retryReasonPattern := range automaticRetryReasonPatterns {
			if isStringFoundInOutput(retryReasonPattern, fullOutputStr) {
				log.Warnf("Automatic retry reason found in log: %s", retryReasonPattern)
				if isAutomaticRetryOnReason {
					log.Printf("isAutomaticRetryOnReason=true - retrying...")
					return runTest(buildTestParams, outputTool, xcprettyOptions, false, false)
				}
				log.Errorf("isAutomaticRetryOnReason=false, no more retry, stopping the test!")
				return fullOutputStr, exitCode, testError
			}
		}

		//
		// Retry on fail
		if isRetryOnFail {
			log.Warnf("Test run failed")
			log.Printf("isRetryOnFail=true - retrying...")
			return runTest(buildTestParams, outputTool, xcprettyOptions, false, false)
		}

		return fullOutputStr, exitCode, testError
	}

	buildParams := buildTestParams.BuildParams

	xcodebuildArgs := []string{buildParams.Action, buildParams.ProjectPath, "-scheme", buildParams.Scheme}
	if buildTestParams.CleanBuild {
		xcodebuildArgs = append(xcodebuildArgs, "clean")
	}
	// the 'build' argument is required *before* the 'test' arg, to prevent
	//  the Xcode bug described in the README, which causes:
	// 'iPhoneSimulator: Timed out waiting 120 seconds for simulator to boot, current state is 1.'
	//  in case the compilation takes a long time.
	// Related Radar link: https://openradar.appspot.com/22413115
	// Demonstration project: https://github.com/bitrise-io/simulator-launch-timeout-includes-build-time

	// for builds < 120 seconds or fixed Xcode versions, one should
	// have the possibility of opting out, because the explicit build arg
	// leads the project to be compiled twice and increase the duration
	// Related issue link: https://github.com/bitrise-io/steps-xcode-test/issues/55
	if buildTestParams.BuildBeforeTest {
		xcodebuildArgs = append(xcodebuildArgs, "build")
	}
	xcodebuildArgs = append(xcodebuildArgs, "test", "-destination", buildParams.DeviceDestination)
	xcodebuildArgs = append(xcodebuildArgs, "-resultBundlePath", buildTestParams.TestOutputDir)

	if buildTestParams.GenerateCodeCoverage {
		xcodebuildArgs = append(xcodebuildArgs, "GCC_INSTRUMENT_PROGRAM_FLOW_ARCS=YES")
		xcodebuildArgs = append(xcodebuildArgs, "GCC_GENERATE_TEST_COVERAGE_FILES=YES")
	}

	if buildTestParams.AdditionalOptions != "" {
		options, err := shellquote.Split(buildTestParams.AdditionalOptions)
		if err != nil {
			return "", 1, fmt.Errorf("failed to parse additional options (%s), error: %s", buildTestParams.AdditionalOptions, err)
		}
		xcodebuildArgs = append(xcodebuildArgs, options...)
	}

	xcprettyArgs := []string{}
	if xcprettyOptions != "" {
		options, err := shellquote.Split(xcprettyOptions)
		if err != nil {
			return "", 1, fmt.Errorf("failed to parse additional options (%s), error: %s", xcprettyOptions, err)
		}
		// get and delete the xcpretty output file, if exists
		xcprettyOutputFilePath := ""
		isNextOptOutputPth := false
		for _, aOpt := range options {
			if isNextOptOutputPth {
				xcprettyOutputFilePath = aOpt
				break
			}
			if aOpt == "--output" {
				isNextOptOutputPth = true
				continue
			}
		}
		if xcprettyOutputFilePath != "" {
			if isExist, err := pathutil.IsPathExists(xcprettyOutputFilePath); err != nil {
				log.Errorf("Failed to check xcpretty output file status (path: %s), error: %s", xcprettyOutputFilePath, err)
			} else if isExist {
				log.Warnf("=> Deleting existing xcpretty output: %s", xcprettyOutputFilePath)
				if err := os.Remove(xcprettyOutputFilePath); err != nil {
					log.Errorf("Failed to delete xcpretty output file (path: %s), error: %s", xcprettyOutputFilePath, err)
				}
			}
		}
		//
		xcprettyArgs = append(xcprettyArgs, options...)
	}

	log.Infof("Running the tests...")

	var rawOutput string
	var err error
	var exit int
	if outputTool == "xcpretty" {
		rawOutput, exit, err = runPrettyXcodeBuildCmd(true, xcprettyArgs, xcodebuildArgs)
	} else {
		rawOutput, exit, err = runXcodeBuildCmd(true, xcodebuildArgs...)
	}

	if err != nil {
		return handleTestError(rawOutput, exit, err)
	}
	return rawOutput, exit, nil
}

func saveRawOutputToLogFile(rawXcodebuildOutput string, isRunSuccess bool) (string, error) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("xcodebuild-output")
	if err != nil {
		return "", fmt.Errorf("Failed to create temp dir, error: %s", err)
	}
	logFileName := "raw-xcodebuild-output.log"
	logPth := filepath.Join(tmpDir, logFileName)
	if err := fileutil.WriteStringToFile(logPth, rawXcodebuildOutput); err != nil {
		return "", fmt.Errorf("Failed to write xcodebuild output to file, error: %s", err)
	}

	if !isRunSuccess {
		deployDir := os.Getenv("BITRISE_DEPLOY_DIR")
		if deployDir == "" {
			return "", errors.New("No BITRISE_DEPLOY_DIR found")
		}
		deployPth := filepath.Join(deployDir, logFileName)

		if err := command.CopyFile(logPth, deployPth); err != nil {
			return "", fmt.Errorf("Failed to copy xcodebuild output log file from (%s) to (%s), error: %s", logPth, deployPth, err)
		}
		logPth = deployPth
	}

	if err := cmd.ExportEnvironmentWithEnvman("BITRISE_XCODE_RAW_TEST_RESULT_TEXT_PATH", logPth); err != nil {
		log.Warnf("Failed to export: BITRISE_XCODE_RAW_TEST_RESULT_TEXT_PATH, error: %s", err)
	}
	return logPth, nil
}

func screenshotName(startTime time.Time, title, uuid string) string {
	formattedDate := startTime.Format("2006-01-02_03-04-05")
	fixedTitle := strings.Replace(title, " ", "_", -1)
	return fmt.Sprintf("%s_%s_%s", formattedDate, fixedTitle, uuid)
}

func updateScreenshotNames(testLogsDir string) (bool, error) {
	testSummariesPath := filepath.Join(testLogsDir, testSummaryFileName)
	if exist, err := pathutil.IsPathExists(testSummariesPath); err != nil {
		return false, fmt.Errorf("Failed to check if file exists: %s", testSummariesPath)
	} else if !exist {
		return false, fmt.Errorf("no TestSummaries file found: %s", testSummariesPath)
	}

	//
	// TestSummaries
	testSummaries, err := xcodeutil.NewTestSummaries(testSummariesPath)
	if err != nil {
		return false, fmt.Errorf("failed to parse %s, error: %s", filepath.Base(testSummariesPath), err)
	}

	log.Debugf("Test items with screenshots: %s", pretty.Object(testSummaries.TestItemsWithScreenshots))
	log.Debugf("TestSummaries version has been set to: %s\n", testSummaries.Type)

	if len(testSummaries.TestItemsWithScreenshots) > 0 {
		log.Printf("Renaming screenshots")
	} else {
		log.Printf("No screenshot found")
		return false, nil
	}

	for _, testItem := range testSummaries.TestItemsWithScreenshots {
		startTimeIntervalObj, found := testItem["StartTimeInterval"]
		if !found {
			return false, fmt.Errorf("missing StartTimeInterval")
		}
		startTimeInterval, casted := startTimeIntervalObj.(float64)
		if !casted {
			return false, fmt.Errorf("StartTimeInterval is not a float64")
		}
		startTime, err := xcodeutil.TimestampToTime(startTimeInterval)
		if err != nil {
			return false, err
		}

		uuidObj, found := testItem["UUID"]
		if !found {
			return false, fmt.Errorf("missing UUID")
		}
		uuid, casted := uuidObj.(string)
		if !casted {
			return false, fmt.Errorf("UUID is not a string")
		}

		// Renaming the screenshots
		{
			var err error
			var origScreenshotPth string
			if testSummaries.Type == xcodeutil.TestSummariesWithScreenshotData { // TestSummariesWithScreenshotData - TestSummaries.plist
				origScreenshotPth, err = updateOldSummaryTypeScreenshotName(testItem, testLogsDir, uuid, startTime)
			} else {
				origScreenshotPth, err = updateNewSummaryTypeScreenshotName(testItem, testLogsDir, uuid, startTime)
			}
			if err != nil {
				log.Warnf("Failed to rename the screenshot: %s - err: %s", filepath.Base(origScreenshotPth), err)
			}
		}
	}

	return true, nil
}

func updateOldSummaryTypeScreenshotName(testItem map[string]interface{}, testLogsDir, uuid string, startTime time.Time) (string, error) {
	var origScreenshotPth string

	for _, ext := range []string{"png", "jpg"} {
		origScreenshotPth = filepath.Join(testLogsDir, "Attachments", fmt.Sprintf("Screenshot_%s.%s", uuid, ext))
		var newScreenshotPth string

		if exist, err := pathutil.IsPathExists(origScreenshotPth); err != nil {
			return "", err
		} else if exist {
			titleObj, found := testItem["Title"]
			if !found {
				return origScreenshotPth, fmt.Errorf("missing Title")
			}
			title, casted := titleObj.(string)
			if !casted {
				return origScreenshotPth, fmt.Errorf("Title is not a string")
			}

			newScreenshotPth = filepath.Join(testLogsDir, "Attachments", screenshotName(startTime, title, uuid)+"."+ext)
			if err := os.Rename(origScreenshotPth, newScreenshotPth); err != nil {
				return origScreenshotPth, err
			}
			log.Printf("%s => %s", filepath.Base(origScreenshotPth), filepath.Base(newScreenshotPth))
		}
	}
	return origScreenshotPth, nil
}

func updateNewSummaryTypeScreenshotName(testItem map[string]interface{}, testLogsDir, uuid string, startTime time.Time) (string, error) {
	var origScreenshotPth string

	attachmentsObj, found := testItem["Attachments"]
	if !found {
		return "", fmt.Errorf("Attachments not found in the TestSummaries.plist")
	}

	attachments, casted := attachmentsObj.([]interface{})
	if !casted {
		return "", fmt.Errorf("Failed to cast attachmentsObj")
	}

	var fileName string
	for _, attachmentObj := range attachments {
		attachment, casted := attachmentObj.(map[string]interface{})
		if !casted {
			return "", fmt.Errorf("Failed to cast attachmentObj")
		}

		fileNameObj, found := attachment["Filename"]
		if found {
			fileName, casted = fileNameObj.(string)
			if casted {
				origScreenshotPth = filepath.Join(testLogsDir, "Attachments", fileName)
			}
		}

		if exist, err := pathutil.IsPathExists(origScreenshotPth); err != nil {
			return "", err
		} else if exist {
			formattedDate := startTime.Format("2006-01-02_03-04-05")
			newScreenshotPth := filepath.Join(testLogsDir, "Attachments", (formattedDate + "_" + fileName))
			if err := os.Rename(origScreenshotPth, newScreenshotPth); err != nil {
				log.Warnf("Failed to rename the screenshot: %s", filepath.Base(origScreenshotPth))
				continue
			}
			log.Printf("Screenshot renamed: %s => %s", filepath.Base(origScreenshotPth), filepath.Base(newScreenshotPth))
		}
	}
	return origScreenshotPth, nil
}

func saveAttachments(scheme, testDir, attachementDir string) error {

	if exist, err := pathutil.IsDirExists(attachementDir); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("no test attachments found at: %s", attachementDir)
	}

	// update screenshot name:
	// Screenshot_uuid.png -> start_date_time_title_uuid.png
	// Screenshot_uuid.jpg -> start_date_time_title_uuid.jpg
	var found bool
	var err error
	if found, err = updateScreenshotNames(testDir); err != nil {
		log.Warnf("Failed to update screenshot names, error: %s", err)
	}

	if !found {
		return nil
	}

	// deploy zipped attachments
	deployDir := os.Getenv("BITRISE_DEPLOY_DIR")
	if deployDir == "" {
		return errors.New("No BITRISE_DEPLOY_DIR found")
	}

	zipedTestsDerivedDataPath := filepath.Join(deployDir, fmt.Sprintf("%s-xc-test-Attachments.zip", scheme))
	if err := cmd.Zip(testDir, "Attachments", zipedTestsDerivedDataPath); err != nil {
		return err
	}

	if err := cmd.ExportEnvironmentWithEnvman("BITRISE_XCODE_TEST_ATTACHMENTS_PATH", zipedTestsDerivedDataPath); err != nil {
		log.Warnf("Failed to export: BITRISE_XCODE_TEST_ATTACHMENTS_PATH, error: %s", err)
	}

	log.Donef("The zipped attachments are available in: %s", zipedTestsDerivedDataPath)
	return nil
}

func getAttachmentDir(testOutputDir string) (string, error) {
	if exist, err := pathutil.IsDirExists(testOutputDir); err != nil {
		return "", err
	} else if !exist {
		return "", fmt.Errorf("no test logs found at: %s", testOutputDir)
	}

	if exist, err := pathutil.IsPathExists(path.Join(testOutputDir, testSummaryFileName)); err != nil {
		return "", err
	} else if !exist {
		return "", fmt.Errorf("no %s found at: %s", testSummaryFileName, testOutputDir)
	}

	var attachementDir string
	{
		attachementDir = filepath.Join(testOutputDir, "Attachments")
		if exist, err := pathutil.IsDirExists(attachementDir); err != nil {
			return "", err
		} else if !exist {
			return "", fmt.Errorf("no test attachments found at: %s", attachementDir)
		}
	}

	log.Debugf("Test output dir: %s", testOutputDir)
	log.Debugf("Attachment dir: %s", attachementDir)
	return attachementDir, nil
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

//--------------------
// Main
//--------------------

func main() {
	var configs Configs
	if err := stepconf.Parse(&configs); err != nil {
		fail("Issue with input: %s", err)
	}

	stepconf.Print(configs)
	fmt.Println()
	log.SetEnableDebugLog(configs.Verbose)

	// Project-or-Workspace flag
	action := ""
	if strings.HasSuffix(configs.ProjectPath, ".xcodeproj") {
		action = "-project"
	} else if strings.HasSuffix(configs.ProjectPath, ".xcworkspace") {
		action = "-workspace"
	} else {
		if err := cmd.ExportEnvironmentWithEnvman("BITRISE_XCODE_TEST_RESULT", "failed"); err != nil {
			log.Warnf("Failed to export: BITRISE_XCODE_TEST_RESULT, error: %s", err)
			fmt.Println()
		}
		fail("Invalid project file (%s), extension should be (.xcodeproj/.xcworkspace)", configs.ProjectPath)
	}

	log.Printf("* action: %s", action)

	// Detect Xcode major version
	xcodebuildVersion, err := utility.GetXcodeVersion()
	if err != nil {
		fail("Failed to determine xcode version, error: %s", err)
	}
	log.Printf("- xcodebuildVersion: %s (%s)", xcodebuildVersion.Version, xcodebuildVersion.BuildVersion)

	if xcodebuildVersion.MajorVersion < 9 && configs.HeadlessMode {
		log.Warnf("Headless mode is enabled but it's only availabe with Xcode 9.x or newer.")
	}

	xcodeMajorVersion := xcodebuildVersion.MajorVersion
	if xcodeMajorVersion < minSupportedXcodeMajorVersion {
		fail("Invalid xcode major version (%d), should not be less then min supported: %d", xcodeMajorVersion, minSupportedXcodeMajorVersion)
	}

	// Detect xcpretty version
	outputTool := configs.OutputTool
	xcprettyVersion, err := InstallXcpretty()
	if err != nil {
		log.Warnf("%s", err)
		log.Printf("Switching to xcodebuild for output tool")
		outputTool = "xcodebuild"
	} else {
		log.Printf("- xcprettyVersion: %s", xcprettyVersion.String())
		fmt.Println()
	}

	// Simulator infos
	simulator, err := xcodeutil.GetSimulator(configs.SimulatorPlatform, configs.SimulatorDevice, configs.SimulatorOsVersion)
	if err != nil {
		if err := cmd.ExportEnvironmentWithEnvman("BITRISE_XCODE_TEST_RESULT", "failed"); err != nil {
			log.Warnf("Failed to export: BITRISE_XCODE_TEST_RESULT, error: %s", err)
		}
		fail("failed to get simulator udid, error: ", err)
	}

	log.Infof("Simulator infos")
	log.Printf("* simulator_name: %s, UDID: %s, status: %s", simulator.Name, simulator.SimID, simulator.Status)

	// Device Destination
	deviceDestination := fmt.Sprintf("id=%s", simulator.SimID)

	log.Printf("* device_destination: %s", deviceDestination)
	fmt.Println()

	// Create temporary directory for test outputs
	var testOutputDir string
	{
		tempDir, err := ioutil.TempDir("", "XCUITestOutput")
		if err != nil {
			fail("Could not create test output temporary directory.")
		}
		// Leaving the output dir in place after exiting
		testOutputDir = path.Join(tempDir, "Test.xcresult")
	}

	buildParams := models.XcodeBuildParamsModel{
		Action:            action,
		ProjectPath:       configs.ProjectPath,
		Scheme:            configs.Scheme,
		DeviceDestination: deviceDestination,
		CleanBuild:        configs.IsCleanBuild,
	}

	buildTestParams := models.XcodeBuildTestParamsModel{
		BuildParams: buildParams,

		TestOutputDir:        testOutputDir,
		BuildBeforeTest:      configs.ShouldBuildBeforeTest,
		AdditionalOptions:    configs.TestOptions,
		GenerateCodeCoverage: configs.GenerateCodeCoverageFiles,
	}

	if configs.IsSingleBuild {
		buildTestParams.CleanBuild = configs.IsCleanBuild
	}

	//
	// If headless mode disabled - Start simulator
	if simulator.Status == "Shutdown" && !configs.HeadlessMode {
		log.Infof("Booting simulator (%s)...", simulator.SimID)

		if err := xcodeutil.BootSimulator(simulator, xcodebuildVersion); err != nil {
			if err := cmd.ExportEnvironmentWithEnvman("BITRISE_XCODE_TEST_RESULT", "failed"); err != nil {
				log.Warnf("Failed to export: BITRISE_XCODE_TEST_RESULT, error: %s", err)
			}
			fail("failed to boot simulator, error: ", err)
		}

		progress.NewDefaultWrapper("Waiting for simulator boot").WrapAction(func() {
			time.Sleep(60 * time.Second)
		})

		fmt.Println()
	}

	//
	// Run build
	if !configs.IsSingleBuild {
		if rawXcodebuildOutput, exitCode, buildErr := runBuild(buildParams, outputTool); buildErr != nil {
			if _, err := saveRawOutputToLogFile(rawXcodebuildOutput, false); err != nil {
				log.Warnf("Failed to save the Raw Output, err: %s", err)
			}

			log.Warnf("xcode build exit code: %d", exitCode)
			log.Warnf("xcode build log:\n%s", rawXcodebuildOutput)
			log.Errorf("xcode build failed with error: %s", buildErr)
			if err := cmd.ExportEnvironmentWithEnvman("BITRISE_XCODE_TEST_RESULT", "failed"); err != nil {
				log.Warnf("Failed to export: BITRISE_XCODE_TEST_RESULT, error: %s", err)
			}
			os.Exit(1)
		}
	}

	//
	// Run test
	rawXcodebuildOutput, exitCode, testErr := runTest(buildTestParams, outputTool, configs.XcprettyTestOptions, true, configs.ShouldRetryTestOnFail)

	logPth, err := saveRawOutputToLogFile(rawXcodebuildOutput, (testErr == nil))

	if err != nil {
		log.Warnf("Failed to save the Raw Output, error %s", err)
	}

	if configs.ExportUITestArtifacts {
		fmt.Println()
		log.Infof("Exporting attachments")

		attachementDir, err := getAttachmentDir(buildTestParams.TestOutputDir)
		if err != nil {
			log.Warnf("Failed to export UI test artifacts, error %s", err)
		}

		if err := saveAttachments(configs.Scheme, buildTestParams.TestOutputDir, attachementDir); err != nil {
			log.Warnf("Failed to export UI test artifacts, error %s", err)
		}
	}

	if testErr != nil {
		log.Warnf("xcode test exit code: %d", exitCode)
		log.Errorf("xcode test failed, error: %s", testErr)
		log.Errorf("\nLast lines of the Xcode's build log:")
		fmt.Println(stringutil.LastNLines(rawXcodebuildOutput, 10))
		log.Warnf(`If you can't find the reason of the error in the log, please check the raw-xcodebuild-output.log
The log file is stored in $BITRISE_DEPLOY_DIR, and its full path
is available in the $BITRISE_XCODE_RAW_TEST_RESULT_TEXT_PATH environment variable.

You can check the full, unfiltered and unformatted Xcode output in the file:
%s
If you have the Deploy to Bitrise.io step (after this step),
that will attach the file to your build as an artifact!`, logPth)
		if err := cmd.ExportEnvironmentWithEnvman("BITRISE_XCODE_TEST_RESULT", "failed"); err != nil {
			log.Warnf("Failed to export: BITRISE_XCODE_TEST_RESULT, error: %s", err)
		}
		os.Exit(1)
	}

	if err := cmd.ExportEnvironmentWithEnvman("BITRISE_XCODE_TEST_RESULT", "succeeded"); err != nil {
		log.Warnf("Failed to export: BITRISE_XCODE_TEST_RESULT, error: %s", err)
	}
}
