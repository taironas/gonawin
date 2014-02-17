## Some thoughts concerning the world cup tournament:

### World cup phases:

* Round of 16
* Quarter-finals
* Semi-finals
* Round 19 (Match for third plac)
* Finals

### world cup views
we want to have different views of the world cup tournament:

* global view of the tournaments.
* group view (first phase) where we show the rank of each group.
* phase view where we show the matches grouped by phases.
* calendar view where we show the matches grouped by dates.


`url: /tournament/:id` 

main url where we display the world cup groups ranked by points and the braket with the current rank.
needs:
* get groups ranked by points
* get the braket based on the current rank

`url: /tournament/:id/calendar`

overall view of the matches. we will display all matches (first phase and second phase).
we should offer the possibility to group the matches by phases of by dates.
needs:
* get all matches grouped by day
* get all matches grouped by phases

`url: /tournament/:id/matches/first_stage`

display all tournament matches grouped by group (first phase)
'/matches' will redirect to 'matches/first_stage'

`url: /tournament/:id/matches/second_stage`

display bracket then remaining matches grouped by phases
