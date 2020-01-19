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
	if c.Execute(createTableServers) != nil {
		eventmanager.Master.ExecEvent("stop")
		return
	}

	if c.Execute(createTableServerParameters) != nil {
		eventmanager.Master.ExecEvent("stop")
		return
	}

	if c.Execute(createTableContainerParameters) != nil {
		eventmanager.Master.ExecEvent("stop")
		return
	}

	if c.Execute(createTableVolumes) != nil {
		eventmanager.Master.ExecEvent("stop")
		return
	}
}

// Execute ...
func (c *Connection) Execute(statement string) error {
	stm, err := c.Database.Prepare(statement)
	if err != nil {
		logger.Logger.Errorw("error preparing statement", "statement", statement)
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return err
	}

	_, err = stm.Exec()
	if err != nil {
		logger.Logger.Errorw("error executing statement", "statement", statement)
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return err
	}

	return nil
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
