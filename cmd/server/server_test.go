package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckNames(t *testing.T) {
	s := newServer()

	s.users["Batman"] = ""

	// Test cases
	tests := []struct {
		test   string
		name   string
		errMsg string
	}{
		{
			test:   "Should output that the name is too short",
			name:   "bd",
			errMsg: "Name too short. ",
		},
		{
			test:   "Should output that the name is too long",
			name:   "Bobby Vance",
			errMsg: "Name too long. ",
		},
		{
			test:   "Should output that the name already exists",
			name:   "Batman",
			errMsg: "Name already exists. ",
		},
		{
			test:   "Should validate the name",
			name:   "Bob",
			errMsg: "",
		},
	}

	// Tests the provided cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err, _ := checkName(test.name, s)
			assert.Equal(t, test.errMsg, err)
		})
	}
}
