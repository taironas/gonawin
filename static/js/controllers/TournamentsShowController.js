'use strict';

purpleWingApp.controller('TournamentsShowController',
		     function TournamentsShowController($scope, tournamentData, $location, $routeParams){
			 console.log("team show controller");
			 console.log($routeParams.tournamentId);
			 console.log('before getTournament');
			 $scope.tournamentData = tournamentData.getTournament($routeParams.tournamentId);
		     }
		    );
