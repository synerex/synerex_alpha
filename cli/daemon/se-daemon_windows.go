package main

// Daemon registration dir for Synerex-Engine.
// for Windows Service

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"strings"
)

func registerCurrentDir() error{
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Synerex\se-daemon`, registry.ALL_ACCESS)
	if err != nil {
		return 	fmt.Errorf("se-daemon: registry error. you need to have administrator privilege: %s",err.Error())
	}
	d, _  := os.Getwd()
	//Todo: currently register current dir. But we need to change?
	err = k.SetStringValue("basedir",d)
	if err != nil {
		return fmt.Errorf("se-daemon: Registry set basedir error: %s", err.Error())
	}
	err = k.SetStringValue("gopath",getGoPath())
	if err != nil {
		return fmt.Errorf("se-daemon: Registry set gopath error: %s", err.Error())
	}
	return nil
}

func getRegisteredDir() (string, error){
	// if the system is not running service/daemon mode, just return currentwd
	if !isDaemon {
		return os.Getwd()
	}

	k, _ , err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Synerex\se-daemon`, registry.ALL_ACCESS)
	if err != nil {
		// cannot continue may panic...
		return "", 	fmt.Errorf("se-daemon: Registry query eeror. you need to have administrator privilege: %s", err.Error())
	}
	st, _ ,e2 := k.GetStringValue("basedir")
	if e2 != nil {
		return "", fmt.Errorf("se-daemon: Registry get error: %s", e2.Error())
	}
	return st, nil
}

func getGoPath() string{
	env := os.Environ()
	for _, ev := range env {
		if strings.Contains(ev,"GOPATH=") {
			return ev
		}
	}
	return ""
}

func getRegisteredGoPath() (string, error){
	if !isDaemon {
		return getGoPath(), nil
	}

	k, _ , err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Synerex\se-daemon`, registry.ALL_ACCESS)
	if err != nil {
		// cannot continue may panic...
		return "", 	fmt.Errorf("se-daemon: Registry query eeror. you need to have administrator privilege: %s", err.Error())
	}
	st, _ ,e2 := k.GetStringValue("gopath")
	if e2 != nil {
		return "", fmt.Errorf("se-daemon: Registry get gopath error: %s", e2.Error())
	}
	return st, nil
}

func removeRegisteredDir() error{
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Synerex\se-daemon`, registry.ALL_ACCESS)
	if err != nil {
		return 	fmt.Errorf("se-daemon: Registry error. you need to have administrator privilege: %s",err.Error())
	}
	err = k.DeleteValue("basedir")
	if err != nil {
		return fmt.Errorf("se-daemon: Registry basedir remove error: %s", err.Error())
	}
	k, _, err = registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\Synerex`, registry.ALL_ACCESS)
	err = registry.DeleteKey(k,"se-daemon")
	if err != nil {
		return fmt.Errorf("se-daemon: Registry synerex remove error: %s", err.Error())
	}

	return nil
}


func binName(bn string) string{
	return bn+".exe"
}