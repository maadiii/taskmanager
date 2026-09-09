package config

import (
	_ "github.com/joho/godotenv/autoload"
)

func NewConfig() *Config {
	return &Config{
		Server: new(Server).init(),
		PgDb:   new(PgDb).init(),
		Redis:  new(Redis).init(),
	}
}

type Config struct {
	Server *Server
	PgDb   *PgDb
	Redis  *Redis
}

type Redis struct {
	Addr string
}

func (r *Redis) init() *Redis {
	r.Addr = getEnvString("REDIS_ADDR")

	return r
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
