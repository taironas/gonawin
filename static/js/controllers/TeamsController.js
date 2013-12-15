'use strict';

purpleWingApp.controller('TeamsController',
		     function TeamsController($scope, teamsData, $location, $routeParams){
			 console.log('teams controller');
			 $scope.teamsData = teamsData.getData();
		     }
		    );
