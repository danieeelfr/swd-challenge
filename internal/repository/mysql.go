package repository

import (
	"fmt"

	"github.com/danieeelfr/swd-challenge/internal/config"
	"github.com/danieeelfr/swd-challenge/internal/models"

	// the mysql driver
	_ "github.com/go-sql-driver/mysql"
	// the mysql gorm dialects
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/jinzhu/gorm"
)

// MySQLRepo holds the MySQL repo implementation
type MySQLRepo struct {
	cfg *config.MySQLRepositoryConfig
	DB  *gorm.DB
}

// NewMySQLRepo returns the repository implementation
func NewMySQLRepo(cfg *config.MySQLRepositoryConfig) *MySQLRepo {
	return &MySQLRepo{cfg: cfg}
}

// Connect to the database
func (r *MySQLRepo) Connect() error {
	db, err := gorm.Open("mysql",
		fmt.Sprintf(
			"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			r.cfg.DBUser, r.cfg.DBPassword, r.cfg.DBHost, r.cfg.DBPort, r.cfg.DBName,
		),
	)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.User{})

	r.DB = db

	return err
}
