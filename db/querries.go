package db

var createTableServers = `
CREATE TABLE IF NOT EXISTS servers
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR,
	port VARCHAR
);
`

var createTableServerParameters = `
CREATE TABLE IF NOT EXISTS server_parameters
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	param_key VARCHAR NOT NULL,
	param_value VARCHAR,
	server_id INTEGER,
	CONSTRAINT fk_servers
		FOREIGN KEY (server_id)
		REFERENCES servers(id)
		ON DELETE CASCADE
);
`
