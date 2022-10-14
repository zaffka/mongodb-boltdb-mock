/*Package db - MongoDB realization for the  globalsign/mgo mongo driver
Driver's page: https://github.com/globalsign/mgo
It's a community developing driver inspired by labix/mgo.v2
*/
package db

import (
	"errors"

	"github.com/globalsign/mgo"
)

//Mongo struct wraps *mgo.Session struct
type Mongo struct {
	*mgo.Session
}

//Connect dials to the mongo server
func (m *Mongo) Connect(resources ...interface{}) (err error) {
	dsn, ok := resources[0].(string)
	if !ok {
		return errors.New("Unexpected resources set, want `dsn string`")
	}

	m.Session, err = mgo.Dial(dsn)
	if err != nil {
		return err
	}
	return nil
}

//Copy makes session copy
func (m *Mongo) Copy() Handler {
	copy := m.Session.Copy()
	return &Mongo{copy}
}

//CopyWithSettings makes session copy with new working mode
func (m *Mongo) CopyWithSettings(settings ...interface{}) (Handler, error) {
	var mode int
	var refresh bool

	mode, ok := settings[0].(int)
	if !ok {
		return nil, errors.New("Unexpected parameters set, want `mode int, refresh bool`")
	}

	refresh, ok = settings[1].(bool)
	if !ok {
		return nil, errors.New("Unexpected parameters set, want `mode int, refresh bool`")
	}

	copy := m.Session.Copy()
	copy.SetMode(mgo.Mode(mode), refresh)
	return &Mongo{copy}, nil
}

//Close shuts down the link to db
func (m *Mongo) Close() {
	m.Session.Close()
}

/*
ExecOn - sets resources for the Mongo driver - databaseName and a collectionName
Is databaseName doesn't set driver will use one from the connection string
or will use "test" as a name.
*/
func (m *Mongo) ExecOn(resources ...interface{}) Querier {
	var databaseName, collectionName string
	var ok bool

	switch len(resources) {
	case 2:
		databaseName, _ = resources[0].(string)
		collectionName, ok = resources[1].(string)
		if !ok {
			collectionName = "test"
		}

		return &MongoCollection{m.Session.DB(databaseName).C(collectionName)}
	case 1:
		collectionName, ok = resources[0].(string)
		if !ok {
			collectionName = "test"
		}

		return &MongoCollection{m.Session.DB("").C(collectionName)}
	default:
		return &MongoCollection{m.Session.DB("").C("test")}
	}
}

//MongoCollection is a wrapper for *mgo.Collection
type MongoCollection struct {
	*mgo.Collection
}

//Insert puts documents to db
func (mc *MongoCollection) Insert(docs ...interface{}) error {
	err := mc.Collection.Insert(docs...)
	if err != nil {
		return err
	}
	return nil
}

//Remove deletes one document according to selector
func (mc *MongoCollection) Remove(selector interface{}) error {
	err := mc.Collection.Remove(selector)
	if err != nil {
		return err
	}
	return nil
}

//RemoveAll deletes all documents according to selector, return number of deleted docs or an error
func (mc *MongoCollection) RemoveAll(selector interface{}) (num int, err error) {
	info, err := mc.Collection.RemoveAll(selector)
	if err != nil {
		return 0, err
	}
	return info.Removed, nil
}

//Update updates one document and return an error if nothing to update
func (mc *MongoCollection) Update(selector interface{}, update interface{}) error {
	err := mc.Collection.Update(selector, update)
	if err != nil {
		return err
	}
	return nil
}

//UpdateAll updates documents and return an error if nothing to update or number of updated docs
func (mc *MongoCollection) UpdateAll(selector interface{}, update interface{}) (num int, err error) {
	info, err := mc.Collection.UpdateAll(selector, update)
	if err != nil {
		return 0, err
	}
	return info.Updated, nil
}

//Upsert updates document and insert a new doc if nothing to update
func (mc *MongoCollection) Upsert(selector interface{}, update interface{}) (num int, err error) {
	info, err := mc.Collection.Upsert(selector, update)
	if err != nil {
		return 0, err
	}
	return info.Updated, nil
}

//Find searches for the docs according to query
func (mc *MongoCollection) Find(query interface{}) Refiner {
	q := mc.Collection.Find(query)
	return &MongoQuery{q}
}

//MongoQuery wrapper for *mgo.Query
type MongoQuery struct {
	*mgo.Query
}

//One refines mongo query and return one record
func (mq *MongoQuery) One(result interface{}) error {
	err := mq.Query.One(result)
	if err != nil {
		return err
	}
	return nil
}

//All refines mongo query and return all records
func (mq *MongoQuery) All(results interface{}) error {
	err := mq.Query.All(results)
	if err != nil {
		return err
	}
	return nil
}

//Distinct selects values set for the key
func (mq *MongoQuery) Distinct(key string, result interface{}) error {
	err := mq.Query.Distinct(key, result)
	if err != nil {
		return err
	}
	return nil
}

//Count returns numbers of the queried records
func (mq *MongoQuery) Count() (num int, err error) {
	num, err = mq.Query.Count()
	if err != nil {
		return 0, err
	}
	return num, nil
}
