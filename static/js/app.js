'use strict';

var purpleWingApp = angular.module('purpleWingApp', ['ngSanitize', 'ngRoute', 'directive.g+signin', 'ngCookies']);

purpleWingApp.config(['$routeProvider',
	function($routeProvider, $httpProvider) {
		$routeProvider.
			when('/', { templateUrl: 'templates/main.html', controller: 'MainController' }).
			when('/about', { templateUrl: 'templates/about.html' }).
			when('/contact', { templateUrl: 'templates/contact.html' }).
			when('/users/:userId', { templateUrl: 'templates/user_show.html', controller: 'UserShowController' }).
			when('/teams', { templateUrl: 'templates/teams.html', controller: 'TeamsController' }).
			when('/teams/new', { templateUrl: 'templates/teams_new.html', controller: 'TeamsNewController' }).
			when('/teams/:teamId', { templateUrl: 'templates/teams_show.html', controller: 'TeamsShowController' }).
			when('/teams/:teamId/edit', { templateUrl: 'templates/teams_edit.html', controller: 'TeamsEditController' }).
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
