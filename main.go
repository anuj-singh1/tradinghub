package main

import (
	"database/sql"
	"tradingdata/internal/api"
	"tradingdata/internal/config"
	db "tradingdata/internal/db/sqlc"
	"tradingdata/internal/log"
)

func main() {
	envConfig := config.GetConfig(config.EnvVariablePath, config.EnvFileName, config.EnvFileExtension)

	log.SetupLogger()
	DB, err := config.SetupDb(envConfig)
	if err != nil {
		log.GetLogger().Fatalln("cannot connect to job db:", err)
	}

	if DB != nil {
		defer func(DB *sql.DB) {
			err = DB.Close()
			if err != nil {
				log.GetLogger().Fatalln(err)
			}
		}(DB)
	}

	migrateErr := db.RunMigrations(DB)

	if migrateErr != nil {
		log.GetLogger().Fatalln("migration script failed: ", migrateErr)
	}

	server, err := api.NewServer(envConfig, db.NewTokenDb(DB))
	if err != nil {
		log.GetLogger().Fatalln("cannot create intialize server and routes:", err)
	}

	err = server.Start()
	if err != nil {
		log.GetLogger().Fatalln("cannot start server:", err)
	}
}
