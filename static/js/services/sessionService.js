'use strict'

var sessionService = angular.module('sessionService', []);

sessionService.factory('SessionService', function() {
	var currentUser = undefined;
	
	return {
		setCurrentUser: function(user) {
			currentUser = user;
		},
		getCurrentUser: function() {
			return currentUser;
		}
	};
});