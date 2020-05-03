# goloca
Handling close to 15,000 requests per minute per server

## Design
Goloca is split into 4 main packages - 
* _main_ - Initializes the app and registers the endpoints
* _db_ - Exposes a _Store_ interface which is used to implement the database functionality
* _pool_ - This package handles all the concurrency-related aspects. It's where the goroutines dance.
* _api_ - This package contains the API handlers behind the endpoints.

I'll describe each package in detail.

---

## Package db
This package implements a `Store` interface like this - 

```go
type Store interface {
	Connect()
	Config() interface{}
	Insert(interface{}) error
}
```
The reason for using an interface is obvious. If someone later wants to change the connection logic or use a different
database altogether, it can be done easily by implementing this interface.

Currently, this interface is implemented by the`PostgresConn` type. The `Connect` method on PostgresConn will try to
initialize the database connection (using credentials from the environment), and if it fails, it'll make 9 more attempts at 10 second intervals. If the database connection is not established within these 10 attempts, the app will exit. Upon connection, it will create the required table in the database if it doesn't already exist.

The `Insert` method on `PostgresConn` executes a simple SQL query to insert data into the table. Inside the application, a data point is stored as the type `RideDataPoint`, as implemented in _db/types.go_. While the `Insert` method accepts anything as an argument, in the `PostgresConn` implementation, the argument must be the `RideDataPoint` type - passing anything else will result in an error.

I did not use any external sql library like sqlx, in order to keep things light and simple.

## Package pool
This package is responsible for spawning and assigning jobs to workers. A worker is a simple goroutine listening to a universal job channel. In Java speak, this is a thread-pool analogue. This is the definition of a pool - 

```go
type Pool struct {
	NumWorkers int
	JobChannel *(chan Job)
	Workers    []*worker
}
```
NumWorkers is the number of workers that this pool will spawn. JobChannel is a channel down which the pool can pass `Job`s. All the spawned workers are 'listening' to this channel, and any `Job` passed down the JobChannel of a queue would be taken up by a free worker as soon as possible. Like I said - its thread-pool analogue.

Speaking of `Job`s, `Job` is a type defined like this - 

```go
type Job struct {
	DataPoint  interface{}
	StatusChan chan interface{}
}
```

`DataPoint` can be anything, but in our application it's restricted to only the `RideDataPoint` type (I didn't give it a fixed type to keep the **pool** package reusable in the future). `StatusChan` is a channel on which you can listen for the status of your job. On completion of the job, the worker sends `true` down that `Job`'s `StatusChan`. If an error arises, that error can also be sent down the `Job`'s `StatusChan`. Keeping things like this allows the goroutine that dispatched the job to the pool to track whats going on with the job. 

In this application, we use this channel in our API handlers in the **api** package to know when to close the http Response and return it.

The `pool` package is dynamic enough for you to be able to spawn any amount of workers to execute any kind of function with any kind of data, while also allowing the user to keep track of what the worker is doing. Thanks to this interview task, I now have a pretty sweet threadpooling library that I'll use in my other go projects.

## Package api
This package is pretty standard. I define two handlers here (using the famous return-the-handler-as-a-closure approach). I used dependency injection to give the handler what it needs from the global app config. Using dependency injection allows me to keep track of how different dependencies are getting used in my code and it will also come in handy if I ever decide to write tests for my code (since I'll be able to mock things easier).

Both the handlers have the same purpose - to allow the client to send some data to add to the database. However, one handler uses a standard HTTP connection while the other uses a WebSocket. The reason for having a WebSocket handler is that during peak traffic hours, the client would be able to connect to my service using the socket and then pass messages to me - this would mean that a new connection need not be established everytime the client sends data, which will shed off some latency. The other handler is a simple HTTP handler which accepts post requests containing the data to be added, validates that data and adds it to the database using the `pool`.

To summarize - 
1. Use a WebSocket during high traffic to save time.
2. Use the HTTP endpoint during low traffic to save resources.

Both the handlers use the `pool`'s `Dispatch()` method to create a job for adding the data recieved to the database.

## Package main
This is where I set up the routes and initialize the app. In the file _app.go_, I load the environment from an `env` file, create a new database connection pool using the credentials stored in the env file, start the `pool` and setup an http upgrader (for websocket functionality). I store all this in a `Config` object and use it to pass dependencies down to my handlers.

---

# Deployment
I've written a Dockerfile to deploy the app as a container, and a docker-compose.yaml file to deploy it on a server with a single command. For the Dockerfile, I used Ubuntu as the base image, since for some reason the GoLang image was too big (I thought of a multi-stage build, but it just seemed easier to set things up with ubuntu instead). With this Dockerfile, this app can be deployed on Google AppEngine Flex or AWS Elastic Beanstalk. The docker-compose file can be used to deploy it on a VPS.

Please note that even after bringing up the service with docker-compose, you'd still need to setup the webserver to make the app properly accessible from the internet. A simple nginx reverse-proxy would do the job.

When deploying with docker-compose, the `DATABASE_HOST` property in config.env should be changed to 'database'.

---
# Benchmarks
On my machine (A Dell Inspiron 5577 with 8GB ram and a quad-core Intel Core i5), I spawned 16 workers. With this configuration - 
* It takes an average of 3.5 seconds to complete 1000 requests with a concurrency of 50. 
* With a concurrency of 1, the above takes an average of 24 seconds.
* A single request took approximately 20 ms to complete.

With this kind of a config, a single server can easily handle _an order of magnitude_ more requests than 200 per minute.

---
# System Architecture
This service can be easily scaled horizontally, if the need arises. Depending on the need to read from the database, we can set up a Postgres cluster to keep the database highly available for any services that depend on it. 

To scale the application in real time, we can spin up a load balancer that sits in front of a cluster of servers hosting this service. When the traffic increases so much that the last server is at 90% capacity for more than a minute, we can simply add another server to the cluster.

This can be done with Kubernetes, both on premise or with cloud provider like Google Cloud, AWS or Azure.

For running realtime ML models on the data recieved, we can setup a Kafka Cluster that streams the data recieved by our server to any clients that subscribe to it. This cluster can also be scaled horizontally when traffic increases.

--- 
# Please note
I have written this with the understanding that I was supposed to write a service that can handle high traffic. For that reason, I have excluded trivial functionality like authentication from the service. While including authentication would slow down the HTTP requests by a bit (given that a database call would have to be made for each HTTP request), it definitely won't slow down the service to less than a 1000 requests per minute. This would still meet the requirements.

Furthermore, if the client uses WebSockets instead of HTTP request, the service wouldn't need to verify it's identity every time it sent data and the server would still be able to manage around 15,000 hits per minute in that case.

I have also made no kind of assumptions related to the constraints on the data that I'm receiving, since it doesn't make much difference in how much traffic the application can manage. The database schema is quite simple.

I have included the config.env file in the repository. This is safe because this is just an interview task and not some actual project. There are no secrets in config.env.
