package pkg

// Client side prefixes - error location
var (
	ClRegister = "Client/Register: "
	ClMessage  = "Client/Message: "
	Cl         = "Client: "
)

// Server side prefixes - error location
var (
	SvRegister = "Server/Register: "
	SvMessage  = "Server/Message: "
	Sv         = "Server: "
)

// Suffixes - error reason
var (
	BadParse = "couldnt parse to / from json."
	BadReq   = "couldnt create request."
	BadRes   = "couldnt receive request."
	BadConn  = "couldnt connect websocket."
	BadRead  = "couldnt read response body."
	BadWrite = "couldnt write body."
	BadOpen  = "couldnt open file."
)
