'use strict';

purpleWingApp.controller('TeamsEditController',
		     function TeamsEditController($scope, teamData, $location, $routeParams){
			 console.log("team edit controller");
			 console.log($routeParams.teamId);
			 console.log('before getTeam');
			 $scope.teamData = teamData.getTeam($routeParams.teamId);
			 console.log('teamData: ',$scope.teamData);
		     }
		    );
