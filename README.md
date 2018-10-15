# synerex_alpha :
Demand/Supply Exchange Services for Synergic Mobility

# Introduction
Synerex alpha is an alpha version of Synergic Exchange and its support systems.
This project is supported by JST.

## Requirements
go 1.10 or later
nodejs / npm / yarn for web client development.

## How to start
Do 'go get' at all source directories to install dependent libraries
Then you can run 'go run xxx.go'

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

