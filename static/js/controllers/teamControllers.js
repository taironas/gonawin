'use strict';

var teamControllers = angular.module('teamControllers', []);

teamControllers.controller('TeamListCtrl', ['$scope', 'Team', '$location', function($scope, Team, $location) {
    console.log('Team list controller:');
    $scope.teams = Team.query();
    $scope.searchTeam = function(){
	console.log('TeamListCtrl: searchTeam');
	console.log('keywords: ', $scope.keywords);
	$location.search('q', $scope.keywords).path('/teams/search');
    };
}]);

teamControllers.controller('TeamSearchCtrl', ['$scope', '$routeParams', 'Team', '$location', function($scope, $routeParams, Team, $location) {
    console.log('Team search controller');
    console.log('routeParams: ', $routeParams);
    $scope.teams = Team.search( {q:$routeParams.q});

    $scope.query = $routeParams.q;
    $scope.isSearching = true;

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
	$scope.team = Team.get({ id:$routeParams.id });
    
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
  $scope.team.$promise.then(function(teamResult){
		console.log('team is admin ready');
		// as it depends of currentUser, make a promise
		var deferred = $q.defer();
		$scope.currentUser.$promise.then(function(currentUserResult){
			console.log('is team admin: ', (teamResult.AdminId == currentUserResult.Id));
	  		deferred.resolve((teamResult.AdminId == currentUserResult.Id));
		});
		return deferred.promise;
	}).then(function(result){
		$scope.isTeamAdmin = result;
	});

    $scope.requestInvitation = function(){
	console.log('team request invitation');
	Team.invite( {id:$routeParams.id},
		     function(){
			 console.log('team invite successful');
		     },
		     function(err){
			 console.log('invite failed ', err);
		     });
    };

    // ToDo: remove function?
    $scope.joinTeam = function(){
	console.log('team join team');
    };
    // ToDo: remove function?
    $scope.leaveTeam = function(){
	console.log('team leave team');
    };
    
}]);

teamControllers.controller('TeamEditCtrl', ['$scope', '$routeParams', 'Team', '$location', function($scope, $routeParams, Team, $location) {
	$scope.team = Team.get({ id:$routeParams.id });
	
	$scope.updateTeam = function() {
		var team = Team.get({ id:$routeParams.id });
		Team.update({ id:$routeParams.id }, $scope.team,
			function(){
				$location.path('/teams/show/' + $routeParams.id);
			},
		function(err) {
			console.log('update failed: ', err.data);
		});
	}
}]);
