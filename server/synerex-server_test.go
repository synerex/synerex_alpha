package main

import (
	"testing"
	"os"
)

var (
	testServInfo *synerexServerInfo
)


func initNewServer(){

	testServInfo = newServerInfo()

}

func closeServer(){

}


func TestSMarketServer_Confirm(t *testing.T) {

}

func TestSMarketServer_SubscribeSupply(t *testing.T) {

}

func TestMain(m *testing.M){

	initNewServer()

	code := m.Run()

	closeServer()

	os.Exit(code)

}