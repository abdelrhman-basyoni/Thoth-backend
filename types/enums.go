package types

type roles struct {
	Admin  string
	Author string
}

var Roles = &roles{Admin: "ADMIN", Author: "AUTHOR"}
