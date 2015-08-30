package extract

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	mdl "github.com/taironas/gonawin/models"
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

// UserId returns a userId.
// It gets the 'userId' from the request and parses it to int64
//
func (c Context) UserId() (int64, error) {

	strUserId, err := route.Context.Get(c.r, "userId")
	if err != nil {
		log.Errorf(c.c, "%s error getting user id, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}

	var userId int64
	userId, err = strconv.ParseInt(strUserId, 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error converting user id from string to int64, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}
	return userId, nil
}

// User returns a User object from the request passed in the context.
//
func (c Context) User() (*mdl.User, error) {

	userId, err := c.UserId()
	if err != nil {
		return nil, err
	}

	var u *mdl.User
	if u, err = mdl.UserById(c.c, userId); err != nil {
		log.Errorf(c.c, "%s user not found", c.desc)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}

	return u, nil
}

// Admin returns a admin mdl.User object with respect to the
// userId passed as param.
//
func (c Context) Admin(userId int64) (*mdl.User, error) {

	a, err := mdl.UserById(c.c, userId)
	if err != nil {
		log.Errorf(c.c, "%s user not found", c.desc)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}
	return a, nil
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

func (c Context) Match(tournament *mdl.Tournament) (*mdl.Tmatch, error) {

	strmatchIdNumber, err := route.Context.Get(c.r, "matchId")
	if err != nil {
		log.Errorf(c.c, "%s error getting match id, err:%v", c.desc, err)
		return nil, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeMatchCannotUpdate)}
	}

	var matchIdNumber int64
	matchIdNumber, err = strconv.ParseInt(strmatchIdNumber, 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error converting match id from string to int64, err:%v", c.desc, err)
		return nil, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeMatchCannotUpdate)}
	}

	match := mdl.GetMatchByIdNumber(c.c, *tournament, matchIdNumber)
	if match == nil {
		log.Errorf(c.c, "%s unable to get match with id number :%v", c.desc, matchIdNumber)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeMatchNotFoundCannotUpdate)}
	}
	return match, nil
}

// Count extracts the 'count' value from the given http.Request
// returns 20 if none is found.
//
func (c Context) Count() int64 {

	defaultValue := 20 // set count to default value
	if len(c.r.FormValue("count")) == 0 {
		return defaultValue
	}

	count, err := strconv.ParseInt(c.r.FormValue("count"), 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s: error during conversion of count parameter: %v", c.desc, err)
		count = defaultValue
	}
	return count
}

// Count extracts the 'count' value from the given http.Request
// returns 20 if none is found.
//
func (c Context) CountOrDefault(d int64) int64 {

	if len(c.r.FormValue("count")) == 0 {
		return d
	}

	count, err := strconv.ParseInt(c.r.FormValue("count"), 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s: error during conversion of count parameter: %v", c.desc, err)
		count = d // set count to default value
	}
	return count
}

// Page extracts the 'page' value from the given http.Request
// returns 1 if none is found.
//
func (c Context) Page() int64 {

	if len(c.r.FormValue("page")) == 0 {
		return 1
	}

	page, err := strconv.ParseInt(c.r.FormValue("page"), 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error during conversion of page parameter: %v", c.desc, err)
		page = 1
	}
	return page
}
