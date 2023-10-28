package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateVariable(t *testing.T) {
	HideProgressBar = false
	gh, err := NewghReleaseDownloader("MoreFind")
	require.Nil(t, err)
	_, err = gh.GetExecutableFromAsset()
	require.Nil(t, err)
}
