package archdif

import (
	"bufio"
	"log"
	"os/exec"
)

const (
	untracked = "\x3F\x3F"
	deleted   = "\x20\x44"
	cdeleted  = "\x44\x20"
)

func GitCmd(keepUntracked *bool) (modifiedFiles []string) {
	cmd := exec.Command("git", "status", "--porcelain")
	statusReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	bufScanner := bufio.NewScanner(statusReader)

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	modifiedFiles = parseScanner(bufScanner, keepUntracked)

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	return
}

func parseScanner(s *bufio.Scanner, keepUntracked *bool) (modifiedFiles []string) {
	for s.Scan() {
		out := s.Text()
		mod := out[:2]

		switch {
		case mod == untracked:
			if *keepUntracked {
				modifiedFiles = append(modifiedFiles, out[3:])
			}
		case mod == deleted || mod == cdeleted:
			continue
		default:
			modifiedFiles = append(modifiedFiles, out[3:])
		}
	}

	return
}
