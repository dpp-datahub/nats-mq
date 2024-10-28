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

package core

import (
	"fmt"

	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
	"github.com/nats-io/nats-mq/nats-mq/conf"
	stan "github.com/nats-io/stan.go"
)

// Stan2QueueConnector connects a STAN channel to an MQ Queue
type Stan2QueueConnector struct {
	BridgeConnector

	queue *ibmmq.MQObject
	sub   stan.Subscription
}

// NewStan2QueueConnector create a new Stan to MQ connector
func NewStan2QueueConnector(bridge *BridgeServer, config conf.ConnectorConfig) Connector {
	connector := &Stan2QueueConnector{}
	connector.init(bridge, config, fmt.Sprintf("STAN:%s to Queue:%s", config.Channel, config.Queue))
	return connector
}

// Start the connector
func (mq *Stan2QueueConnector) Start() error {
	mq.Lock()
	defer mq.Unlock()

	if !mq.bridge.CheckStan() {
		return fmt.Errorf("%s connector requires nats streaming to be available", mq.String())
	}

	mq.bridge.Logger().Tracef("starting connection %s", mq.String())

	err := mq.connectToMQ()
	if err != nil {
		return err
	}

	// Create the Object Descriptor that allows us to give the queue name
	qObject, err := mq.connectToQueue(mq.config.Queue, ibmmq.MQOO_OUTPUT)
	if err != nil {
		return err
	}

	mq.queue = qObject

	sub, err := mq.subscribeToChannel(mq.queue)
	if err != nil {
		return err
	}
	mq.sub = sub

	mq.stats.AddConnect()
	mq.bridge.Logger().Tracef("opened and reading %s", mq.config.Queue)
	mq.bridge.Logger().Noticef("started connection %s", mq.String())

	return nil
}

// Shutdown the connector
func (mq *Stan2QueueConnector) Shutdown() error {
	mq.Lock()
	defer mq.Unlock()
	mq.stats.AddDisconnect()

	mq.bridge.Logger().Noticef("shutting down connection %s", mq.String())

	if mq.sub != nil && mq.config.DurableName == "" { // Don't unsubscribe from durables
		mq.sub.Unsubscribe()
		mq.sub = nil
	}

	var err error

	queue := mq.queue
	mq.queue = nil

	if queue != nil {
		err = queue.Close(0)
	}

	if mq.qMgr != nil {
		_ = mq.qMgr.Disc()
		mq.qMgr = nil
		mq.bridge.Logger().Tracef("disconnected from queue manager for %s", mq.String())
	}
	return err // ignore the disconnect error
}

// CheckConnections ensures the nats/stan connection and report an error if it is down
func (mq *Stan2QueueConnector) CheckConnections() error {
	if !mq.bridge.CheckStan() {
		return fmt.Errorf("%s connector requires nats streaming to be available", mq.String())
	}
	return nil
}
