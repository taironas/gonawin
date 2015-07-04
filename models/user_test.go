package models

import (
	"testing"

	"appengine/aetest"
)

func TestCreateUser(t *testing.T) {
	var c aetest.Context
	var err error
	if c, err = aetest.NewContext(nil); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tests := []struct {
		email    string
		username string
		name     string
		alias    string
		isAdmin  bool
		auth     string
	}{
		{
			email:    "foo@bar.com",
			username: "john.snow",
			name:     "john snow",
			alias:    "",
			isAdmin:  false,
			auth:     "",
		},
	}

	for i, test := range tests {
		var got *User
		if got, err = CreateUser(c, test.email, test.username, test.name, test.alias, test.isAdmin, test.auth); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if got.Email != test.email {
			t.Errorf("test %v - Error; want Email == %s, got %s", i, test.email, got.Email)
		}
	}
}
