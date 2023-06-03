package database

import "os"

const (
	driver = "mysql"
	net    = "tcp"
	ip     = "127.0.0.1"
	port   = "3306"
	dBName = "netBlast"
)

var dataSourceName = os.Getenv("DBUSER") + ":" + os.Getenv("DBPASS") + "@" + net + "(" + ip + ":" + port + ")/" + dBName
