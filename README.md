# Penlook Socket
Real-time library for Go based on Long-Polling

[![Build Status](https://travis-ci.org/penlook/socket.svg?branch=master)](https://travis-ci.org/penlook/socket) [![GoDoc](https://godoc.org/github.com/penlook/socket?status.png)](https://godoc.org/github.com/penlook/socket) [![Software License](https://img.shields.io/badge/license-GNU-blue.svg?style=flat)](LICENSE.md) [![Author](http://img.shields.io/badge/author-penlook-red.svg?style=flat)](https://github.com/penlook)

# Documentation
Server
```go
import (
	. "github.com/penlook/socket"
)

socket := Socket {
	Port: 3000,
	Token: "acbz@3345123124567",
	Transport: LongPolling,
	Template: "asset/*",
}

socket.Initialize()
socket.Static("/static", "./asset")

socket.On("connection", func(client Client) {
	client.On("init", func(data Json) {
		// TODO
	})
})

socket.Listen()
```

Client
```javascript
socket = new Socket();
socket.connect();

socket.on('test', function(data) {
	socket.emit('test2', {
    	data : 'Package 2 from client'
	})
});

socket.emit('init', {
	data: 'Package from client'
})
```
