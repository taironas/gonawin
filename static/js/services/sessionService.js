'use strict'

var sessionService = angular.module('sessionService', []);

sessionService.factory('SessionService', function() {
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
		}
	};
});