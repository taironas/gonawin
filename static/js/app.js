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
