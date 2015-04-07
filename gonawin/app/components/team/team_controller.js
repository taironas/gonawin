'use strict';

// Team controllers manage team entities (creation, update, deletion) by getting
// data from REST service (resource).
// Handle also user subscription to a team (join/leave).
var teamControllers = angular.module('teamControllers', []);
// TeamListCtrl: fetch all teams data
teamControllers.controller('TeamListCtrl', ['$rootScope', '$scope', 'Team', 'User', '$location',
  function($rootScope, $scope, Team, User, $location) {
    console.log('Team list controller:');

    $rootScope.title = 'gonawin - Teams';

    $scope.countTeams = 25; // counter for the number of teams to display in view.
    $scope.pageTeams = 1;   // page counter for teams, to know which page to display next.

    $rootScope.currentUser.$promise.then(function(currentUser) {
      $scope.showMyTeams();
    });

    $scope.showMyTeams = function() {
      var userData = User.get({ id:$rootScope.currentUser.User.Id, including: "Teams", count:$scope.countTeams, page:$scope.pageTeams});
    	console.log('user data = ', userData);
    	userData.$promise.then(function(response) {
        $scope.teams = response.Teams;
        if(!$scope.teams || ($scope.teams && !$scope.teams.length)) {
    		    $scope.noTeamsMessage = 'You haven\'t joined a team yet';
  	    } else if($scope.teams !== undefined) {
    		    $scope.showMoreTeams = (response.Teams.length == $scope.countTeams);
  	    }
    	});
    };

    $scope.showOtherTeams = function() {
      Team.query({count:$scope.countTeams, page:$scope.pageTeams}).$promise.then(function(teams) {
        $scope.teams = teams;
        if(!$scope.teams || ($scope.teams && !$scope.teams.length)) {
      	    $scope.noTeamsMessage = 'There are no teams yet';
      	} else if($scope.teams !== undefined) {
      	    $scope.showMoreTeams = (teams.length == $scope.countTeams);
      	}
      });
    };

    // show more teams function:
    // retrieve teams by page and increment page.
    $scope.moreTeams = function() {
    	console.log('more teams');
    	$scope.pageTeams = $scope.pageTeams + 1;
    	Team.query({count:$scope.countTeams, page:$scope.pageTeams}).$promise.then(function(response) {
  	    console.log('response: ', response);
  	    $scope.teams = $scope.teams.concat(response);
  	    $scope.showMoreTeams = (response.length == $scope.countTeams);
    	});
    };

    $scope.moreJoinedTeams = function() {
    	console.log('more joined teams');
    	$scope.pageTeams = $scope.pageTeams + 1;
    	User.teams({id:$rootScope.currentUser.User.Id, count:$scope.countTeams, page:$scope.pageTeams}).$promise.then(function(response) {
  	    console.log('response: ', response);
  	    $scope.teams = $scope.teams.concat(response.Teams);
  	    $scope.showMoreTeams = (response.Teams.length == $scope.countTeams);
    	});
    };

    // Search function redirects to main search url /search.
    $scope.searchTeam = function(){
      console.log('TeamListCtrl: searchTeam');
      console.log('keywords: ', $scope.keywords);
      $location.search('q', $scope.keywords).path('/search');
     };
}]);

