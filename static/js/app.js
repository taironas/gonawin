'use strict';

var purpleWingApp = angular.module('purpleWingApp', [
	'ngSanitize',
	'ngRoute',
	'ngResource',
	'ngCookies',
	'directive.g+signin',
	
	'teamControllers',
	'teamServices'
]);

purpleWingApp.config(['$routeProvider',
	function($routeProvider) {
		$routeProvider.
			when('/', { templateUrl: 'templates/main.html', controller: 'MainController' }).
			when('/about', { templateUrl: 'templates/about.html' }).
			when('/contact', { templateUrl: 'templates/contact.html' }).
			when('/users/:userId', { templateUrl: 'templates/user_show.html', controller: 'UserShowController' }).
			when('/teams', { templateUrl: 'templates/teams/index.html', controller: 'TeamListCtrl' }).
			when('/teams/new', { templateUrl: 'templates/teams/new.html', controller: 'TeamNewCtrl' }).
			when('/teams/show/:id', { templateUrl: 'templates/teams/show.html', controller: 'TeamShowCtrl' }).
			when('/teams/edit/:id', { templateUrl: 'templates/teams/edit.html', controller: 'TeamEditCtrl' }).
			when('/tournaments', { templateUrl: 'templates/tournaments.html', controller: 'TournamentsController' }).
			when('/tournaments/new', { templateUrl: 'templates/tournaments_new.html', controller: 'TournamentsNewController' }).
			when('/tournaments/:tournamentId', { templateUrl: 'templates/tournaments_show.html', controller: 'TournamentsShowController' }).
			when('/tournaments/:tournamentId/edit', { templateUrl: 'templates/tournaments_edit.html', controller: 'TournamentsEditController' }).
			when('/settings/edit-profile', { templateUrl: 'templates/settings_edit-profile.html', controller: 'SettingsEditProfileController' }).
			when('/settings/networks', { templateUrl: 'templates/settings_networks.html' }).
			when('/settings/email', { templateUrl: 'templates/settings_email.html' }).
			when('/invite', { templateUrl: 'templates/invite.html', controller: 'InviteController' }).
			otherwise( {redirectTo: '/'});
}]);
