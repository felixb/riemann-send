package main

import (
	"fmt"
)

type ServerConfig struct {
	Host    string `json:"host,omitempty"`
	Port    uint   `json:"port,omitempty"`
	Mode    string `json:"mode,omitempty"`
	Timeout int64  `json:"timeout,omitempty"`
}

func NewServerConfig() ServerConfig {
	return ServerConfig{
		Host: "localhost",
		Port: 5555,
		Mode: "tcp",
		Timeout: 0,
	}
}

func (config ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
