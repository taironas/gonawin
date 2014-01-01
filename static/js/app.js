'use strict';

var purpleWingApp = angular.module('purpleWingApp', [
	'ngSanitize',
	'ngRoute',
	'ngResource',
	'ngCookies',
	'directive.g+signin',
	
	'sessionController',
	'userControllers',
	'teamControllers',
	'tournamentControllers',
	'inviteController',
	
	'sessionService',
	'dataServices'
]);

purpleWingApp.config(['$routeProvider',
	function($routeProvider) {
		$routeProvider.
			when('/', { templateUrl: 'templates/main.html', controller: 'MainController' }).
			when('/about', { templateUrl: 'templates/about.html' }).
			when('/contact', { templateUrl: 'templates/contact.html' }).
			when('/users/', { templateUrl: 'templates/users/index.html', controller: 'UserListCtrl' }).
			when('/users/show/:id', { templateUrl: 'templates/users/show.html', controller: 'UserShowCtrl' }).
			when('/teams', { templateUrl: 'templates/teams/index.html', controller: 'TeamListCtrl' }).
			when('/teams/new', { templateUrl: 'templates/teams/new.html', controller: 'TeamNewCtrl' }).
			when('/teams/show/:id', { templateUrl: 'templates/teams/show.html', controller: 'TeamShowCtrl' }).
			when('/teams/edit/:id', { templateUrl: 'templates/teams/edit.html', controller: 'TeamEditCtrl' }).
			when('/tournaments', { templateUrl: 'templates/tournaments/index.html', controller: 'TournamentListCtrl' }).
			when('/tournaments/new', { templateUrl: 'templates/tournaments/new.html', controller: 'TournamentNewCtrl' }).
			when('/tournaments/show/:id', { templateUrl: 'templates/tournaments/show.html', controller: 'TournamentShowCtrl' }).
			when('/tournaments/edit/:id', { templateUrl: 'templates/tournaments/edit.html', controller: 'TournamentEditCtrl' }).
			when('/settings/edit-profile', { templateUrl: 'templates/users/edit.html', controller: 'UserEditCtrl' }).
			when('/settings/networks', { templateUrl: 'templates/settings_networks.html' }).
			when('/settings/email', { templateUrl: 'templates/settings_email.html' }).
			when('/invite', { templateUrl: 'templates/invite.html', controller: 'InviteCtrl' }).
			otherwise( {redirectTo: '/'});
}]);
