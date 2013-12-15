'use strict';

purpleWingApp.controller('TeamsShowController',
		     function TeamsShowController($scope, teamData, $location, $routeParams){
			 console.log("team show controller");
			 console.log($routeParams.teamId);
			 console.log('before getTeam');
			 $scope.teamData = teamData.getTeam($routeParams.teamId);
		     }
		    );