// TeamCardCtrl: handles team card
teamControllers.controller('TeamCardCtrl', ['$rootScope', '$scope', '$q', 'Team',
  function($rootScope, $scope, $q, Team) {
    console.log('Team Card controller: team = ', $scope.team);

    $scope.teamData = Team.get({ id:$scope.team.Id });

    // set isTeamAdmin boolean:
    // This variable defines if the current user is admin of the current team.
    $scope.teamData.$promise.then(function(teamResult) {
      console.log('team is admin ready');
      // as it depends of currentUser, make a promise
      var deferred = $q.defer();
      deferred.resolve((teamResult.Team.AdminIds.indexOf($scope.currentUser.User.Id)>=0));
      return deferred.promise;
    }).then(function(result){
      console.log('is team admin:', result);
      $scope.isTeamAdmin = result;
    });

    $scope.requestInvitation = function() {
      console.log('team request invitation');
      Team.requestInvite( {id:$scope.team.Id}, function() {
        Team.get({ id:$scope.team.Id }).$promise.then(function(teamDataResult) {
          $scope.teamData = teamDataResult;
        });
      }, function(err){
        $scope.messageDanger = err;
        console.log('invite failed ', err);
      });
    };

    // This function makes a user join a team.
    // It does so by caling Join on a Team.
    // This will update members data and join button name.
    $scope.joinTeam = function() {
      Team.join({ id:$scope.team.Id }).$promise.then(function(response) {
        console.log('joinTeam response = ', response);
        $scope.joinButtonName = 'Leave';
        $scope.joinButtonMethod = $scope.leaveTeam;
        $scope.messageInfo = response.MessageInfo;
        Team.get({ id:$scope.team.Id }).$promise.then(function(teamDataResult) {
          console.log('teamDataResult = ', teamDataResult);
          $scope.teamData = teamDataResult;
        });
        $rootScope.$broadcast('setUpdatedDashboard');
      });
    };
    // This function makes a user leave a team.
    // It does so by caling Leave on a Team.
    // This will update members data and leave button name.
    $scope.leaveTeam = function() {
      if(confirm('Are you sure?')) {
        Team.leave({ id:$scope.team.Id }).$promise.then(function(response) {
          console.log('leaveTeam response = ', response);
          $scope.joinButtonName = 'Join';
          $scope.joinButtonMethod = $scope.joinTeam;
          $scope.messageInfo = response.MessageInfo;
          Team.get({ id:$scope.team.Id }).$promise.then(function(teamDataResult) {
            console.log('teamDataResult = ', teamDataResult);
            $scope.teamData = teamDataResult;
          });
          $rootScope.$broadcast('setUpdatedDashboard');
        });
      }
    };

    $scope.teamData.$promise.then(function(teamResult) {
      var deferred = $q.defer();
    	if (teamResult.Joined) {
        deferred.resolve('Leave');
    	}
    	else {
        deferred.resolve('Join');
    	}
      return deferred.promise;
    }).then(function(result) {
      $scope.joinButtonName = result;
    });

    $scope.teamData.$promise.then(function(teamResult) {
      var deferred = $q.defer();
      if (teamResult.Joined) {
        deferred.resolve($scope.leaveTeam);
      }
      else {
        deferred.resolve($scope.joinTeam);
      }
      return deferred.promise;
    }).then(function(result) {
      $scope.joinButtonMethod = result;
    });
}]);

// TeamNewCtrl: use this controller to create a team.
teamControllers.controller('TeamNewCtrl', ['$rootScope', '$scope', 'Team', '$location', function($rootScope, $scope, Team, $location) {
  console.log('Team new controller:');

  $rootScope.title = 'gonawin - New Team';

  $scope.addTeam = function() {
    console.log('TeamNewCtrl: AddTeam');
    Team.save($scope.team,
	      function(response) {
		// set message information in root scope to retreive it in team show controller.
		// http://stackoverflow.com/questions/13740885/angularjs-location-scope
		$rootScope.messageInfo = response.MessageInfo;
		$location.path('/teams/' + response.Team.Id);
	      },
	      function(err) {
		$scope.messageDanger = err.data;
		console.log('save failed: ', err.data);
	      });
  };
}]);

