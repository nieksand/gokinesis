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
func Run(c *RecordConsumer) {

	checkpointer := &Checkpointer{true}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		// decode json action request
		line := scanner.Bytes()
		var req KclAction
		if err := json.Unmarshal(line, &req); err != nil {
			fmt.Fprintf(os.Stderr, "Could not understand line read from input: %s", line)
			continue
		}

		// dispatch based on action
		var err error
		switch {
		case req.Action == "processRecords":
			err = (*c).ProcessRecords(req.Records, checkpointer)

		case req.Action == "initialize":
			err = (*c).Init(*req.ShardID)

		case req.Action == "shutdown":
			shutdownType := GracefulShutdown
			if req.Reason == nil || *req.Reason == "ZOMBIE" {
				checkpointer.isAllowed = false
				shutdownType = ZombieShutdown
			}
			err = (*c).Shutdown(shutdownType, checkpointer)

		default:
			err = fmt.Errorf("unsupported KCL action: %s", req.Action)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		// respond with ack
		fmt.Printf(`\n{"action": "status", "responseFor": "%s"}\n`, req.Action)
	}
	if err := scanner.Err(); err != nil {
		panic("unable to read stdin")
		fmt.Fprintf(os.Stderr, "unable to read stdin")
		os.Exit(1)
	}
}

type Checkpointer struct {
	// isAllowed controls whether checkpointing triggers a panic.
	isAllowed bool
}

func (cp *Checkpointer) CheckpointAll() {
	msg := `\n{"action": "checkpoint", "checkpoint": null}\n`
	cp.doCheckpoint(msg)
}

func (cp *Checkpointer) CheckpointSeq(seqNum int64) {
	msg := fmt.Sprintf(`\n{"action": "checkpoint", "checkpoint": %d}\n`, seqNum)
	cp.doCheckpoint(msg)
}

func (cp *Checkpointer) doCheckpoint(msg string) {
	if !cp.isAllowed {
		panic("attempted to checkpoint on ZOMBIE termination")
	}

	// send checkpoint req
	fmt.Print(msg)

	// receive checkpoint ack
	// FIXME
}

type KclRecord struct {
	DataB64        string `json:"data"`
	PartitionKey   string `json:"partitionKey"`
	SequenceNumber int64  `json:"sequenceNumber"`
}

type KclAction struct {
	Action  string       `json:"action"`
	ShardID *string      `json:"shardId"`
	Records []*KclRecord `json:"records"`
	Reason  *string      `json:"reason"`
	Error   *string      `json:"error"`
}
