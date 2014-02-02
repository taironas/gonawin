'use strict'
var dataServices = angular.module('dataServices', ['ngResource']);

dataServices.factory('User', function($http, $resource, $cookieStore) {
	$http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');
    
	return $resource('j/users/:id', {id:'@id', including:'@including'}, {
		get: { method: 'GET', params: {including: '@including'}, url: 'j/users/show/:id' },
		update: { method: 'POST', url: 'j/users/update/:id' },
	})
});

dataServices.factory('Team', function($http, $resource, $cookieStore) {
	$http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');
	console.log("Team data service");
	return $resource('j/teams/:id', {id:'@id', q:'@q', requestId: '@requestId'}, {
		get: { method: 'GET', url: 'j/teams/show/:id' },
		save: { method: 'POST', url: 'j/teams/new' },
		update: { method: 'POST', url: 'j/teams/update/:id' },
		delete: { method: 'POST', url: 'j/teams/destroy/:id' },
		search: { method: 'GET', url: 'j/teams/search?q=:q', isArray: true},
		members: { method: 'GET', url:'j/teams/:id/members' },
		join: {method: 'POST', url: 'j/teamrels/create/:id'},
		leave: {method: 'POST', url: 'j/teamrels/destroy/:id'},
		invite: {method: 'POST', url: 'j/teams/invite/:id'},
		allowRequest : {method: 'POST', url: 'j/teams/allow/:requestId'},
		denyRequest : {method: 'POST', url: 'j/teams/deny/:requestId'}
	})
});

dataServices.factory('Tournament', function($http, $resource, $cookieStore) {
	$http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');
    
	return $resource('j/tournaments/:id', {id:'@id', q:'@q', teamId:'@teamId'}, {
		get: { method: 'GET', url: 'j/tournaments/show/:id' },
		save: { method: 'POST', url: 'j/tournaments/new' },
		update: { method: 'POST', url: 'j/tournaments/update/:id' },
		delete: { method: 'POST', url: 'j/tournaments/destroy/:id' },
		search: { method: 'GET', url: 'j/tournaments/search?q=:q', isArray: true},
		participants: { method: 'GET', url:'j/tournaments/:id/participants' },
		join: {method: 'POST', url: 'j/tournamentrels/create/:id'},
		leave: {method: 'POST', url: 'j/tournamentrels/destroy/:id'},
		joinAsTeam: {method: 'POST', url: 'j/tournamentteamrels/create/:id/:teamId'},
		leaveAsTeam: {method: 'POST', url: 'j/tournamentteamrels/destroy/:id/:teamId'},
		candidates: {method: 'GET', url: 'j/tournaments/candidates/:id', isArray: true}
	})
});

dataServices.factory('Invite', function($http, $cookieStore, $resource){
	$http.defaults.headers.common['Authorization'] = $cookieStore.get('auth');
	    
	return $resource('j/invite', {emails: '@emails'}, {
		send: {method: 'POST', params: {emails: '@emails'}, url: 'j/invite'}
	})
});

dataServices.factory('Session', function($http, $resource, User) {
  
  return $resource('/j/auth/', {access_token:'@access_token', id:'@id', name:'@name', email:'@email'}, {
    fetchUserInfo: { method:'GET', params: {access_token:'@access_token'}, url: 'https://www.googleapis.com/plus/v1/people/me' },
    fetchUser: { method:'GET', params: {access_token:'@access_token', id:'@id', name:'@name', email:'@email'}, url: '/j/auth/google' },
    logout: { method:'JSONP', params: {token:'@token', callback: 'JSON_CALLBACK'}, url: 'https://accounts.google.com/o/oauth2/revoke' }
  });
});
