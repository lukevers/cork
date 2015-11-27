/* vim: set autoindent noexpandtab tabstop=4 shiftwidth=4: */
package main

import (
	"log"
	"sync"
)

var (
	Cork *Config
	waitgp sync.WaitGroup
	Keys   map[string]string = make(map[string]string)
	
)

func init() {
	log.Println("Reading configuration files")
	Cork = LoadAllConfigs("resources/config/")
	waitgp.Add(1)
}

func main() {
	log.Println("Starting http server on", Cork.Http.Host, Cork.Http.Port)
	go Cork.Http.Run()
	log.Println("Starting ssh server on", Cork.Ssh.Host, Cork.Ssh.Port)
	go Cork.Ssh.Run()
	waitgp.Wait()
}
