'use strict';

var tournamentControllers = angular.module('tournamentControllers', []);

tournamentControllers.controller('TournamentListCtrl', ['$scope', 'Tournament', '$location', function($scope, Tournament, $location) {
    console.log('Tournament list controller');
    $scope.tournaments = Tournament.query();
    $scope.searchTournament = function(){
	console.log('TournamentListCtrl: searchTournament');
	console.log('keywords: ', $scope.keywords)
	$location.search('q', $scope.keywords).path('/tournaments/search');
    };
}]);


tournamentControllers.controller('TournamentSearchCtrl', ['$scope', '$routeParams', 'Tournament', '$location', function($scope, $routeParams, Tournament, $location) {
    console.log('Tournament search controller');
    console.log('routeParams: ', $routeParams);
    $scope.tournaments = Tournament.search( {q:$routeParams.q});
    $scope.searchTournament = function(){
	console.log('TournamentListCtrl: searchTournament');
	console.log('keywords: ', $scope.keywords)
	$location.search('q', $scope.keywords).path('/tournaments/search');
    };
}]);

tournamentControllers.controller('TournamentNewCtrl', ['$scope', 'Tournament', '$location', function($scope, Tournament, $location) {
    console.log('Tournament New controller');

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

tournamentControllers.controller('TournamentShowCtrl', ['$scope', '$routeParams', 'Tournament', '$location',function($scope, $routeParams, Tournament, $location) {
    console.log('Tournament Show controller');
    
    $scope.tournamentData =  Tournament.get({ id:$routeParams.id });
    console.log('tournamentData', $scope.tournamentData);
    $scope.deleteTournament = function() {
	Tournament.delete({ id:$routeParams.id },
			  function(){
			      $location.path('/');
			  },
			  function(err) {
			      console.log('delete failed: ', err.data);
			  });
    };
    $scope.joinTournament = function(){
	console.log('join tournament');
	console.log('routeParams: ', $routeParams);
	Tournament.join( {id:$routeParams.id});
    };
    $scope.addTeam = function(){
	console.log('add Team');
	//Tournament.addTeam with tournament.Id and team.Id
	//action="/ng#/tournamentteamrels/{{tournamentData.Tournament.Id}}" 
    };
    $scope.isTournamentAdmin = function(){
	console.log('is tournament admin ?');
	console.log('admin id: ', $scope.tournamentData.Tournament.AdminId);
	// Todo: waiting for #194: change 28001 by getCurrentUser().id service
	return $scope.tournamentData.Tournament.AdminId == 28001;
    }
    $scope.joined = function(){
	console.log('joined ?');
	// query tournamentrel for pair(tournament.id, user.id);
	// or add information to json of show.
	return false;
    }
}]);

tournamentControllers.controller('TournamentEditCtrl', ['$scope', '$routeParams', 'Tournament', '$location',function($scope, $routeParams, Tournament, $location) {
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
