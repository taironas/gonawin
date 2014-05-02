package models

import (
	"testing"
	"time"

	"appengine/aetest"

	"github.com/santiaago/gonawin/helpers/log"
)

func TestSaveActivity(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	log.Infof(c, "Test Save Activity")

	tests := []struct {
		name     string
		activity Activity
		want     *Activity
	}{
		{
			name: "Simple save",
			activity: Activity{
				Type:      "Team",
				Verb:      "created",
				Actor:     ActivityEntity{1, "user", "John Smith"},
				Object:    ActivityEntity{10, "team", "TeamFoo"},
				Target:    ActivityEntity{100, "foo", "TargetFoo"},
				Published: time.Now(),
				UserID:    1,
			},
			want: &Activity{
				Type:      "Team",
				Verb:      "created",
				Actor:     ActivityEntity{1, "user", "John Smith"},
				Object:    ActivityEntity{10, "team", "TeamFoo"},
				Target:    ActivityEntity{100, "foo", "TargetFoo"},
				Published: time.Now(),
				UserID:    1,
			},
		},
	}
	for _, test := range tests {
		err := test.activity.save(c)
		if err != nil {
			t.Errorf("TestSaveActivity(%q): got '%v' error wanted no err", test.name, err)
		}
	}
}
