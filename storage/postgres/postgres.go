package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)


type PostgresStorage struct {
	DB *sqlx.DB
}

func NewPostgresStorage(configg *viper.Viper) (*PostgresStorage,error) {
	db,err := ConnectDB(configg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &PostgresStorage{
		DB: db,
	},nil
}

func ConnectDB(configg *viper.Viper) (*sqlx.DB ,error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	configg.GetString("database.host"),
	configg.GetString("database.port"),
	configg.GetString("database.user"),
	configg.GetString("database.password"),
	configg.GetString("database.dbname"),
	configg.GetString("database.sslmode"),
	))
	if err != nil {
		return nil, err
	}
	return db , nil
}