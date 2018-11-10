# synerex_alpha :  [![CircleCI](https://circleci.com/gh/synerex/synerex_alpha.svg?style=shild)](https://circleci.com/gh/synerex/synerex_alpha) 
Demand/Supply Exchange Services for Synergic Mobility

# Introduction
Synerex alpha is an alpha version of Synergic Exchange and its support systems.
This project is supported by JST.

## Requirements
go 1.11 or later (we use go.mod files for module dependencies)
nodejs(10.13.0) / npm(6.4.1) / yarn(1.12.1) for web client development.

## How to start
Do 'go get' at all source directories to install dependent libraries.

Starting from SynerexEngine.
```
  cd cli/daemon
  go build
  ./se-daemon
```

Then move to se directory.
```
cd ../se
go build
./se build all
```

## Source Directories

### cli
#### deamon
 se-daemon for cli service
  It can start all server.
 ```
 go build se-daemon.go se-daemon_[os].go
 ```


#### se
 command line client for Synerex Engine
```
 go build se.go  // build se command
 
 se run all     // start all servers and providers
 se stop all    // stop all servers and providers
 se ps -l       // list current running server and providers
```

#### api

Protocl Buffer / gRPC definition of Synergex API

#### server

Synerex Server draft version

#### provider

Synerex Service Providers

#####    ad

#####    taxi

#####    multi

#####    user

#####    fleet

#### sxutil

Synerex Utility Package Both server and provider package will
use this.

monitor Monitoring Server

