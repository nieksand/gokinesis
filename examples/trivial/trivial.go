/*
 * A trivial example consuming from Kinesis using the KCL daemon.
 *
 * Just echos messages it sees to stderr.
 */
package main

import (
	"fmt"
	"kinesis"
	"os"
)

type EchoConsumer struct {
}

func (ec *EchoConsumer) Init(shardId string) error {
	fmt.Fprintf(os.Stderr, "init: %s\n", shardId)
	return nil
}

func (ec *EchoConsumer) ProcessRecords(records []*kinesis.KclRecord,
	checkpointer *kinesis.Checkpointer) error {
	for i := range records {
		fmt.Fprintf(os.Stderr, "process: %s\n", records[i].DataB64)
	}
	checkpointer.CheckpointAll()
	return nil
}

func (ec *EchoConsumer) Shutdown(shutdownType kinesis.ShutdownType,
	checkpointer *kinesis.Checkpointer) error {
	fmt.Fprintf(os.Stderr, "shutdown: %s\n", shutdownType)
	if shutdownType == kinesis.GracefulShutdown {
		checkpointer.CheckpointAll()
	}
	return nil
}

func main() {
	var ec EchoConsumer
	kinesis.Run(&ec)
}
