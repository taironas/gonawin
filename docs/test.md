Issue #200 Testing
===================

| Num.  |          m path                  |              ng# path                | status |     comments                    |
-------|:----------------------------------|:-------------------------------------|:-----|:----------------------------
| 1     |    /                             | /ng#                                 | &#x2713;   | login & logout ok
| 2     |    /m/about                      | /ng#/about                           | &#x2713;   |
| 3     |    /m/contact                    | /ng#/contact                         | &#x2713;   |
| 4     |    /m/auth                       | none                                 | N.A. | done at root
| 5     |    /m/auth/facebook              | none                                 | N.A. | handled by javascript
| 6     |    /m/auth/facebook/callback     | none                                 | N.A. | handled by javascript
| 7     |    /m/auth/google                | none                                 | &#x2713; | handled by javascript
| 8     |    /m/auth/google/callback       | none                                 | &#x2713; | handled by javascript
| 9     |    /m/auth/twitter               | none				                          | &#x2713; | handled by server
| 10    |    /m/auth/twitter/callback      | none				                          | &#x2713; | handled by server & client
| 11    |    /m/logout                     | none	 			                          | &#x2713; | handled by javascript
| 12    |    /m/users/[0-9]+      	       | /ng#/users/show/:id		              |	&#x2713;	 | [issue #249](https://github.com/santiaago/purple-wing/issues/249) -> closed
| 13    |    /m/a			                     | none				                          | N.A. | path is no longer necessary
| 14    |    /m/a/users			               | /ng#/users/			                    |	&#x2713;	 | [issue #250](https://github.com/santiaago/purple-wing/issues/250) -> closed
| 15    |    /m/teams			                 | /ng#/teams/			                    | &#x2713;	 |
| 16    |    /m/teams/new		               | /ng#/teams/new			                  |	&#x2713;   | 
| 17    |    /m/teams/[0-9]+		           | /ng#/teams/show/:id	  	            |	&#x2713;	 | [issue #251](https://github.com/santiaago/purple-wing/issues/251) -> closed
| 18    |    /m/teams/[0-9]+/edit	         | /ng#/teams/edit/:id		              |	&#x2713;	 | [issue #251](https://github.com/santiaago/purple-wing/issues/251) -> closed
| 19    |    not defined                   | /ng#/teams/search			              |	&#x2713;	 |
| 20    |    /m/teams/destroy/[0-9]+	     | replaced by delete team				      | &#x2713;   |
| 21    |    /m/teams/[0-9]+/invite        | invite button on team page		        | &#x2713;	 |
| 22    |    /m/teams/[0-9]+/request	     | allow/deny button		                | &#x2713;	 |
| 22    |    /m/tournaments		             | /ng#/tournaments/			              | &#x2713;   |
| 23    |    /m/tournaments/new		         | /ng#/tournaments/new		              |	&#x2713;   |
| 24    |    /m/tournaments/[0-9]+         | /ng#/tournaments/show/:id		        |	&#x2713;	 | [issue #252](https://github.com/santiaago/purple-wing/issues/252) -> closed
| 25    |    /m/tournaments/[0-9]+/edit    | /ng#/tournaments/edit/:id		        |	&#x2713;	 | [issue #252](https://github.com/santiaago/purple-wing/issues/252) -> closed
| 26    |    not defined		               | /ng#/tournaments/search		          | &#x2713;	 |
| 27    |    /m/tournaments/destroy/[0-9]+ | replaced by delete tournament        | &#x2713;   |
| 28    |    /m/teamrels/create            | none				                          | &#x2713;	 | relation created when join team or when create team
| 29    |    /m/teamrels/destroy	         | none				                          | &#x2713;   | relation created when leave team or when delete team
| 30    |    /m/tournamentrels/create      | none				                          | &#x2713;	 | relation created when join tournament
| 31    |    /m/tournamentrels/destroy     | none				                          | &#x2713;	 | relation created when leave tournament or when delete tournament
| 32    |    /m/tournamentteamrels/create  | none				                          | &#x2713;	 | relation created when join tournament as team
| 33    |    /m/tournamentteamrels/destroy | none				                          | &#x2713;	 | relation created when leave tournament as team or when delete tournament
| 34    |    /m/settings/edit-profile      | /ng#/settings/edit-profile           |	&#x2713;	 |
| 35    |    /m/settings/networks          | /ng#/settings/network		            | &#x2713; 	 |
| 36    |    /m/settings/email		         | /ng#/settings/email		              |	&#x2713;   |
| 37    |    /m/invite                     | /ng#/invite			                    |	&#x2713;   |
