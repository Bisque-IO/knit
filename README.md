# knit

A fault-tolerant service for vertx clustering using
[knit](https://github.com/bisque-io/knit) framework



## Commands

```
LOCK ACQUIRE
LOCK UPDATE
LOCK RELEASE
MULTIMAP GET
MULTIMAP PUT
MULTIMAP REMOVE
MAP GET
MAP PUT
MAP REMOVE
PING
```

## Building

Using the source file from the examples directory, we'll build an application
named "knit"

```
go build -o knit main.go
```

## Running

It's ideal to have three, five, or seven nodes in your cluster.

Let's create the first node.

```
./knit -n 1 -a :11001
```

This will create a node named 1 and bind the address to :11001

Now let's create two more nodes and add them to the cluster.

```
./knit -n 2 -a :11002 -j :11001
./knit -n 3 -a :11003 -j :11001
```

Now we have a fault-tolerant three node cluster up and running.

### Using

You can use any Redis compatible client, such as the redis-cli, telnet,
or netcat.

I'll use the redis-cli in the example below.

Connect to the leader. This will probably be the first node you created.

```
redis-cli -p 11001
```

Send the server a knit command and receive the first knit.

```
> RAFT HELP
1) RAFT LEADER
2) RAFT INFO [pattern]
3) RAFT SERVER LIST
4) RAFT SERVER ADD id address
5) RAFT SERVER REMOVE id
6) RAFT SNAPSHOT NOW
7) RAFT SNAPSHOT LIST
8) RAFT SNAPSHOT FILE id
9) RAFT SNAPSHOT READ id [RANGE start end]
```


For other information check out the [knit README](https://github.com/tidwall/knit).