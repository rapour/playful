package cassandra

import (
	"fmt"
	"time"

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
	cluster.Timeout = 3 * time.Second

	initialSession, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer initialSession.Close()

	// dev
	// err = initialSession.Query(fmt.Sprintf("DROP KEYSPACE IF EXISTS %s", c.Keyspace)).Exec()
	// if err != nil {
	// 	return nil, err
	// }

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
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("main session creation failed: %v", err)
	}

	return &Client{
		Config:  c,
		Session: session,
	}, nil
}
