package config

import "time"

var Config *Cfg

var TimeLocation *time.Location

type Cfg struct {
	App  AppConfig
	Db   DB
	Data Data
}

type AppConfig struct {
	Name     string
	Url      string
	Port     int
	Env      string
	Debug    bool
	Timezone string
}

type DB struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Name       string
	Connection DbConn
}

type DbConn struct {
	Open int
	TTL  int
	Idle int
}

type Data struct {
	MaxRows uint
}
