## jakob

Jakob is fault-tolerant, distributed cluster of [Redis](http://redis.io) servers with built-in load-balancing and fall-backs 
to provide data availability. Jakob is specifically meant to start a cluster of [Tile38](http://tile38.com) servers to store 
geo-spatial data. 

Jakob relies on [Apache Kafka](ttps://kafka.apache.org/) to store the logs and for log replication. It also relies on
the amazing [Machinery](github.com/RichardKnop/machinery/)  to sync data(logs) between servers in background. 
Jakob's Machinery setup uses [RabbitMQ](http://rabbitmq.com) as broker and [Redis](http://redis.io) as result-backend.

Jakob has two types of servers - `setters` and `getters`.

A setter server will receive all the Tile38 setter commands like (`SET`, `NEARBY`, `FENCE`, etc), while a getter server will always receive
Tile38 getter commands like (`GET`, `MATCH`, etc). 

The two clusters are servers are arranged in a consistent-hashed ring. A setter or getter server is selected using consistent hashing.


It just exposes two HTTP endpoints. Both of them are POST HTTP endpoints.

| Endpoint | Purpose | Data |
| ------- |:--------:|:-------------------------------------------:|
| `/init` | initializes the cluster with a setter and getter peer | e.g., `{"setter": "setter:8080", "getter":"getter:8081"}`
| `/join` | joins the cluster with a setter and getter peer | e.g., `{"setter": "setter:9000", "getter":"getter:9001"}`

Installation:

`go get -u github.com/younisshah/jakob`

or using `glide`

`glide get github.com/younisshah/jakob`

and then `glide up`

Next install and start Apache Kafka. Refer to this for a quick start. [Apache Kafka Quick Start](https://kafka.apache.org/quickstart).
Download and install RabbitMQ from here [RabbitMQ](http://www.rabbitmq.com/download.html).
Finally download and start Redis server from here [Redis](https://redis.io/download).


__Sample usage:__

Boot jakob:

```go
boot.Up()
```

Initialize the cluster, use `/init`

`curl -X POST localhost:30000/init -d '{"setter":"localhost:9000", "getter":"localhost:9001"}'`

To join a cluster, use `/join`
 
`curl -X POST localhost:30000/join -d '{"setter":"localhost:9002", "getter":"localhost:9003"}'` 

To get a setter peer from the cluster

```go
ring := jring.New(jfs.SETTER)
peer, err := ring.Get("some_id")
if err != nil {
    log.Println(err)
}
```

To get getter peer from the cluster

```go
ring := jring.New(jfs.GETTER)
peer, err := ring.Get("some_id")
if err != nil {
    log.Println(err)
}
```

To send a Tile38 command to jakob

```go
cmd = command.NewCommand(peer, "SET", "my", "home", "POINT", "12.4", "23.45")
cmd.Execute()
if cmd.Error != nil {
    log.Println(cmd.Error)
} else {
    log.Println(cmd.Result)
}
```

---

**Jakob is still in infancy which means API is subject change.**

---

TODO

> more comprehensive documentation
