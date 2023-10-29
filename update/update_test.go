package update

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateVariable(t *testing.T) {
	_, err := NewghReleaseDownloader("morefind")
	require.Nil(t, err)
}
