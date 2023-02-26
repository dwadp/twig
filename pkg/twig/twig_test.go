package twig

import (
	"fmt"
	"github.com/dwadp/twig/pkg/config"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
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

	defer func(t *testing.T) {
		assert.NoError(t, config.CleanupTestFiles("test-config"))
	}(t)

	RunPHP(cfg, []string{"artisan", "route:list"})

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()
	out, err := runPhpExec("docker/whalesay")
	if err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}
	if string(out) != "something" {
		t.Errorf("Expected %q, got %q", "someting", out)
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
