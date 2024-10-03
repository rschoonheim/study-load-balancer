package linux

import "os/exec"

// NetworkInterfaceCreate - creates a virtual network interface
func NetworkInterfaceCreate(name, ip string) error {
	// Execute the commands
	//
	cmd := exec.Command("ip", "link", "add", "dev", name, "type", "bridge")

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	// Add ip address to net plan
	//
	cmd = exec.Command("ip", "addr", "add", ip+"/24", "dev", name)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	// Bring the interface up
	//
	cmd = exec.Command("ip", "link", "set", "dev", name, "up")
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	return nil
}
