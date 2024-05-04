package model

type SystemHealthInterface interface {
	PingDatabase() error
}

func (model *Model) PingDatabase() error {
	return model.db.Ping()
}
