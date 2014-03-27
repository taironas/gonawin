'use strict';

// Team controllers manage team entities (creation, update, deletion) by getting
// data from REST service (resource).
// Handle also user subscription to a team (join/leave).
var teamControllers = angular.module('teamControllers', []);
// TeamListCtrl: fetch all teams data 
teamControllers.controller('TeamListCtrl', ['$rootScope', '$scope', 'Team', 'User', '$location',
  function($rootScope, $scope, Team, User, $location) {
    console.log('Team list controller:');
    $scope.teams = Team.query();
    
    $scope.teams.$promise.then(function(result){
    if(!$scope.teams || ($scope.teams && !$scope.teams.length))
      $scope.noTeamsMessage = 'No team has been created';
    });

    $rootScope.currentUser.$promise.then(function(currentUser){
      var userData = User.get({ id:currentUser.User.Id, including: "Teams" });
      console.log('user data = ', userData);
      userData.$promise.then(function(result){
        $scope.joinedTeams = result.Teams;
        if(!$scope.joinedTeams || ($scope.joinedTeams && !$scope.joinedTeams.length))
          $scope.noJoinedTeamsMessage = 'You didn\'t join a team';
      });
    });

    $scope.searchTeam = function(){
      console.log('TeamListCtrl: searchTeam');
      console.log('keywords: ', $scope.keywords);
      $location.search('q', $scope.keywords).path('/teams/search');
    };
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
// TeamSearchCtrl: returns an array of teams based on a search query.
teamControllers.controller('TeamSearchCtrl', ['$scope', '$routeParams', 'Team', '$location', function($scope, $routeParams, Team, $location) {
  console.log('Team search controller:');
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
  // This function makes a user leave a team.
  // It does so by caling Leave on a Team.
  // This will update members data and leave button name.
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
// TeamEditCtrl: collects data to update an existing team.
teamControllers.controller('TeamEditCtrl', ['$rootScope', '$scope', '$routeParams', 'Team', '$location', function($rootScope, $scope, $routeParams, Team, $location) {
  console.log('Team edit controller:');
  $scope.teamData = Team.get({ id:$routeParams.id });

  $scope.updateTeam = function() {
    var teamData = Team.get({ id:$routeParams.id });
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
