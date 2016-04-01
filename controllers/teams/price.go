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

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// Prices handler, use it to get the team's prices.
//  GET	/j/teams/[0-9]+/prices/     Get the prices of a team team with the given team id.
// Response: array of JSON formatted prices.
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

	return templateshlp.RenderJSON(w, c, pvm)
}

type teamPricesViewModel struct {
	Prices []*mdl.Price
}

func buildTeamPricesViewModel(prices []*mdl.Price) teamPricesViewModel {
	return teamPricesViewModel{Prices: prices}
}

// PriceByTournament handler, use it to get the price of a team for a specific tournament.
//	GET	/j/teams/[0-9]+/prices/[0-9]+/		Retrieves price of a team with the given id for the specified tournament.
// Response: JSON formatted price.
//
func PriceByTournament(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method == "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Price by tournament Handler:"
	extract := extract.NewContext(c, desc, r)

	var t *mdl.Team
	var err error
	t, err = extract.Team()
	if err != nil {
		return err
	}

	var tournamentID int64
	tournamentID, err = extract.TournamentID()
	if err != nil {
		return err
	}

	p := t.PriceByTournament(c, tournamentID)

	pvm := buildTeamPriceViewModel(p)

	return templateshlp.RenderJSON(w, c, pvm)
}

// UpdatePrice handler, use it to update the price of a team for a specific tournament.
//	POST	/j/teams/[0-9]+/prices/update/[0-9]+/		Updates price of a team with the given id for the specified tournament.
// Response: JSON formatted price.
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

	var tournamentID int64
	tournamentID, err = extract.TournamentID()
	if err != nil {
		return err
	}

	p := t.PriceByTournament(c, tournamentID)

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

	pvm := buildTeamPriceViewModel(p)

	return templateshlp.RenderJSON(w, c, pvm)
}

type teamPriceViewModel struct {
	Price *mdl.Price
}

func buildTeamPriceViewModel(price *mdl.Price) teamPriceViewModel {
	return teamPriceViewModel{Price: price}
}
