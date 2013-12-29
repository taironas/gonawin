'use strict';

var teamControllers = angular.module('teamControllers', []);

teamControllers.controller('TeamListCtrl', ['$scope', 'Team', function($scope, Team) {
	$scope.teams = Team.query();
}]);

teamControllers.controller('TeamNewCtrl', ['$scope', 'Team', function($scope, Team) {
	$scope.addTeam = function() {
		Team.save($scope.team);
	}
}]);

teamControllers.controller('TeamShowCtrl', ['$scope', '$routeParams', 'Team', function($scope, $routeParams, Team) {
	$scope.team = Team.get({ id:$routeParams.id });
	
	$scope.deleteTeam = function() {
		Team.delete({ id:$routeParams.id });
	};
}]);

teamControllers.controller('TeamEditCtrl', ['$scope', '$routeParams', 'Team', function($scope, $routeParams, Team) {
	$scope.team = Team.get({ id:$routeParams.id });
	
	$scope.updateTeam = function() {
		var team = Team.get({ id:$routeParams.id });
		
		Team.update({ id:$routeParams.id }, $scope.team);
	}
}]);
