# GMQ
#### _Golang Message Queue_

GMQ is an asynchronous message dispatcher, following the publish/subscribe queue architecture; a FIFO model regulates the queue.
Clients will either generate a message on the server queue (publish) or retrieve from the server the first element in the queue (subscribe).
If a queue doesn't exist the first message published on the queue will create it; if there are no messages in the queue, the server will return an empty message.

Messages are in JSON format, with literal attributes to identify the queue and the result of operation, plus a base64 encoded payload.
You decouple application layers by exchanging binary messages within your infrastructure deploying a server and a client on each componenent; this model is especially good for distributed systems.
At the moment a GMQ Server does not support distribution accross multiple nodes, but this is a definitive to-be implementation.

##### Deploy GMQ Server

First download the code and build it: Go 1.4 is recommended.
	go get github.com/inge4pres/gmq
	cd gmq/server
	go build gmqserver.go

Create a configuration file following example in gmq/test/configuration.
Start the server with
	./gmqserver -f /path/to/your/config/file

##### Deploy GMQ Client


##### _Security Notice_

GMQ server does not support TLS communication (at the moment) and does not use any form of  authentication mechanism, so be sure to have set up proper security with firewall and access to the server only from clients you trust.
