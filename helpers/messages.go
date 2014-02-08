/*
* Copyright (c) 2013 Santiago Arias | Remy Jourde | Carlos Bernal
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

package helpers

import ()


// errors
const (
	ErrorCodeNotSupported = "Not Supported"
	ErrorCodeInternal = "Internal error"

	ErrorCodeNameCannotBeEmpty = "Name field cannot be empty"

	// teams
	ErrorCodeTeamAlreadyExists        = "Sorry, that team already exists"
	ErrorCodeTeamCannotCreate         = "Could not create the team"
	ErrorCodeTeamNotFound             = "Team not found"
	ErrorCodeTeamNotFoundCannotUpdate = "Team not found, unable to update"
	ErrorCodeTeamNotFoundCannotDelete = "Team not found, unable to delete"
	ErrorCodeTeamNotFoundCannotInvite = "Team not found, unable to send invitation"
	ErrorCodeTeamUpdateForbiden       = "Team can only be updated by the team administrator"
	ErrorCodeTeamDeleteForbiden       = "Team can only be deleted by the team administrator"
	ErrorCodeTeamCannotUpdate         = "Could not update team"
	ErrorCodeTeamCannotInvite         = "Could not send invitation"
	ErrorCodeTeamRequestNotFound      = "Request not found"
	ErrorCodeTeamMemberNotFound       = "Member not found"

	//tournaments
	ErrorCodeTournamentAlreadyExists        = "Sorry, that tournament already exists"
	ErrorCodeTournamentCannotCreate         = "Could not create the team"
	ErrorCodeTournamentNotFound             = "Tournament not found"
	ErrorCodeTournamentNotFoundCannotUpdate = "Tournament not found, unable to update"
	ErrorCodeTournamentNotFoundCannotDelete = "Tournament not found, unable to delete"
	ErrorCodeTournamentUpdateForbiden       = "Tournament can only be updated by the team administrator"
	ErrorCodeTournamentDeleteForbiden       = "Tournament can only be deleted by the team administrator"
	ErrorCodeTournamentCannotUpdate         = "Could not update tournament"
	ErrorCodeTournamentCannotSearch = "Something went wrong, we are unable to perform search query"
)

//info 
//const ()










