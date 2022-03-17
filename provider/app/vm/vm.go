package vm

import (
	"log"
	"os/exec"
	"strconv"
)

func StartVM(id string, appName string, videoRelayPort, audioRelayPort, syncPort int) error {
	log.Printf("[%s] Spinning off VM\n", id)

	params := []string{
		id,
		strconv.Itoa(videoRelayPort),
		strconv.Itoa(audioRelayPort),
		strconv.Itoa(syncPort),
		appName,
	}
	cmd := exec.Command("./startVM.sh", params...)
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}

func StopVM(id, appName string) error {
	log.Printf("[%s] Stopping VM\n", id)

	params := []string{
		id,
		appName,
	}
	cmd := exec.Command("./stopVM.sh", params...)
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
