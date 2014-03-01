package models

import (
	"testing"

	"appengine/aetest"

	"github.com/santiaago/purple-wing/helpers/log"
)

func TestTeamById(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// Run code and tests requiring the appengine.Context using c.
	log.Infof(c, "Team Test")

	if _, err := TeamById(c, int64(0)); err != nil {
		t.Log(err)
	}
}
