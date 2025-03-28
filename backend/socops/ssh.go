package nixops

import (
	"backend/model"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func loadSSHKey() (ssh.Signer, error) {
	privateKeyB64 := os.Getenv("SOCOPS_SSH_KEY_BASE64")
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		return nil, fmt.Errorf("Error decoding SSH key: %s", err)
	}
	key, err := ssh.ParsePrivateKey((privateKey)) // TODO: check if byte or string
	return key, err
}

func SendCommand(command string, clients ...model.Client) map[string]string {
	results := make(map[string]string)

	// Load the SSH key once
	key, err := loadSSHKey()
	if err != nil {
		for _, client := range clients {
			results[client.Name] = fmt.Sprintf("Error loading SSH key: %s", err)
		}
		return results
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout: 0, // Adjust if needed
		// HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: test without
	}

	for _, client := range clients {
		conn, err := ssh.Dial("tcp", client.IndexerIP+":22", config)
		if err != nil {
			results[client.Name] = fmt.Sprintf("SSH connection error: %s", err)
			continue
		}
		defer conn.Close()

		session, err := conn.NewSession()
		if err != nil {
			results[client.Name] = fmt.Sprintf("SSH session error: %s", err)
			continue
		}
		defer session.Close()

		output, err := session.CombinedOutput(command)
		if err != nil {
			results[client.Name] = fmt.Sprintf("Command execution error: %s", err)
		} else {
			results[client.Name] = string(output)
		}
	}

	return results
}

func Reboot(clients ...model.Client) {
	SendCommand("reboot", clients...)
}

func GitPull(clients ...model.Client) {
	SendCommand("git pull", clients...)
}

func OpenStatus(clients ...model.Client) {
	SendCommand("fastfetch", clients...)
}

func NixOSRebuildSwitch(clients ...model.Client) {
	SendCommand(
		"sudo nixos-rebuild switch --flake ${configDirectory}#${hostname}",
		clients...) // TODO: adjust this
}
