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
	// generic
	ErrorCodeNotSupported      = "Not Supported"
	ErrorCodeInternal          = "Internal error"
	ErrorCodeNotFound          = "Not Found"
	ErrorCodeNameCannotBeEmpty = "Name field cannot be empty"

	// sessions
	ErrorCodeSessionsAccessTokenNotValid      = "Access token is not valid"
	ErrorCodeSessionsForbiden                 = "You are not authorized to log in to gonawin"
	ErrorCodeSessionsUnableToSignin           = "Error occurred during signin process"
	ErrorCodeSessionsCannotGetTempCredentials = "Error getting temporary credentials"
	ErrorCodeSessionsCannotGetSecretValue     = "Error getting 'secret' value"
	ErrorCodeSessionsCannotGetRequestToken    = "Error getting request token"
	ErrorCodeSessionsCannotGetUserInfo        = "Error getting user info from Twitter"

	// users
	ErrorCodeUserNotFound             = "User not found"
	ErrorCodeUserNotFoundCannotUpdate = "User not found, unable to update"
	ErrorCodeUserCannotUpdate         = "Could not update user"
	ErrorCodeUsersCannotUpdate        = "Could not update users"
	ErrorCodeUsersCannotPublishScore  = "Could not pusblish score activities"
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
	ErrorCodeTeamAdminCannotLeave     = "Team administrator cannot leave the team"
	//tournaments
	ErrorCodeTournamentAlreadyExists          = "Sorry, that tournament already exists"
	ErrorCodeTournamentCannotCreate           = "Could not create the team"
	ErrorCodeTournamentNotFound               = "Tournament not found"
	ErrorCodeTournamentNotFoundCannotUpdate   = "Tournament not found, unable to update"
	ErrorCodeTournamentNotFoundCannotDelete   = "Tournament not found, unable to delete"
	ErrorCodeTournamentUpdateForbiden         = "Tournament can only be updated by the team administrator"
	ErrorCodeTournamentDeleteForbiden         = "Tournament can only be deleted by the team administrator"
	ErrorCodeTournamentCannotUpdate           = "Could not update tournament"
	ErrorCodeTournamentCannotSearch           = "Something went wrong, we are unable to perform search query"
	ErrorCodeMatchCannotUpdate                = "Something went wrong, unable to update match"
	ErrorCodeMatchesCannotUpdate              = "Something went wrong, unable to update matches"
	ErrorCodeMatchNotFoundCannotUpdate        = "Match not found, unable to update match"
	ErrorCodeMatchNotFound                    = "Match not found"
	ErrorCodeMatchNotFoundCannotSetPrediction = "Match not found, unable to set prediction"
	ErrorCodeCannotSetPrediction              = "Something went wrong, unable to set prediction"
	ErrorCodeNotAllowedToSetPrediction        = "You have to join the tournament to be able to set a predict for this match"
	ErrorCodeTeamsCannotUpdate                = "Could not update teams"

	// invite
	ErrorCodeInviteNoEmailAddr     = "No email address has been entered"
	ErrorCodeInviteEmailsInvalid   = "Emails list is not properly formatted"
	ErrorCodeInviteEmailCannotSend = "Sorry, we were unable to send the Email"

	// relations

)

//info
//const ()
