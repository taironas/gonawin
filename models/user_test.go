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
		title    string
		email    string
		username string
		name     string
		alias    string
		isAdmin  bool
		auth     string
		err      string
	}{
		{
			title:    "can create user",
			email:    "foo@bar.com",
			username: "john.snow",
			name:     "john snow",
			alias:    "",
			isAdmin:  false,
			auth:     "",
			err:      "",
		},
	}

	for i, test := range tests {
		t.Log(test.title)
		var got *User
		if got, err = CreateUser(c, test.email, test.username, test.name, test.alias, test.isAdmin, test.auth); err != nil {
			t.Errorf("test %v - Error: %v", i, err)
		}
		if got.Email != test.email {
			t.Errorf("test %v - Error; want Email == %s, got %s", i, test.email, got.Email)
		}
		if got.Username != test.username {
			t.Errorf("test %v - Error; want Username == %s, got %s", i, test.username, got.Username)
		}
		if got.Name != test.name {
			t.Errorf("test %v - Error; want Name == %s, got %s", i, test.name, got.Name)
		}
		if got.Alias != test.alias {
			t.Errorf("test %v - Error; want Name == %s, got %s", i, test.alias, got.Alias)
		}
		if got.IsAdmin != test.isAdmin {
			t.Errorf("test %v - Error; want isAdmin == %s, got %s", i, test.isAdmin, got.IsAdmin)
		}
		if got.Auth != test.auth {
			t.Errorf("test %v - Error; want auth == %s, got %s", i, test.auth, got.Auth)
		}
	}
}
