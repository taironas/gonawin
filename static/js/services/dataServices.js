'use strict'
var dataServices = angular.module('dataServices', ['ngResource']);

dataServices.factory('Team', function($resource) {
	return $resource('j/teams/:id', {id:'@id'}, {
		get: { method: 'GET', url: 'j/teams/show/:id' },
		save: { method: 'POST', url: 'j/teams/new' },
		update: { method: 'POST', url: 'j/teams/update/:id' },
		delete: { method: 'POST', url: 'j/teams/destroy/:id' }
	})
});
