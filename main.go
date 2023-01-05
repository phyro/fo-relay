package main

import (
	"encoding/json"
	"log"

	"github.com/fiatjaf/relayer"
	"github.com/nbd-wtf/go-nostr"
)

type DoNothingStore struct{}

// Implement relayer's Storage interface
func (d *DoNothingStore) Init() error                                { return nil }
func (d *DoNothingStore) DeleteEvent(id string, pubkey string) error { return nil }
func (d *DoNothingStore) SaveEvent(event *nostr.Event) error         { return nil }
func (d *DoNothingStore) QueryEvents(filter *nostr.Filter) ([]nostr.Event, error) {
	return []nostr.Event{}, nil
}

type Relay struct{}

// Implement relays's Relay interface
func (r *Relay) Name() string                  { return "ForwardOnlyRelay" }
func (r *Relay) Storage() relayer.Storage      { return &DoNothingStore{} }
func (r *Relay) OnInitialized(*relayer.Server) {}
func (r *Relay) Init() error                   { return nil }
func (r *Relay) BeforeSave(evt *nostr.Event)   {}
func (r *Relay) AfterSave(evt *nostr.Event)    {}
func (r *Relay) AcceptEvent(evt *nostr.Event) bool {
	// block events that are too large
	jsonb, _ := json.Marshal(evt)
	return len(jsonb) <= 10000
}

func main() {
	if err := relayer.Start(&Relay{}); err != nil {
		log.Fatalf("server terminated: %v", err)
	}
}
