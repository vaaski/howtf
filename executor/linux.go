//go:build linux

package executor

import "log"

func Execute(command string) {
	log.Println("Executor linux", command)
}
