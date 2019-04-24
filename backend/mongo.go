package backend

import (
	"github.com/globalsign/mgo"
)

// MongoDriver represents a "stacktrace.fun" MongoDB driver.
type MongoDriver struct {
	Session    *mgo.Session
	Message    *mgo.Collection
	Punishment *mgo.Collection
}

// Connect attempts to connect with the remote MongoDB server.
func (driver *MongoDriver) Connect(uri string, databaseName string) error {
	session, err := mgo.Dial(uri)
	if err != nil {
		return err
	}
	driver.Session = session

	database := driver.Session.DB(databaseName)
	driver.Message = database.C("messages")
	driver.Punishment = database.C("punishments")

	return nil
}
