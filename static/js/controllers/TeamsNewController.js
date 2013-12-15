'use strict';

purpleWingApp.controller('TeamsNewController',
		     function TeamsNewController($scope, teamsNewData, $location, $routeParams){
			 console.log('teams new controller');
			 $scope.teamsNewData = teamsNewData.getData();
		     }
		    );
