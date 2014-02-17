## Some thoughts concerning the JSON API:

### Ranking:

* Global
* scope of a team
* scope of a tournament
* dedicated page for leaderboards (ranking sorted by geographical location)

### API:

### Ranking: 
    url:
        * j/tournaments/:id/ranking?count=50
        * j/teams/:id/ranking
        * j/users/ranking
    parameters:
        * count: specify the number of user to retrieve. If count is not present, then default value is 10.
    description: ranking urls where an array of user (id, username, score) sorted by score will be returned.
