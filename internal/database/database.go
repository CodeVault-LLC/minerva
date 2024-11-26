package database

import (
	contentSchema "github.com/codevault-llc/minerva/internal/contents/models"
	coreSchema "github.com/codevault-llc/minerva/internal/core/models"
	networkSchema "github.com/codevault-llc/minerva/internal/network/models"
	"github.com/gocql/gocql"
)

func turnStringTogether(strings ...string) string {
	var result string
	for _, str := range strings {
		result += str
	}
	return result
}

type Database struct {
	Db *gocql.Session
}

func NewDatabase() (*Database, error) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Consistency = gocql.Quorum

	tempSession, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer tempSession.Close()

	err = createSchema(tempSession)
	if err != nil {
		return nil, err
	}

	cluster.Keyspace = "minerva"
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &Database{Db: session}, nil
}

func createSchema(session *gocql.Session) error {
	err := session.Query(turnStringTogether(coreSchema.CoreSchema, networkSchema.NetworkSchema, contentSchema.ContentSchema)).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) Close() {
	d.Db.Close()
}

func (d *Database) GetDatabase() *gocql.Session {
	return d.Db
}
