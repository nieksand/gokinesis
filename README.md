gokinesis
=========
A Golang client library for AWS Kinesis.

The Golang KCL supports everything needed to write Kinesis stream consumers.  It
offers the same functionality as the official Python KCL:

    http://aws.amazon.com/blogs/aws/speak-to-kinesis-in-python/

This is an unofficial KCL created outside of AWS.


## Current State
The Golang KCL has been tested for basic functionality.  It can consume records,
checkpoint, and handle shard creation and merges.

However, it has not been used in any production applications.  Feedback and bug
reports can be submitted through Github Issues and are greatly appreciated.


## Quickstart
To build the example:

    cd gokinesis
    export GOPATH=`pwd`
    cd examples/trivial
    go build .

See examples/trivial/README.md for instructions on how to invoke the KCL
MultiLangDaemon on the trivial consumer demo.


## Usage
To write a consumer, implement a struct satisfying the kinesis.RecordConsumer
interface.

1. ``Init(shardId string)`` is called at start.
2. ``ProcessRecords(batch []*KclRecord, cp *Checkpointer)`` is invoked for each
   batch of messages to be processed.
3. ``Shutdown(t ShutdownType, cp *Checkpointer)`` is called at finish.  No
   further calls to the RecordConsumer are made after this.  Implementations
   must check the ShutdownType.  Invoke a checkpoint if it's graceful; do NOT
   checkpoint if it's ZOMBIE.

If the implementation returns an non-nil error on any of these three methods,
the Golang KCL will exit immediately with a non-zero exit code.

In the Checkpointer, CheckpointAll() marks to the end of the most recently
delivered batch of records.  CheckpointSeq() marks to a specific sequence
number.  Checkpointing returns a non-nil error on conditions which can be
retried.

To begin processing, instantiate the struct and pass it to kinesis.Run().

An example of this is shown in examples/trivial/trivial.go

**Note:** there should be only a single RecordProcessor and invocation of
kinesis.Run().  The MultiLangDaemon protocol uses stdin/stdout to communicate
and does not support multiple consumers in one process.  Instead, it spawns
multiple independent processes which are each assigned a single shard.

**Note:** the KCL MultiLangDaemon uses stdin and stdout to communicate with the
client.  Any consumer output should go elsewhere (e.g. stderr or a file).


## Design
The KCL MultiLangDaemon protocol restricts when checkpointing may happen.
Accordingly, an interface and controller function rather than a Go channel is
used for processing records.


## License
The MIT License (MIT)

Copyright (c) 2014 Niek J. Sanders

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
