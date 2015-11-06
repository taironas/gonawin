package gonawintest

import (
	"testing"
	"time"

	"appengine/aetest"

	"github.com/taironas/gonawin/helpers/log"

	mdl "github.com/taironas/gonawin/models"
)

func TestActivitySaveActivities(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	log.Infof(c, "Test Activity.SaveActivities")

	tests := []struct {
		name     string
		activity mdl.Activity
		want     *mdl.Activity
	}{
		{
			name: "Simple save",
			activity: mdl.Activity{
				Type:      "Team",
				Verb:      "created",
				Actor:     mdl.ActivityEntity{1, "user", "John Smith"},
				Object:    mdl.ActivityEntity{10, "team", "TeamFoo"},
				Target:    mdl.ActivityEntity{100, "foo", "TargetFoo"},
				Published: time.Now(),
				CreatorID: 1,
			},
			want: &mdl.Activity{
				Type:      "Team",
				Verb:      "created",
				Actor:     mdl.ActivityEntity{1, "user", "John Smith"},
				Object:    mdl.ActivityEntity{10, "team", "TeamFoo"},
				Target:    mdl.ActivityEntity{100, "foo", "TargetFoo"},
				Published: time.Now(),
				CreatorID: 1,
			},
		},
	}

	for _, test := range tests {
		activities := []*mdl.Activity{&test.activity}
		err := mdl.SaveActivities(c, activities)
		if err != nil {
			t.Errorf("TestActivitySaveActivities(%s): got '%v' error wanted no err", test.name, err)
		}
	}
}
