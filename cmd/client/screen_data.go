package client

import "netBlast/pkg"

type option struct {
	name  string
	value string
}

// Additional model for userlist screen
type userList struct {
	users []pkg.User
}

// Additional model for settings screen
type settings struct {
	options []option
}

// Additional model for help screen
type help struct {
	options []option
}
