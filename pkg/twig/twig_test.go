package twig

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dwadp/twig/pkg/config"
	"github.com/stretchr/testify/assert"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestRunPHPCommand(t *testing.T) {
	cfg, err := config.NewConfig("test-config", config.WithStub(`php:
  - version: 7.2
    executable_path: /your/path/to/php7.2
    default: true
  - version: 7.4
    executable_path: /your/path/to/php7.4
  - version: 8.0
    executable_path: /your/path/to/php8.0
  - version: 8.1
    executable_path: /your/path/to/php8.1
composer:
  executable_path: /your/composer.phar/path`))

	assert.NoError(t, err)
	assert.NoError(t, cfg.Init())
	assert.NoError(t, cfg.Read())

	RunPHP(cfg, []string{"artisan", "route:list"})

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()
	out, err := runPhpExec()
	if err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}
	if string(out) != "something" {
		t.Errorf("Expected %q, got %q", "someting", out)
	}

	t.Cleanup(func() {
		os.RemoveAll(cfg.BasePath())
	})
}

func TestIsRootProjectDir(t *testing.T) {
	tempDir := t.TempDir()

	assert.False(t, isRootProjectDir(tempDir))

	tests := []struct {
		name   string
		object string
		isDir  bool
		want   bool
	}{
		{
			name:   "determine by composer.json",
			object: "composer.json",
			want:   true,
		},
		{
			name:   "determine by vendor directory",
			object: "vendor",
			isDir:  true,
			want:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.isDir {
				f, err := os.OpenFile(filepath.Join(tempDir, test.object), os.O_CREATE, 0644)
				assert.NoError(t, err)
				defer f.Close()
			} else {
				err := os.Mkdir(filepath.Join(tempDir, test.object), 0755)
				assert.NoError(t, err)
			}

			assert.True(t, isRootProjectDir(tempDir))

			t.Cleanup(func() {
				os.RemoveAll(filepath.Join(tempDir, test.object))
			})
		})
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, "something")
	os.Exit(0)
}
