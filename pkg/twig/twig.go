package twig

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/creack/pty"
	"github.com/dwadp/twig/pkg/config"
	"github.com/dwadp/twig/pkg/schema"
	"golang.org/x/term"
)

var execCommand = exec.Command

func runPhpExec() ([]byte, error) {
	cmd := execCommand("php", "-v")
	return cmd.CombinedOutput()
}

func runInTTY(cmd *exec.Cmd) error {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	defer func() { _ = ptmx.Close() }()

	resizeSignalCh := make(chan os.Signal, 1)
	signal.Notify(resizeSignalCh, syscall.SIGWINCH)

	go func() {
		for range resizeSignalCh {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()

	resizeSignalCh <- syscall.SIGWINCH
	defer func() {
		signal.Stop(resizeSignalCh)
		close(resizeSignalCh)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		_ = cmd.Process.Signal(sig)
	}()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}

func isRootProjectDir(cwd string) bool {
	_, err := os.Stat(filepath.Join(cwd, "composer.json"))

	if err != nil && errors.Is(err, os.ErrNotExist) {
		vendorDir, err := os.Stat(filepath.Join(cwd, "vendor"))

		if err != nil && errors.Is(err, os.ErrNotExist) {
			return false
		}

		return vendorDir.IsDir()
	}

	return err == nil
}

func preparePHP(cfg *config.Config) (*config.PHP, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if isRootProjectDir(wd) {
		file, err := os.Open(filepath.Join(wd, "composer.json"))
		if err != nil {
			return nil, err
		}

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		versions, err := schema.GetVersion(b)
		if err != nil {
			return nil, err
		}

		php, err := cfg.GetPreferredPHPVersion(versions)
		if err != nil {
			return nil, err
		}

		return php, nil
	}

	return cfg.GetDefaultPHPVersion()
}

func RunPHP(cfg *config.Config, args []string) error {
	php, err := preparePHP(cfg)
	if err != nil {
		return err
	}

	cmd := exec.Command(php.ExecutablePath, args...)

	if err := runInTTY(cmd); err != nil {
		return err
	}

	return nil
}

func RunComposer(cfg *config.Config, args []string) error {
	php, err := preparePHP(cfg)
	if err != nil {
		return err
	}

	args = append([]string{cfg.Composer.ExecutablePath}, args...)

	cmd := exec.Command(php.ExecutablePath, args...)

	if err := runInTTY(cmd); err != nil {
		return err
	}

	return nil
}
