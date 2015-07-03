# GMQ
#### _Golang Message Queue_

GMQ is an asynchronous message dispatcher, following the publish/sbubscribe queue architecture; a FIFO model regulates the queue.
Clients will either generate a message on the server queue (publish) or retrieve from the server the first element in the queue (subscribe).
If a queue doesn't exist the first message published on the queue will create it; if there are no messages in the queue, the server will return an empty message.

Messages are in JSON format, with literal attributes to identify the queue, and the result of operation and a base64 encoded payload.
You decouple application layers by exchanging binary messages within your infrastructure deploying a server and a client on each componenent; this model is especially good for distributed systems.
At the moment a GMQ Server does not support distribution accross multiple nodes, but this is a definitive to-be.
