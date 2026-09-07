package config

type Config struct {
	PgDb PgDb
}

type PgDb struct {
	DSN string
}
