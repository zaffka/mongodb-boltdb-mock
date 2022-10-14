/*Package db - пакет враппер для работы с базами

Состав пакета:
	db.go - набор общих интерфейсов для использования в коде
	mgo.go - реализация для драйвера globalsign/mgo
	mock.go - набор mock-структур для проведения тестирования
*/
package db

//New - экспортируемая функция-обёртка, для упрощения создания переменной интерфейсного типа
func New(self Handler) Handler {
	return self
}

//Handler - интерфейс, содержащий набор методов для работы с корневой сессией(коннектом) Монго
type Handler interface {
	Connect(resources ...interface{}) error
	Copy() Handler                                             //Метод возвращает интерфейсный тип с копией сессии, полученной от корневой
	CopyWithSettings(settings ...interface{}) (Handler, error) //Метод возвращает интерфейсный тип с копией сессии в нужном режиме
	Close()

	ExecOn(resources ...interface{}) Querier
}

//Querier - набор методов для работы с конечным сетом данных
type Querier interface {
	Insert(docs ...interface{}) error
	Remove(selector interface{}) error
	RemoveAll(selector interface{}) (num int, err error)
	Update(selector interface{}, update interface{}) error
	UpdateAll(selector interface{}, update interface{}) (num int, err error)
	Upsert(selector interface{}, update interface{}) (num int, err error)

	Find(query interface{}) Refiner //Метод позволяет уточнить и разобрать полученные данные
}

//Refiner - набор методов для уточнения запроса
type Refiner interface {
	One(result interface{}) error                  //One принимает в качестве параметра ссылку на структуру для анмаршалинга
	All(results interface{}) error                 //All принимает в качестве параметра ссылку на слайс структур для анмаршалинга
	Distinct(key string, result interface{}) error //Distinct распаковывает в result значения, полученные по ключу key
	Count() (num int, err error)
}
