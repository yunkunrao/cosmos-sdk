localgaiad
==========

It is assumed that you have already `setup docker <https://docs.docker.com/engine/installation/>`__.

Description
-----------
Image for local gaiad testnets.

Add the gaiad binary to the image by attaching it in a folder to the `/gaiad` mount point.

It assumes that the configuration was created by the `tendermint testnet` command and it is also attached to the `/gaiad` mount point.

Example:
This example builds a linux gaiad binary under the `build/` folder, creates tendermint configuration for a single-node validator and runs the node:
```
cd $GOPATH/src/github.com/tendermint/tendermint

#Build binary
make build-linux

#Create configuration
docker run -v `pwd`/build:/tendermint tendermint/tendermint testnet --o . --v 1

#Run the node
docker run -v `pwd`/build:/tendermint tendermint/localgaiad
```

Logging
-------
Log is saved under the attached volume, in the `gaiad.log` file or to the file described in the `LOG` environment variable.

Special binaries
----------------
If you have multiple binaries with different names, you can specify which one to run with the `BINARY` environment variable. The path of the binary is relative to the attached volume.

docker-compose.yml
==================
This file creates a 4-node network using the localgaiad image. The nodes of the network are exposed to the host machine on ports 46656-46657, 46659-46660, 46661-46662, 46663-46664 respectively.

