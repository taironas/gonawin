## Some thoughts concerning the JSON API:

### Ranking:

* Global
* scope of a team
* scope of a tournament
* dedicated page for leaderboards (ranking sorted by geographical location)

### API:

`url: /ranking/
parameters: 
* scope: global, team, tournament
* id: specify the team/tournament id
* count: specify the number of user to retrieve`

ranking url where an array of user (id, username, score) sorted by score will be returned.
