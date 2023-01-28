package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"playful/app/pkg/controller"
	"playful/app/pkg/service"
	"playful/app/tools/http"

	"github.com/gin-gonic/gin"
)

type httpController struct {
	config  http.Config
	router  *gin.Engine
	service service.PlayfulService
}

// It keeps a list of clients those are currently attached
// and broadcasting events to those clients.
type Event struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

// New event messages are broadcast to all registered client connection channels
type ClientChan chan string

func (c *httpController) Manager() chan error {

	errChan := make(chan error)
	// Initialize new streaming server
	stream := NewServer()

	// We are streaming data to clients in the interval 1 seconds
	go func() {
		for {
			time.Sleep(time.Millisecond * 15)

			l, err := c.service.GetLoaction(context.TODO())
			if err != nil {
				log.Printf("failed to get location: %v", err)
				continue
			}

			sl, err := json.Marshal(&l)
			if err != nil {
				log.Printf("failed to marshalize location: %v", err)
				continue
			}

			// Send to clients message channel
			stream.Message <- string(sl)
		}
	}()

	go func() {

		c.router.GET("/loc", HeadersMiddleware(), stream.serveHTTP(), func(c *gin.Context) {
			v, ok := c.Get("clientChan")
			if !ok {
				return
			}
			clientChan, ok := v.(ClientChan)
			if !ok {
				return
			}
			c.Stream(func(w io.Writer) bool {
				// Stream message to client from message channel
				if msg, ok := <-clientChan; ok {
					c.SSEvent("message", msg)
					return true
				}
				return false
			})
		})

		err := c.router.Run(fmt.Sprintf(":%s", c.config.Port))
		errChan <- err
	}()

	return errChan
}

// Initialize event and Start procnteessing requests
func NewServer() (event *Event) {
	event = &Event{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go event.listen()

	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *Event) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (stream *Event) serveHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize client channel
		clientChan := make(ClientChan)

		// Send new connection to event server
		stream.NewClients <- clientChan

		defer func() {
			// Send closed connection to event server
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}

func corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

func NewHttpController(c http.Config, s service.PlayfulService) controller.HttpController {

	r := gin.Default()
	r.Use(corsMiddleware)

	return &httpController{
		config:  c,
		service: s,
		router:  r,
	}
}
