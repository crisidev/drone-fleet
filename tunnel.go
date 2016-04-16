package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func (t *Tunnel) HandleTunnel() (args []string) {
	args = t.ParseTunnel()
	t.CreateSSHDir()
	if err := t.WriteSSHKey(); err != nil {
		os.Exit(1)
	}

	log.Notice("starting ssh-agent for tunnel authentication")
	cmd := Cmd{CmdName: "ssh-agent"}
	if err := cmd.ExecCmd(); err != nil {
		log.Error("error starting ssh-agent for tunnel authentication")
		os.Exit(1)
	}

	log.Notice("adding repo SSH private key to ssh-agent")
	cmd = Cmd{CmdName: "ssh-add", CmdArgs: []string{t.SSHKey}}
	if err := cmd.ExecCmd(); err != nil {
		log.Error("error adding repo SSH private key to ssh-agent")
		os.Exit(1)
	}

	return args
}

func (t *Tunnel) ParseTunnel() []string {
	tunnelSplit := strings.Split(t.Tunnel, "@")
	if len(tunnelSplit) == 1 {
		return []string{fmt.Sprintf("--tunnel=%s", t.Tunnel)}
	} else {
		log.Debugf("user %s will be used for the SSH tunnel", tunnelSplit[0])
		t.Tunnel = tunnelSplit[1]
		t.User = tunnelSplit[0]
		return []string{fmt.Sprintf("--tunnel=%s", tunnelSplit[1]), fmt.Sprintf("--ssh-username=%s", tunnelSplit[0])}
	}
}

func (t *Tunnel) CreateSSHDir() {
	log.Debugf("creating SSH home dir %s", t.SSHDir)
	os.MkdirAll(t.SSHDir, 0744)
	t.SSHKey = path.Join(t.SSHDir, "id_rsa")
}

func (t *Tunnel) RemoveSSHKey() {
	log.Debugf("removing old SSH private key")
	os.Remove(t.SSHKey)
}

func (t *Tunnel) WriteSSHKey() (err error) {
	t.RemoveSSHKey()
	log.Debugf("writing repo SSH private key to %s", t.SSHKey)
	err = ioutil.WriteFile(t.SSHKey, []byte(t.Key), 0400)
	if err != nil {
		log.Errorf("error writing %s: %s", t.SSHKey, err.Error())
		return err
	}
	return nil
}

func (t *Tunnel) EvalSSHAgent(line string) {
	if os.Getenv("SSH_AUTH_SOCK") == "" {
		log.Debug("matched string SSH_AUTH_SOCK")
		agentSockSplit := strings.Fields(line)
		agentSockSplit = strings.Split(agentSockSplit[0], "=")
		if len(agentSockSplit) == 2 {
			envVar := agentSockSplit[0]
			agentSock := strings.Replace(agentSockSplit[1], ";", "", 1)
			if _, err := os.Stat(agentSock); err == nil {
				log.Infof("setting up environment variable %s=%s", envVar, agentSock)
				os.Setenv(envVar, agentSock)
			}
		}
	} else {
		log.Info("environment variable SSH_AUTH_SOCK already present")
	}
}
