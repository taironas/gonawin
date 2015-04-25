package teams

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

// team returns a pointer to the requested team.
func (rc requestContext) team() (*mdl.Team, error) {

	teamId, err := rc.teamId()
	if err != nil {
		return nil, err
	}

	var t *mdl.Team
	t, err = mdl.TeamById(rc.c, teamId)
	if err != nil {
		log.Errorf(rc.c, "%s team with id:%v was not found %v", rc.desc, teamId, err)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}
	return t, nil
}
