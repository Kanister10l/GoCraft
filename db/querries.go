package db

var createTableServers = `
CREATE TABLE IF NOT EXISTS servers
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR,
	port VARCHAR,
	image VARCHAR
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

var createTableVolumes = `
CREATE TABLE IF NOT EXISTS volumes
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR,
	location VARCHAR,
	server_id INTEGER,
	CONSTRAINT fk_servers
		FOREIGN KEY (server_id)
		REFERENCES servers(id)
		ON DELETE CASCADE
)
`

var createTableContainerParameters = `
CREATE TABLE IF NOT EXISTS container_parameters
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	param_key VARCHAR,
	param_value VARCHAR,
	server_id INTEGER,
	CONSTRAINT fk_servers
		FOREIGN KEY (server_id)
		REFERENCES servers(id)
		ON DELETE CASCADE
)
`

// SelectVolumesByServer ...
var SelectVolumesByServer = `
SELECT v.id, v.name, v.location FROM volumes v
JOIN servers s ON s.id = v.server_id
WHERE s.id = %d
`

// SelectServerParametersByServer ...
var SelectServerParametersByServer = `
SELECT sp.id, sp.param_key, sp.param_value FROM server_parameters sp
JOIN servers s ON s.id = sp.server_id
WHERE s.id = %d
`

// SelectContainerParametersByServer ...
var SelectContainerParametersByServer = `
SELECT cp.id, cp.param_key, cp.param_value FROM container_parameters cp
JOIN servers s ON s.id = cp.server_id
WHERE s.id = %d
`

// SelectServers ...
var SelectServers = `
SELECT s.id, s.image, s.name, s.port FROM servers s
`
