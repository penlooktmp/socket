package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "strconv"
    "github.com/oleiade/lane"
    //"container/list"
    "net/http"
    //"fmt"
    //"time"
)

const LongPolling int = 0

type Json map[string] interface {}

type Context struct {
    Context *gin.Context
    Channel chan Context
    Output chan Json
    Handshake string
    Polling bool
}

type Socket struct {
    Port int
    Token string
    Transport int
 	Event map[string] func(client Client)
    Clients map[string] Client
    Context chan Context
    Output chan Json
    Router *gin.Engine
    Template string
}

func (socket *Socket) Initialize() Socket {

    // Route
    gin.SetMode(gin.DebugMode)
    socket.Router = gin.Default()

    // Context
    socket.Context = make(chan Context, 10)

 	// Event
    socket.Event = make(map[string] func(client Client))

    // Clients
    socket.Clients = make(map[string] Client)

    // Output
    socket.Output = make(chan Json, 100)

    // Socket template
    socket.Router.LoadHTMLGlob(socket.Template)

    return *socket
}

func (socket Socket) Debug(message string) {
    fmt.Println(message)
}

func (socket Socket) ParseContext(context Context, callback func(client Client)) Client {

	client := socket.Clients[context.Handshake]
    client.Context = context.Context

	if ! context.Polling {
        fmt.Println("Emit connection")
        socket.Emit("connection", Json {
	        "handshake" : context.Handshake,
	    })
        callback(client)
	    return client
	}

	if data := <- client.Output ; data != nil {
        client.Output <- data
    }

    return client
}

func (socket Socket) On(event string, callback func(client Client)) {
	socket.Event[event] = callback
}

func (socket Socket) Emit(event string, data Json) {
    fmt.Println("Socket emit handshake")
	socket.Output <- Json {
    	"event": event,
    	"data" : data,
    }
    fmt.Println("Done emit")
}

func (socket Socket) Static(route string, directory string) Socket {
    socket.Router.Static(route, directory)
    return socket
}

// Check polling request per connection
func (socket Socket) LoopSocketEvent(context Context) {
    fmt.Println("Loop Socket Event")
	go func(callback func(client Client), context Context) {
		for {
			select {
				case context := <- context.Channel:
                    fmt.Println("Parse Context")
					socket.ParseContext(context, callback)
			}
		}
	} (socket.Event["connection"], context)
}

func (socket Socket) RegisterClientEvent(event Event) {

}

func (socket Socket) LoopClientEvent(context Context) {

	client := socket.Clients[context.Handshake]
    event  := client.Event
    fmt.Println("Node List")
    fmt.Println(event.Pop())
    fmt.Println("Done")

    // Scanning linked list
    //events := client.Events
}

func (socket Socket) GetConnection(context *gin.Context) Context {

	handshake := random()
	output    := make(chan Json, 10)
	channel   := make(chan Context, 10)
	event     := lane.NewStack()

	client := Client {
       	Context: context,
       	Output : output,
       	Channel: channel,
        Event: event,
        MaxNode: 0,
    }

    socket.Clients[handshake] = client

	return Context {
		Context   : context,
		Output    : output,
		Channel   : channel,
		Handshake : handshake,
		Polling   : false,
	}
}

func (socket Socket) GetPolling(context *gin.Context) Context {
    fmt.Println("Get Polling")
	handshake := context.Params.ByName("handshake")

	client := socket.Clients[handshake]
	client.Context = context

    fmt.Println("Return context")

	return Context {
		Context : context,
		Channel : client.Channel,
		Output : client.Output,
		Handshake: handshake,
		Polling: true,
	}
}

func (socket Socket) Response(context Context) {
	select {
        case data := <- socket.Output:
            fmt.Println("Response socket output")
            context.Context.JSON(200, data)
        case data := <- context.Output:
            fmt.Println("Recieve output")
            context.Context.JSON(200, data)
	}
}

func (socket Socket) Listen() Socket {

    socket.Router.GET("/polling", func(_context *gin.Context) {
        fmt.Println("Start connection")
    	context := socket.GetConnection(_context)
    	socket.LoopSocketEvent(context)
    	context.Channel <- context
    	socket.Response(context)
    })

    socket.Router.GET("/polling/:handshake", func(_context *gin.Context) {
    	context := socket.GetPolling(_context)
        fmt.Println(context)
        fmt.Println("Loop Client Event")
    	socket.LoopClientEvent(context)
        fmt.Println("Push to context channel")
    	context.Channel <- context
        fmt.Println("Response")
    	socket.Response(context)
    })

    http.ListenAndServe(":" + strconv.Itoa(socket.Port), socket.Router)
    return socket
}
