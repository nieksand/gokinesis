Trivial Consumer Demo
=====================
This is a trivial example using the Golang KCL.

On Init() it creates a file: ``/tmp/<shardId>.demo``

On ProcessRecords() it writes the base64-encoded message bodies to the file.  It
checkpoints after each completed batch of work.

On Shutdown() it verifies that we are not a zombie and calls a final checkpoint.


## Usage
Install the KCL MultiLangDaemon using the instructions from the official AWS
Python KCL readme.

* Ensure the Golang KCL is in your GOPATH.  

* ``go build .``

* From the Python KCL samples directory, use the helper script to generate the
  launch command.  Be sure to set this ``--java`` path properly.

    python amazon_kclpy_helper.py --print_command --java /usr/bin/java --properties ./sample.properties

* Copy/paste the output to execute..
