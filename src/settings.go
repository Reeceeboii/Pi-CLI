package main

import "fmt"

const (
	defaultAddress  = "127.0.0.1"
	defaultPort     = 80
	defaultRefreshS = 1
)

type Settings struct {
	PiHoleAddress string `json:"pi_hole_address"`
	PiHolePort    int64  `json:"pi_hole_port"`
	RefreshS      int64  `json:"refresh_s"`
}

func (settings *Settings) loadFromFile() {
	fmt.Println("looking for file...")
}
