package supabase

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// TODO this is assuming that we're going to use the rest api to interact with the db
// but really we will just go through the direct connection...

type Client interface {
	NoOp()
	GetTestEvents() ([]TestEvent, error)
	CheckSession(string) (bool, error)
}

type DBClientFactory func() Client

type pgClient struct {
	url string
}

func NewPGClientFactory(dbUrl string) DBClientFactory {
	return func() Client {
		return &pgClient{url: dbUrl}
	}
}

var _ Client = (*pgClient)(nil)

func (c *pgClient) NoOp() {}

func (c *pgClient) GetTestEvents() ([]TestEvent, error) {
	conn, err := pgx.Connect(context.Background(), c.url)
	if err != nil {
		log.Err(err).Msgf("Failed to connect to the database: %v", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	// Example query to test connection
	var events []TestEvent
	rows, err := conn.Query(context.Background(), "SELECT * from test_events")
	if err != nil {
		log.Err(err).Msgf("Query failed: %v", err)
		return nil, err
	}

	for rows.Next() {
		var id int
		var content string
		var createdAt time.Time
		err := rows.Scan(&id, &createdAt, &content)
		if err != nil {
			log.Err(err).Msgf("Query failed: %v", err)

			return nil, err
		}
		events = append(events, TestEvent{Id: id, Content: content})
	}
	if err = rows.Err(); err != nil {
		log.Err(err).Msgf("Query failed: %v", err)
		return nil, err
	}

	log.Info().Msgf("got events: %v", events)

	return events, nil
}

func (c *pgClient) CheckSession(sessionId string) (bool, error) {
	conn, err := pgx.Connect(context.Background(), c.url)
	if err != nil {
		log.Err(err).Msgf("Failed to connect to the database: %v", err)
		return false, err
	}
	defer conn.Close(context.Background())

	var id, userId string
	err = conn.QueryRow(context.Background(), "select id, user_id from auth.sessions where id=$1", sessionId).Scan(&id, &userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Just means not a real session
			return false, nil
		}
		log.Err(err).Msgf("error querying user session", err)
		return false, err
	}

	return true, nil
}
