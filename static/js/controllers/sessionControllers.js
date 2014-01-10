'use strict';

var sessionControllers = angular.module('sessionControllers', []);

sessionControllers.controller('SessionCtrl', ['$scope', '$http', '$cookieStore', '$location', 'SessionService',
	function ($scope, $http, $cookieStore, $location, SessionService) {
		console.log('SessionController module');
		$scope.currentUser = undefined;
		$scope.loggedIn = false;
		
		$scope.$on('event:google-plus-signin-success', function (event, authResult) {
			// User successfully authorized the G+ App!
			console.log('Signed in!');
			$scope.checkAuth(authResult)
		});
		$scope.$on('event:google-plus-signin-failure', function (event, authResult) {
			// User has not authorized the G+ App!
			console.log('Not signed into Google Plus.');
		});
	  $scope.checkAuth = function(authResult){
			if(authResult.access_token){
		    $scope.login(authResult.access_token);
		    return true;
			}
			else {
		    return false;
			}
	  };
	  $scope.login = function(accessToken){
			var peopleUrl = 'https://www.googleapis.com/plus/v1/people/me?access_token='+accessToken;
			$http({ method: 'GET', url: peopleUrl, contentType: 'application/json' }).
		    success(function(result) { $scope.completeLogin(accessToken, result) });
	  };
	  $scope.completeLogin = function(accessToken, userInfo){
			var authUrl = '/j/auth/google';
			$http({
				method: 'GET',
				url: authUrl,
				contentType: 'application/json',
				params:{ access_token: accessToken, id: userInfo.id, name: userInfo.displayName, email: userInfo.emails[0].value} }).
		    success(function(data, status, headers, config) {
					console.log('completeLogin successfully');
					SessionService.setCurrentUser(data);
					$scope.currentUser = data;
					SessionService.setUserLoggedIn(true);
					$scope.loggedIn = true;
					$cookieStore.put('access_token', accessToken);
					$cookieStore.put('auth', $scope.currentUser.Auth);
				}).
				error(function(result) { console.log('completeLogin failed') });
	  };
	
	  $scope.disconnect = function(){
			var revokeUrl = 'https://accounts.google.com/o/oauth2/revoke?token='+$cookieStore.get('access_token');
			$.ajax({
				type: 'GET',
				url: revokeUrl,
				async: false,
				contentType: 'application/json',
				dataType: 'jsonp',
		    success: function(result){
					console.log('disconnected!');
					$cookieStore.remove('auth');
					$cookieStore.remove('access_token');
					SessionService.setCurrentUser(undefined);
					$scope.currentUser = undefined;
					SessionService.setUserLoggedIn(false);
					$scope.loggedIn = false;
					$scope.$apply();
					// redirect to home page
					$location.path('/');
				}
			});
		};
}]);
