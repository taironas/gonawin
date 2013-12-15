'use strict';

purpleWingApp.controller('UserShowController',
		     function UserShowController($scope, userData, $location, $routeParams){
			 console.log("user show controller");
			 console.log($routeParams.userId);
			 console.log('before getUser');
			 $scope.userData = userData.getUser($routeParams.userId);
			 console.log($scope.userData);
			 console.log('after getUser');
			 
		     }
		    );
