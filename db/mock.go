package db

import "errors"

//Mock - структура для проверки методов db.Handler
type Mock struct {
	Msg     string
	Mode    int
	Refresh bool
	Closed  bool
}

//Connect - присваивает resource в поле Msg структуры
func (mk *Mock) Connect(resources ...interface{}) (err error) {
	dsn, ok := resources[0].(string)
	if !ok {
		return errors.New("Unexpected resources set, want `dsn string`")
	}
	mk.Msg = dsn
	return nil
}

//Copy - возвращает db.Handler со структурой &Mock{Msg: "session copied"}
func (mk *Mock) Copy() Handler {
	m := &Mock{}
	m.Msg = "session copied"
	return m
}

//CopyWithSettings - возвращает db.Handler со структурой &Mock{Msg: "session copied w settings"}
func (mk *Mock) CopyWithSettings(settings ...interface{}) (Handler, error) {
	m := &Mock{}
	m.Msg = "session copied w settings"
	return m, nil
}

//Close присваивает полю Closed значение true
func (mk *Mock) Close() {
	mk.Closed = true
}

//ExecOn - возвращает db.Querier со структурой &MockCollection{Msg: "ExecOn called"}
func (mk *Mock) ExecOn(resources ...interface{}) Querier {
	m := &MockCollection{}
	m.Msg = "ExecOn called"
	return m
}

//MockCollection - - структура для проверки методов db.Querier
type MockCollection struct {
	Msg      string
	DocsNum  int
	Selector int
	Upd      int
}

//Insert - считает количество переданных документов, пишет число в поле DocsNum
func (mc *MockCollection) Insert(docs ...interface{}) error {
	mc.DocsNum = len(docs)
	return nil
}

//Remove - пишет число 111 в поле Selector
func (mc *MockCollection) Remove(selector interface{}) error {
	mc.Selector = 111
	return nil
}

//RemoveAll - пишет число 333 в поле Selector
func (mc *MockCollection) RemoveAll(selector interface{}) (num int, err error) {
	return 333, nil
}

//Update - пишет число 555 в поле Selector и 777 в поле Upd
func (mc *MockCollection) Update(selector interface{}, update interface{}) error {
	mc.Selector = 555
	mc.Upd = 777
	return nil
}

//UpdateAll - возвращает число 888 и nil для ошибки
func (mc *MockCollection) UpdateAll(selector interface{}, update interface{}) (num int, err error) {
	return 888, nil
}

//Upsert - возвращает число 999 и nil для ошибки
func (mc *MockCollection) Upsert(selector interface{}, update interface{}) (num int, err error) {
	return 999, nil
}

//Find - возвращает db.Refiner со структурой &MockQuery{}
func (mc *MockCollection) Find(query interface{}) Refiner {
	return &MockQuery{}
}

//MockQuery - структура для проверки методов db.Refiner
type MockQuery struct {
	Res     string
	DistKey string
}

//One - пишет "result" в поле Res
func (mq *MockQuery) One(result interface{}) error {
	mq.Res = "result"
	return nil
}

//All - пишет "results" в поле Res
func (mq *MockQuery) All(results interface{}) error {
	mq.Res = "results"
	return nil
}

//Distinct - пишет key в поле DistKey
func (mq *MockQuery) Distinct(key string, result interface{}) error {
	mq.DistKey = key
	return nil
}

//Count - возвращает 999 и nil в качестве ошибки
func (mq *MockQuery) Count() (num int, err error) {
	return 999, nil
}
