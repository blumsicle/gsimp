package bcliconfig

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeExpandsRootPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	currentUser, err := user.Current()
	require.NoError(t, err)

	t.Setenv("BCLI_TEST_ROOT", "/tmp/bcli-root")
	t.Setenv("BCLI_TEST_RELATIVE_ROOT", "~/from-env")

	tests := []struct {
		name string
		root string
		want string
	}{
		{
			name: "environment variable",
			root: "$BCLI_TEST_ROOT/src",
			want: "/tmp/bcli-root/src",
		},
		{
			name: "tilde",
			root: "~/src",
			want: filepath.Join(homeDir, "src"),
		},
		{
			name: "bare tilde",
			root: "~",
			want: homeDir,
		},
		{
			name: "environment variable to tilde",
			root: "$BCLI_TEST_RELATIVE_ROOT/project",
			want: filepath.Join(homeDir, "from-env", "project"),
		},
		{
			name: "non-leading tilde",
			root: "/tmp/~/src",
			want: "/tmp/~/src",
		},
		{
			name: "named user tilde",
			root: "~" + currentUser.Username + "/src",
			want: filepath.Join(homeDir, "src"),
		},
		{
			name: "unknown named user tilde",
			root: "~bcli-user-that-should-not-exist/src",
			want: "~bcli-user-that-should-not-exist/src",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{RootPath: tt.root}

			cfg.Normalize()

			assert.Equal(t, tt.want, cfg.RootPath)
		})
	}
}
