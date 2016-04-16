package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

// Takes in input a file path, split it and return the file name
// and a list of strings representing representing the path
func (f *Fleet) SplitUnitPath(unitPath string) (string, []string) {
	unitPathSplit := []string{"/"}
	unitPathSplit = append(unitPathSplit, DeleteEmptyString(strings.Split(unitPath, "/"))...)
	if len(unitPathSplit) >= 3 {
		return unitPathSplit[len(unitPathSplit)-1], unitPathSplit[:len(unitPathSplit)-1]
	}
	return "", []string{}
}

// A scalable unit need to have and identifier in its name. This function
// append the identifier to a unit name.
// Example: myunit.mydomain.com@.service -> myunit.mydomain.com@1.service
// The identifier is the parameter idx.
func (f *Fleet) HandleScalableUnit(idx int, unitName string) (string, error) {
	if strings.Contains(unitName, "@") {
		unitNameSplit := strings.Split(unitName, "@")
		if len(unitNameSplit) == 2 {
			return fmt.Sprintf("%s@%d%s", unitNameSplit[0], idx, unitNameSplit[1]), nil
		}
	}
	log.Errorf("unit is scaled but unit file %s has wrong unit format, should be like unit@.service", unitName)
	return "", errors.New("scaled unit wrong file name")
}

func (f *Fleet) Deploy(idx int, unitPath string) (err error) {
	unitName, workPath := f.SplitUnitPath(unitPath)
	if vargs.Scale > 1 {
		unitName, err = f.HandleScalableUnit(idx, unitName)
		if err != nil {
			return err
		}
	}
	workPath = append(workPath, unitName)
	unitPath = path.Join(workPath...)

	f.Cmd("stop", unitName)

	if vargs.Destroy {
		f.Cmd("destroy", unitName)
	}

	if !vargs.Stop {
		if err := f.Cmd("start", unitPath); err != nil {
			log.Errorf("error starting unit %s", unitName)
			os.Exit(1)
		} else if err = f.CheckRunningUnit(unitName); err != nil {
			log.Errorf("unit %s has not started within %d seconds, deploy failed, stopping unit", unitName, vargs.StartTimeout)
			f.Cmd("stop", unitName)
			os.Exit(1)
		}
	} else {
		log.Critical("requested emergency stop")
	}

	return nil
}

func (f *Fleet) CheckRunningUnit(unitName string) (err error) {
	log.Infof("checking every 30 seconds if unit %s is started successfully, timeout %d", unitName, vargs.StartTimeout)
	for idx := 0; idx <= vargs.StartTimeout; idx += 30 {
		cmd := Cmd{CmdName: "/bin/fleetctl", CmdArgs: append(fleetArgs, "list-units")}
		output, err := cmd.ExecCmdOutput()
		if err != nil {
			return err
		}
		for _, line := range strings.Split(output, "\n") {
			if strings.Contains(line, unitName) {
				lineSplit := strings.Fields(line)
				if len(lineSplit) == 4 && lineSplit[2] == "active" {
					if lineSplit[3] == "running" {
						log.Infof("unit %s marked as running and started successfully", unitName)
					} else {
						log.Warningf("unit %s marked as active but not running, status is %s", unitName, lineSplit[3])
					}
					return nil
				}
			}
		}
		time.Sleep(time.Second * 30)
	}
	return errors.New("time excedeed")
}

// Wrapper around execCmd to run fleetctl with arguments
func (f *Fleet) Cmd(action, unitName string) (err error) {
	cmd := Cmd{CmdName: "/bin/fleetctl", CmdArgs: append(fleetArgs, action, unitName)}
	err = cmd.ExecCmd()
	DroneSleep(vargs.Sleep)
	return err
}
