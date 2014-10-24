package kinesis

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	//"encoding/base64"
)


var respCheckpointSeq = `\n{"action": "checkpoint", "checkpoint": %d}\n`
var respActionDone    = `\n{"action": "status", "responseFor": "%s"}\n`




// {"action":"initialize","shardId":"shardId-123"}
// {"action":"processRecords","records":[{"data":"bWVvdw==","partitionKey":"cat","sequenceNumber":"456"}]}
// {"action":"shutdown","reason":"TERMINATE"}
type kclRecord struct {
	DataB64        string `json:"data"`
	PartitionKey   string `json:"partitionKey"`
	SequenceNumber int64  `json:"sequenceNumber"`
}

type kclAction struct {
	Action  string      `json:"action"`
	ShardID *string     `json:"shardId"`
	Records []kclRecord `json:"records"`
	Reason  *string     `json:"reason"`
	Error   *string     `json:"error"`
}

type Record struct {
	PartitionKey   string
	SequenceNumber int64
	Body           []byte
}

type RecordProcessor struct {
	RecordChan <-chan *Record
}

func NewRecordProcessor() *RecordProcessor {

	rc := make(chan *Record)
	go func() {
		// read lines from stdin
		stdinScanner := bufio.NewScanner(os.Stdin)
		for stdinScanner.Scan() {

			// decode from json
			line := stdinScanner.Bytes()
			var rec kclAction
			if err := json.Unmarshal(line, &rec); err != nil {
				fmt.Fprintf(os.Stderr, "Could not understand line read from input: %s", line)
				continue
			}

			// dispatch based on action
			if rec.Action == "processRecords" {
				for _ = range rec.Records {
				}
			} else if rec.Action == "initialize" {

			} else if rec.Action == "shutdown" {

			} else {
				panic(fmt.Sprintf("unsupported KCL action: %s", rec.Action))
			}
		}

		if err := stdinScanner.Err(); err != nil {
			panic("unable to read stdin")
		}
	}()

	return &RecordProcessor{RecordChan: rc}
}

// CheckpointAt marks the completed work at the indicated sequence number.
func (*rp RecordProcessor) CheckpointAt(seqNum int64) error {

	fmt.Printf(`\n{"action" : "checkpoint", "checkpoint" : %d}\n`, seqNum)

	// read stdin line
	// parse json




}
