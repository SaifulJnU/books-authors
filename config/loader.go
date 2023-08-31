package config

import "os"

var (
	ENV      string
	MongoURL string
)

func GetEnvDefault(key string, defVal string) string {

	val, ex := os.LookupEnv(key)
	if !ex {
		val = defVal
	}
	return val

}

func SetEnvionment() {
	ENV = GetEnvDefault("ENV", "local")
	MongoURL = GetEnvDefault("Mongo_URL", "mongodb://admin:secret@localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs0")
}
