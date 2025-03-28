package nixops

import (
	"backend/model"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func loadSSHKey() (ssh.Signer, error) {
	//ask for base64 of key

	privateKeyB64 := os.Getenv("NIXOPS_SSH_KEY_BASE64")
	// decode
	key, err := ssh.ParsePrivateKey([]byte(privateKeyPath))
	return key, err
}

// Do the same than SendOne without loading the key
// del sendone except loadKey, same for sendall
func send(...) {}

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
		Timeout:         0,                           //to adjust
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: test without
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

// TODO: NixOS rebuild switch...

// TODO: Git pull (automatically git stash git pull git unstash)

// TODO: Reboot

// TODO: OpenStatus (fastfetch)

// test("", "", "")
// test("")
// Instead of One, ALl:
func Send(command string, client ...model.Client) {
}

func Reboot(clients ...model.Client) {
	Send("reboot", clients...)
}

// ok := []string{"One", "Two"}
// x  test(ok)
// v  test(ok...)