// TeamShowCtrl: fetch data of specific team.
// // Handle also deletion of this same team and join/leave.
teamControllers.controller('TeamShowCtrl', ['$scope', '$routeParams', 'Team', '$location', '$q', '$rootScope', function($scope, $routeParams, Team, $location, $q, $rootScope) {
  console.log('Team show controller:');
  $scope.teamData = Team.get({ id:$routeParams.id });
  // get message info from redirects.
  $scope.messageInfo = $rootScope.messageInfo;
  // reset to nil var message info in root scope.
  $rootScope.messageInfo = undefined;

  $scope.teamData.$promise.then(function(response){
    $rootScope.title = 'gonawin - ' + response.Team.Name;
  });

  $scope.deleteTeam = function() {
    if(confirm('Are you sure?')){
        Team.delete({ id:$routeParams.id },
        function(response){
            $rootScope.messageInfo = response.MessageInfo;
            $location.path('/');
        },
        function(err) {
            $scope.messageDanger = err.data;
            console.log('delete failed: ', err.data);
        });
    }
  };

  // set admin candidates and array of functions.
  $scope.teamData.$promise.then(function(response){
    $scope.adminCandidates = response.Players;
    var len = 0;
    if(response.Players){
      len = response.Players.length;
    }
    $scope.addAdminButtonName = new Array(len);
    $scope.addAdminButtonMethod = new Array(len);

    for (var i=0 ; i<len; i++){
      // check if user is admin already here.
      if(response.Team.AdminIds.indexOf(response.Players[i].Id)>=0){
        $scope.addAdminButtonName[response.Players[i].Id] = 'Remove Admin';
        $scope.addAdminButtonMethod[response.Players[i].Id] = $scope.removeAdmin;
      } else{
        $scope.addAdminButtonName[response.Players[i].Id] = 'Add Admin';
        $scope.addAdminButtonMethod[response.Players[i].Id] = $scope.addAdmin;
      }
    }
  });

  // admin modal add buttons.
  // add admin state.
  $scope.addAdmin = function(userId){
    Team.addAdmin({id:$routeParams.id, userId:userId}).$promise.then(function(response){
      $scope.addAdminButtonName[userId] = 'Remove admin';
      $scope.addAdminButtonMethod[userId] = $scope.removeAdmin;
      $scope.messageInfo = response.MessageInfo;
    }, function(err) {
      $scope.messageDanger = err.data;
      console.log('save failed: ', err.data);
    });
  };
  // remove admin state.
  $scope.removeAdmin = function(userId){
    Team.removeAdmin({id:$routeParams.id, userId:userId}).$promise.then(function(response){
      $scope.addAdminButtonName[userId] = 'Add admin';
      $scope.addAdminButtonMethod[userId] = $scope.addAdmin;
      $scope.messageInfo = response.MessageInfo;
    }, function(err){
      console.log('save failed: ', err.data);
      $scope.messageDanger = err.data;
    });
  };

  // set isTeamAdmin boolean:
  // This variable defines if the current user is admin of the current team.
  $scope.teamData.$promise.then(function(teamResult){
    console.log('team is admin ready');
    // as it depends of currentUser, make a promise
    var deferred = $q.defer();
    deferred.resolve((teamResult.Team.AdminIds.indexOf($scope.currentUser.User.Id)>=0));
    return deferred.promise;
  }).then(function(result){
    console.log('is team admin:', result);
    $scope.isTeamAdmin = result;
  });

  // set tournament ids with "values" so that angular understands:
  // http://stackoverflow.com/questions/15488342/binding-inputs-to-an-array-of-primitives-using-ngrepeat-uneditable-inputs
  $scope.teamData.$promise.then(function(teamresp){
    var len = 0;
    if(teamresp.Team.TournamentIds){
      len = teamresp.Team.TournamentIds.length;
    }
    var tournamentIds = new Array();
    for(var i = 0; i < len; i++){
      tournamentIds.push({value: teamresp.Team.TournamentIds[i]});
    }
    $scope.teamData.Team.TournamentIds = tournamentIds;
    console.log('new tournament ids:', $scope.teamData.Team.TournamentIds);
  });

    $scope.requestInvitation = function(){
      console.log('team request invitation');
      Team.requestInvite( {id:$routeParams.id}, function(){
        Team.get({ id:$routeParams.id }).$promise.then(function(teamDataResult){
          $scope.teamData = teamDataResult;
        });
      }, function(err){
        $scope.messageDanger = err
        console.log('invite failed ', err);
      });
    };

    // This function makes a user join a team.
    // It does so by caling Join on a Team.
    // This will update members data and join button name.
    $scope.joinTeam = function(){
      Team.join({ id:$routeParams.id }).$promise.then(function(response){
        console.log('joinTeam response = ', response);
        $scope.joinButtonName = 'Leave';
        $scope.joinButtonMethod = $scope.leaveTeam;
        $scope.messageInfo = response.MessageInfo;
        Team.get({ id:$routeParams.id }).$promise.then(function(teamDataResult){
          console.log('teamDataResult = ', teamDataResult);
          $scope.teamData = teamDataResult;
        });
        $rootScope.$broadcast('setUpdatedDashboard');
      });
    };
    // This function makes a user leave a team.
    // It does so by caling Leave on a Team.
    // This will update members data and leave button name.
    $scope.leaveTeam = function(){
      if(confirm('Are you sure?')){
        Team.leave({ id:$routeParams.id }).$promise.then(function(response){
          console.log('leaveTeam response = ', response);
          $scope.joinButtonName = 'Join';
          $scope.joinButtonMethod = $scope.joinTeam;
          $scope.messageInfo = response.MessageInfo;
          Team.get({ id:$routeParams.id }).$promise.then(function(teamDataResult){
            console.log('teamDataResult = ', teamDataResult);
            $scope.teamData = teamDataResult;
          });
          $rootScope.$broadcast('setUpdatedDashboard');
        });
      }
    };

    $scope.teamData.$promise.then(function(teamResult){
	var deferred = $q.defer();
	if (teamResult.Joined) {
	    deferred.resolve('Leave');
	}
	else {
	    deferred.resolve('Join');
	}
	return deferred.promise;
    }).then(function(result){
	$scope.joinButtonName = result;
    });

    $scope.teamData.$promise.then(function(teamResult){
	var deferred = $q.defer();
	if (teamResult.Joined) {
	    deferred.resolve($scope.leaveTeam);
	}
	else {
	    deferred.resolve($scope.joinTeam);
	}
	return deferred.promise;
    }).then(function(result){
	$scope.joinButtonMethod = result;
    });

  // tab at undefined means you are in tournament/:id url.
  // So calendar should be active by default
  if($routeParams.tab == undefined){
    $scope.tab = 'members';
  } else {
    // Initialize tab variable to handle views:
    $scope.tab = $routeParams.tab;
  }

  $scope.tabs = {
    "members":      { title: 'Members',     url: 'components/team/tab_members.html' },
    "tournaments":  { title: 'Tournaments', url: 'components/team/tab_tournaments.html' },
    "ranking":      { title: 'Ranking',     url: 'components/team/tab_ranking.html' },
    "accuracies":   { title: 'Accuracies',  url: 'components/team/tab_accuracies.html' },
    "prizes":       { title: 'Prizes',      url: 'components/team/tab_prices.html' }
  };

  // set the current tab based on the 'tab' parameter
  if($scope.tab == undefined) {
    $scope.currentTab = $scope.tabs["members"].url;
  } else {
    $scope.currentTab = $scope.tabs[$scope.tab].url;
  }

  $scope.onClickTab = function (tab) {
    console.log('teamControllers: onClickTab, tab = ', tab);
    $scope.currentTab = tab.url;
  }

  $scope.rankingData = Team.ranking({id:$routeParams.id, rankby:$routeParams.rankby});
  // predicate is updated for ranking tables
  $scope.predicate = 'Score';
  $scope.accuracyData = Team.accuracies({id:$routeParams.id});

  $scope.pricesData = Team.prices({id:$routeParams.id});

  $scope.updatePrice = function(index) {
    var price = $scope.pricesData.Prices[index];
    console.log('update, prize = ', price);
    Team.updatePrice({id:price.TeamId, tournamentId:price.TournamentId}, price,
		function(response){
		  $rootScope.messageInfo = response.MessageInfo;
		  $location.path('/teams/' + price.TeamId);
		},
		function(err) {
		  $scope.messageDanger = err.data;
		  console.log('update failed: ', err.data);
		});
  }
}]);

