package git

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/simplesurance/baur/exec"
)

// CommitID return the commit id of HEAD by running git rev-parse in the passed
// directory
func CommitID(dir string) (string, error) {
	out, exitCode, err := exec.Command(dir, "git rev-parse HEAD")
	if err != nil {
		return "", errors.Wrap(err, "executing git rev-parse HEAD failed")
	}

	if exitCode != 0 {
		return "", errors.Wrapf(err, "executing git rev-parse HEAD failed, output: %q", out)
	}

	commitID := strings.TrimSpace(out)
	if len(commitID) == 0 {
		return "", errors.Wrap(err, "executing git rev-parse HEAD failed, no Stdout output")
	}

	return commitID, err
}

// LsFiles runs git ls-files in dir, passes args as argument and returns the
// output
func LsFiles(dir, args string) (string, error) {
	cmd := "git ls-files --error-unmatch " + args

	out, exitCode, err := exec.Command(dir, cmd)
	if err != nil {
		return "", errors.Wrapf(err, "executing %q failed", cmd)
	}

	if exitCode != 0 {
		if strings.Contains(out, "did not match any file(s)") {
			splt := strings.Split(out, "Did you forget to 'git add'")
			return "", errors.New(strings.Replace(splt[0], "\n", " ", -1))
		}

		return "", fmt.Errorf("%q exited with code %d, output: %q", cmd, exitCode, out)
	}

	return out, nil
}

// WorkTreeIsDirty returns true if the repository contains modified files,
// untracked files are considered, files in .gitignore are ignored
func WorkTreeIsDirty(dir string) (bool, error) {
	const cmd = "git status -s"

	out, exitCode, err := exec.Command(dir, cmd)
	if err != nil {
		return false, errors.Wrapf(err, "executing %q failed", cmd)
	}

	if exitCode != 0 {
		return false, fmt.Errorf("%q exited with code %d, output: %q", cmd, exitCode, out)
	}

	if len(out) == 0 {
		return false, nil
	}

	return true, nil
}
