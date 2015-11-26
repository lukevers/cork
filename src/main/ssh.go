/* vim: set autoindent noexpandtab tabstop=4 shiftwidth=4: */
package main

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"log"
	"fmt"
	"errors"
	encssh "github.com/ianmcmahon/encoding_ssh"
	"reflect"
)

type SshServer struct {
	Auth     bool
	Port     int
	Host     string
	Version  string
	Rsa      Rsa
	config   *ssh.ServerConfig
	listener net.Listener
}

type Rsa struct {
	Public  string
	Private string
}

func (s *SshServer) Run() {
	s.config = &ssh.ServerConfig{
		ServerVersion: "SSH-2.0-Cork",
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			good := false
			for id, k := range Keys {
				if id != conn.User() {
					continue
				}

				pub_key, _ := encssh.DecodePublicKey(k)
				pk, _ := ssh.NewPublicKey(pub_key)

				good = reflect.DeepEqual(pk, key)
				if good {
					UserHasBeenAuthorized(id)
					delete(Keys, id)
					break
				}
			}

			if !good {
				return nil, errors.New("Bad")
			}

			return nil, nil
		},
	}

	bytes, err := ioutil.ReadFile(s.Rsa.Private)
	if err != nil {
		log.Fatal("Could not read private key:", err)
	}

	private, err := ssh.ParsePrivateKey(bytes)
	if err != nil {
		log.Fatal("Could not parse private key:", err)
	}

	s.config.AddHostKey(private)

	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		log.Fatal("Could not listen on connection:", err)
	}

	defer s.listener.Close()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Could not accept connections:", err)
		}

		go func(c net.Conn) {
			_, channels, requests, _ := ssh.NewServerConn(c, s.config)
			go ssh.DiscardRequests(requests)
			handle(channels)
			c.Close()
		}(conn)
	}
}

func handle(channels <-chan ssh.NewChannel) {
	for ch := range channels {
		if ch.ChannelType() != "session" {
			ch.Reject(ssh.UnknownChannelType, "Unknown channel type")
			continue
		}

		channel, _, err := ch.Accept()
		if err != nil {
			log.Println("Could not accept channel", err)
		}

		channel.Write([]byte("Authorized! "))
		channel.Close()
	}
}
