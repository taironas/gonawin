package tournaments

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	mdl "github.com/santiaago/gonawin/models"
	"github.com/taironas/route"

	"appengine"
)

// requestContext type holds the information needed read the request and log any errors.
type requestContext struct {
	c    appengine.Context // appengine context
	desc string            // handler description
	r    *http.Request     // the HTTP request
}

// tournament returns a tournament instance.
// It gets the 'tournamentId' from the request and queries the datastore to get
// the tournament.
func (rc requestContext) tournament() (*mdl.Tournament, error) {

	strTournamentId, err := route.Context.Get(rc.r, "tournamentId")
	if err != nil {
		log.Errorf(rc.c, "%s error getting tournament id, err:%v", rc.desc, err)
		return nil, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
	}

	var tournamentId int64
	tournamentId, err = strconv.ParseInt(strTournamentId, 0, 64)
	if err != nil {
		log.Errorf(rc.c, "%s error converting tournament id from string to int64, err:%v", rc.desc, err)
		return nil, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
	}

	var tournament *mdl.Tournament
	if tournament, err = mdl.TournamentById(rc.c, tournamentId); err != nil {
		log.Errorf(rc.c, "%s tournament not found: %v", rc.desc, err)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
	}
	return tournament, nil
}

// userId returns a userId.
// It gets the 'userId' from the request and parses it to int64
func (rc requestContext) userId() (int64, error) {

	strUserId, err := route.Context.Get(rc.r, "userId")
	if err != nil {
		log.Errorf(rc.c, "%s error getting user id, err:%v", rc.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}

	var userId int64
	userId, err = strconv.ParseInt(strUserId, 0, 64)
	if err != nil {
		log.Errorf(rc.c, "%s error converting user id from string to int64, err:%v", rc.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}
	return userId, nil
}

// teamId returns the team identifier.
// It gets the 'teamId' from the request and parses it to int64
func (rc requestContext) teamId() (int64, error) {
	strTeamId, err := route.Context.Get(rc.r, "teamId")
	if err != nil {
		log.Errorf(rc.c, "%s error getting team id, err:%v", rc.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}

	var teamId int64
	teamId, err = strconv.ParseInt(strTeamId, 0, 64)
	if err != nil {
		log.Errorf(rc.c, "%s error converting team id from string to int64, err:%v", rc.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}
	return teamId, nil
}

// admin returns a admin mdl.User object with respect to the
// userId passed as param.
func (rc requestContext) admin(userId int64) (*mdl.User, error) {

	a, err := mdl.UserById(rc.c, userId)
	log.Infof(rc.c, "%s User: %v", rc.desc, a)
	if err != nil {
		log.Errorf(rc.c, "%s user not found", rc.desc)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}
	return a, nil
}

func (rc requestContext) match(tournament *mdl.Tournament) (*mdl.Tmatch, error) {

	strmatchIdNumber, err := route.Context.Get(rc.r, "matchId")
	if err != nil {
		log.Errorf(rc.c, "%s error getting match id, err:%v", rc.desc, err)
		return nil, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeMatchCannotUpdate)}
	}

	var matchIdNumber int64
	matchIdNumber, err = strconv.ParseInt(strmatchIdNumber, 0, 64)
	if err != nil {
		log.Errorf(rc.c, "%s error converting match id from string to int64, err:%v", rc.desc, err)
		return nil, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeMatchCannotUpdate)}
	}

	match := mdl.GetMatchByIdNumber(rc.c, *tournament, matchIdNumber)
	if match == nil {
		log.Errorf(rc.c, "%s unable to get match with id number :%v", rc.desc, matchIdNumber)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeMatchNotFoundCannotUpdate)}
	}
	return match, nil
}
