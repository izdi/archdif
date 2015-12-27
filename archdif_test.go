package archdif

import (
	"bufio"
	"os/exec"
	"testing"
)

var (
	ku bool
	status string = "?? .gitignore\nA  file\nAM file2\nD  file3"
)


func TestGitCmd(t *testing.T) {
	cmd, scanner := fakeGitOutput()
	cmd.Start()

	files := parseScanner(scanner, &ku)
	if len(files) > 2 {
		t.Fatalf("parseScanner returned wrong number, expected 2, got: %s", len(files))
	}

	cmd.Wait()
}

func TestGitCmdWithUntracked(t *testing.T)  {
	cmd, scanner := fakeGitOutput()
	cmd.Start()

	files := parseScanner(scanner, &ku)
	if len(files) > 3 {
		t.Fatalf("parseScanner returned wrong number, expected 3, got: %s", len(files))
	}

	cmd.Wait()
}

func fakeGitOutput() (*exec.Cmd, *bufio.Scanner) {
	cmd := exec.Command("echo", status)
	readcloser, _ := cmd.StdoutPipe()
	bufScanner := bufio.NewScanner(readcloser)

	return cmd, bufScanner
}

// to be continued...