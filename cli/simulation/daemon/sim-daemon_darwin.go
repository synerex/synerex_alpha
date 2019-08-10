package main

import "os"
import "strings"

func registerCurrentDir() error{
	return nil
}

func getRegisteredDir() (string, error) {
	d, e := os.Getwd()
	return d,e
}
func removeRegisteredDir() error{
	return nil
}

func binName(bn string) string{
	return bn
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

func getRegisteredGoPath() (string, error) {
	if !isDaemon {
		return getGoPath(), nil
	}
	return getGoPath(), nil
}
