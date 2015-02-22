package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
	"fmt"
)

func TestSocket(t *testing.T) {

	assert := assert.New(t)
	assert.Equal("Test", "Test")

	socket := Socket {
		Port: 3000,
		Token: "acbz@3345123124567",
		Transport: LongPolling,
		Template: "asset/*",
	}

	fmt.Println("Initialize")
	socket.Initialize()

	socket.Static("/static", "./asset")
	socket.Router.GET("/", func(context *gin.Context) {
		context.HTML(200, "index.html", Json {})
	})

	fmt.Println("On connection")
	socket.On("connection", func(client Client) {
		fmt.Println("Emit abc")
		client.Emit("abc", Json {
			"key1" : "value1",
			"key2" : "value2",
			"key3" : "value3",
		})

		client.Emit("abc", Json {
			"key1" : "value1",
			"key2" : "value2",
			"key3" : "value3",
		})

		client.Emit("abc", Json {
			"key1" : "value1",
			"key2" : "value2",
			"key3" : "value3",
		})

		fmt.Println("Client on")
		client.On("test", func(client Client) {
			fmt.Println("Enter test event")
			client.Emit("test", Json {
				"abce" : "Hello",
				"abcf" : "Yes",
			})

			client.On("test_2", func(client Client) {
				client.Emit("abc", Json {
					"test" : "1234",
				})
				client.On("test_3", func(client Client) {
					client.On("test_4", func(client Client) {
						client.Emit("abc", Json {
							"test" : "1234",
						})
					})
					client.Emit("abc", Json {
						"test" : "test",
					})
				})
				client.On("test_5", func(client Client) {
					client.Emit("abc", Json {
						"abc3" : "234556",
					})
					client.On("test7", func(client Client) {
						client.On("test8", func(client Client) {
							client.Emit("abc", Json {
								"abc" : "avc",
							})
						})
					})
				})
				client.Emit("abc", Json {
					"abc" : "abc",
				})
			})
			client.On("test6", func(client Client) {
				client.Emit("abcer", Json {
					"abc" : "1245667",
				})
			})
		})
	})

	socket.Listen()
}


