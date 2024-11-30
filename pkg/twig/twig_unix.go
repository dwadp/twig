//go:build !windows
// +build !windows

package twig

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

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
