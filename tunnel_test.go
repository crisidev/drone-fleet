package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseTunnel(t *testing.T) {
	testTunnel := "bastion.mydomain.com"
	tunnel := Tunnel{Tunnel: testTunnel}
	testTunnelSplit := tunnel.Parse()
	if testTunnelSplit[0] != fmt.Sprintf("--tunnel=%s", testTunnel) {
		t.Fatalf("expected --tunnel=%s, got %s", testTunnel, testTunnelSplit[0])
	}

	testTunnel = "bastion.mydomain.com"
	testUser := "myuser"
	tunnel = Tunnel{Tunnel: "myuser@bastion.mydomain.com"}
	testTunnelSplit = tunnel.Parse()
	if testTunnelSplit[0] != fmt.Sprintf("--tunnel=%s", testTunnel) {
		t.Fatalf("expected --tunnel=%s, got %s", testTunnel, testTunnelSplit[0])
	}
	if testTunnelSplit[1] != fmt.Sprintf("--ssh-username=%s", testUser) {
		t.Fatalf("expected --ssh-username=%s, got %s", testUser, testTunnelSplit[1])
	}
}

func TestCreateSSHDir(t *testing.T) {
	dir, _ := ioutil.TempDir("", "atestdir")
	tunnel := Tunnel{SSHDir: dir}
	tunnel.CreateSSHDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("expected sshDir %s to exists, does not exist", dir)
		_ = os.Remove(dir)
	}
	_ = os.Remove(dir)
}

func TestWriteSSHKeyOk(t *testing.T) {
	dir, _ := ioutil.TempDir("", "atestdir")
	tunnel := Tunnel{SSHDir: dir, Key: "avalue"}
	tunnel.CreateSSHDir()
	_ = tunnel.WriteSSHKey()
	if _, err := os.Stat(tunnel.SSHKey); os.IsNotExist(err) {
		t.Fatalf("expected keyFile %s to exists, does not exist", tunnel.SSHKey)
		_ = os.RemoveAll(dir)
	}
	_ = os.RemoveAll(dir)
}

func TestWriteSSHKeyError(t *testing.T) {
	tunnel := Tunnel{SSHKey: "anonexistentdir/id_rsa", Key: "avalue"}
	err := tunnel.WriteSSHKey()
	if err == nil {
		t.Fatalf("expected err != nil, got %s", err)
	}
}
