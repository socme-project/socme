package wazuhapi

type WazuhAPI struct {
	Host     string
	Port     string
	IndexerPort string
	Username string
	Password string
	Token    string
	Insecure bool
}
