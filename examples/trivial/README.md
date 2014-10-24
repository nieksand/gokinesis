Trivial Consumer Demo
=====================
This is a trivial example using the Golang KCL.

On Init() it creates a file: ``/tmp/<shardId>.demo``

On ProcessRecords() it writes the base64-encoded message bodies to the file.  It
checkpoints after each completed batch of work.

On Shutdown() it verifies that we are not a zombie and calls a final checkpoint.


## Usage
1. Ensure the Golang KCL is in your GOPATH.  
2. ``go build .``
3. complicated magic to invoke KCL MultiLangDaemon
