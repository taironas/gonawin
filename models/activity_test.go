package models

import (
	"testing"
	"time"

	"appengine/aetest"

	"github.com/santiaago/purple-wing/helpers/log"
)

func TestPublishActivity(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	log.Infof(c, "Test Publish Activity")

	tests := []struct {
		name      string
		activity  Activity
		want      *Activity
	}{
		{
			name: "Simple publish",
			activity: Activity{
				Type:       "Team",
				Verb:       "created",
				Actor:      ActivityEntity { 1, "user", "John Smith"},
				Object:     ActivityEntity { 10, "team", "TeamFoo"},
				Target:     ActivityEntity { 100, "foo", "TargetFoo"},
        Published:  time.Now(),
        UserID:     1,
			},
			want: &Activity{
				Type:       "Team",
				Verb:       "created",
				Actor:      ActivityEntity { 1, "user", "John Smith"},
				Object:     ActivityEntity { 10, "team", "TeamFoo"},
				Target:     ActivityEntity { 100, "foo", "TargetFoo"},
				Published:  time.Now(),
				UserID:     1,
			},
		},
	}
	for _, test := range tests {
		err := Publish(c, test.activity.Type, test.activity.Verb, test.activity.Actor, test.activity.Object, test.activity.Target, test.activity.UserID)
		if err != nil {
			t.Errorf("TestPublishActivity(%q): got '%v' error wanted no err", test.name, err)
		}
	}
}
