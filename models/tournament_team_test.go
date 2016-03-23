package models

import (
	"testing"

	"appengine/aetest"

	"github.com/taironas/gonawin/helpers/log"
)

func TestTTeamByID(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	log.Infof(c, "Tteam TeamByID")

	tests := []struct {
		name string
		id   int64
		want *Tteam
	}{
		{
			name: "Nil entity",
			id:   int64(0),
			want: nil,
		},
	}
	for _, test := range tests {
		got, _ := TTeamByID(c, test.id)
		if got == nil && test.want != nil {
			t.Errorf("TTeamByID(%q): got nil wanted %v", test.name, *test.want)
		} else if got != nil && test.want == nil {
			t.Errorf("TTeamByID(%q): got %v wanted nil", test.name, *got)
		} else if got == nil && test.want == nil {
			// This is OK
		} else if *got != *test.want {
			t.Errorf("TTeamByID(%q): got %v wanted %v", test.name, *got, *test.want)
		}
	}
}
