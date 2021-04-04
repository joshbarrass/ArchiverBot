package uarchiver

import (
	"fmt"
	"io"
	"os/exec"
)

type StdPipe io.ReadCloser

type WaitFunc func() error

// DownloadAuto calls uarchiver and lets it run to completion without
// any access to stdin or stdout
func DownloadAuto(url string) error {
	cmd := exec.Command("uarchiver", "-a", url)
	err := cmd.Run()
	return err
}

// DownloadAutoOutput starts uarchiver and returns access to stdout
// and stderr channels
func DownloadAutoOutput(url string) (WaitFunc, StdPipe, StdPipe, error) {
	cmd := exec.Command("uarchiver", "-a", url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to bind stdout: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to bind stderr: %w", err)
	}
	err = cmd.Start()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to start to command: %w", err)
	}
	return cmd.Wait, stdout, stderr, nil
}
