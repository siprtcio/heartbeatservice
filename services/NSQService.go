package services

import (
	"fmt"

	nsq "github.com/nsqio/go-nsq"
)

type NSQServices struct {
	producer           *nsq.Producer
	consumer           *nsq.Consumer
	maxInFlight        int
	concurrentHandlers int
}

func (this *NSQServices) initNSQServices(nsqdAddr string, maxInFlight int, concurrentHandlers int) error {
	this.maxInFlight = maxInFlight
	this.concurrentHandlers = concurrentHandlers
	// Create the configuration object and set the maxInFlight
	cfg := nsq.NewConfig()

	cfg.MaxInFlight = this.maxInFlight

	// Create the producer
	p, err := nsq.NewProducer(nsqdAddr, cfg)
	if err != nil {
		return err
	}
	this.producer = p
	return nil
}

func (this *NSQServices) IsQueueEmpty() bool {
	stats := this.consumer.Stats()
	if stats.MessagesFinished == (stats.MessagesReceived + stats.MessagesRequeued) {
		return true
	}
	return false
}

func (this *NSQServices) nsqSubscribe(nsqldTcpAddr, topicName, channelName string, hdlr nsq.HandlerFunc) error {
	fmt.Printf("Subscribe on %s/%s\n", topicName, channelName)

	// Create the configuration object and set the maxInFlight
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = this.maxInFlight
	var err error
	// Create the consumer with the given topic and chanel names
	this.consumer, err = nsq.NewConsumer(topicName, channelName, cfg)
	this.consumer.SetLoggerLevel(nsq.LogLevelError)
	if err != nil {
		return err
	}

	// Set the handler
	this.consumer.AddConcurrentHandlers(hdlr, this.concurrentHandlers)

	nsqlds := []string{nsqldTcpAddr}

	// Connect to the NSQ daemon
	if err := this.consumer.ConnectToNSQLookupds(nsqlds); err != nil {
		return err
	}

	// Wait for the consumer to stop.
	<-this.consumer.StopChan
	return nil
}

func (this *NSQServices) StopConsumer() {
	if this.consumer != nil {
		this.consumer.Stop()
	}
}

func (this *NSQServices) nsqPublish(topicName string, message []byte) error {
	return this.producer.Publish(topicName, message)
}
