'use strict';

var sessionControllers = angular.module('sessionControllers', []);

sessionControllers.controller('SessionCtrl', ['$scope', '$location', 'SessionService',function ($scope, $location, SessionService) {
    console.log('SessionController module');
    $scope.currentUser  = undefined;
    $scope.loggedIn = false;
    
    $scope.initSession = function(){
	console.log('SessionController module:: initSession');
	
	SessionService.getUserLoggedIn().then(function(result) {
	    $scope.loggedIn = result;
	    
	    if($scope.loggedIn) {
		SessionService.getCurrentUser().then(function(result) {
		    $scope.currentUser = result;
		});
	    }
	});
    }
    
    $scope.$on('event:google-plus-signin-success', function (event, authResult) {
	// User successfully authorized the G+ App!
	console.log('SessionController module:: User successfully authorized the G+ App!');
	SessionService.fetchUserInfo(authResult.access_token).then(function(promise1) {
	    SessionService.fetchUser(authResult.access_token, promise1.data).then(function(promise2) {
		$scope.currentUser = promise2.data;
		$scope.loggedIn = true;
	    });
	});
    });
    $scope.$on('event:google-plus-signin-failure', function (event, authResult) {
	// User has not authorized the G+ App!
	console.log('Not signed into Google Plus.');
    });
    
    $scope.disconnect = function(){
	console.log('SessionController module:: disconnect');

	SessionService.logout();
				
	$scope.currentUser = undefined;
	$scope.loggedIn = false;
	
	$location.path('/');
    };
}]);
