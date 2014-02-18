## Some thoughts concerning the JSON API:

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
