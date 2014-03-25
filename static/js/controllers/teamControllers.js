'use strict';

var teamControllers = angular.module('teamControllers', []);

teamControllers.controller('TeamListCtrl', ['$rootScope', '$scope', 'Team', 'User', '$location',
  function($rootScope, $scope, Team, User, $location) {
    console.log('Team list controller:');
    $scope.teams = Team.query();

    $scope.teams.$promise.then(function(result){
      if(!$scope.teams || ($scope.teams && !$scope.teams.length))
        $scope.noTeamsMessage = 'You have no teams';
    });

    $rootScope.currentUser.$promise.then(function(currentUser){
      var userData = User.get({ id:currentUser.User.Id, including: "Teams" });
      console.log('user data = ', userData);
      userData.$promise.then(function(result){
        $scope.joinedTeams = result.Teams;
      });
    });

    $scope.searchTeam = function(){
      console.log('TeamListCtrl: searchTeam');
      console.log('keywords: ', $scope.keywords);
      $location.search('q', $scope.keywords).path('/teams/search');
    };

    $scope.showMoreTeams = false;
    $scope.moreTeams = function() {
      $scope.showMoreTeams = true;
    };

}]);

teamControllers.controller('TeamCardCtrl', ['$scope', 'Team',
  function($scope, Team) {
    console.log('Team card controller:');

    $scope.teamData = Team.get({ id:$scope.$parent.team.Id});

    $scope.teamData.$promise.then(function(teamData){
      $scope.team = teamData.Team;
      console.log('team card controller, teamData = ', teamData);
      $scope.membersCount = teamData.Players.length;
    });
}]);

teamControllers.controller('TeamSearchCtrl', ['$scope', '$routeParams', 'Team', '$location', function($scope, $routeParams, Team, $location) {
  console.log('Team search controller');
  console.log('routeParams: ', $routeParams);
  // get teams result from search query
  $scope.teamsData = Team.search( {q:$routeParams.q});

  $scope.query = $routeParams.q;
  // use the isSearching mode to differientiate:
  // no teams in app AND no teams found using query search
  $scope.isSearching = true;

  $scope.teamsData.$promise.then(function(result){
    $scope.teams = result.Teams;
    $scope.messageInfo = result.MessageInfo;
  });

  $scope.searchTeam = function(){
    console.log('TeamSearchCtrl: searchTeam');
    console.log('keywords: ', $scope.keywords);
    $location.search('q', $scope.keywords).path('/teams/search');
  };
}]);

// New Team controller. Use this controller to create a team.
teamControllers.controller('TeamNewCtrl', ['$scope', 'Team', '$location', '$rootScope', function($scope, Team, $location, $rootScope) {
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

teamControllers.controller('TeamShowCtrl', ['$scope', '$routeParams', 'Team', '$location', '$q', '$rootScope', function($scope, $routeParams, Team, $location, $q, $rootScope) {
  $scope.teamData = Team.get({ id:$routeParams.id });
  // get message info from redirects.
  $scope.messageInfo = $rootScope.messageInfo;
  $scope.deleteTeam = function() {
    Team.delete({ id:$routeParams.id },
		function(){
		  $location.path('/');
		},
		function(err) {
		  $scope.messageDanger = err.data;
		  console.log('delete failed: ', err.data);
		});
  };

  // set isTeamAdmin boolean:
  // This variable defines if the current user is admin of the current team.
  $scope.teamData.$promise.then(function(teamResult){
    console.log('team is admin ready');
    // as it depends of currentUser, make a promise
    var deferred = $q.defer();
    deferred.resolve((teamResult.Team.AdminId == $scope.currentUser.User.Id));
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
    Team.join({ id:$routeParams.id }).$promise.then(function(result){
      Team.members({ id:$routeParams.id }).$promise.then(function(membersResult){
        $scope.teamData.Members = membersResult.Members;
        $scope.joinButtonName = 'Leave';
        $scope.joinButtonMethod = $scope.leaveTeam;
      } );

    });
  };

  $scope.leaveTeam = function(){
    Team.leave({ id:$routeParams.id }).$promise.then(function(result){
      Team.members({ id:$routeParams.id }).$promise.then(function(membersResult){
        $scope.teamData.Members = membersResult.Members;
        $scope.joinButtonName = 'Join';
        $scope.joinButtonMethod = $scope.joinTeam;
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
}]);

teamControllers.controller('TeamEditCtrl', ['$scope', '$routeParams', 'Team', '$location', function($scope, $routeParams, Team, $location) {
  $scope.teamData = Team.get({ id:$routeParams.id });

  $scope.updateTeam = function() {
    var teamData = Team.get({ id:$routeParams.id });
    Team.update({ id:$routeParams.id }, $scope.teamData.Team,
		function(err){
		  $location.path('/teams/' + $routeParams.id);
		},
		function(err) {
		  $scope.messageDanger = err.data;
		  console.log('update failed: ', err.data);
		});
  }
}]);

// Ranking controller
teamControllers.controller('TeamRankingCtrl', ['$scope', '$routeParams', 'Team', '$location',function($scope, $routeParams, Team, $location) {
  console.log('Team ranking controller');
  console.log('route params', $routeParams);
  $scope.teamData = Team.get({ id:$routeParams.id });

  $scope.rankingData = Team.ranking({id:$routeParams.id, rankby:$routeParams.rankby});
  // predicate is udate for ranking tables
  $scope.predicate = 'Score';

}]);

//Team Accuracies controller
teamControllers.controller('TeamAccuraciesCtrl', ['$scope', '$routeParams', 'Team', '$location',function($scope, $routeParams, Team, $location) {
  console.log('Team accuracies controller');
  console.log('route params', $routeParams);
  $scope.teamData = Team.get({ id:$routeParams.id });

  $scope.accuracyData = Team.accuracies({id:$routeParams.id});
}]);

// team Accuracy in tournament
teamControllers.controller('TeamAccuracyByTournamentCtrl', ['$scope', '$routeParams', 'Team', '$location', function($scope, $routeParams, Team, $location){
  console.log('Team accuracy for team controller');
  console.log('route params', $routeParams);
  $scope.teamData = Team.get({ id:$routeParams.id });
  $scope.accuracyData = Team.accuracy({id:$routeParams.id, tournamentId:$routeParams.tournamentId});

}]);