// TeamEditCtrl: collects data to update an existing team.
teamControllers.controller('TeamEditCtrl', ['$rootScope', '$scope', '$routeParams', 'Team', '$location', function($rootScope, $scope, $routeParams, Team, $location) {
  console.log('Team edit controller:');
  $scope.teamData = Team.get({ id:$routeParams.id });
  $scope.visibility = 'public';
  console.log('$scope.teamData = ', $scope.teamData);
  $scope.teamData.$promise.then(function(response){
    $rootScope.title = 'gonawin - ' + response.Team.Name;
    if(response.Team.Private){
      $scope.visibility = 'private';
    }
  });

  $scope.updateTeam = function() {
    $scope.teamData.Team.Visibility = $scope.visibility;
    console.log('team data at update', $scope.teamData.Team);

    Team.update({ id:$routeParams.id }, $scope.teamData.Team,
		function(response){
		  $rootScope.messageInfo = response.MessageInfo;
		  $location.path('/teams/' + $routeParams.id);
		},
		function(err) {
		  $scope.messageDanger = err.data;
		  console.log('update failed: ', err.data);
		});
  }
}]);

// TeamInviteCtrl:
teamControllers.controller('TeamInviteCtrl', ['$rootScope', '$scope', '$routeParams', 'Team', 'User', '$location',
  function($rootScope, $scope, $routeParams, Team, User, $location) {
    console.log('Team invite controller:');
    $scope.teamData = Team.get({ id:$routeParams.id });

    $scope.teamData.$promise.then(function(response) {
      $rootScope.title = 'gonawin - ' + response.Team.Name;
    });

    $scope.inviteData = Team.invited({id:$routeParams.id });
    $scope.inviteData.$promise.then(function(response) {
    	console.log('invite data ', response);
      if(response.Users.length === 0 ) {
        $scope.noInvitationsMessage = 'No invitations sent.';
      }
      $scope.invitedUsers = response.Users;
    });

    if($routeParams.q !== undefined) {
      $scope.teamData.$promise.then(function(teamResponse) {
        // get teams result from search query
        $scope.keywords = $routeParams.q;
        $scope.usersData = User.search( {q:$routeParams.q});

        // users
        $scope.usersData.$promise.then(function(usersResult) {
          $scope.users = usersResult.Users;
          var len  = 0;
          if($scope.users) {
              len = $scope.users.length;
          }
          $scope.inviteData.$promise.then(function(inviteResponse) {
            for(var i = 0; i < len; i++) {
              $scope.users[i].invitationSent = false;
              // set invitation sent data.
              for( var j = 0; j < inviteResponse.Users.length;j++) {
                if($scope.users[i].Id == inviteResponse.Users[j].Id) {
                  $scope.users[i].invitationSent = true;
                }
              }
              var lenPlayers = 0;
              if(teamResponse.Players) {
                lenPlayers = teamResponse.Players.length;
              }
              for(var k = 0; k < lenPlayers; k++) {
                if($scope.users[i].Id == teamResponse.Players[k].Id) {
                  $scope.users[i].isMember = true;
                } else {
                  $scope.users[i].isMember = false;
                }
              }
            }
          });
          $scope.messageInfo = usersResult.MessageInfo;
          if(usersResult.Users === undefined) {
              $scope.noUsersMessage = 'No users found.';
          }
        });
      });
    }

    // Search function
    $scope.searchTeam = function(){
      console.log('TeamInviteCtrl: searchTeam');
      console.log('keywords: ', $scope.keywords);
      $location.search('q', $scope.keywords);
    };

    $scope.invite = function(userId, index){
	console.log('invite user id: ', userId);
	console.log('invite index: ', index);
	console.log('user: ', $scope.users[index]);
	$scope.users[index].invitationSent = true;

	console.log('sending invitation');
	Team.sendInvite({ id:$routeParams.id, userId: $scope.users[index].Id}).$promise.then(function(r){
	    console.log('invitation sent.');
	    $scope.invitedUsers.push($scope.users[index]);
	});
	$scope.noInvitationsMessage = '';
    };

}]);
