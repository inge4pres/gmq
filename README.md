# GMQ
#### _Golang Message Queue_

GMQ is an asynchronous message dispatcher implementing the publish/subscribe queue architecture; a FIFO model regulates the queue.
Clients will either generate a message on the server queue (publish) or retrieve from the server the first element in the queue (subscribe).
If a queue doesn't exist the first message published on the queue will create it; if there are no messages in the queue, the server will return an empty message.

Messages are in JSON format, with literal attributes to identify the queue and the result of operation and a base64 encoded payload (the message to be stored in the queue).

You decouple application layers by exchanging binary messages within your infrastructure deploying a server and a client on each componenent; this model is especially good for distributed systems.
At the moment a GMQ Server does not support memory sinchronization accross multiple nodes, but this is a definitive to-be implementation.

GMQ should be configred to store messages in:
 - memory, very fast but unreliable for production: memory synchronization between multiple nodes not supported yet
 - database (MySQL), Postgres support to come, messages durability relies on database
 - filesystem, one file per queue will be created, the slowest but most reliable mode for production use

##### Deploy GMQ Server

Download the code and build it: Go 1.4 is recommended.
	go get github.com/inge4pres/gmq
	cd gmq/server
	go build gmqserver.go

Create a configuration file following example in gmq/test/configuration.
Start the server with
	./gmqserver -f /path/to/your/config/file

##### Deploy GMQ Client

Build the client
	cd gmq/client
	go build gmqclient.go

Start the client with parameters from command line (basic usage)
**Publish**
	./gmqclient -s gmqserver.mydomain.com:4567 -q myqueue -a P -m "base64messagetobesendtotheserver"
**Subscribe**
	./gmqclient -s gmqserver.mydomain.com:4567 -q myqueue -a S

##### _Security Notice_

GMQ server does not support TLS communication (at the moment) and does not use any authentication mechanism; be sure to have set up proper security with firewall and access to the server only from clients you trust.
