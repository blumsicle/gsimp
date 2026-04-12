package completion

import (
	"bytes"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCLI struct {
	Completion Command `cmd:""`
}

func TestRunWritesShellCompletionScript(t *testing.T) {
	tests := []struct {
		shell    string
		contains []string
	}{
		{
			shell: "zsh",
			contains: []string{
				"#compdef bcli",
				"compdef _bcli bcli",
				"_bcli() {",
			},
		},
		{
			shell: "bash",
			contains: []string{
				"_bcli_completions()",
				"complete -F _bcli_completions bcli",
			},
		},
		{
			shell: "fish",
			contains: []string{
				"# fish shell completion for bcli",
				"complete -c bcli -f",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			cli := &testCLI{}
			command := Command{Shell: tt.shell}
			var stdout bytes.Buffer
			parser, err := kong.New(
				cli,
				kong.Name("bcli"),
				kong.Writers(&stdout, &bytes.Buffer{}),
			)
			require.NoError(t, err)

			err = command.Run(&kong.Context{Kong: parser})
			require.NoError(t, err)

			for _, want := range tt.contains {
				assert.Contains(t, stdout.String(), want)
			}
		})
	}
}

func TestRunRejectsUnsupportedShell(t *testing.T) {
	command := Command{Shell: "tcsh"}

	err := command.Run(nil)

	assert.EqualError(t, err, `unsupported shell "tcsh"`)
}
