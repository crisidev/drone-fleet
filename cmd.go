package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// Run binary with arguments, check for errors, log output.
// Logging is done firing a gorouting and a WaitGroup is used
// to syncronize functions exit.
func (c *Cmd) ExecCmd() (err error) {
	log.Debugf("running command %s %s", c.CmdName, strings.Join(c.CmdArgs, " "))
	c.Cmd = exec.Command(c.CmdName, c.CmdArgs...)
	c.Wg = &sync.WaitGroup{}

	outScanner, errScanner, err := c.CreatePipes()
	if err != nil {
		return err
	}

	c.Wg.Add(2)
	go c.LogOutput(outScanner)
	go c.LogError(errScanner)

	err = c.Cmd.Start()
	if err != nil {
		log.Errorf("error starting cmd %s: %s", c.Cmd.Args, err)
		return err
	}

	err = c.Cmd.Wait()
	if err != nil {
		log.Errorf("error waiting for cmd %s: %s", c.Cmd.Args, err)
		return err
	}

	c.Wg.Wait()
	return nil
}

func (c *Cmd) CreatePipes() (*bufio.Scanner, *bufio.Scanner, error) {
	cmdReader, err := c.Cmd.StdoutPipe()
	if err != nil {
		log.Errorf("error creating StdoutPipe for cmd %s: %s", c.Cmd.Args, err)
		return nil, nil, err
	}
	outScanner := bufio.NewScanner(cmdReader)

	cmdErrReader, err := c.Cmd.StderrPipe()
	if err != nil {
		log.Errorf("error creating StderrPipe for cmd %s: %s", c.Cmd.Args, err)
		return nil, nil, err
	}
	errScanner := bufio.NewScanner(cmdErrReader)
	return outScanner, errScanner, nil
}

func (c *Cmd) ExecCmdOutput() (string, error) {
	log.Debugf("running command %s %s", c.CmdName, strings.Join(c.CmdArgs, " "))
	c.Cmd = exec.Command(c.CmdName, c.CmdArgs...)
	output, err := c.Cmd.Output()
	if err != nil {
		log.Errorf("error getting output for cmd %s : %s", c.Cmd.Args, err)
		return "", err
	}
	return string(output), nil
}

func (c *Cmd) LogOutput(scanner *bufio.Scanner) {
	defer c.Wg.Done()

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "SSH_AUTH_SOCK") {
			tunnel.EvalSSHAgent(line)
		}
		log.Info(line)
	}
}

func (c *Cmd) LogError(scanner *bufio.Scanner) {
	defer c.Wg.Done()

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "DEBUG") || strings.Contains(line, fmt.Sprintf("%s/id_rsa", sshDir)) {
			log.Debug(line)
		} else {
			log.Error(line)
		}
	}
}
