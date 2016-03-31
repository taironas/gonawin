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

// NewContext returns a new Context.
//
func NewContext(c appengine.Context, desc string, r *http.Request) Context {
	return Context{c, desc, r}
}

// UserID returns a userId.
// It gets the 'userId' from the request and parses it to int64
//
func (c Context) UserID() (int64, error) {

	strUserID, err := route.Context.Get(c.r, "userId")
	if err != nil {
		log.Errorf(c.c, "%s error getting user id, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}

	var userID int64
	userID, err = strconv.ParseInt(strUserID, 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error converting user id from string to int64, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}
	return userID, nil
}

// User returns a User object from the request passed in the context.
//
func (c Context) User() (*mdl.User, error) {

	userID, err := c.UserID()
	if err != nil {
		return nil, err
	}

	var u *mdl.User
	if u, err = mdl.UserById(c.c, userID); err != nil {
		log.Errorf(c.c, "%s user not found", c.desc)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}

	return u, nil
}

// Admin returns a admin mdl.User object with respect to the
// userId passed as param.
//
func (c Context) Admin(userID int64) (*mdl.User, error) {

	a, err := mdl.UserById(c.c, userID)
	if err != nil {
		log.Errorf(c.c, "%s user not found", c.desc)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}
	return a, nil
}

// TeamID returns the team identifier.
// It gets the 'teamID' from the request and parses it to int64
//
func (c Context) TeamID() (int64, error) {
	strTeamID, err := route.Context.Get(c.r, "teamId")
	if err != nil {
		log.Errorf(c.c, "%s error getting team id, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}

	var teamID int64
	teamID, err = strconv.ParseInt(strTeamID, 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error converting team id from string to int64, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}
	return teamID, nil
}

// Team returns a team object from the request passed in the context.
//
func (c Context) Team() (*mdl.Team, error) {

	teamID, err := c.TeamID()
	if err != nil {
		return nil, err
	}

	var t *mdl.Team
	t, err = mdl.TeamByID(c.c, teamID)
	if err != nil {
		log.Errorf(c.c, "%s team with id:%v was not found %v", c.desc, teamID, err)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}
	return t, nil
}

// RequestID returns a int64 requestId from the HTTP request.
//
func (c Context) RequestID() (int64, error) {

	strRequestID, err := route.Context.Get(c.r, "requestId")
	if err != nil {
		log.Errorf(c.c, "%s error getting request id, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
	}

	var requestID int64
	requestID, err = strconv.ParseInt(strRequestID, 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error converting request id from string to int64, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
	}
	return requestID, nil
}

// TeamRequest returns a a request to join a team from an HTTP request
//
func (c Context) TeamRequest() (*mdl.TeamRequest, error) {

	requestID, err := c.RequestID()
	if err != nil {
		return nil, err
	}

	var teamRequest *mdl.TeamRequest
	if teamRequest, err = mdl.TeamRequestByID(c.c, requestID); err != nil {
		log.Errorf(c.c, "%s teams.DenyRequest, team request not found: %v", c.desc, err)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
	}
	return teamRequest, nil
}

// TournamentID returns the ID of the tournament that the request holds.
//
func (c Context) TournamentID() (int64, error) {
	strTournamentID, err := route.Context.Get(c.r, "tournamentId")
	if err != nil {
		log.Errorf(c.c, "%s error getting tournament id, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
	}

	var tournamentID int64
	tournamentID, err = strconv.ParseInt(strTournamentID, 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error converting tournament id from string to int64, err:%v", c.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
	}
	return tournamentID, nil
}

// Tournament returns a tournament instance.
// It gets the 'tournamentId' from the request and queries the datastore to get
// the tournament.
//
func (c Context) Tournament() (*mdl.Tournament, error) {

	tournamentID, err := c.TournamentID()
	if err != nil {
		return nil, err
	}

	var tournament *mdl.Tournament
	if tournament, err = mdl.TournamentById(c.c, tournamentID); err != nil {
		log.Errorf(c.c, "%s tournament not found: %v", c.desc, err)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
	}
	return tournament, nil
}

// Match returns a match instance.
//
func (c Context) Match(tournament *mdl.Tournament) (*mdl.Tmatch, error) {

	strmatchIDNumber, err := route.Context.Get(c.r, "matchId")
	if err != nil {
		log.Errorf(c.c, "%s error getting match id, err:%v", c.desc, err)
		return nil, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeMatchCannotUpdate)}
	}

	var matchIDNumber int64
	matchIDNumber, err = strconv.ParseInt(strmatchIDNumber, 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error converting match id from string to int64, err:%v", c.desc, err)
		return nil, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeMatchCannotUpdate)}
	}

	match := mdl.GetMatchByIDNumber(c.c, *tournament, matchIDNumber)
	if match == nil {
		log.Errorf(c.c, "%s unable to get match with id number :%v", c.desc, matchIDNumber)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeMatchNotFoundCannotUpdate)}
	}
	return match, nil
}

// Count extracts the 'count' value from the given http.Request
// returns 20 if none is found.
//
func (c Context) Count() int64 {

	defaultValue := int64(20) // set count to default value
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

// CountOrDefault extracts the 'count' value from the given http.Request
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
		return int64(1)
	}

	page, err := strconv.ParseInt(c.r.FormValue("page"), 0, 64)
	if err != nil {
		log.Errorf(c.c, "%s error during conversion of page parameter: %v", c.desc, err)
		page = 1
	}
	return page
}
