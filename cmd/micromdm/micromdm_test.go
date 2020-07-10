package main

import (
	"bytes"
	"strings"
	"testing"
)

type mainIOFunc func(t *testing.T, stdin, stdout, stderr *bytes.Buffer)

func mainUsage(t *testing.T, stdin, stdout, stderr *bytes.Buffer) {
	output := stderr.String()
	words := []string{"USAGE", "pidfile"}
	for _, word := range words {
		if !strings.Contains(output, word) {
			t.Errorf("expected %q in output, got:\n%s", word, output)
		}
	}
}

func TestMain(t *testing.T) {
	tests := []struct {
		name   string
		stdin  *bytes.Buffer
		stdout *bytes.Buffer
		stderr *bytes.Buffer
		args   []string
		exit   int
		check  mainIOFunc
	}{
		{
			name:   "short usage",
			stdin:  new(bytes.Buffer),
			stdout: new(bytes.Buffer),
			stderr: new(bytes.Buffer),
			args:   []string{"micromdm", "-h"},
			exit:   2,
			check:  mainUsage,
		},
		{
			name:   "long usage",
			stdin:  new(bytes.Buffer),
			stdout: new(bytes.Buffer),
			stderr: new(bytes.Buffer),
			args:   []string{"micromdm", "-help"},
			exit:   2,
			check:  mainUsage,
		},
		{
			name:   "subcommand help usage",
			stdin:  new(bytes.Buffer),
			stdout: new(bytes.Buffer),
			stderr: new(bytes.Buffer),
			args:   []string{"micromdm", "help"},
			exit:   2,
			check:  mainUsage,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got, want := micromdm(tt.args, tt.stdin, tt.stdout, tt.stderr), tt.exit; got != want {
				t.Fatalf("exit code: got %d, want %d", got, want)
			}

			tt.check(t, tt.stdin, tt.stdout, tt.stderr)
		})
	}
}
