package socket

import (
    "container/list"
    "github.com/gin-gonic/gin"
    "strconv"
    "net/http"
    "encoding/json"
)

const LongPolling int = 0

// Request context structure
type Context struct {
    Context *gin.Context
    Channel chan Context
    Output chan Json
    Handshake string
    Polling bool
}

// Socket structure
type Socket struct {
    Port int
    Token string
    Transport int
 	Event map[string] func(client Client)
    Clients map[string] Client
    Context chan Context
    Router *gin.Engine
    Template string
}

// Initiazlie for socket
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

    // Socket template
    socket.Router.LoadHTMLGlob(socket.Template)

    return *socket
}

// Update new context for new or exist client per handshake
func (socket Socket) UpdateContext(context Context) Client {

	client := socket.Clients[context.Handshake]
    client.Context = context.Context

	if ! context.Polling {
        context.Context.JSON(200, Json {
            "event" : "connection",
            "data"  : Json {
                "handshake" : context.Handshake,
            },
        })
	}

    return client
}

// Socket listen client event
func (socket Socket) On(event string, callback func(client Client)) {
	socket.Event[event] = callback
}

// Static resources
func (socket Socket) Static(route string, directory string) Socket {
    socket.Router.Static(route, directory)
    return socket
}

// Check polling request per connection
func (socket Socket) LoopSocketEvent(context Context) {
    socket.UpdateContext(context)
}

// Scan client events in the first handshake
func (socket Socket) InitClientEvent(context Context) {
    client := socket.Clients[context.Handshake]

    if ! client.HandshakeFlag {
        client.HandshakeFlag = true
        socket.Clients[context.Handshake] = client
        callback := socket.Event["connection"]
        callback(client)
    }
}

// Submit emit package from client
func (socket Socket) SubmitClientEvent(context Context) {

    decoder := json.NewDecoder(context.Context.Request.Body)
    var pkg Json
    decoder.Decode(&pkg)
    event_name := pkg["event"]
    client := socket.Clients[context.Handshake]

    var event Event
    for cursor := client.Event.Front(); cursor != nil; cursor = cursor.Next() {
        event = cursor.Value.(Event)
        if event.Name == event_name {
            event.Callback(pkg)
        }
    }

    context.Context.JSON(200, Json {
         "status" : "OK",
    })
}

// Get fresh connection
func (socket Socket) GetConnection(context *gin.Context) Context {

	handshake := random()
	output    := make(chan Json, 10)
	channel   := make(chan Context, 10)
	event     := list.New()

	client := Client {
        Socket: socket,
       	Context: context,
       	Output : output,
       	Channel: channel,
        Event: event,
        Handshake: handshake,
        HandshakeFlag: false,
        MaxEvent: 0,
    }

    socket.Clients[handshake] = client

	return Context {
		Context   : client.Context,
		Output    : client.Output,
		Channel   : client.Channel,
		Handshake : handshake,
		Polling   : false,
	}
}

// Get polling connection
func (socket Socket) GetPolling(context *gin.Context) Context {
	handshake := context.Params.ByName("handshake")
	client := socket.Clients[handshake]
	client.Context = context

	return Context {
		Context : context,
		Channel : client.Channel,
		Output : client.Output,
		Handshake: handshake,
		Polling: true,
	}
}

// Waiting for long-polling response
func (socket Socket) Response(context Context) {
    select {
        case data := <- context.Output:
            context.Context.JSON(200, data)
    }
}

// Listen socket
func (socket Socket) Listen() Socket {

    socket.Router.GET("/polling", func(_context *gin.Context) {
        context := socket.GetConnection(_context)
        socket.LoopSocketEvent(context)
    })

    socket.Router.GET("/polling/:handshake", func(_context *gin.Context) {
        handshake := _context.Params.ByName("handshake")
        if handshake == "" {
            return
        }
        context := socket.GetPolling(_context)
        socket.InitClientEvent(context)
	    context.Channel <- context
	    socket.Response(context)
    })

    socket.Router.POST("/polling/:handshake", func(_context *gin.Context) {
        handshake := _context.Params.ByName("handshake")
        if handshake == "" {
            return
        }
        context := socket.GetPolling(_context)
        socket.SubmitClientEvent(context)
        context.Channel <- context
    })

    http.ListenAndServe(":" + strconv.Itoa(socket.Port), socket.Router)
    return socket
}