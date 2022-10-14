package db_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zaffka/mongodb-boltdb-mock/db"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	mgoTestDSN = "mongodb://localhost:27017"
)

func TestRealMongo(t *testing.T) {
	mongo := db.New(&db.Mongo{})
	err := mongo.Connect(mgoTestDSN)
	if err != nil {
		t.Skipf("Failed to connect to db %s", err.Error())
	}
	defer mongo.Close()

	sess := mongo.Copy()
	defer sess.Close()

	var tt []bson.M

	for i := 0; i < 5; i++ {
		tt = append(tt, bson.M{"_id": bson.NewObjectId(), "msg": i, "created_at": time.Now()})
	}

	t.Run("Copy session", func(t *testing.T) {
		mgoStruct := sess.(*db.Mongo)
		assert.IsType(t, &db.Mongo{}, sess)
		assert.Equal(t, mgo.Mode(2), mgoStruct.Session.Mode())
	})

	t.Run("Copy session w mode", func(t *testing.T) {
		newsess, err := mongo.CopyWithSettings(3, true)
		if err != nil {
			t.Error(err)
		}
		mgoStruct, ok := newsess.(*db.Mongo)
		if !ok {
			t.Error("Failed at the type assertion")
		}
		assert.Equal(t, mgo.Mode(3), mgoStruct.Session.Mode())
	})

	t.Run("Insert one", func(t *testing.T) {
		err = sess.ExecOn("ctest").Insert(tt[0])
		if err != nil {
			t.Errorf("Insert failed w %s", err.Error())
		}
	})

	t.Run("Insert one w dbname", func(t *testing.T) {
		err = sess.ExecOn("test2", "ctest").Insert(tt[0])
		if err != nil {
			t.Errorf("Insert failed w %s", err.Error())
		}
	})

	t.Run("Insert one w excessive params", func(t *testing.T) {
		err = sess.ExecOn("excessive", "test2", "ctest").Insert(tt[0])
		if err != nil {
			t.Errorf("Insert failed w %s", err.Error())
		}
	})

	t.Run("Insert many", func(t *testing.T) {
		err = sess.ExecOn("ctest").Insert(tt[1], tt[2])
		if err != nil {
			t.Errorf("Insert failed w %s", err.Error())
		}
	})

	t.Run("Find one", func(t *testing.T) {
		var i interface{}
		err := sess.ExecOn("ctest").Find(bson.M{}).One(&i)
		if err != nil {
			t.Errorf("Got error %s", err.Error())
		}
		assert.NotZero(t, i)
	})

	t.Run("Find one w wrong dbname", func(t *testing.T) {
		var i interface{}
		err := sess.ExecOn("wrong", "ctest").Find(bson.M{}).One(&i)
		assert.Error(t, err)
	})

	t.Run("Find all", func(t *testing.T) {
		var i []interface{}
		sess.ExecOn("ctest").Find(bson.M{}).All(&i)
		if err != nil {
			t.Errorf("Got error %s", err.Error())
		}
		assert.NotZero(t, i)
	})

	t.Run("Distinct", func(t *testing.T) {
		var ints []int
		err := sess.ExecOn("ctest").Find(bson.M{}).Distinct("msg", &ints)
		if err != nil {
			t.Errorf("Got error %s", err.Error())
		}
		assert.Equal(t, 3, len(ints))
	})

	t.Run("Count", func(t *testing.T) {
		num, err := sess.ExecOn("ctest").Find(bson.M{}).Count()
		if err != nil {
			t.Errorf("Got error %s", err.Error())
		}
		assert.NotZero(t, num)
		t.Logf("...records found: %v", num)
	})

	t.Run("Update one", func(t *testing.T) {
		err := sess.ExecOn("ctest").Update(bson.M{"msg": 0}, bson.M{"$set": bson.M{"msg": 999}})
		if err != nil {
			t.Errorf("Got error %s", err.Error())
		}
	})

	t.Run("Update one w unexisted document", func(t *testing.T) {
		err := sess.ExecOn("ctest").Update(bson.M{"msg": 111}, bson.M{"$set": bson.M{"msg": 999}})
		assert.Error(t, err)
	})

	t.Run("Update all", func(t *testing.T) {
		num, err := sess.ExecOn("ctest").UpdateAll(bson.M{}, bson.M{"$set": bson.M{"updated": true}})
		if err != nil {
			t.Errorf("Got error %s", err.Error())
		}
		assert.NotZero(t, num)
	})

	t.Run("Upsert", func(t *testing.T) {
		num, err := sess.ExecOn("ctest").Upsert(bson.M{"msg": 333}, bson.M{"$set": bson.M{"upserted": true}})
		if err != nil {
			t.Errorf("Got error %s", err.Error())
		}
		assert.Zero(t, num)
	})

	t.Run("Remove one", func(t *testing.T) {
		err := sess.ExecOn("ctest").Remove(bson.M{"msg": 1})
		if err != nil {
			t.Errorf("Got error %s", err.Error())
		}
	})

	t.Run("Remove all", func(t *testing.T) {
		num, err := sess.ExecOn("ctest").RemoveAll(bson.M{})
		if err != nil {
			t.Errorf("Got error %s", err.Error())
		}
		assert.NotZero(t, num)
	})

	t.Run("Remove all", func(t *testing.T) {
		num, err := sess.ExecOn("test2", "ctest").RemoveAll(bson.M{})
		if err != nil {
			t.Errorf("Got error %s", err.Error())
		}
		assert.NotZero(t, num)
	})
}

