package gonawintest

import (
	"appengine/aetest"

	mdl "github.com/taironas/gonawin/models"
)

// CreateWorldCup creates a world cup tournament into the datastore
func CreateWorldCup(c aetest.Context, adminID int64) *mdl.Tournament {
	tournament, _ := mdl.CreateWorldCup(c, adminID)

	return tournament
}

// PredictOnMatch add a predict to given match and for a given user
func PredictOnMatch(c aetest.Context, user *mdl.User, matchID, result1, result2 int64) error {
	if predict, err := mdl.CreatePredict(c, user.Id, result1, result2, matchID); err != nil {
		return err
	}
	// add p.Id to User predict table.
	if err = user.AddPredictId(c, predict.Id); err != nil {
		return err
	}

	return nil
}
