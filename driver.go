package golangNeo4jBoltDriver

import (
	"database/sql"
	"database/sql/driver"
)

var (
	magicPreamble     = []byte{0x60, 0x60, 0xb0, 0x17}
	supportedVersions = []byte{
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}
	handShake          = append(magicPreamble, supportedVersions...)
	noVersionSupported = [4]byte{0x00, 0x00, 0x00, 0x00}
)

const (
	// Version is the current version of this driver
	Version = "1.0"
	// ClientID is the id of this client
	ClientID = "GolangNeo4jBolt/" + Version
)

// Driver allows connections to Neo4j. It implements driver.Driver and also
// includes its own more Neo-friendly interfaces. Some of the features of this
// interface implement neo-specific features unavailable in the driver.Driver
// compatible interface
//
// Driver objects should be thread safe, so you can use them to open connections
// in multiple threads. The connection objects themselves, and any prepared
// statements/transactions within, are not thread safe.
type Driver interface {
	// Open opens a sql.driver compatible connection. Used internally
	// by the go sql interface
	Open(string) (driver.Conn, error)

	// OpenNeo opens a Neo-specific connection. This should be used
	// directly when not using the golang sql interface
	OpenNeo(string) (Conn, error)
}

type boltDriver struct{}

// NewDriver creates a new Driver object
func NewDriver() Driver {
	return &boltDriver{}
}

// Open opens a new Bolt connection to the Neo4J database
func (d *boltDriver) Open(connStr string) (driver.Conn, error) {
	return newBoltConn(connStr, nil) // Never use pooling when using SQL driver
}

// Open opens a new Bolt connection to the Neo4J database. Implements a Neo-friendly alternative to sql/driver.
func (d *boltDriver) OpenNeo(connStr string) (Conn, error) {
	return newBoltConn(connStr, nil)
}

// DriverPool is a driver allowing connection to Neo4j with support for connection pooling
// The driver allows you to open a new connection to Neo4j
//
// Driver objects should be THREAD SAFE, so you can use them
// to open connections in multiple threads. The connection objects
// themselves, and any prepared statements/transactions within ARE NOT
// THREAD SAFE.
type DriverPool interface {
	// OpenPool opens a Neo-specific connection.
	OpenPool() (Conn, error)
	reclaim(*boltConn)
}

type boltDriverPool struct {
	connStr  string
	maxConns int
	pool     chan *boltConn
}

// NewDriverPool creates a new Driver object with connection pooling
func NewDriverPool(connStr string, max int) (DriverPool, error) {
	d := &boltDriverPool{
		connStr:  connStr,
		maxConns: max,
		pool:     make(chan *boltConn, max),
	}
	for i := 0; i < max; i++ {
		conn, err := newPooledBoltConn(connStr, d)
		if err != nil {
			return nil, err
		}
		d.pool <- conn
	}
	return d, nil
}

// OpenNeo opens a new Bolt connection to the Neo4J database.
func (d *boltDriverPool) OpenPool() (Conn, error) {
	conn := <-d.pool
	if conn.conn == nil {
		if err := conn.initialize(nil); err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func (d *boltDriverPool) reclaim(conn *boltConn) {
	// sneakily swap out connection so a reference to
	// it isn't held on to
	newConn := &boltConn{}
	*newConn = *conn
	d.pool <- newConn
	conn = nil
}

func init() {
	sql.Register("neo4j-bolt", &boltDriver{})
	sql.Register("neo4j-bolt-recorder", &recorder{})
}
