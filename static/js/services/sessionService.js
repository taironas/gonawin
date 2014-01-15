'use strict'

var sessionService = angular.module('sessionService', []);

sessionService.factory('SessionService', ['$cookieStore', '$http', '$q', 'User', function($cookieStore, $http, $q, User) {
	var currentUser = undefined;
	var userIsLoggedIn = false;
	
	return {
		getCurrentUser: function() {
			var deferred = $q.defer();
			
			if(userIsLoggedIn && !currentUser)
			{
				currentUser = User.get({ id:$cookieStore.get('user_id') });
			}
			deferred.resolve(currentUser);
			
			return deferred.promise;
		},
		getUserLoggedIn: function() {
			var deferred = $q.defer();

			if(userIsLoggedIn) {
				deferred.resolve(true);
			}

			if(!$cookieStore.get('access_token') || !$cookieStore.get('user_id')) {
				userIsLoggedIn = false;
				deferred.resolve(false);
			} 
			else {
				User.get({ id:$cookieStore.get('user_id') }).$promise.then(function(result){
					currentUser = result;

					if(currentUser.Auth == $cookieStore.get('auth'))
					{
						userIsLoggedIn = true;
						deferred.resolve(true);
					} 
					else 
					{
						userIsLoggedIn = false;
						deferred.resolve(false);
					}
				});
			}
			
			return deferred.promise;
		},
		
		fetchUserInfo: function(access_token) {
			var service = this;
			if(access_token){
		    var peopleUrl = 'https://www.googleapis.com/plus/v1/people/me?access_token='+access_token;
				var promise = $http({ method: 'GET', url: peopleUrl, contentType: 'application/json' });
				return promise;
			}
			else {
		    return undefined;
			}
	  },
		
		fetchUser: function(accessToken, userInfo) {
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
					$cookieStore.put('user_id', currentUser.Id);
				}).
				error(function(result) { console.log('fetchUser failed') });
			
			return promise;
	  },
		
		logout: function() {
			var revokeUrl = 'https://accounts.google.com/o/oauth2/revoke?token='+$cookieStore.get('access_token');
			$.ajax({
				type: 'GET',
				url: revokeUrl,
				async: true,
				contentType: 'application/json',
				dataType: 'jsonp'
			});
			
			console.log('user has been logged out!');
			$cookieStore.remove('auth');
			$cookieStore.remove('access_token');
			$cookieStore.remove('user_id');

			currentUser = undefined;
			userIsLoggedIn = false;
		}
	};
}]);