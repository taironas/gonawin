'use strict';

// Search controllers manage the search of entities.
var searchControllers = angular.module('searchControllers', []);
// SearchCtrl: fetch all entities (teams, tournaments, users) on demand.
searchControllers.controller('SearchCtrl', ['$rootScope', '$scope', '$routeParams', 'Team', 'User', 'Tournament', '$location', function($rootScope, $scope, $routeParams, Team, User, Tournament, $location) {
    console.log('Search Controller:');
    console.log('Search Controller:', $routeParams.q);

    $scope.noTournamentsMessage = '';
    $scope.noTeamsMessage = '';
    $scope.noUsersMessage = '';

    if($routeParams.q != undefined){
	// get teams result from search query
	$scope.teamsData = Team.search( {q:$routeParams.q});
	$scope.tournamentsData = Tournament.search( {q:$routeParams.q});
	$scope.usersData = User.search( {q:$routeParams.q});
    }

    $scope.query = $routeParams.q;
    if($routeParams.q != undefined){
	// use the isSearching mode to differientiate:
	// no teams in app AND no teams found using query search
	$scope.isSearching = true;
	
	// get teams, tournaments and users via promises:
	// teams
	$scope.teamsData.$promise.then(function(result){
	    $scope.teams = result.Teams;
	    $scope.messageInfo = result.MessageInfo;
	    if(result.Teams == undefined){
		$scope.noTeamsMessage = 'No teams found.';
	    }
	});
	// tournaments
	$scope.tournamentsData.$promise.then(function(result){
	    $scope.tournaments = result.Tournaments;
	    $scope.messageInfo = result.MessageInfo;
	    if(result.Tournaments == undefined){
		$scope.noTournamentsMessage = 'No tournaments found.';
	    }
	});
	// users
	$scope.usersData.$promise.then(function(result){
	    $scope.users = result.Users;
	    $scope.messageInfo = result.MessageInfo;
	    if(result.Users == undefined){
		$scope.noUsersMessage = 'No users found.';
	    }
	});
    }
    // search function:
    $scope.search = function(){
	console.log('SearchCtrl: search');
	console.log('keywords: ', $scope.keywords);
	$location.search('q', $scope.keywords).path('/search');
    };

}]);
