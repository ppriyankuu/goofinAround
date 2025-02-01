package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresClient struct {
	DB *gorm.DB
}

func NewPostgresClient(dsn string) (*PostgresClient, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &PostgresClient{DB: db}, nil
}

func (p *PostgresClient) AutoMigrate(models ...interface{}) error {
	return p.DB.AutoMigrate(models...)
}

func (p *PostgresClient) Save(task interface{}) error {
	return p.DB.Save(task).Error
}

func (p *PostgresClient) First(task interface{}, id uint) error {
	return p.DB.First(task, id).Error
}
