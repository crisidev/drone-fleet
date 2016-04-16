package main

import (
	"bufio"
	"fmt"
	"io"
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

	cmdReader, err := c.CmdGetStdOutPipe()
	if err != nil {
		return err
	}
	outScanner := bufio.NewScanner(cmdReader)

	cmdErrReader, err := c.CmdGetStdErrPipe()
	if err != nil {
		return err
	}
	errScanner := bufio.NewScanner(cmdErrReader)

	c.Wg.Add(2)
	go c.CmdLogOutput(outScanner)
	go c.CmdLogError(errScanner)

	err = c.CmdStart()
	if err != nil {
		return err
	}

	err = c.CmdWait()
	if err != nil {
		return err
	}

	c.Wg.Wait()
	return nil
}

func (c *Cmd) CmdStart() (err error) {
	err = c.Cmd.Start()
	if err != nil {
		log.Errorf("error starting cmd %s: %s", c.Cmd.Args, err)
		return err
	}
	return nil
}

func (c *Cmd) CmdWait() (err error) {
	err = c.Cmd.Wait()
	if err != nil {
		log.Errorf("error waiting for cmd %s: %s", c.Cmd.Args, err)
		return err
	}
	return nil
}

func (c *Cmd) CmdGetStdErrPipe() (io.Reader, error) {
	cmdErrReader, err := c.Cmd.StderrPipe()
	if err != nil {
		log.Errorf("error creating StderrPipe for cmd %s: %s", c.Cmd.Args, err)
		return nil, err
	}
	return cmdErrReader, nil
}

func (c *Cmd) CmdGetStdOutPipe() (io.Reader, error) {
	cmdReader, err := c.Cmd.StdoutPipe()
	if err != nil {
		log.Errorf("error creating StdoutPipe for cmd %s: %s", c.Cmd.Args, err)
		return nil, err
	}
	return cmdReader, nil
}

func (c *Cmd) CmdLogOutput(scanner *bufio.Scanner) {
	defer c.Wg.Done()

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "SSH_AUTH_SOCK") {
			tunnel.EvalSSHAgent(line)
		}
		log.Info(line)
	}
}

func (c *Cmd) CmdLogError(scanner *bufio.Scanner) {
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
