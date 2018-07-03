package xcodeutil

import (
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/stretchr/testify/require"
)

func TestWalkXcodeTestSummaries(t *testing.T) {
	t.Log()
	{
		log, err := fileutil.ReadStringFromFile("../_samples/TestSummaries.plist")
		require.NoError(t, err)

		var testSummaries TestSummaries
		testSummaries.Content = log
		testSummaries, err = testSummaries.collectTestItemsWithScreenshotAndSetType()
		require.NoError(t, err)
		require.Equal(t, 2, len(testSummaries.TestItemsWithScreenshots))
	}

	t.Log()
	{
		log, err := fileutil.ReadStringFromFile("../_samples/TestSummaries2.plist")
		require.NoError(t, err)

		var testSummaries TestSummaries
		testSummaries.Content = log
		testSummaries, err = testSummaries.collectTestItemsWithScreenshotAndSetType()
		require.NoError(t, err)
		require.Equal(t, 2, len(testSummaries.TestItemsWithScreenshots))
	}
}

func TestTimestampToTime(t *testing.T) {
	time, err := TimestampStrToTime("522675441.31045401")
	require.NoError(t, err)

	require.Equal(t, 2017, time.Year())
	require.Equal(t, 7, int(time.Month()))
	require.Equal(t, 25, time.Day())
	require.Equal(t, 11, time.Hour())
	require.Equal(t, 37, time.Minute())
	require.Equal(t, 21, time.Second())
}
