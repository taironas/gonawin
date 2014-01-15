'use strict';

var purpleWingApp = angular.module('purpleWingApp', [
    'ngSanitize',
    'ngRoute',
    'ngResource',
    'ngCookies',
    'directive.g+signin',
    'directive.formValidation',
    
    'mainControllers',
    'sessionControllers',
    'userControllers',
    'teamControllers',
    'tournamentControllers',
    'inviteControllers',
    
    'sessionService',
    'dataServices'
]);

purpleWingApp.config(['$routeProvider',
	function($routeProvider) {
		$routeProvider.
			when('/', { templateUrl: 'templates/main.html', controller: 'MainCtrl', requireLogin: false }).
			when('/about', { templateUrl: 'templates/about.html', requireLogin: false }).
			when('/contact', { templateUrl: 'templates/contact.html', requireLogin: false }).
			when('/users/', { templateUrl: 'templates/users/index.html', controller: 'UserListCtrl', requireLogin: true }).
			when('/users/show/:id', { templateUrl: 'templates/users/show.html', controller: 'UserShowCtrl', requireLogin: true }).
			when('/teams', { templateUrl: 'templates/teams/index.html', controller: 'TeamListCtrl', requireLogin: true }).
			when('/teams/new', { templateUrl: 'templates/teams/new.html', controller: 'TeamNewCtrl', requireLogin: true }).
			when('/teams/show/:id', { templateUrl: 'templates/teams/show.html', controller: 'TeamShowCtrl', requireLogin: true }).
			when('/teams/edit/:id', { templateUrl: 'templates/teams/edit.html', controller: 'TeamEditCtrl', requireLogin: true }).
			when('/teams/search', { templateUrl: 'templates/teams/index.html', controller: 'TeamSearchCtrl', requireLogin: true}).
			when('/tournaments', { templateUrl: 'templates/tournaments/index.html', controller: 'TournamentListCtrl', requireLogin: true }).
			when('/tournaments/new', { templateUrl: 'templates/tournaments/new.html', controller: 'TournamentNewCtrl', requireLogin: true }).
			when('/tournaments/show/:id', { templateUrl: 'templates/tournaments/show.html', controller: 'TournamentShowCtrl', requireLogin: true }).
			when('/tournaments/edit/:id', { templateUrl: 'templates/tournaments/edit.html', controller: 'TournamentEditCtrl', requireLogin: true }).
			when('/tournaments/edit/:id', { templateUrl: 'templates/tournaments/edit.html', controller: 'TournamentEditCtrl', requireLogin: true }).
			when('/tournaments/search', { templateUrl: 'templates/tournaments/index.html', controller: 'TournamentSearchCtrl', requireLogin: true }).
			when('/settings/edit-profile', { templateUrl: 'templates/users/edit.html', controller: 'UserEditCtrl', requireLogin: true }).
			when('/settings/networks', { templateUrl: 'templates/settings/networks.html', requireLogin: true }).
			when('/settings/email', { templateUrl: 'templates/settings/email.html', requireLogin: true }).
			when('/invite', { templateUrl: 'templates/invite.html', controller: 'InviteCtrl', requireLogin: true }).
			otherwise( {redirectTo: '/'});
}]);

purpleWingApp.run(['$rootScope', 'SessionService', '$location',
	function($rootScope, SessionService, $location){
		$rootScope.$on("$routeChangeStart", function(event, next, current) {
		// Everytime the route in our app changes check authentication status
		if (next.requireLogin) {
			SessionService.getUserLoggedIn().then(function(result) {
			}, function(reason) {
				// if you're logged out send to home page.
				$location.path('/');
				event.preventDefault();
			});
		}
	});
}]);
