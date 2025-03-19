package domain

const CollectionNameSSH = "ssh"

type SSH struct {
	Meta
	Label    string `json:"label" db:"label"`
	Host     string `json:"host" db:"host"`
	Port     int    `json:"port" db:"port"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
