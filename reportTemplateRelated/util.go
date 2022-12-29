package reportTemplateRelated

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"runtime"
)

func RunCommand(command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	}
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return "", errors.New(stderr.String())
	}
	return out.String(), nil
}
