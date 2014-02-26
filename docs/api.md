## gonawin API:

### Ranking:

* global
* scope of a team
* scope of a tournament
* dedicated page for leaderboards (ranking sorted by geographical location)

### Ranking API: 
####urls:

* `j/tournaments/:id/ranking?count=50`
* `j/teams/:id/ranking`
* `j/users/ranking`

####parameters:

`count`: specify the number of user to retrieve. If `count` is not present, the default value is `10`.

####description:

The ranking urls will return an array of entities (tournaments, teams, users) sorted by the score.

-------------

### Predict API

You can predict a match who is part of a tournament. To do this you set two parameters `result1` and `result2`. This two parameters are the scores that the user predicts for a specific match. A match is between a `Team1` and a `Team2`. So the results go respectively to each team.

Use the following URL to post a predict on a match:
* `/j/tournaments/:id/matches/:matchId/predict?result1=:result1&result2=:result2`



