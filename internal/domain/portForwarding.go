package domain

const CollectionNamePortForwarding = "port_forwarding"

type PortForwarding struct {
	Meta
	LocalAddress  string `json:"localAddress" db:"localAddress"`
	LocalPort     string `json:"localPort" db:"localPort"`
	RemoteAddress string `json:"remoteAddress" db:"remoteAddress"`
	RemotePort    string `json:"remotePort" db:"remotePort"`
	SSHId         string `json:"sshId" db:"sshId"`
	Running       bool   `json:"running" db:"running"`
}
