package pkg

// Client side prefixes - error location
var (
	ClRegister = "Client/Register: "
	ClMessage  = "Client/Message: "
	ClTest     = "Client/Test: "
	ClList     = "Client/List: "
	Cl         = "Client: "
)

// Server side prefixes - error location
var (
	SvRegister = "Server/Register: "
	SvMessage  = "Server/Message: "
	SvTest     = "server/Test: "
	SvList     = "server/List: "
	Sv         = "Server: "
)

var (
	DB = "Database: "
)

// Suffixes - error reason
var (
	BadParseTo   = "couldnt parse to json."
	BadParseFrom = "couldnt parse from json."
	BadReq       = "couldnt create request."
	BadRes       = "couldnt receive request."
	BadConn      = "couldnt connect."
	BadRead      = "couldnt read response body."
	BadWrite     = "couldnt write body."
	BadOpen      = "couldnt open."
	BadClose     = "couldnt close."
	BadDir       = "couldnt get directory."
)
