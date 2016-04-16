package main

import "testing"

func TestExecCmdOk(t *testing.T) {
	cmd := Cmd{CmdName: "echo", CmdArgs: []string{"test"}}
	err := cmd.ExecCmd()
	if err != nil {
		t.Fatalf("expected err == nil, got %s", err)
	}
}

func TestExecCmdErrorNonFound(t *testing.T) {
	cmd := Cmd{CmdName: "anonexistentcommand"}
	err := cmd.ExecCmd()
	if err == nil {
		t.Fatalf("expected err =! nil, got %s", err)
	}
}

func TestExecCmdBadReturnCode(t *testing.T) {
	cmd := Cmd{CmdName: "cd", CmdArgs: []string{"/anonexistentpath"}}
	err := cmd.ExecCmd()
	if err == nil {
		t.Fatalf("expected err =! nil, got %s", err)
	}
}
