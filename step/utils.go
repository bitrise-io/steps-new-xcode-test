package step

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/stringutil"
)

func saveRawOutputToLogFile(rawXcodebuildOutput string) (string, error) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("xcodebuild-output")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir, error: %s", err)
	}
	logFileName := "raw-xcodebuild-output.log"
	logPth := filepath.Join(tmpDir, logFileName)
	if err := fileutil.WriteStringToFile(logPth, rawXcodebuildOutput); err != nil {
		return "", fmt.Errorf("failed to write xcodebuild output to file, error: %s", err)
	}

	return logPth, nil
}

func saveAttachments(scheme, testSummariesPath, attachementDir string) (string, error) {
	if exist, err := pathutil.IsDirExists(attachementDir); err != nil {
		return "", err
	} else if !exist {
		return "", fmt.Errorf("no test attachments found at: %s", attachementDir)
	}

	if found, err := UpdateScreenshotNames(testSummariesPath, attachementDir); err != nil {
		log.Warnf("Failed to update screenshot names, error: %s", err)
	} else if !found {
		return "", nil
	}

	// deploy zipped attachments
	deployDir := os.Getenv("BITRISE_DEPLOY_DIR")
	if deployDir == "" {
		return "", errors.New("no BITRISE_DEPLOY_DIR found")
	}

	zipedTestsDerivedDataPath := filepath.Join(deployDir, fmt.Sprintf("%s-xc-test-Attachments.zip", scheme))
	if err := Zip(filepath.Dir(attachementDir), filepath.Base(attachementDir), zipedTestsDerivedDataPath); err != nil {
		return "", err
	}

	return zipedTestsDerivedDataPath, nil
}

func getSummariesAndAttachmentPath(testOutputDir string) (testSummariesPath string, attachmentDir string, err error) {
	const testSummaryFileName = "TestSummaries.plist"
	if exist, err := pathutil.IsDirExists(testOutputDir); err != nil {
		return "", "", err
	} else if !exist {
		return "", "", fmt.Errorf("no test logs found at: %s", testOutputDir)
	}

	testSummariesPath = path.Join(testOutputDir, testSummaryFileName)
	if exist, err := pathutil.IsPathExists(testSummariesPath); err != nil {
		return "", "", err
	} else if !exist {
		return "", "", fmt.Errorf("no test summaries found at: %s", testSummariesPath)
	}

	var attachementDir string
	{
		attachementDir = filepath.Join(testOutputDir, "Attachments")
		if exist, err := pathutil.IsDirExists(attachementDir); err != nil {
			return "", "", err
		} else if !exist {
			return "", "", fmt.Errorf("no test attachments found at: %s", attachementDir)
		}
	}

	log.Debugf("Test summaries path: %s", testSummariesPath)
	log.Debugf("Attachment dir: %s", attachementDir)
	return testSummariesPath, attachementDir, nil
}

func printLastLinesOfXcodebuildTestLog(rawXcodebuildOutput string, isRunSuccess bool) {
	const lastLines = "\nLast lines of the build log:"
	if !isRunSuccess {
		log.Errorf(lastLines)
	} else {
		log.Infof(lastLines)
	}

	fmt.Println(stringutil.LastNLines(rawXcodebuildOutput, 20))

	if !isRunSuccess {
		log.Warnf("If you can't find the reason of the error in the log, please check the xcodebuild_test.log.")
	}

	log.Infof(colorstring.Magenta(`
The log file is stored in $BITRISE_DEPLOY_DIR, and its full path
is available in the $BITRISE_XCODEBUILD_TEST_LOG_PATH environment variable.

If you have the Deploy to Bitrise.io step (after this step),
that will attach the file to your build as an artifact!`))
}

// Zip ...
func Zip(targetDir, targetRelPathToZip, zipPath string) error {
	zipCmd := exec.Command("/usr/bin/zip", "-rTy", zipPath, targetRelPathToZip)
	zipCmd.Dir = targetDir
	if out, err := zipCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Zip failed, out: %s, err: %#v", out, err)
	}
	return nil
}
