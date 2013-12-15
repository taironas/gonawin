'use strict';

var purpleWingApp = angular.module('purpleWingApp', ['ngSanitize', 'directive.g+signin', 'ngCookies'])
    .config(function($routeProvider){
	$routeProvider.when('/',
			    {
				templateUrl: 'templates/main.html', 
				controller: 'MainController'
			    });
	$routeProvider.when('/about',
			    {
				templateUrl: 'templates/about.html'
			    });
	$routeProvider.when('/contact',
			    {
				templateUrl: 'templates/contact.html'
			    });
	$routeProvider.when('/users/:userId',
			    {
				templateUrl: 'templates/user_show.html',
				controller: 'UserShowController'
			    });
	$routeProvider.when('/teams',
			    {
				templateUrl: 'templates/teams.html',
				controller: 'TeamsController'
			    });
	$routeProvider.when('/teams/new',
			    {
				templateUrl: 'templates/teams_new.html',
				controller: 'TeamsNewController'
			    });
	$routeProvider.when('/teams/:teamId',
			    {
				templateUrl: 'templates/teams_show.html',
				controller: 'TeamsShowController'
			    });
	$routeProvider.when('/teams/:teamId/edit',
			    {
				templateUrl: 'templates/teams_edit.html',
				controller: 'TeamsEditController'
			    });
	$routeProvider.when('/tournaments',
			    {
				templateUrl: 'templates/tournaments.html',
				controller: 'TournamentsController'
			    });
	$routeProvider.when('/tournaments/new',
			    {
				templateUrl: 'templates/tournaments_new.html',
				controller: 'TournamentsNewController'
			    });
	$routeProvider.when('/tournaments/:tournamentId',
			    {
				templateUrl: 'templates/tournaments_show.html',
				controller: 'TournamentsShowController'
			    });
	$routeProvider.when('/tournaments/:tournamentId/edit',
			    {
				templateUrl: 'templates/tournaments_edit.html',
				controller: 'TournamentsEditController'
			    });
	$routeProvider.when('/settings/edit-profile',
			    {
				templateUrl: 'templates/settings_edit-profile.html',
				controller: 'SettingsEditProfileController'
			    });
	$routeProvider.when('/settings/networks',
			    {
				templateUrl: 'templates/settings_networks.html'
			    });
	$routeProvider.when('/settings/email',
			    {
				templateUrl: 'templates/settings_email.html'
			    });

	$routeProvider.when('/invite',
			    {
				templateUrl: 'templates/invite.html',
				controller: 'InviteController'
			    });

	$routeProvider.otherwise( {redirectTo: '/'});
    })
    .factory('myCache', function($cacheFactory){
	return $cacheFactory('myCache', {capacity:3})
    });
