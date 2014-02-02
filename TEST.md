Issue #200 Testing
===================

| Num.  |          m path                  |              ng# path                | status |     comments                    |
-------|:----------------------------------|:-------------------------------------|:-----|:----------------------------
| 1     |    /                             | /ng#                                 | OK   | login & logout ok
| 2     |    /m/about                      | /ng#/about                           | OK   |
| 3     |    /m/contact                    | /ng#/contact                         | OK   |
| 4     |    /m/auth                       | none                                 | N.A. | done at root
| 5     |    /m/auth/facebook              | none                                 | N.A. | handled by javascript
| 6     |    /m/auth/facebook/callback     | none                                 | N.A. | handled by javascript
| 7     |    /m/auth/google                | none                                 | N.A. | handled by javascript
| 8     |    /m/auth/google/callback       | none                                 | N.A. | handled by javascript
| 9     |    /m/auth/twitter               | none				                          | N.A. | handled by javascript
| 10    |    /m/auth/twitter/callback      | none				                          | N.A. | handled by javascript
| 11    |    /m/logout                     | none	 			                          | N.A. | handled by javascript
| 12    |    /m/users/[0-9]+      	       | /ng#/users/show/:id		              |	KO	 | issue #249
| 13    |    /m/a			                     | none				                          | N.A. | path is no longer necessary
| 14    |    /m/a/users			               | /ng#/users/			                    |	KO	 | issue #250
| 15    |    /m/teams			                 | /ng#/teams/			                    | OK	 |
| 16    |    /m/teams/new		               | /ng#/teams/new			                  |	OK   | 
| 17    |    /m/teams/[0-9]+		           | /ng#/teams/show/:id	  	            |	KO	 | issue #251
| 18    |    /m/teams/[0-9]+/edit	         | /ng#/teams/edit/:id		              |	KO	 | issue #251
| 19    |    not defined                   | /ng#/teams/search			              |	OK	 |
| 20    |    /m/teams/destroy/[0-9]+	     | replaced by delete team				      | OK   |
| 21    |    /m/teams/[0-9]+/invite        | invite button on team page		        | OK	 |
| 22    |    /m/teams/[0-9]+/request	     | allow/deny button		                | OK	 |
| 22    |    /m/tournaments		             | /ng#/tournaments/			              | OK   |
| 23    |    /m/tournaments/new		         | /ng#/tournaments/new		              |	OK   |
| 24    |    /m/tournaments/[0-9]+         | /ng#/tournaments/show/:id		        |	KO	 | issue #252
| 25    |    /m/tournaments/[0-9]+/edit    | /ng#/tournaments/edit/:id		        |	KO	 | issue #252
| 26    |    not defined		               | /ng#/tournaments/search		          | OK	 |
| 27    |    /m/tournaments/destroy/[0-9]+ | replaced by delete tournament        | OK   |
| 28    |    /m/teamrels/create            | none				                          | OK	 | relation created when join team or when create team
| 29    |    /m/teamrels/destroy	         | none				                          | OK   | relation created when leave team or when delete team
| 30    |    /m/tournamentrels/create      | none				                          | OK	 | relation created when join tournament
| 31    |    /m/tournamentrels/destroy     | none				                          | OK	 | relation created when leave tournament or when delete tournament
| 32    |    /m/tournamentteamrels/create  | none				                          | OK	 | relation created when join tournament as team
| 33    |    /m/tournamentteamrels/destroy | none				                          | OK	 | relation created when leave tournament as team or when delete tournament
| 34    |    /m/settings/edit-profile      | /ng#/settings/edit-profile           |	OK	 |
| 35    |    /m/settings/networks          | /ng#/settings/network		            | OK 	 |
| 36    |    /m/settings/email		         | /ng#/settings/email		              |	OK   |
| 37    |    /m/invite                     | /ng#/invite			                    |	OK   |
