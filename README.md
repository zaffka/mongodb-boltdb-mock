# NOTICE

This repo is restored.
I've deleted it from github, which broke InVisionApp/go-health:
https://github.com/InVisionApp/go-health/issues/79

# Wrapper with mocks for MongoDB and BoltDB

Package shows how to mock MongoDB or BoltDB using a single Go interface.  
It's not cover all db methods from the [globalsign/mgo](https://github.com/globalsign/mgo) driver's package.  
Feel free to extend the package according to your needs.

Package inspired by this article - [Mocking Mongo in Go](http://thylong.com/golang/2016/mocking-mongo-in-golang/)

**Table of contents**

- [How to install](#how-to-install)
- [Package contents](#package-contents)
- [Interface methods being in use by realization](#interface-methods-being-in-use-by-realization)
- [MongoDB examples](#mongodb-examples)
- [BoltDB examples](#boltdb-examples)
- [Mocking](#mocking)

## How to install

`go get github.com/zaffka/mongodb-boltdb-mock`

## Package contents

- Wrapper itself (db/db.go) with:

  - db.Handler - main interface for opening, closing and controling connection to DB
  - db.Querier - interface for selecting data from DB
  - db.Refiner - interface for refining query object obtained by the Find() method

- Realization for the MongoDB (db/mgo.go)
- Realization for the BoltDB (db/bolt.go)
- A set of mocks (db/mock.go)
- Unit tests (db/db_test.go)

## Interface methods being in use by realization

| Interface/Function | MongoDB | Mocks  | BoltDB |
| ------------------ | ------- | ------ | ------ |
| **db.New()**       | +       | +      | +      |
| **db.Handler**     | &nbsp;  | &nbsp; | &nbsp; |
| Connect            | +       | +      | +      |
| Copy               | +       | +      | -      |
| CopyWithSettings   | +       | +      | -      |
| Close              | +       | +      | +      |
| ExecOn             | +       | +      | +      |
| **db.Querier**     | &nbsp;  | &nbsp; | &nbsp; |
| Insert             | +       | +      | +      |
| Remove             | +       | +      | +      |
| RemoveAll          | +       | +      | -      |
| Update             | +       | +      | -      |
| UpdateAll          | +       | +      | -      |
| Upsert             | +       | +      | -      |
| Find               | +       | +      | +      |
| **db.Refiner**     | &nbsp;  | &nbsp; | &nbsp; |
| One                | +       | +      | +      |
| All                | +       | +      | -      |
| Distinct           | +       | +      | -      |
| Count              | +       | +      | -      |

## MongoDB examples

### ...up the db

```go
mongo := db.New(&db.Mongo{})
err := mongo.Connect("mongo://localhost:/27017")
defer mongo.Close()
```

#### ...inserting data

```go
func main() {
mongo := db.New(&db.Mongo{})
err := mongo.Connect("mongo://localhost:/27017")
defer mongo.Close()

os := new(OurStuct)
os.Save(mongo)
}

type OurStruct struct {}
func (o *OurStruct) Save(db db.Handler, databaseName, collectionName string) error {
    sess := db.Copy() //Pay attention globalsign's driver whant us to use
	defer sess.Close() //copied or clonned session of the *mgo.Session object

	err := sess.ExecOn(databaseName, collectionName).Insert(o)
	if err != nil {
		return err
	}

	return nil
}
```

#### ...reading

```go
func main() {
mongo := db.New(&db.Mongo{})
err := mongo.Connect("mongo://localhost:/27017")
defer mongo.Close()

os := new(OurStuct)
os.Save(mongo)
}

type OurStruct struct {
    ID bson.ObjectID `bson:"_id"`
}
func (o *OurStruct) Read(db db.Handler, databaseName, collectionName string) error {
	sess := db.Copy()
	defer sess.Close()

	err := sess.ExecOn(databaseName, collectionName).Find(o.ID).One(o)
	if err != nil {
		return err
	}

	return nil
}
```

### ...updating

```go
func main() {
mongo := db.New(&db.Mongo{})
err := mongo.Connect("mongo://localhost:/27017")
defer mongo.Close()

os := new(OurStuct)
os.Save(mongo)
}

type OurStruct struct {
    ID bson.ObjectID `bson:"_id"`
}
func (o *OurStruct) Update(db db.Handler, databaseName, collectionName, update string) error {
	sess := db.Copy()
	defer sess.Close()

	_, err := sess.ExecOn(databaseName, collectionName).Upsert(bson.M{"_id": o.ID}, bson.M{"$set": bson.M{"somefield": update}})
	if err != nil {
		return err
	}

	return nil
}
```

and so on...

## BoltDB examples

Features:

- Database opens at the system's temp directory (/tmp @ linux)
- Working directory with db file will look like `/tmp/[basename][random numbers]/[basename]`
- bolt.Close() removes the working directory and a db file
- Use [boltbrowser](https://github.com/br0xen/boltbrowser) to work with bolt's files
- Any structs and data types can be used as keys and values to store in BoltDB (Gob marshaling\unmarshaling inside)
- BoltDB uses buckets as Mongo's collections analogues

### ...up the db

```go
bolt := db.New(&db.Bolt{})
err := bolt.Connect("bolt", "bucketOne", "bucketTwo")
defer bolt.Close()
```

Where `"bolt"` - it's a `[basename]` and others are variadic set of buckets names.  
If the list scipped bucket with the `default` name being in use.

#### ...inserting data

```go
bolt := db.New(&db.Bolt{})
err := bolt.Connect("bolt", "bucketOne", "bucketTwo")
defer bolt.Close()

err = bolt.ExecOn("bucketOne").Insert("key", &db.Mock{Msg: "test"})
if err != nil {
	log.Error(err)
}
```

#### ...reading

```go
...
var res db.Mock
err := bolt.ExecOn("bucketOne").Find("key").One(&res)
...
```

### ...deleting

```go
...
err := bolt.ExecOn("bucketOne").Remove("key")
...
```

## Mocking

Just replace `&db.Mongo{}` (or `&db.Bolt{}`) with `&db.Mock{}` and cover your functions by unit tests with ease.  
No real database needed.

```go
func main() {
mongo := db.New(&db.Mongo{})
mongo.ExecOn(...).Find(...)...
...
```

replace with

```go
func main() {
mock := db.New(&db.Mock{})
mongo.ExecOn(...).Find(...)...
...
```
