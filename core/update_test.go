package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateVariable(t *testing.T) {
	gh, err := NewghReleaseDownloader("MoreFind")
	require.Nil(t, err)
	_, err = gh.GetExecutableFromAsset()
}
