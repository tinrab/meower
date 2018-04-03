package main

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

type config struct {
	Port   string
	DBPort int
}

var (
	cfg             config
	dbClusterConfig *gocql.ClusterConfig
	dbSession       *gocql.Session
)

func init() {
	log.SetLevel(log.DebugLevel)
	cfg = config{
		Port: os.Getenv("APP_PORT"),
	}
	cfg.DBPort, _ = strconv.Atoi(os.Getenv("APP_DB_PORT"))

	dbClusterConfig = gocql.NewCluster("db")
	dbClusterConfig.DisableInitialHostLookup = true
	dbClusterConfig.Port = cfg.DBPort
	dbClusterConfig.ProtoVersion = 4
	dbClusterConfig.Keyspace = "meower"
}

func main() {
	var err error

	dbSession, err = dbClusterConfig.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer dbSession.Close()

	r := gin.Default()
	r.POST("/meows", createMeowEndpoint)
	if err = r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
