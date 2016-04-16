package main

import (
	"os/exec"
	"sync"
)

type Params struct {
	Image        string   `json:"image"`
	Endpoint     string   `json:"endpoint"`
	Tunnel       string   `json:"tunnel"`
	Units        []string `json:"units"`
	Scale        int      `json:"scale"`
	Timeout      int      `json:"timeout"`
	Sleep        int      `json:"sleep"`
	RollingSleep int      `json:"rolling_sleep"`
	Destroy      bool     `json:"destroy"`
	Stop         bool     `json:"stop"`
	RollBack     bool     `json:"roll_back"`
	Debug        bool     `json:"debug"`
}

type Tunnel struct {
	Tunnel string
	User   string
	Key    string
	SSHDir string
	SSHKey string
}

type Fleet struct{}

type Cmd struct {
	CmdName string
	CmdArgs []string
	Wg      *sync.WaitGroup
	Cmd     *exec.Cmd
}
