package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var databaseConnection *sqlx.DB

func openConnection(dataSourceName, driverName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	log.Debug("Successfully connected to database")
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

func ExecuteQuery(model interface{}, query Query, parameters[] string){
	dbConn := getConnection()

}

func Check() {
	dbConn := getConnection()
	rows, err := dbConn.Exec("CREATE TABLE IF NOT EXISTS test (A INT,B INT,C INT)")
	if err != nil {
		log.Error(err)
	}
	rows, err = dbConn.Exec("INSERT into test values (1,2,3)")
	if err != nil {
		log.Debug(err)
	}
	log.Debug(rows.RowsAffected())
	//print(rows.LastInsertId())
	e := dbConn.Close()
	log.Debug(e)
}
