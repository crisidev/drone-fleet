package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
	"github.com/op/go-logging"
)

const (
	APP_NAME    = "docker-fleet"
	APP_VERSION = "0.2"
	APP_SITE    = "https://github.com/crisidev/docker-fleet"
	sshDir      = "/root/.ssh"
)

var (
	log       = logging.MustGetLogger(APP_NAME)
	fleetArgs []string
	vargs     *Params
	workspace *drone.Workspace
	tunnel    *Tunnel
	repo      *drone.Repo
)

func main() {
	// system is used only in main
	system := drone.System{}

	// load json from stdin and parse it into variables
	plugin.Param("workspace", &workspace)
	plugin.Param("repo", &repo)
	plugin.Param("system", &system)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	LogSetup(vargs.Debug)
	log.Noticef("drone fleet scheduler deploy plugin, version %s", APP_VERSION)
	log.Noticef("deploying repo %s on %s (drone version %s)", repo.Name, system.Link, system.Version)
	log.Debugf("workspace path is %s", workspace.Path)

	// run fleet plugin
	RunFleetPlugin()
	log.Notice("done")
}

func RunFleetPlugin() {
	SetFleetConfig()
	fleetArgs = SetFleetArgs()

	if vargs.Scale == 0 {
		vargs.Scale = 1
		log.Warning("disabling rolling_update, this unit is not scaled")
		log.Debugf("deploying %d units", len(vargs.Units))
		RunFleetDeploy(0, vargs.Units)
	} else {
		log.Noticef("enable rolling_update, this unit is scaled to %d instances", vargs.Scale)
		log.Debugf("deploying %d units", len(vargs.Units)*vargs.Scale)
		for idx := 0; idx < vargs.Scale; idx++ {
			RunFleetDeploy(idx, vargs.Units)
		}
	}

	log.Noticef("deploy of repo %s succeded", repo.Name)
	os.Exit(0)
}

func SetFleetArgs() (args []string) {
	args = []string{fmt.Sprintf("--request-timeout=%s", strconv.Itoa(vargs.Timeout)), "--strict-host-key-checking=false"}
	if vargs.Debug {
		args = append(args, "--debug")
	}

	if vargs.Endpoint == "" {
		log.Error("etcd endpoint is mandatory, please add it to .drone.yml or via env variable DRONE_ETCD_ENDPOINT")
		os.Exit(1)
	} else {
		log.Infof("using etcd endpoint %s", vargs.Endpoint)
		args = append(args, fmt.Sprintf("--endpoint=%s", vargs.Endpoint))
	}

	if vargs.Tunnel != "" {
		log.Infof("tunnel specified, using SSH to %s", vargs.Tunnel)
		tunnel = &Tunnel{Tunnel: vargs.Tunnel, Key: workspace.Keys.Private, SSHDir: sshDir}
		fleetArgs = append(fleetArgs, tunnel.Handle()...)
	}

	return args
}

func SetFleetConfig() {
	if vargs.RollBack {
		log.Warning("roll_back functionality not developed yet")
	}

	if vargs.Timeout == 0 {
		vargs.Timeout = 10
	}

	if vargs.StartTimeout == 0 {
		vargs.StartTimeout = 60
	}

	if vargs.Sleep == 0 {
		vargs.Sleep = 5
	}

	if vargs.RollingSleep == 0 {
		vargs.RollingSleep = 120
	}

	// In case no units are passed, use repo name as unit name
	if vargs.Units == nil {
		if vargs.Scale > 1 {
			vargs.Units = []string{fmt.Sprintf("%s@.service", repo.Name)}
		} else {
			vargs.Units = []string{fmt.Sprintf("%s.service", repo.Name)}
		}
		log.Infof("parameter units not specified, defaulting to %s", vargs.Units[0])
	}
}

// Build a valid unit path, check the existence of the unit file and call the deploy function
func RunFleetDeploy(idx int, units []string) {
	for _, unit := range units {
		unitPath := path.Join(workspace.Path, unit)
		if _, err := os.Stat(unitPath); err == nil {
			fleet := Fleet{}
			fleet.Deploy(idx, unitPath)
		} else {
			log.Errorf("file %s not found", unitPath)
			os.Exit(1)
		}
		if vargs.Scale > 1 {
			DroneSleep(vargs.RollingSleep)
		}
	}
}
