#		Pub-Sub using Redis

## What is Pub-Sub?

Pub/Sub is a messaging paradigm that consists of defining Publishers and Subscribers which Channels in between them, where Publishers act as “message senders” and Subscribers act as “message receivers”. The Publishers do not send messages directly to Subscribers but rather to Channels; those Channels act as an intermediary between Publishers and Subscribers, the idea is to have those Subscribers to only receive the messages they are interested to, and having them decoupled from the Publishers.


## Redis
Redis (Remote Dictionary Server) is an in-memory data structure store, used as a distributed, in-memory key–value database, cache and message broker, with optional durability. Redis supports different kinds of abstract data structures, such as strings, lists, maps, sets, sorted sets, HyperLogLogs, bitmaps, streams, and spatial indices. It is an open source (BSD licensed), in-memory data structure store, used as a database, cache, and message broker. Redis supports Pub-Sub paradigm and defines different commands for publish, subscribing and unsubscribing; it even supports pattern-matching subscriptions.


<p align="center">
  <img src="https://miro.medium.com/max/786/1*9ryv3g78fC0ocNgae2HPzg.png">
Pub-Sub using Redis 
</p>



## Basic Pub-Sub Application using Redis and golang

### Redis Docker-compose set up
Below is the docker-compose file we have used to setup our redis:

    version: '3.6'
    services:
        redis:
            image: 'bitnami/redis:latest'
            ports:
                - "6379:6379"
            environment:
                - ALLOW_EMPTY_PASSWORD=yes


In this docker-compose file the redis image used is:

    bitnami/redis:latest

Now that we have set up out docker-compose file we need to establish  to redis with our server using golang

For this we would be using golang library:

    github.com/gomodule/redigo/redis

To implement pubsub using redis we need to a publisher and consumer object which can connect to redis and send and recieve data


    var (
        RedisPub redis.Conn
        psc      redis.PubSubConn
    )

    func NewRedisProvider() {
        var err error
        RedisPub, err = redis.Dial("tcp", "localhost:6379")
        if err != nil {
            fmt.Printf("Unable to connect to redis client")
        }
        RedisSub, err := redis.Dial("tcp", "localhost:6379")
        if err != nil {
            fmt.Printf("Unable to connect to redis client")
        }
        psc = redis.PubSubConn{Conn: RedisSub}
        err = psc.Subscribe("order")
        if err != nil {
            fmt.Printf("Unable to subscribe to order channel")
        }
    }


Below implemented the publish method for our pubsub application:

    func Publish(key string, value interface{}) error {
        _, err := RedisPub.Do("PUBLISH", key, val)
        if err != nil {
            fmt.Printf("Unable to publish to order channel")
            return err
        }
        return nil
    }

Once our publisher is ready, our second method would be to create a Consumer to receive messages:

    func Consume() error {
        for {
            switch v := psc.Receive().(type) {
            case redis.Message:
                fmt.Printf("Recieved: %v", string(v.Data))
            case error:
                return v
            }
        }
        return nil
    }

Below is our ``` main()```method

    func main() {
        NewRedisProvider()
        time.Sleep(time.Second * 2)
        Publish("order", "Hello")
        go Consume()
        time.Sleep(time.Minute * 2)
    }

Below is the output of above-mentioned code:

    GOROOT=/home/manubhav/sdk/go1.19 #gosetup
    GOPATH=/home/manubhav/go #gosetup
    /home/manubhav/sdk/go1.19/bin/go build -o /tmp/GoLand/___go_build_main_go /home/manubhav/GolandProjects/redisPubSub/main.go #gosetup
    /tmp/GoLand/___go_build_main_go
    
    Recieved: Hello
    
    Process finished with the exit code 0



## Strengths of PubSub Pattern:

Let's discuss some advantages of the Pub/Sub Pattern:
- Loose Coupling Between System Components
- Better View of the System-wide Workflow
- Enables Better & Faster Integration
- Ensures Smoother Scalability
- Guaranteed Consistent Reliability
- Builds Elasticity
- Software Modularization
- Language Agnostic Software Development
- The clarity in Business Logic
- Improves Responsiveness

## Conclusion
In this blog, we learned about the Publish/Subscribe design pattern. And explored how the Redis pub/sub works. We also explored what are the best use cases of Redis pub/sub, real-time messaging. 


