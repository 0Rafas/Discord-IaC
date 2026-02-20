package generator

const MainTemplate = `package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"{{.Name}}/internal/bot"
{{if eq .Database "postgres"}}	"{{.Name}}/internal/database"{{end}}
{{if eq .Cache "redis"}}	"{{.Name}}/internal/cache"{{end}}
{{if .Observability.Prometheus}}	"{{.Name}}/internal/metrics"{{end}}
)

func main() {
	log.Println("Starting {{.Name}}...")

	// Initialize discord client
	client, err := bot.NewClient()
	if err != nil {
		log.Fatalf("Error creating discord session: %v", err)
	}

{{if eq .Database "postgres"}}
	// Initialize Database
	err = database.Init()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
{{end}}

{{if eq .Cache "redis"}}
	// Initialize Redis Cache
	err = cache.Init()
	if err != nil {
		log.Fatalf("Cache initialization failed: %v", err)
	}
{{end}}

{{if .Observability.Prometheus}}
	// Start Prometheus Metrics Server
	go metrics.StartServer(":2112")
{{end}}

	// Open connection to discord
	err = client.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Shutting down cleanly...")
	client.Close()
}
`

const ClientTemplate = `package bot

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"{{.Name}}/internal/handlers"
)

func NewClient() (*discordgo.Session, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN environment variable is not set")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	// Register Intents
	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

{{if .Sharding}}
	// Sharding is enabled
	// In a real production bot, this would be dynamically allocated by a master broker
	// Here we set up local sharding configuration.
	session.ShardCount = 1
	session.ShardID = 0
{{end}}

	// Initialize the Event Router (Worker Pool)
	router := handlers.NewRouter(session, {{.Concurrency.WorkerPoolSize}})
	router.StartWorkers()

	// Register event handlers
	session.AddHandler(router.OnReady)
	session.AddHandler(router.OnMessageCreate)

	return session, nil
}
`

const HandlersTemplate = `package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
{{if .Observability.Prometheus}}	"{{.Name}}/internal/metrics"{{end}}
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
		events:     make(chan interface{}, 10000), 
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
{{if .Observability.Prometheus}}
		metrics.MessagesReceived.Inc()
{{end}}
		if e.Content == "{{.Prefix}}ping" {
			r.session.ChannelMessageSend(e.ChannelID, "Pong!")
		}
	}
}
`

const MetricsTemplate = `package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	MessagesReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bot_messages_received_total",
		Help: "The total number of messages received by the bot",
	})
)

// StartServer starts the prometheus metrics endpoint
func StartServer(addr string) {
	log.Printf("Starting Prometheus metrics server on http://localhost%s/metrics", addr)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Metrics server failed: %v", err)
	}
}
`

const GoModTemplate = `module {{.Name}}

go 1.22.0

require (
	github.com/bwmarrin/discordgo v0.27.1
{{if .Observability.Prometheus}}	github.com/prometheus/client_golang v1.19.0{{end}}
)
`

const DockerfileTemplate = `FROM golang:1.22-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bot .
CMD ["./bot"]
`

const DockerComposeTemplate = `version: '3.8'

services:
  bot:
    build: .
    environment:
      - DISCORD_TOKEN=${DISCORD_TOKEN}
{{if eq .Database "postgres"}}      - DB_URL=postgres://user:pass@db:5432/botdb{{end}}
{{if eq .Cache "redis"}}      - REDIS_URL=redis:6379{{end}}
{{if .Observability.Prometheus}}    ports:
      - "2112:2112"{{end}}
{{if eq .Database "postgres"}}
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: botdb
    ports:
      - "5432:5432"
{{end}}
{{if eq .Cache "redis"}}
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
{{end}}
{{if .Observability.Prometheus}}
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"
{{end}}
`

const DatabaseTemplate = `package database

import "log"

// Init initializes the database connection
func Init() error {
	log.Println("Database connection initialized")
	return nil
}
`

const CacheTemplate = `package cache

import "log"

// Init initializes the cache connection
func Init() error {
	log.Println("Cache connection initialized")
	return nil
}
`

const GithubActionsTemplate = `name: Discord Bot CI

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.22.0'

    - name: Build
      run: go build -v ./...

    - name: Test
      run: go test -v ./...
`

const K8sDeploymentTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}-bot
  labels:
    app: {{.Name}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.Name}}
  template:
    metadata:
      labels:
        app: {{.Name}}
{{if .Observability.Prometheus}}      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "2112"{{end}}
    spec:
      containers:
      - name: bot
        image: {{.Name}}-bot:latest # Replace with your registry
        env:
        - name: DISCORD_TOKEN
          valueFrom:
            secretKeyRef:
              name: bot-secrets
              key: discord-token
`
