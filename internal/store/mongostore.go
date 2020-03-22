package store

import (
	"time"

	"github.com/sebastianrosch/couchconnections/internal/db"

	// "github.com/sebastianrosch/couchconnections/pkg/strutil"
	// "github.com/sebastianrosch/couchconnections/pkg/types"

	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
)

const (
	// EventsCollection the collection name of the events collection
	EventsCollection = "lrp.events"
	// EventsIndex the index name for the events.from index
	EventsIndex = "index.events.id"
)

type Event struct {
	ID          string    `bson:"id"`
	Topic       string    `bson:"topic"`
	Description string    `bson:"description"`
	Host        string    `bson:"host"`
	ZoomLink    string    `bson:"zoomLink"`
	Start       time.Time `bson:"start"`
}

// MongoStore is the service store for MongoDB.
type MongoStore struct {
	db     *mgo.Database
	events *mgo.Collection
}

// NewMongoStore returns an instance of MongoStore connected to a mongo database.
func NewMongoStore(uri, databaseName, username, password string) (*MongoStore, error) {
	session, err := db.Connect(uri, username, password)
	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to db")
	}

	db := session.DB(databaseName)
	events := db.C(EventsCollection)
	if err := events.EnsureIndex(mgo.Index{
		Key:        []string{"id"},
		Unique:     false,
		Name:       EventsIndex,
		Background: true,
	}); err != nil {
		return nil, errors.Wrapf(err, "could not ensure index")
	}

	return &MongoStore{
		db:     db,
		events: events,
	}, nil
}

// Name gets the name of this db implementation.
func (s *MongoStore) Name() string {
	return "mongo"
}

// CheckReadiness checks the readiness of the db and returns an error if it's
// not ready.
func (s *MongoStore) CheckReadiness() error {
	return s.db.Session.Ping()
}

// GetAllEvents returns all events.
func (s *MongoStore) GetAllEvents() ([]Event, error) {
	var results []Event

	err := s.events.Find(nil).All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// CreateEvent adds a new event.
func (s *MongoStore) CreateEvent(topic, description string) (*Event, error) {
	event := &Event{
		Topic:       topic,
		Description: description,
		Start:       time.Now(),
	}
	err := s.events.Insert(event)

	return event, err
}
