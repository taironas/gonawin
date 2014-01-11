'use strict';

var teamControllers = angular.module('teamControllers', []);

teamControllers.controller('TeamListCtrl', ['$scope', 'Team', '$location',function($scope, Team, $location) {
    console.log('TeamListCtrl:');
    $scope.teams = Team.query();
    $scope.searchTeam = function(){
	console.log('TeamListCtrl: searchTeam');
	console.log('keywords: ', $scope.keywords);
	$location.search('q', $scope.keywords).path('/teams/search');
    };
}]);

teamControllers.controller('TeamSearchCtrl', ['$scope', '$routeParams', 'Team', '$location', function($scope, $routeParams, Team, $location) {
    console.log('search controller');
    console.log('routeParams: ', $routeParams);
    $scope.teams = Team.search( {q:$routeParams.q});
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
