package postgresqlDriver

import (
	"database/sql"
	"fmt"
	"github.com/lobart/go_geoserver.git/pkg/models"
)

type PostgreSqlDriver struct {
	Cfg models.DBConfig
	db *sql.DB
}

func (d PostgreSqlDriver) Connect() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		d.Cfg.Server.Host, d.Cfg.Server.Port, d.Cfg.Database.Username, d.Cfg.Database.Password, d.Cfg.Database.DBName)
	d.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
}

func (d PostgreSqlDriver) Close() {
	d.db.Close()
}


func (d PostgreSqlDriver) Push(kick models.KickConfig) {
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
	}

	fmt.Println("Successfully push!")
}