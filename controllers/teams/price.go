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

package teams

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"appengine"

	"github.com/santiaago/gonawin/extract"
	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

// Prices handler, use it to get the team's prices.
//  POST	/j/teams/[0-9]+/prices/     Get the prices of a team team with the given team id.
// Reponse: array of JSON formatted prices.
//
func Prices(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Prices Handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error

	team, err = extract.Team()
	if err != nil {
		return err
	}

	prices := team.Prices(c)

	pvm := buildTeamPricesViewModel(prices)

	return templateshlp.RenderJson(w, c, pvm)
}

type teamPricesViewModel struct {
	Prices []*mdl.Price
}

func buildTeamPricesViewModel(prices []*mdl.Price) teamPricesViewModel {
	return teamPricesViewModel{Prices: prices}
}

// PriceByTournament, use it to list the prices of a team for a specific tournament.
//
// Use this handler to get the price of a team for a specific tournament.
//	GET	/j/teams/[0-9]+/prices/[0-9]+/	Retreives price of a team with the given id for the specified tournament.
//
func PriceByTournament(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method == "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Prices by tournament Handler:"
	extract := extract.NewContext(c, desc, r)

	var t *mdl.Team
	var err error
	t, err = extract.Team()
	if err != nil {
		return err
	}

	var tournamentId int64
	tournamentId, err = extract.TournamentId()
	if err != nil {
		return err
	}

	log.Infof(c, "%s ready to get price", desc)
	p := t.PriceByTournament(c, tournamentId)

	data := struct {
		Price *mdl.Price
	}{
		p,
	}
	return templateshlp.RenderJson(w, c, data)
}

// UpdatePrice handler, use it to update the price of a team for a specific tournament.
//
// Use this handler to get the price of a team for a specific tournament.
//	GET	/j/teams/[0-9]+/prices/update/[0-9]+/	Update Retreives price of a team with the given id for the specified tournament.
//
func UpdatePrice(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team update price Handler:"
	extract := extract.NewContext(c, desc, r)

	var t *mdl.Team
	var err error
	t, err = extract.Team()
	if err != nil {
		return err
	}

	var tournamentId int64
	tournamentId, err = extract.TournamentId()
	if err != nil {
		return err
	}

	log.Infof(c, "%s ready to get price", desc)
	p := t.PriceByTournament(c, tournamentId)

	// only work on name and private. Other values should not be editable
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "%s Error when reading request body err: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
	}

	var priceData PriceData
	err = json.Unmarshal(body, &priceData)
	if err != nil {
		log.Errorf(c, "%s Error when decoding request body err: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
	}

	if helpers.IsStringValid(priceData.Description) && (p.Description != priceData.Description) {
		// update data
		p.Description = priceData.Description
		p.Update(c)
	} else {
		log.Errorf(c, "%s Cannot update because updated is not valid.", desc)
		log.Errorf(c, "%s Update description = %s", desc, priceData.Description)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
	}

	data := struct {
		Price *mdl.Price
	}{
		p,
	}
	return templateshlp.RenderJson(w, c, data)
}
