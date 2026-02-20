package handlers

import (
	"log"

	"TitanBot/internal/metrics"
	"github.com/bwmarrin/discordgo"
)

// Router handles concurrent event processing
type Router struct {
	session    *discordgo.Session
	workerSize int
	events     chan interface{}
}

func NewRouter(s *discordgo.Session, workers int) *Router {
	return &Router{
		session:    s,
		workerSize: workers,
		// Buffered channel to prevent blocking the Gateway
		events: make(chan interface{}, 10000),
	}
}

// StartWorkers spawns goroutines to handle events concurrently
func (r *Router) StartWorkers() {
	log.Printf("Starting %d event workers...", r.workerSize)
	for i := 0; i < r.workerSize; i++ {
		go func() {
			for event := range r.events {
				r.processEvent(event)
			}
		}()
	}
}

func (r *Router) OnReady(s *discordgo.Session, event *discordgo.Ready) {
	r.events <- event
}

func (r *Router) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	r.events <- m
}

// processEvent routes the event to the correct handler
func (r *Router) processEvent(event interface{}) {
	switch e := event.(type) {
	case *discordgo.Ready:
		log.Printf("Bot logged in as %s", e.User.String())
	case *discordgo.MessageCreate:
		if e.Author.ID == r.session.State.User.ID {
			return
		}

		metrics.MessagesReceived.Inc()

		if e.Content == "!ping" {
			r.session.ChannelMessageSend(e.ChannelID, "Pong!")
		}
	}
}
