## some thoughts concerning the world cup tournament:

### world cup phases:

* Round of 16
* Quarter-finals
* Semi-finals
* Round 19 (Match for third place)
* Finals

### world cup views
we want to have different views of the world cup tournament:

* global view of the tournaments.
* group view (first phase) where we show the rank of each group.
* phase view where we show the matches grouped by phases.
* calendar view where we show the matches grouped by dates.

#### global view:
`url: /tournament/:id`

The main url where we display the world cup groups ranked by points and the braket with the current rank.

needs:
* get groups ranked by points
* get the braket based on the current rank

#### matches view:
`url: /tournament/:id/matches`

In the matches view will display all matches (first phase and second phase).

We can group the matches by `phases` or by `dates`.

By default we group by days. We you want group by phase use parameter `groupby` with values `phase` or `day`

- example: `/tournament/:id/matches?groupby=phase`

#### group view:
`url: /tournament/:id/matches/first_stage`

Display all tournament matches grouped by group (first phase)
`/matches` will redirect to `matches/first_stage`


#### phase view
`url: /tournament/:id/matches/second_stage`

Display bracket then remaining matches grouped by phases
