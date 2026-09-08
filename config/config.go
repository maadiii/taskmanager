package config

import (
	_ "github.com/joho/godotenv/autoload"
)

func NewConfig() *Config {
	return &Config{
		Server: new(Server).init(),
		PgDb:   new(PgDb).init(),
	}
}

type Config struct {
	Server *Server
	PgDb   *PgDb
}

type Server struct {
	Port int
}

func (s *Server) init() *Server {
	s.Port = getEnvInt("PORT")

	return s
}

type PgDb struct {
	DSN string
}

func (pg *PgDb) init() *PgDb {
	pg.DSN = getEnvString("PG_DSN")

	return pg
}
