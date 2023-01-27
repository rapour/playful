package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

type Client struct {
	Config  Config
	Session *gocql.Session
}

func (c *Client) Close() {
	c.Session.Close()
}

func NewCassandraClient(c Config) (*Client, error) {

	cluster := gocql.NewCluster(c.Address)
	initialSession, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer initialSession.Close()

	createNamespace := initialSession.Query(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
    WITH replication = {
        'class' : 'SimpleStrategy',
        'replication_factor' : %d
    }`, c.Keyspace, c.ReplicationFactor))

	err = createNamespace.Exec()
	if err != nil {
		return nil, err
	}

	cluster.Keyspace = c.Keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("main session creation failed: %v", err)
	}

	return &Client{
		Config:  c,
		Session: session,
	}, nil
}
