package nixops

import (
	"backend/model"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// CommandOutput represents the output of an executed command on a client
type CommandOutput struct {
	ClientName string
	Message    string
	Error      bool
}

func loadSSHKey() (ssh.Signer, error) {
	privateKeyB64 := os.Getenv("OPSME_SSH_KEY_BASE64")
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		return nil, fmt.Errorf("Error decoding SSH key: %s", err)
	}
	key, err := ssh.ParsePrivateKey(privateKey)
	return key, err
}

func SendCommand(command string, clients ...model.Client) []CommandOutput {
	var results []CommandOutput

	// Load the SSH key once
	key, err := loadSSHKey()
	if err != nil {
		for _, client := range clients {
			results = append(results, CommandOutput{
				ClientName: client.Name,
				Message:    fmt.Sprintf("Error loading SSH key: %s", err),
				Error:      true,
			})
		}
		return results
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout: 0, // Adjust if needed
		// HostKeyCallback: ssh.InsecureIgnoreHostKey(), //TODO: TESTING WITHOUT
	}

	for _, client := range clients {
		conn, err := ssh.Dial("tcp", client.IndexerIP+":22", config)
		if err != nil {
			results = append(results, CommandOutput{
				ClientName: client.Name,
				Message:    fmt.Sprintf("SSH connection error: %s", err),
				Error:      true,
			})
			continue
		}
		defer conn.Close()

		session, err := conn.NewSession()
		if err != nil {
			results = append(results, CommandOutput{
				ClientName: client.Name,
				Message:    fmt.Sprintf("SSH session error: %s", err),
				Error:      true,
			})
			continue
		}
		defer session.Close()

		output, err := session.CombinedOutput(command)
		if err != nil {
			results = append(results, CommandOutput{
				ClientName: client.Name,
				Message:    fmt.Sprintf("Command execution error: %s", err),
				Error:      true,
			})
		} else {
			results = append(results, CommandOutput{
				ClientName: client.Name,
				Message:    string(output),
				Error:      false,
			})
		}
	}

	return results
}

func Reboot(clients ...model.Client) []CommandOutput {
	return SendCommand("reboot", clients...)
}

func FetchMe(clients ...model.Client) []CommandOutput {
	return SendCommand("fastfetch", clients...)
}

func Update(clients ...model.Client) []CommandOutput {
	return SendCommand(
		"cd /etc/nixos && git stash && git pull && nix flake update && sudo nixos-rebuild switch --upgrade --flake /etc/nixos#artemis && git stash pop",
		clients...)
}

