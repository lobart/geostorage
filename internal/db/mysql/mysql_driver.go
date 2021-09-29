package mysql

import (
	"database/sql"
	"fmt"
	"github.com/lobart/go_geoserver.git/internal/models"
)

type MySqlDriver struct {
	Cfg *models.DBConfig
	db  *sql.DB
}

func (d *MySqlDriver) Connect() error {
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		d.Cfg.Database.Username, d.Cfg.Database.Password, d.Cfg.Server.Host, d.Cfg.Server.Port , d.Cfg.Database.DBName)
	var err error

	d.db, err = sql.Open("mysql", mysqlInfo)

	// if there is an error opening the connection, handle it
	if err != nil {
		return err
	}
	return nil
}

func (d *MySqlDriver) Close() error{
	err:= d.db.Close()
	if err!=nil{
		return err
	}
	return nil
}

func (d *MySqlDriver) Push(kick *models.KickConfig) error{

	sqlStatement := fmt.Sprintf("INSERT INTO kicks (companyname, kickname, longitude, latitude, speed, status) VALUES (%s, %s, %s, %s, %s, %s)", kick.CompanyName, kick.KickName, kick.Longitude, kick.Latitude, kick.Speed, kick.Status)
	// perform a db.Query insert
	insert, err := d.db.Query(sqlStatement)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
		return err
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

	fmt.Println("Successfully push!")
	return nil
}
