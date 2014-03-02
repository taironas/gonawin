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

-------------

### Score API

#### User
A __User__ gains points each time he wins a prediction either by hitting the _exact score_ or by hitting the _trend score_. No points are give in any other way.
As a user can belong to multiple teams, even if it is for the same match, we centralize the __predicts__ so that a user can only have a single predict for a match.

##### User score formula:

[Todo]

User's __score__ is available in the __User__ url:

* `j/users/show/:id`

#### Team
In the same way as user,  a __Team__ will have a __score__.

#### Team score formula:

[]

Team's score is available in the __Teams__ url:
* `j/teams/show/:id`

A __Tournament__ will have a __User__ ranking and a __Team__ ranking.

* `j/tournaments/:id/rank?with=users`
* `j/tournaments/show/rank?with=teams`




