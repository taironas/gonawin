'use strict';

purpleWingApp.controller('TournamentsEditController',
		     function TournamentsEditController($scope, tournamentData, $location, $routeParams){
			 console.log("tournament edit controller");
			 console.log($routeParams.teamId);
			 console.log('before getTournament');
			 $scope.tournamentData = teamData.getTournament($routeParams.tournamentId);
			 console.log('teamData: ',$scope.tournamentData);
		     }
		    );
