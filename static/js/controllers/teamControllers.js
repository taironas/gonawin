'use strict';

var teamControllers = angular.module('teamControllers', []);

teamControllers.controller('TeamListCtrl', ['$scope', 'Team', function($scope, Team) {
	$scope.teams = Team.query();
}]);

teamControllers.controller('TeamNewCtrl', ['$scope', 'Team', '$location', function($scope, Team, $location) {
	$scope.addTeam = function() {
		Team.save($scope.team,
			function(team) {
				$location.path('/teams/show/' + team.Id);
			},
			function(err) {
				console.log('save failed: ', err.data);
			});
	};
}]);

teamControllers.controller('TeamShowCtrl', ['$scope', '$routeParams', 'Team', '$location', function($scope, $routeParams, Team, $location) {
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
