# drone-fleet

[![Build Status](https://drone.crisidev.org/api/badges/crisidev/drone-fleet/status.svg)](https://drone.crisidev.org/crisidev/drone-fleet)
[![](https://badge.imagelayers.io/crisidev/drone-fleet:latest.svg)](https://imagelayers.io/?images=crisidev/drone-fleet:latest 'Get your own badge on imagelayers.io')

Drone plugin to deploy unitfiles on [CoreOS](https://coreos.com) using [Fleet](https://github.com/coreos/fleet)

For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Binary

Build the binary using `make`:

```
make
make test
```

### Example

```sh
$ ./build/drone-fleet <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "owner": "drone",
        "name": "drone",
        "full_name": "drone/drone"
    },
    "system": {
        "link": "http://beta.drone.io",
        "version": "0.4"
    },
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com",
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone",
        "keys": {
          "private": "lololol"
        }
    },
    "vargs": {
        "image": "crisidev/drone-fleet",
        "endpoint": "http://etcd.mydomain.com:4001",
        "units": [
            "fleet/unit.service",
            "fleet/unit.timer"
        ],
        "timeout": 2,
        "destroy_unit": true,
        "debug": true
    }
}
EOF
```

## Docker

Build the container using `make`:

```
$ make docker
```

### Example

```sh
docker run -i plugins/drone-fleet <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "owner": "drone",
        "name": "drone",
        "full_name": "drone/drone"
    },
    "system": {
        "link": "http://beta.drone.io",
        "version": "0.4"
    },
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com",
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone",
        "keys": {
          "private": "lololol"
        }
    },
    "vargs": {
        "image": "crisidev/drone-fleet",
        "endpoint": "http://etcd.mydomain.com:4001",
        "units": [
            "fleet/unit.service",
            "fleet/unit.timer"
        ],
        "timeout": 2,
        "destroy_unit": true,
        "debug": true
    }
}
EOF
```

