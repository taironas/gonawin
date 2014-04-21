'use strict';

// Team controllers manage team entities (creation, update, deletion) by getting
// data from REST service (resource).
// Handle also user subscription to a team (join/leave).
var teamControllers = angular.module('teamControllers', []);
// TeamListCtrl: fetch all teams data
teamControllers.controller('TeamListCtrl', ['$rootScope', '$scope', 'Team', 'User', '$location', function($rootScope, $scope, Team, User, $location) {
    console.log('Team list controller:');

    $scope.countTeams = 25;            // counter for the number of teams to display in view.
    $scope.countJoinedTeams = 25;      // counter for the number of teams joined by user to display in view.

    $scope.pageTeams = 1;         // page counter for teams, to know which page to display next.
    $scope.pageJoinedTeams = 1;   // page counter for joined teams, to know which page to display next.

    // main query to /j/teams to get not joined teams.
    $scope.teams = Team.query({count:$scope.countTeams, page:$scope.pageTeams});

    // initilize team message and button visibility.
    $scope.teams.$promise.then(function(response){
	if(!$scope.teams || ($scope.teams && !$scope.teams.length)){
	    $scope.noTeamsMessage = 'No team has been created';
	}else if($scope.teams != undefined){
	    $scope.showMoreTeams = (response.length == $scope.countTeams);
	}
    });
    
    $rootScope.currentUser.$promise.then(function(currentUser){
	var userData = User.get({ id:currentUser.User.Id, including: "Teams", count:$scope.countJoinedTeams, page:$scope.pageJoinedTeams});
	console.log('user data = ', userData);
	userData.$promise.then(function(result){
            $scope.joinedTeams = result.Teams;
            if(!$scope.joinedTeams || ($scope.joinedTeams && !$scope.joinedTeams.length)){
		$scope.noJoinedTeamsMessage = 'You didn\'t join a team';
	    }else if($scope.joinedTeams != undefined){
		$scope.showMoreJoinedTeams = (result.Teams.length == $scope.countJoinedTeams);
	    }
	});
    });
    
    // Search function redirects to main search url /search.
    $scope.searchTeam = function(){
	console.log('TeamListCtrl: searchTeam');
	console.log('keywords: ', $scope.keywords);
	$location.search('q', $scope.keywords).path('/search');
    };

    // show more teams function:
    // retreive teams by page and increment page.
    $scope.moreTeams = function(){
	console.log('more teams');
	$scope.pageTeams = $scope.pageTeams + 1;
	Team.query({count:$scope.countTeams, page:$scope.pageTeams}).$promise.then(function(response){
	    console.log('response: ', response);
	    $scope.teams = $scope.teams.concat(response);
	    $scope.showMoreTeams = (response.length == $scope.countTeams);
	});
    };

    $scope.moreJoinedTeams = function(){
	console.log('more joined teams');
	$scope.pageJoinedTeams = $scope.pageJoinedTeams + 1;
	User.joinedTeams({id:$rootScope.currentUser.User.Id, count:$scope.countJoinedTeams, page:$scope.pageJoinedTeams}).$promise.then(function(response){
	    console.log('response: ', response);
	    $scope.joinedTeams = $scope.joinedTeams.concat(response.Teams);
	    console.log('jjajajjd', (response.Teams.length));
	    console.log('jjajajjd', ($scope.countJoinedTeams));
	    console.log('jjajajjd', (response.Teams.length == $scope.countJoinedTeams));
	    $scope.showMoreJoinedTeams = (response.Teams.length == $scope.countJoinedTeams);
	});

    }
}]);

// TeamCardCtrl: fetch data of a particular team.
teamControllers.controller('TeamCardCtrl', ['$scope', 'Team',
  function($scope, Team) {
    console.log('Team card controller:');
    console.log('team ID: ', $scope.$parent.team.Id);
    $scope.teamData = Team.get({ id:$scope.$parent.team.Id});

    $scope.teamData.$promise.then(function(teamData){
      $scope.team = teamData.Team;
      $scope.membersCount = teamData.Players.length;
    });
}]);

