package main

import (
	"kinesis"
	"fmt"
)

type EchoConsumer struct {
}

func (ec *EchoConsumer) Init(shardId string) error {
	fmt.Println("init: %s", shardId)
	return nil
}

func (ec *EchoConsumer) ProcessRecords(records []*kinesis.KclRecord,
			 checkpointer *kinesis.Checkpointer) error {
	for i := range records {
		fmt.Println("process: %s", records[i].DataB64)
	}

	checkpointer.CheckpointAll()
	return nil
}

func (ec *EchoConsumer) Shutdown(shutdownType kinesis.ShutdownType,
		checkpointer *kinesis.Checkpointer) error {
	fmt.Println("shutdown: %s", )
	if shutdownType == kinesis.GracefulShutdown {
		checkpointer.CheckpointAll()
	}
	return nil
}

func main() {
	var ec EchoConsumer
	kinesis.Run(&ec)
}
