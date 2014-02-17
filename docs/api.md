## Some thoughts concerning the JSON API:

### Ranking:

* Global
* scope of a team
* scope of a tournament
* dedicated page for leaderboards (ranking sorted by geographical location)

### API:

### Ranking: 
    url: /j/ranking/
    parameters:
      * scope: global, team, tournament
      * id: specify the team/tournament id
      * count: specify the number of user to retrieve
    description: ranking url where an array of user (id, username, score) sorted by score will be returned.
    ex1: /j/ranking?scope=global&count=50 (return 50 first users sorted by score)
    ex2: /j/ranking?scope=tournament&id=5481&count=10 (return 10 first users for a given tournament)
