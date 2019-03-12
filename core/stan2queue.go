package core

import (
	"fmt"
	"github.com/nats-io/go-nats-streaming"
	"sync"

	"github.com/ibm-messaging/mq-golang/ibmmq"
)

// Stan2QueueConnector connects a STAN channel to an MQ Queue
type Stan2QueueConnector struct {
	sync.Mutex

	config ConnectionConfig
	bridge *BridgeServer

	qMgr  *ibmmq.MQQueueManager
	queue *ibmmq.MQObject

	sub stan.Subscription
}

// NewStan2QueueConnector create a new Stan to MQ connector
func NewStan2QueueConnector(bridge *BridgeServer, config ConnectionConfig) Connector {
	return &Stan2QueueConnector{
		config: config,
		bridge: bridge,
	}
}

func (mq *Stan2QueueConnector) String() string {
	return fmt.Sprintf("STAN:%s to Queue:%s", mq.config.Channel, mq.config.Queue)
}

// Config returns the configuraiton for this connector
func (mq *Stan2QueueConnector) Config() ConnectionConfig {
	return mq.config
}

// Start the connector
func (mq *Stan2QueueConnector) Start() error {
	mq.Lock()
	defer mq.Unlock()

	if mq.bridge.stan == nil {
		return fmt.Errorf("%s connector requires nats streaming to be available", mq.String())
	}

	mqconfig := mq.config.MQ
	queueName := mq.config.Queue

	mq.bridge.Logger.Tracef("starting connection %s", mq.String())

	qMgr, err := ConnectToQueueManager(mqconfig)
	if err != nil {
		return err
	}

	mq.bridge.Logger.Tracef("connected to queue manager %s at %s as %s for %s", mqconfig.QueueManager, mqconfig.ConnectionName, mqconfig.ChannelName, mq.String())

	mq.qMgr = qMgr

	// Create the Object Descriptor that allows us to give the queue name
	mqod := ibmmq.NewMQOD()
	openOptions := ibmmq.MQOO_OUTPUT
	mqod.ObjectType = ibmmq.MQOT_Q
	mqod.ObjectName = queueName

	qObject, err := mq.qMgr.Open(mqod, openOptions)

	if err != nil {
		return err
	}

	mq.queue = &qObject

	sub, err := mq.bridge.stan.Subscribe(mq.config.Channel, mq.messageHandler)

	if err != nil {
		return err
	}

	mq.sub = sub

	mq.bridge.Logger.Tracef("opened and reading %s", queueName)
	mq.bridge.Logger.Noticef("started connection %s", mq.String())

	return nil
}

func (mq *Stan2QueueConnector) messageHandler(m *stan.Msg) {
	qmgrFlag := mq.qMgr

	if mq.config.ExcludeHeaders {
		qmgrFlag = nil
	}

	mqmd, handle, buffer, err := mq.bridge.natsToMQMessage(m.Data, "", qmgrFlag)

	pmo := ibmmq.NewMQPMO()
	pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT
	pmo.OriginalMsgHandle = handle

	// Now put the message to the queue
	err = mq.queue.Put(mqmd, pmo, buffer)

	if err != nil {
		mq.bridge.Logger.Noticef("MQ publish failure, %s, %s", mq.String(), err.Error())
	}
}

// Shutdown the connector
func (mq *Stan2QueueConnector) Shutdown() error {
	mq.Lock()
	defer mq.Unlock()

	mq.bridge.Logger.Noticef("shutting down connection %s", mq.String())

	var err error

	queue := mq.queue
	mq.queue = nil

	if queue != nil {
		err = queue.Close(0)
	}

	if mq.qMgr != nil {
		_ = mq.qMgr.Disc()
		mq.bridge.Logger.Tracef("disconnected from queue manager for %s", mq.String())
	}

	if mq.sub != nil {
		mq.sub.Unsubscribe()
		mq.sub = nil
	}

	return err // ignore the disconnect error
}