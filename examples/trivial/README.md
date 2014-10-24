This is a trivial example using the Golang KCL.

On Init() it creates a file: ``/tmp/<shardId>.demo``

On ProcessRecords() it writes the base64-encoded message bodies to the file.  It
checkpoints after each completed batch of work.

On Shutdown() it verifies that we are not a zombie and calls a final checkpoint.
