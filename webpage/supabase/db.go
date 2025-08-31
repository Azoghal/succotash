package supabase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/supabase-community/supabase-go"
)

// TODO this is assuming that we're going to use the rest api to interact with the db
// but really we will just go through the direct connection...

type Client interface {
	NoOp()
	GetTestEvents() ([]TestEvent, error)
}

type postgrestClient struct {
	supabaseClient *supabase.Client
}

func (c *postgrestClient) NoOp() {
}

func (c *postgrestClient) GetTestEvents() ([]TestEvent, error) {
	data, _, err := c.supabaseClient.From("test_events").Select("*", "", false).Execute()
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to select testEvents"))
	}

	resp := string(data)
	fmt.Printf("The data: %s\n", resp)

	var events []TestEvent
	err = json.Unmarshal(data, &events)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to unmarshall json"))
	}

	return events, nil
}

type DBClientFactory func() Client

func NewRestDBClientFactory(url, key string) DBClientFactory {
	return func() Client {

		fmt.Printf("%s: %s\n", url, key)

		c, err := supabase.NewClient(url, key, &supabase.ClientOptions{})
		if err != nil {
			fmt.Println("cannot initalize client", err)
		}
		return &postgrestClient{
			supabaseClient: c,
		}
	}

}

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