// TeamNewCtrl: use this controller to create a team.
teamControllers.controller('TeamNewCtrl', ['$rootScope', '$scope', 'Team', '$location', function($rootScope, $scope, Team, $location) {
  console.log('Team new controller:');
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
	    }else{
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
    	},function(err){
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
	var len  = 0
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

    // get prices for current team:
    $scope.pricesData = Team.prices({ id:$routeParams.id });

    $scope.requestInvitation = function(){
	console.log('team request invitation');
	Team.invite( {id:$routeParams.id}, function(){
	    console.log('team invite successful');
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
	    $scope.joinButtonName = 'Leave';
	    $scope.joinButtonMethod = $scope.leaveTeam;
	    $scope.messageInfo = response.MessageInfo;
	    Team.members({ id:$routeParams.id }).$promise.then(function(membersResult){
		$scope.teamData.Members = membersResult.Members;
	    } );

	});
    };
    // This function makes a user leave a team.
    // It does so by caling Leave on a Team.
    // This will update members data and leave button name.
    $scope.leaveTeam = function(){
	Team.leave({ id:$routeParams.id }).$promise.then(function(response){
	    $scope.joinButtonName = 'Join';
	    $scope.joinButtonMethod = $scope.joinTeam;
	    $scope.messageInfo = response.MessageInfo;
	    Team.members({ id:$routeParams.id }).$promise.then(function(membersResult){
		$scope.teamData.Members = membersResult.Members;
	    });
	});
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

    // Action triggered when 'Add Admin button' is clicked, modal window will be hidden.
    $scope.newTeam = function(){
	$('#team-modal').modal('hide');
    };

    // listen 'hidden.bs.modal' event to redirect to new team page
    // Only redirect if flag 'redirectToNewTeam' is set.
    $('#team-modal').on('hidden.bs.modal', function (e) {
	// need to have scope for $location to work. So add 'apply' function
	// inside js listener
	// if($scope.redirectToNewTeam == true){
	//     $scope.$apply(function(){
	// 	$location.path('/teams/new/');
	//     });
	// }
    })

    $scope.tabs = [{
      title: 'Members',
      url: 'templates/teams/players.html'
    }, {
      title: 'Ranking',
      url: 'templates/teams/partials/rankingData.html'
    }, {
      title: 'Accuracies',
      url: 'templates/teams/partials/accuraciesData.html'
    }, {
      title: 'Prices',
      url: 'templates/teams/partials/pricesData.html'
    }];

    $scope.currentTab = 'templates/teams/players.html';

    $scope.onClickTab = function (tab) {
        $scope.currentTab = tab.url;
    }

    $scope.isActiveTab = function(tabUrl) {
        return tabUrl == $scope.currentTab;
    }

    $scope.rankingData = Team.ranking({id:$routeParams.id, rankby:$routeParams.rankby});
    // predicate is udate for ranking tables
    $scope.predicate = 'Score';
    $scope.accuracyData = Team.accuracies({id:$routeParams.id});
    $scope.priceData = Team.prices({id:$routeParams.id});

}]);

// TeamEditCtrl: collects data to update an existing team.
teamControllers.controller('TeamEditCtrl', ['$rootScope', '$scope', '$routeParams', 'Team', '$location', function($rootScope, $scope, $routeParams, Team, $location) {
  console.log('Team edit controller:');
  $scope.teamData = Team.get({ id:$routeParams.id });
  $scope.visibility = 'public';
  console.log('$scope.teamData = ', $scope.teamData);
  $scope.teamData.$promise.then(function(response){
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

// TeamRankingCtrl: fetch ranking data of a specific team.
teamControllers.controller('TeamRankingCtrl', ['$scope', '$routeParams', 'Team', '$location',function($scope, $routeParams, Team, $location) {
  console.log('Team ranking controller:');
  console.log('route params', $routeParams);
  $scope.teamData = Team.get({ id:$routeParams.id });

  $scope.rankingData = Team.ranking({id:$routeParams.id, rankby:$routeParams.rankby});
  // predicate is udate for ranking tables
  $scope.predicate = 'Score';

}]);

// TeamAccuraciesCtrl: fetch accuracies data of a specific team.
teamControllers.controller('TeamAccuraciesCtrl', ['$scope', '$routeParams', 'Team', '$location',function($scope, $routeParams, Team, $location) {
  console.log('Team accuracies controller:');
  console.log('route params', $routeParams);
  $scope.teamData = Team.get({ id:$routeParams.id });

  $scope.accuracyData = Team.accuracies({id:$routeParams.id});
}]);

// TeamAccuracyByTournamentCtrl: fetch accuracy data by tournament of a specific team.
teamControllers.controller('TeamAccuracyByTournamentCtrl', ['$scope', '$routeParams', 'Team', '$location', function($scope, $routeParams, Team, $location){
  console.log('Team accuracy by tournament controller:');
  console.log('route params', $routeParams);
  $scope.teamData = Team.get({ id:$routeParams.id });
  $scope.accuracyData = Team.accuracy({id:$routeParams.id, tournamentId:$routeParams.tournamentId});
}]);


// TeamPricesCtrl: fetch prices data for a specific team.
teamControllers.controller('TeamPricesCtrl', ['$scope', '$routeParams', 'Team', '$location',function($scope, $routeParams, Team, $location) {
  console.log('Team prices controller:');
  console.log('route params', $routeParams);
  $scope.teamData = Team.get({ id:$routeParams.id });

  $scope.priceData = Team.prices({id:$routeParams.id});
}]);

// TeamPriceByTournamentCtrl: fetch accuracy data by tournament of a specific team.
teamControllers.controller('TeamPriceByTournamentCtrl', ['$scope', '$routeParams', 'Team', '$location', function($scope, $routeParams, Team, $location){
  console.log('Team price by tournament controller:');
  console.log('route params', $routeParams);
  $scope.teamData = Team.get({ id:$routeParams.id });
  $scope.priceData = Team.price({id:$routeParams.id, tournamentId:$routeParams.tournamentId});
}]);

// TeamPriceEditByTournamentCtrl: collects data to update an existing price.
teamControllers.controller('TeamPriceEditByTournamentCtrl', ['$rootScope', '$scope', '$routeParams', 'Team', '$location', function($rootScope, $scope, $routeParams, Team, $location) {
  console.log('Team price edit controller:');
  $scope.teamData = Team.get({ id:$routeParams.id });
  $scope.priceData = Team.price({id:$routeParams.id, tournamentId:$routeParams.tournamentId});

  $scope.updatePrice = function() {
    console.log('update, ', $scope.priceData)
    Team.updatePrice({id:$routeParams.id, tournamentId:$routeParams.tournamentId}, $scope.priceData.Price,
		function(response){
		  $rootScope.messageInfo = response.MessageInfo;
		  $location.path('/teams/' + $routeParams.id + '/prices/' + $routeParams.tournamentId);
		},
		function(err) {
		  $scope.messageDanger = err.data;
		  console.log('update failed: ', err.data);
		});
  }
}]);
