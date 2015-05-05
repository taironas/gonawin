package extract

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

// Context type holds the information needed to read the request and log any errors.
//
type Context struct {
	c    appengine.Context // appengine context
	desc string            // handler description
	r    *http.Request     // the HTTP request
}

func NewContext(c appengine.Context, desc string, r *http.Request) Context {
	return Context{c, desc, r}
}

// User returns a User object from the request passed in the context.
//
func (c Context) User() (*mdl.User, error) {

	strUserId, err := route.Context.Get(c.r, "userId")
	if err != nil {
		log.Errorf(c.c, "%s error getting user id, err:%v", c.desc, err)
		return nil, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}

	var userId int64
	userId, err = strconv.ParseInt(strUserId, 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error converting user id from string to int64, err:%v", c.desc, err)
		return nil, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}

	var u *mdl.User
	u, err = mdl.UserById(c.c, userId)
	log.Infof(c.c, "%s User: %v", c.desc, u)
	if err != nil {
		log.Errorf(c.c, "%s user not found", c.desc)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}

	return u, nil
}

// TeamId returns the team identifier.
// It gets the 'teamId' from the request and parses it to int64
//
func (c Context) TeamId() (int64, error) {
	strTeamId, err := route.Context.Get(c.r, "teamId")
	if err != nil {
		log.Errorf(c.c, "%s error getting team id, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}

	var teamId int64
	teamId, err = strconv.ParseInt(strTeamId, 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error converting team id from string to int64, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}
	return teamId, nil
}

// Team returns a team object from the request passed in the context.
//
func (c Context) Team() (*mdl.Team, error) {

	teamId, err := c.TeamId()
	if err != nil {
		return nil, err
	}

	var t *mdl.Team
	t, err = mdl.TeamById(c.c, teamId)
	if err != nil {
		log.Errorf(c.c, "%s team with id:%v was not found %v", c.desc, teamId, err)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}
	return t, nil
}

// RequestId returns a int64 requestId from the HTTP request.
//
func (c Context) RequestId() (int64, error) {

	strRequestId, err := route.Context.Get(c.r, "requestId")
	if err != nil {
		log.Errorf(c.c, "%s error getting request id, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
	}

	var requestId int64
	requestId, err = strconv.ParseInt(strRequestId, 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error converting request id from string to int64, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
	}
	return requestId, nil
}

// TeamRequest returns a a request to join a team from an HTTP request
//
func (c Context) TeamRequest() (*mdl.TeamRequest, error) {

	requestId, err := c.RequestId()
	if err != nil {
		return nil, err
	}

	var teamRequest *mdl.TeamRequest
	if teamRequest, err = mdl.TeamRequestById(c.c, requestId); err != nil {
		log.Errorf(c.c, "%s teams.DenyRequest, team request not found: %v", c.desc, err)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
	}
	return teamRequest, nil
}

// TournamentId returns the id of the tournament that the request holds.
//
func (c Context) TournamentId() (int64, error) {
	strTournamentId, err := route.Context.Get(c.r, "tournamentId")
	if err != nil {
		log.Errorf(c.c, "%s error getting tournament id, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
	}

	var tournamentId int64
	tournamentId, err = strconv.ParseInt(strTournamentId, 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error converting tournament id from string to int64, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
	}
	return tournamentId, nil
}

// Tournament returns a tournament instance.
// It gets the 'tournamentId' from the request and queries the datastore to get
// the tournament.
//
func (c Context) Tournament() (*mdl.Tournament, error) {

	tournamentId, err := c.TournamentId()
	if err != nil {
		return nil, err
	}

	var tournament *mdl.Tournament
	if tournament, err = mdl.TournamentById(c.c, tournamentId); err != nil {
		log.Errorf(c.c, "%s tournament not found: %v", c.desc, err)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
	}
	return tournament, nil
}
