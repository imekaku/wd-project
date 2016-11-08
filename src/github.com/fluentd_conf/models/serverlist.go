package models

type ServerInfo struct {
	Data []HostInfo `json:data`
}

type HostInfo struct {
	Hostname string `json:hostname`
	Id       int64  `json:id`
}
