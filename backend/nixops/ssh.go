package nixops

import (
	"backend/model"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func loadSSHKey() (ssh.Signer, error) {
	privateKeyPath := os.Getenv("NIXOPS_SSH_KEY_PATH")
	key, err := ssh.ParsePrivateKey([]byte(privateKeyPath))
	if err != nil {
		return nil, fmt.Errorf("Error parsing private key: %s", err)
	}
	return key, nil
}

func SendOne(client model.Client, command string) (string, error) {
	key, err := loadSSHKey()
	if err != nil {
		return "", fmt.Errorf("Error loading SSH key: %s", err)
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout:         0, //to adjust
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", client.IndexerIP+":22", config)
	if err != nil {
		return "", fmt.Errorf("SSH connection error : %s", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("SSH session error : %s", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	return string(output), err
}

func SendAll(clients []model.Client, command string) map[string]string {
	results := make(map[string]string) // Initialize the map to store results (to remove)
	for _, client := range clients {
		output, err := SendOne(client, command)
		if err != nil {
			results[client.Name] = fmt.Sprintf("Error: %s", err)
		} else {
			results[client.Name] = output
		}
	}
	return results
}
