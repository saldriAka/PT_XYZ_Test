package config

type Config struct {
	Server      Server
	Database    Database
	Jwt         Jwt
	Storage     Storage
	ServiceMode ServiceMode
}

type Server struct {
	Host   string
	Port   string
	Assets string
}

type Database struct {
	Host string
	Port string
	Name string
	User string
	Pass string
	Tz   string
}

type Jwt struct {
	Key string
	Exp int
}

type Storage struct {
	BasePath string
}

type ServiceMode struct {
	ServiceMode string
}
