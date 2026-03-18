package config

import "os"

func GetDBHost() string {
	return os.Getenv("DB_HOST")
}

func GetDBPort() string {
	return os.Getenv("DB_PORT")
}

func GetDBUser() string {
	return os.Getenv("DB_USER")
}

func GetDBPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func GetDBName() string {
	return os.Getenv("DB_NAME")
}

func GetDBSSLMode() string {
	return os.Getenv("DB_SSL_MODE")
}
