package headscale

import (
	"time"

	hsClient "github.com/hibare/headscale-client-go"
	"golang.org/x/net/context"
)

type Server struct {
	Server hsClient.HeadscaleClientInterface
}

func ConnectServer(url, token string) (*Server, error) {
	client, err := hsClient.NewClient(url, token, hsClient.HeadscaleClientOptions{})
	return &Server{Server: client}, err
}

func (s *Server) NewClient(name string) (preauthkey string, err error) {
	client, err := s.Server.Users().Create(context.Background(), name)
	if err != nil {
		return "", err
	}
	keys, err := s.Server.PreAuthKeys().
		Create(context.Background(), client.Name, true, false, time.Now().Add(24*time.Hour), []string{"client"})

	if len(keys.PreAuthKey) == 0 {
		return "", err
	}

	return keys.PreAuthKey[0].Key, err
}
