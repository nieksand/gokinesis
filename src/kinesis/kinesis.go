/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2014 Niek J. Sanders
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */
package kinesis

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// RecordConsumer interface must be implemented and then used via Run().  A Go
// KCL program may have only a single RecordConsumer launched.  The KCL daemon
// will spawn additional programs as needed.
type RecordConsumer interface {
	// Init is called before record processing with the shardId.
	Init(string) error

	// ProcessRecords is called for each batch of records to be processed.
	ProcessRecords([]*KclRecord, *Checkpointer) error

	// Shutdown is called before termination.
	Shutdown(ShutdownType, *Checkpointer) error
}

// ShutdownType indicates whether we have a graceful or zombie shutdown.
type ShutdownType int

const (
	unknownShutdown ShutdownType = iota

	// GracefulShutdown means checkpoint should be called.
	GracefulShutdown

	// ZombieShutdown means checkpoint may NOT be called.  Another record
	// processor may own the shard now.
	ZombieShutdown
)

// Run consumes from the local Kinesis KCL daemon.
func Run(c RecordConsumer) {

	checkpointer := &Checkpointer{true}
	for {
		// read next daemon request
		req := getAction()
		if req == nil {
			break
		}

		// dispatch based on action
		var err error
		switch {
		case req.Action == "processRecords":
			err = c.ProcessRecords(req.Records, checkpointer)

		case req.Action == "initialize":
			err = c.Init(*req.ShardID)

		case req.Action == "shutdown":
			shutdownType := GracefulShutdown
			if req.Reason == nil || *req.Reason == "ZOMBIE" {
				checkpointer.isAllowed = false
				shutdownType = ZombieShutdown
			}
			err = c.Shutdown(shutdownType, checkpointer)

		default:
			err = fmt.Errorf("Unsupported KCL action: %s", req.Action)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		// respond with ack
		fmt.Printf("\n{\"action\": \"status\", \"responseFor\": \"%s\"}\n", req.Action)
	}
}

// Checkpointer marks a consumers progress.
type Checkpointer struct {
	// isAllowed controls whether checkpointing triggers a panic.
	isAllowed bool
}

// CheckpointAll marks all consumed messages as processed.
func (cp *Checkpointer) CheckpointAll() {
	msg := "\n{\"action\": \"checkpoint\", \"checkpoint\": null}\n"
	cp.doCheckpoint(msg)
}

// CheckpointSeq marks messages up to sequence number as processed.
func (cp *Checkpointer) CheckpointSeq(seqNum int64) {
	msg := fmt.Sprintf("\n{\"action\": \"checkpoint\", \"checkpoint\": %d}\n", seqNum)
	cp.doCheckpoint(msg)
}

func (cp *Checkpointer) doCheckpoint(msg string) {
	if !cp.isAllowed {
		panic("attempted to checkpoint on ZOMBIE termination")
	}

	// send checkpoint req
	fmt.Print(msg)

	// receive checkpoint ack
	ack := getAction()

	if ack == nil {
		fmt.Fprintf(os.Stderr, "Received EOF rather than checkpoint ack\n")
		os.Exit(1)
	} else if ack.Action != "checkpoint" {
		fmt.Fprintf(os.Stderr, "Received invalid checkpoint ack: %s\n", ack.Action)
		os.Exit(1)
	} else if ack.Error != nil {
		fmt.Fprintf(os.Stderr, "Checkpoint error: %s\n", ack.Error)
		os.Exit(1)
	}
}

// KclRecord is an individual kinesis record.  Note that the body is always
// base64 encoded.
type KclRecord struct {
	DataB64        string `json:"data"`
	PartitionKey   string `json:"partitionKey"`
	SequenceNumber int64  `json:"sequenceNumber,string"`
}

// KclAction is a request from the local KCL daemon.
type KclAction struct {
	Action  string       `json:"action"`
	ShardID *string      `json:"shardId"`
	Records []*KclRecord `json:"records"`
	Reason  *string      `json:"reason"`
	Error   *string      `json:"error"`
}

var scanner = bufio.NewScanner(os.Stdin)

// getAction reads a request from the KCL daemon.
func getAction() *KclAction {
	for scanner.Scan() {
		// decode json action request
		line := scanner.Bytes()
		var req KclAction
		if err := json.Unmarshal(line, &req); err != nil {
			fmt.Fprintf(os.Stderr, "Could not understand line read from input: %s\n", line)
			continue
		}
		return &req
	}

	if err := scanner.Err(); err != nil {
		panic("unable to read stdin")
	}

	return nil
}
