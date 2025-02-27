//go:build windows

package executor

import "log"

func Execute(command string) {
	log.Println("Executor windows", command)
}
