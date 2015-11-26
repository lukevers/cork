/* vim: set autoindent noexpandtab tabstop=4 shiftwidth=4: */
package main

import (
	"sync"
	"log"
)

var (
	Conf   *Config
	waitgp sync.WaitGroup
	Keys  map[string]string = make(map[string]string)
)

func init() {
	log.Println("Reading configuration files")
	Conf = LoadAllConfigs("resources/config/")
	waitgp.Add(1)
}

func main() {
	log.Println("Starting http server on", Conf.Http.Host, Conf.Http.Port)
	go Conf.Http.Run()
	log.Println("Starting ssh server on", Conf.Ssh.Host, Conf.Ssh.Port)
	go Conf.Ssh.Run()
	waitgp.Wait()
}
