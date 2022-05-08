// Command gw finds a Dradle wrapper (gradlew) script in the current or parent directories and runs it.
package main // import "rmazur.io/gw"

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const gradlew = "gradlew"

func main() {
	cwd, err := filepath.Abs(".")
	if err != nil {
		log.Fatal("unable to resolve current path:", err)
	}
	name := gradlew
	if runtime.GOOS == "windows" {
		name = name + ".bat"
	}

	for !checkPath(cwd, name) {
		cwd = filepath.Dir(cwd)
		if len(cwd) < 2 {
			log.Fatal("gradle wrapper not found")
		}
	}

	gw := exec.Command(filepath.Join(cwd, name))
	gw.Args = os.Args[1:]
	gw.Stdout = os.Stdout
	gw.Stdin = os.Stdin
	err = gw.Run()
	if err != nil {
		if eErr, ok := err.(*exec.ExitError); ok {
			os.Exit(eErr.ExitCode())
		} else {
			log.Fatalf("problems running %s: %s", gw.Path, err)
		}
	}
}

func checkPath(p, name string) bool {
	if fi, err := os.Stat(filepath.Join(p, name)); err == nil {
		return !fi.IsDir()
	}
	return false
}
