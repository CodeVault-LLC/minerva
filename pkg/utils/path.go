package utils

import "os/exec"

// isNmapInstalled checks if Nmap is installed on the system.
func IsNmapInstalled() bool {
	cmd := exec.Command("nmap", "-v")
	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}
