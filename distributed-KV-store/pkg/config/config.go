package config

type Config struct {
	BindAddr string
	NodeID   string
	Peers    []string
	RaftDir  string
	DataDir  string
}
