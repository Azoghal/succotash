package supabase

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/supabase-community/supabase-go"
)

// TODO this is assuming that we're going to use the rest api to interact with the db
// but really we will just go through the direct connection...

type Client interface {
	NoOp()
	GetTestEvents() ([]TestEvent, error)
}

type client struct {
	supabaseClient *supabase.Client
}

func (c *client) NoOp() {
}

func (c *client) GetTestEvents() ([]TestEvent, error) {
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

type RestDBClientFactory func(url, key string) Client

func NewRestDBClientFactory() RestDBClientFactory {
	return func(url, key string) Client {

		fmt.Printf("%s: %s\n", url, key)

		c, err := supabase.NewClient(url, key, &supabase.ClientOptions{})
		if err != nil {
			fmt.Println("cannot initalize client", err)
		}
		return &client{
			supabaseClient: c,
		}
	}

}
