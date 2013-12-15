'use strict';

purpleWingApp.controller('TournamentsController',
		     function TournamentsController($scope, tournamentsData, $location, $routeParams){
			 console.log('tournaments controller');
			 $scope.tournamentsData = tournamentsData.getData();
			 console.log('after');
		     }
		    );
