package database

import (
	log "adeia-api/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var databaseConnection *sqlx.DB

func openConnection(dataSourceName, driverName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	log.Debug("Successfully connected to database: "+ driverName)
	return db, nil
}

func getConnection() *sqlx.DB {
	if databaseConnection == nil {
		datasourceName := getValidDriverName("localhost", "5432", "dyanesh", "varun", "mydb", "disable")
		databaseConnection, err := openConnection(datasourceName, "postgres")
		if err == nil {
			return databaseConnection
		}
		print(err)
		panic(err)
	}
	return databaseConnection
}

func getValidDriverName(host, port, user, password, dbname, sslmode string) string {
	return "host=" + host + " " +
		"port=" + port + " " +
		"user=" + user + " " +
		"password=" + password + " " +
		"dbname=" + dbname + " " +
		"sslmode=" + sslmode
}

func ExecuteQuery(query Query, parameters[] string) int64 {
	dbConn := getConnection()
	rows, err := dbConn.Exec(string(query), parameters)
	if err!=nil {
		log.Error(err)
		return 0
	}
	rowsCount, _ := rows.RowsAffected()
	return rowsCount
}

func Check() {
	dbConn := getConnection()
	rows := ExecuteQuery("INSERT into test values (1,2,3)", []string{})
	log.Debug(rows)
	//print(rows.LastInsertId())
	e := dbConn.Close()
	log.Debug(e)
}
