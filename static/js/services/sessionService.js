'use strict'

var sessionService = angular.module('sessionService', []);

sessionService.factory('SessionService', ['$cookieStore', '$http', function($cookieStore, $http) {
	var currentUser = undefined;
	var userIsLoggedIn = false;
	
	return {
		setCurrentUser: function(user) {
			currentUser = user;
		},
		getCurrentUser: function() {
			return currentUser;
		},
		setUserLoggedIn: function(value) {
			userIsLoggedIn = value;
		},
		getUserLoggedIn: function() {
			return userIsLoggedIn;
		},
		
		login: function(authResult) {
			var service = this;
			if(authResult.access_token){
		    var peopleUrl = 'https://www.googleapis.com/plus/v1/people/me?access_token='+authResult.access_token;
				var promise = $http({ method: 'GET', url: peopleUrl, contentType: 'application/json' });
				return promise;
			}
			else {
		    return undefined;
			}
	  },
		
		completeLogin: function(accessToken, userInfo) {
			var authUrl = '/j/auth/google';
			var promise = $http({
				method: 'GET',
				url: authUrl,
				contentType: 'application/json',
				params:{ access_token: accessToken, id: userInfo.id, name: userInfo.displayName, email: userInfo.emails[0].value} }).
		    success(function(data, status, headers, config) {					
					currentUser = data;
					userIsLoggedIn = true;
					
					$cookieStore.put('access_token', accessToken);
					$cookieStore.put('auth', currentUser.Auth);
				}).
				error(function(result) { console.log('completeLogin failed') });
			
			return promise;
	  },
		
		logout: function() {
			var revokeUrl = 'https://accounts.google.com/o/oauth2/revoke?token='+$cookieStore.get('access_token');
			$.ajax({
				type: 'GET',
				url: revokeUrl,
				async: false,
				contentType: 'application/json',
				dataType: 'jsonp',
		    success: function(result){
					console.log('user has been logged out!');
					$cookieStore.remove('auth');
					$cookieStore.remove('access_token');

					currentUser = undefined;
					userIsLoggedIn = false;
				}
			});
		}
	};
}]);