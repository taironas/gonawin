/*
 * Copyright (c) 2014 Santiago Arias | Remy Jourde
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package users

import (
	"errors"
	"net/http"
	"strconv"

	"appengine"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	"github.com/taironas/route"

	mdl "github.com/santiaago/gonawin/models"
)

type requestContext struct {
	c    appengine.Context
	desc string
	r    *http.Request
}

func (rc requestContext) userId() (int64, error) {

	strUserId, err := route.Context.Get(rc.r, "userId")
	if err != nil {
		log.Errorf(rc.c, "%s error getting user id, err:%v", rc.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFoundCannotUpdate)}
	}

	var userId int64
	userId, err = strconv.ParseInt(strUserId, 0, 64)
	if err != nil {
		log.Errorf(rc.c, "%s error converting user id from string to int64, err:%v", rc.desc, err)
		return 0, &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFoundCannotUpdate)}
	}
	return userId, nil
}

func (rc requestContext) user() (*mdl.User, error) {

	userId, err := rc.userId()
	if err != nil {
		return nil, err
	}

	var user *mdl.User
	user, err = mdl.UserById(rc.c, userId)
	if err != nil {
		log.Errorf(rc.c, "%s user not found", rc.desc)
		return nil, &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}
	return user, nil
}