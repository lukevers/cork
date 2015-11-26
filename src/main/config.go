/* vim: set autoindent noexpandtab tabstop=4 shiftwidth=4: */
package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
)

type Config struct {
	Http HttpServer
	Ssh  SshServer
}

func LoadAllConfigs(folder string) *Config {
	c := Config{}

	// Httpd
	if _, err := toml.Decode(file(folder, "http.toml"), &c.Http); err != nil {
		log.Fatal("Could not decode config file:", err)
	}

	// Sshd
	if _, err := toml.Decode(file(folder, "ssh.toml"), &c.Ssh); err != nil {
		log.Fatal("Could not decode config file:", err)
	}

	return &c
}

func file(folder, file string) string {
	data, err := ioutil.ReadFile(folder + file)
	if err != nil {
		log.Fatal("Could not read file:", err)
	}

	return string(data)
}

