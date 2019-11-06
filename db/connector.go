package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/kanister10l/GoCraft/logger"

	"github.com/kanister10l/GoCraft/eventmanager"

	"github.com/ztrue/tracerr"

	//Database driver import
	_ "github.com/mattn/go-sqlite3"
)

// Connection ...
type Connection struct {
	Database *sql.DB
	Location string
}

// Connect ...
func Connect(location string) *Connection {
	database, err := sql.Open("sqlite3", location)
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		eventmanager.Master.ExecEvent("stop")
	}

	connection := &Connection{
		Database: database,
		Location: location,
	}

	eventmanager.Master.Register("stop", func() {
		connection.Database.Close()
		logger.Logger.Infow("database connection closed")
	})

	logger.Logger.Infow("succesfully connected to database", "db_location", location)

	return connection
}

// InitDatabase ...
func (c *Connection) InitDatabase() {
	stm, err := c.Database.Prepare(createTableServers)
	if err != nil {
		logger.Logger.Errorw("error preparing statement", "table name", "servers")
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		eventmanager.Master.ExecEvent("stop")
	}

	_, err = stm.Exec()
	if err != nil {
		logger.Logger.Errorw("error creating table", "table name", "servers")
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		eventmanager.Master.ExecEvent("stop")
	}

	stm, err = c.Database.Prepare(createTableServerParameters)
	if err != nil {
		logger.Logger.Errorw("error preparing statement", "table name", "server_parameters")
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		eventmanager.Master.ExecEvent("stop")
	}

	_, err = stm.Exec()
	if err != nil {
		logger.Logger.Errorw("error creating table", "table name", "server_parameters")
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		eventmanager.Master.ExecEvent("stop")
	}
}

// Test ...
func Test() {
	database, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		os.Exit(127)
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		os.Exit(127)
	}

	_, err = statement.Exec()
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		os.Exit(127)
	}

	statement, err = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		os.Exit(127)
	}

	_, err = statement.Exec("Rob", "Gronkowski")
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		os.Exit(127)
	}

	rows, err := database.Query("SELECT id, firstname, lastname FROM people")
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		os.Exit(127)
	}

	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}

	err = database.Close()
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		os.Exit(127)
	}
}
