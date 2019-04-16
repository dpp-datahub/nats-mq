/*
 * Copyright 2012-2019 The NATS Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/nats-io/nats-mq/nats-mq/core"
)

var configFile string

func main() {
	var server *core.BridgeServer
	var err error

	flag.StringVar(&configFile, "c", "", "configuration filepath")
	flag.Parse()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGHUP)

		for {
			signal := <-sigChan

			if signal == os.Interrupt {
				if server.Logger() != nil {
					fmt.Println() // clear the line for the control-C
					server.Logger().Noticef("received sig-interrupt, shutting down")
				}
				server.Stop()
				os.Exit(0)
			}

			if signal == syscall.SIGHUP {
				if server.Logger() != nil {
					server.Logger().Errorf("received sig-hup, restarting")
				}
				server.Stop()
				server := core.NewBridgeServer()
				server.LoadConfigFile(configFile)
				err = server.Start()

				if err != nil {
					if server.Logger() != nil {
						server.Logger().Errorf("error starting bridge, %s", err.Error())
					} else {
						log.Printf("error starting bridge, %s", err.Error())
					}
					server.Stop()
					os.Exit(0)
				}
			}
		}
	}()

	server = core.NewBridgeServer()
	server.LoadConfigFile(configFile)
	err = server.Start()

	if err != nil {
		if server.Logger() != nil {
			server.Logger().Errorf("error starting bridge, %s", err.Error())
		} else {
			log.Printf("error starting bridge, %s", err.Error())
		}
		server.Stop()
		os.Exit(0)
	}

	// exit main but keep running goroutines
	runtime.Goexit()
}
