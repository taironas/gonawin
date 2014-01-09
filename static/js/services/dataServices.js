'use strict'
var dataServices = angular.module('dataServices', ['ngResource']);

dataServices.factory('User', function($http, $resource, $cookieStore) {
	$http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');
	
	return $resource('j/users/:id', {id:'@id'}, {
		get: { method: 'GET', url: 'j/users/show/:id' },
		update: { method: 'POST', url: 'j/users/update/:id' },
	})
});

dataServices.factory('Team', function($http, $resource, $cookieStore) {
	$http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');
	
	return $resource('j/teams/:id', {id:'@id'}, {
		get: { method: 'GET', url: 'j/teams/show/:id' },
		save: { method: 'POST', url: 'j/teams/new' },
		update: { method: 'POST', url: 'j/teams/update/:id' },
		delete: { method: 'POST', url: 'j/teams/destroy/:id' }
	})
});

dataServices.factory('Tournament', function($http, $resource, $cookieStore) {
	$http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');
	
	return $resource('j/tournaments/:id', {id:'@id'}, {
		get: { method: 'GET', url: 'j/tournaments/show/:id' },
		save: { method: 'POST', url: 'j/tournaments/new' },
		update: { method: 'POST', url: 'j/tournaments/update/:id' },
		delete: { method: 'POST', url: 'j/tournaments/destroy/:id' }
	})
});

dataServices.factory('Invite', function($http, $log, $q, $cookieStore){
	$http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');
	
	return {
		send: function(currentUser, emails){
			var deferred = $q.defer();
				$http({
					method: 'POST',
					url: '/j/invite',
					contentType: 'application/json',
					params:{ emails: emails, name: currentUser.Name } }).
					success(function(data,status,headers,config){
						deferred.resolve(data);
						$log.info(data, status, headers() ,config);
					}).
					error(function (data, status, headers, config){
						$log.warn(data, status, headers(), config);
						deferred.reject(status);
					});
				return deferred.promise;
		}
	};
});
