'use strict';

var tournamentControllers = angular.module('tournamentControllers', []);

tournamentControllers.controller('TournamentListCtrl', ['$scope', 'Tournament', function($scope, Tournament) {
	$scope.tournaments = Tournament.query();
}]);

tournamentControllers.controller('TournamentNewCtrl', ['$scope', 'Tournament', '$location', function($scope, Tournament, $location) {
	$scope.addTournament = function() {
		Tournament.save($scope.tournament,
			function(tournament) {
				$location.path('/tournaments/show/' + tournament.Id);
			},
			function(err) {
				console.log('save failed: ', err.data);
			});
	};
}]);

tournamentControllers.controller('TournamentShowCtrl', ['$scope', '$routeParams', 'Tournament', '$location',
	function($scope, $routeParams, Tournament, $location) {
		$scope.tournament = Tournament.get({ id:$routeParams.id });
	
		$scope.deleteTournament = function() {
			Tournament.delete({ id:$routeParams.id },
				function(){
					$location.path('/');
				},
				function(err) {
					console.log('delete failed: ', err.data);
				});
		};
}]);

tournamentControllers.controller('TournamentEditCtrl', ['$scope', '$routeParams', 'Tournament', '$location',
	function($scope, $routeParams, Tournament, $location) {
		$scope.tournament = Tournament.get({ id:$routeParams.id });
	
		$scope.updateTournament = function() {
			var tournament = Tournament.get({ id:$routeParams.id });
			Tournament.update({ id:$routeParams.id }, $scope.tournament,
				function(){
					$location.path('/tournaments/show/' + $routeParams.id);
				},
			function(err) {
				console.log('update failed: ', err.data);
			});
		}
}]);
