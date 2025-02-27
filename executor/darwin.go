//go:build darwin

package executor

import (
	"log"
	"os"
	"os/exec"
)

func Execute(command string) {
	log.Println("Executor darwin", command)

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		log.Println("Executor darwin error", err)
		return
	}

	log.Println("Executor darwin done")
}