func TestMock(t *testing.T) {
	mock := db.New(&db.Mock{})

	t.Run("Handler Connect", func(t *testing.T) {
		mock.Connect("dsnstring")
		m := mock.(*db.Mock)
		assert.Equal(t, "dsnstring", m.Msg)
	})

	t.Run("Handler Copy", func(t *testing.T) {
		copy := mock.Copy()
		m := copy.(*db.Mock)
		assert.Equal(t, "session copied", m.Msg)
	})

	t.Run("Handler Copy w mode", func(t *testing.T) {
		copy, err := mock.CopyWithSettings(true)
		if err != nil {
			t.Errorf("Got error %s", err.Error())
		}
		m := copy.(*db.Mock)
		assert.Equal(t, "session copied w settings", m.Msg)
	})

	t.Run("Handler Close", func(t *testing.T) {
		mock.Close()
		m := mock.(*db.Mock)
		assert.True(t, m.Closed)
	})

	t.Run("Handler ExecOn", func(t *testing.T) {
		colHandler := mock.ExecOn()
		d := colHandler.(*db.MockCollection)
		assert.IsType(t, &db.MockCollection{}, colHandler)
		assert.Equal(t, "ExecOn called", d.Msg)
	})

	t.Run("Querier Insert", func(t *testing.T) {
		colHandler := mock.ExecOn()
		colHandler.Insert("testdoc1", "testdoc2")
		c := colHandler.(*db.MockCollection)
		assert.Equal(t, 2, c.DocsNum)
	})

	t.Run("Querier Remove", func(t *testing.T) {
		colHandler := mock.ExecOn()
		colHandler.Remove("testdoc1")
		c := colHandler.(*db.MockCollection)
		assert.Equal(t, 111, c.Selector)
	})

	t.Run("Querier RemoveAll", func(t *testing.T) {
		colHandler := mock.ExecOn()
		num, _ := colHandler.RemoveAll("testdoc1")
		assert.Equal(t, 333, num)
	})

	t.Run("Querier Update", func(t *testing.T) {
		colHandler := mock.ExecOn()
		colHandler.Update("selector", "update")
		c := colHandler.(*db.MockCollection)
		assert.Equal(t, 555, c.Selector)
		assert.Equal(t, 777, c.Upd)
	})

	t.Run("Querier UpdateAll", func(t *testing.T) {
		colHandler := mock.ExecOn()
		num, _ := colHandler.UpdateAll("selector", "update")
		assert.Equal(t, 888, num)
	})

	t.Run("Querier Upsert", func(t *testing.T) {
		colHandler := mock.ExecOn()
		num, _ := colHandler.Upsert("selector", "update")
		assert.Equal(t, 999, num)
	})

	t.Run("Querier Find, returns Refiner", func(t *testing.T) {
		colHandler := mock.ExecOn()
		queryHandler := colHandler.Find("query")
		assert.IsType(t, &db.MockQuery{}, queryHandler)
	})

	t.Run("Refiner One", func(t *testing.T) {
		var res string
		queryHandler := mock.ExecOn().Find("query")
		queryHandler.One(&res)
		str := queryHandler.(*db.MockQuery)
		assert.Equal(t, "result", str.Res)
	})

	t.Run("Refiner All", func(t *testing.T) {
		var res string
		queryHandler := mock.ExecOn().Find("query")
		queryHandler.All(&res)
		str := queryHandler.(*db.MockQuery)
		assert.Equal(t, "results", str.Res)
	})

	t.Run("Refiner Distinct", func(t *testing.T) {
		var res string
		queryHandler := mock.ExecOn().Find("query")
		queryHandler.Distinct("key", &res)
		str := queryHandler.(*db.MockQuery)
		assert.Equal(t, "key", str.DistKey)
	})

	t.Run("Refiner Count", func(t *testing.T) {
		num, _ := mock.ExecOn().Find("query").Count()
		assert.Equal(t, 999, num)
	})

}

func TestBoltDB(t *testing.T) {

	bolt := db.New(&db.Bolt{})
	err := bolt.Connect("bolt", "bucketOne", "bucketTwo")
	if err != nil {
		t.Fatalf("Failed to open bolt file, %v", err)
	}
	defer bolt.Close()

	t.Run("Insert", func(t *testing.T) {
		err := bolt.ExecOn("bucketOne").Insert("key", &db.Mock{Msg: "test"})
		assert.NoError(t, err)
	})

	t.Run("Insert", func(t *testing.T) {
		err := bolt.ExecOn("bucketOne").Insert("key2", &db.Mock{Msg: "test"})
		assert.NoError(t, err)
	})

	t.Run("Insert wo bucket", func(t *testing.T) {
		err := bolt.ExecOn(nil).Insert("key3", &db.Mock{Msg: "test"})
		assert.NoError(t, err)
	})

	t.Run("Insert w non-existed bucket", func(t *testing.T) {
		err := bolt.ExecOn("nobucket").Insert("key3", &db.Mock{Msg: "test"})
		assert.Error(t, err)
	})

	t.Run("Remove", func(t *testing.T) {
		err := bolt.ExecOn("bucketOne").Remove("key2")
		assert.NoError(t, err)
	})

	t.Run("Read inserted value by key", func(t *testing.T) {
		var res db.Mock
		err := bolt.ExecOn("bucketOne").Find("key").One(&res)
		if err != nil {
			t.Errorf("Failed w %v", err)
		}
		assert.Equal(t, "test", res.Msg)
		t.Logf("Result: %v", res)
	})

	t.Run("Read from nil bucket", func(t *testing.T) {
		var res db.Mock
		err := bolt.ExecOn(nil).Find("key").One(&res)
		assert.Error(t, err)
	})

	t.Run("Read from unexisted bucket", func(t *testing.T) {
		var res db.Mock
		err := bolt.ExecOn("nobucket").Find("key").One(&res)
		assert.Error(t, err)
	})
}
