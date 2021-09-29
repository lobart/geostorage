package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/lobart/go_geoserver.git/internal/models"
)

type PostgreSqlDriver struct {
	Cfg *models.DBConfig
	db  *sql.DB
}

func (d *PostgreSqlDriver) Connect() error{
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		d.Cfg.Server.Host, d.Cfg.Server.Port, d.Cfg.Database.Username, d.Cfg.Database.Password, d.Cfg.Database.DBName)
	d.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	return nil
}

func (d *PostgreSqlDriver) Close() error{
	err:=d.db.Close()
	if err!=nil{
		return err
	}
	return nil
}


func (d *PostgreSqlDriver) Push(kick *models.KickConfig) error{
	var err error
	err = d.db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	sqlStatement := `
		INSERT INTO kicks (companyname, kickname, longitude, latitude, speed, status)
		VALUES ($1, $2, $3, $4, &5, &6)`

	_, err = d.db.Exec(sqlStatement, kick.CompanyName, kick.KickName, kick.Longitude, kick.Latitude, kick.Speed, kick.Status )
	if err != nil {
		panic(err)
		return err
	}

	fmt.Println("Successfully push!")
	return nil
}