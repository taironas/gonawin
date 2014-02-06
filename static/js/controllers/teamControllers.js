'use strict';

var teamControllers = angular.module('teamControllers', []);

teamControllers.controller('TeamListCtrl', ['$scope', 'Team', '$location', function($scope, Team, $location) {
  console.log('Team list controller:');
  $scope.teams = Team.query();

  $scope.teams.$promise.then(function(result){
    if(!$scope.teams || ($scope.teams && !$scope.teams.length))
      $scope.noTeamsMessage = 'You have no teams';
  });

  $scope.searchTeam = function(){
    console.log('TeamListCtrl: searchTeam');
    console.log('keywords: ', $scope.keywords);
    $location.search('q', $scope.keywords).path('/teams/search');
  };
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
  });
   
  $scope.searchTeam = function(){
    console.log('TeamSearchCtrl: searchTeam');
    console.log('keywords: ', $scope.keywords);
    $location.search('q', $scope.keywords).path('/teams/search');
  };
}]);

teamControllers.controller('TeamNewCtrl', ['$scope', 'Team', '$location', function($scope, Team, $location) {
  $scope.addTeam = function() {
    console.log('TeamNewCtrl: AddTeam');
    Team.save($scope.team,
	      function(team) {
		$location.path('/teams/show/' + team.Id);
	      },
	      function(err) {
		console.log('save failed: ', err.data);
	      });
  };
}]);

teamControllers.controller('TeamShowCtrl', ['$scope', '$routeParams', 'Team', '$location', '$q', function($scope, $routeParams, Team, $location, $q) {
  $scope.teamData = Team.get({ id:$routeParams.id });
  
  $scope.deleteTeam = function() {
    Team.delete({ id:$routeParams.id },
		function(){
		  $location.path('/');
		},
		function(err) {
		  console.log('delete failed: ', err.data);
		});
  };

  // set isTeamAdmin boolean
  $scope.teamData.$promise.then(function(teamResult){
    console.log('team is admin ready');
    // as it depends of currentUser, make a promise
    var deferred = $q.defer();
    deferred.resolve((teamResult.Team.AdminId == $scope.currentUser.Id));
    return deferred.promise;
  }).then(function(result){
    $scope.isTeamAdmin = result;
  });

  $scope.requestInvitation = function(){
    console.log('team request invitation');
    Team.invite( {id:$routeParams.id}, function(){
      console.log('team invite successful');
    }, function(err){
      console.log('invite failed ', err);
    });
  };

  $scope.joinTeam = function(){
    Team.join({ id:$routeParams.id }).$promise.then(function(result){
      Team.members({ id:$routeParams.id }).$promise.then(function(membersResult){
	$scope.teamData.Players = membersResult.Members;
	$scope.joinButtonName = 'Leave';
	$scope.joinButtonMethod = $scope.leaveTeam;
      } );
      
    });
  };

  $scope.leaveTeam = function(){
    Team.leave({ id:$routeParams.id }).$promise.then(function(result){
      Team.members({ id:$routeParams.id }).$promise.then(function(membersResult){
	$scope.teamData.Players = membersResult.Members;
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
		function(){
		  $location.path('/teams/show/' + $routeParams.id);
		},
		function(err) {
		  console.log('update failed: ', err.data);
		});
  }
}]);
