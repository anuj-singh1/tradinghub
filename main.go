package main

import (
	"github.com/joho/godotenv"
	"tradingdata/internal/api"
	"tradingdata/internal/config"
	"tradingdata/internal/log"
)

func loadEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.GetLogger().Fatalln("Error loading app.env file", err)
	}
}

func main() {
	envConfig := config.GetConfig(config.EnvVariablePath, config.EnvFileName, config.EnvFileExtension)

	//loadEnv("app.env")
	log.SetupLogger()
	//DB, err := config.SetupDb(envConfig)
	//if err != nil {
	//	log.Fatal("cannot connect to job db:", err)
	//}
	//
	//if DB != nil {
	//	defer func(DB *sql.DB) {
	//		err = DB.Close()
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//	}(DB)
	//}
	//
	//migrateErr := db.RunMigrations(DB)
	//
	//if migrateErr != nil {
	//	log.Fatal("migration script failed: ", migrateErr)
	//}

	server, err := api.NewServer(envConfig)
	if err != nil {
		log.Fatal("cannot create intialize server and routes:", err)
	}

	err = server.Start(envConfig.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}


