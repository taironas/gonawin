'use strict';

purpleWingApp.controller('TournamentsNewController',
		     function TournamentsNewController($scope, tournamentsNewData, $location, $routeParams){
			 console.log('tournaments new controller');
			 $scope.tournamentsNewData = tournamentsNewData.getData();
		     }
		    );
