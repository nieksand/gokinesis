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
	outfile		*os.File
}

func (ec *EchoConsumer) Init(shardId string) error {
	var err error
	ec.outfile, err = os.Create( fmt.Sprintf("/tmp/%s.demo", shardId) )
	if err != nil {
		return err
	}

	fmt.Fprintf(ec.outfile, "init: %s\n", shardId)
	return nil
}

func (ec *EchoConsumer) ProcessRecords(records []*kinesis.KclRecord,
	checkpointer *kinesis.Checkpointer) error {
	for i := range records {
		fmt.Fprintf(ec.outfile, "process: %s\n", records[i].DataB64)
	}
	checkpointer.CheckpointAll()
	return nil
}

func (ec *EchoConsumer) Shutdown(shutdownType kinesis.ShutdownType,
	checkpointer *kinesis.Checkpointer) error {
	fmt.Fprintf(ec.outfile, "shutdown: %s\n", shutdownType)
	if shutdownType == kinesis.GracefulShutdown {
		checkpointer.CheckpointAll()
	}
	return nil
}

func main() {
	var ec EchoConsumer
	kinesis.Run(&ec)
}
