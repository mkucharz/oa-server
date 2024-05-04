package db

type SystemHealthInterface interface {
	Ping() error
}
